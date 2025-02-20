package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"gdp8-backend/internal/firebase"
	"gdp8-backend/internal/middleware"
	"gdp8-backend/internal/routes"
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

	routes.RegisterAllRoutes(firebaseAuth)

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
