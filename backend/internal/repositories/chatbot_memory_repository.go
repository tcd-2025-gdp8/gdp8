package repositories

import (
	"database/sql"

	"gdp8-backend/internal/models"
)

type ChatBotMemoryRepo interface {
	GetFileContexts(tx *sql.Tx, studyGroupID models.StudyGroupID) (map[string]string, error)
	InsertMessage(tx *sql.Tx, studyGroupID models.StudyGroupID, message *models.FileContext) error
	DeleteMemory(tx *sql.Tx, studyGroupID models.StudyGroupID, contextSrc string) error
}

type SQLChatBotMemoryRepo struct{}

func (s *SQLChatBotMemoryRepo) GetFileContexts(
	tx *sql.Tx,
	studyGroupID models.StudyGroupID) (map[string]string, error) {
	rows, err := tx.Query("SELECT context_src, context_data FROM chatbot_memory WHERE study_group_id = ?",
		studyGroupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	contexts := make(map[string]string)
	for rows.Next() {
		var src string
		var data string
		if err := rows.Scan(&src, &data); err != nil {
			return nil, err
		}
		contexts[src] = data
	}
	return contexts, nil
}

func (s *SQLChatBotMemoryRepo) InsertMessage(tx *sql.Tx,
	studyGroupID models.StudyGroupID,
	message *models.FileContext) error {
	_, err := tx.Exec("INSERT INTO chatbot_memory (study_group_id, context_src, context_data) VALUES (?, ?, ?)", studyGroupID, message.Name, message.Data)
	return err
}

func (s *SQLChatBotMemoryRepo) DeleteMemory(tx *sql.Tx,
	studyGroupID models.StudyGroupID,
	contextSrc string) error {
	_, err := tx.Exec("DELETE FROM chatbot_memory WHERE study_group_id = ? AND context_src = ?", studyGroupID, contextSrc)
	return err
}
