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
	notificationRepo := repositories.SQLNotificationRepository{}
	studySessionRepo := repositories.SQLStudySessionRepository{}
	studySessionAvailabilityRepo := repositories.SQLStudySessionAvailabilityRepository{}

	userService := services.NewUserService(txManager, &userRepo)
	moduleService := services.NewModuleService(txManager, &moduleRepo)
	notificationService := services.NewNotificationService(txManager, &notificationRepo)
	studyGroupService := services.NewStudyGroupService(txManager, &studyGroupRepo, notificationService)
	studySessionService := services.NewStudySessionService(txManager, &studySessionAvailabilityRepo, &studySessionRepo)

	RegisterStudyGroupRoutes(firebaseAuth, studyGroupService, studySessionService)
	RegisterModuleRoutes(firebaseAuth, moduleService)
	RegisterNotificationRoutes(firebaseAuth, notificationService)
	RegisterUserRoutes(firebaseAuth, userService)
	RegisterChatRoutes(firebaseAuth)
	RegisterFileRoutes(firebaseAuth, studyGroupService)

	authHandler := handlers.NewAuthHandler(firebaseAuth)
	http.HandleFunc("/api/auth/verify", authHandler.VerifyHandler)
}
