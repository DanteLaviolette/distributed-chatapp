package dao

import (
	"context"

	"go.violettedev.com/eecs4222/shared/constants"
	"go.violettedev.com/eecs4222/shared/database"
	"go.violettedev.com/eecs4222/shared/database/schema"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
Writes the refresh token to the database, returning the objects ID or an
error if one occurred.
*/
func WriteRefreshToken(refreshToken schema.RefreshTokenSchema) (string, error) {
	refreshToken.ID = primitive.NewObjectID()
	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), constants.DatabaseTimeout)
	defer cancel()
	collection := database.GetRefreshTokenCollection()
	// Insert document
	res, err := collection.InsertOne(ctx, refreshToken)
	if err != nil {
		return "", err
	}
	// Return objects ID
	return res.InsertedID.(primitive.ObjectID).Hex(), err
}

/*
Gets & deletes the refresh token of the given userId & refreshId.
Returns the refresh tokens secret (hashed) upon success. Returns an empty string
with an error otherwise.
*/
func GetAndDeleteRefreshTokenSecret(userIdHex string, refreshIdHex string) (string, error) {
	var res schema.RefreshTokenSchema
	refreshId, err := primitive.ObjectIDFromHex(refreshIdHex)
	if err != nil {
		return "", err
	}
	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), constants.DatabaseTimeout)
	defer cancel()
	collection := database.GetRefreshTokenCollection()
	// Find & delete refresh token
	err = collection.FindOneAndDelete(ctx, bson.M{
		"_id":    refreshId,
		"userid": userIdHex,
	}).Decode(&res)
	if err != nil {
		return "", err
	}
	return res.Secret, nil
}

/*
Deletes the refresh token of the given id. Returns an error upon failure,
nil otherwise.
*/
func DeleteRefreshTokenById(id string) error {
	refreshId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), constants.DatabaseTimeout)
	defer cancel()
	collection := database.GetRefreshTokenCollection()
	// Delete refresh token
	_, err = collection.DeleteOne(ctx, bson.M{
		"_id": refreshId,
	})
	if err != nil {
		return err
	}
	return nil
}
