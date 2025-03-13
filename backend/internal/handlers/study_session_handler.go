package handlers

import (
	"encoding/json"
	"net/http"
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

	sessions, err := h.service.GetAllStudySessionsByStudyGroup(groupID)
	if err != nil {
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

	requests, err := h.service.GetCurrentStudySessionAvailabilityRequests(groupID)
	if err != nil {
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

	details := &models.StudySessionAvailabilityRequestDetails{
		Title:                   req.Title,
		AvailabilityPeriodStart: req.AvailabilityPeriodStart,
		AvailabilityPeriodEnd:   req.AvailabilityPeriodEnd,
	}

	if err := h.service.CreateStudySessionAvailabilityRequest(groupID, details); err != nil {
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

	if err := h.service.DeleteStudySessionAvailabilityRequest(requestID); err != nil {
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
