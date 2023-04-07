package presentation

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.violettedev.com/eecs4222/historical_messaging/business"
)

/*
Page endpoint for getting a page of messages. Takes PageRequest as input.
- 400 on bad input
- 500 if something goes wrong
- 200 with json of messages on success
*/
func GetMessagesEndpoint(c *fiber.Ctx) error {
	// Parse query param
	lastTimestampString := c.Query("lastTimestamp", "0")
	lastTimestamp, err := strconv.ParseInt(lastTimestampString, 10, 64)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	// Get messages
	messages, err := business.GetMessages(lastTimestamp)
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}
	// Parse to JSON
	messagesJson, err := json.Marshal(messages)
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}
	// Send messages
	c.Send(messagesJson)
	// Return status code
	return c.SendStatus(200)
}
