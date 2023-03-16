package livechat

import (
	"github.com/gofiber/websocket/v2"
	"go.violettedev.com/eecs4222/livechat/structs"
)

func handleAuthMessage(message string, c *websocket.Conn) {

}

func handlePing(message string, c *websocket.Conn) {
	c.WriteJSON(structs.Message{
		Type: "pong",
	})
}
