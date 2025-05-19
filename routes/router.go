package routes

import (
	"encoding/json"
	"net/http"

	"chat/handlers"

	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(rdb *redis.Client) {
	// /ws routes
	handlers.Mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleWebSocket(w, r, rdb)
	})
	// /stats routes
	handlers.Mux.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		count := handlers.GetClientCount(w, r)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{"connected_clients": count})
	})
}
