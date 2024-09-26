package main

import (
	"fmt"
	"os"

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
	// // Establish a single gRPC connection to the server

	log.Infof("APPDIRPATH: %s", os.Getenv("APP_DIR_PATH"))
	rpcClientConfig := app.RpcClientConfig{
		ServerAddress: "localhost:50051", // Replace with your server address
		Logger:        log,
		AppDirPath:    os.Getenv("APP_DIR_PATH"),
	}
	// Initialize the gRPC client using RpcClient
	rpcClient, err := app.NewRpcClient(rpcClientConfig)
	if err != nil {
		log.Fatalf("Failed to initialize RPC clients: %v", err)
	}
	defer rpcClient.CloseConnections() // Ensure the connection is closed when the application exits
	log.Info("gRPC clients initialized.")

	// Attempt auto login
	err = rpcClient.AuthClient.TokenManager.TryAutoLogin()
	if err != nil {
		log.Infof("Log in failed: %v", err)
	} else {
		log.Info("User automatically logged in with stored tokens.")
	}
	isLoggedIn := err == nil

	// Use the UI manager to run the appropriate UI based on auth status
	ui.RunUIBasedOnAuthStatus(isLoggedIn, log, rpcClient)

}
