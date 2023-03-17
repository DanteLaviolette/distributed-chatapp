package cache

import (
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

/*
Returns the redis client singleton
*/
func GetRedisClientSingleton() *redis.Client {
	if redisClient == nil {
		redisDb, err := strconv.Atoi(os.Getenv("REDIS_DB"))
		if err != nil {
			log.Fatal(err)
		}
		redisClient = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_URL"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       redisDb,
		})
	}
	return redisClient
}
