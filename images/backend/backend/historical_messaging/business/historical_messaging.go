package business

import (
	"log"

	"go.violettedev.com/eecs4222/historical_messaging/persistence"
	"go.violettedev.com/eecs4222/shared/database/schema"
)

/*
Returns a page (list) of messages based on lastTimestamp. List will be
empty if there are no more messages.
Returns nil, error if anything goes wrong.
*/
func GetMessages(lastTimestamp int64) ([]schema.MessageSchema, error) {
	messages, err := persistence.GetMessagesByPage(lastTimestamp)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return messages, nil
}
