package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"gdp8-backend/internal/models"
	"gdp8-backend/internal/services"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "User id missing in URL", http.StatusBadRequest)
		return
	}
	id := parts[3]

	user, err := h.userService.GetUser(models.UserID(id))
	if err != nil {
		log.Printf("Error fetching user: %v\n", err)
		http.Error(w, fmt.Sprintf("Error fetching user: %v", err), http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, user)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	createdUser, err := h.userService.CreateUser(user.ID, user.UserDetails)
	if err != nil {
		log.Printf("Error creating user: %v\n", err)
		http.Error(w, fmt.Sprintf("Error creating user: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdUser); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode user: %v", err), http.StatusInternalServerError)
	}
}

type ModulePreferences struct {
	IDs []models.ModuleID `json:"selectedModules"`
}

func (h *UserHandler) SetModules(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		http.Error(w, "User id missing in URL", http.StatusBadRequest)
		return
	}
	id := parts[3]

	var prefs ModulePreferences
	if err := json.NewDecoder(r.Body).Decode(&prefs); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.userService.SetModules(models.UserID(id), prefs.IDs); err != nil {
		log.Printf("Error setting modules: %v\n", err)
		http.Error(w, fmt.Sprintf("Error setting modules: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) GetModules(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		http.Error(w, "User id missing in URL", http.StatusBadRequest)
		return
	}
	id := parts[3]

	user, err := h.userService.GetUser(models.UserID(id))
	if err != nil {
		log.Printf("Error fetching user: %v\n", err)
		http.Error(w, fmt.Sprintf("Error fetching user: %v", err), http.StatusInternalServerError)
		return
	}
	modules := user.Modules

	sendJSONResponse(w, modules)
}
