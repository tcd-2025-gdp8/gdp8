package routes

import (
	"net/http"

	"firebase.google.com/go/v4/auth"

	"gdp8-backend/internal/handlers"
	"gdp8-backend/internal/middleware"
	"gdp8-backend/internal/persistence"
	"gdp8-backend/internal/repositories"
	"gdp8-backend/internal/services"
)

func RegisterUserRoutes(firebaseAuth *auth.Client) {
	txManager := persistence.MockTransactionManager{}
	userRepo := repositories.NewMockUserRepository()
	userService := services.NewUserService(&txManager, userRepo)
	handler := handlers.NewUserHandler(userService, ModuleRepo)

	http.HandleFunc("GET /api/user/{id}", middleware.WithFirebaseAuth(firebaseAuth, handler.GetUser))
	http.HandleFunc("POST /api/user/{id}/modules", middleware.WithFirebaseAuth(firebaseAuth, handler.SetModules))
	http.HandleFunc("GET /api/user/{id}/modules", middleware.WithFirebaseAuth(firebaseAuth, handler.GetModules))
	http.HandleFunc("POST /api/user", middleware.WithFirebaseAuth(firebaseAuth, handler.CreateUser))
}
