package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"gdp8-backend/internal/models"
)

type StudySessionRepository interface {
	GetAllStudySessionsByStudyGroup(tx *sql.Tx, studyGroupID models.StudyGroupID) ([]models.StudySession, error)
	GetAllStudySessionsByUser(tx *sql.Tx, userID models.UserID) ([]models.StudySession, error)
	GetStudySession(tx *sql.Tx, studySessionID models.StudySessionID) (*models.StudySession, error)
	CreateStudySession(tx *sql.Tx, studyGroupID models.StudyGroupID,
		creatorID models.UserID, studySessionDetails *models.StudySessionDetails) (*models.StudySession, error)
	UpdateStudySession(tx *sql.Tx, studySessionID models.StudySessionID,
		studySessionDetails *models.StudySessionDetails) (*models.StudySession, error)
	DeleteStudySession(tx *sql.Tx, studySessionID models.StudySessionID) error
	GetUpcomingSessions(
		tx *sql.Tx,
		windowStart time.Time,
		windowEnd time.Time,
	) ([]models.StudySession, error)
}

var ErrStudySessionNotFound = errors.New("study session not found")

type SQLStudySessionRepository struct {
}

func (s *SQLStudySessionRepository) GetAllStudySessionsByStudyGroup(tx *sql.Tx,
	studyGroupID models.StudyGroupID) ([]models.StudySession, error) {

	query := `
        SELECT id, study_group_id, creator_id, title, start_time, duration_minutes, end_time
        FROM study_sessions
        WHERE study_group_id = ?
        ORDER BY start_time DESC`

	rows, err := tx.Query(query, studyGroupID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving study sessions: %w", err)
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var sessions []models.StudySession
	for rows.Next() {
		var session models.StudySession
		err := rows.Scan(
			&session.ID,
			&session.StudyGroupID,
			&session.CreatorID,
			&session.Title,
			&session.StartTime,
			&session.DurationMinutes,
			&session.EndTime,
		)
		if err != nil {
			return nil, fmt.Errorf("error retrieving study session: %w", err)
		}
		sessions = append(sessions, session)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error retrieving study sessions: %w", err)
	}

	return sessions, nil
}

func (s *SQLStudySessionRepository) GetAllStudySessionsByUser(_ *sql.Tx,
	_ models.UserID) ([]models.StudySession, error) {

	// TODO implement
	panic("not implemented")
}

func (s *SQLStudySessionRepository) GetStudySession(tx *sql.Tx,
	studySessionID models.StudySessionID) (*models.StudySession, error) {

	query := `
        SELECT id, study_group_id, creator_id, title, start_time, duration_minutes, end_time
        FROM study_sessions
        WHERE id = ?`

	var session models.StudySession
	err := tx.QueryRow(query, studySessionID).Scan(
		&session.ID,
		&session.StudyGroupID,
		&session.CreatorID,
		&session.Title,
		&session.StartTime,
		&session.DurationMinutes,
		&session.EndTime,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrStudySessionNotFound
		}
		return nil, fmt.Errorf("error retrieving study session: %w", err)
	}

	return &session, nil
}

func (s *SQLStudySessionRepository) CreateStudySession(tx *sql.Tx, studyGroupID models.StudyGroupID,
	creatorID models.UserID, studySessionDetails *models.StudySessionDetails) (*models.StudySession, error) {

	query := `
        INSERT INTO study_sessions (study_group_id, creator_id, title, start_time, duration_minutes)
        VALUES (?, ?, ?, ?, ?)`

	result, err := tx.Exec(query,
		studyGroupID,
		creatorID,
		studySessionDetails.Title,
		studySessionDetails.StartTime,
		studySessionDetails.DurationMinutes)
	if err != nil {
		return nil, fmt.Errorf("error creating study session: %w", err)
	}

	studySessionID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error creating study session: %w", err)
	}

	return s.GetStudySession(tx, models.StudySessionID(studySessionID))
}

func (s *SQLStudySessionRepository) UpdateStudySession(tx *sql.Tx, studySessionID models.StudySessionID,
	studySessionDetails *models.StudySessionDetails) (*models.StudySession, error) {

	query := `
        UPDATE study_sessions 
        SET title = ?, start_time = ?, duration_minutes = ?
        WHERE id = ?`

	_, err := tx.Exec(query,
		studySessionDetails.Title,
		studySessionDetails.StartTime,
		studySessionDetails.DurationMinutes,
		studySessionID)

	if err != nil {
		return nil, fmt.Errorf("error updating study session: %w", err)
	}

	return s.GetStudySession(tx, studySessionID)
}

func (s *SQLStudySessionRepository) DeleteStudySession(tx *sql.Tx, studySessionID models.StudySessionID) error {
	query := `
        DELETE FROM study_sessions 
        WHERE id = ?`

	_, err := tx.Exec(query, studySessionID)

	if err != nil {
		return fmt.Errorf("error deleting study session: %w", err)
	}

	return nil
}

func (s *SQLStudySessionRepository) GetUpcomingSessions(
	tx *sql.Tx,
	windowStart time.Time,
	windowEnd time.Time,
) ([]models.StudySession, error) {
	query := `
        SELECT id, study_group_id, creator_id, title, start_time, duration_minutes, end_time
        FROM study_sessions
        WHERE start_time >= ? AND start_time < ?`
	rows, err := tx.Query(query, windowStart, windowEnd)
	if err != nil {
		return nil, fmt.Errorf("error querying upcoming sessions: %w", err)
	}
	defer rows.Close()

	var sessions []models.StudySession
	for rows.Next() {
		var session models.StudySession
		err := rows.Scan(
			&session.ID,
			&session.StudyGroupID,
			&session.CreatorID,
			&session.Title,
			&session.StartTime,
			&session.DurationMinutes,
			&session.EndTime,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning session: %w", err)
		}
		sessions = append(sessions, session)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sessions: %w", err)
	}

	return sessions, nil
}
