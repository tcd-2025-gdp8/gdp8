package routes

import (
	"net/http"

	"gdp8-backend/internal/firebase"
	"gdp8-backend/internal/handlers"
)

func RegisterCalendarRoutes(firebaseAuth *firebase.FirebaseAuth, calendarHandler *handlers.CalendarHandler) {
	http.HandleFunc("/api/calendar/invite", calendarHandler.CreateInviteHandler)
}
