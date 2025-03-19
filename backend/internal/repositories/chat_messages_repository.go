package repositories

import (
	"database/sql"

	"gdp8-backend/internal/models"
)

type ChatMessagesRepository interface {
	GetChatMessagesByGroupID(tx *sql.Tx, studyGroupID models.StudyGroupID) ([]models.ChatMessageView, error)
	InsertMessage(tx *sql.Tx, studyGroupID models.StudyGroupID, message *models.ChatMessageDetails) error
}

type SQLChatMessagesRepository struct {
}

func (s *SQLChatMessagesRepository) GetChatMessagesByGroupID(tx *sql.Tx,
	studyGroupID models.StudyGroupID) ([]models.ChatMessageView, error) {

	rows, err := tx.Query(
		`SELECT m.id, m.user_id, u.name as user_name, m.text, m.timestamp
		FROM study_group_chat_messages m
		JOIN users u ON m.user_id = u.id
		WHERE study_group_id = ?
		ORDER BY timestamp ASC
		LIMIT 1000`, studyGroupID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var messages []models.ChatMessageView
	for rows.Next() {
		var message models.ChatMessageView
		if err := rows.Scan(&message.ID, &message.UserID, &message.UserName, &message.Text, &message.Timestamp); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (s *SQLChatMessagesRepository) InsertMessage(tx *sql.Tx,
	studyGroupID models.StudyGroupID, message *models.ChatMessageDetails) error {

	_, err := tx.Exec(
		`INSERT INTO study_group_chat_messages (study_group_id, user_id, text, timestamp)
		VALUES (?, ?, ?, ?)`,
		studyGroupID, message.UserID, message.Text, message.Timestamp)

	return err
}
