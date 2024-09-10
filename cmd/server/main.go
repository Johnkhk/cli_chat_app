package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/johnkhk/cli_chat_app/server/app"
	"github.com/johnkhk/cli_chat_app/server/logger"
	"github.com/johnkhk/cli_chat_app/server/storage"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Check the port value
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("PORT environment variable is not set.")
	}

	// Initialize the logger
	log := logger.InitLogger()
	log.Info("Server application started")

	// Connect to the database
	db, err := storage.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize the database: %v", err)
	}
	defer db.Close()

	// Create and run the gRPC server
	if err := app.RunGRPCServer(port, db, log); err != nil {
		log.Fatalf("Failed to run gRPC server: %v", err)
	}
}
