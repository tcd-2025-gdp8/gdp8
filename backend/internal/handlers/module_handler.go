package handlers

import (
	"encoding/json"
	"net/http"

	"gdp8-backend/internal/services"
)

type ModuleHandler struct {
	moduleService services.ModuleService
}

func NewModuleHandler(moduleService services.ModuleService) *ModuleHandler {
	return &ModuleHandler{moduleService: moduleService}
}

// GetAllModules returns a list of all available modules
func (h *ModuleHandler) GetAllModules(w http.ResponseWriter, r *http.Request) {
	modules, err := h.moduleService.GetAllModules()
	if err != nil {
		http.Error(w, "Failed to fetch modules", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(modules)
}

// CreateModule handles the creation of a new module
func (h *ModuleHandler) CreateModule(w http.ResponseWriter, r *http.Request) {
	var newModule struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

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
