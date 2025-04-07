package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"gdp8-backend/internal/firebase"
	"gdp8-backend/internal/integrations"
	"gdp8-backend/internal/middleware"
	"gdp8-backend/internal/persistence"
	"gdp8-backend/internal/repositories"
	"gdp8-backend/internal/routes"
	"gdp8-backend/internal/services"
)

func main() {
	credentialsPath := "credentials/serviceAccountKey.json"
	if os.Getenv("FIREBASE_AUTH_EMULATOR_HOST") != "" {
		credentialsPath = "credentials/mockServiceAccountKey.json"
	}

	firebaseAuth, err := firebase.InitializeFirebase(credentialsPath)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase Admin SDK: %v", err)
	}

	db, err := persistence.OpenDB()
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	err = persistence.ExecuteMigrations(db)
	if err != nil {
		log.Fatalf("Failed to execute migrations: %v", err)
	}

	txManager := persistence.NewSQLTransactionManager(db)

	userRepo := repositories.SQLUserRepository{}
	userService := services.NewUserService(txManager, &userRepo)

	eventInvitesService := initializeEventInvitesService(userService)

	routes.RegisterAllRoutes(firebaseAuth, txManager, userService, eventInvitesService)

	corsHandler := middleware.SimpleCORS(http.DefaultServeMux)

	server := http.Server{
		Addr:              ":8080",
		Handler:           corsHandler,
		ReadHeaderTimeout: 3 * time.Second,
	}

	log.Println("Server running on http://localhost" + server.Addr)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

func initializeEventInvitesService(userService services.UserService) services.EventInvitesService {
	googleCalendarService, err := integrations.GetCalendarService()
	if err != nil {
		log.Printf("Failed to initialize Google Calendar service: %v. "+
			"Google Calendar features will not be available.", err)
		return &services.NoOpEventInvitesService{}
	}

	return services.NewEventInvitesService(googleCalendarService, userService)
}
