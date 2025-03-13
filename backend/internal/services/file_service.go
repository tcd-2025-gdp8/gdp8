package services

import (
	"database/sql"
	"fmt"

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

// func (s *fileServiceImpl) HasDeletionRights(fileName string,
// 	groupID models.StudyGroupID,
// 	userID models.UserID) (bool, error) {
// 	return persistence.WithTransaction(s.txManager, func(tx *sql.Tx) (bool, error) {
// 		file, err := s.fileRepo.GetFileByNameAndGroupID(tx, fileName, groupID)
// 		if err != nil {
// 			return false, err
// 		}
// 		if file.UserID == userID {
// 			return true, nil
// 		}
// 		group, err := s.studyGroupService.GetStudyGroupByID(groupID)
// 		if err != nil {
// 			return false, err
// 		}
// 		for _, member := range group.Members {
// 			if member.UserID == userID && member.Role == models.RoleAdmin {
// 				return true, nil
// 			}
// 		}
// 		return false, nil
// 	})
// }

func (s *fileServiceImpl) HasDeletionRights(fileName string,
	groupID models.StudyGroupID,
	userID models.UserID) (bool, error) {
	return persistence.WithTransaction(s.txManager, func(tx *sql.Tx) (bool, error) {
		file, err := s.fileRepo.GetFileByNameAndGroupID(tx, fileName, groupID)
		if err != nil {
			fmt.Printf("[HasDeletionRights] Error fetching file '%s' in group %d: %v\n", fileName, groupID, err)
			return false, err
		}
		fmt.Printf("[HasDeletionRights] File '%s' fetched: owner=%s, requestedUser=%s\n", fileName, file.UserID, userID)
		if file.UserID == userID {
			fmt.Printf("[HasDeletionRights] User %s is the owner of file '%s'\n", userID, fileName)
			return true, nil
		}
		group, err := s.studyGroupService.GetStudyGroupByID(groupID)
		if err != nil {
			fmt.Printf("[HasDeletionRights] Error fetching study group %d: %v\n", groupID, err)
			return false, err
		}
		for _, member := range group.Members {
			if member.UserID == userID && member.Role == models.RoleAdmin {
				fmt.Printf("[HasDeletionRights] User %s is an admin in group %d\n", userID, groupID)
				return true, nil
			}
		}
		fmt.Printf("[HasDeletionRights] User %s does not have deletion rights for file '%s' in group %d\n", userID, fileName, groupID)
		return false, nil
	})
}
