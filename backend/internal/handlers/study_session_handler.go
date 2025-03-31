package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"gdp8-backend/internal/models"
	"gdp8-backend/internal/services"
	"gdp8-backend/internal/utils"
)

type StudySessionResponse struct {
	ID              int64     `json:"id"`
	StudyGroupID    int64     `json:"studyGroupId"`
	CreatorID       string    `json:"creatorId"`
	Title           string    `json:"title"`
	StartTime       time.Time `json:"startTime"`
	DurationMinutes int64     `json:"durationMinutes"`
	EndTime         time.Time `json:"endTime"`
}

type StudySessionDetailsDTO struct {
	Title           string    `json:"title"`
	StartTime       time.Time `json:"startTime"`
	DurationMinutes int64     `json:"durationMinutes"`
}

type StudySessionAvailabilityRequestResponse struct {
	ID                      int64                              `json:"id"`
	StudyGroupID            int64                              `json:"studyGroupId"`
	Title                   string                             `json:"title"`
	AvailabilityPeriodStart time.Time                          `json:"availabilityPeriodStart"`
	AvailabilityPeriodEnd   time.Time                          `json:"availabilityPeriodEnd"`
	Entries                 []StudySessionAvailabilityEntryDTO `json:"entries"`
}

type StudySessionAvailabilityEntryDTO struct {
	UserID            string    `json:"userId"`
	AvailabilityStart time.Time `json:"availabilityStart"`
	AvailabilityEnd   time.Time `json:"availabilityEnd"`
}

const (
	minTitleLength = 3
	maxTitleLength = 255
)

type StudySessionHandler struct {
	service services.StudySessionService
}

func NewStudySessionHandler(service services.StudySessionService) *StudySessionHandler {
	return &StudySessionHandler{service: service}
}

