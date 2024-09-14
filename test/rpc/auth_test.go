package rpc

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/johnkhk/cli_chat_app/test"
	"github.com/johnkhk/cli_chat_app/test/setup"
)

// Test the registration and login flow
func TestRegisterLoginFlow(t *testing.T) {
	t.Parallel() // Allow this test to run in parallel

	// Initialize resources using default configuration
	rpcClients, _, cleanup := setup.InitializeTestResources(t, nil, 1)
	rpcClient := rpcClients[0]
	defer cleanup()
	log := rpcClient.Logger

	// Get the directory path for storing JWT tokens
	tokenDir := filepath.Join(os.TempDir(), fmt.Sprintf(".test_cli_chat_app_%s_client_0", t.Name()))
	tokenFile := filepath.Join(tokenDir, "jwt_tokens") // Assuming the token file is named "jwt_tokens"

	// Test the login of an unregistered/wrong credentials user
	log.Infof("Testing unregistered user login")
	err := rpcClient.AuthClient.LoginUser("unregistered", "testpassword")
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
	err = rpcClient.AuthClient.RegisterUser("unregistered", "testpassword")
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
	err = rpcClient.AuthClient.LoginUser("unregistered", "testpassword")
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	// Ensure the JWT token file was created after successful login
	log.Infof("Checking JWT token file after successful login")
	if _, err := os.Stat(tokenFile); err != nil {
		t.Fatalf("JWT token file was not created after successful login")
	}
}

