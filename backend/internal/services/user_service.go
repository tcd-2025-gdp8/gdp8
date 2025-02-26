package services

import (
	"database/sql"

	"gdp8-backend/internal/models"
	"gdp8-backend/internal/persistence"
	"gdp8-backend/internal/repositories"
)

type UserService interface {
	GetUser(id models.UserID) (*models.User, error)
	CreateUser(id models.UserID, userDetails models.UserDetails) (*models.User, error)
	SetModules(userID models.UserID, modules []models.ModuleID) error
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

func (s *userServiceImpl) GetUser(id models.UserID) (*models.User, error) {
	return persistence.WithTransaction(s.txManager, func(tx *sql.Tx) (*models.User, error) {
		return s.userRepo.GetUserByID(tx, id)
	})
}

func (s *userServiceImpl) CreateUser(id models.UserID, userDetails models.UserDetails) (*models.User, error) {
	return persistence.WithTransaction(s.txManager, func(tx *sql.Tx) (*models.User, error) {
		return s.userRepo.CreateUser(tx, id, userDetails)
	})
}

func (s *userServiceImpl) SetModules(userID models.UserID, modules []models.ModuleID) error {
	return persistence.WithTransactionNoReturnVal(s.txManager, func(tx *sql.Tx) error {
		return s.userRepo.SetUserModules(tx, userID, modules)
	})
}
