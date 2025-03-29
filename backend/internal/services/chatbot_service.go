package services

import (
	"database/sql"
	"fmt"

	"gdp8-backend/internal/models"
	"gdp8-backend/internal/persistence"
	"gdp8-backend/internal/repositories"
)

type ChatBotService interface {
	GetMemory(studyGroupID models.StudyGroupID, userID models.UserID) ([]models.ChatMessageView, error)
	StoreMemory(studyGroupID models.StudyGroupID, message *models.ChatMessageDetails) error
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

func (s *ChatBotServiceImpl) GetMemory(
	studyGroupID models.StudyGroupID,
	userID models.UserID) ([]models.ChatMessageView, error) {
	messages, err := s.chatService.GetChatMessagesByGroupID(studyGroupID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get memory: %w", err)
	}

	return messages, nil
}

func (s *ChatBotServiceImpl) StoreMemory(studyGroupID models.StudyGroupID, message *models.ChatMessageDetails) error {
	return persistence.WithTransactionNoReturnVal(s.txManager, func(tx *sql.Tx) error {
		return s.memoryRepo.InsertMessage(tx, studyGroupID, message)
	})
}
