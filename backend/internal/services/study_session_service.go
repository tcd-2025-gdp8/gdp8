package services

import (
	"database/sql"
	"errors"
	"log"

	"gdp8-backend/internal/models"
	"gdp8-backend/internal/persistence"
	"gdp8-backend/internal/repositories"
)

type StudySessionService interface {
	GetAllStudySessionsByStudyGroup(studyGroupID models.StudyGroupID,
		requesterID models.UserID) ([]models.StudySession, error)
	GetAllStudySessionsByUser(userID models.UserID) ([]models.StudySession, error)
	CreateStudySession(studyGroupID models.StudyGroupID, creatorID models.UserID,
		studySessionDetails *models.StudySessionDetails) (*models.StudySession, error)
	UpdateStudySession(studySessionID models.StudySessionID,
		studySessionDetails *models.StudySessionDetails, requesterID models.UserID) (*models.StudySession, error)
	DeleteStudySession(studySessionID models.StudySessionID, requesterID models.UserID) error

	GetCurrentStudySessionAvailabilityRequests(
		studyGroupID models.StudyGroupID, requesterID models.UserID) ([]models.StudySessionAvailabilityRequest, error)
	CreateStudySessionAvailabilityRequest(studyGroupID models.StudyGroupID,
		availabilityRequestDetails *models.StudySessionAvailabilityRequestDetails, creatorID models.UserID) error
	DeleteStudySessionAvailabilityRequest(
		availabilityRequestID models.StudySessionAvailabilityRequestID, requesterID models.UserID) error
	UpsertUserAvailabilityEntries(availabilityRequestID models.StudySessionAvailabilityRequestID,
		userID models.UserID, availabilityEntries []models.AvailabilityEntry) error
}

var ErrUnauthorizedStudySessionOperation = errors.New("unauthorized study session operation")
var ErrInvalidStudySessionAvailabilityEntry = errors.New("invalid study session availability entry")

type studySessionServiceImpl struct {
	txMgr                        persistence.TransactionManager
	studySessionAvailabilityRepo repositories.StudySessionAvailabilityRepository
	studySessionRepository       repositories.StudySessionRepository
	studyGroupService            StudyGroupService
	notificationService          NotificationService
}

func NewStudySessionService(txMgr persistence.TransactionManager,
	studySessionAvailabilityRepo repositories.StudySessionAvailabilityRepository,
	studySessionRepository repositories.StudySessionRepository,
	studyGroupService StudyGroupService,
	notificationService NotificationService) StudySessionService {

	return &studySessionServiceImpl{
		txMgr:                        txMgr,
		studySessionAvailabilityRepo: studySessionAvailabilityRepo,
		studySessionRepository:       studySessionRepository,
		studyGroupService:            studyGroupService,
		notificationService:          notificationService,
	}
}

func (s *studySessionServiceImpl) GetAllStudySessionsByStudyGroup(
	studyGroupID models.StudyGroupID, requesterID models.UserID) ([]models.StudySession, error) {

	isMember, err := s.studyGroupService.IsGroupMember(studyGroupID, requesterID)
	if err != nil {
		return nil, err
	}

	if !isMember {
		return nil, ErrUnauthorizedStudySessionOperation
	}

	return persistence.WithTransaction(s.txMgr, func(tx *sql.Tx) ([]models.StudySession, error) {
		return s.studySessionRepository.GetAllStudySessionsByStudyGroup(tx, studyGroupID)
	})
}

func (s *studySessionServiceImpl) GetAllStudySessionsByUser(_ models.UserID) ([]models.StudySession, error) {
	// TODO implement me
	panic("implement me")
}

func (s *studySessionServiceImpl) CreateStudySession(studyGroupID models.StudyGroupID,
	creatorID models.UserID, studySessionDetails *models.StudySessionDetails) (*models.StudySession, error) {

	isMember, err := s.studyGroupService.IsGroupMember(studyGroupID, creatorID)
	if err != nil {
		return nil, err
	}

	if !isMember {
		return nil, ErrUnauthorizedStudySessionOperation
	}

	studySession, err := persistence.WithTransaction(s.txMgr, func(tx *sql.Tx) (*models.StudySession, error) {
		return s.studySessionRepository.CreateStudySession(tx, studyGroupID, creatorID, studySessionDetails)
	})

	if err != nil {
		return nil, err
	}

	go func() {
		notificationErr := s.notificationService.AddStudySessionNotification(
			models.NotificationTypeStudySessionScheduled, studySession, creatorID)
		if notificationErr != nil {
			log.Printf("Error sending notification: %v\n", notificationErr)
		}
	}()
	// TODO send a calendar invite

	return studySession, err
}

func (s *studySessionServiceImpl) UpdateStudySession(studySessionID models.StudySessionID,
	studySessionDetails *models.StudySessionDetails, requesterID models.UserID) (*models.StudySession, error) {

	studySession, err := persistence.WithTransaction(s.txMgr, func(tx *sql.Tx) (*models.StudySession, error) {
		studySession, err := s.studySessionRepository.GetStudySession(tx, studySessionID)
		if err != nil {
			return nil, err
		}
		if studySession.CreatorID != requesterID {
			return nil, ErrUnauthorizedStudySessionOperation
		}

		return s.studySessionRepository.UpdateStudySession(tx, studySessionID, studySessionDetails)
	})

	if err != nil {
		return nil, err
	}

	go func() {
		notificationErr := s.notificationService.AddStudySessionNotification(
			models.NotificationTypeStudySessionUpdated, studySession, requesterID)
		if notificationErr != nil {
			log.Printf("Error sending notification: %v\n", notificationErr)
		}
	}()
	// TODO update the calendar invite

	return studySession, err
}

