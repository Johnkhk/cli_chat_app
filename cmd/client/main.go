package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"

	"github.com/johnkhk/cli_chat_app/client/app"
	"github.com/johnkhk/cli_chat_app/client/logger"
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

	// Initialize the gRPC client
	app.InitializeRPCClient()

	// // Initialize token storage
	// tokenStorage := initializeTokenStorage()

	// // Create a new TokenManager instance with the storage
	// tokenManager := storage.NewTokenManager(tokenStorage)

	// Start client loop
	runClientLoop()
}

func runClientLoop() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("\n> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		command := strings.TrimSpace(input)
		handleCommand(command)
	}
}

func handleCommand(command string) {
	switch command {
	case "/register":
		app.RegisterUser()
	case "/login":
		app.LoginUser()
	case "/help":
		app.DisplayBanner()
	case "/quit":
		fmt.Println("Exiting the application.")
		os.Exit(0)
	default:
		fmt.Println("Unknown command:", command)
		fmt.Println("Type '/help' to see a list of available commands.")
	}
}
