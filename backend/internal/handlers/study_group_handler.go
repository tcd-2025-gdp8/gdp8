package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"gdp8-backend/internal/middleware"
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

type MemberAdminOperationDTO struct {
	TargetUserID string `json:"targetUserId"`
}

var ErrInvalidRequestPayload = errors.New("invalid request payload")
var ErrInvalidCommand = errors.New("invalid command")

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
	uid, ok := ctx.Value(middleware.UIDCtxKey{}).(string)
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

func (h *StudyGroupHandler) HandleStudyMemberOperation(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	studyGroupID, err := utils.ConvertToType[models.StudyGroupID](idString)
	if err != nil {
		http.Error(w, "Invalid study group ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	uid, ok := ctx.Value(middleware.UIDCtxKey{}).(string)
	if !ok || uid == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := models.UserID(uid)

	var memberOperationDetails MemberAdminOperationDTO
	_ = json.NewDecoder(r.Body).Decode(&memberOperationDetails)
	targetUserID := models.UserID(memberOperationDetails.TargetUserID)

	command := r.PathValue("command")
	err = h.handleCommand(command, studyGroupID, userID, targetUserID)

	switch {
	case err == nil:
		w.WriteHeader(http.StatusOK)
	case errors.Is(err, ErrInvalidRequestPayload):
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
	case errors.Is(err, services.ErrStudyGroupNotFound):
		http.Error(w, "Study group not found", http.StatusNotFound)
	case errors.Is(err, services.ErrUnauthorizedMemberOperation):
		http.Error(w, "Unauthorized study group operation", http.StatusForbidden)
	case errors.Is(err, services.ErrInvalidMemberOperation):
		http.Error(w, "Invalid study group operation", http.StatusBadRequest)
	default:
		http.Error(w, "Error processing study group operation", http.StatusInternalServerError)
	}
}

func (h *StudyGroupHandler) handleCommand(command string, studyGroupID models.StudyGroupID,
	userID models.UserID, targetUserID models.UserID) error {
	switch command {
	case "accept-invite":
		return h.service.HandleSelfMemberOperation(services.AcceptStudyGroupInviteCommand, studyGroupID, userID)
	case "reject-invite":
		return h.service.HandleSelfMemberOperation(services.RejectStudyGroupInviteCommand, studyGroupID, userID)
	case "request-to-join":
		return h.service.HandleSelfMemberOperation(services.RequestToJoinStudyGroupCommand, studyGroupID, userID)
	case "leave":
		return h.service.HandleSelfMemberOperation(services.LeaveStudyGroupCommand, studyGroupID, userID)
	case "invite":
		if targetUserID == "" {
			return ErrInvalidRequestPayload
		}
		return h.service.HandleAdminMemberOperation(
			services.InviteMemberToStudyGroupCommand, studyGroupID, targetUserID, userID)
	case "accept-request-to-join":
		if targetUserID == "" {
			return ErrInvalidRequestPayload
		}
		return h.service.HandleAdminMemberOperation(
			services.AcceptRequestToJoinStudyGroupCommand, studyGroupID, targetUserID, userID)
	case "reject-request-to-join":
		if targetUserID == "" {
			return ErrInvalidRequestPayload
		}
		return h.service.HandleAdminMemberOperation(
			services.RejectRequestToJoinStudyGroupCommand, studyGroupID, targetUserID, userID)
	case "remove-member":
		if targetUserID == "" {
			return ErrInvalidRequestPayload
		}
		return h.service.HandleAdminMemberOperation(
			services.RemoveMemberFromStudyGroupCommand, studyGroupID, targetUserID, userID)
	default:
		return ErrInvalidCommand
	}
}
