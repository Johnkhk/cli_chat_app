package rpc

import (
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
