package app

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/johnkhk/cli_chat_app/genproto/auth"
)

// AuthClient encapsulates the gRPC client and logger for authentication services.
type AuthClient struct {
	Client       auth.AuthServiceClient
	Logger       *logrus.Logger
	TokenManager *TokenManager
}

// RegisterUser sends a registration request to the server.
func (c *AuthClient) RegisterUser(username, password string) error {
	req := &auth.RegisterRequest{
		Username: username,
		Password: password,
	}

	resp, err := c.Client.RegisterUser(context.Background(), req)
	if err != nil {
		c.Logger.Errorf("Failed to register user: %v", err)
		return fmt.Errorf("Failed to register user: %v", err)
	}

	if resp.Success {
		c.Logger.Infof("Registration successful: %s", resp.Message)
		return nil
	} else {
		c.Logger.Infof("Registration failed: %s", resp.Message)
		return fmt.Errorf("Registration failed: %s", resp.Message)
	}
}

// LoginUser sends a login request to the server.
func (c *AuthClient) LoginUser(username, password string) error {
	req := &auth.LoginRequest{
		Username: username,
		Password: password,
	}

	resp, err := c.Client.LoginUser(context.Background(), req)
	if err != nil {
		c.Logger.Warnf("Failed to login: %v", err)
		return fmt.Errorf("Failed to login: %v", err)
	}

	if resp.Success {
		c.Logger.Infof("Login successful: %s", resp.Message)

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
