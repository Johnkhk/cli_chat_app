// client/app/storage/token_manager.go

package storage

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/johnkhk/cli_chat_app/client/app"
	"github.com/johnkhk/cli_chat_app/client/logger"
)

// TokenManager struct holds the reference to TokenStorage.
type TokenManager struct {
	storage TokenStorage
}

// NewTokenManager creates a new TokenManager with the provided TokenStorage.
func NewTokenManager(storage TokenStorage) *TokenManager {
	return &TokenManager{storage: storage}
}

// TryAutoLogin checks if valid tokens exist and logs the user in automatically if possible.
func (tm *TokenManager) TryAutoLogin() bool {
	// Read the current tokens
	accessToken, refreshToken, err := tm.storage.ReadTokens()
	if err != nil {
		logger.Log.Errorf("Failed to read tokens: %v", err)
		return false
	}

	// Validate the access token locally
	if isTokenExpired(accessToken) {
		logger.Log.Info("Access token expired, attempting to refresh.")
		// If access token is expired, attempt to refresh it
		newAccessToken, err := app.RefreshAccessToken(refreshToken)
		if err != nil {
			logger.Log.Errorf("Failed to refresh access token: %v", err)
			return false
		}

		// Store the new access token and reuse the same refresh token
		if err := tm.storage.StoreTokens(newAccessToken, refreshToken); err != nil {
			logger.Log.Errorf("Failed to store refreshed tokens: %v", err)
			return false
		}

		logger.Log.Info("Access token refreshed successfully, user logged in automatically.")
		return true
	}

	logger.Log.Info("Access token is valid, user logged in automatically.")
	return true
}

// Helper function to check if a token is expired by decoding the JWT payload.
func isTokenExpired(tokenString string) bool {
	// Split the JWT into its parts: header, payload, signature
	parts := strings.Split(tokenString, ".")
	if len(parts) < 2 {
		return true // Invalid token format
	}

	// Decode the payload part (Base64URL encoded)
	payload, err := decodeSegment(parts[1])
	if err != nil {
		return true
	}

	// Parse the JSON payload to extract the "exp" field
	var claims struct {
		Exp int64 `json:"exp"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return true
	}

	// Check if the token is expired
	return time.Unix(claims.Exp, 0).Before(time.Now())
}

// Helper function to decode a Base64URL-encoded segment.
func decodeSegment(seg string) ([]byte, error) {
	return json.Marshal(seg) // Decode base64URL without padding
}

// LoginUserWithStorage logs in the user and stores the tokens using the provided storage.
func (tm *TokenManager) LoginUserWithStorage() {
	// Attempt to log in and obtain tokens
	accessToken, refreshToken, err := app.LoginUser()
	if err != nil {
		logger.Log.Errorf("Login failed: %v", err)
		return
	}

	// Store tokens securely
	if err := tm.storage.StoreTokens(accessToken, refreshToken); err != nil {
		logger.Log.Errorf("Failed to store tokens: %v", err)
		return
	}

	logger.Log.Info("User logged in successfully, tokens stored.")
}

// RefreshToken refreshes the access token using the stored refresh token.
func (tm *TokenManager) RefreshToken() {
	// Read the current tokens
	accessToken, refreshToken, err := tm.storage.ReadTokens()
	if err != nil {
		logger.Log.Errorf("Failed to read tokens: %v", err)
		return
	}

	// Use the refresh token to get a new access token
	newAccessToken, err := app.RefreshAccessToken(refreshToken)
	if err != nil {
		logger.Log.Errorf("Failed to refresh access token: %v", err)
		return
	}

	// Store the new access token and reuse the same refresh token
	if err := tm.storage.StoreTokens(newAccessToken, refreshToken); err != nil {
		logger.Log.Errorf("Failed to store refreshed tokens: %v", err)
		return
	}

	logger.Log.Info("Access token refreshed successfully.")
}
