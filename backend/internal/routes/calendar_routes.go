package routes

import (
	"log"
	"net/http"

	"gdp8-backend/internal/handlers"
)

func RegisterCalendarRoutes(calendarHandler *handlers.CalendarHandler) {
	http.HandleFunc("/api/calendar/invite", calendarHandler.CreateInviteHandler)
	log.Println("✅ Registered /api/calendar/invite route")
}
