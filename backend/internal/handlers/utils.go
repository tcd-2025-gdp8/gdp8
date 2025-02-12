package handlers

import (
	"encoding/json"
	"net/http"
)

func sendJSONResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "Error processing response", http.StatusInternalServerError)
		return
	}
}
