package app

import (
	"context"
	"fmt"

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

func (c *AuthClient) RegisterUser(username, password string) error {

	// Create a new request with user details.
	req := &auth.RegisterRequest{
		Username: username,
		Password: password, // In a real application, passwords should be hashed before sending.
	}

	// Send the request to the server.
	resp, err := c.Client.RegisterUser(context.Background(), req)
	if err != nil {
		c.Logger.Errorf("Failed to register user: %v", err)
		return fmt.Errorf("Failed to register user: %v", err)
	}

	// Check if registration was successful
	if resp.Success {
		c.Logger.Infof("Registration successful: %s", resp.Message)
		return nil
	} else {
		c.Logger.Infof("Registration failed: %s", resp.Message)
		return nil
	}
}

// LoginUser sends a login request to the server.
func (c *AuthClient) LoginUser(username, password string) error {
	// Create a new request with user details.
	req := &auth.LoginRequest{
		Username: username,
		Password: password,
	}

	// Send the request to the server.
	resp, err := c.Client.LoginUser(context.Background(), req)
	if err != nil {
		c.Logger.Errorf("Failed to login: %v", err)
		return fmt.Errorf("Failed to login: %v", err)
	}

	// Check if login was successful
	if resp.Success {
		c.Logger.Infof("Login successful: %s", resp.Message)

		// Store the tokens in the TokenManager
		if err := c.TokenManager.StoreTokens(resp.AccessToken, resp.RefreshToken); err != nil {
			c.Logger.Errorf("Failed to store tokens: %v", err)
			return fmt.Errorf("Failed to store tokens: %v", err)
		}

		c.Logger.Info("Tokens stored successfully.")
		return nil
	} else {
		c.Logger.Infof("Login failed: %s", resp.Message)
		return fmt.Errorf("Login failed: %s", resp.Message)
	}
}

// CloseConnection closes the gRPC connection.
func (c *AuthClient) CloseConnection() {
	if err := c.Connection.Close(); err != nil {
		c.Logger.Errorf("Failed to close the connection: %v", err)
	}
}
