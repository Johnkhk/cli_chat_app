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

// // StreamMessages handles bidirectional message streaming between users
// // Called by the client to establish a bidirectional stream for sending and receiving messages
// func (s *ChatServiceServer) StreamMessages(stream chat.ChatService_StreamMessagesServer) error {

// 	// Extract sender's userID from the stream context
// 	ctx := stream.Context()
// 	senderID, err := s.extractSenderIDFromContext(ctx)
// 	if err != nil {
// 		s.Logger.Errorf("Failed to extract sender ID: %v", err)
// 		return err
// 	}

// 	// Register the sender's stream in the active clients map
// 	// s.registerClient(senderID, stream)
// 	// defer s.unregisterClient(senderID)

// 	s.Logger.Infof("User %d connected with stream", senderID)

// 	for {
// 		// Receive a message from the client stream
// 		req, err := stream.Recv()
// 		if err == io.EOF {
// 			return nil // End of stream
// 		}
// 		if err != nil {
// 			s.Logger.Errorf("Error receiving message from user %d: %v", senderID, err)
// 			return err
// 		}

// 		// Log received message details
// 		s.Logger.Infof("Received message with ID %s from user %d to recipient %d", req.MessageId, senderID, req.RecipientId)

// 		// Special handling for registration messages
// 		if req.RecipientId == math.MaxUint32 { // Assuming a `RecipientId` of `0` indicates a registration message
// 			s.Logger.Infof("Handling registration message from user %d", senderID)

// 			// Mark user as active or perform additional registration logic here
// 			// For example, update user status in a database or broadcast user presence.

// 			// Send a response back indicating successful registration
// 			registrationResponse := &chat.MessageResponse{
// 				MessageId: req.MessageId,
// 				Status:    "registered",
// 				Timestamp: time.Now().Format(time.RFC3339),
// 			}
// 			if err := stream.Send(registrationResponse); err != nil {
// 				s.Logger.Errorf("Failed to send registration response to sender %d: %v", senderID, err)
// 				return err
// 			}

// 			s.Logger.Infof("Registration response sent for user %d", senderID)
// 			continue
// 		}

// 		// Directly send the message to the recipient's stream if they are connected
// 		if err := s.sendMessageToRecipient(req); err != nil {
// 			s.Logger.Errorf("Failed to send message ID %s to recipient %d: %v", req.MessageId, req.RecipientId, err)
// 			// Optionally, send an error response back to the sender
// 			response := &chat.MessageResponse{
// 				MessageId:    req.MessageId,
// 				Status:       "error",
// 				Timestamp:    time.Now().Format(time.RFC3339),
// 				ErrorMessage: fmt.Sprintf("Failed to deliver message to recipient %d: %v", req.RecipientId, err),
// 			}
// 			if err := stream.Send(response); err != nil {
// 				s.Logger.Errorf("Failed to send error response to sender %d: %v", senderID, err)
// 				return err
// 			}
// 		} else {
// 			// Send a response back to the sender confirming message delivery
// 			response := &chat.MessageResponse{
// 				MessageId: req.MessageId,
// 				Status:    "delivered",                     // Indicate that the message was delivered to the recipient
// 				Timestamp: time.Now().Format(time.RFC3339), // Use ISO 8601 format
// 			}

// 			// Send the response back to the sender
// 			if err := stream.Send(response); err != nil {
// 				s.Logger.Errorf("Failed to send delivery confirmation to sender %d: %v", senderID, err)
// 				return err
// 			}

// 			s.Logger.Infof("Sent delivery confirmation for message ID %s", req.MessageId)
// 		}
// 	}
// }

