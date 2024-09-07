package app

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/johnkhk/cli_chat_app/genproto/auth"
)

// AuthClient encapsulates the gRPC client and logger.
type AuthClient struct {
	Client       auth.AuthServiceClient
	Connection   *grpc.ClientConn
	Logger       *logrus.Logger
	TokenManager *TokenManager
}

// NewAuthClient initializes a new AuthClient with the necessary dependencies.
func NewAuthClient(serverAddress string, logger *logrus.Logger, tokenManager *TokenManager) (*AuthClient, error) {
	// Establish a gRPC connection to the server
	// conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Errorf("Failed to connect to server: %v", err)
		return nil, err
	}

	// Create a new AuthService client using the established connection
	client := auth.NewAuthServiceClient(conn)

	return &AuthClient{
		Client:       client,
		Connection:   conn,
		Logger:       logger,
		TokenManager: tokenManager,
	}, nil
}

// RegisterUser sends a registration request to the server.
func (c *AuthClient) RegisterUser() {
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
	resp, err := c.Client.RegisterUser(context.Background(), req)
	if err != nil {
		c.Logger.Errorf("Failed to register user: %v", err)
		return
	}

	// Check if registration was successful
	if resp.Success {
		c.Logger.Infof("Registration successful: %s", resp.Message)
	} else {
		c.Logger.Infof("Registration failed: %s", resp.Message)
	}
}

// LoginUser sends a login request to the server.
func (c *AuthClient) LoginUser() {
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
	resp, err := c.Client.LoginUser(context.Background(), req)
	if err != nil {
		c.Logger.Errorf("Failed to login: %v", err)
		return
	}

	// Check if login was successful
	if resp.Success {
		c.Logger.Infof("Login successful: %s", resp.Message)

		// Store the tokens in the TokenManager
		if err := c.TokenManager.StoreTokens(resp.AccessToken, resp.RefreshToken); err != nil {
			c.Logger.Errorf("Failed to store tokens: %v", err)
			return
		}

		c.Logger.Info("Tokens stored successfully.")
	} else {
		c.Logger.Infof("Login failed: %s", resp.Message)
	}
}

// CloseConnection closes the gRPC connection.
func (c *AuthClient) CloseConnection() {
	if err := c.Connection.Close(); err != nil {
		c.Logger.Errorf("Failed to close the connection: %v", err)
	}
}
