package services

import (
	"gdp8-backend/internal/repositories"
	"gdp8-backend/internal/persistence"
)

type ModuleService interface {
	GetAllModules() ([]repositories.Module, error)
	CreateModule(id, name string) error
}

type moduleServiceImpl struct {
	txManager   persistence.TransactionManager
	moduleRepo  repositories.ModuleRepository
}

func NewModuleService(txManager persistence.TransactionManager, moduleRepo repositories.ModuleRepository) ModuleService {
	return &moduleServiceImpl{txManager: txManager, moduleRepo: moduleRepo}
}

func (s *moduleServiceImpl) GetAllModules() ([]repositories.Module, error) {
	return persistence.WithTransaction(s.txManager, func(tx persistence.Transaction) ([]repositories.Module, error) {
		return s.moduleRepo.GetAllModules(tx)
	})
}

func (s *moduleServiceImpl) CreateModule(id, name string) error {
	return persistence.WithTransactionNoReturnVal(s.txManager, func(tx persistence.Transaction) error {
		return s.moduleRepo.AddModule(tx, repositories.Module{ID: id, Name: name})
	})
}
