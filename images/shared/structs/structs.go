package structs

import "go.mongodb.org/mongo-driver/bson/primitive"

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

type RefreshToken struct {
	UserId string
	Secret string
}

type RefreshTokenWithId struct {
	ID     primitive.ObjectID `bson:"_id" json:"id"`
	UserId string
	Secret string
}
