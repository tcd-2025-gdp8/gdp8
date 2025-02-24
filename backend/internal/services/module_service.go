package services

import (
	"database/sql"

	"gdp8-backend/internal/models"
	"gdp8-backend/internal/persistence"
	"gdp8-backend/internal/repositories"
)

type ModuleService interface {
	GetAllModules() ([]models.Module, error)
	CreateModule(moduleDetails models.ModuleDetails) (*models.Module, error)
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
	return persistence.WithTransaction(s.txManager, func(tx *sql.Tx) ([]models.Module, error) {
		return s.moduleRepo.GetAllModules(tx)
	})
}

func (s *moduleServiceImpl) CreateModule(moduleDetails models.ModuleDetails) (*models.Module, error) {
	return persistence.WithTransaction(s.txManager, func(tx *sql.Tx) (*models.Module, error) {
		return s.moduleRepo.CreateModule(tx, moduleDetails)
	})
}
