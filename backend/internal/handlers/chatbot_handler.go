package handlers

import (
	"net/http"

	"gdp8-backend/internal/chatbot"
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
	_ = chatID
	// needMemory := r.FormValue("memory")

	var resp string
	var err error

	// if needMemory == "true" {
	// TODO (SCRUM 153):
	// - query for all the context_data in the group

	// TODO: query db and add the memory in the prompt
	resp, err = h.client.Prompt(r.Context(), userMsg)
	// } else {
	// resp, err = h.client.Prompt(r.Context(), userMsg)
	// }

	if err != nil {
		http.Error(w, "Couldn't generate the response", http.StatusInternalServerError)
	}

	sendJSONResponse(w, map[string]string{"message": resp})
}

// func (h ChatBotHandler) isUserAllowed(_ *http.Request) bool {
// 	return true
// }
