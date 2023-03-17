// Handles the usage of credentials to handle authentication in the
// context of a fiber request
package auth

import (
	"errors"
	"log"
	"net/http"

	"go.violettedev.com/eecs4222/shared/auth/internal/auth_token"
	"go.violettedev.com/eecs4222/shared/auth/internal/refresh_token"
	"go.violettedev.com/eecs4222/shared/auth/internal/structs"
	"go.violettedev.com/eecs4222/shared/constants"

	"go.violettedev.com/eecs4222/shared/database/schema"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/golang-jwt/jwt/v5"
)

const AUTH_HEADER_NAME = "authorization"
const REFRESH_COOKIE_NAME = "refresh"

const REFRESH_NEEDED_WS_ERROR = "refresh"

// Local keys
const LOCALS_USER_ID = "userId"
const LOCALS_USER_NAME = "userName"
const LOCALS_USER_EMAIL = "userEmail"
const LOCALS_REFRESH_TOKEN_ID = "refreshTokenId"

type AuthProvider struct {
	/*
		Generates and adds credentials to the response. Returns true on success,
		false on failure.
	*/
	GenerateAndSetCredentials func(schema.UserSchema, *fiber.Ctx) bool
	/*
		To be used as fiber middleware. Proceeds if user is logged in (potentially
		refreshing their credentials). Fails the request with 401 error if not logged
		in.
		Upon success, adds userId, userName, userEmail & refreshTokenId to c.Locals
	*/
	IsAuthenticatedFiberMiddleware func(*fiber.Ctx) error
	/*
		Returns auth info if its valid.
		Error "refresh" means the token is expired.
		Returns: (id, name, email, error)
	*/
	GetAuthContextWebSocket func(*websocket.Conn, string) (string, string, string, error)
}

func Initialize(authPrivateKey string, refreshPrivateKey string) *AuthProvider {
	return &AuthProvider{
		IsAuthenticatedFiberMiddleware: func(c *fiber.Ctx) error {
			return isAuthenticatedFiberMiddleware(c, authPrivateKey, refreshPrivateKey)
		},
		GenerateAndSetCredentials: func(user schema.UserSchema, c *fiber.Ctx) bool {
			return generateAndSetAuthHeaderAndRefreshToken(user, c,
				authPrivateKey, refreshPrivateKey) == nil
		},
		GetAuthContextWebSocket: func(c *websocket.Conn, authToken string) (string, string, string, error) {
			auth, err := getAuthContextWebSocket(c, authToken, authPrivateKey)
			if err != nil {
				return "", "", "", err
			}
			return auth.Data.Id, auth.Data.Name, auth.Data.Email, nil
		},
	}
}

