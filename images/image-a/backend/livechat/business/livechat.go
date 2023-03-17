package business

import (
	"os"
	"time"

	"github.com/gofiber/websocket/v2"
	"go.violettedev.com/eecs4222/livechat/structs"
	"go.violettedev.com/eecs4222/shared/auth"
)

var socketIdsToConnection = make(map[string]*websocket.Conn)

func HandleConnectionOpened(c *websocket.Conn, authCtx *structs.AuthContext) {
	socketIdsToConnection[authCtx.SocketId] = c
}

func HandleConnectionClosed(authCtx *structs.AuthContext) {
	delete(socketIdsToConnection, authCtx.SocketId)
}

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
		// TODO: Write message to db
		// TODO: Notify of messages via redis
		message := structs.ChatMessage{
			Type:    "message",
			Message: content,
			Name:    name,
			Email:   email,
			Ts:      time.Now().UnixMilli(),
		}
		sendChatMessageToEveryone(message)
	}
}

// Send a chat message to all websockets
func sendChatMessageToEveryone(message structs.ChatMessage) {
	for _, conn := range socketIdsToConnection {
		if conn != nil {
			conn.WriteJSON(message)
		}
	}
}

// Send pong on ping recipient (heartbeat)
func HandlePing(c *websocket.Conn, content string) {
	c.WriteJSON(structs.Message{
		Type: "pong",
	})
}
