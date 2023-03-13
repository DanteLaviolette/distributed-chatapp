package persistence

import (
	"context"
	"shared/constants"
	"shared/structs"

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
