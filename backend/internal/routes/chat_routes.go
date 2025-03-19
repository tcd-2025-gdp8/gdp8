package routes

import (
	"net/http"

	"firebase.google.com/go/v4/auth"

	"gdp8-backend/internal/handlers"
	"gdp8-backend/internal/middleware"
	"gdp8-backend/internal/services"
)

func RegisterChatRoutes(firebaseAuth *auth.Client, chatService services.ChatService) {
	hub := handlers.NewChatHub()
	go hub.Run()

	handler := handlers.NewChatHandler(hub, chatService)

	http.HandleFunc("GET /api/chat/{chatID}", middleware.WithWebSocketAuth(firebaseAuth, handler.ServeWs))
	http.HandleFunc("GET /api/chat/{chatID}/past", middleware.WithFirebaseAuth(firebaseAuth, handler.GetStudyGroup))
}