// StreamMessages handles bidirectional message streaming between users
func (s *ChatServiceServer) StreamMessages(stream chat.ChatService_StreamMessagesServer) error {
	s.Logger.Info("StreamMessages called")
	// Extract sender's userID from the stream context
	ctx := stream.Context()
	senderID, err := s.extractSenderIDFromContext(ctx)
	if err != nil {
		s.Logger.Errorf("Failed to extract sender ID: %v", err)
		return err
	}
	s.Logger.Infof("User %d connected with stream", senderID)

	// Register the sender's stream in the active clients map when the stream is established
	s.registerClient(senderID, stream)
	defer s.unregisterClient(senderID)

	s.Logger.Infof("User %d connected with stream", senderID)

	// Perform any additional initial registration logic here
	// For example, mark the user as active or send a welcome message.
	// Example: update user status in a database or notify other users
	s.Logger.Infof("User %d is now active and registered with the chat service", senderID)

	// // Send a welcome message or registration response to the user (if needed)
	// welcomeResponse := &chat.MessageResponse{
	// 	MessageId: fmt.Sprintf("welcome-%d", senderID),
	// 	Status:    "registered",
	// 	Timestamp: time.Now().Format(time.RFC3339),
	// }
	// if err := stream.Send(welcomeResponse); err != nil {
	// 	s.Logger.Errorf("Failed to send welcome message to user %d: %v", senderID, err)
	// 	return err
	// }
	s.Logger.Infof("Welcome message sent to user %d", senderID)

	for {
		// Receive a message from the client stream
		req, err := stream.Recv()
		if err == io.EOF {
			return nil // End of stream
		}
		if err != nil {
			s.Logger.Errorf("Error receiving message from user %d: %v", senderID, err)
			return err
		}

		// Log received message details
		s.Logger.Infof("Received message with ID %s from user %d to recipient %d", req.MessageId, senderID, req.RecipientId)

		// Directly send the message to the recipient's stream if they are connected
		if err := s.sendMessageToRecipient(req); err != nil {
			s.Logger.Errorf("Failed to send message ID %s to recipient %d: %v", req.MessageId, req.RecipientId, err)
			// Optionally, send an error response back to the sender
			response := &chat.MessageResponse{
				MessageId:    req.MessageId,
				Status:       "error",
				Timestamp:    time.Now().Format(time.RFC3339),
				ErrorMessage: fmt.Sprintf("Failed to deliver message to recipient %d: %v", req.RecipientId, err),
			}
			if err := stream.Send(response); err != nil {
				s.Logger.Errorf("Failed to send error response to sender %d: %v", senderID, err)
				return err
			}
		} else {
			// Send a response back to the sender confirming message delivery
			response := &chat.MessageResponse{
				MessageId: req.MessageId,
				Status:    "delivered",                     // Indicate that the message was delivered to the recipient
				Timestamp: time.Now().Format(time.RFC3339), // Use ISO 8601 format
			}

			// Send the response back to the sender
			if err := stream.Send(response); err != nil {
				s.Logger.Errorf("Failed to send delivery confirmation to sender %d: %v", senderID, err)
				return err
			}

			s.Logger.Infof("Sent delivery confirmation for message ID %s", req.MessageId)
		}
	}
}

// Extracts the sender ID from the context (assuming userID is set in context)
func (s *ChatServiceServer) extractSenderIDFromContext(ctx context.Context) (uint32, error) {
	// Assuming sender ID is stored as a string in context
	senderIDStr, ok := ctx.Value("userID").(string)
	if !ok {
		return 0, fmt.Errorf("failed to get user ID from context")
	}

	// Convert sender ID from string to uint32
	var senderID uint32
	_, err := fmt.Sscanf(senderIDStr, "%d", &senderID)
	if err != nil {
		return 0, fmt.Errorf("invalid sender ID format: %v", err)
	}

	return senderID, nil
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

// sendMessageToRecipient attempts to send a message directly to the recipient if they are connected
func (s *ChatServiceServer) sendMessageToRecipient(msg *chat.MessageRequest) error {
	s.mu.RLock()
	recipientStream, recipientConnected := s.ActiveClients[msg.RecipientId]
	s.mu.RUnlock()

	if !recipientConnected {
		s.Logger.Warnf("Recipient %d is not connected", msg.RecipientId)
		return fmt.Errorf("recipient %d is not connected", msg.RecipientId)
	}

	// Forward the message to the recipient's stream
	s.Logger.Infof("Forwarding message ID %s to recipient %d", msg.MessageId, msg.RecipientId)
	response := &chat.MessageResponse{
		MessageId: msg.MessageId,
		Status:    "received",
		Timestamp: time.Now().Format(time.RFC3339), // Timestamp for when the recipient received it
	}

	if err := recipientStream.Send(response); err != nil {
		return fmt.Errorf("failed to send message ID %s to recipient %d: %v", msg.MessageId, msg.RecipientId, err)
	}

	s.Logger.Infof("Message ID %s successfully delivered to recipient %d", msg.MessageId, msg.RecipientId)
	return nil
}
