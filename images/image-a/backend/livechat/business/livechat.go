package business

import (
	"os"
	"time"

	"github.com/gofiber/websocket/v2"
	"go.violettedev.com/eecs4222/livechat/structs"
	"go.violettedev.com/eecs4222/shared/auth"
)

/*
Handle authentication method:
- add credentials to authCtx if signed in
- return "refresh" message if credentials are expired
- do nothing if not signed in
*/
func HandleAuthMessage(c *websocket.Conn, authCtx *structs.AuthContext, content string) {
	authProvider := auth.Initialize(os.Getenv("AUTH_PRIVATE_KEY"),
		os.Getenv("REFRESH_PRIVATE_KEY"))
	name, email, err := authProvider.GetAuthContextWebSocket(c, content)
	if err != nil && err.Error() == "refresh" {
		// Request refresh
		c.WriteJSON(structs.Message{
			Type: "refresh",
		})
	} else {
		authCtx.Name = name
		authCtx.Email = email
		c.WriteJSON(structs.Message{
			Type: "signed_in",
		})
	}
	// Note: We don't send anything on no auth -- as it's default behavior
}

// Sends a message to all clients -- can only be performed by logged in users
func HandleChatMessage(c *websocket.Conn, authCtx *structs.AuthContext, content string) {
	name := authCtx.Name
	email := authCtx.Email
	if name != "" && email != "" {
		ts := time.Now().UnixMilli()
		c.WriteJSON(structs.ChatMessage{
			Type:    "message",
			Message: content,
			Name:    name,
			Email:   email,
			Ts:      ts,
		})
	}
}

// Send pong on ping recipient (heartbeat)
func HandlePing(c *websocket.Conn, content string) {
	c.WriteJSON(structs.Message{
		Type: "pong",
	})
}
