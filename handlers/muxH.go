package handlers

import (
	"net/http"

	"github.com/rs/cors"
)

var Mux = http.NewServeMux()

func MuxHandler() http.Handler {
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "https://chat-sv.netlify.app"}, // Svelte dev server
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization" },
		AllowCredentials: true,
	}).Handler(Mux)

	return handler
}
