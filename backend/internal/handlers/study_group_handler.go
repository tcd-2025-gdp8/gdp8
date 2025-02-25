package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"slices"

	"gdp8-backend/internal/models"
	"gdp8-backend/internal/services"
	"gdp8-backend/internal/utils"
)

type StudyGroupHandler struct {
	service services.StudyGroupService
}

type StudyGroupMemberDTO struct {
	ID   models.UserID `json:"id"`
	Name string        `json:"name"`
	Role string        `json:"role"`
}

type StudyGroupDTO struct {
	ID          models.StudyGroupID `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Type        string              `json:"type"`
	ModuleID    models.ModuleID     `json:"moduleId"`
}

type StudyGroupWithMembersDTO struct {
	StudyGroupDTO
	Members []StudyGroupMemberDTO `json:"members,omitempty"`
}

var ErrInvalidRequestPayload = errors.New("invalid request payload")
var ErrInvalidCommand = errors.New("invalid command")

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

	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	studyGroup, err := h.service.GetStudyGroupByID(id)
	if err != nil {
		if errors.Is(err, services.ErrStudyGroupNotFound) {
			http.Error(w, "Study group not found", http.StatusNotFound)
			return
		}
		log.Printf("Error fetching study group: %v\n", err)
		http.Error(w, "Error fetching study group", http.StatusInternalServerError)
		return
	}

	if !canSeeGroup(userID, studyGroup) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if canSeeMembers(userID, studyGroup) {
		sendJSONResponse(w, mapStudyGroupWithMembersToDTO(studyGroup))
	} else {
		sendJSONResponse(w, mapStudyGroupToDTO(studyGroup))
	}
}

func (h *StudyGroupHandler) GetRelevantStudyGroups(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	studyGroups, err := h.service.GetAllRelevantStudyGroups(userID)
	if err != nil {
		log.Printf("Error fetching study groups: %v\n", err)
		http.Error(w, "Error fetching study groups", http.StatusInternalServerError)
		return
	}

	dtoList := make([]StudyGroupWithMembersDTO, 0, len(studyGroups))
	for _, studyGroup := range studyGroups {
		if !canSeeGroup(userID, &studyGroup) {
			continue
		}

		if canSeeMembers(userID, &studyGroup) {
			dtoList = append(dtoList, mapStudyGroupWithMembersToDTO(&studyGroup))
		} else {
			dtoList = append(dtoList, StudyGroupWithMembersDTO{
				StudyGroupDTO: mapStudyGroupToDTO(&studyGroup),
				Members:       nil,
			})
		}

	}

	sendJSONResponse(w, dtoList)
}

func (h *StudyGroupHandler) CreateStudyGroup(w http.ResponseWriter, r *http.Request) {
	var createDTO struct {
		Name        string          `json:"name"`
		Description string          `json:"description"`
		Type        string          `json:"type"`
		ModuleID    models.ModuleID `json:"moduleId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&createDTO); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// TODO implement validation (group type in particular)
	studyGroupDetails := models.StudyGroupDetails{
		Name:        createDTO.Name,
		Description: createDTO.Description,
		Type:        models.StudyGroupType(createDTO.Type),
		ModuleID:    createDTO.ModuleID,
	}

	createdStudyGroup, err := h.service.CreateStudyGroup(studyGroupDetails, userID)
	if err != nil {
		log.Printf("Error creating study group: %v\n", err)
		http.Error(w, "Error creating study group", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, mapStudyGroupWithMembersToDTO(createdStudyGroup))
}

func (h *StudyGroupHandler) HandleStudyMemberOperation(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	studyGroupID, err := utils.ConvertToType[models.StudyGroupID](idString)
	if err != nil {
		http.Error(w, "Invalid study group ID", http.StatusBadRequest)
		return
	}

	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var memberOperationDetails struct {
		TargetUserID string `json:"targetUserId"`
	}
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
		log.Printf("Error processing study group operation: %v\n", err)
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

func canSeeGroup(userID models.UserID, studyGroup *models.StudyGroupView) bool {
	if studyGroup.Type == models.TypeInviteOnly {
		return slices.ContainsFunc(studyGroup.Members, func(m models.StudyGroupMemberView) bool {
			return m.UserID == userID
		})
	}
	return true
}

func canSeeMembers(userID models.UserID, studyGroup *models.StudyGroupView) bool {
	if studyGroup.Type == models.TypePublic {
		return true
	}
	return slices.ContainsFunc(studyGroup.Members, func(m models.StudyGroupMemberView) bool {
		return m.UserID == userID
	})
}

func mapStudyGroupMemberToDTO(member *models.StudyGroupMemberView) StudyGroupMemberDTO {
	return StudyGroupMemberDTO{
		ID:   member.UserID,
		Name: member.Name,
		Role: string(member.Role),
	}
}

func mapStudyGroupWithMembersToDTO(studyGroup *models.StudyGroupView) StudyGroupWithMembersDTO {
	members := make([]StudyGroupMemberDTO, len(studyGroup.Members))
	for i, member := range studyGroup.Members {
		members[i] = mapStudyGroupMemberToDTO(&member)
	}

	return StudyGroupWithMembersDTO{
		StudyGroupDTO: mapStudyGroupToDTO(studyGroup),
		Members:       members,
	}
}

func mapStudyGroupToDTO(studyGroup *models.StudyGroupView) StudyGroupDTO {
	return StudyGroupDTO{
		ID:          studyGroup.ID,
		Name:        studyGroup.Name,
		Description: studyGroup.Description,
		Type:        string(studyGroup.Type),
		ModuleID:    studyGroup.ModuleID,
	}
}
