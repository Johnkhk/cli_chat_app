package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"github.com/johnkhk/cli_chat_app/genproto/auth"
	"github.com/johnkhk/cli_chat_app/server/storage"
)

// AuthServer implements the AuthService.
type AuthServer struct {
	auth.UnimplementedAuthServiceServer
	DB                     *sql.DB
	Logger                 *logrus.Logger
	AccessTokenExpiration  time.Duration
	RefreshTokenExpiration time.Duration
}

// NewAuthServer creates a new AuthServer with the given dependencies.
func NewAuthServer(db *sql.DB, logger *logrus.Logger, accessTokenExpiration, refreshTokenExpiration time.Duration) *AuthServer {
	return &AuthServer{
		DB:                     db,
		Logger:                 logger,
		AccessTokenExpiration:  accessTokenExpiration,
		RefreshTokenExpiration: refreshTokenExpiration,
	}
}

// RegisterUser handles user registration requests.
func (s *AuthServer) RegisterUser(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	s.Logger.Infof("Registering new user: %s", req.Username)

	// Check if the user already exists
	var count int
	err := s.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", req.Username).Scan(&count)
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
	_, err = s.DB.Exec("INSERT INTO users (username, password_hash, created_at) VALUES (?, ?, NOW())",
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
	s.Logger.Infof("User login attempt: %s", req.Username)

	// Retrieve the user data from the database
	var user storage.User
	err := s.DB.QueryRow("SELECT id, username, password_hash FROM users WHERE username = ?", req.Username).Scan(&user.ID, &user.Username, &user.Password)
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
	accessToken, err := generateAccessToken(user.ID, user.Username, s.AccessTokenExpiration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %v", err)
	}

	refreshToken, err := generateRefreshToken(user.ID, user.Username, s.RefreshTokenExpiration)
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
	s.Logger.Info("Received refresh token request")
	refreshToken := req.RefreshToken

	// Validate and parse the refresh token
	userID, username, err := parseAndValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %v", err)
	}

	// Generate a new access token using the extracted user ID
	newAccessToken, err := generateAccessToken(userID, username, s.AccessTokenExpiration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %v", err)
	}

	// Return the new access token
	return &auth.RefreshTokenResponse{AccessToken: newAccessToken}, nil
}

// UploadPublicKey handles the RPC request to upload a public key.
func (s *AuthServer) UploadPublicKey(ctx context.Context, req *auth.UploadPublicKeyRequest) (*auth.UploadPublicKeyResponse, error) {
	// Step 1: Validate the input.
	if req.Username == "" || len(req.PublicKey) == 0 {
		return &auth.UploadPublicKeyResponse{
			Success: false,
			Message: "Invalid input: Username and public key are required.",
		}, fmt.Errorf("invalid input: missing username or public key")
	}

	// Step 2: Store the public key.
	err := s.saveUserPublicKey(req.Username, req.PublicKey)
	if err != nil {
		return &auth.UploadPublicKeyResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to upload public key: %v", err),
		}, err
	}

	// Step 3: Return a success response.
	return &auth.UploadPublicKeyResponse{
		Success: true,
		Message: "Public key uploaded successfully.",
	}, nil
}

func (s *AuthServer) saveUserPublicKey(username string, publicKey []byte) error {
	// Find the user ID based on the username
	var userID int
	err := s.DB.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
	if err != nil {
		return fmt.Errorf("failed to find user: %v", err)
	}

	// Insert or update the public key for the user
	_, err = s.DB.Exec(`
        INSERT INTO user_keys (user_id, identity_public_key) 
        VALUES (?, ?)
        ON DUPLICATE KEY UPDATE identity_public_key = VALUES(identity_public_key), updated_at = CURRENT_TIMESTAMP
    `, userID, publicKey)
	if err != nil {
		return fmt.Errorf("failed to store public key: %v", err)
	}

	log.Printf("Public key for user %s stored successfully.", username)
	return nil
}

// GetPublicKey handles the RPC request to retrieve a public key using user_id.
func (s *AuthServer) GetPublicKey(ctx context.Context, req *auth.GetPublicKeyRequest) (*auth.GetPublicKeyResponse, error) {
	// Step 1: Validate the input.
	if req.UserId == 0 {
		return &auth.GetPublicKeyResponse{
			Success: false,
			Message: "Invalid input: User ID is required.",
		}, fmt.Errorf("invalid input: missing user ID")
	}

	// Step 2: Retrieve the public key from the database.
	publicKey, err := s.getUserPublicKeyByID(req.UserId)
	if err != nil {
		return &auth.GetPublicKeyResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to retrieve public key: %v", err),
		}, err
	}

	// Step 3: Return the public key in the response.
	return &auth.GetPublicKeyResponse{
		Success:   true,
		PublicKey: publicKey,
		Message:   "Public key retrieved successfully.",
	}, nil
}

// getUserPublicKeyByID retrieves the public key for a user from the database using user_id.
func (s *AuthServer) getUserPublicKeyByID(userID int32) ([]byte, error) {
	// Find the public key based on the user ID.
	var publicKey []byte
	err := s.DB.QueryRow("SELECT identity_public_key FROM user_keys WHERE user_id = ?", userID).Scan(&publicKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("public key not found for user ID: %d", userID)
		}
		return nil, fmt.Errorf("failed to retrieve public key: %v", err)
	}

	log.Printf("Public key for user ID %d retrieved successfully.", userID)
	return publicKey, nil
}
