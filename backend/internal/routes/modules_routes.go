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

// TODO: share the state of moduleRepo in a more elegant way
// I have made it global right now so Users and Modules share the same module repo instance
var ModuleRepo repositories.ModuleRepository

func RegisterModuleRoutes(firebaseAuth *auth.Client) {
	txManager := persistence.MockTransactionManager{}
	ModuleRepo = repositories.NewMockModuleRepository()
	moduleService := services.NewModuleService(&txManager, ModuleRepo)
	handler := handlers.NewModuleHandler(moduleService)

	http.HandleFunc("GET /api/modules", middleware.WithFirebaseAuth(firebaseAuth, handler.GetAllModules))
	http.HandleFunc("POST /api/modules", middleware.WithFirebaseAuth(firebaseAuth, handler.CreateModule))
}
