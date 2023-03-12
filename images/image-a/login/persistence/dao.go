package persistence

import (
	"context"
	"shared/constants"
	"shared/persistence"
	"shared/structs"

	"go.mongodb.org/mongo-driver/bson"
)

/*
Gets user with the given email, returning an error upon failure.
*/
func GetUser(email string) (structs.User, error) {
	var res structs.User
	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), constants.DatabaseTimeout)
	defer cancel()
	collection := persistence.GetUserCollection()
	err := collection.FindOne(ctx, bson.M{
		"email": email,
	}).Decode(&res)
	return res, err
}
