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

	// Initialize the gRPC client and create an AuthClient instance. (lazy dial)
	client, err := app.NewAuthClient("localhost:50051", log, tokenManager)
	if err != nil {
		log.Fatalf("Failed to initialize gRPC client: %v", err)
	}
	defer client.CloseConnection() // Ensure the connection is closed when the application exits
	log.Info("gRPC client initialized.")

	// Set the AuthClient in the TokenManager
	tokenManager.SetClient(client.Client)

	err = tokenManager.TryAutoLogin()
	if err != nil {
		log.Infof("Log in failed: %v", err)
	} else {
		log.Info("User automatically logged in with stored tokens.")
	}
	isLoggedIn := err == nil

	// Use the UI manager to run the appropriate UI based on auth status
	ui.RunUIBasedOnAuthStatus(isLoggedIn, log, client)

	// Start client loop
	// runClientLoop(client, log)
}

// // Main client loop to handle user input
// func runClientLoop(client *app.AuthClient, log *logrus.Logger) {
// 	reader := bufio.NewReader(os.Stdin)

// 	for {
// 		fmt.Print("\n> ")
// 		input, err := reader.ReadString('\n')
// 		if err != nil {
// 			log.Error("Error reading input:", err)
// 			continue
// 		}

// 		command := strings.TrimSpace(input)
// 		handleCommand(command, client, log)
// 	}
// }

// // Handle user commands
// func handleCommand(command string, client *app.AuthClient, log *logrus.Logger) {
// 	switch command {
// 	// Unauthenticated commands
// 	case "/register":
// 		client.RegisterUser()
// 	case "/login":
// 		client.LoginUser()
// 	case "/help":
// 		app.DisplayBanner() // Display help or banner
// 	case "/quit":
// 		fmt.Println("Exiting the application.")
// 		os.Exit(0)

// 	// Authenticated commands
// 	// case "/logout":
// 	// case "/add_friend":
// 	// client.AddFriend() // Example usage with AuthClient

// 	/*
// 		TODO:
// 		1. With the new auth, try to streamline the client CLI's entry point (main.go)
// 		2. Implement /logout and /add_friend commands
// 	*/

// 	default:
// 		log.Warn("Unknown command: ", command)
// 		fmt.Println("Type '/help' to see a list of available commands.")
// 	}
// }
