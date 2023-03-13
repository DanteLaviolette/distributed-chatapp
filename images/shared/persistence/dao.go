package persistence

import (
	"context"
	"shared/constants"
	"shared/structs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
Writes the refresh token to the database, returning the objects ID or an
error if one occurred.
*/
func WriteRefreshToken(refreshToken structs.RefreshToken) (string, error) {
	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), constants.DatabaseTimeout)
	defer cancel()
	collection := GetRefreshTokenCollection()
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
	var res structs.RefreshTokenWithId
	refreshId, err := primitive.ObjectIDFromHex(refreshIdHex)
	if err != nil {
		return "", err
	}
	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), constants.DatabaseTimeout)
	defer cancel()
	collection := GetRefreshTokenCollection()
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