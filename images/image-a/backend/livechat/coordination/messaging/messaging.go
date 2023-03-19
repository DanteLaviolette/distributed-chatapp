package messaging

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/gofiber/websocket/v2"
	"github.com/redis/go-redis/v9"
	"go.violettedev.com/eecs4222/livechat/cache"
	"go.violettedev.com/eecs4222/livechat/structs"
	"go.violettedev.com/eecs4222/shared/constants"
)

var redisClient *redis.Client
var messagingPubSub *redis.PubSub
var socketIdsToConnection *(map[string]*websocket.Conn)

/*
Sets up messaging coordination between other instances.
connectionMap is a pointer to a map containing [user_id, ws_connection]
Must be called before any other functions in this class.
*/
func SetupMessagingPubSub(connectionMap *map[string]*websocket.Conn) {
	if redisClient == nil {
		// Save local variables
		redisClient = cache.GetRedisClientSingleton()
		socketIdsToConnection = connectionMap
		// Setup pubsub
		ctx, cancel := context.WithCancel(context.Background())
		messagingPubSub = redisClient.Subscribe(ctx, os.Getenv("REDIS_MESSAGING_CHANNEL"))
		// Receive message in background
		go func() {
			// Cleanup redis client
			defer messagingPubSub.Close()
			defer cancel()
			// Handle subscription
			for {
				msg, err := messagingPubSub.ReceiveMessage(ctx)
				if err != nil {
					log.Print(err)
					continue
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

// Publishes message all instances of this service (including this one)
// Returns error on failure
func PublishMessage(message structs.ChatMessage) error {
	messageString := chatMessageToJsonString(message)
	if messageString == "" {
		return errors.New("failed")
	}
	ctx, cancel := context.WithTimeout(context.Background(), constants.DatabaseTimeout)
	defer cancel()
	err := redisClient.Publish(ctx, os.Getenv("REDIS_MESSAGING_CHANNEL"), messageString).Err()
	if err != nil {
		return err
	}
	return nil
}

// Send a chat message to all websockets
func sendChatMessageToEveryone(message structs.ChatMessage) {
	for _, conn := range *socketIdsToConnection {
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
