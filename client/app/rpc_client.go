package app

import (
	"context"
	"os"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/connectivity"

	"github.com/johnkhk/cli_chat_app/client/logger"
	"github.com/johnkhk/cli_chat_app/genproto/auth" // Import the generated protobuf package
)

var (
	conn       *grpc.ClientConn
	authClient auth.AuthServiceClient
)

// InitializeRPCClient sets up the gRPC client connection asynchronously.
func InitializeRPCClient() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		logger.Log.Fatalf("Error loading .env file: %v", err)
	}

	// Get the server address from environment variables
	serverAddress := os.Getenv("SERVER_ADDRESS")

	// Create a context for the connection attempt
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Set up a backoff strategy for reconnection attempts
	backoffConfig := backoff.Config{
		BaseDelay:  1.0 * time.Second,  // Initial delay
		Multiplier: 1.6,                // Multiplier for each attempt
		MaxDelay:   15.0 * time.Second, // Maximum delay
	}

	// Initialize the gRPC client connection with backoff
	var err error
	conn, err = grpc.DialContext(ctx, serverAddress, grpc.WithInsecure(), grpc.WithConnectParams(grpc.ConnectParams{Backoff: backoffConfig}))
	if err != nil {
		logger.Log.Fatalf("Failed to connect to server: %v", err)
	}

	// Check the connection state periodically
	go monitorConnectionState()

	authClient = auth.NewAuthServiceClient(conn)
	logger.Log.Infof("Attempting to connect to gRPC server at %s", serverAddress)
}

func monitorConnectionState() {
	for {
		state := conn.GetState()
		logger.Log.Infof("gRPC connection state: %s", state.String())

		if state == connectivity.Ready {
			logger.Log.Info("Connected to gRPC server")
			return
		}

		// Wait before checking again
		time.Sleep(2 * time.Second)
	}
}

// GetAuthClient returns the initialized AuthServiceClient.
func GetAuthClient() auth.AuthServiceClient {
	return authClient
}

// RegisterUser sends a register request to the server.
func RegisterUser() {
	// Create a new request with user details.
	req := &auth.RegisterRequest{
		Username: "exampleUser",
		Password: "examplePassword", // Hash passwords in a real application
	}

	// Send the request to the server.
	resp, err := authClient.RegisterUser(context.Background(), req)
	if err != nil {
		logger.Log.Errorf("Failed to register user: %v", err)
		return
	}

	if resp.Success {
		logger.Log.Info("Registration successful:", resp.Message)
	} else {
		logger.Log.Info("Registration failed:", resp.Message)
	}
}

// LoginUser sends a login request to the server.
func LoginUser() {
	// Create a new request with user details.
	req := &auth.LoginRequest{
		Username: "exampleUser",
		Password: "examplePassword",
	}

	// Send the request to the server.
	resp, err := authClient.LoginUser(context.Background(), req)
	if err != nil {
		logger.Log.Errorf("Failed to login: %v", err)
		return
	}

	if resp.Success {
		logger.Log.Info("Login successful:", resp.Message)
	} else {
		logger.Log.Info("Login failed:", resp.Message)
	}
}
