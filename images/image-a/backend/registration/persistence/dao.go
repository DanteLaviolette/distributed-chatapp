package persistence

import (
	"context"

	"go.violettedev.com/eecs4222/shared/constants"
	"go.violettedev.com/eecs4222/shared/database"
	"go.violettedev.com/eecs4222/shared/database/schema"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
Adds user to the database, returning the db error if the insertion failed, or
null otherwise.
*/
func InsertUser(user schema.UserSchema) error {
	// Set random ID
	user.ID = primitive.NewObjectID()
	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), constants.DatabaseTimeout)
	defer cancel()
	collection := database.GetUserCollection()
	_, err := collection.InsertOne(ctx, user)
	return err
}

/*
Updates the password for the given userId. Returns nil upon success, error otherwise.
*/
func UpdatePasswordForUserId(userId string, passwordHash string) error {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), constants.DatabaseTimeout)
	defer cancel()
	collection := database.GetUserCollection()
	_, err = collection.UpdateOne(ctx, bson.M{
		"_id": id,
	}, bson.M{
		"$set": bson.M{"password": passwordHash},
	})
	return err
}
