package routes

import (
	"net/http"

	"gdp8-backend/internal/handlers"
	"gdp8-backend/internal/repositories"
	"gdp8-backend/internal/services"
)

func RegisterStudyGroupRoutes() {
	studyGroupRepo := repositories.NewMockStudyGroupRepository()
	studyGroupService := services.NewStudyGroupService(studyGroupRepo)
	handler := handlers.NewStudyGroupHandler(studyGroupService)

	http.HandleFunc("GET /api/study-groups/{id}", handler.GetStudyGroup)
	http.HandleFunc("GET /api/study-groups", handler.GetAllStudyGroups)
	http.HandleFunc("GET /api/study-groups/", handler.GetAllStudyGroups)
}
