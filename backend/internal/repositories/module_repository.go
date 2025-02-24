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
}

var ErrModuleAlreadyExists = errors.New("module already exists")
var ErrModuleNotFound = errors.New("module not found")

type SQLModuleRepository struct {
}

func (s *SQLModuleRepository) GetAllModules(tx *sql.Tx) ([]models.Module, error) {

	rows, err := tx.Query("SELECT id, code, name FROM modules")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

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
