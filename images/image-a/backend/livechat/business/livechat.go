package business

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"

	"github.com/gofiber/websocket/v2"
	"github.com/redis/go-redis/v9"
	"go.violettedev.com/eecs4222/livechat/cache"
	"go.violettedev.com/eecs4222/livechat/structs"
	"go.violettedev.com/eecs4222/shared/auth"
	"go.violettedev.com/eecs4222/shared/constants"
)

var socketIdsToConnection = make(map[string]*websocket.Conn)
var redisClient *redis.Client
var pubsub *redis.PubSub

const ANONYMOUS_USERS_REDIS_KEY = "anonymousUsers"
const AUTHORIZED_USERS_REDIS_KEY = "authorizedUsers"

// Initialized redis pub/sub for distributed messaging
func InitializeDistributedMessaging() {
	// Only initialize 1 instance
	if pubsub == nil {
		// Setup pubsub
		ctx, cancel := context.WithCancel(context.Background())
		redisClient = cache.GetRedisClientSingleton()
		pubsub = redisClient.Subscribe(ctx, os.Getenv("REDIS_MESSAGING_CHANNEL"))
		// Receive message in background
		go func() {
			// Cleanup redis client
			defer pubsub.Close()
			defer cancel()
			// Handle subscription
			for {
				msg, err := pubsub.ReceiveMessage(ctx)
				log.Print(msg)
				if err != nil {
					panic(err)
				}
				// On successful message recipient, send message
				// to all connected users on this server
				chatMessage := stringJsonToChatMessage(msg.Payload)
				if chatMessage != nil {
					sendChatMessageToEveryone(*chatMessage)
				}
			}
		}()
	}
}

// Store socket in memory -- required authCtx.socketId to be set
func HandleConnectionOpened(c *websocket.Conn, authCtx *structs.AuthContext) {
	socketIdsToConnection[authCtx.SocketId] = c
}

// Removes socket from in-memory store
func HandleConnectionClosed(authCtx *structs.AuthContext) {
	delete(socketIdsToConnection, authCtx.SocketId)
	// TODO: Notify redis
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
// Uses redis pub/sub to notify all server instances of message (including the
// current server)
func HandleChatMessage(c *websocket.Conn, authCtx *structs.AuthContext, subject string, content string) {
	name := authCtx.Name
	email := authCtx.Email
	if name != "" && email != "" {
		// TODO: Write message to db
		message := structs.ChatMessage{
			Type:    "message",
			Subject: subject,
			Message: content,
			Name:    name,
			Email:   email,
			Ts:      time.Now().UnixMilli(),
		}
		// Publish message to all servers
		publishMessage(message)
	}
}

// Send pong on ping recipient (heartbeat)
func HandlePing(c *websocket.Conn, content string) {
	c.WriteJSON(structs.Message{
		Type: "pong",
	})
}

// Publishes message to pubsub (all other instances of this service)
// Panics if redis publish fails
// Doesn't publish if converting json to string fails
func publishMessage(message structs.ChatMessage) error {
	messageString := chatMessageToJsonString(message)
	if messageString == "" {
		return errors.New("failed")
	}
	ctx, cancel := context.WithTimeout(context.Background(), constants.DatabaseTimeout)
	defer cancel()
	err := redisClient.Publish(ctx, os.Getenv("REDIS_MESSAGING_CHANNEL"), messageString).Err()
	if err != nil {
		panic(err)
	}
	return nil
}

// Send a chat message to all websockets
func sendChatMessageToEveryone(message structs.ChatMessage) {
	for _, conn := range socketIdsToConnection {
		if conn != nil {
			conn.WriteJSON(message)
		}
	}
}

// Converts json string to chat message. Returns nil on error.
func stringJsonToChatMessage(message string) *structs.ChatMessage {
	var res structs.ChatMessage
	err := json.Unmarshal([]byte(message), &res)
	if err != nil {
		log.Print(err)
		return nil
	}
	return &res
}

// Converts chat message to json string. Returns "" on error
func chatMessageToJsonString(message structs.ChatMessage) string {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Print(err)
		return ""
	}
	return string(messageBytes)
}
