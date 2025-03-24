package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"google.golang.org/api/calendar/v3"

	"gdp8-backend/internal/services"
)

type CalendarHandler struct {
	calendarService services.GoogleCalendarService
}

func NewCalendarHandler(calendarService services.GoogleCalendarService) *CalendarHandler {
	return &CalendarHandler{calendarService: calendarService}
}

// CreateInviteHandler handles HTTP requests for creating a Google Calendar invite.
func (h *CalendarHandler) CreateInviteHandler(w http.ResponseWriter, r *http.Request) {
	// Log that we received a request (for debugging)
	// (You can remove this once everything works.)
	// log.Println("🔥 Received request at /api/calendar/invite")

	var req struct {
		Summary        string   `json:"summary"`
		Location       string   `json:"location"`
		Description    string   `json:"description"`
		StartTime      string   `json:"startTime"` // ISO8601 formatted
		EndTime        string   `json:"endTime"`   // ISO8601 formatted
		AttendeeEmails []string `json:"attendeeEmails"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		http.Error(w, "Invalid startTime", http.StatusBadRequest)
		return
	}
	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		http.Error(w, "Invalid endTime", http.StatusBadRequest)
		return
	}

	event := &calendar.Event{
		Summary:     req.Summary,
		Location:    req.Location,
		Description: req.Description,
		Start: &calendar.EventDateTime{
			DateTime: startTime.Format(time.RFC3339),
			TimeZone: "UTC", // Adjust as needed.
		},
		End: &calendar.EventDateTime{
			DateTime: endTime.Format(time.RFC3339),
			TimeZone: "UTC",
		},
		Attendees: func(emails []string) []*calendar.EventAttendee {
			attendees := make([]*calendar.EventAttendee, len(emails))
			for i, email := range emails {
				attendees[i] = &calendar.EventAttendee{Email: email}
			}
			return attendees
		}(req.AttendeeEmails),
		Reminders: &calendar.EventReminders{
			UseDefault: true,
		},
	}

	calendarID := "primary"
	createdEvent, err := h.calendarService.CreateEvent(calendarID, event)
	if err != nil {
		http.Error(w, "Failed to create calendar event: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(createdEvent); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
