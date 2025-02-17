package store

import (
	"fmt"
	"time"

	"github.com/johnkhk/cli_chat_app/client/lib"
)

type ChatMessage struct {
	MessageID  string    `json:"messageId"`
	SenderID   uint32    `json:"senderId"`
	ReceiverID uint32    `json:"receiverId"`
	Message    string    `json:"message"`
	FileType   string    `json:"fileType"`
	FileSize   uint64    `json:"fileSize"`
	FileName   string    `json:"fileName"`
	Timestamp  time.Time `json:"timestamp"`
	Delivered  int       `json:"delivered"`
}

// SaveChatMessage inserts a new chat message with the specified messageId into the `chat_history` table.
func (s *SQLiteStore) SaveChatMessage(messageID string, senderID, receiverID uint32, message []byte, delivered int, fileOpts *lib.SendMessageOptions) error {
	// Prepare the SQL query for inserting a new chat message.
	if fileOpts == nil {
		fileOpts = &lib.SendMessageOptions{
			FileType: "text",
			FileSize: 0,
			FileName: "",
		}
	}

	query := `
		INSERT INTO chat_history (messageId, sender_id, receiver_id, message, delivered, file_type, file_size, file_name)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?);` // `delivered` is set to 0 (false) initially.
	_, err := s.DB.Exec(query, messageID, senderID, receiverID, message, delivered, fileOpts.FileType, fileOpts.FileSize, fileOpts.FileName)
	if err != nil {
		return fmt.Errorf("failed to save chat message: %v", err)
	}

	return nil
}

// UpdateMessageDeliveryStatus updates the `delivered` status of a message in the `chat_history` table.
func (s *SQLiteStore) UpdateMessageDeliveryStatus(messageID string, delivered bool) error {
	query := `
		UPDATE chat_history
		SET delivered = ?
		WHERE messageId = ?;`
	_, err := s.DB.Exec(query, delivered, messageID)
	if err != nil {
		return fmt.Errorf("failed to update delivery status for message ID %s: %v", messageID, err)
	}
	return nil
}

// GetChatHistory retrieves all chat messages between a sender and receiver.
func (s *SQLiteStore) GetChatHistory(senderID, receiverID uint32) ([]ChatMessage, error) {
	query := `
		SELECT messageId, sender_id, receiver_id, message, file_type, file_size, file_name, timestamp, delivered
		FROM chat_history
		WHERE (sender_id = ? AND receiver_id = ?)
		   OR (sender_id = ? AND receiver_id = ?)
		ORDER BY timestamp ASC;`

	rows, err := s.DB.Query(query, senderID, receiverID, receiverID, senderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat history: %v", err)
	}
	defer rows.Close()

	var messages []ChatMessage
	for rows.Next() {
		var msg ChatMessage
		err := rows.Scan(&msg.MessageID, &msg.SenderID, &msg.ReceiverID, &msg.Message, &msg.FileType, &msg.FileSize, &msg.FileName, &msg.Timestamp, &msg.Delivered)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

// DeleteMessage deletes a specific message from the `chat_history` table.
func (s *SQLiteStore) DeleteMessage(messageID int64) error {
	query := `
		DELETE FROM chat_history
		WHERE id = ?;`
	_, err := s.DB.Exec(query, messageID)
	if err != nil {
		return fmt.Errorf("failed to delete message with ID %d: %v", messageID, err)
	}
	return nil
}
