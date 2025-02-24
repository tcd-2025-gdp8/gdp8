package repositories

import (
	"database/sql"
	"encoding/json"
	"errors"

	"gdp8-backend/internal/models"
)

type StudyGroupRepository interface {
	GetStudyGroupByID(tx *sql.Tx, id models.StudyGroupID) (*models.StudyGroupView, error)
	GetAllStudyGroups(tx *sql.Tx) ([]models.StudyGroupView, error)
	CreateStudyGroup(tx *sql.Tx, studyGroupDetails models.StudyGroupDetails,
		adminUserID models.UserID) (*models.StudyGroupView, error)
	UpdateStudyGroupDetails(tx *sql.Tx, id models.StudyGroupID,
		details models.StudyGroupDetails) (*models.StudyGroupView, error)
	DeleteStudyGroup(tx *sql.Tx, id models.StudyGroupID) error
	UpdateStudyGroupMember(tx *sql.Tx, id models.StudyGroupID, userID models.UserID, role *models.StudyGroupRole) error
}

var ErrStudyGroupNotFound = errors.New("study group not found")

type SQLStudyGroupRepository struct {
}

func (s *SQLStudyGroupRepository) GetStudyGroupByID(tx *sql.Tx,
	id models.StudyGroupID) (*models.StudyGroupView, error) {
	query := `
		SELECT 
			s.id AS study_group_id,
			s.name,
			s.description,
			s.type,
			s.module_id,
			JSON_ARRAYAGG(
				JSON_OBJECT(
					'user_id', u.id,
					'name', u.name,
					'role', usg.type
				)
			) AS members
		FROM study_groups s
		LEFT JOIN user_study_groups usg ON s.id = usg.study_group_id
		LEFT JOIN users u ON usg.user_id = u.id
		WHERE s.id = ?
		GROUP BY s.id
	`

	var studyGroup models.StudyGroupView
	row := tx.QueryRow(query, id)
	var membersJSON string

	err := row.Scan(
		&studyGroup.ID,
		&studyGroup.StudyGroupDetails.Name,
		&studyGroup.StudyGroupDetails.Description,
		&studyGroup.StudyGroupDetails.Type,
		&studyGroup.StudyGroupDetails.ModuleID,
		&membersJSON,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrStudyGroupNotFound
		}
		return nil, err
	}

	//nolint:musttag
	err = json.Unmarshal([]byte(membersJSON), &studyGroup.Members)
	if err != nil {
		return nil, err
	}

	return &studyGroup, nil
}

func (s *SQLStudyGroupRepository) GetAllStudyGroups(tx *sql.Tx) ([]models.StudyGroupView, error) {

	query := `
		SELECT 
			s.id AS study_group_id,
			s.name,
			s.description,
			s.type,
			s.module_id,
			JSON_ARRAYAGG(
				JSON_OBJECT(
					'user_id', u.id,
					'name', u.name,
					'role', usg.type
				)
			) AS members
		FROM study_groups s
		LEFT JOIN user_study_groups usg ON s.id = usg.study_group_id
		LEFT JOIN users u ON usg.user_id = u.id
		GROUP BY s.id
	`

	rows, err := tx.Query(query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var studyGroups []models.StudyGroupView

	for rows.Next() {
		var studyGroup models.StudyGroupView
		var membersJSON string

		err := rows.Scan(
			&studyGroup.ID,
			&studyGroup.StudyGroupDetails.Name,
			&studyGroup.StudyGroupDetails.Description,
			&studyGroup.StudyGroupDetails.Type,
			&studyGroup.StudyGroupDetails.ModuleID,
			&membersJSON,
		)
		if err != nil {
			return nil, err
		}

		//nolint:musttag
		err = json.Unmarshal([]byte(membersJSON), &studyGroup.Members)
		if err != nil {
			return nil, err
		}

		studyGroups = append(studyGroups, studyGroup)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return studyGroups, nil
}

func (s *SQLStudyGroupRepository) CreateStudyGroup(tx *sql.Tx,
	studyGroupDetails models.StudyGroupDetails, adminUserID models.UserID) (*models.StudyGroupView, error) {
	queryStudyGroup := `
		INSERT INTO study_groups (name, description, type, module_id, max_members)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := tx.Exec(
		queryStudyGroup,
		studyGroupDetails.Name,
		studyGroupDetails.Description,
		studyGroupDetails.Type,
		studyGroupDetails.ModuleID,
		10, // TODO fix max group members
	)
	if err != nil {
		return nil, err
	}

	studyGroupID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	queryUserStudyGroup := `
		INSERT INTO user_study_groups (user_id, study_group_id, type)
		VALUES (?, ?, ?)
	`

	_, err = tx.Exec(queryUserStudyGroup, adminUserID, studyGroupID, "admin")
	if err != nil {
		return nil, err
	}

	return s.GetStudyGroupByID(tx, models.StudyGroupID(studyGroupID))
}

func (s *SQLStudyGroupRepository) UpdateStudyGroupDetails(tx *sql.Tx,
	id models.StudyGroupID, details models.StudyGroupDetails) (*models.StudyGroupView, error) {

	query := `
		UPDATE study_groups
		SET name = ?, description = ?, type = ?, module_id = ?, max_members = ?
		WHERE id = ?
	`

	_, err := tx.Exec(query,
		details.Name,
		details.Description,
		details.Type,
		details.ModuleID,
		10, // TODO fix max study group members
		id,
	)
	if err != nil {
		return nil, err
	}

	return s.GetStudyGroupByID(tx, id)
}

func (s *SQLStudyGroupRepository) DeleteStudyGroup(tx *sql.Tx, id models.StudyGroupID) error {
	query := `
		DELETE FROM study_groups
		WHERE id = ?
	`

	_, err := tx.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

// UpdateStudyGroupMember updates the role of a member in a study group or removes the member if the role is nil.
// The operation is idempotent - it will not return an error if the user already has the requested role.
// Returns an error if the study group does not exist.
func (s *SQLStudyGroupRepository) UpdateStudyGroupMember(tx *sql.Tx,
	id models.StudyGroupID, userID models.UserID, role *models.StudyGroupRole) error {
	checkQuery := `
		SELECT COUNT(1)
		FROM study_groups
		WHERE id = ?
	`
	var count int
	err := tx.QueryRow(checkQuery, id).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrStudyGroupNotFound
	}

	if role == nil {
		deleteQuery := `
			DELETE FROM user_study_groups
			WHERE study_group_id = ? AND user_id = ?
		`
		_, err = tx.Exec(deleteQuery, id, userID)
		return err
	}

	upsertQuery := `
		INSERT INTO user_study_groups (user_id, study_group_id, type)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE type = VALUES(type)
	`
	_, err = tx.Exec(upsertQuery, userID, id, role)
	return err
}
