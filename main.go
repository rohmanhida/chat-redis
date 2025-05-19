package main

import (
	"context"
	"log"
	"net/http"

	"chat/handlers"
	"chat/routes"

	"github.com/gorilla/websocket"
)

// context
var ctx = context.Background()

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {

	// initialize Redis client
	rdb := handlers.InitRedis()

	// Subscribe to Redis channel
	go handlers.SubscribeRedis(rdb)

	// routes
	routes.RegisterRoutes(rdb)

	// mux handler
	httpHandler := handlers.MuxHandler()

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", httpHandler))
}
