package repositories

import (
	"database/sql"
	"errors"

	"gdp8-backend/internal/models"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	GetUserByID(tx *sql.Tx, id models.UserID) (*models.User, error)
	CreateUser(tx *sql.Tx, id models.UserID, userDetails models.UserDetails) (*models.User, error)
	SetUserModules(tx *sql.Tx, id models.UserID, modules []models.ModuleID) error
}

type SQLUserRepository struct {
}

func (s *SQLUserRepository) GetUserByID(tx *sql.Tx, id models.UserID) (*models.User, error) {

	row := tx.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id)

	var user models.User
	if err := row.Scan(&user.ID, &user.Name, &user.Email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	rows, err := tx.Query(`
		SELECT m.id, m.code, m.name
		FROM user_modules um
		INNER JOIN modules m ON um.module_id = m.id
		WHERE um.user_id = ?
	`, id)
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
	if err = rows.Err(); err != nil {
		return nil, err
	}

	user.Modules = modules
	return &user, nil
}

func (s *SQLUserRepository) CreateUser(tx *sql.Tx,
	id models.UserID, userDetails models.UserDetails) (*models.User, error) {
	_, err := tx.Exec("INSERT INTO users (id, name, email) VALUES (?, ?, ?)",
		id, userDetails.Name, userDetails.Email)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID: id,
		UserDetails: models.UserDetails{
			Name:  userDetails.Name,
			Email: userDetails.Email,
		},
		Modules: []models.Module{},
	}
	return user, nil
}

func (s *SQLUserRepository) SetUserModules(tx *sql.Tx, id models.UserID, modules []models.ModuleID) error {
	_, err := tx.Exec("DELETE FROM user_modules WHERE user_id = ?", id)
	if err != nil {
		return err
	}

	for _, moduleID := range modules {
		_, err := tx.Exec("INSERT INTO user_modules (user_id, module_id) VALUES (?, ?)", id, moduleID)
		if err != nil {
			return err
		}
	}

	return nil
}
