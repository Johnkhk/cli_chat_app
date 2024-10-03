package rpc

import (
	"context"
	"testing"
	"time"

	"github.com/johnkhk/cli_chat_app/test/setup"
)

// TestRegisterLoginAndStreamFlow tests the full flow of user registration, login, and persistent stream.
func TestRegisterLoginAndStreamFlow(t *testing.T) {
	// Initialize resources using default configuration
	rpcClients, _, cleanup, server := setup.InitializeTestResources(t, nil, 1)
	rpcClient := rpcClients[0]
	defer cleanup()
	log := rpcClient.Logger

	// Register the user
	log.Infof("Registering user")
	err := rpcClient.AuthClient.RegisterUser("unregistered", "testpassword")
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	// Login of the registered user
	log.Infof("Testing registered user login")
	err, _ = rpcClient.AuthClient.LoginUser("unregistered", "testpassword")
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	// After login, we should wait a moment to ensure that the persistent stream is established.
	// Add a small delay to allow the stream to register the user.
	time.Sleep(2 * time.Second) // Adjust as needed, depending on gRPC setup

	// Check that the user is registered in ActiveClients on the chat server
	activeClients := server.ChatServer.ActiveClients
	log.Infof("Active clients count: %d", len(activeClients))
	for k, v := range activeClients {
		log.Infof("Active Client Key: %d, Value: %v", k, v)
	}

	// Get the user ID (we assume the registered user has ID = 1 for this test)
	userID := uint32(1) // Replace this with the actual user ID based on your logic

	// Check if the user is in the ActiveClients map
	if _, exists := activeClients[userID]; !exists {
		t.Fatalf("User %d is not found in ActiveClients map after opening stream", userID)
	}

	log.Infof("User %d successfully registered in ActiveClients map", userID)
}

func TestLoginWelcomeMessageReceived(t *testing.T) {
	// Initialize resources using default configuration
	rpcClients, _, cleanup, _ := setup.InitializeTestResources(t, nil, 1)
	rpcClient := rpcClients[0]
	defer cleanup()
	log := rpcClient.Logger

	// Register a new user
	log.Infof("Registering new user")
	err := rpcClient.AuthClient.RegisterUser("newuser", "testpassword")
	if err != nil {
		t.Fatalf("Failed to register new user: %v", err)
	}

	// Login the user
	log.Infof("Logging in user")
	err, _ = rpcClient.AuthClient.LoginUser("newuser", "testpassword")
	if err != nil {
		t.Fatalf("Failed to login user: %v", err)
	}

	// Wait for the welcome message to be received
	select {
	case msg := <-rpcClient.ChatClient.MessageChannel:
		expectedMessage := "welcome"
		if msg.MessageId != expectedMessage {
			t.Errorf("Expected welcome message, but got: %s", msg.MessageId)
		} else {
			log.Infof("Received welcome MessageID: %s", msg.MessageId)
		}
	case <-time.After(3 * time.Second): // Adjust timeout as necessary
		t.Fatalf("Did not receive welcome message within timeout period")
	}
}

func TestTwoUsersSendUnencrypteedMessages(t *testing.T) {
	rpcClients, _, cleanup, _ := setup.InitializeTestResources(t, nil, 2)
	defer cleanup()

	client1 := rpcClients[0] // Represents User1
	client2 := rpcClients[1] // Represents User2

	// Register and login two users
	if err := client1.AuthClient.RegisterUser("user1", "password"); err != nil {
		t.Fatalf("Failed to register user1: %v", err)
	}
	if err, _ := client1.AuthClient.LoginUser("user1", "password"); err != nil {
		t.Fatalf("Failed to login user1: %v", err)
	}
	if err := client2.AuthClient.RegisterUser("user2", "password"); err != nil {
		t.Fatalf("Failed to register user2: %v", err)
	}
	if err, _ := client2.AuthClient.LoginUser("user2", "password"); err != nil {
		t.Fatalf("Failed to login user2: %v", err)
	}

	// user1ID, err := client1.AuthClient.TokenManager.GetUserIdFromAccessToken()
	// if err != nil {
	// 	t.Fatalf("Failed to get user1 ID: %v", err)
	// }

	user2ID, err := client2.AuthClient.TokenManager.GetUserIdFromAccessToken()
	if err != nil {
		t.Fatalf("Failed to get user2 ID: %v", err)
	}

	// User1 sends a message to User2
	t.Log("User 1 sending unencrypted message to User 2")
	messageFromUser1 := "Hello from User 1 to User 2"
	if err := client1.ChatClient.SendUnencryptedMessage(context.Background(), user2ID, messageFromUser1); err != nil {
		t.Fatalf("Failed to send message from User 1 to User 2: %v", err)
	}

	// Consume the first message (e.g., the welcome message)
	select {
	case msg := <-client2.ChatClient.MessageChannel:
		t.Logf("User 2 received first message (MessageID: %s): %s", msg.MessageId, msg.EncryptedMessage)
	case <-time.After(3 * time.Second):
		t.Fatalf("User 2 did not receive the first message in time")
	}

	// Consume the second message (e.g., the actual chat message)
	select {
	case msg := <-client2.ChatClient.MessageChannel:
		t.Logf("User 2 received second message (MessageID: %s): %s", msg.MessageId, msg.EncryptedMessage)
	case <-time.After(3 * time.Second):
		t.Fatalf("User 2 did not receive the second message in time")
	}

}