func (h *StudySessionHandler) GetStudySessionsByGroup(w http.ResponseWriter, r *http.Request) {
	groupIDString := r.PathValue("groupID")
	groupID, err := utils.ConvertToType[models.StudyGroupID](groupIDString)
	if err != nil {
		http.Error(w, "Invalid study group ID", http.StatusBadRequest)
		return
	}

	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sessions, err := h.service.GetAllStudySessionsByStudyGroup(groupID, userID)
	if err != nil {
		log.Printf("Error getting study sessions for group ID %s, user ID %s: %v", groupIDString, userID, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := make([]StudySessionResponse, len(sessions))
	for i, session := range sessions {
		response[i] = toStudySessionResponse(session)
	}

	sendJSONResponse(w, response)
}

func (h *StudySessionHandler) GetStudySessionsByUser(_ http.ResponseWriter, _ *http.Request) {
	// TODO implement
	panic("Not implemented")
}

func (h *StudySessionHandler) CreateStudySession(w http.ResponseWriter, r *http.Request) {
	var req StudySessionDetailsDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := validateTitle(req.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	groupIDString := r.PathValue("groupID")
	groupID, err := utils.ConvertToType[models.StudyGroupID](groupIDString)
	if err != nil {
		http.Error(w, "Invalid study group ID", http.StatusBadRequest)
		return
	}

	creatorID, err := getUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	details := &models.StudySessionDetails{
		Title:           req.Title,
		StartTime:       req.StartTime,
		DurationMinutes: req.DurationMinutes,
	}

	session, err := h.service.CreateStudySession(groupID, creatorID, details)
	if err != nil {
		log.Printf("Error creating study session for group ID %s, creator ID %s: %v", groupIDString, creatorID, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, toStudySessionResponse(*session))
}

func (h *StudySessionHandler) UpdateStudySession(w http.ResponseWriter, r *http.Request) {
	studySessionIDString := r.PathValue("studySessionID")
	studySessionID, err := utils.ConvertToType[models.StudySessionID](studySessionIDString)
	if err != nil {
		http.Error(w, "Invalid study session ID", http.StatusBadRequest)
		return
	}

	var req StudySessionDetailsDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = validateTitle(req.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	details := &models.StudySessionDetails{
		Title:           req.Title,
		StartTime:       req.StartTime,
		DurationMinutes: req.DurationMinutes,
	}

	session, err := h.service.UpdateStudySession(studySessionID, details, userID)
	if err != nil {
		log.Printf("Error updating study session with ID %s by user %s: %v", studySessionIDString, userID, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, toStudySessionResponse(*session))
}

func (h *StudySessionHandler) DeleteStudySession(w http.ResponseWriter, r *http.Request) {
	studySessionIDString := r.PathValue("studySessionID")
	studySessionID, err := utils.ConvertToType[models.StudySessionID](studySessionIDString)
	if err != nil {
		http.Error(w, "Invalid study session ID", http.StatusBadRequest)
		return
	}

	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.service.DeleteStudySession(studySessionID, userID); err != nil {
		log.Printf("Error deleting study session with ID %s by user %s: %v", studySessionIDString, userID, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *StudySessionHandler) GetCurrentAvailabilityRequests(w http.ResponseWriter, r *http.Request) {
	groupIDString := r.PathValue("groupID")
	groupID, err := utils.ConvertToType[models.StudyGroupID](groupIDString)
	if err != nil {
		http.Error(w, "Invalid study group ID", http.StatusBadRequest)
		return
	}

	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	requests, err := h.service.GetCurrentStudySessionAvailabilityRequests(groupID, userID)
	if err != nil {
		log.Printf("Error getting current availability requests for group ID %s, user ID %s: %v", groupIDString, userID, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := make([]StudySessionAvailabilityRequestResponse, len(requests))
	for i, req := range requests {
		response[i] = toAvailabilityRequestResponse(req)
	}

	sendJSONResponse(w, response)
}

func (h *StudySessionHandler) CreateAvailabilityRequest(w http.ResponseWriter, r *http.Request) {
	groupIDString := r.PathValue("groupID")
	groupID, err := utils.ConvertToType[models.StudyGroupID](groupIDString)
	if err != nil {
		http.Error(w, "Invalid study group ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Title                   string    `json:"title"`
		AvailabilityPeriodStart time.Time `json:"availabilityPeriodStart"`
		AvailabilityPeriodEnd   time.Time `json:"availabilityPeriodEnd"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = validateTitle(req.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	details := &models.StudySessionAvailabilityRequestDetails{
		Title:                   req.Title,
		AvailabilityPeriodStart: req.AvailabilityPeriodStart,
		AvailabilityPeriodEnd:   req.AvailabilityPeriodEnd,
	}

	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.service.CreateStudySessionAvailabilityRequest(groupID, details, userID); err != nil {
		log.Printf("Error creating availability request for group ID %s, user ID %s: %v", groupIDString, userID, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *StudySessionHandler) DeleteAvailabilityRequest(w http.ResponseWriter, r *http.Request) {
	requestIDString := r.PathValue("availabilityRequestID")
	requestID, err := utils.ConvertToType[models.StudySessionAvailabilityRequestID](requestIDString)
	if err != nil {
		http.Error(w, "Invalid study group ID", http.StatusBadRequest)
		return
	}

	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.service.DeleteStudySessionAvailabilityRequest(requestID, userID); err != nil {
		log.Printf("Error deleting availability request with ID %s by user %s: %v", requestIDString, userID, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *StudySessionHandler) UpsertAvailabilityEntries(w http.ResponseWriter, r *http.Request) {
	requestIDString := r.PathValue("availabilityRequestID")
	requestID, err := utils.ConvertToType[models.StudySessionAvailabilityRequestID](requestIDString)
	if err != nil {
		http.Error(w, "Invalid study group ID", http.StatusBadRequest)
		return
	}

	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Entries []struct {
			AvailabilityStart time.Time `json:"availabilityStart"`
			AvailabilityEnd   time.Time `json:"availabilityEnd"`
		} `json:"entries"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	entries := make([]models.AvailabilityEntry, len(req.Entries))
	for i, entry := range req.Entries {
		entries[i] = models.AvailabilityEntry{
			AvailabilityEntryStart: entry.AvailabilityStart,
			AvailabilityEntryEnd:   entry.AvailabilityEnd,
		}
	}

	if err := h.service.UpsertUserAvailabilityEntries(requestID, userID, entries); err != nil {
		log.Printf("Error updating availability entries for request ID %s by user %s: %v", requestIDString, userID, err)
		log.Printf("Detailed error: %+v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func toStudySessionResponse(session models.StudySession) StudySessionResponse {
	return StudySessionResponse{
		ID:              int64(session.ID),
		StudyGroupID:    int64(session.StudyGroupID),
		CreatorID:       string(session.CreatorID),
		Title:           session.Title,
		StartTime:       session.StartTime,
		DurationMinutes: session.DurationMinutes,
		EndTime:         session.EndTime,
	}
}

func toAvailabilityRequestResponse(
	request models.StudySessionAvailabilityRequest) StudySessionAvailabilityRequestResponse {

	entries := make([]StudySessionAvailabilityEntryDTO, len(request.Entries))
	for i, entry := range request.Entries {
		entries[i] = StudySessionAvailabilityEntryDTO{
			UserID:            string(entry.UserID),
			AvailabilityStart: entry.AvailabilityEntryStart,
			AvailabilityEnd:   entry.AvailabilityEntryEnd,
		}
	}

	return StudySessionAvailabilityRequestResponse{
		ID:                      int64(request.ID),
		StudyGroupID:            int64(request.StudyGroupID),
		Title:                   request.Title,
		AvailabilityPeriodStart: request.AvailabilityPeriodStart,
		AvailabilityPeriodEnd:   request.AvailabilityPeriodEnd,
		Entries:                 entries,
	}
}

func validateTitle(title string) error {
	if len(strings.TrimSpace(title)) < minTitleLength {
		return fmt.Errorf("title must be at least %d characters long", minTitleLength)
	}
	if len(strings.TrimSpace(title)) > maxTitleLength {
		return fmt.Errorf("title must be at most %d characters long", maxTitleLength)
	}
	return nil
}
