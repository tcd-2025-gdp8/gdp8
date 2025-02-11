package services

import (
	"errors"
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
		return nil, err
	}

	return studyGroup, nil
}

func (s StudyGroupServiceImpl) GetAllStudyGroups() ([]models.StudyGroup, error) {
	return s.studyGroupRepo.GetAllStudyGroups()
}
