package db

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var RedisClient *redis.Client
var RedisCtx = context.Background()

func InitRedis() *redis.Client {
	_ = godotenv.Load()

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + "6379",
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	_, err := RedisClient.Ping(RedisCtx).Result()
	if err != nil {
		log.Fatalf("Redis接続失敗: %v", err)
	}

	return RedisClient
}