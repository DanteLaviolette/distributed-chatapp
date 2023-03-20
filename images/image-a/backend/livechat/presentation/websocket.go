package presentation

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"go.violettedev.com/eecs4222/livechat/business"
	"go.violettedev.com/eecs4222/livechat/structs"
)

func InitializeDistributedMessaging() {
	business.InitializeDistributedMessaging()
}

func CanUpgradeToWebSocket(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func LiveChatWebSocket(c *websocket.Conn) {
	// Create auth context -- setting uuid for SocketId
	var authCtx *structs.AuthContext = &structs.AuthContext{
		SocketId: uuid.NewString(),
	}
	// On connection close, remove connection
	defer func() {
		business.HandleConnectionClosed(authCtx)
	}()
	// Handle connection open
	business.HandleConnectionOpened(c, authCtx)
	// Handle messages
	for {
		// Parse message to JSON
		var message structs.Message
		err := c.ReadJSON(&message)
		if err != nil {
			log.Print(err)
			return
		}
		// Call business layer based on message type
		if message.Type == "ping" {
			business.HandlePing(authCtx, message.Content)
		} else if message.Type == "auth" {
			business.HandleAuthMessage(c, authCtx, message.Content)
		} else if message.Type == "message" {
			business.HandleChatMessage(authCtx, message.Subject, message.Content)
		}
	}
}
