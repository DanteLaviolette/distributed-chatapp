package persistence

import (
	"context"

	"go.violettedev.com/eecs4222/shared/constants"
	"go.violettedev.com/eecs4222/shared/database"
	"go.violettedev.com/eecs4222/shared/database/schema"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
Writes the given message to the database.
Returns an error if something goes wrong.
Returns the messages id otherwise.
*/
func WriteMessage(message schema.MessageSchema) (primitive.ObjectID, error) {
	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), constants.DatabaseTimeout)
	defer cancel()
	collection := database.GetMessageCollection()
	// Insert message
	res, err := collection.InsertOne(ctx, message)
	if err != nil {
		return message.ID, err
	}
	// Return id
	return res.InsertedID.(primitive.ObjectID), err
}
