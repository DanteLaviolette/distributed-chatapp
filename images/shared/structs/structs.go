package structs

import (
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Data type representing User to be created in DB
type User struct {
	Email     string
	FirstName string
	LastName  string
	Password  string
}

// Data type representing User stored in DB
type UserWithId struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Email     string
	FirstName string
	LastName  string
	Password  string
}

// Data type representing user log in request
type LoginRequest struct {
	Email    string
	Password string
}

// Data type representing change password request
type ChangePasswordRequest struct {
	Password string
}

// Refresh token for writing to db
type RefreshToken struct {
	UserId string
	Secret string
}

// Refresh token with id (for responses from db)
type RefreshTokenWithId struct {
	ID     primitive.ObjectID `bson:"_id" json:"id"`
	UserId string
	Secret string
}

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
