package business

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
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
var messagingPubSub *redis.PubSub
var userCountPubSub *redis.PubSub

var localAnonymousUserCount = 0
var localAuthenticatedUserMap = make(map[string]int)
var isCleanup = false

// Initialized redis pub/sub for distributed messaging
func InitializeDistributedMessaging() {
	// Only initialize 1 instance
	if redisClient == nil {
		redisClient = cache.GetRedisClientSingleton()
		setupMessagingPubSub()
		setupUserCountPubSub()
		go cleanupOnExit()
	}
}

func setupMessagingPubSub() {
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

func setupUserCountPubSub() {
	// Setup pubsub
	ctx, cancel := context.WithCancel(context.Background())
	redisClient = cache.GetRedisClientSingleton()
	userCountPubSub = redisClient.Subscribe(ctx, os.Getenv("REDIS_USER_COUNT_CHANNEL"))
	// Receive message in background
	go func() {
		// Cleanup redis client
		defer userCountPubSub.Close()
		defer cancel()
		// Handle subscription
		for {
			msg, err := userCountPubSub.ReceiveMessage(ctx)
			if err != nil {
				panic(err)
			}
			// On successful message recipient, send message
			// to all connected users on this server
			userCountMessage := stringJsonToUserCountMessage(msg.Payload)
			if userCountMessage != nil {
				sendUserCountMessageToEveryone(*userCountMessage)
			}
		}
	}()
}

// Store socket in memory & update user count -- requires authCtx.socketId to be set
func HandleConnectionOpened(c *websocket.Conn, authCtx *structs.AuthContext) {
	socketIdsToConnection[authCtx.SocketId] = c
	incrementAnonymousUserCount()
	publishUserCountMessage()
}

// Removes socket from in-memory store & updates user count
func HandleConnectionClosed(authCtx *structs.AuthContext) {
	if !isCleanup {
		delete(socketIdsToConnection, authCtx.SocketId)
		if authCtx.UserId != "" {
			decrementAuthorizedUserCount(authCtx.UserId)
		} else {
			decrementAnonymousUserCount()
		}
		publishUserCountMessage()
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
		authCtx.Name = name
		authCtx.Email = email
		authCtx.UserId = userId
		c.WriteJSON(structs.Message{
			Type: "signed_in",
		})
		// Update user count -- was previously an unregistered user
		decrementAnonymousUserCount()
		incrementAuthorizedUserCount(userId)
		publishUserCountMessage()
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

func publishUserCountMessage() error {
	// Get counts
	authorizedUserCount, err := getAuthorizedUserCount()
	if err != nil {
		return err
	}
	anonymousUserCount, err := getAnonymousUserCount()
	if err != nil {
		return err
	}
	// Create user count message string
	userCountMessageString := userCountMessageToJsonString(structs.UserCountMessage{
		Type:            "user_count",
		AuthorizedUsers: authorizedUserCount,
		AnonymousUsers:  anonymousUserCount,
	})
	if userCountMessageString == "" {
		return errors.New("failed")
	}
	// Publish message
	ctx, cancel := context.WithTimeout(context.Background(), constants.DatabaseTimeout)
	defer cancel()
	err = redisClient.Publish(ctx, os.Getenv("REDIS_USER_COUNT_CHANNEL"), userCountMessageString).Err()
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

// Sends user chat message to all websockets
func sendUserCountMessageToEveryone(userCountMessage structs.UserCountMessage) {
	for _, conn := range socketIdsToConnection {
		if conn != nil {
			conn.WriteJSON(userCountMessage)
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

// Increments anonymous user count & notifies pub/sub
func incrementAnonymousUserCount() {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DatabaseTimeout)
	defer cancel()
	redisClient.Incr(ctx, os.Getenv("ANONYMOUS_USERS_REDIS_KEY"))
	// Locally increase count
	localAnonymousUserCount += 1
}

// Decrements anonymous user count & notifies pub/sub
func decrementAnonymousUserCount() {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DatabaseTimeout)
	defer cancel()
	redisClient.Decr(ctx, os.Getenv("ANONYMOUS_USERS_REDIS_KEY"))
	// Locally reduce count
	localAnonymousUserCount -= 1
}

// Increments authorized users count (by their id) & notifies pub/sub
func incrementAuthorizedUserCount(userId string) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DatabaseTimeout)
	defer cancel()
	// Increment users session count
	redisClient.HIncrBy(ctx, os.Getenv("AUTHORIZED_USERS_REDIS_KEY"), userId, 1)
	// Increment count for the user locally
	localAuthenticatedUserMap[userId] = localAuthenticatedUserMap[userId] + 1
}

// Decrements authorized users count (by their id) & notifies pub/sub
// User is only counted as removed once all their sockets are ended
func decrementAuthorizedUserCount(userId string) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DatabaseTimeout*2)
	defer cancel()
	// Remove 1 from users session count
	count, err := redisClient.HIncrBy(ctx, os.Getenv("AUTHORIZED_USERS_REDIS_KEY"), userId, -1).Result()
	if err != nil {
		log.Print(err)
	}
	// Delete user if that was the last logout
	if count <= 0 {
		redisClient.HDel(ctx, os.Getenv("AUTHORIZED_USERS_REDIS_KEY"), userId)
	}
	// De-increment count for the user locally
	localAuthenticatedUserMap[userId] = localAuthenticatedUserMap[userId] - 1
	if localAuthenticatedUserMap[userId] <= 0 {
		delete(localAuthenticatedUserMap, userId)
	}
}

// Returns anonymous user count on success. (-1, err) on failure
func getAnonymousUserCount() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DatabaseTimeout)
	defer cancel()
	// Get value from cache
	count, err := redisClient.Get(ctx, os.Getenv("ANONYMOUS_USERS_REDIS_KEY")).Int()
	if err != nil {
		return -1, nil
	}
	return count, nil

}

