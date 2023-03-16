// Handles creation/usage of refresh JWTs
package refresh_token

import (
	"crypto/rand"
	"log"
	"math/big"
	"time"

	"go.violettedev.com/eecs4222/shared/auth/internal/structs"

	"go.violettedev.com/eecs4222/shared/auth/internal/dao"
	"go.violettedev.com/eecs4222/shared/auth/internal/token_util"
	"go.violettedev.com/eecs4222/shared/constants"

	"go.violettedev.com/eecs4222/shared/database/schema"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const SECRET_LENGTH = 32

/*
Parses & returns refresh token from request.
Returns an error if anything goes wrong. Returns the token and the error if the
error is related due to the token being invalid (ie. expired or bad signature)
*/
func ParseRefreshToken(refreshJWT string, refreshPrivateKey string) (*structs.RefreshTokenClaim, error) {
	// Validate & parse tokens
	refreshToken, err := token_util.ParseJWT(refreshJWT, &structs.RefreshTokenClaim{}, refreshPrivateKey)
	// Convert claims to specific jwt claims
	return claimToRefreshTokenClaim(refreshToken), err
}

/*
Given the userId, refresh token id & refresh token secret, deletes the
database record for the refresh token, and checks if the refreshSecret
is equal to the hashed one in the database.
Returns true if the token is successfully used, false otherwise.
*/
func UseRefreshToken(userId string, refreshId string, refreshSecret string) bool {
	// NOTE: We don't need to validate refresh token expiry, as db will automatically
	// delete records once the key expires
	refreshSecretHashed, err := dao.GetAndDeleteRefreshTokenSecret(userId, refreshId)
	if err != nil {
		log.Print(err)
		return false
	}
	// Validate refresh secret
	return bcrypt.CompareHashAndPassword([]byte(refreshSecretHashed), []byte(refreshSecret)) == nil
}

/*
Given a user & key returns a JWT refresh token signed by the key.
Returns nil if an error occurs.
Note: Also writes a hashed secret to the DB to validate the refresh token later.
JWT is defined in structs.RefreshTokenClaim
*/
func CreateRefreshToken(userId string, key string) (string, error) {
	refreshSecret, secretId, err := generateAndPrepareRefreshSecret(userId)
	if err != nil {
		return "", err
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
		return "", err
	}
	return signedToken, nil
}

/*
Deletes the refresh token from the database.
*/
func InvalidateRefreshToken(refreshId string) {
	err := dao.DeleteRefreshTokenById(refreshId)
	if err != nil {
		log.Print(err)
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
	id, err := dao.WriteRefreshToken(schema.RefreshTokenSchema{
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
Returns the claim as a RefreshToken claim if successful. Returns nil if unsuccessful
*/
func claimToRefreshTokenClaim(claim jwt.Claims) *structs.RefreshTokenClaim {
	refreshTokenClaim, refreshParsed := claim.(*structs.RefreshTokenClaim)
	if !refreshParsed {
		return nil
	}
	return refreshTokenClaim
}