func TestSenderAndReceiverSavesChatHistory(t *testing.T) {
	rpcClients, _, cleanup, _ := setup.InitializeTestResources(t, nil, 2)
	defer cleanup()

	client1 := rpcClients[0] // Represents User1
	client2 := rpcClients[1] // Represents User2

	// Register and login two users
	if err := client1.AuthClient.RegisterUser("user1", "password"); err != nil {
		t.Fatalf("Failed to register user1: %v", err)
	}
	if err, _ := client1.AuthClient.LoginUser("user1", "password"); err != nil {
		t.Fatalf("Failed to login user1: %v", err)
	}
	if err := client2.AuthClient.RegisterUser("user2", "password"); err != nil {
		t.Fatalf("Failed to register user2: %v", err)
	}
	if err, _ := client2.AuthClient.LoginUser("user2", "password"); err != nil {
		t.Fatalf("Failed to login user2: %v", err)
	}

	user1ID, err := client1.AuthClient.TokenManager.GetUserIdFromAccessToken()
	if err != nil {
		t.Fatalf("Failed to get user1 ID: %v", err)
	}

	user2ID, err := client2.AuthClient.TokenManager.GetUserIdFromAccessToken()
	if err != nil {
		t.Fatalf("Failed to get user2 ID: %v", err)
	}

	// User1 sends a message to User2
	t.Log("User 1 sending unencrypted message to User 2")
	messageFromUser1 := "custom message from User 1 to User 2"
	if err := client1.ChatClient.SendUnencryptedMessage(context.Background(), user2ID, messageFromUser1); err != nil {
		t.Fatalf("Failed to send message from User 1 to User 2: %v", err)
	}

	time.Sleep(2 * time.Second) // Sleep to allow the message to be processed (before context cancel from test cleanup)

	// Make sure the message is saved in the chat history of User1
	chatMessages, err := client1.ChatClient.Store.GetChatHistory(user1ID, user2ID)

	message := chatMessages[0]
	if message.Message != messageFromUser1 {
		t.Fatalf("Expected message %s, but got: %s", messageFromUser1, message.Message)
	}

	// Delievered status should be 1 since user 1 should have received the message
	if message.Delivered != 1 {
		t.Fatalf("Expected delivered status for sender to be 1, but got: %d", message.Delivered)
	}

	// Make sure the message is saved in the chat history of User2
	chatMessages, err = client2.ChatClient.Store.GetChatHistory(user1ID, user2ID)

	message = chatMessages[0]
	if message.Message != messageFromUser1 {
		t.Fatalf("Expected message %s, but got: %s", messageFromUser1, message.Message)
	}

	// Delivered status should be 1 since user 1 should have received the message
	if message.Delivered != 1 {
		t.Fatalf("Expected delivered status 1, but got: %d", message.Delivered)
	}
}

func TestLogoutUser(t *testing.T) {
	rpcClients, _, cleanup, server := setup.InitializeTestResources(t, nil, 2)
	defer cleanup()

	client1 := rpcClients[0] // Represents User1
	client2 := rpcClients[1] // Represents User2

	// Register and login two users
	if err := client1.AuthClient.RegisterUser("user1", "password"); err != nil {
		t.Fatalf("Failed to register user1: %v", err)
	}
	if err, _ := client1.AuthClient.LoginUser("user1", "password"); err != nil {
		t.Fatalf("Failed to login user1: %v", err)
	}
	if err := client2.AuthClient.RegisterUser("user2", "password"); err != nil {
		t.Fatalf("Failed to register user2: %v", err)
	}
	if err, _ := client2.AuthClient.LoginUser("user2", "password"); err != nil {
		t.Fatalf("Failed to login user2: %v", err)
	}

	// sleep some time before we logout user2
	time.Sleep(2 * time.Second)

	// Make sure active clients map is updated
	activeClients := server.ChatServer.ActiveClients
	if len(activeClients) != 2 {
		t.Fatalf("Expected 2 active client, but got: %d", len(activeClients))
	}

	// User2 logs out
	if err := client2.AuthClient.LogoutUser(); err != nil {
		t.Fatalf("Failed to logout user2: %v", err)
	}
	time.Sleep(2 * time.Second) // Sleep to allow logout to cancel context, stream

	// Make sure active clients map is updated
	activeClients = server.ChatServer.ActiveClients
	if len(activeClients) != 1 {
		t.Fatalf("Expected 1 active client, but got: %d", len(activeClients))
	}
}

