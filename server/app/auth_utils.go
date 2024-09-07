// server/app/auth_utils.go

package app

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// isValidRefreshToken validates the provided refresh token.
func isValidRefreshToken(token string) bool {
	// Implement your refresh token validation logic here.
	// For example, you could parse the token, check its signature, and verify its claims.
	return true // Simplified for demonstration
}

// Helper function to generate a new access token
func generateAccessToken(userID int64, username string) (string, error) {
	// Define token claims, using the user ID as the subject
	claims := jwt.MapClaims{
		"sub":      fmt.Sprintf("%d", userID),            // Use user ID as subject
		"username": username,                             // Add username to claims
		"exp":      time.Now().Add(time.Hour * 1).Unix(), // Token expires in 1 hour
		// "exp": time.Now().Add(time.Second * 5).Unix(), // Token expires in 1 hour WORKS!
	}

	// Create a new token object using the signing method and claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with your secret key
	accessToken, err := token.SignedString([]byte(os.Getenv("CLI_CHAT_APP_JWT_SECRET_KEY"))) // Replace with your actual secret key
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

// Helper function to generate a new refresh token
func generateRefreshToken(userID int64, username string) (string, error) {
	// Define refresh token claims using the user ID as the subject
	claims := jwt.MapClaims{
		"sub":      fmt.Sprintf("%d", userID),                 // Use user ID as subject
		"username": username,                                  // Add username to claims
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(), // Refresh token expires in 7 days
	}

	// Create a new token object using the signing method and claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with your secret key
	refreshToken, err := token.SignedString([]byte(os.Getenv("CLI_CHAT_APP_JWT_SECRET_KEY"))) // Replace with your actual secret key
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

// Helper function to validate and parse the refresh token
func parseAndValidateRefreshToken(tokenString string) (int64, string, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is correct
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the key for verification
		return []byte(os.Getenv("CLI_CHAT_APP_JWT_SECRET_KEY")), nil
	})

	if err != nil {
		return 0, "", err
	}

	// Check if the token is valid
	if !token.Valid {
		return 0, "", fmt.Errorf("invalid token")
	}

	// Extract the claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, "", fmt.Errorf("invalid token claims")
	}

	// Extract the user ID from the claims
	userID, ok := claims["sub"].(string)
	if !ok {
		return 0, "", fmt.Errorf("invalid user ID in token claims")
	}

	// Parse the user ID string to int64
	parsedUserID, err := parseInt64(userID)
	if err != nil {
		return 0, "", fmt.Errorf("failed to parse user ID: %v", err)
	}

	// Extract the email from the claims (if added)
	username, ok := claims["username"].(string)
	if !ok {
		return 0, "", fmt.Errorf("invalid username in token claims")
	}

	return parsedUserID, username, nil
}

// Helper function to parse string to int64
func parseInt64(s string) (int64, error) {
	var id int64
	_, err := fmt.Sscanf(s, "%d", &id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
