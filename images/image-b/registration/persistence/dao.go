package persistence

import (
	"context"
	"shared/constants"
	"shared/persistence"
	"shared/structs"
)

/*
Adds user to the database, returning the db error if the insertion failed, or
null otherwise.
*/
func InsertUser(user structs.User) error {
	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), constants.DatabaseTimeout)
	defer cancel()
	collection := persistence.GetUserCollection()
	_, err := collection.InsertOne(ctx, user)
	return err
}
