package services

import (
	"database/sql"
	"errors"

	"gdp8-backend/internal/models"
	"gdp8-backend/internal/persistence"
	"gdp8-backend/internal/repositories"
)

type NotificationService interface {
	GetUserNotifications(userID models.UserID) ([]models.NotificationView, error)
	MarkNotificationAsRead(userID models.UserID, notificationID models.NotificationID) error
	AddStudyGroupEventNotification(
		notificationType models.NotificationType,
		triggeringUserID models.UserID,
		targetUserID *models.UserID,
		studyGroupID models.StudyGroupID,
		studyGroupMembers []models.StudyGroupMemberView,
	) error

	// NEW: AddStudySessionReminderNotification sends a reminder for a study session.
	AddStudySessionReminderNotification(session models.StudySession, studyGroupMembers []models.StudyGroupMemberView) error
}

var ErrInvalidNotificationType = errors.New("invalid notification type")

type notificationServiceImpl struct {
	txMgr            persistence.TransactionManager
	notificationRepo repositories.NotificationRepository
}

func NewNotificationService(txMgr persistence.TransactionManager,
	notificationRepo repositories.NotificationRepository) NotificationService {

	return &notificationServiceImpl{
		txMgr:            txMgr,
		notificationRepo: notificationRepo,
	}
}

func (s *notificationServiceImpl) GetUserNotifications(userID models.UserID) ([]models.NotificationView, error) {
	return persistence.WithTransaction(s.txMgr, func(tx *sql.Tx) ([]models.NotificationView, error) {
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
	session models.StudySession,
	studyGroupMembers []models.StudyGroupMemberView,
) error {
	notification := models.NotificationDetails{
		Type:             models.NotificationTypeStudySessionReminder,
		TriggeringUserID: session.CreatorID,
		StudyGroupID:     session.StudyGroupID,
		MessageID:        nil,
	}

	var usersToBeNotified []models.UserID
	for _, member := range studyGroupMembers {
		if member.Role == models.RoleMember || member.Role == models.RoleAdmin {
			usersToBeNotified = append(usersToBeNotified, member.UserID)
		}
	}

	return persistence.WithTransactionNoReturnVal(s.txMgr, func(tx *sql.Tx) error {
		return s.notificationRepo.AddNotification(tx, notification, usersToBeNotified)
	})
}

func (s *notificationServiceImpl) AddStudyGroupEventNotification(
	notificationType models.NotificationType,
	triggeringUserID models.UserID, targetUserID *models.UserID,
	studyGroupID models.StudyGroupID, studyGroupMembers []models.StudyGroupMemberView) error {

	notification := models.NotificationDetails{
		Type:             notificationType,
		TriggeringUserID: triggeringUserID,
		TargetUserID:     targetUserID,
		StudyGroupID:     studyGroupID,
		MessageID:        nil,
	}

	actualStudyGroupMembers := make([]models.UserID, 0, len(studyGroupMembers))
	for _, member := range studyGroupMembers {
		if member.Role == models.RoleMember || member.Role == models.RoleAdmin {
			actualStudyGroupMembers = append(actualStudyGroupMembers, member.UserID)
		}
	}

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
