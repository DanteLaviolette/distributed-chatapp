package persistence

import (
	"context"

	"go.violettedev.com/eecs4222/shared/constants"
	"go.violettedev.com/eecs4222/shared/database"
	"go.violettedev.com/eecs4222/shared/database/schema"

	"go.mongodb.org/mongo-driver/bson"
)

/*
Gets user with the given email, returning an error upon failure.
*/
func GetUser(email string) (schema.UserSchema, error) {
	var res schema.UserSchema
	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), constants.DatabaseTimeout)
	defer cancel()
	collection := database.GetUserCollection()
	err := collection.FindOne(ctx, bson.M{
		"email": email,
	}).Decode(&res)
	return res, err
}
