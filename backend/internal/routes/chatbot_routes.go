package routes

import (
	"net/http"

	"firebase.google.com/go/v4/auth"

	"gdp8-backend/internal/handlers"
	"gdp8-backend/internal/middleware"
	"gdp8-backend/internal/services"
)

func RegisterChatBotRoutes(
	firebaseAuth *auth.Client,
	studyGroupService services.StudyGroupService,
	fileService services.FileService,
	chatService services.ChatService,
	chatbotService services.ChatBotService) {
	handler := handlers.NewChatBotHandler(
		studyGroupService,
		fileService,
		chatService,
		chatbotService)

	http.HandleFunc("POST /api/chatbot/{chatID}", middleware.WithWebSocketAuth(firebaseAuth, handler.Test))
}
