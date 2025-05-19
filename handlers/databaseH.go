package handlers

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func InitRedis() *redis.Client {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get Redis URL
	redisURL := os.Getenv("REDIS_URI")
	if redisURL == "" {
		log.Fatal("REDIS_URI not set")
	}

	// Parse the URL into Redis options
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("Failed to parse Redis URI: %v", err)
	}

	// Create Redis client
	rdb := redis.NewClient(opt)

	// Ping test
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis")

	return rdb
}
