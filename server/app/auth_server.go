package app

import (
	"context"
	"database/sql"
	"fmt"
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
		UserId:       user.ID,
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

// Implement the server-side handler for UploadPublicKeys
func (s *AuthServer) UploadPublicKeys(ctx context.Context, req *auth.PublicKeyUploadRequest) (*auth.PublicKeyUploadResponse, error) {

	// Retrieve userID from the context
	// userID, ok := ctx.Value("userID").(int)
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return nil, fmt.Errorf("failed to retrieve userID from context")
	}

	// Begin a new transaction
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}

	// Insert the primary prekey bundle information, including user_id
	insertQuery := `
        INSERT INTO prekey_bundle (user_id, registration_id, device_id, identity_key, pre_key_id, pre_key, signed_pre_key_id, signed_pre_key, signed_pre_key_signature)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
    `
	_, err = tx.ExecContext(ctx, insertQuery,
		userID,                    // user_id from context
		req.RegistrationId,        // registration_id
		req.DeviceId,              // device_id
		req.IdentityKey,           // identity_key
		req.PreKeyId,              // pre_key_id
		req.PreKey,                // pre_key
		req.SignedPreKeyId,        // signed_pre_key_id
		req.SignedPreKey,          // signed_pre_key
		req.SignedPreKeySignature) // signed_pre_key_signature
	if err != nil {
		tx.Rollback() // Rollback in case of error
		s.Logger.Errorf("failed to insert prekey bundle: %v", err)
		return nil, fmt.Errorf("failed to insert prekey bundle: %v", err)
	}

	// Uncomment this part if you're handling One-Time PreKeys in the future
	// if len(req.OneTimePreKeys) > 0 {
	// 	oneTimePreKeyQuery := `
	//         INSERT INTO one_time_prekeys (pre_key_id, pre_key, bundle_id)
	//         VALUES (?, ?, LAST_INSERT_ID())  -- LAST_INSERT_ID() gets the last inserted bundle_id from prekey_bundle
	//     `
	// 	for _, oneTimePreKey := range req.OneTimePreKeys {
	// 		_, err := tx.ExecContext(ctx, oneTimePreKeyQuery,
	// 			oneTimePreKey.PreKeyId, // One-Time PreKey ID
	// 			oneTimePreKey.PreKey)   // One-Time PreKey Public Key
	// 		if err != nil {
	// 			tx.Rollback() // Rollback in case of error
	// 			return nil, fmt.Errorf("failed to insert one-time prekey: %v", err)
	// 		}
	// 	}
	// }

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	s.Logger.Infof("Public keys uploaded successfully for user: %s", userID)
	// Return a success response
	return &auth.PublicKeyUploadResponse{
		Success: true,
	}, nil
}

func (s *AuthServer) GetPublicKeyBundle(ctx context.Context, req *auth.PublicKeyBundleRequest) (*auth.PublicKeyBundleResponse, error) {
	// Fetch the prekey bundle from the database using the user_id (and optionally device_id)
	var identityKey, preKey, signedPreKey, signedPreKeySignature []byte
	var preKeyID, signedPreKeyID uint32
	// var oneTimePreKeys []auth.OneTimePreKey

	// Query the database to get the public keys for the requested user/device
	query := `SELECT identity_key, pre_key_id, pre_key, signed_pre_key_id, signed_pre_key, signed_pre_key_signature 
              FROM prekey_bundle WHERE user_id = ? AND device_id = ?`
	err := s.DB.QueryRow(query, req.GetUserId(), req.GetDeviceId()).Scan(
		&identityKey, &preKeyID, &preKey, &signedPreKeyID, &signedPreKey, &signedPreKeySignature)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch prekey bundle: %v", err)
	}

	// // Fetch one-time prekeys (if any)
	// // One-time prekeys can be fetched from a separate table and added to the response
	// oneTimePreKeysQuery := `SELECT pre_key_id, pre_key FROM one_time_prekeys WHERE user_id = ? AND device_id = ?`
	// rows, err := s.DB.Query(oneTimePreKeysQuery, req.GetUserId(), req.GetDeviceId())
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to fetch one-time prekeys: %v", err)
	// }
	// defer rows.Close()

	// for rows.Next() {
	// 	var oneTimePreKey auth.OneTimePreKey
	// 	if err := rows.Scan(&oneTimePreKey.PreKeyId, &oneTimePreKey.PreKey); err != nil {
	// 		return nil, fmt.Errorf("failed to scan one-time prekey: %v", err)
	// 	}
	// 	oneTimePreKeys = append(oneTimePreKeys,oneTimePreKey)
	// }

	// Return the public key bundle
	return &auth.PublicKeyBundleResponse{
		IdentityKey:           identityKey,
		PreKeyId:              preKeyID,
		PreKey:                preKey,
		SignedPreKeyId:        signedPreKeyID,
		SignedPreKey:          signedPreKey,
		SignedPreKeySignature: signedPreKeySignature,
		// OneTimePreKeys:        oneTimePreKeys,
		OneTimePreKeys: nil,
	}, nil
}
