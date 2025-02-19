package routes

import (
	"net/http"

	"firebase.google.com/go/v4/auth"

	"gdp8-backend/internal/handlers"
	"gdp8-backend/internal/middleware"
	"gdp8-backend/internal/persistence"
	"gdp8-backend/internal/repositories"
	"gdp8-backend/internal/services"
)

func RegisterStudyGroupRoutes(firebaseAuth *auth.Client) {
	txManager := persistence.MockTransactionManager{}
	studyGroupRepo := repositories.NewMockStudyGroupRepository()
	studyGroupService := services.NewStudyGroupService(&txManager, studyGroupRepo)
	handler := handlers.NewStudyGroupHandler(studyGroupService)

	http.HandleFunc("GET /api/study-groups", middleware.WithFirebaseAuth(firebaseAuth, handler.GetAllStudyGroups))
	http.HandleFunc("GET /api/study-groups/", middleware.WithFirebaseAuth(firebaseAuth, handler.GetAllStudyGroups))
	http.HandleFunc("GET /api/study-groups/{id}", middleware.WithFirebaseAuth(firebaseAuth, handler.GetStudyGroup))
	http.HandleFunc("POST /api/study-groups", middleware.WithFirebaseAuth(firebaseAuth, handler.CreateStudyGroup))
	http.HandleFunc("POST /api/study-groups/", middleware.WithFirebaseAuth(firebaseAuth, handler.CreateStudyGroup))
	// TODO endpoint for deleting and updating study groups
	http.HandleFunc("POST /api/study-groups/{id}/{command}", middleware.WithFirebaseAuth(firebaseAuth, handler.HandleStudyMemberOperation))
	http.HandleFunc("POST /api/study-groups/{id}/{command}", middleware.WithFirebaseAuth(firebaseAuth, handler.HandleStudyMemberOperation))
}
