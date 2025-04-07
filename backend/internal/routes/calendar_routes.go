package routes

import (
	"net/http"

	"gdp8-backend/internal/handlers"
	"gdp8-backend/internal/middleware"
	"gdp8-backend/internal/services"

	"firebase.google.com/go/v4/auth"
)

func RegisterCalendarRoutes(firebaseAuth *auth.Client, calendarService *services.CalendarService) {
	handler := handlers.NewCalendarHandler(calendarService, firebaseAuth)

	http.HandleFunc("/api/calendar/invite", middleware.WithFirebaseAuth(firebaseAuth, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.Invite(w, r)
	}))
}
