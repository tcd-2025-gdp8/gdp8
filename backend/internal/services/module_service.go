package services

import (
	"gdp8-backend/internal/models"
	"gdp8-backend/internal/persistence"
	"gdp8-backend/internal/repositories"
)

type ModuleService interface {
	GetAllModules() ([]models.Module, error)
	CreateModule(id, name string) error
}

type moduleServiceImpl struct {
	txManager  persistence.TransactionManager
	moduleRepo repositories.ModuleRepository
}

func NewModuleService(txManager persistence.TransactionManager,
	moduleRepo repositories.ModuleRepository) ModuleService {
	return &moduleServiceImpl{txManager: txManager, moduleRepo: moduleRepo}
}

func (s *moduleServiceImpl) GetAllModules() ([]models.Module, error) {
	return persistence.WithTransaction(s.txManager, func(tx persistence.Transaction) ([]models.Module, error) {
		return s.moduleRepo.GetAllModules(tx)
	})
}

func (s *moduleServiceImpl) CreateModule(id, name string) error {
	return persistence.WithTransactionNoReturnVal(s.txManager, func(tx persistence.Transaction) error {
		return s.moduleRepo.AddModule(tx, models.Module{ID: id, Name: name})
	})
}
