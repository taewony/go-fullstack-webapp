package main

import (
	"log"
	"net/http"

	"github.com/taewony/go-fullstack-webapp/internal/models"
	"github.com/taewony/go-fullstack-webapp/internal/router"
)

func main() {
	// Initialize the in-memory SQLite3 database connection
	models.InitDB()
	// models.InsertInitialDB()

	// Initialize the router
	r := router.NewRouter()

	// Start the server
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
