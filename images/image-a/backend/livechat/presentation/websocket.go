package presentation

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
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
		messageType, message, err := c.ReadMessage()
		if err != nil {
			log.Print(err)
			return
		}
		println(messageType, string(message))
		c.WriteMessage(websocket.TextMessage, message)
	}
}
