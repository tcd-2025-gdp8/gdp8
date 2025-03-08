package routes

import (
	"net/http"

	"firebase.google.com/go/v4/auth"

	"gdp8-backend/internal/handlers"
	"gdp8-backend/internal/middleware"
	"gdp8-backend/internal/services"
)

func RegisterFileRoutes(firebaseAuth *auth.Client, studyGroupService services.StudyGroupService) {
	handler := handlers.NewFileHandler(studyGroupService)
	// To get all files from a chat room—depending on how we change the study groups could be group id
	// do not think is worth to have an endpoint to get a single file
	http.HandleFunc("GET /api/files/{chatID}", middleware.WithFirebaseAuth(firebaseAuth, handler.GetFiles))
	http.HandleFunc("POST /api/file", middleware.WithFirebaseAuth(firebaseAuth, handler.UploadFile))
}
