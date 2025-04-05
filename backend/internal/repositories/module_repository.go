package repositories

import (
	"database/sql"
	"errors"

	"gdp8-backend/internal/models"
)

type ModuleRepository interface {
	GetAllModules(tx *sql.Tx) ([]models.Module, error)
	CreateModule(tx *sql.Tx, moduleDetails models.ModuleDetails) (*models.Module, error)
	GetModuleByID(tx *sql.Tx, id models.ModuleID) (*models.Module, error)
	GetModuleStudyGroupStats(tx *sql.Tx) (models.StudyGroupsMap, error)
}

var ErrModuleAlreadyExists = errors.New("module already exists")
var ErrModuleNotFound = errors.New("module not found")

type SQLModuleRepository struct {
}

func (s *SQLModuleRepository) FetchModules(tx *sql.Tx) ([]models.Module, error) {
	rows, err := tx.Query("SELECT id, code, name FROM modules")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var modules []models.Module
	for rows.Next() {
		var module models.Module
		if err := rows.Scan(&module.ID, &module.Code, &module.Name); err != nil {
			return nil, err
		}
		modules = append(modules, module)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return modules, nil
}

func (s *SQLModuleRepository) FetchStudyGroupsForModule(tx *sql.Tx,
	moduleID models.ModuleID) ([]models.StudyGroup, error) {
	groupRows, err := tx.Query(`
		SELECT id, name FROM study_groups WHERE module_id = ?
	`, moduleID)
	if err != nil {
		return nil, err
	}
	defer groupRows.Close()

	var groups []models.StudyGroup
	for groupRows.Next() {
		var group models.StudyGroup
		if err := groupRows.Scan(&group.ID, &group.Name); err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	if err := groupRows.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}

func (s *SQLModuleRepository) FetchStudyGroupMembersCount(tx *sql.Tx, group models.StudyGroup) (int, error) {
	var members int
	err := tx.QueryRow(`
		SELECT COUNT(*) FROM user_study_groups
		WHERE study_group_id = ? AND type IN ('admin', 'member')
	`, group.ID).Scan(&members)
	if err != nil {
		return 0, err
	}

	return members, nil
}

func (s *SQLModuleRepository) FetchTotalStudyTimeForGroup(tx *sql.Tx, group models.StudyGroup) (int64, error) {
	var totalMinutes int64
	err := tx.QueryRow(`
		SELECT COALESCE(SUM(duration_minutes), 0)
		FROM study_sessions
		WHERE study_group_id = ? AND start_time < NOW()
	`, group.ID).Scan(&totalMinutes)
	if err != nil {
		return 0, err
	}

	return totalMinutes, nil
}

func (s *SQLModuleRepository) FetchWeeklyStudyTimeForGroup(tx *sql.Tx, group models.StudyGroup) ([]int64, error) {
	weekly := make([]int64, 4)
	for i := range 4 {
		err := tx.QueryRow(`
			SELECT COALESCE(SUM(duration_minutes), 0)
			FROM study_sessions
			WHERE study_group_id = ?
			  AND start_time >= NOW() - INTERVAL ? WEEK
			  AND start_time < NOW() - INTERVAL ? WEEK
		`, group.ID, 4-i, 3-i).Scan(&weekly[i])
		if err != nil {
			return nil, err
		}
		weekly[i] /= 60 // Convert minutes to hours
	}

	return weekly, nil
}

func (s *SQLModuleRepository) GetAllModules(tx *sql.Tx) ([]models.Module, error) {
	return s.FetchModules(tx)
}

func (s *SQLModuleRepository) CreateModule(tx *sql.Tx, moduleDetails models.ModuleDetails) (*models.Module, error) {
	var count int
	err := tx.QueryRow("SELECT COUNT(1) FROM modules WHERE code = ?", moduleDetails.Code).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrModuleAlreadyExists
	}

	result, err := tx.Exec("INSERT INTO modules (code, name) VALUES (?, ?)",
		moduleDetails.Code, moduleDetails.Name)
	if err != nil {
		return nil, err
	}

	moduleID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &models.Module{
		ID:            models.ModuleID(moduleID),
		ModuleDetails: moduleDetails,
	}, nil
}

func (s *SQLModuleRepository) GetModuleByID(tx *sql.Tx, id models.ModuleID) (*models.Module, error) {
	var module models.Module
	err := tx.QueryRow("SELECT id, code, name FROM modules WHERE id = ?", id).
		Scan(&module.ID, &module.Code, &module.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrModuleNotFound
		}
		return nil, err
	}

	return &module, nil
}

func (s *SQLModuleRepository) GetModuleStudyGroupStats(tx *sql.Tx) (models.StudyGroupsMap, error) {
	modules, err := s.FetchModules(tx)
	if err != nil {
		return nil, err
	}

	statsMap := make(models.StudyGroupsMap)

	for _, module := range modules {
		groups, err := s.FetchStudyGroupsForModule(tx, module.ID)
		if err != nil {
			return nil, err
		}

		var groupStats []models.StudyGroupStatsDTO
		var idx int64 = 1

		for _, group := range groups {
			members, err := s.FetchStudyGroupMembersCount(tx, group)
			if err != nil {
				return nil, err
			}

			totalMinutes, err := s.FetchTotalStudyTimeForGroup(tx, group)
			if err != nil {
				return nil, err
			}

			weekly, err := s.FetchWeeklyStudyTimeForGroup(tx, group)
			if err != nil {
				return nil, err
			}

			groupStats = append(groupStats, models.StudyGroupStatsDTO{
				ID:          idx,
				Name:        group.Name,
				Members:     members,
				TotalHours:  totalMinutes / 60,
				WeeklyHours: weekly,
			})
			idx++
		}

		statsMap[module.Code] = groupStats
	}

	return statsMap, nil
}
