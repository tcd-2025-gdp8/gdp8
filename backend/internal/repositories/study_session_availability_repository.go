package repositories

import (
	"database/sql"
	"fmt"

	"gdp8-backend/internal/models"
)

type StudySessionAvailabilityRepository interface {
	GetCurrentStudySessionAvailabilityRequests(tx *sql.Tx,
		studyGroupID models.StudyGroupID) ([]models.StudySessionAvailabilityRequest, error)
	CreateStudySessionAvailabilityRequests(tx *sql.Tx, studyGroupID models.StudyGroupID,
		availabilityRequestDetails *models.StudySessionAvailabilityRequestDetails) error
	DeleteStudySessionAvailabilityRequests(tx *sql.Tx,
		availabilityRequestID models.StudySessionAvailabilityRequestID) error
	UpsertUserAvailabilityEntries(tx *sql.Tx, availabilityRequestID models.StudySessionAvailabilityRequestID,
		userID models.UserID, availabilityEntries []models.AvailabilityEntry) error
}

type SQLStudySessionAvailabilityRepository struct {
}

func (s *SQLStudySessionAvailabilityRepository) GetCurrentStudySessionAvailabilityRequests(tx *sql.Tx,
	studyGroupID models.StudyGroupID) ([]models.StudySessionAvailabilityRequest, error) {

	query := `
		SELECT id, study_group_id, availability_period_start, availability_period_end
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
			&request.AvailabilityPeriodStart,
			&request.AvailabilityPeriodEnd,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan availability request: %w", err)
		}
		requests = append(requests, request)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating availability requests: %w", err)
	}

	return requests, nil
}

func (s *SQLStudySessionAvailabilityRepository) CreateStudySessionAvailabilityRequests(tx *sql.Tx,
	studyGroupID models.StudyGroupID, availabilityRequestDetails *models.StudySessionAvailabilityRequestDetails) error {

	query := `
		INSERT INTO study_session_availability_requests 
		(study_group_id, availability_period_start, availability_period_end)
		VALUES (?, ?, ?)`

	result, err := tx.Exec(query,
		studyGroupID,
		availabilityRequestDetails.AvailabilityPeriodStart,
		availabilityRequestDetails.AvailabilityPeriodEnd,
	)
	if err != nil {
		return fmt.Errorf("failed to create availability request: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("no rows were affected when creating availability request")
	}

	return nil
}

func (s *SQLStudySessionAvailabilityRepository) DeleteStudySessionAvailabilityRequests(tx *sql.Tx,
	availabilityRequestID models.StudySessionAvailabilityRequestID) error {

	query := `
		DELETE FROM study_session_availability_requests
		WHERE id = ?`

	result, err := tx.Exec(query, availabilityRequestID)
	if err != nil {
		return fmt.Errorf("failed to delete availability request: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("no availability request found with ID %d", availabilityRequestID)
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
