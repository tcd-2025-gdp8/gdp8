package handlers

import (
	"net/http"
	"strconv"

	"gdp8-backend/internal/chatbot"
	"gdp8-backend/internal/models"
	"gdp8-backend/internal/services"
)

type ChatBotHandler struct {
	studyGroupService services.StudyGroupService
	fileService       services.FileService
	chatService       services.ChatService
	chatbotService    services.ChatBotService
	client            chatbot.GeminiClient
}

func NewChatBotHandler(
	studyGroupService services.StudyGroupService,
	fileService services.FileService,
	chatService services.ChatService,
	chatbotService services.ChatBotService,
	client chatbot.GeminiClient,
) *ChatBotHandler {
	return &ChatBotHandler{
		studyGroupService: studyGroupService,
		fileService:       fileService,
		chatService:       chatService,
		chatbotService:    chatbotService,
		client:            client,
	}
}

func (h ChatBotHandler) Prompt(w http.ResponseWriter, r *http.Request) {
	userMsg := r.FormValue("text")
	chatID := r.FormValue("chatID")
	memType := r.FormValue("memory")
	groupID, convErr := strconv.Atoi(chatID)
	if convErr != nil {
		http.Error(w, "Invalid chatID", http.StatusBadRequest)
		return
	}
	var resp string
	if memType == "file" {
		memory, err := h.chatbotService.GetMemory(models.StudyGroupID(groupID))
		if err != nil {
			http.Error(w, "Failed to retrieve memory", http.StatusInternalServerError)
			return
		}
		resp, err = h.client.PromptWithContext(r.Context(), userMsg, memory)
		if err != nil {
			http.Error(w, "Couldn't generate the response", http.StatusInternalServerError)
			return
		}
	} else {
		var err error
		resp, err = h.client.Prompt(r.Context(), userMsg)
		if err != nil {
			http.Error(w, "Couldn't generate the response", http.StatusInternalServerError)
			return
		}
	}
	sendJSONResponse(w, map[string]string{"message": resp})
}
