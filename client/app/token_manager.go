// client/app/token_manager.go

package app

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/johnkhk/cli_chat_app/client/logger"
	"github.com/johnkhk/cli_chat_app/genproto/auth"
)

// TokenManager handles all operations related to managing tokens.
type TokenManager struct {
	filePath string
	client   auth.AuthServiceClient // Use the gRPC client passed during initialization
}

// NewTokenManager creates a new TokenManager with the specified file path and gRPC client.
func NewTokenManager(filePath string, client auth.AuthServiceClient) *TokenManager {
	return &TokenManager{filePath: filePath, client: client}
}

// TryAutoLogin attempts to automatically log in the user using stored tokens.
func (tm *TokenManager) TryAutoLogin() bool {
	// Read the current tokens
	accessToken, refreshToken, err := tm.ReadTokens()
	if err != nil {
		logger.Log.Errorf("Failed to read tokens: %v", err)
		return false
	}

	// Check if the access token is expired
	if isTokenExpired(accessToken) {
		logger.Log.Info("Access token expired, attempting to refresh.")

		// Attempt to refresh the access token
		newAccessToken, err := tm.RefreshAccessToken(refreshToken)
		if err != nil {
			return false
		}

		// Store the new access token and reuse the refresh token
		if err := tm.StoreTokens(newAccessToken, refreshToken); err != nil {
			logger.Log.Errorf("Failed to store refreshed tokens: %v", err)
			return false
		}

		logger.Log.Info("Access token refreshed successfully, user logged in automatically.")
		return true
	}

	logger.Log.Info("Access token is valid, user logged in automatically.")
	return true
}

// StoreTokens stores the access and refresh tokens in a local file.
func (tm *TokenManager) StoreTokens(accessToken, refreshToken string) error {
	data := fmt.Sprintf("access_token:%s\nrefresh_token:%s", accessToken, refreshToken)

	// Write the tokens to a file with secure permissions
	return ioutil.WriteFile(tm.filePath, []byte(data), 0600)
}

// ReadTokens retrieves the access and refresh tokens from the local file.
func (tm *TokenManager) ReadTokens() (string, string, error) {
	data, err := ioutil.ReadFile(tm.filePath)
	if err != nil {
		return "", "", err
	}

	var accessToken, refreshToken string
	_, err = fmt.Sscanf(string(data), "access_token:%s\nrefresh_token:%s", &accessToken, &refreshToken)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

// Helper function to check if a token is expired by decoding the JWT payload.
func isTokenExpired(tokenString string) bool {
	// Split the JWT into its parts: header, payload, signature
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		logger.Log.Errorf("Invalid token format: expected 3 parts but got %d", len(parts))
		return true // Invalid token format
	}

	// Decode the payload part (Base64URL encoded)
	payload, err := decodeSegment(parts[1])
	if err != nil {
		logger.Log.Errorf("Failed to decode token payload: %v", err)
		return true
	}

	// Define a struct with a custom Unmarshal function to handle both string and int
	var claims struct {
		Exp interface{} `json:"exp"` // Use interface{} to handle both int and string
	}

	// Parse the JSON payload to extract the "exp" field
	if err := json.Unmarshal(payload, &claims); err != nil {
		logger.Log.Errorf("Failed to unmarshal token claims: %v", err)
		return true
	}

	// Convert the "exp" field to int64, handling both string and int cases
	var exp int64
	switch v := claims.Exp.(type) {
	case float64: // JSON numbers are unmarshaled as float64
		exp = int64(v)
	case string:
		var err error
		exp, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			logger.Log.Errorf("Failed to parse expiration time from string: %v", err)
			return true
		}
	default:
		logger.Log.Errorf("Unexpected type for expiration time: %T", v)
		return true
	}

	// Check if the token is expired
	expirationTime := time.Unix(exp, 0)
	if expirationTime.Before(time.Now()) {
		logger.Log.Warnf("Token is expired: expiration time was %v", expirationTime)
		return true
	}

	logger.Log.Infof("Token is valid: expiration time is %v", expirationTime)
	return false
}

// Helper function to decode a Base64URL-encoded segment.
func decodeSegment(seg string) ([]byte, error) {
	// Base64URL decode the segment
	decoded, err := base64.RawURLEncoding.DecodeString(seg)
	if err != nil {
		return nil, fmt.Errorf("error decoding segment: %v", err)
	}
	return decoded, nil
}

// RefreshAccessToken uses the refresh token to obtain a new access token from the server.
func (tm *TokenManager) RefreshAccessToken(refreshToken string) (string, error) {
	// Create a context with a timeout to avoid hanging indefinitely
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Make the gRPC call to refresh the token
	resp, err := tm.client.RefreshToken(ctx, &auth.RefreshTokenRequest{RefreshToken: refreshToken})
	if err != nil {
		logger.Log.Errorf("Failed to refresh access token: %v", err)
		return "", err
	}

	// Return the new access token received from the server
	return resp.AccessToken, nil
}
