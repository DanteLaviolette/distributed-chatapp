package persistence

import (
	"context"

	"go.violettedev.com/eecs4222/constants"
	"go.violettedev.com/eecs4222/database"
	"go.violettedev.com/eecs4222/structs"

	"go.mongodb.org/mongo-driver/bson"
)

/*
Gets user with the given email, returning an error upon failure.
*/
func GetUserWithId(email string) (structs.UserWithId, error) {
	var res structs.UserWithId
	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), constants.DatabaseTimeout)
	defer cancel()
	collection := database.GetUserCollection()
	err := collection.FindOne(ctx, bson.M{
		"email": email,
	}).Decode(&res)
	return res, err
}
