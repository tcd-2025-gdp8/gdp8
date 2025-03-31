package services

import (
	"database/sql"

	"gdp8-backend/internal/models"
	"gdp8-backend/internal/persistence"
	"gdp8-backend/internal/repositories"
)

type ChatBotService interface {
	GetMemory(studyGroupID models.StudyGroupID) (map[string]string, error)
	StoreMemory(studyGroupID models.StudyGroupID, ctx *models.FileContext) error
	DeleteMemory(studyGroupID models.StudyGroupID, contextSrc string) error
}

type ChatBotServiceImpl struct {
	txManager         persistence.TransactionManager
	studyGroupService StudyGroupService
	fileService       FileService
	chatService       ChatService
	memoryRepo        repositories.ChatBotMemoryRepo
}

func NewChatBotService(
	txManager persistence.TransactionManager,
	studyGroupService StudyGroupService,
	fileService FileService,
	chatService ChatService,
	memoryRepo repositories.ChatBotMemoryRepo,
) ChatBotService {
	return &ChatBotServiceImpl{
		txManager:         txManager,
		studyGroupService: studyGroupService,
		fileService:       fileService,
		chatService:       chatService,
		memoryRepo:        memoryRepo,
	}
}

func (s *ChatBotServiceImpl) GetMemory(studyGroupID models.StudyGroupID) (map[string]string, error) {
	return persistence.WithTransaction(s.txManager, func(tx *sql.Tx) (map[string]string, error) {
		return s.memoryRepo.GetFileContexts(tx, studyGroupID)
	})
}

func (s *ChatBotServiceImpl) StoreMemory(
	studyGroupID models.StudyGroupID,
	ctx *models.FileContext) error {
	return persistence.WithTransactionNoReturnVal(s.txManager, func(tx *sql.Tx) error {
		return s.memoryRepo.InsertMessage(tx, studyGroupID, ctx)
	})
}

func (s *ChatBotServiceImpl) DeleteMemory(studyGroupID models.StudyGroupID, contextSrc string) error {
	return persistence.WithTransactionNoReturnVal(s.txManager, func(tx *sql.Tx) error {
		return s.memoryRepo.DeleteMemory(tx, studyGroupID, contextSrc)
	})
}
