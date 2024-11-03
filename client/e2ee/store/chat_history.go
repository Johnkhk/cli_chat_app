package store

import (
	"fmt"
	"time"
)

type ChatMessage struct {
	MessageID  string    `json:"messageId"`
	SenderID   uint32    `json:"senderId"`
	ReceiverID uint32    `json:"receiverId"`
	Message    string    `json:"message"`
	Timestamp  time.Time `json:"timestamp"`
	Delivered  int       `json:"delivered"`
}

// SaveChatMessage inserts a new chat message with the specified messageId into the `chat_history` table.
func (s *SQLiteStore) SaveChatMessage(messageID string, senderID, receiverID uint32, message string, delivered int) error {
	// Prepare the SQL query for inserting a new chat message.
	query := `
		INSERT INTO chat_history (messageId, sender_id, receiver_id, message, delivered)
		VALUES (?, ?, ?, ?, ?);` // `delivered` is set to 0 (false) initially.
	_, err := s.DB.Exec(query, messageID, senderID, receiverID, message, delivered)
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
		SELECT messageId, sender_id, receiver_id, message, timestamp, delivered
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
		err := rows.Scan(&msg.MessageID, &msg.SenderID, &msg.ReceiverID, &msg.Message, &msg.Timestamp, &msg.Delivered)
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

// GetAllChatHistory retrieves all stored chat messages.
func (s *SQLiteStore) GetAllChatHistory() ([]ChatMessage, error) {
	query := `
		SELECT messageId, sender_id, receiver_id, message, timestamp, delivered
		FROM chat_history
		ORDER BY timestamp ASC;`

	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all chat history: %v", err)
	}
	defer rows.Close()

	var messages []ChatMessage
	for rows.Next() {
		var msg ChatMessage
		err := rows.Scan(&msg.MessageID, &msg.SenderID, &msg.ReceiverID, &msg.Message, &msg.Timestamp, &msg.Delivered)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func (s *SQLiteStore) GetChatHistoryBetweenUsers(senderID, receiverID uint32) ([]ChatMessage, error) {
	query := `
		SELECT messageId, sender_id, receiver_id, message, timestamp, delivered
		FROM chat_history
		WHERE (sender_id = ? AND receiver_id = ?)
		   OR (sender_id = ? AND receiver_id = ?)
		ORDER BY timestamp ASC;`

	rows, err := s.DB.Query(query, senderID, receiverID, receiverID, senderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat history between users: %v", err)
	}
	defer rows.Close()

	var messages []ChatMessage
	for rows.Next() {
		var msg ChatMessage
		err := rows.Scan(&msg.MessageID, &msg.SenderID, &msg.ReceiverID, &msg.Message, &msg.Timestamp, &msg.Delivered)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
