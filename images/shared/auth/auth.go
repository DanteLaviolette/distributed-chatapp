package auth

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"

	"shared/constants"
	"shared/persistence"
	"shared/structs"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const SECRET_LENGTH = 32
const AUTH_COOKIE_NAME = "auth"
const REFRESH_COOKIE_NAME = "refresh"

// Local keys
const LOCALS_USER_ID = "userId"
const LOCALS_USER_NAME = "userName"
const LOCALS_USER_EMAIL = "userEmail"
const LOCALS_REFRESH_TOKEN_ID = "refreshTokenId"

type AuthProvider struct {
	/*
	   Given a user & key returns a JWT auth token signed by the key as a fiber cookie.
	   Returns nil if an error occurs.
	   JWT Cookie contains the fields id, email, name, iat & exp.
	*/
	CreateAuthCookie func(structs.UserWithId) *fiber.Cookie
	/*
	   Given a user & key returns a JWT refresh token signed by the key as a fiber cookie.
	   Returns nil if an error occurs.
	   Note: Also writes a hashed secret to the DB to validate the refresh token later.
	   JWT Cookie contains the fields userId, secret, secretId, iat & exp.
	*/
	CreateRefreshCookie func(structs.UserWithId) *fiber.Cookie
	/*
		To be used as fiber middleware. Proceeds if user is logged in (potentially
		refreshing their credentials). Fails the request with 401 error if not logged
		in.
		Upon success, adds userId, userName, userEmail & refreshTokenId to c.Locals
	*/
	IsAuthenticatedFiberMiddleware func(*fiber.Ctx) error
}

func Initialize(authPrivateKey string, refreshPrivateKey string) *AuthProvider {
	return &AuthProvider{
		CreateAuthCookie: func(user structs.UserWithId) *fiber.Cookie {
			return createAuthCookieFromUser(user, authPrivateKey)
		},
		CreateRefreshCookie: func(user structs.UserWithId) *fiber.Cookie {
			return createRefreshCookie(user.ID.Hex(), refreshPrivateKey)
		},
		IsAuthenticatedFiberMiddleware: func(c *fiber.Ctx) error {
			return isAuthenticatedFiberMiddleware(c, authPrivateKey, refreshPrivateKey)
		},
	}
}

/*
To be used as fiber middleware. Proceeds if user is logged in (potentially
refreshing their credentials). Fails the request with 401 error if not logged
in.
Upon success, adds userId, userName, userEmail & refreshTokenId to c.Locals
*/
func isAuthenticatedFiberMiddleware(c *fiber.Ctx, authPrivateKey string,
	refreshPrivateKey string) error {
	// Get cookies
	auth := c.Cookies(AUTH_COOKIE_NAME)
	refresh := c.Cookies(REFRESH_COOKIE_NAME)
	// Fail if cookies weren't found
	if auth == "" || refresh == "" {
		return c.SendStatus(failUnauthenticatedRequest(c))
	}

	// Validate & parse tokens
	refreshToken, err := parseJWT(refresh, &structs.RefreshTokenClaim{}, refreshPrivateKey)
	if err != nil || refreshToken == nil {
		// Invalid token case
		return c.SendStatus(http.StatusUnauthorized)
	}
	authToken, err := parseJWT(auth, &structs.AuthTokenClaim{}, authPrivateKey)

	isAuthenticated := false // Assume token is invalid
	if errors.Is(err, jwt.ErrTokenExpired) {
		// Expired auth token -- attempt to refresh token
		isAuthenticated = refreshCookies(c, refreshToken, authToken,
			authPrivateKey, refreshPrivateKey)
	} else if err == nil && authToken != nil {
		// Success case where auth token is valid and not expired
		isAuthenticated = true
	}
	// Convert claims to specific jwt claims
	authTokenClaim, authParsed := authToken.(*structs.AuthTokenClaim)
	refreshTokenClaim, refreshParsed := refreshToken.(*structs.RefreshTokenClaim)

	if isAuthenticated && authParsed && refreshParsed {
		// Add auth token to context for user later
		c.Locals(LOCALS_USER_ID, authTokenClaim.Data.Id)
		c.Locals(LOCALS_USER_NAME, authTokenClaim.Data.Name)
		c.Locals(LOCALS_USER_EMAIL, authTokenClaim.Data.Email)
		c.Locals(LOCALS_REFRESH_TOKEN_ID, refreshTokenClaim.Data.UserId)
		// Proceed with request
		return c.Next()
	} else {
		return c.SendStatus(failUnauthenticatedRequest(c))
	}
}

/*
Given a user & key returns a JWT auth token signed by the key as a fiber cookie.
Returns nil if an error occurs.
JWT Cookie is defined in structs.AuthTokenClaim
*/
func createAuthCookieFromUser(user structs.UserWithId, key string) *fiber.Cookie {
	return createAuthCookie(user.Email, user.FirstName+" "+user.LastName, user.ID.Hex(), key)
}

