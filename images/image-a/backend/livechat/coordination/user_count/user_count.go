package user_count

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"

	"github.com/redis/go-redis/v9"
	"go.violettedev.com/eecs4222/livechat/cache"
	"go.violettedev.com/eecs4222/livechat/coordination"
	"go.violettedev.com/eecs4222/livechat/structs"
	"go.violettedev.com/eecs4222/shared/constants"
)

var userCountPubSub *redis.PubSub

var redisClient *redis.Client

var localAnonymousUserCount = atomic.Int64{}

// TODO: Not thread safe
var localAuthenticatedUserMap = make(map[string]int)
var localAuthenticatedUserMutex = sync.Mutex{}

var didExit = false

/*
Sets up user count coordination between other instances.
connectionMap is a pointer to a map containing [user_id, ws_connection]
Must be called before any other functions in this class.
*/
func SetupUserCountPubSub() {
	// Only run once
	if redisClient == nil {
		// Set local vars
		redisClient = cache.GetRedisClientSingleton()
		// Run cleanup code on exit
		go cleanupOnExit()
		// Setup pubsub
		ctx, cancel := context.WithCancel(context.Background())
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
					log.Print(err)
					continue
				}
				// On successful message recipient, send message
				// to all connected users on this server
				userCountMessage := stringJsonToUserCountMessage(msg.Payload)
				if userCountMessage != nil {
					coordination.MessageEveryone(userCountMessage)
				}
			}
		}()
	}
}

// Publishes the updated user count message (for this & other servers to consume)
// and then notify users of.
func PublishUserCountMessage() error {
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
		log.Print(err)
		return err
	}
	return nil
}

// Increments anonymous user count & notifies pub/sub
func IncrementAnonymousUserCount() {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DatabaseTimeout)
	defer cancel()
	err := redisClient.Incr(ctx, os.Getenv("ANONYMOUS_USERS_REDIS_KEY")).Err()
	if err == nil {
		// Locally increase count
		localAnonymousUserCount.Add(1)
	}
}

// Decrements anonymous user count & notifies pub/sub
func DecrementAnonymousUserCount() {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DatabaseTimeout)
	defer cancel()
	err := redisClient.Decr(ctx, os.Getenv("ANONYMOUS_USERS_REDIS_KEY")).Err()
	if err == nil {
		// Locally reduce count
		localAnonymousUserCount.Add(-1)
	}
}

// Increments authorized users count (by their id) & notifies pub/sub
func IncrementAuthorizedUserCount(userId string) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DatabaseTimeout)
	defer cancel()
	// Increment users session count
	err := redisClient.HIncrBy(ctx, os.Getenv("AUTHORIZED_USERS_REDIS_KEY"), userId, 1).Err()
	if err == nil {
		// Increment count for the user locally
		localAuthenticatedUserMutex.Lock()
		localAuthenticatedUserMap[userId] = localAuthenticatedUserMap[userId] + 1
		localAuthenticatedUserMutex.Unlock()
	}
}

// Decrements authorized users count (by their id) & notifies pub/sub
// User is only counted as removed once all their sockets are ended
func DecrementAuthorizedUserCount(userId string) {
	localAuthenticatedUserMutex.Lock()
	defer localAuthenticatedUserMutex.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), constants.DatabaseTimeout*2)
	defer cancel()
	// Remove 1 from users session count
	count, err := redisClient.HIncrBy(ctx, os.Getenv("AUTHORIZED_USERS_REDIS_KEY"), userId, -1).Result()
	if err != nil {
		log.Print(err)
		return
	} else {
		// De-increment count for the user locally
		localAuthenticatedUserMap[userId] = localAuthenticatedUserMap[userId] - 1
	}
	// Delete user if that was the last logout
	if count <= 0 {
		redisClient.HDel(ctx, os.Getenv("AUTHORIZED_USERS_REDIS_KEY"), userId)
	}
	if localAuthenticatedUserMap[userId] <= 0 {
		delete(localAuthenticatedUserMap, userId)
	}
}

// Returns true if the service is exiting, false otherwise
func DidExit() bool {
	return didExit
}

// Returns anonymous user count on success. (-1, err) on failure
func getAnonymousUserCount() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DatabaseTimeout)
	defer cancel()
	// Get value from cache
	count, err := redisClient.Get(ctx, os.Getenv("ANONYMOUS_USERS_REDIS_KEY")).Int()
	if err != nil {
		return -1, err
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
		return -1, err
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
	redisClient.DecrBy(ctx, os.Getenv("ANONYMOUS_USERS_REDIS_KEY"), localAnonymousUserCount.Load())
	// Clean up authorized users
	localAuthenticatedUserMutex.Lock()
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
	localAuthenticatedUserMutex.Unlock()
	PublishUserCountMessage()
}

/*
Runs a background function to cleanup user counts on exit.
*/
func cleanupOnExit() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	// Run exit func in background waiting for exit signal
	go func() {
		<-sigs
		didExit = true
		// Close all sockets so we can cleanup in peace
		coordination.Cleanup()
		// Clean up user count
		cleanupUserCount()
		os.Exit(0)
	}()
}
