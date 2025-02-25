package routes

import (
	"net/http"

	"firebase.google.com/go/v4/auth"

	"gdp8-backend/internal/handlers"
	"gdp8-backend/internal/middleware"
)

func RegisterChatRoutes(firebaseAuth *auth.Client) {
	hub := handlers.NewChatHub()
	go hub.Run()
	handler := handlers.NewChatHandler(hub)
	http.HandleFunc("GET /api/chat/{chatID}", middleware.WithWebSocketAuth(firebaseAuth, handler.ServeWs))
}
