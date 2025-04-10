package routes

import (
	"net/http"

	"firebase.google.com/go/v4/auth"

	"gdp8-backend/internal/handlers"
	"gdp8-backend/internal/middleware"
	"gdp8-backend/internal/services"
)

func RegisterFileRoutes(firebaseAuth *auth.Client,
	studyGroupService services.StudyGroupService,
	fileService services.FileService,
	chatbotService services.ChatBotService) {
	handler := handlers.NewFileHandler(studyGroupService, fileService, chatbotService)
	http.HandleFunc("GET /api/files/{chatID}", middleware.WithFirebaseAuth(firebaseAuth, handler.GetFiles))
	http.HandleFunc("POST /api/file/delete", middleware.WithFirebaseAuth(firebaseAuth, handler.DeleteFile))
	http.HandleFunc("POST /api/file", middleware.WithFirebaseAuth(firebaseAuth, handler.UploadFile))
}
