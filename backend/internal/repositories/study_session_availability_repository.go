package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"gdp8-backend/internal/models"
)

type StudySessionAvailabilityRepository interface {
	GetCurrentStudySessionAvailabilityRequests(tx *sql.Tx,
		studyGroupID models.StudyGroupID) ([]models.StudySessionAvailabilityRequest, error)
	GetStudySessionAvailabilityRequest(tx *sql.Tx,
		id models.StudySessionAvailabilityRequestID) (*models.StudySessionAvailabilityRequest, error)
	CreateStudySessionAvailabilityRequest(tx *sql.Tx, studyGroupID models.StudyGroupID,
		creatorID models.UserID, availabilityRequestDetails *models.StudySessionAvailabilityRequestDetails) error
	DeleteStudySessionAvailabilityRequest(tx *sql.Tx,
		availabilityRequestID models.StudySessionAvailabilityRequestID) error
	UpsertUserAvailabilityEntries(tx *sql.Tx, availabilityRequestID models.StudySessionAvailabilityRequestID,
		userID models.UserID, availabilityEntries []models.AvailabilityEntry) error
}

type SQLStudySessionAvailabilityRepository struct {
}

func (s *SQLStudySessionAvailabilityRepository) GetCurrentStudySessionAvailabilityRequests(tx *sql.Tx,
	studyGroupID models.StudyGroupID) ([]models.StudySessionAvailabilityRequest, error) {

	query := `
		SELECT id, study_group_id, creator_id, availability_period_start, availability_period_end
		FROM current_study_session_availability_requests
		WHERE study_group_id = ?
		ORDER BY availability_period_start`

	rows, err := tx.Query(query, studyGroupID)
	if err != nil {
		return nil, fmt.Errorf("failed to query current availability requests: %w", err)
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var requests []models.StudySessionAvailabilityRequest
	for rows.Next() {
		var request models.StudySessionAvailabilityRequest
		err := rows.Scan(
			&request.ID,
			&request.StudyGroupID,
			&request.CreatorID,
			&request.AvailabilityPeriodStart,
			&request.AvailabilityPeriodEnd,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan availability request: %w", err)
		}

		requests = append(requests, request)
	}

	for i := range requests {
		entries, err := s.getEntriesForRequest(tx, requests[i].ID)
		if err != nil {
			return nil, err
		}
		requests[i].Entries = entries
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating availability requests: %w", err)
	}

	return requests, nil
}

func (s *SQLStudySessionAvailabilityRepository) GetStudySessionAvailabilityRequest(tx *sql.Tx,
	id models.StudySessionAvailabilityRequestID) (*models.StudySessionAvailabilityRequest, error) {

	query := `
		SELECT id, study_group_id, creator_id, availability_period_start, availability_period_end
		FROM study_session_availability_requests
		WHERE id = ?`

	var request models.StudySessionAvailabilityRequest
	err := tx.QueryRow(query, id).Scan(
		&request.ID,
		&request.StudyGroupID,
		&request.CreatorID,
		&request.AvailabilityPeriodStart,
		&request.AvailabilityPeriodEnd,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan availability request: %w", err)
	}

	entries, err := s.getEntriesForRequest(tx, request.ID)
	if err != nil {
		return nil, err
	}

	request.Entries = entries
	return &request, nil
}

func (s *SQLStudySessionAvailabilityRepository) CreateStudySessionAvailabilityRequest(tx *sql.Tx,
	studyGroupID models.StudyGroupID, creatorID models.UserID,
	availabilityRequestDetails *models.StudySessionAvailabilityRequestDetails) error {

	query := `
		INSERT INTO study_session_availability_requests 
		(study_group_id, creator_id, availability_period_start, availability_period_end)
		VALUES (?, ?, ?, ?)`

	_, err := tx.Exec(query,
		studyGroupID,
		creatorID,
		availabilityRequestDetails.AvailabilityPeriodStart,
		availabilityRequestDetails.AvailabilityPeriodEnd,
	)
	if err != nil {
		return fmt.Errorf("failed to create availability request: %w", err)
	}

	return nil
}

func (s *SQLStudySessionAvailabilityRepository) DeleteStudySessionAvailabilityRequest(tx *sql.Tx,
	availabilityRequestID models.StudySessionAvailabilityRequestID) error {

	query := `
		DELETE FROM study_session_availability_requests
		WHERE id = ?`

	_, err := tx.Exec(query, availabilityRequestID)
	if err != nil {
		return fmt.Errorf("failed to delete availability request: %w", err)
	}

	return nil
}

func (s *SQLStudySessionAvailabilityRepository) UpsertUserAvailabilityEntries(tx *sql.Tx,
	availabilityRequestID models.StudySessionAvailabilityRequestID, userID models.UserID,
	availabilityEntries []models.AvailabilityEntry) error {

	deleteQuery := `
		DELETE FROM study_session_availability_entries
		WHERE availability_request_id = ? AND user_id = ?`

	_, err := tx.Exec(deleteQuery, availabilityRequestID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete existing availability entries: %w", err)
	}

	if len(availabilityEntries) == 0 {
		return nil
	}

	insertQuery := `
		INSERT INTO study_session_availability_entries
		(availability_request_id, user_id, availability_entry_start, availability_entry_end)
		VALUES (?, ?, ?, ?)`

	stmt, err := tx.Prepare(insertQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement for availability entries: %w", err)
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	for _, entry := range availabilityEntries {
		_, err = stmt.Exec(
			availabilityRequestID,
			userID,
			entry.AvailabilityEntryStart,
			entry.AvailabilityEntryEnd,
		)
		if err != nil {
			return fmt.Errorf("failed to insert availability entry: %w", err)
		}
	}

	return nil
}

func (s *SQLStudySessionAvailabilityRepository) getEntriesForRequest(tx *sql.Tx,
	requestID models.StudySessionAvailabilityRequestID) ([]models.StudySessionAvailabilityEntry, error) {

	entryQuery := `
		SELECT user_id, availability_entry_start, availability_entry_end
		FROM study_session_availability_entries
		WHERE availability_request_id = ?`
	entryRows, err := tx.Query(entryQuery, requestID)
	if err != nil {
		return nil, fmt.Errorf("failed to query availability entries: %w", err)
	}

	var entries []models.StudySessionAvailabilityEntry
	for entryRows.Next() {
		var entry models.StudySessionAvailabilityEntry
		err := entryRows.Scan(
			&entry.UserID,
			&entry.AvailabilityEntryStart,
			&entry.AvailabilityEntryEnd,
		)
		if err != nil {
			_ = entryRows.Close()
			return nil, fmt.Errorf("failed to scan availability entry: %w", err)
		}
		entries = append(entries, entry)
	}
	_ = entryRows.Close()

	if err = entryRows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating availability entries: %w", err)
	}

	return entries, nil
}
