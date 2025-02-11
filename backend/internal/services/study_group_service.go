package services

import (
	"errors"
	"fmt"

	"gdp8-backend/internal/models"
	"gdp8-backend/internal/repositories"
)

type StudyGroupService interface {
	GetStudyGroupByID(id models.StudyGroupID) (*models.StudyGroup, error)
	GetAllStudyGroups() ([]models.StudyGroup, error)
}

var ErrStudyGroupNotFound = errors.New("study group not found")

type StudyGroupServiceImpl struct {
	studyGroupRepo repositories.StudyGroupRepository
}

func NewStudyGroupService(studyGroupRepo repositories.StudyGroupRepository) StudyGroupService {
	return &StudyGroupServiceImpl{studyGroupRepo: studyGroupRepo}
}

func (s StudyGroupServiceImpl) GetStudyGroupByID(id models.StudyGroupID) (*models.StudyGroup, error) {
	studyGroup, err := s.studyGroupRepo.GetStudyGroupByID(id)
	if err != nil {
		if errors.Is(err, repositories.ErrStudyGroupNotFound) {
			return nil, ErrStudyGroupNotFound
		}
		return nil, fmt.Errorf("error fetching study group: %w", err)
	}

	return studyGroup, nil
}

func (s StudyGroupServiceImpl) GetAllStudyGroups() ([]models.StudyGroup, error) {
	studyGroups, err := s.studyGroupRepo.GetAllStudyGroups()
	if err != nil {
		return nil, fmt.Errorf("error fetching study groups: %w", err)
	}
	return studyGroups, nil
}
