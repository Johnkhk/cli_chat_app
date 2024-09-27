package app

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"time"

	"github.com/Johnkhk/libsignal-go/protocol/address"
	"github.com/Johnkhk/libsignal-go/protocol/curve"
	"github.com/Johnkhk/libsignal-go/protocol/identity"
	"github.com/Johnkhk/libsignal-go/protocol/prekey"
	"github.com/Johnkhk/libsignal-go/protocol/session"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/johnkhk/cli_chat_app/client/e2ee/store"
	"github.com/johnkhk/cli_chat_app/genproto/auth"
	"github.com/johnkhk/cli_chat_app/genproto/chat"
)

// ChatClient encapsulates the gRPC client for chat services.
type ChatClient struct {
	Client     chat.ChatServiceClient                // gRPC client for chat service
	AuthClient *AuthClient                           // Reference to AuthClient for authentication purposes
	store      *store.SQLiteStore                    // Access to session and identity stores
	Logger     *logrus.Logger                        // Logger for logging messages and errors
	Stream     chat.ChatService_StreamMessagesClient // Persistent gRPC stream for sending messages
}

// OpenPersistentStream opens a persistent gRPC stream for sending and receiving messages.
func (cc *ChatClient) OpenPersistentStream(ctx context.Context) error {
	// Open a new gRPC stream to the chat server for message handling.
	stream, err := cc.Client.StreamMessages(ctx)
	if err != nil {
		return fmt.Errorf("failed to open message stream: %v", err)
	}

	// Save the stream to the ChatClient instance for future use.
	cc.Stream = stream

	cc.Logger.Info("Persistent gRPC stream successfully opened.")
	return nil
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
	// Check if a session already exists with the recipient and device
	remoteAddress := address.Address{
		Name:     fmt.Sprintf("%d", recipientID),
		DeviceID: address.DeviceID(deviceID),
	}

	_, exists, err := cc.store.SessionStore().Load(ctx, remoteAddress)
	if err != nil {
		return fmt.Errorf("failed to load existing session: %v", err)
	}

	// If a session already exists, return as there's no need to initialize again
	if exists {
		cc.Logger.Infof("Session already exists with recipient %d and device %d", recipientID, deviceID)
		return nil
	}

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
		RemoteAddress:    remoteAddress,
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

// SendMessage encrypts a message and sends it to the recipient through the chat service.
func (cc *ChatClient) SendMessage(ctx context.Context, recipientID, deviceID uint32, plaintext []byte) error {
	// Ensure that the persistent stream is open
	if cc.Stream == nil {
		return fmt.Errorf("no active stream found. Ensure that openPersistentStream has been called.")
	}

	// Load the session with the recipient and device
	remoteAddress := address.Address{
		Name:     fmt.Sprintf("%d", recipientID),
		DeviceID: address.DeviceID(deviceID),
	}

	_, exists, err := cc.store.SessionStore().Load(ctx, remoteAddress)
	if err != nil {
		return fmt.Errorf("failed to load session: %v", err)
	}
	if !exists {
		return fmt.Errorf("no session found with recipient %d and device %d", recipientID, deviceID)
	}

	// Create a session object from the record
	session := &session.Session{
		RemoteAddress:    remoteAddress,
		SessionStore:     cc.store.SessionStore(),
		IdentityKeyStore: cc.store.IdentityStore(),
	}

	// Encrypt the plaintext message using the session
	ciphertext, err := session.EncryptMessage(ctx, plaintext)
	if err != nil {
		return fmt.Errorf("failed to encrypt message: %v", err)
	}

	// Create a new message request with the encrypted content
	msgRequest := &chat.MessageRequest{
		RecipientId:      recipientID,                     // Set recipient ID
		EncryptedMessage: ciphertext.Bytes(),              // Set encrypted message
		MessageId:        uuid.NewString(),                // Generate a unique message ID
		Timestamp:        time.Now().Format(time.RFC3339), // Timestamp in ISO 8601 format
	}

	// Send the message using the persistent stream
	if err := cc.Stream.Send(msgRequest); err != nil {
		return fmt.Errorf("failed to send message request: %v", err)
	}

	cc.Logger.Infof("Message sent to recipient %d successfully", recipientID)
	return nil
}

// listenForMessages continuously listens for messages on the open stream.
func (cc *ChatClient) listenForMessages() {
	for {
		resp, err := cc.Stream.Recv()
		if err == io.EOF {
			cc.Logger.Info("Stream closed by server")
			return
		}
		if err != nil {
			cc.Logger.Errorf("Failed to receive message: %v", err)
			return
		}

		cc.Logger.Infof("Received message response: %s", resp.MessageId)

		// Here you can handle different types of responses.
		// If it’s a delivery confirmation, update the UI or logs.
		// If it’s a new message from a sender, decrypt and process it.

		// For example, handle a delivered message:
		if resp.Status == "delivered" {
			cc.Logger.Infof("Message %s was delivered successfully at %s", resp.MessageId, resp.Timestamp)
		}
	}
}
