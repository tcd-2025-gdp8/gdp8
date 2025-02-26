package services

import (
	"database/sql"
	"errors"

	"gdp8-backend/internal/models"
	"gdp8-backend/internal/persistence"
	"gdp8-backend/internal/repositories"
)

type NotificationService interface {
	GetUserNotifications(userID models.UserID) ([]models.Notification, error)

	MarkNotificationAsRead(userID models.UserID, notificationID models.NotificationID) error

	AddStudyGroupEventNotification(notificationType models.NotificationType,
		triggeringUserID models.UserID, targetUserID *models.UserID,
		studyGroupID models.StudyGroupID, studyGroupMembers []models.UserID) error
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

func (s *notificationServiceImpl) AddStudyGroupEventNotification(
	notificationType models.NotificationType,
	triggeringUserID models.UserID, targetUserID *models.UserID,
	studyGroupID models.StudyGroupID, studyGroupMembers []models.UserID) error {

	notification := models.NotificationDetails{
		Type:             notificationType,
		TriggeringUserID: triggeringUserID,
		TargetUserID:     targetUserID,
		StudyGroupID:     studyGroupID,
		MessageID:        nil,
	}

	var usersToBeNotified []models.UserID

	switch notificationType {
	case models.NotificationTypeStudyGroupJoined:
		usersToBeNotified = studyGroupMembers
	case models.NotificationTypeStudyGroupRequestedToJoin:
		usersToBeNotified = studyGroupMembers
	case models.NotificationTypeStudyGroupLeft:
		usersToBeNotified = studyGroupMembers
	case models.NotificationTypeStudyGroupAcceptedInvite:
		usersToBeNotified = studyGroupMembers
	case models.NotificationTypeStudyGroupRejectedInvite:
		usersToBeNotified = studyGroupMembers
	case models.NotificationTypeStudyGroupInvited:
		usersToBeNotified = append([]models.UserID{triggeringUserID}, studyGroupMembers...)
	case models.NotificationTypeStudyGroupAcceptedJoinRequest:
		usersToBeNotified = append([]models.UserID{triggeringUserID}, studyGroupMembers...)
	case models.NotificationTypeStudyGroupRejectedJoinRequest:
		usersToBeNotified = append([]models.UserID{triggeringUserID}, studyGroupMembers...)
	case models.NotificationTypeStudyGroupRemovedMember:
		usersToBeNotified = append([]models.UserID{triggeringUserID}, studyGroupMembers...)
	case models.NotificationTypeStudyGroupChatMessage:
		return ErrInvalidNotificationType
	default:
		return ErrInvalidNotificationType
	}

	return persistence.WithTransactionNoReturnVal(s.txMgr, func(tx *sql.Tx) error {
		return s.notificationRepo.AddNotification(tx, notification, usersToBeNotified)
	})
}
