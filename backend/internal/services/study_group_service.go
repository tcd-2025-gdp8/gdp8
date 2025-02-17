package services

import (
	"errors"
	"fmt"

	"gdp8-backend/internal/models"
	"gdp8-backend/internal/persistence"
	"gdp8-backend/internal/repositories"
)

type StudyGroupService interface {
	GetStudyGroupByID(id models.StudyGroupID) (*models.StudyGroup, error)
	GetAllStudyGroups() ([]models.StudyGroup, error)
}

var ErrStudyGroupNotFound = errors.New("study group not found")

type studyGroupServiceImpl struct {
	txMgr          persistence.TransactionManager
	studyGroupRepo repositories.StudyGroupRepository
}

func NewStudyGroupService(
	txMgr persistence.TransactionManager,
	studyGroupRepo repositories.StudyGroupRepository) StudyGroupService {

	return &studyGroupServiceImpl{
		txMgr:          txMgr,
		studyGroupRepo: studyGroupRepo,
	}
}

func (s *studyGroupServiceImpl) GetStudyGroupByID(id models.StudyGroupID) (*models.StudyGroup, error) {
	studyGroup, err := persistence.WithTransaction(s.txMgr, func(tx persistence.Transaction) (*models.StudyGroup, error) {
		return s.studyGroupRepo.GetStudyGroupByID(tx, id)
	})

	if err != nil {
		if errors.Is(err, repositories.ErrStudyGroupNotFound) {
			return nil, ErrStudyGroupNotFound
		}
		return nil, fmt.Errorf("error fetching study group: %w", err)
	}

	return studyGroup, nil
}

func (s *studyGroupServiceImpl) GetAllStudyGroups() ([]models.StudyGroup, error) {
	studyGroups, err := persistence.WithTransaction(s.txMgr, func(tx persistence.Transaction) ([]models.StudyGroup, error) {
		return s.studyGroupRepo.GetAllStudyGroups(tx)
	})

	if err != nil {
		return nil, fmt.Errorf("error fetching study groups: %w", err)
	}

	return studyGroups, nil
}
