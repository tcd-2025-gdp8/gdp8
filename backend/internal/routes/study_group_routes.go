package routes

import (
	"net/http"

	"firebase.google.com/go/v4/auth"

	"gdp8-backend/internal/handlers"
	"gdp8-backend/internal/middleware"
	"gdp8-backend/internal/services"
)

func RegisterStudyGroupRoutes(firebaseAuth *auth.Client,
	studyGroupService services.StudyGroupService, studySessionService services.StudySessionService) {
	handler := handlers.NewStudyGroupHandler(studyGroupService)
	studySessionHandler := handlers.NewStudySessionHandler(studySessionService)

	// Study Groups CRUD
	http.HandleFunc("GET /api/study-groups", middleware.WithFirebaseAuth(firebaseAuth, handler.GetRelevantStudyGroups))
	http.HandleFunc("GET /api/study-groups/", middleware.WithFirebaseAuth(firebaseAuth, handler.GetRelevantStudyGroups))
	http.HandleFunc("GET /api/study-groups/{id}", middleware.WithFirebaseAuth(firebaseAuth, handler.GetStudyGroup))
	http.HandleFunc("POST /api/study-groups", middleware.WithFirebaseAuth(firebaseAuth, handler.CreateStudyGroup))
	http.HandleFunc("POST /api/study-groups/", middleware.WithFirebaseAuth(firebaseAuth, handler.CreateStudyGroup))
	// TODO endpoint for deleting and updating study groups

	// Study Groups commands
	http.HandleFunc("POST /api/study-groups/{id}/{command}", middleware.WithFirebaseAuth(firebaseAuth,
		handler.HandleStudyMemberOperation))

	// Study Sessions CRUD
	http.HandleFunc("GET /api/study-groups/{groupID}/study-sessions", middleware.WithFirebaseAuth(firebaseAuth,
		studySessionHandler.GetStudySessionsByGroup))
	http.HandleFunc("POST /api/study-groups/{groupID}/study-sessions", middleware.WithFirebaseAuth(firebaseAuth,
		studySessionHandler.CreateStudySession))
	http.HandleFunc("PUT /api/study-groups/{groupID}/study-sessions/{studySessionID}",
		middleware.WithFirebaseAuth(firebaseAuth, studySessionHandler.UpdateStudySession))
	http.HandleFunc("DELETE /api/study-groups/{groupID}/study-sessions/{studySessionID}",
		middleware.WithFirebaseAuth(firebaseAuth, studySessionHandler.DeleteStudySession))

	// Study SessionAvailability Requests handlers
	http.HandleFunc("GET /api/study-groups/{groupID}/availability-requests",
		middleware.WithFirebaseAuth(firebaseAuth, studySessionHandler.GetCurrentAvailabilityRequests))
	http.HandleFunc("POST /api/study-groups/{groupID}/availability-requests",
		middleware.WithFirebaseAuth(firebaseAuth, studySessionHandler.CreateAvailabilityRequest))
	http.HandleFunc("DELETE /api/study-groups/{groupID}/availability-requests/{availabilityRequestID}",
		middleware.WithFirebaseAuth(firebaseAuth, studySessionHandler.DeleteAvailabilityRequest))
	http.HandleFunc("PUT /api/study-groups/{groupID}/availability-requests/{availabilityRequestID}/entries",
		middleware.WithFirebaseAuth(firebaseAuth, studySessionHandler.UpsertAvailabilityEntries))
}
