package handlers

import (
	"log"
	"net/http"
	"time"

	"gdp8-backend/internal/models"
	"gdp8-backend/internal/services"
	"gdp8-backend/internal/utils"
)

type NotificationHandler struct {
	service services.NotificationService
}

type NotificationDTO struct {
	ID               models.NotificationID   `json:"id"`
	Type             models.NotificationType `json:"type"`
	TriggeringUserID models.UserID           `json:"triggeringUserId"`
	TargetUserID     *models.UserID          `json:"targetUserId,omitempty"`
	StudyGroupID     models.StudyGroupID     `json:"studyGroupId"`
	MessageID        *int64                  `json:"messageId,omitempty"`
	CreatedAt        time.Time               `json:"createdAt"`
}

func NewNotificationHandler(service services.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

func (h *NotificationHandler) GetUserNotifications(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	notifications, err := h.service.GetUserNotifications(userID)
	if err != nil {
		log.Printf("Error fetching notifications: %v\n", err)
		http.Error(w, "Failed to fetch notifications", http.StatusInternalServerError)
	}

	dtoList := make([]NotificationDTO, 0, len(notifications))
	for _, notification := range notifications {
		dtoList = append(dtoList, mapNotificationToDTO(&notification))
	}

	sendJSONResponse(w, dtoList)
}

func (h *NotificationHandler) MarkNotificationAsRead(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idString := r.PathValue("id")
	id, err := utils.ConvertToType[models.NotificationID](idString)
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	err = h.service.MarkNotificationAsRead(userID, id)
	if err != nil {
		log.Printf("Error marking notification as read: %v\n", err)
		http.Error(w, "Failed to mark notification as read", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func mapNotificationToDTO(notification *models.Notification) NotificationDTO {
	return NotificationDTO{
		ID:               notification.ID,
		Type:             notification.Type,
		TriggeringUserID: notification.TriggeringUserID,
		TargetUserID:     notification.TargetUserID,
		StudyGroupID:     notification.StudyGroupID,
		MessageID:        notification.MessageID,
		CreatedAt:        notification.CreatedAt,
	}
}
