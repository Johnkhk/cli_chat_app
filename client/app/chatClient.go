package app

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/johnkhk/cli_chat_app/genproto/chat"
)

// ChatClient encapsulates the gRPC client for chat services.
type ChatClient struct {
	Client chat.ChatServiceClient // gRPC client for chat service
	Logger *logrus.Logger         // Logger for logging messages and errors
}

// SendMessage handles creating and sending a message via the gRPC stream.
func (c *ChatClient) SendMessage(ctx context.Context, recipientID int32, messageContent string, filePath string) error {
	// Create a new stream for sending messages
	stream, err := c.Client.StreamMessages(ctx)
	if err != nil {
		c.Logger.Errorf("Failed to establish stream: %v", err)
		return err
	}

	// Generate a unique message ID
	messageID := uuid.New().String()
	timestamp := time.Now().Format(time.RFC3339)

	// Encrypt the message content (this should be done using your encryption method)
	encryptedMessage := encryptMessage([]byte(messageContent)) // Implement your encryption method here

	// Prepare the file if provided
	var fileContent []byte
	var fileName, fileType string
	var fileSize int64

	if filePath != "" {
		// Load and encrypt file content
		fileContent, fileName, fileType, fileSize, err = prepareFile(filePath)
		if err != nil {
			c.Logger.Errorf("Failed to prepare file: %v", err)
			return err
		}
	}

	// Create a MessageRequest struct with the required fields
	req := &chat.MessageRequest{
		RecipientId:      recipientID,
		EncryptedMessage: encryptedMessage,
		MessageId:        messageID,
		Timestamp:        timestamp,
		FileContent:      fileContent,
		FileName:         fileName,
		FileType:         fileType,
		FileSize:         fileSize,
	}

	// Send the message request through the stream
	if err := stream.Send(req); err != nil {
		c.Logger.Errorf("Failed to send message: %v", err)
		return err
	}

	c.Logger.Infof("Message sent successfully: %s", messageID)

	// Listen for the server's response
	go func() {
		for {
			resp, err := stream.Recv()
			if err != nil {
				c.Logger.Errorf("Failed to receive response: %v", err)
				return
			}

			// Handle the response from the server
			if resp.Status == "delivered" {
				c.Logger.Infof("Message delivered: %s", resp.MessageId)
			} else if resp.Status == "error" {
				c.Logger.Errorf("Error delivering message: %s, error: %s", resp.MessageId, resp.ErrorMessage)
			}
		}
	}()

	return nil
}

// EncryptMessage is a placeholder function for encrypting the message content.
func encryptMessage(content []byte) []byte {
	// Replace this with your actual encryption logic
	return content // For now, just returning the original content
}

// PrepareFile handles reading the file, encrypting it, and preparing it for sending.
func prepareFile(filePath string) ([]byte, string, string, int64, error) {
	// Implement file loading and encryption logic here
	// Example: Load file from filePath, encrypt the content, and return file details
	return nil, "", "", 0, nil // Placeholder implementation
}
