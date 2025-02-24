package routes

import (
	"net/http"

	"firebase.google.com/go/v4/auth"

	"gdp8-backend/internal/handlers"
	"gdp8-backend/internal/persistence"
	"gdp8-backend/internal/repositories"
	"gdp8-backend/internal/services"
)

func RegisterAllRoutes(firebaseAuth *auth.Client, txManager persistence.TransactionManager) {
	studyGroupRepo := repositories.SQLStudyGroupRepository{}
	userRepo := repositories.SQLUserRepository{}
	moduleRepo := repositories.SQLModuleRepository{}

	studyGroupService := services.NewStudyGroupService(txManager, &studyGroupRepo)
	userService := services.NewUserService(txManager, &userRepo)
	moduleService := services.NewModuleService(txManager, &moduleRepo)

	RegisterStudyGroupRoutes(firebaseAuth, studyGroupService)
	RegisterModuleRoutes(firebaseAuth, moduleService)
	RegisterUserRoutes(firebaseAuth, userService)

	authHandler := handlers.NewAuthHandler(firebaseAuth)
	http.HandleFunc("/api/auth/verify", authHandler.VerifyHandler)
}
