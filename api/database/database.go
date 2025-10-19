package database

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

func CreateRedisClient(dbNo int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "db:6379",
		Password: "",
		DB:       dbNo,
	})
}
