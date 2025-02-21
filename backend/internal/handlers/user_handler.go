package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"gdp8-backend/internal/models"
	"gdp8-backend/internal/repositories"
	"gdp8-backend/internal/services"
)

type UserHandler struct {
	userService services.UserService
	moduleRepo  repositories.ModuleRepository
}

func NewUserHandler(userService services.UserService, moduleRepo repositories.ModuleRepository) *UserHandler {
	return &UserHandler{userService: userService, moduleRepo: moduleRepo}
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "User id required", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUser(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching user: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode user: %v", err), http.StatusInternalServerError)
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	createdUser, err := h.userService.CreateUser(user)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating user: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdUser); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode user: %v", err), http.StatusInternalServerError)
	}
}

type ModulePreferences struct {
	IDs []string `json:"selectedModules"`
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

	fmt.Println("Received module IDs:", prefs.IDs)

	modules := make([]models.Module, len(prefs.IDs))
	for i, modID := range prefs.IDs {
		mod, err := h.moduleRepo.GetModuleByID(nil, modID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Module %s not found: %v", modID, err), http.StatusBadRequest)
			return
		}
		modules[i] = mod
	}

	if err := h.userService.SetModules(id, modules); err != nil {
		http.Error(w, fmt.Sprintf("Error setting modules: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
