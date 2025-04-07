package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/calendar/v3"
)

type CalendarHandler struct {
	FirebaseAuth    *auth.Client
	CalendarService *calendar.Service
}

func NewCalendarHandler(calendarService *calendar.Service, firebaseAuth *auth.Client) *CalendarHandler {
	return &CalendarHandler{
		FirebaseAuth:    firebaseAuth,
		CalendarService: calendarService,
	}
}

func (h *CalendarHandler) Invite(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Summary     string `json:"summary"`
		Location    string `json:"location"`
		Description string `json:"description"`
		StartTime   string `json:"startTime"`
		EndTime     string `json:"endTime"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	event := &calendar.Event{
		Summary:     req.Summary,
		Location:    req.Location,
		Description: req.Description,
		Start: &calendar.EventDateTime{
			DateTime: req.StartTime,
			TimeZone: "UTC",
		},
		End: &calendar.EventDateTime{
			DateTime: req.EndTime,
			TimeZone: "UTC",
		},
	}

	createdEvent, err := h.CalendarService.Events.Insert("primary", event).Do()
	if err != nil {
		log.Printf(" Failed to insert event: %v", err)
		http.Error(w, "Failed to insert event: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{
		"message":   "Event created successfully",
		"eventLink": createdEvent.HtmlLink,
	}); err != nil {
		log.Printf("❌ Failed to write JSON response: %v", err)
	}
}
