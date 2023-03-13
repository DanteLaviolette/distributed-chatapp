package auth

import (
	"crypto/rand"
	"math/big"

	"shared/constants"
	"shared/persistence"
	"shared/structs"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const SECRET_LENGTH = 32

/*
Given a user & key returns a JWT auth token signed by the key as a fiber cookie.
Returns nil if an error occurs.
*/
func CreateAuthCookie(user structs.UserWithId, key string) *fiber.Cookie {
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"name":  user.FirstName + " " + user.LastName,
		"id":    user.ID.Hex(),
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Unix() + constants.AuthTokenExpirySeconds,
	})
	// Sign token
	signedToken, err := token.SignedString([]byte(key))
	if err != nil {
		return nil
	}
	// Return cookie w/ jwt
	return &fiber.Cookie{
		Name:     "auth",
		Value:    signedToken,
		HTTPOnly: false,
		SameSite: "strict",
		MaxAge:   constants.RefreshTokenExpirySeconds,
	}
}

/*
Given a user & key returns a JWT refresh token signed by the key as a fiber cookie.
Returns nil if an error occurs.
Note:
- Also writes a secret to the DB to validate the refresh token later.
*/
func CreateRefreshCookie(user structs.UserWithId, key string) *fiber.Cookie {
	refreshSecret, secretId, err := generateAndPrepareRefreshSecret(user.ID.Hex())
	if err != nil {
		return nil
	}
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   user.ID.Hex(),
		"secret":   refreshSecret,
		"secretId": secretId,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Unix() + constants.RefreshTokenExpirySeconds,
	})
	// Sign token
	signedToken, err := token.SignedString([]byte(key))
	if err != nil {
		return nil
	}
	// Return cookie w/ JWT
	return &fiber.Cookie{
		Name:     "refresh",
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
		return "", "", nil
	}
	// Write secret to db
	id, err := persistence.WriteRefreshToken(structs.RefreshToken{
		Secret: hashedSecret,
		UserId: userId,
	})
	if err != nil {
		return "", "", nil
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
