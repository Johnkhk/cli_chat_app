package app

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/johnkhk/cli_chat_app/client/logger"
	"github.com/johnkhk/cli_chat_app/genproto/auth"
)

// InitializeRPCClient initializes the gRPC client and establishes a connection to the server.
func InitializeRPCClient() (auth.AuthServiceClient, *grpc.ClientConn, error) {
	// Establish a gRPC connection to the server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		logger.Log.Errorf("Failed to connect to server: %v", err)
		return nil, nil, err
	}

	// Create a new AuthService client using the established connection
	client := auth.NewAuthServiceClient(conn)

	return client, conn, nil
}

// RegisterUser sends a registration request to the server.
func RegisterUser(client auth.AuthServiceClient) {
	// Prompt the user for credentials
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	// Create a new request with user details.
	req := &auth.RegisterRequest{
		Username: username,
		Password: password, // In a real application, passwords should be hashed before sending.
	}

	// Send the request to the server.
	resp, err := client.RegisterUser(context.Background(), req)
	if err != nil {
		logger.Log.Errorf("Failed to register user: %v", err)
		return
	}

	// Check if registration was successful
	if resp.Success {
		logger.Log.Infof("Registration successful: %s", resp.Message)
	} else {
		logger.Log.Infof("Registration failed: %s", resp.Message)
	}
}

// LoginUser sends a login request to the server.
func LoginUser(client auth.AuthServiceClient, tokenManager *TokenManager) {
	// Prompt the user for credentials
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	// Create a new request with user details.
	req := &auth.LoginRequest{
		Username: username,
		Password: password,
	}

	// Send the request to the server.
	resp, err := client.LoginUser(context.Background(), req)
	if err != nil {
		logger.Log.Errorf("Failed to login: %v", err)
		return
	}

	// Check if login was successful
	if resp.Success {
		logger.Log.Infof("Login successful: %s", resp.Message)

		// Store the tokens in the TokenManager
		if err := tokenManager.StoreTokens(resp.AccessToken, resp.RefreshToken); err != nil {
			logger.Log.Errorf("Failed to store tokens: %v", err)
			return
		}

		logger.Log.Info("Tokens stored successfully.")
	} else {
		logger.Log.Infof("Login failed: %s", resp.Message)
	}
}
