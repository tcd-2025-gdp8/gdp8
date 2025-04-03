package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gdp8-backend/internal/models"
	"gdp8-backend/internal/persistence"
	"gdp8-backend/internal/repositories"
)

type NotificationService interface {
	GetUserNotifications(userID models.UserID) ([]models.Notification, error)
	MarkNotificationAsRead(userID models.UserID, notificationID models.NotificationID) error

	AddStudyGroupEventNotification(
		notificationType models.NotificationType,
		triggeringUserID models.UserID,
		targetUserID *models.UserID,
		studyGroupID models.StudyGroupID,
		studyGroupName string,
		studyGroupMembers []models.StudyGroupMemberView,
	) error

	AddStudySessionReminderNotification(session *models.StudySession, studyGroupMembers []models.StudyGroupMemberView) error
}

var ErrInvalidNotificationType = errors.New("invalid notification type")

type notificationServiceImpl struct {
	txMgr            persistence.TransactionManager
	notificationRepo repositories.NotificationRepository
	userService      UserService
}

func NewNotificationService(txMgr persistence.TransactionManager,
	notificationRepo repositories.NotificationRepository, userService UserService) NotificationService {

	return &notificationServiceImpl{
		txMgr:            txMgr,
		notificationRepo: notificationRepo,
		userService:      userService,
	}
}

func (s *notificationServiceImpl) GetUserNotifications(userID models.UserID) ([]models.Notification, error) {
	return persistence.WithTransaction(s.txMgr, func(tx *sql.Tx) ([]models.Notification, error) {
		return s.notificationRepo.GetUserNotifications(tx, userID)
	})
}

func (s *notificationServiceImpl) MarkNotificationAsRead(userID models.UserID,
	notificationID models.NotificationID) error {

	return persistence.WithTransactionNoReturnVal(s.txMgr, func(tx *sql.Tx) error {
		return s.notificationRepo.DeleteUserNotification(tx, notificationID, userID)
	})
}

func (s *notificationServiceImpl) AddStudySessionReminderNotification(
	session *models.StudySession,
	studyGroupMembers []models.StudyGroupMemberView,
) error {

	payload, err := json.Marshal(struct {
		StudySessionID        models.StudySessionID `json:"studySessionId"`
		StudySessionTitle     string                `json:"studySessionTitle"`
		StudySessionStartTime time.Time             `json:"studySessionStartTime"`
		StudySessionEndTime   time.Time             `json:"studySessionEndTime"`
	}{
		StudySessionID:        session.ID,
		StudySessionTitle:     session.Title,
		StudySessionStartTime: session.StartTime,
		StudySessionEndTime:   session.EndTime,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal notification payload: %w", err)
	}

	notification := models.NotificationDetails{
		Type:    models.NotificationTypeStudySessionReminder,
		Payload: payload,
	}

	usersToBeNotified := getActualStudyGroupMembers(studyGroupMembers)

	return persistence.WithTransactionNoReturnVal(s.txMgr, func(tx *sql.Tx) error {
		return s.notificationRepo.AddNotification(tx, notification, usersToBeNotified)
	})
}

func (s *notificationServiceImpl) AddStudyGroupEventNotification(
	notificationType models.NotificationType,
	triggeringUserID models.UserID, targetUserID *models.UserID,
	studyGroupID models.StudyGroupID, studyGroupName string, studyGroupMembers []models.StudyGroupMemberView) error {

	triggeringUser, err := s.userService.GetUser(triggeringUserID)
	if err != nil {
		return fmt.Errorf("failed to get triggering user: %w", err)
	}
	if triggeringUser == nil {
		return fmt.Errorf("failed to get triggering user: user not found")
	}

	var targetUserName *string
	if targetUserID != nil {
		targetUser, err := s.userService.GetUser(*targetUserID)
		if err != nil {
			return fmt.Errorf("failed to get target user: %w", err)
		}
		if targetUser == nil {
			return fmt.Errorf("failed to get target user: user not found")
		}
		targetUserName = &targetUser.Name
	}

	payload, err := json.Marshal(struct {
		TriggeringUserID   models.UserID       `json:"triggeringUserId"`
		TriggeringUserName string              `json:"triggeringUserName"`
		TargetUserID       *models.UserID      `json:"targetUserId,omitempty"`
		TargetUserName     *string             `json:"targetUserName,omitempty"`
		StudyGroupID       models.StudyGroupID `json:"studyGroupId"`
		StudyGroupName     string              `json:"studyGroupName"`
	}{
		TriggeringUserID:   triggeringUserID,
		TriggeringUserName: triggeringUser.Name,
		TargetUserID:       targetUserID,
		TargetUserName:     targetUserName,
		StudyGroupID:       studyGroupID,
		StudyGroupName:     studyGroupName,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal notification payload: %w", err)
	}

	notification := models.NotificationDetails{
		Type:    notificationType,
		Payload: payload,
	}

	actualStudyGroupMembers := getActualStudyGroupMembers(studyGroupMembers)

	var usersToBeNotified []models.UserID

	switch notificationType {
	case models.NotificationTypeStudyGroupJoined:
		usersToBeNotified = actualStudyGroupMembers
	case models.NotificationTypeStudyGroupRequestedToJoin:
		usersToBeNotified = actualStudyGroupMembers
	case models.NotificationTypeStudyGroupLeft:
		usersToBeNotified = actualStudyGroupMembers
	case models.NotificationTypeStudyGroupAcceptedInvite:
		usersToBeNotified = actualStudyGroupMembers
	case models.NotificationTypeStudyGroupRejectedInvite:
		usersToBeNotified = actualStudyGroupMembers
	case models.NotificationTypeStudyGroupInvited:
		usersToBeNotified = append([]models.UserID{triggeringUserID}, actualStudyGroupMembers...)
	case models.NotificationTypeStudyGroupAcceptedJoinRequest:
		usersToBeNotified = append([]models.UserID{triggeringUserID}, actualStudyGroupMembers...)
	case models.NotificationTypeStudyGroupRejectedJoinRequest:
		usersToBeNotified = append([]models.UserID{triggeringUserID}, actualStudyGroupMembers...)
	case models.NotificationTypeStudyGroupRemovedMember:
		usersToBeNotified = append([]models.UserID{triggeringUserID}, actualStudyGroupMembers...)
	case models.NotificationTypeStudyGroupChatMessage:
		return ErrInvalidNotificationType
	case models.NotificationTypeStudySessionReminder:
		// This type is handled separately via AddStudySessionReminderNotification.
		return ErrInvalidNotificationType
	default:
		return ErrInvalidNotificationType
	}

	return persistence.WithTransactionNoReturnVal(s.txMgr, func(tx *sql.Tx) error {
		return s.notificationRepo.AddNotification(tx, notification, usersToBeNotified)
	})
}

func getActualStudyGroupMembers(members []models.StudyGroupMemberView) []models.UserID {
	filteredMembers := make([]models.UserID, 0, len(members))
	for _, member := range members {
		if member.Role == models.RoleMember || member.Role == models.RoleAdmin {
			filteredMembers = append(filteredMembers, member.UserID)
		}
	}
	return filteredMembers
}
