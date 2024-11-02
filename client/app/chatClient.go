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
	"github.com/Johnkhk/libsignal-go/protocol/message"
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
	Client           chat.ChatServiceClient                // gRPC client for chat service
	AuthClient       *AuthClient                           // Reference to AuthClient for authentication purposes
	Store            *store.SQLiteStore                    // Access to session and identity stores
	Logger           *logrus.Logger                        // Logger for logging messages and errors
	Stream           chat.ChatService_StreamMessagesClient // Persistent gRPC stream for sending messages
	ListenCancelFunc context.CancelFunc                    // Cancel function for stopping the message listener
	MessageChannel   chan *chat.MessageResponse            // Channel to send received messages
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

	_, exists, err := cc.Store.SessionStore().Load(ctx, remoteAddress)
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
		SessionStore:     cc.Store.SessionStore(),
		IdentityKeyStore: cc.Store.IdentityStore(),
	}

	// Call ProcessPreKeyBundle to initialize the session using Bob's bundle
	if err := aliceSession.ProcessPreKeyBundle(ctx, rand.Reader, bobPreKeyBundle); err != nil {
		return fmt.Errorf("failed to process pre-key bundle: %v", err)
	}

	cc.Logger.Infof("Session initialized with recipient %d and device %d", recipientID, deviceID)
	return nil
}

// SendMessage encrypts a message and sends it to the recipient through the chat service.
func (cc *ChatClient) EncryptMessage(ctx context.Context, recipientID, deviceID uint32, plaintext []byte) (message.Ciphertext, error) {
	// Ensure that the persistent stream is open
	if cc.Stream == nil {
		return nil, fmt.Errorf("no active stream found. Ensure that openPersistentStream has been called.")
	}

	// Load the session with the recipient and device
	remoteAddress := address.Address{
		Name:     fmt.Sprintf("%d", recipientID),
		DeviceID: address.DeviceID(deviceID),
	}

	_, exists, err := cc.Store.SessionStore().Load(ctx, remoteAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to load session: %v", err)
	}
	if !exists {
		cc.Logger.Infof("No session found with recipient %d. Initializing session...", recipientID)
		err = cc.InitializeSessionForRecipient(ctx, recipientID)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize session with recipient %d: %v", recipientID, err)
		}
	}

	// Create a session object from the record
	session := &session.Session{
		RemoteAddress:    remoteAddress,
		SessionStore:     cc.Store.SessionStore(),
		IdentityKeyStore: cc.Store.IdentityStore(),
	}

	cc.Logger.Infof("Using session to encrypt message for recipient %d", recipientID)

	// Encrypt the plaintext message using the session
	ciphertext, err := session.EncryptMessage(ctx, plaintext)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt message: %v", err)
	}

	return ciphertext, nil
}

func (cc *ChatClient) SendMessage(ctx context.Context, recipientID, deviceID uint32, plaintext []byte) error {

	ciphertext, err := cc.EncryptMessage(ctx, recipientID, deviceID, plaintext)
	if err != nil {
		return fmt.Errorf("failed to encrypt message: %v", err)
	}

	// Determine type of message (Signal or PreKey)
	var encryptionType chat.EncryptionType
	switch ciphertext.(type) {
	case *message.PreKey:
		encryptionType = chat.EncryptionType_PREKEY
	case *message.Signal:
		encryptionType = chat.EncryptionType_SIGNAL
	default:
		return fmt.Errorf("unknown encryption type: %T", ciphertext)
	}

	// Ensure that the persistent stream is open
	if cc.Stream == nil {
		return fmt.Errorf("no active stream found. Ensure that openPersistentStream has been called.")
	}

	// Create a new message request with the encrypted content
	msgRequest := &chat.MessageRequest{
		RecipientId:      recipientID,                     // Set recipient ID
		EncryptedMessage: ciphertext.Bytes(),              // Set encrypted message
		MessageId:        uuid.NewString(),                // Generate a unique message ID
		Timestamp:        time.Now().Format(time.RFC3339), // Timestamp in ISO 8601 format
		EncryptionType:   encryptionType,                  // Set the message type
	}

	// Send the message using the persistent stream
	if err := cc.Stream.Send(msgRequest); err != nil {
		return fmt.Errorf("failed to send message request: %v", err)
	}

	// Store the message in the sender's local chat history with delivered status set to 0 (false) after successfully sending
	if err := cc.Store.SaveChatMessage(msgRequest.MessageId, cc.AuthClient.ParentClient.CurrentUserID, recipientID, string(plaintext), 0); err != nil {
		cc.Logger.Errorf("Failed to store sent message in chat history: %v", err)
	}

	cc.Logger.Infof("Message sent to recipient %d successfully", recipientID)
	return nil
}

