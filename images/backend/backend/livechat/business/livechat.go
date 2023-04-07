package business

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/websocket/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.violettedev.com/eecs4222/livechat/business/coordination"
	"go.violettedev.com/eecs4222/livechat/business/coordination/messaging"
	"go.violettedev.com/eecs4222/livechat/business/coordination/user_count"
	"go.violettedev.com/eecs4222/livechat/persistence"
	"go.violettedev.com/eecs4222/livechat/structs"
	"go.violettedev.com/eecs4222/shared/auth"
	"go.violettedev.com/eecs4222/shared/database/schema"
)

// Initializes coordinators to handle distributed messaging
func InitializeDistributedMessaging() {
	coordination.InitializeThreadSafeSocketHandling()
	messaging.SetupMessagingPubSub()
	user_count.SetupUserCountPubSub()
}

// Store socket in memory & update user count -- requires authCtx.socketId to be set
func HandleConnectionOpened(c *websocket.Conn, authCtx *structs.AuthContext) {
	coordination.AddConnection(authCtx.SocketId, c)
	// User joined, increment count
	user_count.IncrementAnonymousUserCount()
	user_count.PublishUserCountMessage()
}

// Removes socket from in-memory store & updates user count
func HandleConnectionClosed(authCtx *structs.AuthContext) {
	if !user_count.DidExit() {
		// User left, decrement count (depending on auth status)
		coordination.RemoveConnection(authCtx.SocketId)
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
func HandleAuthMessage(authCtx *structs.AuthContext, content string) {
	authProvider := auth.Initialize(os.Getenv("AUTH_PRIVATE_KEY"),
		os.Getenv("REFRESH_PRIVATE_KEY"))
	userId, name, email, err := authProvider.GetAuthContextWebSocket(content)
	if err != nil && err.Error() == "refresh" {
		// Request refresh
		coordination.WriteMessage(authCtx.SocketId, structs.Message{
			Type: "refresh",
		})
	} else {
		// Success, set auth context and notify them that they're signed in
		authCtx.Name = name
		authCtx.Email = email
		authCtx.UserId = userId
		coordination.WriteMessage(authCtx.SocketId, structs.Message{
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
func HandleChatMessage(authCtx *structs.AuthContext, subject string, content string) {
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
				Ts:      ts.UnixNano(),
			},
		}
		// Send message in sub-routine so we don't block the sockets thread
		go func() {
			// Publish message to all servers first
			err := messaging.PublishMessage(message)
			if err != nil {
				// Notify user of failed message
				notifyFailure(authCtx)
				log.Print(err)
				return
			}
			// Write message to DB -- ensures consistency on new page loads
			_, err = persistence.WriteMessage(message.MessageSchema)
			if err != nil {
				// Notify user of failed message
				notifyFailure(authCtx)
				log.Print(err)
				return
			}
		}()
	} else {
		// Notify user of failed message
		notifyFailure(authCtx)
	}
}

// Send pong on ping recipient (heartbeat)
func HandlePing(authCtx *structs.AuthContext, content string) {
	coordination.WriteMessage(authCtx.SocketId, structs.Message{
		Type: "pong",
	})
}

// Notify connection of message that failed to send
func notifyFailure(authCtx *structs.AuthContext) {
	coordination.WriteMessage(authCtx.SocketId, structs.Message{
		Type: "message_failed",
	})
}
