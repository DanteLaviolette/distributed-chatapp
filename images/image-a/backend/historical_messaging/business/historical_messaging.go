package business

import (
	"log"

	"go.violettedev.com/eecs4222/historical_messaging/persistence"
	"go.violettedev.com/eecs4222/historical_messaging/structs"
	"go.violettedev.com/eecs4222/shared/database/schema"
)

/*
Returns a page (list) of messages based on pageRequest. List will be
empty if there are no more messages.
Returns nil, error if anything goes wrong.
*/
func GetMessages(pageRequest structs.PageRequest) ([]schema.MessageSchema, error) {
	messages, err := persistence.GetMessagesByPage(pageRequest.LastTimestamp)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return messages, nil
}
