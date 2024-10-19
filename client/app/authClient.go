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
	ParentClient *RpcClient // Reference to the parent RpcClient
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
func (c *AuthClient) LoginUser(username, password string) (error, uint32) {
	req := &auth.LoginRequest{
		Username: username,
		Password: password,
	}

	resp, err := c.Client.LoginUser(context.Background(), req)
	if err != nil {
		c.Logger.Warnf("Failed to login: %v", err)
		return fmt.Errorf("failed to login: %v", err), 0
	}

	if resp.Success {
		c.Logger.Infof("Login successful: %s", resp.Message)

		// Store jwt tokens after successful login
		if err := c.TokenManager.StoreTokens(resp.AccessToken, resp.RefreshToken); err != nil {
			c.Logger.Errorf("Failed to store tokens: %v", err)
			return fmt.Errorf("failed to store tokens: %v", err), 0
		}
		c.Logger.Infof("JWT Tokens stored successfully to %s", c.TokenManager.filePath)

		// Create local identity
		err = c.CreateLocalIdentityIfNewUserDevice(resp.UserId)
		if err != nil {
			return fmt.Errorf("failed to create local identity: %v", err), 0
		}

		// // Task A: Open the persistent stream
		// if err := c.ParentClient.ChatClient.OpenPersistentStream(context.Background()); err != nil {
		// 	return fmt.Errorf("failed to open persistent stream: %v", err), 0
		// }

		// // Task B: Create a context with cancel function to control lifecycle of message listening
		// listenCtx, cancelFunc := context.WithCancel(context.Background())
		// c.ParentClient.ChatClient.ListenCancelFunc = cancelFunc // Store cancel function in ChatClient for later use

		// // Task C: Listen for incoming messages
		// go c.ParentClient.ChatClient.listenForMessages(listenCtx)
		c.PostLoginTasks()

		return nil, resp.UserId
	} else {
		c.Logger.Infof("Login failed: %s", resp.Message)
		return fmt.Errorf("login failed: %s", resp.Message), 0
	}
}

// ////////////////////////// Encryption key management////
// CreateLocalIdentityIfNewUserDevice checks if the user-device is already registered, and if not, generates keys and uploads them.
func (c *AuthClient) CreateLocalIdentityIfNewUserDevice(userID uint32) error {
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

	// Check if the registration ID already exists in the local_identity database
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

func (c *AuthClient) GetPublicKeyBundle(userID, deviceID uint32) (*auth.PublicKeyBundleResponse, error) {
	req := &auth.PublicKeyBundleRequest{
		UserId:   userID,
		DeviceId: deviceID,
	}
	return c.Client.GetPublicKeyBundle(context.Background(), req)
}

// LogoutUser gracefully logs out the user by canceling the context and closing the stream.
func (c *AuthClient) LogoutUser() error {
	c.Logger.Info("Logging out user...")

	// Call the cancel function to stop listening for messages.
	if c.ParentClient.ChatClient.ListenCancelFunc != nil {
		c.Logger.Info("Canceling the message listener context")
		c.ParentClient.ChatClient.ListenCancelFunc() // Stop the listener
	} else {
		c.Logger.Warn("No message listener to stop.")
	}

	// Close the gRPC stream explicitly.
	if c.ParentClient.ChatClient.Stream != nil {
		c.Logger.Info("Closing the gRPC stream explicitly on logout")
		if err := c.ParentClient.ChatClient.Stream.CloseSend(); err != nil {
			c.Logger.Errorf("Failed to close gRPC stream: %v", err)
		} else {
			c.Logger.Info("Stream closed successfully on logout")
		}
	} else {
		c.Logger.Warn("No stream to close.")
	}

	c.Logger.Info("User logged out successfully.")
	return nil
}

// PostLoginTasks opens the stream and starts listening for messages.
func (c *AuthClient) PostLoginTasks() error {
	var err error
	c.ParentClient.CurrentUserID, err = c.TokenManager.GetUserIdFromAccessToken()
	if err != nil {
		return fmt.Errorf("failed to get user ID from access token: %v", err)
	}
	// Task A: Open the persistent stream.
	if err := c.ParentClient.ChatClient.OpenPersistentStream(context.Background()); err != nil {
		return fmt.Errorf("failed to open persistent stream: %v", err)
	}

	// Task B: Create a context with cancel function to control lifecycle of message listening.
	listenCtx, cancelFunc := context.WithCancel(context.Background())
	c.ParentClient.ChatClient.ListenCancelFunc = cancelFunc // Store cancel function in ChatClient for later use.

	// Task C: Listen for incoming messages.
	go c.ParentClient.ChatClient.listenForMessages(listenCtx)
	return nil
}
