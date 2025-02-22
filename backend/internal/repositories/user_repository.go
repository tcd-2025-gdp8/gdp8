package repositories

import (
	"errors"
	"sync"

	"gdp8-backend/internal/models"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	GetUserByID(id string) (*models.User, error)
	CreateUser(user models.User) (*models.User, error)
	SetUserModules(id string, modules []models.Module) error
}

type MockUserRepository struct {
	mu    sync.Mutex
	users map[string]models.User
}

func NewMockUserRepository() UserRepository {
	return &MockUserRepository{
		users: make(map[string]models.User),
	}
}

func (r *MockUserRepository) GetUserByID(id string) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	user, ok := r.users[id]
	if !ok {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

func (r *MockUserRepository) CreateUser(user models.User) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.users[user.ID]; exists {
		return nil, errors.New("user already exists")
	}
	r.users[user.ID] = user
	return &user, nil
}

func (r *MockUserRepository) SetUserModules(id string, modules []models.Module) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	user, ok := r.users[id]
	if !ok {
		return ErrUserNotFound
	}
	user.Modules = modules
	r.users[id] = user
	return nil
}
