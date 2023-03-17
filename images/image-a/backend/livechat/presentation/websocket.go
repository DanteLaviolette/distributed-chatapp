package presentation

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"go.violettedev.com/eecs4222/livechat/business"
	"go.violettedev.com/eecs4222/livechat/structs"
)

func CanUpgradeToWebSocket(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func LiveChatWebSocket(c *websocket.Conn) {
	var authCtx *structs.AuthContext = &structs.AuthContext{}
	for {
		var message structs.Message
		err := c.ReadJSON(&message)
		if err != nil {
			log.Print(err)
			return
		}
		if message.Type == "ping" {
			business.HandlePing(c, message.Content)
		} else if message.Type == "auth" {
			business.HandleAuthMessage(c, authCtx, message.Content)
		}
	}
}
