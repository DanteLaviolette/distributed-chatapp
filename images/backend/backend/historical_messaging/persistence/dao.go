package persistence

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.violettedev.com/eecs4222/shared/constants"
	"go.violettedev.com/eecs4222/shared/database"
	"go.violettedev.com/eecs4222/shared/database/schema"
)

var PAGE_LIMIT int64 = 50

/*
Returns up to PAGE_LIMIT messages with a timestamp lower than lastTimestamp.
Returns nil, error if anything goes wrong.
Returns [], nil if nothing is found
*/
func GetMessagesByPage(lastTimestamp int64) ([]schema.MessageSchema, error) {
	var res []schema.MessageSchema
	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), constants.DatabaseTimeout)
	defer cancel()
	collection := database.GetMessageCollection()
	// Find PAGE_LIMIT messages w/ lower timestamp sorted in descending order
	cursor, err := collection.Find(ctx, bson.M{
		"ts": bson.M{"$lt": lastTimestamp},
	}, &options.FindOptions{
		Limit: &PAGE_LIMIT,
		Sort:  bson.M{"ts": -1},
	})
	if err != nil {
		return nil, err
	}
	// Decode all messages into res
	err = cursor.All(ctx, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
