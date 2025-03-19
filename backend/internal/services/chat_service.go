package services

import (
	"database/sql"
	"fmt"

	"gdp8-backend/internal/models"
	"gdp8-backend/internal/persistence"
	"gdp8-backend/internal/repositories"
)

type ChatService interface {
	GetChatMessagesByGroupId(studyGroupID models.StudyGroupID) ([]models.ChatMessageView, error)
	InsertMessage(studyGroupID models.StudyGroupID, message *models.ChatMessageDetails) error
}

type chatServiceImpl struct {
	txMgr    persistence.TransactionManager
	chatRepo repositories.ChatMessagesRepository
}

func NewChatService(txMgr persistence.TransactionManager, chatRepo repositories.ChatMessagesRepository) ChatService {
	return &chatServiceImpl{
		txMgr:    txMgr,
		chatRepo: chatRepo,
	}
}

func (s chatServiceImpl) GetChatMessagesByGroupId(studyGroupID models.StudyGroupID) ([]models.ChatMessageView, error) {
	messages, err := persistence.WithTransaction(s.txMgr, func(tx *sql.Tx) ([]models.ChatMessageView, error) {
		return s.chatRepo.GetChatMessagesByGroupId(tx, studyGroupID)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get chat messages: %w", err)
	}
	return messages, nil
}

func (s chatServiceImpl) InsertMessage(studyGroupID models.StudyGroupID, message *models.ChatMessageDetails) error {
	err := persistence.WithTransactionNoReturnVal(s.txMgr, func(tx *sql.Tx) error {
		return s.chatRepo.InsertMessage(tx, studyGroupID, message)
	})

	if err != nil {
		return fmt.Errorf("failed to insert message: %w", err)
	}
	return nil
}
