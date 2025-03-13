package repositories

import (
	"database/sql"
	"errors"

	"gdp8-backend/internal/models"
)

var ErrFileAlreadyExists = errors.New("file already exists")
var ErrFileNotFound = errors.New("file not found")

type FileRepository interface {
	GetAllFiles(tx *sql.Tx) ([]models.File, error)
	GetAllFilesByGroupID(tx *sql.Tx, groupID models.StudyGroupID) ([]models.File, error)
	CreateFile(tx *sql.Tx, file models.File) (*models.File, error)
	GetFileByNameAndGroupID(tx *sql.Tx, fileName string, groupID models.StudyGroupID) (*models.File, error)
	DeleteFile(tx *sql.Tx, fileName string, groupID models.StudyGroupID) error
}

type SQLFileRepository struct {
}

func (s *SQLFileRepository) GetAllFiles(tx *sql.Tx) ([]models.File, error) {
	rows, err := tx.Query("SELECT name, file_owner_id, study_group_id FROM files")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var files []models.File
	for rows.Next() {
		var file models.File
		if err := rows.Scan(&file.Name, &file.UserID, &file.GroupID); err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return files, rows.Err()
}

func (s *SQLFileRepository) GetAllFilesByGroupID(tx *sql.Tx, groupID models.StudyGroupID) ([]models.File, error) {
	rows, err := tx.Query("SELECT name, file_owner_id, study_group_id FROM files WHERE study_group_id = ?", groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var files []models.File
	for rows.Next() {
		var file models.File
		if err := rows.Scan(&file.Name, &file.UserID, &file.GroupID); err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return files, rows.Err()
}

func (s *SQLFileRepository) CreateFile(tx *sql.Tx, file models.File) (*models.File, error) {
	var count int
	err := tx.QueryRow("SELECT COUNT(1) FROM files WHERE name = ? AND study_group_id = ?", file.Name, file.GroupID).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrFileAlreadyExists
	}
	_, err = tx.Exec("INSERT INTO files (name, file_owner_id, study_group_id) VALUES (?, ?, ?)",
		file.Name, file.UserID, file.GroupID)
	if err != nil {
		return nil, err
	}
	return &models.File{
		Name:    file.Name,
		UserID:  file.UserID,
		GroupID: file.GroupID,
	}, nil
}

func (s *SQLFileRepository) GetFileByNameAndGroupID(
	tx *sql.Tx,
	fileName string,
	groupID models.StudyGroupID) (*models.File, error) {
	var file models.File
	err := tx.QueryRow(`SELECT name, file_owner_id, study_group_id
        FROM files
        WHERE name = ? AND study_group_id = ?`, fileName, groupID).
		Scan(&file.Name, &file.UserID, &file.GroupID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrFileNotFound
		}
		return nil, err
	}
	return &file, nil
}

func (s *SQLFileRepository) DeleteFile(tx *sql.Tx, fileName string, groupID models.StudyGroupID) error {
	result, err := tx.Exec("DELETE FROM files WHERE name = ? AND study_group_id = ?", fileName, groupID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrFileNotFound
	}
	return nil
}
