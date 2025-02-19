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

func RegisterModuleRoutes(firebaseAuth *auth.Client) {
	txManager := persistence.MockTransactionManager{}
	moduleRepo := repositories.NewMockModuleRepository() // Mocked for now
	moduleService := services.NewModuleService(&txManager, moduleRepo)
	handler := handlers.NewModuleHandler(moduleService)

	http.HandleFunc("GET /api/modules", middleware.WithFirebaseAuth(firebaseAuth, handler.GetAllModules))
	http.HandleFunc("POST /api/modules", middleware.WithFirebaseAuth(firebaseAuth, handler.CreateModule))
}