// Returns authorized user count on success. (-1, err) on failure
func getAuthorizedUserCount() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DatabaseTimeout)
	defer cancel()
	// Get # of authorized users keys
	count, err := redisClient.HLen(ctx, os.Getenv("AUTHORIZED_USERS_REDIS_KEY")).Result()
	if err != nil {
		return -1, nil
	}
	return int(count), nil
}

// Converts json string to user count msg. Returns nil on error.
func stringJsonToUserCountMessage(message string) *structs.UserCountMessage {
	var res structs.UserCountMessage
	err := json.Unmarshal([]byte(message), &res)
	if err != nil {
		log.Print(err)
		return nil
	}
	return &res
}

// Converts UserCountMessage to json string. Returns "" on error
func userCountMessageToJsonString(message structs.UserCountMessage) string {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Print(err)
		return ""
	}
	return string(messageBytes)
}

// Removes the local user counts from the redis cache & publishes it
func cleanupUserCount() {
	// Clean up anonymous users
	ctx, cancel := context.WithTimeout(context.Background(), constants.DatabaseTimeout)
	defer cancel()
	redisClient.DecrBy(ctx, os.Getenv("ANONYMOUS_USERS_REDIS_KEY"), int64(localAnonymousUserCount))
	// Clean up authorized users
	for userId, count := range localAuthenticatedUserMap {
		ctx, cancel := context.WithTimeout(context.Background(), constants.DatabaseTimeout*2)
		defer cancel()
		count, err := redisClient.HIncrBy(ctx, os.Getenv("AUTHORIZED_USERS_REDIS_KEY"), userId, int64(-1*count)).Result()
		if err != nil {
			log.Print(err)
		}
		// Delete user if that was the last logout
		if count <= 0 {
			redisClient.HDel(ctx, os.Getenv("AUTHORIZED_USERS_REDIS_KEY"), userId)
		}
	}
	publishUserCountMessage()
}

/*
Runs a background function to cleanup the database connection on
exit.
*/
func cleanupOnExit() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	// Run exit func in background waiting for exit signal
	go func() {
		<-sigs
		isCleanup = true
		// Close all sockets so we can cleanup in peace
		for _, conn := range socketIdsToConnection {
			if conn != nil {
				conn.Close()
			}
		}
		// Clean up user count
		cleanupUserCount()
		os.Exit(0)
	}()
}
