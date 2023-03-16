package structs

import "github.com/golang-jwt/jwt/v5"

// Auth Token JWT definition
type AuthTokenJWT struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Auth token claim (to be used by jwt lib)
type AuthTokenClaim struct {
	Data AuthTokenJWT `json:"data"`
	jwt.RegisteredClaims
}

// Refresh token JWT definition
type RefreshTokenJWT struct {
	UserId   string `json:"userId"`
	Secret   string `json:"secret"`
	SecretId string `json:"secretId"`
}

// Refresh token claim (to be used by jwt lib)
type RefreshTokenClaim struct {
	Data RefreshTokenJWT `json:"data"`
	jwt.RegisteredClaims
}
