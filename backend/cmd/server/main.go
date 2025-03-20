package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"gdp8-backend/internal/firebase"
	"gdp8-backend/internal/handlers"
	"gdp8-backend/internal/middleware"
	"gdp8-backend/internal/persistence"
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

	routes.RegisterAllRoutes(firebaseAuth, txManager)

	ctx := context.Background()
	calendarService, err := services.NewGoogleCalendarService(ctx, credentialsPath, "user@example.com")
	if err != nil {
		log.Fatalf("Failed to initialize Google Calendar service: %v", err)
	}

	calendarHandler := handlers.NewCalendarHandler(calendarService)
	routes.RegisterCalendarRoutes(calendarHandler)

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
