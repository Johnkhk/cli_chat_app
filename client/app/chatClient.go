package app

import (
	"context"
	"crypto/rand"
	"fmt"

	"github.com/Johnkhk/libsignal-go/protocol/address"
	"github.com/Johnkhk/libsignal-go/protocol/curve"
	"github.com/Johnkhk/libsignal-go/protocol/identity"
	"github.com/Johnkhk/libsignal-go/protocol/prekey"
	"github.com/Johnkhk/libsignal-go/protocol/session"
	"github.com/sirupsen/logrus"

	"github.com/johnkhk/cli_chat_app/client/e2ee/store"
	"github.com/johnkhk/cli_chat_app/genproto/auth"
	"github.com/johnkhk/cli_chat_app/genproto/chat"
)

// ChatClient encapsulates the gRPC client for chat services.
type ChatClient struct {
	Client     chat.ChatServiceClient // gRPC client for chat service
	AuthClient *AuthClient            // gRPC client for authentication service
	store      *store.SQLiteStore     // Access to session and identity stores
	Logger     *logrus.Logger         // Logger for logging messages and errors
}

// /////////////////////////////////////////////////////////////
func (cc *ChatClient) InitializeSessionForRecipient(ctx context.Context, recipientID uint32) error {
	// Fetch the recipient's pre-key bundle (assuming single device or unknown device)
	bundle, err := cc.AuthClient.GetPublicKeyBundle(recipientID, 0) // Passing 0 to use the first available device
	if err != nil {
		return fmt.Errorf("failed to fetch recipient's pre-key bundle: %v", err)
	}

	// Initialize the session with the recipient using the pre-key bundle
	err = cc.initializeSessionWithDevice(ctx, recipientID, 0, bundle)
	if err != nil {
		return fmt.Errorf("failed to initialize session with recipient %d: %v", recipientID, err)
	}

	cc.Logger.Infof("Successfully initialized session with recipient %d", recipientID)
	return nil
}

func (cc *ChatClient) initializeSessionWithDevice(ctx context.Context, recipientID, deviceID uint32, bundle *auth.PublicKeyBundleResponse) error {
	// // Retrieve our identity key pair (Alice's key pair) from the store
	// aliceIdentityKeyPair := cc.store.IdentityStore().KeyPair(ctx)

	// // Generate a new ephemeral base key for Alice (ourBaseKey)
	// ourBaseKeyPair, err := curve.GenerateKeyPair(rand.Reader)
	// if err != nil {
	//     return fmt.Errorf("failed to generate base key pair: %v", err)
	// }

	// Prepare Bob's pre-key bundle for session initialization
	theirIdentityKey, err := identity.NewKey(bundle.IdentityKey)
	if err != nil {
		return fmt.Errorf("failed to create identity key: %v", err)
	}

	theirSignedPreKey, err := curve.NewPublicKey(bundle.SignedPreKey)
	if err != nil {
		return fmt.Errorf("failed to create signed pre-key: %v", err)
	}

	var theirOneTimePreKey curve.PublicKey
	if len(bundle.PreKey) > 0 {
		theirOneTimePreKey, err = curve.NewPublicKey(bundle.PreKey)
		if err != nil {
			return fmt.Errorf("failed to create one-time pre-key: %v", err)
		}
	}

	// Create a pre-key bundle for Bob using the values from the fetched bundle
	bobPreKeyBundle := &prekey.Bundle{
		RegistrationID:        bundle.RegistrationId,
		DeviceID:              address.DeviceID(deviceID),
		PreKeyID:              store.To(prekey.ID(bundle.PreKeyId)),
		PreKeyPublic:          theirOneTimePreKey, // Use the optional one-time pre-key if available
		SignedPreKeyID:        prekey.ID(bundle.SignedPreKeyId),
		SignedPreKeyPublic:    theirSignedPreKey,
		SignedPreKeySignature: bundle.SignedPreKeySignature,
		IdentityKey:           theirIdentityKey,
	}

	// Initialize Alice's session with Bob using Bob's pre-key bundle
	aliceSession := &session.Session{
		RemoteAddress: address.Address{
			Name:     fmt.Sprintf("%d", recipientID),
			DeviceID: address.DeviceID(deviceID),
		},
		SessionStore:     cc.store.SessionStore(),
		IdentityKeyStore: cc.store.IdentityStore(),
	}

	// Call ProcessPreKeyBundle to initialize the session using Bob's bundle
	if err := aliceSession.ProcessPreKeyBundle(ctx, rand.Reader, bobPreKeyBundle); err != nil {
		return fmt.Errorf("failed to process pre-key bundle: %v", err)
	}

	cc.Logger.Infof("Session initialized with recipient %d and device %d", recipientID, deviceID)
	return nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////

// // SendMessage handles creating and sending a message via the gRPC stream.
// func (c *ChatClient) SendMessage(ctx context.Context, recipientID int32, messageContent string, filePath string) error {
// 	// Create a new stream for sending messages
// 	stream, err := c.Client.StreamMessages(ctx)
// 	if err != nil {
// 		c.Logger.Errorf("Failed to establish stream: %v", err)
// 		return err
// 	}

// 	// Generate a unique message ID
// 	messageID := uuid.New().String()
// 	timestamp := time.Now().Format(time.RFC3339)

// 	// Encrypt the message content (this should be done using your encryption method)
// 	encryptedMessage := encryptMessage([]byte(messageContent)) // Implement your encryption method here

// 	// Prepare the file if provided
// 	var fileContent []byte
// 	var fileName, fileType string
// 	var fileSize int64

// 	if filePath != "" {
// 		// Load and encrypt file content
// 		fileContent, fileName, fileType, fileSize, err = prepareFile(filePath)
// 		if err != nil {
// 			c.Logger.Errorf("Failed to prepare file: %v", err)
// 			return err
// 		}
// 	}

// 	// Create a MessageRequest struct with the required fields
// 	req := &chat.MessageRequest{
// 		RecipientId:      recipientID,
// 		EncryptedMessage: encryptedMessage,
// 		MessageId:        messageID,
// 		Timestamp:        timestamp,
// 		FileContent:      fileContent,
// 		FileName:         fileName,
// 		FileType:         fileType,
// 		FileSize:         fileSize,
// 	}

// 	// Send the message request through the stream
// 	if err := stream.Send(req); err != nil {
// 		c.Logger.Errorf("Failed to send message: %v", err)
// 		return err
// 	}

// 	c.Logger.Infof("Message sent successfully: %s", messageID)

// 	// Listen for the server's response
// 	go func() {
// 		for {
// 			resp, err := stream.Recv()
// 			if err != nil {
// 				c.Logger.Errorf("Failed to receive response: %v", err)
// 				return
// 			}

// 			// Handle the response from the server
// 			if resp.Status == "delivered" {
// 				c.Logger.Infof("Message delivered: %s", resp.MessageId)
// 			} else if resp.Status == "error" {
// 				c.Logger.Errorf("Error delivering message: %s, error: %s", resp.MessageId, resp.ErrorMessage)
// 			}
// 		}
// 	}()

// 	return nil
// }

// // EncryptMessage is a placeholder function for encrypting the message content.
// func encryptMessage(content []byte) []byte {
// 	// Replace this with your actual encryption logic
// 	return content // For now, just returning the original content
// }

// // PrepareFile handles reading the file, encrypting it, and preparing it for sending.
// func prepareFile(filePath string) ([]byte, string, string, int64, error) {
// 	// Implement file loading and encryption logic here
// 	// Example: Load file from filePath, encrypt the content, and return file details
// 	return nil, "", "", 0, nil // Placeholder implementation
// }
