package presentation

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"go.violettedev.com/eecs4222/livechat/structs"
)

func CanUpgradeToWebSocket(c *fiber.Ctx) error {
	println("hit")
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func LiveChatWebSocket(c *websocket.Conn) {
	println("started")
	for {
		var message structs.Message
		err := c.ReadJSON(&message)
		if err != nil {
			log.Print(err)
			return
		}
		log.Println(message)
		c.WriteJSON(message)
	}
}
