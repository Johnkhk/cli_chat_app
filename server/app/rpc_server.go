package app

import (
	"context"
	"database/sql"
	"fmt"

	// Import the storage package

	"golang.org/x/crypto/bcrypt"

	"github.com/johnkhk/cli_chat_app/genproto/auth"
	"github.com/johnkhk/cli_chat_app/server/logger"
	"github.com/johnkhk/cli_chat_app/server/storage"
)

// AuthServer implements the AuthService.
type AuthServer struct {
	auth.UnimplementedAuthServiceServer
	DB *sql.DB
}

// RegisterUser handles user registration requests.
func (s *AuthServer) RegisterUser(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	// Check if the user already exists
	var count int
	err := s.DB.QueryRow("SELECT COUNT(*) FROM chat_users WHERE username = ?", req.Username).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("error checking user existence: %w", err)
	}
	if count > 0 {
		return &auth.RegisterResponse{
			Success: false,
			Message: "Username already exists",
		}, nil
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	// Save the user data to the database
	_, err = s.DB.Exec("INSERT INTO chat_users (username, password_hash, created_at) VALUES (?, ?, NOW())",
		req.Username, string(hashedPassword))
	if err != nil {
		return nil, fmt.Errorf("error saving user to database: %w", err)
	}

	return &auth.RegisterResponse{
		Success: true,
		Message: "User registered successfully",
	}, nil
}

// LoginUser handles user login requests.
func (s *AuthServer) LoginUser(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	// Retrieve the user data from the database
	var user storage.User
	err := s.DB.QueryRow("SELECT id, username, password_hash FROM chat_users WHERE username = ?", req.Username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("error retrieving user data: %w", err)
	}

	// Compare the password hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return &auth.LoginResponse{
			Success: false,
			Message: "Invalid username or password",
		}, nil
	}

	// Generate new access and refresh tokens using user ID as the subject
	accessToken, err := generateAccessToken(user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %v", err)
	}

	refreshToken, err := generateRefreshToken(user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %v", err)
	}

	return &auth.LoginResponse{
		Success:      true,
		Message:      "Login successful",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// RefreshToken handles token refresh requests.
func (s *AuthServer) RefreshToken(ctx context.Context, req *auth.RefreshTokenRequest) (*auth.RefreshTokenResponse, error) {
	logger.Log.Info("Received refresh token request")
	refreshToken := req.RefreshToken

	// Validate and parse the refresh token
	userID, username, err := parseAndValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %v", err)
	}

	// Generate a new access token using the extracted user ID
	newAccessToken, err := generateAccessToken(userID, username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %v", err)
	}

	// Return the new access token
	return &auth.RefreshTokenResponse{AccessToken: newAccessToken}, nil
}
