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

	AddStudySessionNotification(notificationType models.NotificationType, session *models.StudySession) error
	AddStudySessionReminderNotification(session *models.StudySession, studyGroupMembers []models.StudyGroupMemberView) error

	AddStudySessionAvailabilityRequestCreatedNotification(studyGroupID models.StudyGroupID) error
	AddStudySessionAvailabilityUpdatedNotification(availabilityRequest *models.StudySessionAvailabilityRequest,
		userID models.UserID) error
}

var ErrInvalidNotificationType = errors.New("invalid notification type")

type studySessionPayload struct {
	ID        models.StudySessionID `json:"id"`
	Title     string                `json:"title"`
	StartTime time.Time             `json:"startTime"`
	EndTime   time.Time             `json:"endTime"`
}
type studySessionAvailabilityRequestPayload struct {
	ID         models.StudySessionAvailabilityRequestID `json:"id"`
	StudyGroup studyGroupPayload                        `json:"studyGroup"`
	Title      string                                   `json:"title"`
}
type studyGroupPayload struct {
	ID   models.StudyGroupID `json:"id"`
	Name string              `json:"name"`
}
type userPayload struct {
	ID   models.UserID `json:"id"`
	Name string        `json:"name"`
}

type notificationServiceImpl struct {
	txMgr             persistence.TransactionManager
	notificationRepo  repositories.NotificationRepository
	userService       UserService
	studyGroupService StudyGroupService
}

func NewNotificationService(txMgr persistence.TransactionManager,
	notificationRepo repositories.NotificationRepository, userService UserService,
	studyGroupService StudyGroupService) NotificationService {

	return &notificationServiceImpl{
		txMgr:             txMgr,
		notificationRepo:  notificationRepo,
		userService:       userService,
		studyGroupService: studyGroupService,
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

func (s *notificationServiceImpl) AddStudySessionNotification(notificationType models.NotificationType,
	session *models.StudySession) error {

	if notificationType != models.NotificationTypeStudySessionScheduled &&
		notificationType != models.NotificationTypeStudySessionUpdated &&
		notificationType != models.NotificationTypeStudySessionCancelled {
		return ErrInvalidNotificationType
	}

	payload, err := json.Marshal(struct {
		StudySession studySessionPayload `json:"studySession"`
	}{
		StudySession: getStudySessionPayloadObject(session),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal notification payload: %w", err)
	}

	studyGroup, err := s.studyGroupService.GetStudyGroupByID(session.StudyGroupID)
	if err != nil {
		return fmt.Errorf("failed to get study group: %w", err)
	}
	if studyGroup == nil {
		return fmt.Errorf("failed to get study group: study group not found")
	}
	usersToBeNotified := getActualStudyGroupMembers(studyGroup.Members)

	return s.saveNotification(notificationType, payload, usersToBeNotified)
}

func (s *notificationServiceImpl) AddStudySessionReminderNotification(
	session *models.StudySession,
	studyGroupMembers []models.StudyGroupMemberView,
) error {

	payload, err := json.Marshal(struct {
		StudySession studySessionPayload `json:"studySession"`
	}{
		StudySession: getStudySessionPayloadObject(session),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal notification payload: %w", err)
	}

	usersToBeNotified := getActualStudyGroupMembers(studyGroupMembers)

	return s.saveNotification(models.NotificationTypeStudySessionReminder, payload, usersToBeNotified)
}

func (s *notificationServiceImpl) AddStudySessionAvailabilityRequestCreatedNotification(
	studyGroupID models.StudyGroupID) error {

	studyGroup, err := s.studyGroupService.GetStudyGroupByID(studyGroupID)
	if err != nil {
		return fmt.Errorf("failed to get study group: %w", err)
	}
	if studyGroup == nil {
		return fmt.Errorf("failed to get study group: study group not found")
	}

	payload, err := json.Marshal(struct {
		StudyGroup studyGroupPayload `json:"studyGroup"`
	}{
		StudyGroup: studyGroupPayload{
			ID:   studyGroupID,
			Name: studyGroup.Name,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to marshal notification payload: %w", err)
	}

	usersToBeNotified := getActualStudyGroupMembers(studyGroup.Members)

	return s.saveNotification(models.NotificationTypeStudySessionAvailabilityRequestCreated, payload, usersToBeNotified)
}

func (s *notificationServiceImpl) AddStudySessionAvailabilityUpdatedNotification(
	availabilityRequest *models.StudySessionAvailabilityRequest, userID models.UserID) error {

	triggeringUser, err := s.userService.GetUser(userID)
	if err != nil {
		return fmt.Errorf("failed to get triggering user: %w", err)
	}
	if triggeringUser == nil {
		return fmt.Errorf("failed to get triggering user: user not found")
	}

	studyGroup, err := s.studyGroupService.GetStudyGroupByID(availabilityRequest.StudyGroupID)
	if err != nil {
		return fmt.Errorf("failed to get study group: %w", err)
	}
	if studyGroup == nil {
		return fmt.Errorf("failed to get study group: study group not found")
	}

	payload, err := json.Marshal(struct {
		StudySessionAvailabilityRequest studySessionAvailabilityRequestPayload `json:"studySessionAvailabilityRequest"`
	}{
		StudySessionAvailabilityRequest: studySessionAvailabilityRequestPayload{
			ID: availabilityRequest.ID,
			StudyGroup: studyGroupPayload{
				ID:   studyGroup.ID,
				Name: studyGroup.Name,
			},
			Title: availabilityRequest.Title,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to marshal notification payload: %w", err)
	}

	usersToBeNotified := getActualStudyGroupMembers(studyGroup.Members)

	return s.saveNotification(models.NotificationTypeStudySessionAvailabilityEntriesUpdated, payload, usersToBeNotified)
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

	var targetUserPayload *userPayload
	if targetUserID != nil {
		targetUser, err := s.userService.GetUser(*targetUserID)
		if err != nil {
			return fmt.Errorf("failed to get target user: %w", err)
		}
		if targetUser == nil {
			return fmt.Errorf("failed to get target user: user not found")
		}
		targetUserPayload = &userPayload{
			ID:   *targetUserID,
			Name: targetUser.Name,
		}
	}

	payload, err := json.Marshal(struct {
		TriggeringUser userPayload       `json:"triggeringUser"`
		TargetUser     *userPayload      `json:"targetUser,omitempty"`
		StudyGroup     studyGroupPayload `json:"studyGroup"`
	}{
		TriggeringUser: userPayload{
			ID:   triggeringUserID,
			Name: triggeringUser.Name,
		},
		TargetUser: targetUserPayload,
		StudyGroup: studyGroupPayload{
			ID:   studyGroupID,
			Name: studyGroupName,
		},
	})

	if err != nil {
		return fmt.Errorf("failed to marshal notification payload: %w", err)
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

	return s.saveNotification(notificationType, payload, usersToBeNotified)
}

func (s *notificationServiceImpl) saveNotification(notificationType models.NotificationType,
	payload models.NotificationPayloadType, usersToBeNotified []models.UserID) error {

	notification := models.NotificationDetails{
		Type:    notificationType,
		Payload: payload,
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

func getStudySessionPayloadObject(session *models.StudySession) studySessionPayload {
	return studySessionPayload{
		ID:        session.ID,
		Title:     session.Title,
		StartTime: session.StartTime,
		EndTime:   session.EndTime,
	}
}
