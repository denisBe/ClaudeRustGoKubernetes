package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
)

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
		log.Println("REDIS_ADDR environment variable is not set, using default value: " + redisAddr)
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	defer redisClient.Close()

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis at %s: %v", redisAddr, err)
	}
	log.Printf("Connected to Redis at %s", redisAddr)

	mux := http.NewServeMux()
	registerRoutes(mux, redisClient)

	fmt.Println("Server starting on :8081")
	err = http.ListenAndServe(":8081", mux)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}
