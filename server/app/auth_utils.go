package app

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// TokenValidator interface defines the method for validating tokens.
type TokenValidator interface {
	ValidateToken(token string) (userID string, username string, err error)
}

// JWTTokenValidator is a struct that implements the TokenValidator interface using JWT.
type JWTTokenValidator struct {
	secretKey string
}

// NewJWTTokenValidator creates a new instance of JWTTokenValidator.
func NewJWTTokenValidator(secretKey string) *JWTTokenValidator {
	return &JWTTokenValidator{secretKey: secretKey}
}

// ValidateToken validates the JWT token and extracts the user ID and username.
func (v *JWTTokenValidator) ValidateToken(tokenString string) (string, string, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is correct
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the key for verification
		return []byte(v.secretKey), nil
	})

	if err != nil {
		return "", "", fmt.Errorf("failed to parse token: %w", err)
	}

	// Validate the token claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check expiration
		if exp, ok := claims["exp"].(float64); ok {
			if time.Unix(int64(exp), 0).Before(time.Now()) {
				return "", "", fmt.Errorf("token is expired")
			}
		}

		// Extract user ID
		userID, ok := claims["sub"].(string)
		if !ok {
			return "", "", fmt.Errorf("invalid user ID in token claims")
		}

		// Extract username
		username, ok := claims["username"].(string)
		if !ok {
			return "", "", fmt.Errorf("invalid username in token claims")
		}

		return userID, username, nil
	}

	return "", "", fmt.Errorf("invalid token")
}

// isValidRefreshToken validates the provided refresh token.
func isValidRefreshToken(token string) bool {
	// Implement your refresh token validation logic here.
	return true // Simplified for demonstration
}

// Helper function to generate a new access token with a specified expiration duration.
func generateAccessToken(userID int64, username string, expirationDuration time.Duration) (string, error) {
	secretKey := os.Getenv("CLI_CHAT_APP_JWT_SECRET_KEY")
	if secretKey == "" {
		return "", fmt.Errorf("JWT secret key is not set")
	}

	// Generate minimal randomness: a single random byte
	randomByte := make([]byte, 1) // 1 byte = 8 bits of randomness
	_, err := rand.Read(randomByte)
	if err != nil {
		return "", fmt.Errorf("failed to generate random value: %v", err)
	}
	randomValue := hex.EncodeToString(randomByte) // Convert to a minimal hex string
	// Define token claims using the user ID as the subject.
	claims := jwt.MapClaims{
		"sub":      fmt.Sprintf("%d", userID),                 // Use user ID as subject
		"username": username,                                  // Add username to claims
		"exp":      time.Now().Add(expirationDuration).Unix(), // Token expires based on the given duration
		"nonce":    randomValue,                               // Add a minimal random claim to ensure uniqueness

	}

	// Create a new token object using the signing method and claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with your secret key.
	accessToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

// Helper function to generate a new refresh token with a specified expiration duration.
func generateRefreshToken(userID int64, username string, expirationDuration time.Duration) (string, error) {
	secretKey := os.Getenv("CLI_CHAT_APP_JWT_SECRET_KEY")
	if secretKey == "" {
		return "", fmt.Errorf("JWT secret key is not set")
	}

	// Generate minimal randomness: a single random byte
	randomByte := make([]byte, 1) // 1 byte = 8 bits of randomness
	_, err := rand.Read(randomByte)
	if err != nil {
		return "", fmt.Errorf("failed to generate random value: %v", err)
	}
	randomValue := hex.EncodeToString(randomByte) // Convert to a minimal hex string
	// Define refresh token claims using the user ID as the subject.
	claims := jwt.MapClaims{
		"sub":      fmt.Sprintf("%d", userID),                 // Use user ID as subject
		"username": username,                                  // Add username to claims
		"exp":      time.Now().Add(expirationDuration).Unix(), // Refresh token expires based on the given duration
		"nonce":    randomValue,                               // Add a minimal random claim to ensure uniqueness
	}

	// Create a new token object using the signing method and claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with your secret key.
	refreshToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

// Helper function to validate and parse the refresh token.
func parseAndValidateRefreshToken(tokenString string) (int64, string, error) {
	secretKey := os.Getenv("CLI_CHAT_APP_JWT_SECRET_KEY")
	if secretKey == "" {
		return 0, "", fmt.Errorf("JWT secret key is not set")
	}

	// Parse the token.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is correct.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the key for verification.
		return []byte(secretKey), nil
	})

	if err != nil {
		return 0, "", err
	}

	// Check if the token is valid.
	if !token.Valid {
		return 0, "", fmt.Errorf("invalid token")
	}

	// Extract the claims from the token.
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, "", fmt.Errorf("invalid token claims")
	}

	// Extract the user ID from the claims.
	userID, ok := claims["sub"].(string)
	if !ok {
		return 0, "", fmt.Errorf("invalid user ID in token claims")
	}

	// Parse the user ID string to int64.
	parsedUserID, err := parseInt64(userID)
	if err != nil {
		return 0, "", fmt.Errorf("failed to parse user ID: %v", err)
	}

	// Extract the username from the claims (if added).
	username, ok := claims["username"].(string)
	if !ok {
		return 0, "", fmt.Errorf("invalid username in token claims")
	}

	return parsedUserID, username, nil
}

// Helper function to parse a string to int64.
func parseInt64(s string) (int64, error) {
	var id int64
	_, err := fmt.Sscanf(s, "%d", &id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