// SendUnencryptedMessage sends a plaintext message to the recipient through the chat service without encryption.
func (cc *ChatClient) SendUnencryptedMessage(ctx context.Context, recipientID uint32, plaintext string) error {

	// Ensure that the persistent stream is open
	if cc.Stream == nil {
		return fmt.Errorf("no active stream found. Ensure that openPersistentStream has been called.")
	}

	messageId := uuid.NewString()
	// Create a new message request with the plaintext content
	msgRequest := &chat.MessageRequest{
		RecipientId:      recipientID,                     // Set recipient ID
		EncryptedMessage: []byte(plaintext),               // Use plaintext directly as the message content
		MessageId:        messageId,                       // Generate a unique message ID
		Timestamp:        time.Now().Format(time.RFC3339), // Timestamp in ISO 8601 format
		EncryptionType:   chat.EncryptionType_PLAIN,       // Set the message type to PLAIN
	}

	// Send the message using the persistent stream
	if err := cc.Stream.Send(msgRequest); err != nil {
		return fmt.Errorf("failed to send message request: %v", err)
	}

	cc.Logger.Infof("Unencrypted message sent to recipient %d successfully", recipientID)

	userId, err := cc.AuthClient.TokenManager.GetUserIdFromAccessToken()
	if err != nil {
		return fmt.Errorf("failed to get user ID from access token: %v", err)
	}
	// Store the message in the sender's local chat history after successfully sending
	err = cc.Store.SaveChatMessage(
		messageId,
		userId,
		recipientID, // Receiver ID
		plaintext,   // The plaintext message
		0,           // delivered status is set to 0 (false) initially
	)
	if err != nil {
		cc.Logger.Errorf("Failed to store sent message in chat history: %v", err)
	}

	cc.Logger.Infof("Sent message stored in chat history with ID %s", messageId)

	return nil
}

// listenForMessages continuously listens for messages on the open stream and handles context cancellations.
func (cc *ChatClient) listenForMessages(ctx context.Context) {
	defer func() {
		cc.Logger.Info("Closing the gRPC stream")
		if err := cc.Stream.CloseSend(); err != nil {
			cc.Logger.Errorf("Failed to close stream: %v", err)
		} else {
			cc.Logger.Info("Stream closed successfully")
		}
	}()

	for {
		select {
		case <-ctx.Done():
			cc.Logger.Info("Stopping message listener due to context cancellation")
			return
		default:
			resp, err := cc.Stream.Recv()
			if err == io.EOF {
				cc.Logger.Info("Stream closed by server")
				return
			}
			if err != nil {
				cc.Logger.Errorf("Failed to receive message: %v", err)
				return
			}

			cc.Logger.Infof("Received message response: %s, with status: %s", resp.EncryptedMessage, resp.Status)

			// Handle different types of responses
			switch resp.Status {
			case "received":
				cc.Logger.Infof("Message %s was received successfully at %s", resp.MessageId, resp.Timestamp)
				// var err error
				// for now
				// unecryptedMessage := string(resp.EncryptedMessage)
				// Decrypt the message
				unecryptedMessage, err := cc.DecryptMessage(ctx, resp)
				if err != nil {
					cc.Logger.Errorf("Failed to decrypt message %s: %v", resp.MessageId, err)
					continue
				}

				err = cc.Store.SaveChatMessage(resp.MessageId, resp.SenderId, resp.RecipientId, unecryptedMessage, 1)
				if err != nil {
					cc.Logger.Errorf("Failed to save message %s in chat history: %v", resp.MessageId, err)
				}

			case "delivered":
				cc.Logger.Infof("Message %s was delivered successfully at %s", resp.MessageId, resp.Timestamp)
				// Update the delivered status in the sender's database.
				err := cc.Store.UpdateMessageDeliveryStatus(resp.MessageId, true)
				if err != nil {
					cc.Logger.Errorf("Failed to update delivery status for message %s: %v", resp.MessageId, err)
				}
				continue
			case "connected":
				cc.Logger.Infof("User Connected at %s", resp.Timestamp)
			case "stored":
				cc.Logger.Infof("Message %s was stored in server buffer for later delivery at %s", resp.MessageId, resp.Timestamp)
				continue
			default:
				cc.Logger.Warnf("Unknown response type: %s", resp.Status)
			}
			// Send the message to the MessageChannel if it exists
			if cc.MessageChannel != nil {
				cc.Logger.Info("Sending received message to message channel")
				cc.MessageChannel <- resp
			} else {
				cc.Logger.Warn("Message channel is not set. Ignoring received message.")
			}
		}
	}
}

func (cc *ChatClient) DecryptMessage(ctx context.Context, resp *chat.MessageResponse) (string, error) {
	remoteAddress := address.Address{
		Name:     fmt.Sprintf("%d", resp.SenderId),
		DeviceID: address.DeviceID(0), // Assuming deviceID 0
	}

	// Reconstruct the Ciphertext object based on messageType
	var ciphertext message.Ciphertext
	var err error

	switch resp.EncryptionType {
	case chat.EncryptionType_SIGNAL:
		ciphertext, err = message.NewSignalFromBytes(resp.EncryptedMessage)
	case chat.EncryptionType_PREKEY:
		ciphertext, err = message.NewPreKeyFromBytes(resp.EncryptedMessage)
	case chat.EncryptionType_PLAIN:
		return string(resp.EncryptedMessage), nil
	default:
		return "", fmt.Errorf("Message has unknown encryption type: %v", resp.EncryptionType)
	}

	if err != nil {
		return "", fmt.Errorf("failed to reconstruct ciphertext: %v", err)
	}
	// Create a session object
	session := &session.Session{
		RemoteAddress:     remoteAddress,
		SessionStore:      cc.Store.SessionStore(),
		PreKeyStore:       cc.Store.PreKeyStore(),
		SignedPreKeyStore: cc.Store.SignedPreKeyStore(),
		IdentityKeyStore:  cc.Store.IdentityStore(),
	}

	// Decrypt the message
	plaintext, err := session.DecryptMessage(ctx, rand.Reader, ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt message: %v", err)
	}

	return string(plaintext), nil
}
