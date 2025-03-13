package services

import (
	"database/sql"
	"errors"

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
		availabilityRequestDetails *models.StudySessionAvailabilityRequestDetails, requesterID models.UserID) error
	DeleteStudySessionAvailabilityRequest(
		availabilityRequestID models.StudySessionAvailabilityRequestID, requesterID models.UserID) error
	UpsertUserAvailabilityEntries(availabilityRequestID models.StudySessionAvailabilityRequestID,
		userID models.UserID, availabilityEntries []models.AvailabilityEntry) error
}

var ErrUnauthorizedStudySessionOperation = errors.New("unauthorized study session operation")

type studySessionServiceImpl struct {
	txMgr                        persistence.TransactionManager
	studySessionAvailabilityRepo repositories.StudySessionAvailabilityRepository
	studySessionRepository       repositories.StudySessionRepository
	studyGroupService            StudyGroupService
}

func NewStudySessionService(txMgr persistence.TransactionManager,
	studySessionAvailabilityRepo repositories.StudySessionAvailabilityRepository,
	studySessionRepository repositories.StudySessionRepository,
	studyGroupService StudyGroupService) StudySessionService {

	return &studySessionServiceImpl{
		txMgr:                        txMgr,
		studySessionAvailabilityRepo: studySessionAvailabilityRepo,
		studySessionRepository:       studySessionRepository,
		studyGroupService:            studyGroupService,
	}
}

func (s *studySessionServiceImpl) GetAllStudySessionsByStudyGroup(
	studyGroupID models.StudyGroupID, requesterID models.UserID) ([]models.StudySession, error) {

	isMember, err := isStudyGroupMember(studyGroupID, requesterID, s.studyGroupService)
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

	isMember, err := isStudyGroupMember(studyGroupID, creatorID, s.studyGroupService)
	if err != nil {
		return nil, err
	}

	if !isMember {
		return nil, ErrUnauthorizedStudySessionOperation
	}

	studySession, err := persistence.WithTransaction(s.txMgr, func(tx *sql.Tx) (*models.StudySession, error) {
		return s.studySessionRepository.CreateStudySession(tx, studyGroupID, creatorID, studySessionDetails)
	})

	// TODO send a notification
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

	// TODO send a notification
	// TODO update the calendar invite

	return studySession, err
}

func (s *studySessionServiceImpl) DeleteStudySession(studySessionID models.StudySessionID,
	requesterID models.UserID) error {

	err := persistence.WithTransactionNoReturnVal(s.txMgr, func(tx *sql.Tx) error {
		studySession, err := s.studySessionRepository.GetStudySession(tx, studySessionID)
		if err != nil {
			return err
		}
		if studySession.CreatorID != requesterID {
			return ErrUnauthorizedStudySessionOperation
		}

		return s.studySessionRepository.DeleteStudySession(tx, studySessionID)
	})

	// TODO send a notification
	// TODO update the calendar invite

	return err
}

func (s *studySessionServiceImpl) GetCurrentStudySessionAvailabilityRequests(
	studyGroupID models.StudyGroupID, requesterID models.UserID) ([]models.StudySessionAvailabilityRequest, error) {

	return persistence.WithTransaction(s.txMgr, func(tx *sql.Tx) ([]models.StudySessionAvailabilityRequest, error) {
		return s.studySessionAvailabilityRepo.GetCurrentStudySessionAvailabilityRequests(tx, studyGroupID)
	})
}

func (s *studySessionServiceImpl) CreateStudySessionAvailabilityRequest(studyGroupID models.StudyGroupID,
	availabilityRequestDetails *models.StudySessionAvailabilityRequestDetails, requesterID models.UserID) error {

	// TODO validate creator

	err := persistence.WithTransactionNoReturnVal(s.txMgr, func(tx *sql.Tx) error {
		return s.studySessionAvailabilityRepo.CreateStudySessionAvailabilityRequest(tx,
			studyGroupID, availabilityRequestDetails)
	})

	// TODO send a notification

	return err
}

func (s *studySessionServiceImpl) DeleteStudySessionAvailabilityRequest(
	availabilityRequestID models.StudySessionAvailabilityRequestID, requesterID models.UserID) error {

	// TODO validate creator

	return persistence.WithTransactionNoReturnVal(s.txMgr, func(tx *sql.Tx) error {
		return s.studySessionAvailabilityRepo.DeleteStudySessionAvailabilityRequest(tx, availabilityRequestID)
	})
}

func (s *studySessionServiceImpl) UpsertUserAvailabilityEntries(
	availabilityRequestID models.StudySessionAvailabilityRequestID, userID models.UserID,
	availabilityEntries []models.AvailabilityEntry) error {

	// TODO verify the entries against the availability request
	// TODO validate user

	err := persistence.WithTransactionNoReturnVal(s.txMgr, func(tx *sql.Tx) error {
		return s.studySessionAvailabilityRepo.UpsertUserAvailabilityEntries(tx,
			availabilityRequestID, userID, availabilityEntries)
	})

	// TODO send a notification

	return err
}

func isStudyGroupMember(studyGroupID models.StudyGroupID, userID models.UserID,
	studyGroupService StudyGroupService) (bool, error) {

	creatorStudyGroupRole, err := studyGroupService.RetrieveGroupRole(studyGroupID, userID)

	if err != nil {
		return false, err
	}

	if creatorStudyGroupRole == nil {
		return false, nil
	}

	if *creatorStudyGroupRole == models.RoleAdmin || *creatorStudyGroupRole == models.RoleMember {
		return true, nil
	}

	return false, nil
}
