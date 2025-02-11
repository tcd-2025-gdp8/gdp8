package main

import (
	"log"
	"net/http"
	"time"

	"gdp8-backend/internal/routes"
)

func main() {
	routes.RegisterAllRoutes()

	server := http.Server{
		Addr:              ":8080",
		ReadHeaderTimeout: 3 * time.Second,
	}

	log.Println("Server running on http://localhost" + server.Addr)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
