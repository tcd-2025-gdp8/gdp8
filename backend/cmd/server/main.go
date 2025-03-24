package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"gdp8-backend/internal/firebase"
	"gdp8-backend/internal/middleware"
	"gdp8-backend/internal/persistence"
	"gdp8-backend/internal/routes"
)

func main() {
	credentialsPath := "credentials/serviceAccountKey.json"
	if os.Getenv("FIREBASE_AUTH_EMULATOR_HOST") != "" {
		credentialsPath = "credentials/mockServiceAccountKey.json"
	}

	// Initialize Firebase
	firebaseAuth, err := firebase.InitializeFirebase(credentialsPath)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase Admin SDK: %v", err)
	}

	// Open DB connection
	db, err := persistence.OpenDB()
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
	}(db)

	if err := persistence.ExecuteMigrations(db); err != nil {
		log.Fatalf("Failed to execute migrations: %v", err)
	}

	txManager := persistence.NewSQLTransactionManager(db)

	// Register all routes including calendar
	routes.RegisterAllRoutes(firebaseAuth, txManager)

	// Set up HTTP server with CORS middleware
	corsHandler := middleware.SimpleCORS(http.DefaultServeMux)

	server := http.Server{
		Addr:              ":8080",
		Handler:           corsHandler,
		ReadHeaderTimeout: 3 * time.Second,
	}

	log.Println("Server running on http://localhost" + server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
