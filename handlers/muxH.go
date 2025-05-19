package handlers

import (
	"net/http"

	"github.com/rs/cors"
)

var Mux = http.NewServeMux()

func MuxHandler() http.Handler {

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // Svelte dev server
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler(Mux)

	return handler
}
