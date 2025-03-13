package services

import (
	"database/sql"

	"gdp8-backend/internal/models"
	"gdp8-backend/internal/persistence"
	"gdp8-backend/internal/repositories"
)

type FileService interface {
	GetAllFilesByGroupID(groupID models.StudyGroupID) ([]models.File, error)
	CreateFile(file models.File) (*models.File, error)
	DeleteFile(fileName string, groupID models.StudyGroupID) error
	HasDeletionRights(fileName string, groupID models.StudyGroupID, userID models.UserID) (bool, error)
}

type fileServiceImpl struct {
	txManager         persistence.TransactionManager
	fileRepo          repositories.FileRepository
	studyGroupService StudyGroupService
}

func NewFileService(txManager persistence.TransactionManager,
	fileRepo repositories.FileRepository,
	studyGroupService StudyGroupService) FileService {
	return &fileServiceImpl{
		txManager:         txManager,
		fileRepo:          fileRepo,
		studyGroupService: studyGroupService,
	}
}

func (s *fileServiceImpl) GetAllFilesByGroupID(groupID models.StudyGroupID) ([]models.File, error) {
	return persistence.WithTransaction(s.txManager, func(tx *sql.Tx) ([]models.File, error) {
		return s.fileRepo.GetAllFilesByGroupID(tx, groupID)
	})
}

func (s *fileServiceImpl) CreateFile(file models.File) (*models.File, error) {
	return persistence.WithTransaction(s.txManager, func(tx *sql.Tx) (*models.File, error) {
		return s.fileRepo.CreateFile(tx, file)
	})
}

func (s *fileServiceImpl) DeleteFile(fileName string, groupID models.StudyGroupID) error {
	return persistence.WithTransactionNoReturnVal(s.txManager, func(tx *sql.Tx) error {
		return s.fileRepo.DeleteFile(tx, fileName, groupID)
	})
}

func (s *fileServiceImpl) HasDeletionRights(fileName string,
	groupID models.StudyGroupID,
	userID models.UserID) (bool, error) {
	return persistence.WithTransaction(s.txManager, func(tx *sql.Tx) (bool, error) {
		file, err := s.fileRepo.GetFileByNameAndGroupID(tx, fileName, groupID)
		if err != nil {
			return false, err
		}
		if file.UserID == userID {
			return true, nil
		}
		group, err := s.studyGroupService.GetStudyGroupByID(groupID)
		if err != nil {
			return false, err
		}
		for _, member := range group.Members {
			if member.UserID == userID && member.Role == models.RoleAdmin {
				return true, nil
			}
		}
		return false, nil
	})
}