/*
Given the params returns a JWT auth token signed by the key as a fiber cookie.
Returns nil if an error occurs.
JWT Cookie is defined in structs.AuthTokenClaim
*/
func createAuthCookie(email string, name string, id string, key string) *fiber.Cookie {
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, structs.AuthTokenClaim{
		Data: structs.AuthTokenJWT{
			Email: email,
			Name:  name,
			Id:    id,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(constants.AuthTokenExpirySeconds * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})
	// Sign token
	signedToken, err := token.SignedString([]byte(key))
	if err != nil {
		log.Print(err)
		return nil
	}
	// Return cookie w/ jwt
	return &fiber.Cookie{
		Name:     AUTH_COOKIE_NAME,
		Value:    signedToken,
		HTTPOnly: false,
		SameSite: "strict",
		MaxAge:   constants.RefreshTokenExpirySeconds,
	}
}

/*
Given a user & key returns a JWT refresh token signed by the key as a fiber cookie.
Returns nil if an error occurs.
Note: Also writes a hashed secret to the DB to validate the refresh token later.
JWT Cookie is defined in structs.RefreshTokenClaim
*/
func createRefreshCookie(userId string, key string) *fiber.Cookie {
	refreshSecret, secretId, err := generateAndPrepareRefreshSecret(userId)
	if err != nil {
		return nil
	}
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, structs.RefreshTokenClaim{
		Data: structs.RefreshTokenJWT{
			UserId:   userId,
			Secret:   refreshSecret,
			SecretId: secretId,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(constants.RefreshTokenExpirySeconds * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})
	// Sign token
	signedToken, err := token.SignedString([]byte(key))
	if err != nil {
		log.Print(err)
		return nil
	}
	// Return cookie w/ JWT
	return &fiber.Cookie{
		Name:     REFRESH_COOKIE_NAME,
		Value:    signedToken,
		HTTPOnly: true,
		SameSite: "strict",
		MaxAge:   constants.RefreshTokenExpirySeconds,
	}
}

/*
Generates a random secret, stores its hash in the refreshToken DB, and returns
the secret in plaintext, along with the DB object id. Returns an error
if anything goes wrong.
Returns (secret, id, err)
*/
func generateAndPrepareRefreshSecret(userId string) (string, string, error) {
	secret, hashedSecret, err := generateSecret()
	if err != nil {
		return "", "", err
	}
	// Write secret to db
	id, err := persistence.WriteRefreshToken(structs.RefreshToken{
		Secret: hashedSecret,
		UserId: userId,
	})
	if err != nil {
		return "", "", err
	}
	return secret, id, nil
}

/*
Generates a secret, returning the (plaintext_secret, hashed_secret).
*/
func generateSecret() (string, string, error) {
	// Generate secret
	secret, err := generateRandomByteArray(SECRET_LENGTH)
	if err != nil {
		return "", "", err
	}
	// Create hash
	hashedPassword, err := bcrypt.GenerateFromPassword(secret, bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}
	return string(secret), string(hashedPassword), err
}

/*
Generates a cryptographically random byte array of length n (capable of being
converted to a string)
Source: https://gist.github.com/dopey/c69559607800d2f2f90b1b1ed4e550fb
*/
func generateRandomByteArray(n int) ([]byte, error) {
	const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz~!@#$%^&*()_+{}|:<>?,./;[=]"
	res := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return nil, err
		}
		res[i] = chars[num.Int64()]
	}
	return res, nil
}

/*
Parses the given tokenString (JWT) and verifies its signature using privateKey.
Returns the token as a map upon success, otherwise returns claims, along with
and error.
*/
func parseJWT(tokenString string, claim jwt.Claims, privateKey string) (jwt.Claims, error) {
	// Validate token
	token, err := jwt.ParseWithClaims(tokenString, claim, func(token *jwt.Token) (interface{}, error) {
		// Check signing method
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// Return private key
		return []byte(privateKey), nil
	})
	// Return error if invalid
	if err != nil || !token.Valid {
		return token.Claims, err
	}
	return token.Claims, err
}

/*
Uses the given refreshToken, and refreshes the auth & refresh cookie, placing
them into the fiber.ctx. Returns true upon success, false otherwise.
*/
func refreshCookies(c *fiber.Ctx, refreshToken jwt.Claims, authToken jwt.Claims, authPrivateKey string, refreshPrivateKey string) bool {
	// Parse the token into a refreshToken struct
	refreshTokenClaim, refreshParsed := refreshToken.(*structs.RefreshTokenClaim)
	authTokenClaim, authParsed := authToken.(*structs.AuthTokenClaim)
	if refreshParsed && authParsed {
		// Attempt to use refresh token
		if useRefreshToken(refreshTokenClaim.Data.UserId,
			refreshTokenClaim.Data.SecretId, refreshTokenClaim.Data.Secret) {
			// Refresh token used -- refresh credentials
			refreshedAuthCookie := createAuthCookie(authTokenClaim.Data.Email,
				authTokenClaim.Data.Name, authTokenClaim.Data.Id, authPrivateKey)
			refreshedRefreshCookie := createRefreshCookie(authTokenClaim.Data.Id, refreshPrivateKey)
			if refreshedAuthCookie != nil && refreshedRefreshCookie != nil {
				c.Cookie(refreshedAuthCookie)
				c.Cookie(refreshedRefreshCookie)
				return true
			}
		}
	}
	return false
}

/*
Given the userId, refresh token id & refresh token secret, deletes the
database record for the refresh token, and checks if the refreshSecret
is equal to the hashed one in the database.
Returns true if the token is successfully used, false otherwise.
*/
func useRefreshToken(userId string, refreshId string, refreshSecret string) bool {
	refreshSecretHashed, err := persistence.GetAndDeleteRefreshTokenSecret(userId, refreshId)
	if err != nil {
		log.Print(err)
		return false
	}
	// Validate refresh secret
	return bcrypt.CompareHashAndPassword([]byte(refreshSecretHashed), []byte(refreshSecret)) == nil
}

// Clears the auth & refresh cookies & returns 401
func failUnauthenticatedRequest(c *fiber.Ctx) int {
	c.ClearCookie(AUTH_COOKIE_NAME)
	c.ClearCookie(REFRESH_COOKIE_NAME)
	return http.StatusUnauthorized
}
