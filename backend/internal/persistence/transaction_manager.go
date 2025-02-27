package persistence

import (
	"database/sql"
	"errors"
	"fmt"
)

// TransactionManager defines the interface for managing database transactions
type TransactionManager interface {
	Begin() (*sql.Tx, error)
}

// WithTransaction executes the given function with a return value T within a transaction context.
// If the function encounters an error, the transaction is rolled back; otherwise, the transaction is committed.
func WithTransaction[T any](txMgr TransactionManager, fn func(*sql.Tx) (T, error)) (T, error) {
	var defaultT T
	if txMgr == nil {
		return defaultT, errors.New("transaction manager cannot be nil")
	}
	if fn == nil {
		return defaultT, errors.New("transaction function cannot be nil")
	}

	tx, err := txMgr.Begin()
	if err != nil {
		return defaultT, err
	}

	result, err := fn(tx)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			//nolint:errorlint // intentionally not wrapping rollbackErr
			return defaultT, fmt.Errorf("transaction error: %w; rollback error: %v", err, rollbackErr)
		}
		return defaultT, err
	}

	if err := tx.Commit(); err != nil {
		return defaultT, err
	}

	return result, nil
}

// WithTransactionNoReturnVal executes the given function with no return value within a transaction context.
// If the function encounters an error, the transaction is rolled back; otherwise, the transaction is committed.
func WithTransactionNoReturnVal(txMgr TransactionManager, fn func(*sql.Tx) error) error {
	_, err := WithTransaction(txMgr, func(tx *sql.Tx) (struct{}, error) {
		return struct{}{}, fn(tx)
	})
	return err
}
