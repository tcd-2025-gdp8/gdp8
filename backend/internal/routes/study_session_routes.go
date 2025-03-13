package routes

import (
	"net/http"

	"firebase.google.com/go/v4/auth"

	"gdp8-backend/internal/handlers"
	"gdp8-backend/internal/middleware"
	"gdp8-backend/internal/services"
)

func RegisterStudySessionRoutes(firebaseAuth *auth.Client, studySessionService services.StudySessionService) {
	studySessionHandler := handlers.NewStudySessionHandler(studySessionService)

	http.HandleFunc("GET /api/study-groups/{groupID}/study-sessions", middleware.WithFirebaseAuth(firebaseAuth,
		studySessionHandler.GetStudySessionsByGroup))
	http.HandleFunc("POST /api/study-groups/{groupID}/study-sessions", middleware.WithFirebaseAuth(firebaseAuth,
		studySessionHandler.CreateStudySession))
	http.HandleFunc("PUT /api/study-sessions/{studySessionID}",
		middleware.WithFirebaseAuth(firebaseAuth, studySessionHandler.UpdateStudySession))
	http.HandleFunc("DELETE /api/study-sessions/{studySessionID}",
		middleware.WithFirebaseAuth(firebaseAuth, studySessionHandler.DeleteStudySession))

	http.HandleFunc("GET /api/study-groups/{groupID}/availability-requests",
		middleware.WithFirebaseAuth(firebaseAuth, studySessionHandler.GetCurrentAvailabilityRequests))
	http.HandleFunc("POST /api/study-groups/{groupID}/availability-requests",
		middleware.WithFirebaseAuth(firebaseAuth, studySessionHandler.CreateAvailabilityRequest))
	http.HandleFunc("DELETE /api/availability-requests/{availabilityRequestID}",
		middleware.WithFirebaseAuth(firebaseAuth, studySessionHandler.DeleteAvailabilityRequest))
	http.HandleFunc("PUT /api/availability-requests/{availabilityRequestID}/entries",
		middleware.WithFirebaseAuth(firebaseAuth, studySessionHandler.UpsertAvailabilityEntries))
}
