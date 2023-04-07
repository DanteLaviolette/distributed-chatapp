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

type MessageSchema struct {
	ID      primitive.ObjectID `bson:"_id" json:"id"`
	Subject string             `bson:"subject" json:"subject"`
	Message string             `bson:"message" json:"message"`
	Name    string             `bson:"name" json:"name"`
	Email   string             `bson:"email" json:"email"`
	Ts      int64              `bson:"ts" json:"ts"`
}
