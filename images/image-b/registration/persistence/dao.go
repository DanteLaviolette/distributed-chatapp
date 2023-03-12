package persistence

import (
	"registration/structs"
	"shared/persistence"
)

/*
Adds user to the database, returning the db error if the insertion failed, or
null otherwise.
*/
func InsertUser(registerInfo structs.RegisterInfo) error {
	collection, ctx := persistence.GetUserCollection()
	_, err := collection.InsertOne(ctx, registerInfo)
	return err
}
