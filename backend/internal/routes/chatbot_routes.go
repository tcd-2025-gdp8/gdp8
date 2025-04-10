package routes

import (
	"log"
	"net/http"

	"firebase.google.com/go/v4/auth"

	"gdp8-backend/internal/chatbot"
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

	client, err := chatbot.NewClient()
	if err != nil {
		log.Println(err)
		return
	}

	handler := handlers.NewChatBotHandler(
		studyGroupService,
		fileService,
		chatService,
		chatbotService,
		client)

	http.HandleFunc("POST /api/chatbot", middleware.WithFirebaseAuth(firebaseAuth, handler.Prompt))
}
