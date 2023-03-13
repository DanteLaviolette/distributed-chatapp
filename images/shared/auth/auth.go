package auth

import (
	"shared/constants"
	"shared/structs"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

/*
Given a user & key returns a JWT auth token signed by the key as a fiber cookie.
Returns nil if an error occurs.
*/
func CreateAuthCookie(user structs.UserWithId, key string) *fiber.Cookie {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"name":  user.FirstName + " " + user.LastName,
		"id":    user.ID.Hex(),
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Unix() + constants.AuthTokenExpirySeconds,
	})
	signedToken, err := token.SignedString([]byte(key))
	if err != nil {
		return nil
	}
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
*/
func CreateRefreshCookie(user structs.UserWithId, key string) *fiber.Cookie {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"name":  user.FirstName + " " + user.LastName,
		"id":    user.ID.Hex(),
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Unix() + constants.RefreshTokenExpirySeconds,
	})
	signedToken, err := token.SignedString([]byte(key))
	if err != nil {
		return nil
	}
	return &fiber.Cookie{
		Name:     "refresh",
		Value:    signedToken,
		HTTPOnly: true,
		SameSite: "strict",
		MaxAge:   constants.RefreshTokenExpirySeconds,
	}
}
