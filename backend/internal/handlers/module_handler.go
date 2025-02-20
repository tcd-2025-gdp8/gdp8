package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gdp8-backend/internal/models"
	"gdp8-backend/internal/services"
)

type ModuleHandler struct {
	moduleService services.ModuleService
}

func NewModuleHandler(moduleService services.ModuleService) *ModuleHandler {
	return &ModuleHandler{moduleService: moduleService}
}

func (h *ModuleHandler) GetAllModules(w http.ResponseWriter, _ *http.Request) {
	modules, err := h.moduleService.GetAllModules()
	if err != nil {
		http.Error(w, "Failed to fetch modules", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(modules); err != nil {
		http.Error(w, "Failed to encode modules", http.StatusInternalServerError)
	}
}

type ModulePreferences struct {
	IDs []string `json:"selectedModules"`
}

func (h *ModuleHandler) SaveUserModules(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return
	}
	fmt.Println("Received payload:", string(bodyBytes))
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var req ModulePreferences
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// TODO: instead of printing we would bind the modules to the user using the id
	for _, moduleID := range req.IDs {
		fmt.Println("Module ID:", moduleID)
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ModuleHandler) CreateModule(w http.ResponseWriter, r *http.Request) {
	var newModule models.Module

	if err := json.NewDecoder(r.Body).Decode(&newModule); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := h.moduleService.CreateModule(newModule.ID, newModule.Name)
	if err != nil {
		http.Error(w, "Failed to create module", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
