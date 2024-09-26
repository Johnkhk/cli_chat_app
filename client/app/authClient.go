package app

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/johnkhk/cli_chat_app/client/e2ee/store"
	"github.com/johnkhk/cli_chat_app/genproto/auth"
)

// AuthClient encapsulates the gRPC client and logger for authentication services.
type AuthClient struct {
	Client       auth.AuthServiceClient
	Logger       *logrus.Logger
	TokenManager *TokenManager
	AppDirPath   string
	SqliteStore  *store.SQLiteStore
}

// RegisterUser sends a registration request to the server.
func (c *AuthClient) RegisterUser(username, password string) error {
	req := &auth.RegisterRequest{
		Username: username,
		Password: password,
	}

	resp, err := c.Client.RegisterUser(context.Background(), req)
	if err != nil {
		c.Logger.Errorf("Failed to register user: %v", err)
		return fmt.Errorf("Failed to register user: %v", err)
	}

	if resp.Success {
		c.Logger.Infof("Registration successful: %s", resp.Message)
		return nil
	} else {
		c.Logger.Infof("Registration failed: %s", resp.Message)
		return fmt.Errorf("Registration failed: %s", resp.Message)
	}
}

// LoginUser sends a login request to the server.
func (c *AuthClient) LoginUser(username, password string) error {
	req := &auth.LoginRequest{
		Username: username,
		Password: password,
	}

	resp, err := c.Client.LoginUser(context.Background(), req)
	if err != nil {
		c.Logger.Warnf("Failed to login: %v", err)
		return fmt.Errorf("failed to login: %v", err)
	}

	if resp.Success {
		c.Logger.Infof("Login successful: %s", resp.Message)

		// Store jwt tokens after successful login
		if err := c.TokenManager.StoreTokens(resp.AccessToken, resp.RefreshToken); err != nil {
			c.Logger.Errorf("Failed to store tokens: %v", err)
			return fmt.Errorf("failed to store tokens: %v", err)
		}
		c.Logger.Infof("JWT Tokens stored successfully to %s", c.TokenManager.filePath)

		// // Check if the user already has a private key stored locally
		// privateKeyPath, err := c.GetPrivateKeyPath(username)
		// if err != nil {
		// 	c.Logger.Errorf("Failed to get private key path: %v", err)
		// 	return fmt.Errorf("failed to get private key path: %v", err)
		// }

		// if !fileExists(privateKeyPath) {
		// 	c.Logger.Infof("Private key not found for user %s. Generating new key pair...", username)

		// 	// Generate a new public-private key pair for the user
		// 	identityKeyPair, err := identity.GenerateKeyPair(rand.Reader)
		// 	if err != nil {
		// 		c.Logger.Errorf("Failed to generate identity key pair: %v", err)
		// 		return fmt.Errorf("failed to generate identity key pair: %v", err)
		// 	}
		// 	publicKeyBytes := identityKeyPair.PublicKey().Bytes()

		// 	// Upload the public key to the server
		// 	err = c.UploadPublicKeyToServer(username, publicKeyBytes)
		// 	if err != nil {
		// 		c.Logger.Errorf("Failed to upload public key: %v", err)
		// 		return fmt.Errorf("failed to upload public key: %v", err)
		// 	}
		// 	// Store the private key locally
		// 	err = c.storePrivateKey(username, identityKeyPair.PrivateKey().Bytes())
		// 	if err != nil {
		// 		c.Logger.Errorf("Failed to store private key: %v", err)
		// 		return fmt.Errorf("failed to store private key: %v", err)
		// 	}
		// 	// Store the public key locally
		// 	// TODO

		// 	c.Logger.Infof("Public key uploaded successfully for user %s", username)
		// } else {
		// 	c.Logger.Infof("Private key already exists for user %s", username)
		// }

		// Create local identity
		err = c.OnLogIn(resp.UserId)
		if err != nil {
			return fmt.Errorf("failed to create local identity: %v", err)
		}

		return nil
	} else {
		c.Logger.Infof("Login failed: %s", resp.Message)
		return fmt.Errorf("login failed: %s", resp.Message)
	}
}

