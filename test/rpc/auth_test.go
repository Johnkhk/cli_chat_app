package rpc

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/Johnkhk/libsignal-go/protocol/prekey"

	"github.com/johnkhk/cli_chat_app/client/e2ee/store"
	"github.com/johnkhk/cli_chat_app/test"
	"github.com/johnkhk/cli_chat_app/test/setup"
)

// Test the registration and login flow
func TestRegisterLoginFlow(t *testing.T) {
	// t.Parallel() // Allow this test to run in parallel

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
	// t.Parallel() // Allow this test to run in parallel

	// Create a mock time provider
	mockTime := &test.MockTimeProvider{CurrentTime: time.Now()}

	// Create a custom server configuration with your desired token durations
	customConfig, err := setup.NewDefaultTestServerConfig(t)
	if err != nil {
		t.Fatalf("Failed to create default test server config: %v", err)
	}
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
	err = rpcClient.AuthClient.RegisterUser("expiringuser", "testpassword")
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
	// t.Parallel() // Allow this test to run in parallel

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

func TestOnLoginUploadKeysAndLocalIdentity(t *testing.T) {
	// t.Parallel() // Allow this test to run in parallel

	// Initialize resources using default configuration
	rpcClients, db, cleanup := setup.InitializeTestResources(t, nil, 1)
	rpcClient := rpcClients[0]
	defer cleanup()
	log := rpcClient.Logger

	// Register a new user
	log.Infof("Registering new user")
	err := rpcClient.AuthClient.RegisterUser("newuser", "testpassword")
	if err != nil {
		t.Fatalf("Failed to register new user: %v", err)
	}

	// Login the user
	log.Infof("Logging in user")
	err = rpcClient.AuthClient.LoginUser("newuser", "testpassword")
	if err != nil {
		t.Fatalf("Failed to login user: %v", err)
	}

	claims, err := rpcClient.AuthClient.TokenManager.GetClaimsFromAccessToken()
	if err != nil {
		t.Fatalf("Failed to get claims from access token: %v", err)
	}
	// Assuming sub is a string and represents the user ID
	userIDStr, ok := claims["sub"].(string)
	if !ok {
		t.Fatalf("Failed to extract user ID from claims")
	}

	// Convert userIDStr to int32 if necessary
	userIDFromJWT, err := strconv.Atoi(userIDStr) // Atoi converts string to int
	if err != nil {
		t.Fatalf("Failed to convert user ID to integer: %v", err)
	}

	fmt.Printf("Extracted User ID: %d\n", userIDFromJWT)

	// Query the database for the registration ID associated with the user
	var registrationID uint32
	err = rpcClient.Store.DB.QueryRow("SELECT registration_id FROM local_identity WHERE user_id = ?", userIDFromJWT).Scan(&registrationID)
	if err != nil {
		t.Fatalf("Failed to query registration ID for user %d: %v", userIDFromJWT, err)
	}
	extractedUserID, deviceID := store.ExtractUserIDAndDeviceID(registrationID)
	log.Infof("Extracted User ID: %d, Device ID: %d from Registration ID", extractedUserID, deviceID)

	// Verify that the extracted userID matches the userID from the JWT
	if extractedUserID != uint32(userIDFromJWT) {
		t.Fatalf("Extracted userID (%d) does not match userID from JWT (%d)", extractedUserID, userIDFromJWT)
	}

	// // Check that the public key was uploaded and stored in the database
	_ = db
	// Verify that the prekey bundle was uploaded and stored in the database
	var dbUserID, dbRegistrationID, dbDeviceID uint32
	var dbIdentityKey, dbPreKey, dbSignedPreKey, dbSignedPreKeySignature []byte
	var dbPreKeyID, dbSignedPreKeyID uint32

	// Query the database for the prekey bundle
	err = db.QueryRow(`
        SELECT user_id, registration_id, device_id, identity_key, pre_key_id, pre_key, signed_pre_key_id, signed_pre_key, signed_pre_key_signature
        FROM prekey_bundle WHERE user_id = ?`, userIDFromJWT).Scan(
		&dbUserID, &dbRegistrationID, &dbDeviceID, &dbIdentityKey, &dbPreKeyID, &dbPreKey, &dbSignedPreKeyID, &dbSignedPreKey, &dbSignedPreKeySignature)
	if err != nil {
		t.Fatalf("Failed to query prekey bundle for user %d: %v", userIDFromJWT, err)
	}

	// Verify that the values stored match what was uploaded
	if dbUserID != uint32(userIDFromJWT) {
		t.Fatalf("Expected userID %d, but got %d", userIDFromJWT, dbUserID)
	}
	if dbRegistrationID != registrationID {
		t.Fatalf("Expected registrationID %d, but got %d", registrationID, dbRegistrationID)
	}
	if dbDeviceID != deviceID {
		t.Fatalf("Expected deviceID %d, but got %d", deviceID, dbDeviceID)
	}

	// Retrieve the PreKey using PreKeyStore
	preKey, ok, err := rpcClient.Store.PreKeyStore().Load(context.Background(), prekey.ID(dbPreKeyID))
	if err != nil {
		t.Fatalf("Failed to load PreKey from store: %v", err)
	}
	if !ok {
		t.Fatalf("PreKey with ID %d not found in store", dbPreKeyID)
	}
	preKeyPair, err := preKey.KeyPair()
	if err != nil {
		t.Fatalf("Failed to get key pair from PreKey: %v", err)
	}
	preKeyPublicBytes := preKeyPair.PublicKey().Bytes()

	// Compare PreKey
	if !bytes.Equal(dbPreKey, preKeyPublicBytes) {
		t.Fatalf("Expected PreKey %x, but got %x", preKeyPublicBytes, dbPreKey)
	}

	// Retrieve the Signed PreKey using SignedPreKeyStore
	signedPreKey, ok, err := rpcClient.Store.SignedPreKeyStore().Load(context.Background(), prekey.ID(dbSignedPreKeyID))
	if err != nil {
		t.Fatalf("Failed to load Signed PreKey from store: %v", err)
	}
	if !ok {
		t.Fatalf("Signed PreKey with ID %d not found in store", dbSignedPreKeyID)
	}
	signedPreKeyPair, err := signedPreKey.KeyPair()
	if err != nil {
		t.Fatalf("Failed to get key pair from Signed PreKey: %v", err)
	}
	signedPreKeyPublicBytes := signedPreKeyPair.PublicKey().Bytes()

	// Compare Signed PreKey
	if !bytes.Equal(dbSignedPreKey, signedPreKeyPublicBytes) {
		t.Fatalf("Expected Signed PreKey %x, but got %x", signedPreKeyPublicBytes, dbSignedPreKey)
	}

	// Compare Signed PreKey Signature
	// if !bytes.Equal(dbSignedPreKeySignature, signedPreKey.Signature()) {
	if !bytes.Equal(dbSignedPreKeySignature, signedPreKey.GetSigned().Signature) {
		t.Fatalf("Expected Signed PreKey Signature %x, but got %x", signedPreKey.GetSigned().Signature, dbSignedPreKeySignature)
	}

	// Log for debugging
	log.Infof("Verified prekey bundle for user %d", dbUserID)
	log.Infof("IdentityKey: %x, PreKeyID: %d, PreKey: %x, SignedPreKeyID: %d, SignedPreKey: %x", dbIdentityKey, dbPreKeyID, dbPreKey, dbSignedPreKeyID, dbSignedPreKey)

	// Optionally add more checks or validations here if necessary

}