func TestReceiverIsOfflineWhenSenderSends(t *testing.T) {
	rpcClients, _, cleanup, server := setup.InitializeTestResources(t, nil, 2)
	defer cleanup()

	client1 := rpcClients[0] // Represents User1
	client2 := rpcClients[1] // Represents User2

	// Register and login two users
	if err := client1.AuthClient.RegisterUser("user1", "password"); err != nil {
		t.Fatalf("Failed to register user1: %v", err)
	}
	if err, _ := client1.AuthClient.LoginUser("user1", "password"); err != nil {
		t.Fatalf("Failed to login user1: %v", err)
	}
	if err := client2.AuthClient.RegisterUser("user2", "password"); err != nil {
		t.Fatalf("Failed to register user2: %v", err)
	}
	if err, _ := client2.AuthClient.LoginUser("user2", "password"); err != nil {
		t.Fatalf("Failed to login user2: %v", err)
	}

	user1ID, err := client1.AuthClient.TokenManager.GetUserIdFromAccessToken()
	if err != nil {
		t.Fatalf("Failed to get user1 ID: %v", err)
	}

	user2ID, err := client2.AuthClient.TokenManager.GetUserIdFromAccessToken()
	if err != nil {
		t.Fatalf("Failed to get user2 ID: %v", err)
	}

	// User2 logs out
	if err := client2.AuthClient.LogoutUser(); err != nil {
		t.Fatalf("Failed to logout user2: %v", err)
	}
	time.Sleep(2 * time.Second) // Sleep to allow logout to cancel context, stream

	// User1 sends a message to loggedout User2
	t.Log("User 1 sending unencrypted message to User 2")
	messageFromUser1 := "custom message from User 1 to User 2"
	if err := client1.ChatClient.SendUnencryptedMessage(context.Background(), user2ID, messageFromUser1); err != nil {
		t.Fatalf("Failed to send message from User 1 to User 2: %v", err)
	}

	time.Sleep(2 * time.Second) // Sleep to allow the message to be processed (before context cancel from test cleanup)

	// Make sure the message is saved in the chat history of User1 (delivered status should be 0)
	chatMessages, err := client1.ChatClient.Store.GetChatHistory(user1ID, user2ID)
	if err != nil {
		t.Fatalf("Failed to get chat history for User 1: %v", err)
	}

	msg := chatMessages[0]
	if msg.Delivered != 0 {
		t.Fatalf("Expected delivered status 0, but got: %d", msg.Delivered)
	}

	// Make sure the User2's chat history is empty
	chatMessages, err = client2.ChatClient.Store.GetChatHistory(user1ID, user2ID)
	if err != nil {
		t.Fatalf("Failed to get chat history for User 2: %v", err)
	}
	if len(chatMessages) != 0 {
		t.Fatalf("Expected empty chat history for User 2, but got: %d", len(chatMessages))
	}

	// Message should be in server buffer for later delivery
	if len(server.ChatServer.UndeliveredMessages) != 1 {
		t.Fatalf("Expected 1 undelivered message in server buffer, but got: %d", len(server.ChatServer.UndeliveredMessages))
	}

	// Log in User2
	if err, _ := client2.AuthClient.LoginUser("user2", "password"); err != nil {
		t.Fatalf("Failed to login user2: %v", err)
	}
	time.Sleep(2 * time.Second)

	// Buffer should be empty after User2 logs in
	if len(server.ChatServer.UndeliveredMessages) != 0 {
		t.Fatalf("Expected 0 undelivered message in server buffer, but got: %d", len(server.ChatServer.UndeliveredMessages))
	}

	// User 2 should have received the message
	chatMessages, err = client2.ChatClient.Store.GetChatHistory(user1ID, user2ID)
	if err != nil {
		t.Fatalf("Failed to get chat history for User 2: %v", err)
	}
	msgReceived := chatMessages[0]
	if msgReceived.Message != messageFromUser1 {
		t.Fatalf("Expected message %s, but got: %s", messageFromUser1, msgReceived.Message)
	}
}
