package routes

import (
	"net/http"

	"gdp8-backend/internal/handlers"
	"firebase.google.com/go/v4/auth"
)

func RegisterChatRoutes(firebaseAuth *auth.Client) {
	hub := handlers.NewChatHub()
	go hub.Run()
	handler := handlers.NewChatHandler(hub, firebaseAuth)
	http.HandleFunc("GET /api/chat/{chatID}", handler.ServeWs);
}
