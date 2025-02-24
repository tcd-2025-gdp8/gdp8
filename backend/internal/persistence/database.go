package persistence

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
)

type SqlTransactionManager struct {
	DB *sql.DB
}

func NewSqlTransactionManager(db *sql.DB) *SqlTransactionManager {
	return &SqlTransactionManager{DB: db}
}

func (s SqlTransactionManager) Begin() (*sql.Tx, error) {
	return s.DB.Begin()
}

func OpenDB() (*sql.DB, error) {
	dsn := os.Getenv("MYSQL_CONNECTION_STRING")
	if dsn == "" {
		return nil, fmt.Errorf("MYSQL_CONNECTION_STRING is not set")
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func ExecuteMigrations(db *sql.DB) error {
	migrationsDir := "migrations"

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".sql" {
			content, err := os.ReadFile(filepath.Join(migrationsDir, file.Name()))
			if err != nil {
				return fmt.Errorf("failed to read file %s: %w", file.Name(), err)
			}

			_, err = db.Exec(string(content))
			if err != nil {
				return fmt.Errorf("failed to execute migration %s: %w", file.Name(), err)
			}
		}
	}

	return nil
}
