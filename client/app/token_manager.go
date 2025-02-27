package app

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/johnkhk/cli_chat_app/genproto/auth"
)

// For testing
type TimeProvider interface {
	Now() time.Time
}

type RealTimeProvider struct{}

func (RealTimeProvider) Now() time.Time {
	return time.Now()
}

// TokenManager handles all operations related to managing tokens.
type TokenManager struct {
	TimeProvider TimeProvider
	filePath     string
	AuthClient   *AuthClient // Reference to the parent AuthClient
}

func NewTokenManager(filePath string, client *AuthClient) *TokenManager {
	// Ensure the directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create directory %s: %v", dir, err))
	}
	return &TokenManager{
		filePath:     filePath,
		AuthClient:   client,
		TimeProvider: RealTimeProvider{}, // Default TimeProvider
	}
}

// GetAccessToken returns a valid access token, refreshing it if necessary.
func (tm *TokenManager) GetAccessToken() (string, error) {
	// Read the current tokens
	accessToken, refreshToken, err := tm.ReadTokens()
	if err != nil {
		return "", fmt.Errorf("failed to read tokens: %w", err)
	}

	// Check if the access token is expired
	expired, err := tm.IsTokenExpired(accessToken)
	if err != nil {
		return "", fmt.Errorf("failed to check if token is expired: %w", err)
	}

	if expired {
		// Attempt to refresh the access token
		accessToken, err = tm.RefreshAccessToken(refreshToken)
		if err != nil {
			return "", fmt.Errorf("failed to refresh access token: %w", err)
		}

		// Store the new access token and reuse the refresh token
		if err := tm.StoreTokens(accessToken, refreshToken); err != nil {
			return "", fmt.Errorf("failed to store refreshed tokens: %w", err)
		}
	}

	// Return the valid access token
	return accessToken, nil
}

// TryAutoLogin attempts to automatically log in the user using stored tokens.
// Actually, only refreshes the access token if it is expired.
// Otherwise, does nothing.
func (tm *TokenManager) TryAutoLogin() error {
	_, err := tm.GetAccessToken() // Attempt to get a valid access token
	if err != nil {
		return fmt.Errorf("failed to auto-login: %w", err)
	}
	err = tm.AuthClient.PostLoginTasks()
	return err
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
func (tm *TokenManager) IsTokenExpired(tokenString string) (bool, error) {
	// Split the JWT into its parts: header, payload, signature
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return true, fmt.Errorf("invalid token format: expected 3 parts but got %d", len(parts))
	}

	// Decode the payload part (Base64URL encoded)
	payload, err := decodeSegment(parts[1])
	if err != nil {
		return true, fmt.Errorf("failed to decode token payload: %w", err)
	}

	// Define a struct with a custom Unmarshal function to handle both string and int
	var claims struct {
		Exp interface{} `json:"exp"` // Use interface{} to handle both int and string
	}

	// Parse the JSON payload to extract the "exp" field
	if err := json.Unmarshal(payload, &claims); err != nil {
		return true, fmt.Errorf("failed to unmarshal token claims: %w", err)
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
			return true, fmt.Errorf("failed to parse expiration time from string: %w", err)
		}
	default:
		return true, fmt.Errorf("unexpected type for expiration time: %T", v)
	}

	// Check if the token is expired
	expirationTime := time.Unix(exp, 0)
	if expirationTime.Before(tm.TimeProvider.Now()) {
		return true, nil // Token is expired
	}

	return false, nil // Token is valid
}

// Helper function to decode a Base64URL-encoded segment.
func decodeSegment(seg string) ([]byte, error) {
	// Base64URL decode the segment
	decoded, err := base64.RawURLEncoding.DecodeString(seg)
	if err != nil {
		return nil, fmt.Errorf("error decoding segment: %w", err)
	}
	return decoded, nil
}

// RefreshAccessToken uses the refresh token to obtain a new access token from the server.
func (tm *TokenManager) RefreshAccessToken(refreshToken string) (string, error) {
	// Create a context with a timeout to avoid hanging indefinitely
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Make the gRPC call to refresh the token
	resp, err := tm.AuthClient.Client.RefreshToken(ctx, &auth.RefreshTokenRequest{RefreshToken: refreshToken})
	if err != nil {
		return "", fmt.Errorf("failed to refresh access token: %w", err)
	}

	// Return the new access token received from the server
	return resp.AccessToken, nil
}

// SetClient allows updating the gRPC client.
func (tm *TokenManager) SetClient(client *AuthClient) {
	tm.AuthClient = client
}

// GetClaimsFromAccessToken reads the stored access token and parses its JWT claims.
func (tm *TokenManager) GetClaimsFromAccessToken() (map[string]interface{}, error) {
	// Retrieve the stored access token using ReadTokens()
	accessToken, _, err := tm.ReadTokens()
	if err != nil {
		return nil, fmt.Errorf("failed to read tokens: %w", err)
	}

	// Split the JWT into its parts: header, payload, signature
	parts := strings.Split(accessToken, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format: expected 3 parts but got %d", len(parts))
	}

	// Decode the payload part (Base64URL encoded)
	payload, err := decodeSegment(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode token payload: %w", err)
	}

	// Parse the JSON payload to extract the claims
	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token claims: %w", err)
	}

	return claims, nil
}

func (tm *TokenManager) GetUserIdFromAccessToken() (uint32, error) {
	claims, err := tm.GetClaimsFromAccessToken()
	if err != nil {
		return 0, fmt.Errorf("failed to get claims from access token: %w", err)
	}

	// The "sub" claim is most likely a string, so you need to handle it as such
	sub, ok := claims["sub"].(string)
	if !ok {
		return 0, fmt.Errorf("UserId (sub) claim not found or not a string")
	}

	// Convert the string to uint32
	userID, err := strconv.ParseUint(sub, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to parse UserId (sub) claim: %w", err)
	}

	return uint32(userID), nil
}
