package persistence

import (
	"context"
	"shared/constants"
	"shared/persistence"
	"shared/structs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
Gets user with the given email, returning an error upon failure.
*/
func GetUserWithId(email string) (structs.UserWithId, error) {
	var res structs.UserWithId
	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), constants.DatabaseTimeout)
	defer cancel()
	collection := persistence.GetUserCollection()
	err := collection.FindOne(ctx, bson.M{
		"email": email,
	}).Decode(&res)
	return res, err
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
	collection := persistence.GetRefreshTokenCollection()
	// Delete refresh token
	_, err = collection.DeleteOne(ctx, bson.M{
		"_id": refreshId,
	})
	if err != nil {
		return err
	}
	return nil
}
