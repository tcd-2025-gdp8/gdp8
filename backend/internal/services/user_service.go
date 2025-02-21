package services

import (
	"gdp8-backend/internal/models"
	"gdp8-backend/internal/repositories"
	"gdp8-backend/internal/persistence"
)

type UserService interface {
	GetUser(id string) (*models.User, error)
	CreateUser(user models.User) (*models.User, error)
	SetModules(userID string, modules []models.Module) error
}

type userServiceImpl struct {
	txManager persistence.TransactionManager
	userRepo  repositories.UserRepository
}

func NewUserService(txManager persistence.TransactionManager, userRepo repositories.UserRepository) UserService {
	return &userServiceImpl{
		txManager: txManager,
		userRepo:  userRepo,
	}
}

func (s *userServiceImpl) GetUser(id string) (*models.User, error) {
	return s.userRepo.GetUserByID(id)
}

func (s *userServiceImpl) CreateUser(user models.User) (*models.User, error) {
	return s.userRepo.CreateUser(user)
}

func (s *userServiceImpl) SetModules(userID string, modules []models.Module) error {
	return s.userRepo.SetUserModules(userID, modules)
}
