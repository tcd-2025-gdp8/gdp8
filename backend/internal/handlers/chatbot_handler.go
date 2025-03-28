package handlers

import (
	"net/http"

	"gdp8-backend/internal/services"
)

type ChatBotHandler struct {
	studyGroupService services.StudyGroupService
	fileService       services.FileService
	chatService       services.ChatService
	chatbotService    services.ChatBotService
}

func NewChatBotHandler(
	studyGroupService services.StudyGroupService,
	fileService services.FileService,
	chatService services.ChatService,
	chatbotService services.ChatBotService,
) *ChatBotHandler {
	return &ChatBotHandler{
		studyGroupService: studyGroupService,
		fileService:       fileService,
		chatService:       chatService,
		chatbotService:    chatbotService,
	}
}

func (h ChatBotHandler) Test(_ http.ResponseWriter, _ *http.Request) {}
