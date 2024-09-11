package rpc

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/johnkhk/cli_chat_app/client/app"
	"github.com/johnkhk/cli_chat_app/test/setup"
)

// Test the registration and login flow
func TestRegisterLoginFlow(t *testing.T) {
	t.Parallel() // Allow this test to run in parallel

	// Initialize resources using default configuration
	authClient, _, cleanup := setup.InitializeTestResources(t, nil)
	defer cleanup()
	log := authClient.Logger

	// Get the directory path for storing JWT tokens
	tokenDir := filepath.Join(os.TempDir(), fmt.Sprintf(".test_cli_chat_app_%s", t.Name()))
	tokenFile := filepath.Join(tokenDir, "jwt_tokens") // Assuming the token file is named "jwt_tokens"

	// Test the login of an unregistered user
	log.Infof("Testing unregistered user login")
	err := authClient.LoginUser("unregistered", "testpassword")
	if err == nil {
		t.Fatalf("Login should fail for unregistered user")
	}

	// Ensure the JWT token file was not created after failed login
	log.Infof("Checking JWT token file after login attempt for unregistered user")
	if _, err := os.Stat(tokenFile); err == nil {
		t.Fatalf("JWT token file was created after failed login")
	}

	// Register the user
	log.Infof("Registering user")
	err = authClient.RegisterUser("unregistered", "testpassword")
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	// Ensure the JWT token file was not created after registration
	log.Infof("Checking JWT token file after registration")
	if _, err := os.Stat(tokenFile); err == nil {
		t.Fatalf("JWT token file was created after registration but before login")
	}

	// Test the login of the registered user
	log.Infof("Testing registered user login")
	err = authClient.LoginUser("unregistered", "testpassword")
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	// Ensure the JWT token file was created after successful login
	log.Infof("Checking JWT token file after successful login")
	if _, err := os.Stat(tokenFile); err != nil {
		t.Fatalf("JWT token file was not created after successful login")
	}
}

// Test token expiration and refresh
// If you have no tokens/both expired, the only way to get tokens is to login
// If you have an expired access token, you can refresh it on any request
// Only way to get a new refresh token is to login
func TestTokenExpirationAndRefresh(t *testing.T) {
	t.Parallel() // Allow this test to run in parallel

	// Create a custom server configuration with a shorter access token duration for testing
	customConfig := setup.NewDefaultTestServerConfig()
	customConfig.AccessTokenDuration = time.Second * 3  // Set a very short access token duration (5 seconds)
	customConfig.RefreshTokenDuration = time.Minute * 6 // Set a short refresh token duration (1 minute)

	// Initialize resources with the custom server configuration
	authClient, _, cleanup := setup.InitializeTestResources(t, customConfig)
	defer cleanup()
	log := authClient.Logger

	// Register and login the user
	log.Infof("Registering user for expiration test")
	err := authClient.RegisterUser("expiringuser", "testpassword")
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	log.Infof("Logging in user for expiration test")
	err = authClient.LoginUser("expiringuser", "testpassword")
	if err != nil {
		t.Fatalf("Failed to login user: %v", err)
	}

	// Sleep to allow the access token to expire
	log.Infof("Waiting for the access token to expire")
	time.Sleep(time.Second * 5) // Wait for 6 seconds to ensure the access token has expired

	// Make sure the access token is expired but not the refresh token
	accessToken, refreshToken, err := authClient.TokenManager.ReadTokens()
	if err != nil {
		t.Fatalf("Failed to read tokens: %v", err)
	}

	// Check if the access token is expired
	if expired, err := app.IsTokenExpired(accessToken); err != nil || !expired {
		t.Fatalf("Access token should be expired, got err: %v", err)
	}
	// Check if the refresh token is still valid
	if expired, err := app.IsTokenExpired(refreshToken); err != nil || expired {
		t.Fatalf("Refresh token should not be expired, got err: %v", err)
	}

	// Attempt to get a new access token
	log.Infof("Attempting to get a new access token")
	err = authClient.LoginUser("expiringuser", "testpassword")
	if err != nil {
		t.Fatalf("Failed to login with expired access token but valid refresh token: %v", err)
	}

	// Check that both tokens are now valid
	newAccessToken, newRefreshToken, err := authClient.TokenManager.ReadTokens()

	// Check that new tokens have replaced the old ones
	if accessToken == newAccessToken || refreshToken == newRefreshToken {
		t.Fatalf("New tokens should have replaced the old ones")
	}

	// Check if the access token is expired
	if expired, err := app.IsTokenExpired(newAccessToken); err != nil || expired {
		t.Fatalf("Access token should be not be expired since it was refreshed, got err: %v", err)
	}
	// Check if the refresh token is still valid
	if expired, err := app.IsTokenExpired(newRefreshToken); err != nil || expired {
		t.Fatalf("Refresh token should not be expired since it was not refreshed, got err: %v", err)
	}

	// Sleep to allow the refresh token to expire
	log.Infof("Waiting for the refresh token to expire")
	time.Sleep(time.Second * 7) // Wait for 4 seconds to ensure the refresh token has expired

	// Attempt to hit a protected endpoint with an expired refresh token. (should fail)
	// TODO: Implement a protected endpoint to test this

	// Attempt to login with an expired refresh token. Should refresh both tokens if successful
	log.Infof("Attempting to login with expired refresh token")
	err = authClient.LoginUser("expiringuser", "testpassword")
	if err != nil {
		t.Fatalf("Failed to login with expired refresh token: %v", err)
	}

	// Check that both tokens are now valid
	newerAccessToken, newerRefreshToken, err := authClient.TokenManager.ReadTokens()

	// Check that new tokens have replaced the old ones
	if newAccessToken == newerAccessToken || newRefreshToken == newerRefreshToken {
		t.Fatalf("New tokens should have replaced the old ones")
	}

	// Check if the access token is expired
	if expired, err := app.IsTokenExpired(newerAccessToken); err != nil || expired {
		t.Fatalf("Access token should be not be expired since it was refreshed, got err: %v", err)
	}
	// Check if the refresh token is still valid
	if expired, err := app.IsTokenExpired(newerRefreshToken); err != nil || expired {
		t.Fatalf("Refresh token should not be expired since it was refreshed, got err: %v", err)
	}

}
