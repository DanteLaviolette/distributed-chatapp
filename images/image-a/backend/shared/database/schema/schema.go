package schema

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Data type representing User stored in DB
type UserSchema struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Email     string
	FirstName string
	LastName  string
	Password  string
}

// Refresh token with id (for responses from db)
type RefreshTokenSchema struct {
	ID     primitive.ObjectID `bson:"_id" json:"id"`
	UserId string
	Secret string
}
