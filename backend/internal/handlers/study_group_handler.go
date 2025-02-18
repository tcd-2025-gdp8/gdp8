package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"gdp8-backend/internal/models"
	"gdp8-backend/internal/services"
	"gdp8-backend/internal/utils"
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

type StudyGroupCreateDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
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

func (h *StudyGroupHandler) GetAllStudyGroups(w http.ResponseWriter, _ *http.Request) {
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

func (h *StudyGroupHandler) CreateStudyGroup(w http.ResponseWriter, r *http.Request) {
	var createDTO StudyGroupCreateDTO
	if err := json.NewDecoder(r.Body).Decode(&createDTO); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	uid, ok := ctx.Value("uid").(string)
	if !ok || uid == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// TODO implement validation (group type in particular)
	studyGroupDetails := models.StudyGroupDetails{
		Name:        createDTO.Name,
		Description: createDTO.Description,
		Type:        models.StudyGroupType(createDTO.Type),
	}

	createdStudyGroup, err := h.service.CreateStudyGroup(studyGroupDetails, models.UserID(uid))
	if err != nil {
		http.Error(w, "Error creating study group", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, MapStudyGroupToDTO(*createdStudyGroup))
}
