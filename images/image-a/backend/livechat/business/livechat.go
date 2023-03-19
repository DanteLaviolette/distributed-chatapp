package business

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/websocket/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.violettedev.com/eecs4222/livechat/coordination/messaging"
	"go.violettedev.com/eecs4222/livechat/coordination/user_count"
	"go.violettedev.com/eecs4222/livechat/persistence"
	"go.violettedev.com/eecs4222/livechat/structs"
	"go.violettedev.com/eecs4222/shared/auth"
	"go.violettedev.com/eecs4222/shared/database/schema"
)

var socketIdsToConnection = make(map[string]*websocket.Conn)

// Initializes coordinators to handle distributed messaging
func InitializeDistributedMessaging() {
	messaging.SetupMessagingPubSub(&socketIdsToConnection)
	user_count.SetupUserCountPubSub(&socketIdsToConnection)
}

// Store socket in memory & update user count -- requires authCtx.socketId to be set
func HandleConnectionOpened(c *websocket.Conn, authCtx *structs.AuthContext) {
	socketIdsToConnection[authCtx.SocketId] = c
	// User joined, increment count
	user_count.IncrementAnonymousUserCount()
	user_count.PublishUserCountMessage()
}

// Removes socket from in-memory store & updates user count
func HandleConnectionClosed(authCtx *structs.AuthContext) {
	if !user_count.DidExit() {
		// User left, decrement count (depending on auth status)
		delete(socketIdsToConnection, authCtx.SocketId)
		if authCtx.UserId != "" {
			user_count.DecrementAuthorizedUserCount(authCtx.UserId)
		} else {
			user_count.DecrementAnonymousUserCount()
		}
		user_count.PublishUserCountMessage()
	}
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
	userId, name, email, err := authProvider.GetAuthContextWebSocket(c, content)
	if err != nil && err.Error() == "refresh" {
		// Request refresh
		c.WriteJSON(structs.Message{
			Type: "refresh",
		})
	} else {
		// Success, set auth context and notify them that they're signed in
		authCtx.Name = name
		authCtx.Email = email
		authCtx.UserId = userId
		c.WriteJSON(structs.Message{
			Type: "signed_in",
		})
		// Update user count -- was previously counted as an anonymous user
		// and is now authorized
		user_count.DecrementAnonymousUserCount()
		user_count.IncrementAuthorizedUserCount(userId)
		user_count.PublishUserCountMessage()
	}
	// Note: We don't send anything on no auth -- as it's default behavior
}

// Sends a message to all clients -- can only be performed by logged in users
// Uses redis pub/sub to notify all server instances of message (including the
// current server)
func HandleChatMessage(c *websocket.Conn, authCtx *structs.AuthContext, subject string, content string) {
	name := authCtx.Name
	email := authCtx.Email
	if name != "" && email != "" {
		ts := time.Now()
		message := structs.ChatMessage{
			Type: "message",
			MessageSchema: schema.MessageSchema{
				ID:      primitive.NewObjectIDFromTimestamp(ts),
				Subject: subject,
				Message: content,
				Name:    name,
				Email:   email,
				Ts:      ts.UnixMilli(),
			},
		}
		// Write message to DB
		id, err := persistence.WriteMessage(message.MessageSchema)
		if err != nil {
			// Notify user of failed message
			notifyFailure(c)
			log.Print(err)
			return
		}
		message.ID = id
		// Publish message to all servers
		err = messaging.PublishMessage(message)
		if err != nil {
			// Notify user of failed message
			notifyFailure(c)
		}
	} else {
		// Notify user of failed message
		notifyFailure(c)
	}
}

// Send pong on ping recipient (heartbeat)
func HandlePing(c *websocket.Conn, content string) {
	c.WriteJSON(structs.Message{
		Type: "pong",
	})
}

// Notify connection of message that failed to send
func notifyFailure(c *websocket.Conn) {
	c.WriteJSON(structs.Message{
		Type: "message_failed",
	})
}
