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
	ActiveClients map[uint32]chat.ChatService_StreamMessagesServer // Map from userID to their active stream
	mu            sync.RWMutex                                     // Protect access to ActiveClients
	Logger        *logrus.Logger
}

func NewChatServiceServer(logger *logrus.Logger) *ChatServiceServer {
	return &ChatServiceServer{
		ActiveClients: make(map[uint32]chat.ChatService_StreamMessagesServer),
		Logger:        logger,
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
		SenderId:       0, // Use 0 or a special ID to indicate server-originated message.
		SenderUsername: "server",
		RecipientId:    senderID,
		MessageId:      "welcome",
		Status:         "connected",
		Timestamp:      time.Now().Format(time.RFC3339),
	}

	// Send the welcome message to the client.
	if err := stream.Send(welcomeResponse); err != nil {
		s.Logger.Errorf("Failed to send welcome message to user %d: %v", senderID, err)
		return err
	}
	s.Logger.Infof("Sent welcome message to user %d", senderID)

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
			}); err != nil {
				s.Logger.Errorf("Failed to send message ID %s to recipient %d: %v", req.MessageId, req.RecipientId, err)

				// Send a response back to the sender indicating a failed delivery.
				failedDeliveryResponse := &chat.MessageResponse{
					SenderId:       senderID,
					SenderUsername: "server", // Indicate that this is a server response
					RecipientId:    req.RecipientId,
					MessageId:      req.MessageId,
					Status:         "delivery_failed",
					Timestamp:      time.Now().Format(time.RFC3339),
				}
				if sendErr := stream.Send(failedDeliveryResponse); sendErr != nil {
					s.Logger.Errorf("Failed to send delivery failure response to sender %d: %v", senderID, sendErr)
					return sendErr
				}
			} else {
				s.Logger.Infof("Message ID %s successfully sent to recipient %d", req.MessageId, req.RecipientId)

				// Send a response back to the sender confirming message delivery.
				deliveryResponse := &chat.MessageResponse{
					SenderId:         senderID,
					SenderUsername:   senderUsername, // Include sender's username in confirmation response
					RecipientId:      req.RecipientId,
					MessageId:        req.MessageId,
					EncryptedMessage: req.EncryptedMessage, // Include the actual message content for confirmation.
					Status:           "delivered",
					Timestamp:        time.Now().Format(time.RFC3339),
				}
				if err := stream.Send(deliveryResponse); err != nil {
					s.Logger.Errorf("Failed to send delivery confirmation to sender %d: %v", senderID, err)
					return err
				}
				s.Logger.Infof("Sent delivery confirmation for message ID %s", req.MessageId)
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
// Accepts a MessageResponse instead of MessageRequest.
func (s *ChatServiceServer) sendMessageToRecipient(resp *chat.MessageResponse) error {
	s.mu.RLock()
	recipientStream, recipientConnected := s.ActiveClients[resp.RecipientId]
	s.mu.RUnlock()

	if !recipientConnected {
		s.Logger.Warnf("Recipient %d is not connected", resp.RecipientId)
		return fmt.Errorf("recipient %d is not connected", resp.RecipientId)
	}

	// Forward the response to the recipient's stream
	s.Logger.Infof("Forwarding message ID %s to recipient %d", resp.MessageId, resp.RecipientId)

	if err := recipientStream.Send(resp); err != nil {
		return fmt.Errorf("failed to send message ID %s to recipient %d: %v", resp.MessageId, resp.RecipientId, err)
	}

	s.Logger.Infof("Message ID %s successfully delivered to recipient %d", resp.MessageId, resp.RecipientId)
	return nil
}
