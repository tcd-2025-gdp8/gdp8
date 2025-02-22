package repositories

import (
	"errors"
	"sync"

	"gdp8-backend/internal/models"
	"gdp8-backend/internal/persistence"
)

type ModuleRepository interface {
	GetAllModules(tx persistence.Transaction) ([]models.Module, error)
	AddModule(tx persistence.Transaction, module models.Module) error
	GetModuleByID(tx persistence.Transaction, id string) (models.Module, error)
}

var ErrModuleAlreadyExists = errors.New("module already exists")

type MockModuleRepository struct {
	modules map[string]models.Module
	mu      sync.Mutex
}

func NewMockModuleRepository() ModuleRepository {
	return &MockModuleRepository{
		modules: map[string]models.Module{
			"CSU44052": {ID: "CSU44052", Name: "Computer Graphics"},
			"CSU44061": {ID: "CSU44061", Name: "Machine Learning"},
			"CSU44051": {ID: "CSU44051", Name: "Human Factors"},
			"CSU44000": {ID: "CSU44000", Name: "Internet Applications"},
			"CSU44012": {ID: "CSU44012", Name: "Topics in Functional Programming"},
			"CSU44099": {ID: "CSU44099", Name: "Final Year Project"},
			"CSU44098": {ID: "CSU44098", Name: "Group Design Project"},
			"CSU44081": {ID: "CSU44081", Name: "Entrepreneurship & High Tech Venture Creation"},
		},
	}
}

func (r *MockModuleRepository) GetAllModules(_ persistence.Transaction) ([]models.Module, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	modulesList := make([]models.Module, 0, len(r.modules))
	for _, module := range r.modules {
		modulesList = append(modulesList, module)
	}
	return modulesList, nil
}

func (r *MockModuleRepository) AddModule(_ persistence.Transaction, module models.Module) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.modules[module.ID]; exists {
		return ErrModuleAlreadyExists
	}

	r.modules[module.ID] = module
	return nil
}

func (r *MockModuleRepository) GetModuleByID(_ persistence.Transaction, id string) (models.Module, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	mod, ok := r.modules[id]
	if !ok {
		return models.Module{}, errors.New("module not found")
	}
	return mod, nil
}
