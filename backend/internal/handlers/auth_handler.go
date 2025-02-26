package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
)

// AuthHandler holds the Firebase Auth client.
type AuthHandler struct {
	firebaseAuth *auth.Client
}

// NewAuthHandler is a constructor for AuthHandler.
func NewAuthHandler(firebaseAuth *auth.Client) *AuthHandler {
	return &AuthHandler{firebaseAuth: firebaseAuth}
}

// VerifyHandler checks the Authorization Bearer token.
func (h *AuthHandler) VerifyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Verify the token with Firebase Admin
	decodedToken, err := h.firebaseAuth.VerifyIDToken(r.Context(), tokenString)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// For demo, just return a JSON response with the UID
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"uid":    decodedToken.UID,
	}); err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}
}
