package dbm

import (
	"github.com/go-redis/redis/v8"
)

var redis_client *redis.Client

func ConnectRedis() {
	redis_client = redis.NewClient(
		&redis.Options{
			Addr: "0.0.0.0:6379",
		},
	)
}

func GetRedisConnection() *redis.Client {
	return redis_client
}
