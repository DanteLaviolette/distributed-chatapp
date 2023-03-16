// Handles creation of auth JWTs
package auth_token

import (
	"log"
	"shared/constants"
	"shared/internal/auth/token_util"
	"shared/structs"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

/*
Parses & returns auth token.
Returns an error if anything goes wrong. Returns the token and the error if the
error is related due to the token being invalid (ie. expired or bad signature)
*/
func ParseAuthToken(authJWT string, authPrivateKey string) (*structs.AuthTokenClaim, error) {
	// Validate & parse token
	authToken, err := token_util.ParseJWT(authJWT, &structs.AuthTokenClaim{}, authPrivateKey)
	// Convert claims to specific jwt claims
	return claimToAuthTokenClaim(authToken), err
}

/*
Given a user & key returns a JWT auth token signed by the key as a
string. Returns empty str with error if an error occurs.
JWT contains the fields id, email, name, iat & exp.
*/
func CreateAuthTokenFromUser(user structs.UserWithId, key string) (string, error) {
	return CreateAuthToken(user.Email, user.FirstName+" "+user.LastName, user.ID.Hex(), key)
}

/*
Given the params returns a JWT auth token signed by the key as a
string. Returns empty str with error if an error occurs.
JWT contains the fields id, email, name, iat & exp.
*/
func CreateAuthToken(email string, name string, id string, key string) (string, error) {
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
		return "", err
	}
	return signedToken, nil
}

/*
Returns the claim as a AuthToken claim if successful. Returns nil if unsuccessful
*/
func claimToAuthTokenClaim(claim jwt.Claims) *structs.AuthTokenClaim {
	authTokenClaim, authParsed := claim.(*structs.AuthTokenClaim)
	if !authParsed {
		return nil
	}
	return authTokenClaim
}
