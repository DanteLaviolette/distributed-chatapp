package token_util

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

/*
Parses the given tokenString (JWT) and verifies its signature using privateKey.
Returns the token as a map upon success, otherwise returns claims, along with
and error.
*/
func ParseJWT(tokenString string, claim jwt.Claims, privateKey string) (jwt.Claims, error) {
	// Validate token
	token, err := jwt.ParseWithClaims(tokenString, claim, func(token *jwt.Token) (interface{}, error) {
		// Check signing method
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
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
