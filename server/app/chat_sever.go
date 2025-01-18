package app

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/johnkhk/cli_chat_app/genproto/chat"
)

type ChatServiceServer struct {
	chat.UnimplementedChatServiceServer
	ActiveClients       map[uint32]chat.ChatService_StreamMessagesServer // Map from userID to their active stream
	UndeliveredMessages map[uint32][]*chat.MessageResponse               // Map from userID to their undelivered messages
	mu                  sync.RWMutex                                     // Protect access to ActiveClients and UndeliveredMessages
	Logger              *logrus.Logger
}

func NewChatServiceServer(logger *logrus.Logger) *ChatServiceServer {
	return &ChatServiceServer{
		ActiveClients:       make(map[uint32]chat.ChatService_StreamMessagesServer),
		UndeliveredMessages: make(map[uint32][]*chat.MessageResponse),
		Logger:              logger,
	}
}

// IsActiveClient checks if a user is in the ActiveClients map.
func (s *ChatServiceServer) IsActiveClient(userID uint32) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.ActiveClients[userID]
	return exists
}

// StreamMessages handles bidirectional message streaming between users.
func (s *ChatServiceServer) StreamMessages(stream chat.ChatService_StreamMessagesServer) error {
	s.Logger.Info("StreamMessages called")

	// Extract sender's userID from the stream context.
	ctx := stream.Context()
	senderID, senderUsername, err := s.extractSenderIDFromContext(ctx)
	if err != nil {
		s.Logger.Errorf("Failed to extract sender ID: %v", err)
		return err
	}
	s.Logger.Infof("User %d connected with stream", senderID)

	// Register the sender's stream in the active clients map when the stream is established.
	s.registerClient(senderID, stream)
	defer s.unregisterClient(senderID)

	// Send a welcome message after the stream is established.
	welcomeResponse := &chat.MessageResponse{
		SenderId:         0, // Use 0 or a special ID to indicate server-originated message.
		SenderUsername:   "server",
		RecipientId:      senderID,
		MessageId:        "welcome",
		Status:           "connected",
		Timestamp:        time.Now().Format(time.RFC3339),
		EncryptionType:   chat.EncryptionType_PLAIN,
		EncryptedMessage: []byte("Welcome to the CLI chat app!"),
	}
	s.Logger.Infof("Sending server message!!!: %s", welcomeResponse.EncryptedMessage)

	// Send the welcome message to the client.
	if err := stream.Send(welcomeResponse); err != nil {
		s.Logger.Errorf("Failed to send welcome message to user %d: %v", senderID, err)
		return err
	}
	s.Logger.Infof("Sent welcome message to user %d", senderID)

	// Deliver any undelivered messages to the client.
	if err := s.deliverUndeliveredMessages(senderID); err != nil {
		s.Logger.Errorf("Failed to deliver undelivered messages to user %d: %v", senderID, err)
	}

	for {
		select {
		case <-ctx.Done(): // Handle client disconnection more explicitly.
			s.Logger.Infof("Client %d context closed: %v", senderID, ctx.Err())
			return nil
		default:
			req, err := stream.Recv() // Blocking call to receive a message.
			if err == io.EOF {
				s.Logger.Infof("Stream closed by client %d", senderID)
				return nil // Client closed the stream.
			}
			if err != nil {
				s.Logger.Errorf("Error receiving message from user %d: %v", senderID, err)
				if err == context.Canceled || err == context.DeadlineExceeded {
					s.Logger.Infof("Client %d disconnected: %v", senderID, err)
					return nil
				}
				return err // Handle other errors.
			}

			s.Logger.Infof("Received message with ID %s from user %d to recipient %d", req.MessageId, senderID, req.RecipientId)

			// Send the message directly to the recipient's stream if they are connected.
			if err := s.sendMessageToRecipient(&chat.MessageResponse{
				SenderId:         senderID,
				SenderUsername:   senderUsername, // Include sender's username
				RecipientId:      req.RecipientId,
				MessageId:        req.MessageId,
				EncryptedMessage: req.EncryptedMessage, // Include the actual message content for the recipient
				Status:           "received",
				Timestamp:        time.Now().Format(time.RFC3339), // Timestamp for when the recipient received it
				EncryptionType:   req.EncryptionType,
			}); err != nil {
				s.Logger.Errorf("Failed to send/store message ID %s to recipient %d: %v", req.MessageId, req.RecipientId, err)

				// Send a response back to the sender indicating a failed delivery.
				failedDeliveryResponse := &chat.MessageResponse{
					SenderId:       senderID,
					SenderUsername: "server", // Indicate that this is a server response
					RecipientId:    req.RecipientId,
					MessageId:      req.MessageId,
					Status:         "delivery_failed",
					Timestamp:      time.Now().Format(time.RFC3339),
					EncryptionType: req.EncryptionType,
				}
				if sendErr := stream.Send(failedDeliveryResponse); sendErr != nil {
					s.Logger.Errorf("Failed to send delivery failure response to sender %d: %v", senderID, sendErr)
					return sendErr
				}
			} else {
				var status string
				if s.IsActiveClient(req.RecipientId) {
					status = "delivered"
				} else {
					status = "stored"
				}
				s.Logger.Infof("Message ID %s successfully processed for recipient %d with status %s", req.MessageId, req.RecipientId, status)

				// Send a response back to the sender indicating message status.
				deliveryResponse := &chat.MessageResponse{
					SenderId:         senderID,
					SenderUsername:   senderUsername, // Include sender's username in confirmation response
					RecipientId:      req.RecipientId,
					MessageId:        req.MessageId,
					EncryptedMessage: req.EncryptedMessage, // Include the actual message content for confirmation.
					Status:           status,
					Timestamp:        time.Now().Format(time.RFC3339),
					EncryptionType:   req.EncryptionType,
				}
				if err := stream.Send(deliveryResponse); err != nil {
					s.Logger.Errorf("Failed to send delivery confirmation to sender %d: %v", senderID, err)
					return err
				}
				s.Logger.Infof("Sent delivery confirmation for message ID %s with status %s", req.MessageId, status)
			}
		}
	}
}