// ////////////////////////// Encryption key management////
// OnLogIn checks if the user-device is already registered, and if not, generates keys and uploads them.
func (c *AuthClient) OnLogIn(userID uint32) error {
	db := c.SqliteStore.DB

	// Get MAC address
	macAddress, err := getMACAddress()
	if err != nil {
		return fmt.Errorf("failed to get MAC address: %v", err)
	}

	// Convert MAC address to uint32 for deviceID
	deviceID, err := store.MacToUint32(macAddress)
	if err != nil {
		return fmt.Errorf("failed to convert MAC address to uint32: %v", err)
	}

	// Generate registration ID using userID and deviceID
	registrationID := store.GenerateRegistrationID(userID, deviceID)

	// Check if the registration ID already exists in the database
	var existingID uint32
	err = db.QueryRow("SELECT registration_id FROM local_identity WHERE registration_id = ?", registrationID).Scan(&existingID)
	if err != nil {
		if err == sql.ErrNoRows {
			// No existing record, proceed to generate keys
			c.Logger.Infoln("No existing registration found. Creating new local identity and generating keys...")

			// Create and store the local identity (private parts)
			c.Logger.Infof("Generating new local identity for registration ID: %d", registrationID)
			local_identity, err := c.SqliteStore.CreateLocalIdentity(registrationID)
			if err != nil {
				return fmt.Errorf("failed to create local identity: %v", err)
			}
			// Create the PublicKeyUploadRequest
			req := &auth.PublicKeyUploadRequest{
				IdentityKey:           local_identity.IdentityPublicKey,
				PreKeyId:              local_identity.PreKeyID,
				PreKey:                local_identity.PreKeyPublicKey,
				SignedPreKeyId:        local_identity.SignedPreKeyID,
				SignedPreKey:          local_identity.SignedPreKeyPublicKey,
				SignedPreKeySignature: local_identity.Signature,
				RegistrationId:        registrationID,
				DeviceId:              deviceID,
				OneTimePreKeys:        nil,
			}

			// Set a timeout for the request
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// Send the request to the server
			res, err := c.Client.UploadPublicKeys(ctx, req)
			if err != nil {
				c.Logger.Fatalf("Failed to upload public keys: %v", err)
			}
			if res.Success {
				c.Logger.Infof("Public keys uploaded successfully for registration ID: %d", registrationID)
			} else {
				c.Logger.Infof("Failed to upload public keys: %s", res)
			}

		} else {
			// Some other error occurred during the query
			return fmt.Errorf("failed to query registration ID: %v", err)
		}
	} else {
		// Registration ID already exists, device is already registered
		fmt.Println("Device is already registered with registration ID:", registrationID)
	}

	return nil
}

///////////////////////////////////// Key Management /////////////

// /////////////// Encryption Key Management ///////////////

// func (c *AuthClient) storePrivateKey(username string, privateKey []byte) error {

// 	// Define the file path for storing the private key
// 	filePath, err := c.GetPrivateKeyPath(username)
// 	if err != nil {
// 		return fmt.Errorf("failed to get private key path: %v", err)
// 	}

// 	// Write the private key bytes to the file
// 	err = os.WriteFile(filePath, privateKey, 0600) // 0600 ensures the file is only readable and writable by the owner
// 	if err != nil {
// 		return fmt.Errorf("failed to write private key to file: %v", err)
// 	}

// 	return nil
// }

// // UploadPublicKeyToServer uploads the public key of the user to the server.
// func (c *AuthClient) UploadPublicKeyToServer(username string, publicKey []byte) error {
// 	req := &auth.UploadPublicKeyRequest{
// 		Username:  username,
// 		PublicKey: publicKey,
// 	}

// 	// Make the RPC call to upload the public key
// 	resp, err := c.Client.UploadPublicKey(context.Background(), req)
// 	if err != nil {
// 		c.Logger.Errorf("Failed to upload public key: %v", err)
// 		return fmt.Errorf("failed to upload public key: %w", err)
// 	}

// 	if !resp.Success {
// 		c.Logger.Infof("Failed to upload public key: %s", resp.Message)
// 		return fmt.Errorf("failed to upload public key: %s", resp.Message)
// 	}

// 	c.Logger.Infof("Public key uploaded successfully for user: %s", username)
// 	return nil
// }

// // fileExists checks if a file exists.
// func fileExists(path string) bool {
// 	_, err := os.Stat(path)
// 	return err == nil
// }

// func (c *AuthClient) GetPrivateKeyPath(username string) (string, error) {
// 	privateKeyPath := filepath.Join(c.AppDirPath, fmt.Sprintf("%s_identity_private_key.pem", username))
// 	return privateKeyPath, nil
// }

// // GetPublicKey sends a request to the gRPC server to retrieve the public key for a user.
// func (c *AuthClient) GetPublicKey(userID int32) (*auth.GetPublicKeyResponse, error) {
// 	// Step 1: Create a request object.
// 	req := &auth.GetPublicKeyRequest{
// 		UserId: userID,
// 	}

// 	// Step 2: Send the request to the gRPC server.
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Set a timeout for the request.
// 	defer cancel()

// 	resp, err := c.Client.GetPublicKey(ctx, req)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get public key: %v", err)
// 	}

// 	// Step 3: Return the response from the server.
// 	return resp, nil
// }
