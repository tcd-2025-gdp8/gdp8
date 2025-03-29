package repositories

import (
	"database/sql"

	"gdp8-backend/internal/models"
)

type ChatBotMemoryRepo interface {
	GetChatMessagesByGroupID(tx *sql.Tx, studyGroupID models.StudyGroupID) ([]models.ChatMessageView, error)
	InsertMessage(tx *sql.Tx, studyGroupID models.StudyGroupID, message *models.ChatMessageDetails) error
}
