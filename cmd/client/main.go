package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"

	"github.com/johnkhk/cli_chat_app/client/app"
	"github.com/johnkhk/cli_chat_app/client/logger"
	"github.com/johnkhk/cli_chat_app/client/ui"
)

func main() {

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
		os.Exit(1)
	}

	// Initialize the client logger
	log := logger.InitLogger()
	log.Info("Client application started")

	// Initialize the TokenManager with file path for token storage
	filePath := filepath.Join(os.Getenv("HOME"), ".cli_chat_app", "jwt_tokens") // For Linux/macOS
	// filePath := filepath.Join(os.Getenv("USERPROFILE"), ".cli_chat_app", "jwt_tokens") // For Windows
	tokenManager := app.NewTokenManager(filePath, nil) // Will set the client later

	// Initialize the gRPC client using RpcClient
	rpcClient, err := app.NewRpcClient("localhost:50051", log, tokenManager)
	if err != nil {
		log.Fatalf("Failed to initialize RPC clients: %v", err)
	}
	defer rpcClient.CloseConnections() // Ensure the connection is closed when the application exits
	log.Info("gRPC clients initialized.")

	// Set the AuthClient in the TokenManager
	tokenManager.SetClient(rpcClient.AuthClient.Client)

	// Attempt auto login
	err = tokenManager.TryAutoLogin()
	if err != nil {
		log.Infof("Log in failed: %v", err)
	} else {
		log.Info("User automatically logged in with stored tokens.")
	}
	isLoggedIn := err == nil

	// Use the UI manager to run the appropriate UI based on auth status
	ui.RunUIBasedOnAuthStatus(isLoggedIn, log, rpcClient)

}
