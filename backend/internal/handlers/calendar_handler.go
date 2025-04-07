package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"gdp8-backend/internal/services"

	"firebase.google.com/go/v4/auth"
)

type CalendarHandler struct {
	Service      *services.CalendarService
	FirebaseAuth *auth.Client
}

func NewCalendarHandler(service *services.CalendarService, auth *auth.Client) *CalendarHandler {
	return &CalendarHandler{Service: service, FirebaseAuth: auth}
}

func (h *CalendarHandler) Invite(w http.ResponseWriter, r *http.Request) {
	// Extract and verify token
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing Authorization Header", http.StatusUnauthorized)
		return
	}

	tokenStr := authHeader[len("Bearer "):]
	token, err := h.FirebaseAuth.VerifyIDToken(r.Context(), tokenStr)
	if err != nil {
		http.Error(w, "Invalid Firebase Token", http.StatusUnauthorized)
		return
	}

	log.Printf("🔐 Authenticated user: %s", token.UID)

	// Decode invite
	var input services.CalendarInviteInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Send calendar invite
	if err := h.Service.SendInvite(input); err != nil {
		http.Error(w, "Failed to send invite: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(`{"message": "Invite sent successfully"}`)); err != nil {
		log.Printf("❌ Failed to write response: %v", err)
	}
}
