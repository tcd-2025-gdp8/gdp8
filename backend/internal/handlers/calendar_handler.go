package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"firebase.google.com/go/v4/auth"

	"gdp8-backend/internal/middleware"
	"gdp8-backend/internal/services"
)

type CalendarHandler struct {
	CalendarService *services.CalendarService
	FirebaseAuth    *auth.Client
}

func NewCalendarHandler(service *services.CalendarService, firebaseAuth *auth.Client) *CalendarHandler {
	return &CalendarHandler{
		CalendarService: service,
		FirebaseAuth:    firebaseAuth,
	}
}

type InviteRequest struct {
	Summary     string   `json:"summary"`
	Description string   `json:"description"`
	Location    string   `json:"location"`
	StartTime   string   `json:"startTime"` // ISO 8601 format: "2025-04-10T15:00:00Z"
	EndTime     string   `json:"endTime"`   // ISO 8601 format
	Attendees   []string `json:"attendees"` // List of validated emails
}

// Invite handles the POST /api/calendar/invite request
func (h *CalendarHandler) Invite(w http.ResponseWriter, r *http.Request) {
	var req InviteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	uid, ok := r.Context().Value(middleware.UIDCtxKey{}).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.FirebaseAuth.GetUser(r.Context(), uid)

	if err != nil || user.Email == "" {
		http.Error(w, "Failed to get user email", http.StatusInternalServerError)
		return
	}

	input := services.CalendarInviteInput{
		Summary:     req.Summary,
		Description: req.Description,
		Location:    req.Location,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Attendees:   req.Attendees,
		Organizer:   user.Email,
	}

	err = h.CalendarService.SendInvite(input)
	if err != nil {
		http.Error(w, "Failed to send invite: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"message": "Invite sent successfully"}`)) // ✅ just assign
	if err != nil {
		log.Printf("Failed to write response: %v", err)
	}

}
