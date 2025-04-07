package routes

import (
	"log"
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
	fileRepo := repositories.SQLFileRepository{}
	chatRepo := repositories.SQLChatMessagesRepository{}
	chatbotRepo := repositories.SQLChatBotMemoryRepo{}

	userService := services.NewUserService(txManager, &userRepo)
	moduleService := services.NewModuleService(txManager, &moduleRepo)
	studyGroupService := services.NewStudyGroupService(txManager, &studyGroupRepo, nil)
	notificationService := services.NewNotificationService(txManager, &notificationRepo, userService, studyGroupService)
	studyGroupService.SetNotificationService(notificationService)
	studySessionService := services.NewStudySessionService(txManager,
		&studySessionAvailabilityRepo, &studySessionRepo, studyGroupService, notificationService)
	fileService := services.NewFileService(txManager, &fileRepo, studyGroupService)
	chatService := services.NewChatService(txManager, &chatRepo, studyGroupService)

	chatbotService := services.NewChatBotService(
		txManager,
		studyGroupService,
		fileService,
		chatService,
		&chatbotRepo,
	)

	RegisterStudyGroupRoutes(firebaseAuth, studyGroupService)
	RegisterStudySessionRoutes(firebaseAuth, studySessionService)
	RegisterModuleRoutes(firebaseAuth, moduleService)
	RegisterNotificationRoutes(firebaseAuth, notificationService)
	RegisterUserRoutes(firebaseAuth, userService)
	RegisterChatRoutes(firebaseAuth, chatService)
	RegisterFileRoutes(firebaseAuth, studyGroupService, fileService, chatbotService)
	RegisterChatBotRoutes(firebaseAuth, studyGroupService, fileService, chatService, chatbotService)

	authHandler := handlers.NewAuthHandler(firebaseAuth)
	http.HandleFunc("/api/auth/verify", authHandler.VerifyHandler)

	calendarService, err := services.NewCalendarService("credentials/serviceAccountKey.json")
	if err != nil {
		log.Printf(" Calendar service not initialized: %v", err)
	} else {
		RegisterCalendarRoutes(firebaseAuth, calendarService)
	}
}
