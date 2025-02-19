package repositories

import (
	"errors"
	"gdp8-backend/internal/persistence"
	"sync"
)

type Module struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ModuleRepository interface {
	GetAllModules(tx persistence.Transaction) ([]Module, error)
	AddModule(tx persistence.Transaction, module Module) error
}

var ErrModuleAlreadyExists = errors.New("module already exists")

type MockModuleRepository struct {
	modules map[string]Module
	mu      sync.Mutex
}

func NewMockModuleRepository() ModuleRepository {
	return &MockModuleRepository{
		modules: map[string]Module{
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

func (r *MockModuleRepository) GetAllModules(_ persistence.Transaction) ([]Module, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	modulesList := make([]Module, 0, len(r.modules))
	for _, module := range r.modules {
		modulesList = append(modulesList, module)
	}
	return modulesList, nil
}

func (r *MockModuleRepository) AddModule(_ persistence.Transaction, module Module) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.modules[module.ID]; exists {
		return ErrModuleAlreadyExists
	}

	r.modules[module.ID] = module
	return nil
}
