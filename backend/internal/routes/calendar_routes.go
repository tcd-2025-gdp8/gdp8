package routes

import (
	"log"
	"net/http"

	"firebase.google.com/go/v4/auth"

	"gdp8-backend/internal/handlers"
	"gdp8-backend/internal/middleware"
	"gdp8-backend/internal/services"
)

func RegisterCalendarRoutes(firebaseAuth *auth.Client, calendarService *services.CalendarService) {
	if calendarService == nil || calendarService.GetService() == nil {
		log.Println("⚠️ CalendarService is nil, skipping /api/calendar/invite route registration")
		return
	}

	handler := handlers.NewCalendarHandler(calendarService.GetService(), firebaseAuth)

	http.HandleFunc(
		"/api/calendar/invite",
		middleware.WithFirebaseAuth(firebaseAuth, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
				return
			}
			handler.Invite(w, r)
		}))
}
