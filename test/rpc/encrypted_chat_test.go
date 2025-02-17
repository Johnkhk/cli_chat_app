package rpc

import (
	"context"
	"testing"
	"time"

	"github.com/johnkhk/cli_chat_app/client/app"
	"github.com/johnkhk/cli_chat_app/genproto/chat"
	utils "github.com/johnkhk/cli_chat_app/test"
	"github.com/johnkhk/cli_chat_app/test/setup"
)

// Helper function for sending and verifying encrypted messages
func sendAndVerifyMessage(t *testing.T, sender *app.RpcClient, receiver *app.RpcClient, message []byte, expectedType chat.EncryptionType) {
	err := sender.ChatClient.SendMessage(context.Background(), receiver.CurrentUserID, 0, message, nil)
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	select {
	case msg := <-receiver.ChatClient.MessageChannel:
		if msg.EncryptionType != expectedType {
			t.Fatalf("Expected %v message, but received: %v", expectedType, msg.EncryptionType)
		}

		decryptedBytes, err := receiver.ChatClient.DecryptMessage(context.Background(), msg)
		decrypted := string(decryptedBytes)
		if err != nil {
			t.Fatalf("Failed to decrypt message: %v", err)
		}
		if string(decrypted) != string(message) {
			t.Fatalf("Decrypted message does not match original message. Got: %s, Want: %s", decrypted, message)
		}
		t.Logf("Successfully received and decrypted message: %s", decrypted)
	case <-time.After(3 * time.Second):
		t.Fatal("Did not receive message within timeout period")
	}
}

// Test for one user sending PreKey messages only
func TestOneUserSendingMessagesPreKeyOnly(t *testing.T) {
	rpcClients, _, cleanup, _ := setup.InitializeTestResources(t, nil, 2)
	defer cleanup()

	client1 := rpcClients[0]
	client2 := rpcClients[1]

	utils.RegisterAndLoginUser(t, client1, "user1")
	utils.RegisterAndLoginUser(t, client2, "user2")

	utils.WaitForWelcomeMessage(t, client1, "user1")
	utils.WaitForWelcomeMessage(t, client2, "user2")

	messageFromUser1 := []byte("Hello, this is an encrypted message from User1 to User2")
	sendAndVerifyMessage(t, client1, client2, messageFromUser1, chat.EncryptionType_PREKEY)
}

// Test for transitioning from PreKey to Signal messages
func TestPreKeyToSignalMessageTransition(t *testing.T) {
	rpcClients, _, cleanup, _ := setup.InitializeTestResources(t, nil, 2)
	defer cleanup()

	client1 := rpcClients[0]
	client2 := rpcClients[1]

	utils.RegisterAndLoginUser(t, client1, "user1")
	utils.RegisterAndLoginUser(t, client2, "user2")

	utils.WaitForWelcomeMessage(t, client1, "user1")
	utils.WaitForWelcomeMessage(t, client2, "user2")

	// First message from User1 to User2 (PreKey)
	message1 := []byte("Hello, this is an encrypted message from User1 to User2")
	sendAndVerifyMessage(t, client1, client2, message1, chat.EncryptionType_PREKEY)

	// Response from User2 to User1 (PreKey)
	message2 := []byte("Hello, this is a response from User2 to User1")
	sendAndVerifyMessage(t, client2, client1, message2, chat.EncryptionType_SIGNAL)

	// Second message from User1 to User2 (Signal)
	message3 := []byte("Hello again, this should be a Signal message")
	sendAndVerifyMessage(t, client1, client2, message3, chat.EncryptionType_SIGNAL)
}

// Test for persisting encrypted chat history for both sender and receiver
func TestSenderAndReceiverSaveEncryptedChatHistory(t *testing.T) {
	rpcClients, _, cleanup, _ := setup.InitializeTestResources(t, nil, 2)
	defer cleanup()

	client1 := rpcClients[0]
	client2 := rpcClients[1]

	utils.RegisterAndLoginUser(t, client1, "user1")
	utils.RegisterAndLoginUser(t, client2, "user2")

	utils.WaitForWelcomeMessage(t, client1, "user1")
	utils.WaitForWelcomeMessage(t, client2, "user2")

	// User1 sends an encrypted message to User2
	message := []byte("Encrypted message from User1 to User2")
	sendAndVerifyMessage(t, client1, client2, message, chat.EncryptionType_PREKEY)

	time.Sleep(2 * time.Second) // Allow processing time

	// Check chat history for User1
	chatHistory1, err := client1.ChatClient.Store.GetChatHistory(client1.CurrentUserID, client2.CurrentUserID)
	if err != nil || string(chatHistory1[0].Message) != string(message) {
		t.Fatalf("Failed to verify message in User1's chat history. Got: %s, Want: %s", chatHistory1[0].Message, message)
	}

	// Check chat history for User2
	chatHistory2, err := client2.ChatClient.Store.GetChatHistory(client1.CurrentUserID, client2.CurrentUserID)
	if err != nil || string(chatHistory2[0].Message) != string(message) {
		t.Fatalf("Failed to verify message in User2's chat history. Got: %s, Want: %s", chatHistory2[0].Message, message)
	}
	t.Logf("Encrypted message successfully saved in chat histories of both User1 and User2")
}