// / TestTokenExpirationAndRefresh tests the token expiration and refresh functionality.
// It verifies that tokens expire correctly and can be refreshed appropriately.
// If you have no tokens or both are expired, the only way to get tokens is to login.
// If you have an expired access token, you can refresh it on any request.
// The only way to get a new refresh token is to login.
func TestTokenExpirationAndRefresh(t *testing.T) {
	t.Parallel() // Allow this test to run in parallel

	// Create a mock time provider
	mockTime := &test.MockTimeProvider{CurrentTime: time.Now()}

	// Create a custom server configuration with your desired token durations
	customConfig := setup.NewDefaultTestServerConfig()
	customConfig.AccessTokenDuration = time.Minute * 15    // Access token valid for 15 minutes
	customConfig.RefreshTokenDuration = time.Hour * 24 * 7 // Refresh token valid for 7 days
	customConfig.TimeProvider = mockTime                   // Inject the mock time provider

	// Initialize resources with the custom server configuration
	rpcClients, _, cleanup := setup.InitializeTestResources(t, customConfig, 1)
	defer cleanup()
	rpcClient := rpcClients[0]
	log := rpcClient.Logger

	// Register and login the user
	log.Infof("Registering user for expiration test")
	err := rpcClient.AuthClient.RegisterUser("expiringuser", "testpassword")
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	log.Infof("Logging in user for expiration test")
	err = rpcClient.AuthClient.LoginUser("expiringuser", "testpassword")
	if err != nil {
		t.Fatalf("Failed to login user: %v", err)
	}

	// Advance time to allow the access token to expire
	log.Infof("Advancing time to expire the access token")
	mockTime.Advance(customConfig.AccessTokenDuration + time.Minute*30) // Advance time past the access token expiration

	// Make sure the access token is expired but not the refresh token
	accessToken, refreshToken, err := rpcClient.AuthClient.TokenManager.ReadTokens()
	if err != nil {
		t.Fatalf("Failed to read tokens: %v", err)
	}

	// Check if the access token is expired
	if expired, err := rpcClient.AuthClient.TokenManager.IsTokenExpired(accessToken); err != nil || !expired {
		t.Fatalf("Access token should be expired, got err: %v", err)
	}

	// Check if the refresh token is still valid
	if expired, err := rpcClient.AuthClient.TokenManager.IsTokenExpired(refreshToken); err != nil || expired {
		t.Fatalf("Refresh token should not be expired, got err: %v", err)
	}

	// Attempt to get a new access token
	log.Infof("Attempting to get a new access token")
	err = rpcClient.AuthClient.LoginUser("expiringuser", "testpassword")
	if err != nil {
		t.Fatalf("Failed to login with expired access token but valid refresh token: %v", err)
	}

	// Check that both tokens are now valid
	newAccessToken, newRefreshToken, err := rpcClient.AuthClient.TokenManager.ReadTokens()
	if err != nil {
		t.Fatalf("Failed to read new tokens: %v", err)
	}

	fmt.Printf("New Access Token: %s\nOld Access Token: %s\n\n", newAccessToken, accessToken)
	fmt.Printf("New Refresh Token: %s\nOld Refresh Token: %s\n\n", newRefreshToken, refreshToken)
	fmt.Printf("==? %v ==? %v\n", newAccessToken == accessToken, newRefreshToken == refreshToken)
	// Check that new tokens have replaced the old ones
	if accessToken == newAccessToken || refreshToken == newRefreshToken {
		t.Fatalf("New tokens should have replaced the old ones")
	}

	// Unadvance time back to the current time
	log.Infof("Unadvancing time to current time")
	mockTime.Advance(-customConfig.AccessTokenDuration - time.Minute*30) // Reset time to current
	// Check if the new access token is not expired
	if expired, err := rpcClient.AuthClient.TokenManager.IsTokenExpired(newAccessToken); err != nil || expired {
		t.Fatalf("Access token should not be expired since it was refreshed, got err: %v", err)
	}

	// Check if the new refresh token is still valid
	if expired, err := rpcClient.AuthClient.TokenManager.IsTokenExpired(newRefreshToken); err != nil || expired {
		t.Fatalf("Refresh token should not be expired since it was not refreshed, got err: %v", err)
	}

	// Advance time to allow the refresh token to expire
	log.Infof("Advancing time to expire the refresh token")
	mockTime.Advance(customConfig.RefreshTokenDuration + time.Minute*30) // Advance time past the refresh token expiration

	// Attempt to hit a protected endpoint with an expired refresh token (should fail)
	// TODO: Implement a protected endpoint to test this

	// Attempt to login with an expired refresh token. Should get new tokens if successful
	log.Infof("Attempting to login with expired refresh token")
	err = rpcClient.AuthClient.LoginUser("expiringuser", "testpassword")
	if err != nil {
		t.Fatalf("Failed to login with expired refresh token: %v", err)
	}

	// Check that both tokens are now valid
	newerAccessToken, newerRefreshToken, err := rpcClient.AuthClient.TokenManager.ReadTokens()
	if err != nil {
		t.Fatalf("Failed to read newer tokens: %v", err)
	}

	// Check that new tokens have replaced the old ones
	if newAccessToken == newerAccessToken || newRefreshToken == newerRefreshToken {
		t.Fatalf("New tokens should have replaced the old ones")
	}

	// Unadvance time back to the current time
	log.Infof("Unadvancing time to current time")
	mockTime.Advance(-customConfig.RefreshTokenDuration - time.Minute*30) // Reset time to current

	// Check if the newer access token is not expired
	if expired, err := rpcClient.AuthClient.TokenManager.IsTokenExpired(newerAccessToken); err != nil || expired {
		t.Fatalf("Access token should not be expired since it was refreshed, got err: %v", err)
	}

	// Check if the newer refresh token is not expired
	if expired, err := rpcClient.AuthClient.TokenManager.IsTokenExpired(newerRefreshToken); err != nil || expired {
		t.Fatalf("Refresh token should not be expired since it was refreshed, got err: %v", err)
	}
}

// TestRegisterUserWithExistingUsername tests the registration of a user with an existing username
func TestRegisterUserWithExistingUsername(t *testing.T) {
	t.Parallel() // Allow this test to run in parallel

	// Initialize resources using default configuration
	rpcClients, _, cleanup := setup.InitializeTestResources(t, nil, 1)
	rpcClient := rpcClients[0]
	defer cleanup()
	log := rpcClient.Logger

	// Register a new user
	log.Infof("Registering new user")
	err := rpcClient.AuthClient.RegisterUser("existinguser", "testpassword")
	if err != nil {
		t.Fatalf("Failed to register new user: %v", err)
	}

	// Attempt to register a user with the same username
	log.Infof("Attempting to register user with existing username")
	err = rpcClient.AuthClient.RegisterUser("existinguser", "newpassword")
	if err == nil {
		t.Fatalf("Registration should fail for an existing username")
	}

	// Ensure the error is related to the username already being taken
	expectedErrMsg := "Registration failed: Username already exists" // Replace with the actual error message returned by your RegisterUser method
	if err.Error() != expectedErrMsg {
		t.Fatalf("Expected error message: %q, got: %q", expectedErrMsg, err.Error())
	}
}
