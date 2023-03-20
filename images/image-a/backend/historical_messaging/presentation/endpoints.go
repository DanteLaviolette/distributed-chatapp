package presentation

import (
	"encoding/json"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.violettedev.com/eecs4222/historical_messaging/business"
	"go.violettedev.com/eecs4222/historical_messaging/structs"
)

/*
Page endpoint for getting a page of messages. Takes PageRequest as input.
- 400 on bad input
- 500 if something goes wrong
- 200 with json of messages on success
*/
func GetMessagesEndpoint(c *fiber.Ctx) error {
	// Parse body to struct
	var pageRequest structs.PageRequest
	if err := c.BodyParser(&pageRequest); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	// Get messages
	messages, err := business.GetMessages(pageRequest)
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
