package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"

	"github.com/johnkhk/cli_chat_app/client/app"
	"github.com/johnkhk/cli_chat_app/client/logger"
	"github.com/johnkhk/cli_chat_app/genproto/auth"
)

func main() {
	// Display banner
	app.DisplayBanner()

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		logger.Log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize the client logger
	logger.InitLogger()
	logger.Log.Info("Client application started")

	// Initialize the gRPC client and connection
	client, conn, err := app.InitializeRPCClient()
	if err != nil {
		logger.Log.Fatalf("Failed to initialize gRPC client: %v", err)
	}
	defer conn.Close() // Ensure the connection is closed when the application exits

	logger.Log.Info("gRPC client initialized and connected to server.")

	// Initialize the TokenManager with file path for token storage
	filePath := filepath.Join(os.Getenv("HOME"), ".cli_chat_app", "jwt_tokens") // For Linux/macOS
	// filePath := filepath.Join(os.Getenv("USERPROFILE"), ".cli_chat_app", "jwt_tokens") // For Windows
	tokenManager := app.NewTokenManager(filePath, client)

	// Attempt to automatically log in using stored tokens
	if tokenManager.TryAutoLogin() {
		logger.Log.Info("User automatically logged in with stored tokens.")
	} else {
		logger.Log.Info("Automatic login failed or no valid token found.")
	}

	// Start client loop
	runClientLoop(client, tokenManager)
}

// Main client loop to handle user input
func runClientLoop(client auth.AuthServiceClient, tokenManager *app.TokenManager) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("\n> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		command := strings.TrimSpace(input)
		handleCommand(command, client, tokenManager)
	}
}

// Handle user commands
func handleCommand(command string, client auth.AuthServiceClient, tokenManager *app.TokenManager) {
	switch command {
	// unauthenticated commands
	case "/register":
		app.RegisterUser(client)
	case "/login":
		app.LoginUser(client, tokenManager) // Token manager used to store token upon successful login
	case "/help":
		app.DisplayBanner() // Display help or banner
	case "/quit":
		fmt.Println("Exiting the application.")
		os.Exit(0) // Exit the application

	// authed commands
	// case "/logout":
	// case "/add_friend":
	// app.AddFriend(client, tokenManager) // Token manager used to retrieve token from local

	/*
		TODO:
		1. with the new auth try to streamline the client cli's entry point (main.go)
		2. implement /logout and /add_friend commands
	*/

	default:
		fmt.Println("Unknown command:", command)
		fmt.Println("Type '/help' to see a list of available commands.")
	}
}