// Extracts the sender ID from the context (assuming userID is set in context)
func (s *ChatServiceServer) extractSenderIDFromContext(ctx context.Context) (uint32, string, error) {
	// Assuming sender ID is stored as a string in context
	senderIDStr, ok := ctx.Value("userID").(string)
	if !ok {
		return 0, "", fmt.Errorf("failed to get userID from context")
	}

	senderUsername, ok := ctx.Value("username").(string)
	if !ok {
		return 0, "", fmt.Errorf("failed to get username from context")
	}

	// Convert sender ID from string to uint32
	var senderID uint32
	_, err := fmt.Sscanf(senderIDStr, "%d", &senderID)
	if err != nil {
		return 0, "", fmt.Errorf("invalid sender ID format: %v", err)
	}

	return senderID, senderUsername, nil
}

// registerClient registers a client's stream with their user ID
func (s *ChatServiceServer) registerClient(userID uint32, stream chat.ChatService_StreamMessagesServer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ActiveClients[userID] = stream
	s.Logger.Infof("User %d has been registered in active clients", userID)
}

// unregisterClient removes a client's stream from the active clients map
func (s *ChatServiceServer) unregisterClient(userID uint32) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.ActiveClients, userID)
	s.Logger.Infof("User %d has been unregistered from active clients", userID)
}

// sendMessageToRecipient attempts to send a message directly to the recipient if they are connected.
// If the recipient is not connected, it stores the message in the undelivered messages buffer.
func (s *ChatServiceServer) sendMessageToRecipient(resp *chat.MessageResponse) error {
	s.mu.RLock()
	recipientStream, recipientConnected := s.ActiveClients[resp.RecipientId]
	s.mu.RUnlock()

	if !recipientConnected {
		s.Logger.Warnf("Recipient %d is not connected, storing message ID %s in buffer", resp.RecipientId, resp.MessageId)
		// Store the message in the undelivered messages buffer
		s.mu.Lock()
		s.UndeliveredMessages[resp.RecipientId] = append(s.UndeliveredMessages[resp.RecipientId], resp)
		s.mu.Unlock()
		return nil
	}

	// Forward the response to the recipient's stream
	s.Logger.Infof("Forwarding message ID %s to recipient %d", resp.MessageId, resp.RecipientId)

	if err := recipientStream.Send(resp); err != nil {
		return fmt.Errorf("failed to send message ID %s to recipient %d: %v", resp.MessageId, resp.RecipientId, err)
	}

	s.Logger.Infof("Message ID %s successfully delivered to recipient %d", resp.MessageId, resp.RecipientId)
	return nil
}

// deliverUndeliveredMessages sends any stored messages to the user upon reconnection.
func (s *ChatServiceServer) deliverUndeliveredMessages(userID uint32) error {
	var messages []*chat.MessageResponse

	// Lock the UndeliveredMessages map and get the messages
	s.mu.Lock()
	if msgs, exists := s.UndeliveredMessages[userID]; exists && len(msgs) > 0 {
		messages = msgs
		// Remove messages from buffer
		delete(s.UndeliveredMessages, userID)
	}
	s.mu.Unlock()

	if len(messages) == 0 {
		s.Logger.Infof("No undelivered messages for user %d", userID)
		return nil
	}

	s.Logger.Infof("Delivering %d undelivered messages to user %d", len(messages), userID)

	s.mu.RLock()
	recipientStream, recipientConnected := s.ActiveClients[userID]
	s.mu.RUnlock()
	if !recipientConnected {
		s.Logger.Warnf("Recipient %d is not connected while trying to deliver undelivered messages", userID)
		// If recipient is not connected, put the messages back into the buffer
		s.mu.Lock()
		s.UndeliveredMessages[userID] = messages // Put messages back
		s.mu.Unlock()
		return fmt.Errorf("recipient %d is not connected", userID)
	}

	for _, msg := range messages {
		if err := recipientStream.Send(msg); err != nil {
			s.Logger.Errorf("Failed to send undelivered message ID %s to user %d: %v", msg.MessageId, userID, err)
			// Optionally, handle error
		} else {
			s.Logger.Infof("Delivered undelivered message ID %s to user %d", msg.MessageId, userID)
		}
	}

	return nil
}
