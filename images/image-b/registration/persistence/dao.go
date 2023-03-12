package persistence

import (
	"shared/persistence"
	"shared/structs"
)

/*
Adds user to the database, returning the db error if the insertion failed, or
null otherwise.
*/
func InsertUser(user structs.User) error {
	collection, ctx := persistence.GetUserCollection()
	_, err := collection.InsertOne(ctx, user)
	return err
}