/*
Invalidates the users credentials, returning empty credentials
*/
func InvalidateCredentials(c *fiber.Ctx) {
	// Clear auth values
	c.Response().Header.Add(AUTH_HEADER_NAME, "")
	c.ClearCookie(REFRESH_COOKIE_NAME)
	// Get refresh token id
	refreshTokenId := c.Locals(LOCALS_REFRESH_TOKEN_ID)
	refreshTokenIdString, ok := refreshTokenId.(string)
	// Invalidate refresh token
	if ok {
		refresh_token.InvalidateRefreshToken(refreshTokenIdString)
	} else {
		log.Println("Refresh token not found in locals")
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
	// Get request credentials
	authString, refreshString, err := getRequestCredentials(c)
	if err != nil {
		return c.SendStatus(failUnauthenticatedRequest(c))
	}
	// Parse refresh token
	refreshToken, err := refresh_token.ParseRefreshToken(refreshString, refreshPrivateKey)
	if err != nil {
		return c.SendStatus(failUnauthenticatedRequest(c))
	}
	// Parse auth token
	authToken, err := auth_token.ParseAuthToken(authString, authPrivateKey)
	isAuthenticated := false // Assume token is invalid
	// Refresh token if needed
	if errors.Is(err, jwt.ErrTokenExpired) {
		// Expired auth token -- attempt to refresh token
		authString, refreshString, err = refreshCredentials(c, refreshToken,
			authToken, authPrivateKey, refreshPrivateKey)
		// Handle case that refresh failed
		if err != nil {
			return c.SendStatus(failUnauthenticatedRequest(c))
		}
		// Refresh was successful -- set response credentials
		setResponseCredentials(c, authString, refreshString)
		// Parse updated tokens
		refreshToken, err = refresh_token.ParseRefreshToken(refreshString, refreshPrivateKey)
		var err2 error
		authToken, err2 = auth_token.ParseAuthToken(authString, authPrivateKey)
		if err != nil || err2 != nil {
			return c.SendStatus(failUnauthenticatedRequest(c))
		}
		isAuthenticated = true
	} else if err == nil && authToken != nil && refreshToken != nil {
		// Success case where auth token is valid and not expired
		isAuthenticated = true
	}
	if isAuthenticated {
		// Populate locals
		setFiberContextAuthLocals(c, *authToken, *refreshToken)
		// Proceed with request
		return c.Next()
	} else {
		return c.SendStatus(failUnauthenticatedRequest(c))
	}
}

/*
Returns auth info, or error if there is no (or invalid) auth info.
Error "refresh" means the token is expired.
*/
func getAuthContextWebSocket(c *websocket.Conn, authTokenString string,
	authPrivateKey string) (*structs.AuthTokenClaim, error) {
	authToken, err := auth_token.ParseAuthToken(authTokenString, authPrivateKey)
	if errors.Is(err, jwt.ErrTokenExpired) {
		return nil, errors.New(REFRESH_NEEDED_WS_ERROR)
	} else if err != nil {
		return nil, err
	}
	return authToken, nil
}

/*
Generates and adds credentials to the response.
*/
func generateAndSetAuthHeaderAndRefreshToken(user schema.UserSchema,
	c *fiber.Ctx, authPrivateKey string, refreshPrivateKey string) error {
	authToken, err := auth_token.CreateAuthTokenFromUser(user, authPrivateKey)
	if err != nil {
		return err
	}
	refreshToken, err := refresh_token.CreateRefreshToken(user.ID.Hex(), refreshPrivateKey)
	if err != nil {
		return err
	}
	setResponseCredentials(c, authToken, refreshToken)
	return nil
}

/*
Uses the given refreshToken, and refreshes the auth & refresh cookie,
Returning the new (authJWT, refreshJWT, error).
Tokens will be empty strings if an error occurs.
*/
func refreshCredentials(c *fiber.Ctx, refreshToken *structs.RefreshTokenClaim,
	authToken *structs.AuthTokenClaim, authPrivateKey string,
	refreshPrivateKey string) (string, string, error) {
	// Attempt to use refresh token
	if !refresh_token.UseRefreshToken(refreshToken.Data.UserId,
		refreshToken.Data.SecretId, refreshToken.Data.Secret) {
		return "", "", errors.New("Failed to use refresh token")
	}
	// Refresh token used -- refresh credentials
	refreshedAuthToken, err := auth_token.CreateAuthToken(authToken.Data.Email,
		authToken.Data.Name, authToken.Data.Id, authPrivateKey)
	if err != nil {
		log.Print(err)
		return "", "", err
	}
	refreshedRefreshToken, err := refresh_token.CreateRefreshToken(authToken.Data.Id, refreshPrivateKey)
	if err != nil {
		log.Print(err)
		return "", "", err
	}
	return refreshedAuthToken, refreshedRefreshToken, nil
}

/*
Get credentials from request as string JWTs. Returns error if anything
goes wrong.
Returns: (auth, refresh, error)
*/
func getRequestCredentials(c *fiber.Ctx) (string, string, error) {
	// Get header
	auth := c.Request().Header.Peek(AUTH_HEADER_NAME)
	// Fail if cookies weren't found
	if auth == nil {
		return "", "", errors.New("credentials not found")
	}
	// Get cookie
	refresh := c.Cookies(REFRESH_COOKIE_NAME)
	// Fail if cookies weren't found
	if refresh == "" {
		return "", "", errors.New("credentials not found")
	}
	return string(auth), refresh, nil
}

/*
Adds the authToken & refresh token to the fiber context response.
Auth token is added to Authorization header
Refresh token is added as cookie
*/
func setResponseCredentials(c *fiber.Ctx, authToken string, refreshJWT string) {
	c.Response().Header.Add(AUTH_HEADER_NAME, authToken)
	c.Cookie(createRefreshCookie(refreshJWT))
}

/*
Creates a refresh cookie from the given JWT string
*/
func createRefreshCookie(refreshJWT string) *fiber.Cookie {
	return &fiber.Cookie{
		Name:     REFRESH_COOKIE_NAME,
		Value:    refreshJWT,
		HTTPOnly: true,
		SameSite: "strict",
		MaxAge:   constants.RefreshTokenExpirySeconds,
	}
}

/*
Adds userId, userName, userEmail & refreshTokenId to c.Locals
*/
func setFiberContextAuthLocals(c *fiber.Ctx, authToken structs.AuthTokenClaim,
	refreshToken structs.RefreshTokenClaim) {
	// Add auth token to context for user later
	c.Locals(LOCALS_USER_ID, authToken.Data.Id)
	c.Locals(LOCALS_USER_NAME, authToken.Data.Name)
	c.Locals(LOCALS_USER_EMAIL, authToken.Data.Email)
	c.Locals(LOCALS_REFRESH_TOKEN_ID, refreshToken.Data.UserId)
}

// Clears the auth & refresh credentials & returns 401
func failUnauthenticatedRequest(c *fiber.Ctx) int {
	c.Response().Header.Add(AUTH_HEADER_NAME, "")
	c.ClearCookie(REFRESH_COOKIE_NAME)
	return http.StatusUnauthorized
}
