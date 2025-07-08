package main

import (
	"log"
	"net/http"

	"github.com/SHIVAM-GOUR/gbt-master-backend/config"
	"github.com/SHIVAM-GOUR/gbt-master-backend/routes"
)

func main() {
	// Initialize database
	config.InitDB()

	// Setup routes
	r := routes.SetupRoutes()

	// Start server
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