func (s *studySessionServiceImpl) DeleteStudySession(studySessionID models.StudySessionID,
	requesterID models.UserID) error {

	deletedStudySession, err := persistence.WithTransaction(s.txMgr, func(tx *sql.Tx) (*models.StudySession, error) {
		studySession, err := s.studySessionRepository.GetStudySession(tx, studySessionID)
		if err != nil {
			return nil, err
		}
		if studySession.CreatorID != requesterID {
			return nil, ErrUnauthorizedStudySessionOperation
		}

		err = s.studySessionRepository.DeleteStudySession(tx, studySessionID)
		if err != nil {
			return nil, err
		}
		return studySession, nil
	})

	if err != nil {
		return err
	}

	go func() {
		notificationErr := s.notificationService.AddStudySessionNotification(
			models.NotificationTypeStudySessionCancelled, deletedStudySession, requesterID)
		if notificationErr != nil {
			log.Printf("Error sending notification: %v\n", notificationErr)
		}
	}()
	// TODO update the calendar invite

	return err
}

func (s *studySessionServiceImpl) GetCurrentStudySessionAvailabilityRequests(
	studyGroupID models.StudyGroupID, requesterID models.UserID) ([]models.StudySessionAvailabilityRequest, error) {

	isMember, err := s.studyGroupService.IsGroupMember(studyGroupID, requesterID)
	if err != nil {
		return nil, err
	}

	if !isMember {
		return nil, ErrUnauthorizedStudySessionOperation
	}

	return persistence.WithTransaction(s.txMgr, func(tx *sql.Tx) ([]models.StudySessionAvailabilityRequest, error) {
		return s.studySessionAvailabilityRepo.GetCurrentStudySessionAvailabilityRequests(tx, studyGroupID)
	})
}

func (s *studySessionServiceImpl) CreateStudySessionAvailabilityRequest(studyGroupID models.StudyGroupID,
	availabilityRequestDetails *models.StudySessionAvailabilityRequestDetails, creatorID models.UserID) error {

	isMember, err := s.studyGroupService.IsGroupMember(studyGroupID, creatorID)
	if err != nil {
		return err
	}

	if !isMember {
		return ErrUnauthorizedStudySessionOperation
	}

	req, err := persistence.WithTransaction(s.txMgr, func(tx *sql.Tx) (*models.StudySessionAvailabilityRequest, error) {
		return s.studySessionAvailabilityRepo.CreateStudySessionAvailabilityRequest(tx,
			studyGroupID, creatorID, availabilityRequestDetails)
	})

	if err != nil {
		return err
	}

	go func() {
		notificationErr := s.notificationService.AddStudySessionAvailabilityRequestCreatedNotification(req, creatorID)
		if notificationErr != nil {
			log.Printf("Error sending notification: %v\n", notificationErr)
		}
	}()

	return err
}

func (s *studySessionServiceImpl) DeleteStudySessionAvailabilityRequest(
	availabilityRequestID models.StudySessionAvailabilityRequestID, _ models.UserID) error {

	// TODO validate user

	return persistence.WithTransactionNoReturnVal(s.txMgr, func(tx *sql.Tx) error {
		return s.studySessionAvailabilityRepo.DeleteStudySessionAvailabilityRequest(tx, availabilityRequestID)
	})
}

func (s *studySessionServiceImpl) UpsertUserAvailabilityEntries(
	availabilityRequestID models.StudySessionAvailabilityRequestID, userID models.UserID,
	availabilityEntries []models.AvailabilityEntry) error {

	req, err := persistence.WithTransaction(s.txMgr, func(tx *sql.Tx) (*models.StudySessionAvailabilityRequest, error) {
		availReq, err := s.studySessionAvailabilityRepo.GetStudySessionAvailabilityRequest(tx, availabilityRequestID)
		if err != nil {
			return nil, err
		}

		isMember, err := s.studyGroupService.IsGroupMember(availReq.StudyGroupID, userID)
		if err != nil {
			return nil, err
		}

		if !isMember {
			return nil, ErrUnauthorizedStudySessionOperation
		}

		for _, entry := range availabilityEntries {
			if entry.AvailabilityEntryStart.Before(availReq.AvailabilityPeriodStart) ||
				entry.AvailabilityEntryEnd.After(availReq.AvailabilityPeriodEnd) {

				return nil, ErrInvalidStudySessionAvailabilityEntry
			}
		}

		return s.studySessionAvailabilityRepo.UpsertUserAvailabilityEntries(tx,
			availabilityRequestID, userID, availabilityEntries)
	})

	if err != nil {
		return err
	}

	go func() {
		notificationErr := s.notificationService.AddStudySessionAvailabilityUpdatedNotification(req, userID)
		if notificationErr != nil {
			log.Printf("Error sending notification: %v\n", notificationErr)
		}
	}()

	return err
}
