package routes

import (
	"net/http"

	"firebase.google.com/go/v4/auth"

	"gdp8-backend/internal/handlers"
)

func RegisterAllRoutes(firebaseAuth *auth.Client) {
	RegisterStudyGroupRoutes(firebaseAuth)

	authHandler := handlers.NewAuthHandler(firebaseAuth)
	http.HandleFunc("/api/auth/verify", authHandler.VerifyHandler)
}
