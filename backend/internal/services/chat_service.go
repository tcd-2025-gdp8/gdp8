package services

import (
	"database/sql"
	"fmt"

	"gdp8-backend/internal/models"
	"gdp8-backend/internal/persistence"
	"gdp8-backend/internal/repositories"
)

type ChatService interface {
	GetChatMessagesByGroupId(studyGroupID models.StudyGroupID, userID models.UserID) ([]models.ChatMessageView, error)
	InsertMessage(studyGroupID models.StudyGroupID, message *models.ChatMessageDetails) error
}

var ErrUnauthorizedChatOperation = fmt.Errorf("unauthorized chat operation")

type chatServiceImpl struct {
	txMgr             persistence.TransactionManager
	chatRepo          repositories.ChatMessagesRepository
	studyGroupService StudyGroupService
}

func NewChatService(
	txMgr persistence.TransactionManager,
	chatRepo repositories.ChatMessagesRepository,
	studyGroupService StudyGroupService,
) ChatService {

	return &chatServiceImpl{
		txMgr:             txMgr,
		chatRepo:          chatRepo,
		studyGroupService: studyGroupService,
	}
}

func (s chatServiceImpl) GetChatMessagesByGroupId(studyGroupID models.StudyGroupID,
	userID models.UserID) ([]models.ChatMessageView, error) {

	isStudyGroupMember, err := s.studyGroupService.IsGroupMember(studyGroupID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat messages: failed to check user group membership: %w", err)
	}
	if !isStudyGroupMember {
		return nil, ErrUnauthorizedChatOperation
	}

	messages, err := persistence.WithTransaction(s.txMgr, func(tx *sql.Tx) ([]models.ChatMessageView, error) {
		return s.chatRepo.GetChatMessagesByGroupId(tx, studyGroupID)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get chat messages: %w", err)
	}
	return messages, nil
}

func (s chatServiceImpl) InsertMessage(studyGroupID models.StudyGroupID, message *models.ChatMessageDetails) error {
	isStudyGroupMember, err := s.studyGroupService.IsGroupMember(studyGroupID, message.UserID)
	if err != nil {
		return fmt.Errorf("failed to get chat messages: failed to check user group membership: %w", err)
	}
	if !isStudyGroupMember {
		return ErrUnauthorizedChatOperation
	}

	err = persistence.WithTransactionNoReturnVal(s.txMgr, func(tx *sql.Tx) error {
		return s.chatRepo.InsertMessage(tx, studyGroupID, message)
	})

	if err != nil {
		return fmt.Errorf("failed to insert message: %w", err)
	}
	return nil
}
