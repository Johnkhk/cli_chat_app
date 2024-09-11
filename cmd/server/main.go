package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/johnkhk/cli_chat_app/server/app"
	"github.com/johnkhk/cli_chat_app/server/logger"
	"github.com/johnkhk/cli_chat_app/server/storage"
)

func main() {
	envPath := os.Getenv("ENV_PATH")

	if envPath == "" {
		// Default to loading the .env file if ENV_PATH is not set
		envPath = ".env"
	}

	// Load environment variables from the specified .env file path
	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("Error loading environment file %s: %v", envPath, err)
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

	// Create a context that is canceled on SIGINT or SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Run the gRPC server in a separate goroutine
	serverErr := make(chan error, 1)
	go func() {
		if err := app.RunGRPCServer(ctx, port, db, log); err != nil {
			serverErr <- err
		}
	}()

	// Wait for the context to be canceled (signal received) or the server to return an error
	select {
	case <-ctx.Done():
		log.Info("Received shutdown signal, shutting down gracefully...")
	case err := <-serverErr:
		log.Fatalf("gRPC server error: %v", err)
	}

	// Allow some time for shutdown tasks
	time.Sleep(1 * time.Second)
	log.Info("Server stopped.")
}
