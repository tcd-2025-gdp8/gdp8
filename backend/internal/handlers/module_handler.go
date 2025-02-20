package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

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

func (h *ModuleHandler) SaveUserModules(w http.ResponseWriter, r *http.Request) {
	fmt.Println("YESS")
	var userModules []string
	if err := json.NewDecoder(r.Body).Decode(&userModules); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// TODO: bind module with user
	// requires the actual DB integration

	for _, module := range userModules {
		fmt.Println(module)

	}

	w.WriteHeader(http.StatusOK)
}

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
