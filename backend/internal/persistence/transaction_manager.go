package persistence

import (
	"errors"
	"fmt"
)

// TransactionManager defines the interface for managing database transactions
type TransactionManager interface {
	Begin() (Transaction, error)
}

type Transaction interface {
	Commit() error
	Rollback() error
}

// WithTransaction executes the given function with a return value T within a transaction context.
// If the function encounters an error, the transaction is rolled back; otherwise, the transaction is committed.
func WithTransaction[T any](txMgr TransactionManager, fn func(Transaction) (T, error)) (T, error) {
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
func WithTransactionNoReturnVal(txMgr TransactionManager, fn func(Transaction) error) error {
	_, err := WithTransaction(txMgr, func(tx Transaction) (struct{}, error) {
		return struct{}{}, fn(tx)
	})
	return err
}

type MockTransactionManager struct {
}

func (m *MockTransactionManager) Begin() (Transaction, error) {
	return &MockTransaction{}, nil
}

type MockTransaction struct {
}

func (m *MockTransaction) Commit() error {
	return nil
}

func (m *MockTransaction) Rollback() error {
	return nil
}
