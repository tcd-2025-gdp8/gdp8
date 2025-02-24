package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"gdp8-backend/internal/middleware"
	"gdp8-backend/internal/models"
)

func getUserID(r *http.Request) (models.UserID, error) {
	uid, ok := r.Context().Value(middleware.UIDCtxKey{}).(string)
	if !ok || uid == "" {
		return "", errors.New("unauthorized")
	}
	return models.UserID(uid), nil
}

func sendJSONResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "Error processing response", http.StatusInternalServerError)
		return
	}
}
