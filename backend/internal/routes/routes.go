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

	calendarService, err := services.NewCalendarService("credentials/serviceAccountKey.json")
	if err != nil {
		log.Fatalf("Failed to initialize calendar service: %v", err)
	}
	// Test Google Calendar Invite at startup
	go func() {
		testInput := services.CalendarInviteInput{
			Summary:     "Startup Test Event",
			Description: "This event was sent from routes.go during backend boot",
			Location:    "Zoom",
			StartTime:   "2025-04-10T12:00:00Z",
			EndTime:     "2025-04-10T13:00:00Z",
			Attendees:   []string{"anandsainbileg@gmail.com"},
			Organizer:   "anandsainbileg@gmail.com",
		}

		log.Println("📨 Sending test calendar invite...")
		err := calendarService.SendInvite(testInput)
		if err != nil {
			log.Printf("❌ Calendar invite test failed: %v", err)
		} else {
			log.Println("✅ Test invite sent successfully")
		}
	}()

	RegisterCalendarRoutes(firebaseAuth, calendarService) // ✅ This must be called!

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

}
