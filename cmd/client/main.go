package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

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
	appDirPath, err := app.GetAppDirPath()
	if err != nil {
		log.Fatalf("Failed to get app directory path: %v", err)
	}

	// Initialize the gRPC client using RpcClient
	rpcClientConfig := app.RpcClientConfig{
		ServerAddress: os.Getenv("SERVER_ADDRESS"),
		Logger:        log,
		AppDirPath:    appDirPath,
	}
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

	// Create a context that is canceled on interrupt signals
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up a channel to listen for OS signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Start a goroutine to listen for signals
	go func() {
		sig := <-sigs
		log.Infof("Received signal: %v, initiating shutdown...", sig)
		cancel()
	}()

	// Run the UI in a separate goroutine
	uiDone := make(chan struct{})
	go func() {
		ui.RunUIBasedOnAuthStatus(isLoggedIn, log, rpcClient)
		close(uiDone)
	}()

	// Wait for either the UI to finish or a signal to be received
	select {
	case <-ctx.Done():
		log.Info("Context canceled, shutting down...")
	case <-uiDone:
		log.Info("UI exited, shutting down...")
	}

	// Deferred CloseConnections will be called here
	log.Info("Application shutdown complete.")
}
