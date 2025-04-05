package handlers

import (
	"encoding/json"
	"log"
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
		log.Printf("Error fetching modules: %v\n", err)
		http.Error(w, "Failed to fetch modules", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, modules)
}

func (h *ModuleHandler) CreateModule(w http.ResponseWriter, r *http.Request) {
	var moduleDetails models.ModuleDetails

	if err := json.NewDecoder(r.Body).Decode(&moduleDetails); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	_, err := h.moduleService.CreateModule(moduleDetails)
	if err != nil {
		log.Printf("Error creating module: %v\n", err)
		http.Error(w, "Failed to create module", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *ModuleHandler) GetModuleStudyGroupStats(w http.ResponseWriter, _ *http.Request) {
	stats, err := h.moduleService.GetModuleStudyGroupStats()
	if err != nil {
		log.Printf("Error fetching study group stats: %v\n", err)
		http.Error(w, "Failed to fetch study group stats", http.StatusInternalServerError)
		return
	}
	sendJSONResponse(w, stats)
}
