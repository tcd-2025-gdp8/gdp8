package main

import (
	"gdp8-backend/internal/routes"
	"log"
	"net/http"
)

func main() {
	routes.RegisterAllRoutes()

	port := ":8080"
	log.Println("Server running on http://localhost" + port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
