package services

import (
	"database/sql"
	"gdp8-backend/internal/models"
	"gdp8-backend/internal/persistence"
	"gdp8-backend/internal/repositories"
)

type StudySessionService interface {
	GetAllStudySessionsByStudyGroup(studyGroupID models.StudyGroupID) ([]models.StudySession, error)
	GetAllStudySessionsByUser(userID models.UserID) ([]models.StudySession, error)
	CreateStudySession(studyGroupID models.StudyGroupID, creatorID models.UserID,
		studySessionDetails *models.StudySessionDetails) (*models.StudySession, error)
	UpdateStudySession(studySessionID models.StudySessionID,
		studySessionDetails *models.StudySessionDetails) (*models.StudySession, error)
	DeleteStudySession(studySessionID models.StudySessionID) error

	GetCurrentStudySessionAvailabilityRequests(
		studyGroupID models.StudyGroupID) ([]models.StudySessionAvailabilityRequest, error)
	CreateStudySessionAvailabilityRequest(studyGroupID models.StudyGroupID,
		availabilityRequestDetails *models.StudySessionAvailabilityRequestDetails) error
	DeleteStudySessionAvailabilityRequest(
		availabilityRequestID models.StudySessionAvailabilityRequestID) error
	UpsertUserAvailabilityEntries(availabilityRequestID models.StudySessionAvailabilityRequestID,
		userID models.UserID, availabilityEntries []models.AvailabilityEntry) error
}

type studySessionServiceImpl struct {
	txMgr                        persistence.TransactionManager
	studySessionAvailabilityRepo repositories.StudySessionAvailabilityRepository
	studySessionRepository       repositories.StudySessionRepository
}

func NewStudySessionService(txMgr persistence.TransactionManager,
	studySessionAvailabilityRepo repositories.StudySessionAvailabilityRepository,
	studySessionRepository repositories.StudySessionRepository) StudySessionService {

	return &studySessionServiceImpl{
		txMgr:                        txMgr,
		studySessionAvailabilityRepo: studySessionAvailabilityRepo,
		studySessionRepository:       studySessionRepository,
	}
}

func (s *studySessionServiceImpl) GetAllStudySessionsByStudyGroup(
	studyGroupID models.StudyGroupID) ([]models.StudySession, error) {

	return persistence.WithTransaction(s.txMgr, func(tx *sql.Tx) ([]models.StudySession, error) {
		return s.studySessionRepository.GetAllStudySessionsByStudyGroup(tx, studyGroupID)
	})
}

func (s *studySessionServiceImpl) GetAllStudySessionsByUser(userID models.UserID) ([]models.StudySession, error) {
	//TODO implement me
	panic("implement me")
}

func (s *studySessionServiceImpl) CreateStudySession(studyGroupID models.StudyGroupID,
	creatorID models.UserID, studySessionDetails *models.StudySessionDetails) (*models.StudySession, error) {

	// TODO validate creator

	studySession, err := persistence.WithTransaction(s.txMgr, func(tx *sql.Tx) (*models.StudySession, error) {
		return s.studySessionRepository.CreateStudySession(tx, studyGroupID, creatorID, studySessionDetails)
	})

	// TODO send a notification
	// TODO send a calendar invite

	return studySession, err
}

func (s *studySessionServiceImpl) UpdateStudySession(studySessionID models.StudySessionID,
	studySessionDetails *models.StudySessionDetails) (*models.StudySession, error) {

	// TODO validate creator

	studySession, err := persistence.WithTransaction(s.txMgr, func(tx *sql.Tx) (*models.StudySession, error) {
		return s.studySessionRepository.UpdateStudySession(tx, studySessionID, studySessionDetails)
	})

	// TODO send a notification
	// TODO update the calendar invite

	return studySession, err
}

func (s *studySessionServiceImpl) DeleteStudySession(studySessionID models.StudySessionID) error {

	// TODO validate creator

	err := persistence.WithTransactionNoReturnVal(s.txMgr, func(tx *sql.Tx) error {
		return s.studySessionRepository.DeleteStudySession(tx, studySessionID)
	})

	// TODO send a notification
	// TODO update the calendar invite

	return err
}

func (s *studySessionServiceImpl) GetCurrentStudySessionAvailabilityRequests(
	studyGroupID models.StudyGroupID) ([]models.StudySessionAvailabilityRequest, error) {

	return persistence.WithTransaction(s.txMgr, func(tx *sql.Tx) ([]models.StudySessionAvailabilityRequest, error) {
		return s.studySessionAvailabilityRepo.GetCurrentStudySessionAvailabilityRequests(tx, studyGroupID)
	})
}

func (s *studySessionServiceImpl) CreateStudySessionAvailabilityRequest(studyGroupID models.StudyGroupID,
	availabilityRequestDetails *models.StudySessionAvailabilityRequestDetails) error {

	err := persistence.WithTransactionNoReturnVal(s.txMgr, func(tx *sql.Tx) error {
		return s.studySessionAvailabilityRepo.CreateStudySessionAvailabilityRequest(tx,
			studyGroupID, availabilityRequestDetails)
	})

	// TODO send a notification

	return err
}

func (s *studySessionServiceImpl) DeleteStudySessionAvailabilityRequest(
	availabilityRequestID models.StudySessionAvailabilityRequestID) error {

	return persistence.WithTransactionNoReturnVal(s.txMgr, func(tx *sql.Tx) error {
		return s.studySessionAvailabilityRepo.DeleteStudySessionAvailabilityRequest(tx, availabilityRequestID)
	})
}

func (s *studySessionServiceImpl) UpsertUserAvailabilityEntries(
	availabilityRequestID models.StudySessionAvailabilityRequestID, userID models.UserID,
	availabilityEntries []models.AvailabilityEntry) error {

	// TODO verify the entries against the availability request

	err := persistence.WithTransactionNoReturnVal(s.txMgr, func(tx *sql.Tx) error {
		return s.studySessionAvailabilityRepo.UpsertUserAvailabilityEntries(tx,
			availabilityRequestID, userID, availabilityEntries)
	})

	// TODO send a notification

	return err
}
