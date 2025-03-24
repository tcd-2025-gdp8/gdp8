package routes

import (
	"gdp8-backend/internal/handlers"
	"net/http"
)

func RegisterCalendarRoutes(calendarHandler *handlers.CalendarHandler) {
	http.HandleFunc("/api/calendar/invite", calendarHandler.CreateInviteHandler)
}
