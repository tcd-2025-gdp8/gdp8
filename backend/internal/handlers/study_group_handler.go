package handlers

import (
	"errors"
	"gdp8-backend/internal/models"
	"gdp8-backend/internal/services"
	"gdp8-backend/internal/utils"
	"net/http"
)

type StudyGroupHandler struct {
	service services.StudyGroupService
}

type StudyGroupDTO struct {
	ID          models.StudyGroupID `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Type        string              `json:"type"`
}

func MapStudyGroupToDTO(studyGroup models.StudyGroup) StudyGroupDTO {
	return StudyGroupDTO{
		ID:          studyGroup.ID,
		Name:        studyGroup.Name,
		Description: studyGroup.Description,
		Type:        string(studyGroup.Type),
	}
}

func NewStudyGroupHandler(studyGroupService services.StudyGroupService) *StudyGroupHandler {
	return &StudyGroupHandler{service: studyGroupService}
}

func (h *StudyGroupHandler) GetStudyGroup(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := utils.ConvertToType[models.StudyGroupID](idString)
	if err != nil {
		http.Error(w, "Invalid study group ID", http.StatusBadRequest)
		return
	}

	studyGroup, err := h.service.GetStudyGroupByID(id)
	if err != nil {
		if errors.Is(err, services.ErrStudyGroupNotFound) {
			http.Error(w, "Study group not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error fetching study group", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, MapStudyGroupToDTO(*studyGroup))
}

func (h *StudyGroupHandler) GetAllStudyGroups(w http.ResponseWriter, r *http.Request) {
	studyGroups, err := h.service.GetAllStudyGroups()
	if err != nil {
		http.Error(w, "Error fetching study groups", http.StatusInternalServerError)
		return
	}

	dtoList := make([]StudyGroupDTO, len(studyGroups))
	for i, studyGroup := range studyGroups {
		dtoList[i] = MapStudyGroupToDTO(studyGroup)
	}

	sendJSONResponse(w, dtoList)
}
