package rpc

import (
	"testing"

	"github.com/johnkhk/cli_chat_app/test/setup"
)

// TestSendFriendRequestAndVerifyStatus tests the flow of sending a friend request and verifying its status.
func TestSendFriendRequestAndVerifyStatus(t *testing.T) {
	t.Parallel()

	// Initialize resources with two clients for two different users
	rpcClients, _, cleanup := setup.InitializeTestResources(t, nil, 2)
	defer cleanup()

	client1 := rpcClients[0] // Represents User1
	client2 := rpcClients[1] // Represents User2

	// Register and login two users
	if err := client1.AuthClient.RegisterUser("user1", "password"); err != nil {
		t.Fatalf("Failed to register user1: %v", err)
	}
	if err := client1.AuthClient.LoginUser("user1", "password"); err != nil {
		t.Fatalf("Failed to login user1: %v", err)
	}
	if err := client2.AuthClient.RegisterUser("user2", "password"); err != nil {
		t.Fatalf("Failed to register user2: %v", err)
	}
	if err := client2.AuthClient.LoginUser("user2", "password"); err != nil {
		t.Fatalf("Failed to login user2: %v", err)
	}

	// User1 sends a friend request to User2
	if err := client1.FriendsClient.SendFriendRequest("user2"); err != nil {
		t.Fatalf("User1 failed to send friend request to user2: %v", err)
	}

	// User1 verifies the friend request is in outgoing requests
	outgoingRequests, err := client1.FriendsClient.GetOutgoingFriendRequests()
	if err != nil {
		t.Fatalf("Failed to get outgoing friend requests for user1: %v", err)
	}
	if len(outgoingRequests) != 1 {
		t.Fatalf("Expected 1 outgoing friend request for user1, got: %d", len(outgoingRequests))
	}

	// User2 verifies the friend request is in incoming requests
	incomingRequests, err := client2.FriendsClient.GetIncomingFriendRequests()
	if err != nil {
		t.Fatalf("Failed to get incoming friend requests for user2: %v", err)
	}
	if len(incomingRequests) != 1 {
		t.Fatalf("Expected 1 incoming friend request for user2, got: %d", len(incomingRequests))
	}
}

// TestAcceptFriendRequestAndVerify tests the flow of accepting a friend request.
func TestAcceptFriendRequestAndVerify(t *testing.T) {
	t.Parallel()

	// Initialize resources with two clients for two different users
	rpcClients, _, cleanup := setup.InitializeTestResources(t, nil, 2)
	defer cleanup()

	client1 := rpcClients[0] // Represents User1
	client2 := rpcClients[1] // Represents User2

	// Register and login two users
	if err := client1.AuthClient.RegisterUser("user1", "password"); err != nil {
		t.Fatalf("Failed to register user1: %v", err)
	}
	if err := client1.AuthClient.LoginUser("user1", "password"); err != nil {
		t.Fatalf("Failed to login user1: %v", err)
	}
	if err := client2.AuthClient.RegisterUser("user2", "password"); err != nil {
		t.Fatalf("Failed to register user2: %v", err)
	}
	if err := client2.AuthClient.LoginUser("user2", "password"); err != nil {
		t.Fatalf("Failed to login user2: %v", err)
	}

	// User1 sends a friend request to User2
	if err := client1.FriendsClient.SendFriendRequest("user2"); err != nil {
		t.Fatalf("User1 failed to send friend request to user2: %v", err)
	}

	// User2 accepts the friend request
	incomingRequests, err := client2.FriendsClient.GetIncomingFriendRequests()
	if err != nil {
		t.Fatalf("Failed to get incoming friend requests for user2: %v", err)
	}
	if len(incomingRequests) == 0 {
		t.Fatalf("No incoming friend requests found for user2")
	}
	if err := client2.FriendsClient.AcceptFriendRequest(incomingRequests[0].RequestId); err != nil {
		t.Fatalf("Failed to accept friend request: %v", err)
	}

	// Verify that no incoming or outgoing requests exist after acceptance
	incomingRequests, err = client2.FriendsClient.GetIncomingFriendRequests()
	if err != nil {
		t.Fatalf("Failed to get incoming friend requests after accepting: %v", err)
	}
	if len(incomingRequests) != 0 {
		t.Fatalf("Expected 0 incoming friend requests for user2 after acceptance, got: %d", len(incomingRequests))
	}

	outgoingRequests, err := client1.FriendsClient.GetOutgoingFriendRequests()
	if err != nil {
		t.Fatalf("Failed to get outgoing friend requests for user1 after acceptance: %v", err)
	}
	if len(outgoingRequests) != 0 {
		t.Fatalf("Expected 0 outgoing friend requests for user1 after acceptance, got: %d", len(outgoingRequests))
	}
}

// TestGetFriendListAfterAcceptingRequest tests retrieving the friend list after a request is accepted.
func TestGetFriendListAfterAcceptingRequest(t *testing.T) {
	t.Parallel()

	// Initialize resources with two clients for two different users
	rpcClients, _, cleanup := setup.InitializeTestResources(t, nil, 2)
	defer cleanup()

	client1 := rpcClients[0] // Represents User1
	client2 := rpcClients[1] // Represents User2

	// Register and login two users
	if err := client1.AuthClient.RegisterUser("user1", "password"); err != nil {
		t.Fatalf("Failed to register user1: %v", err)
	}
	if err := client1.AuthClient.LoginUser("user1", "password"); err != nil {
		t.Fatalf("Failed to login user1: %v", err)
	}
	if err := client2.AuthClient.RegisterUser("user2", "password"); err != nil {
		t.Fatalf("Failed to register user2: %v", err)
	}
	if err := client2.AuthClient.LoginUser("user2", "password"); err != nil {
		t.Fatalf("Failed to login user2: %v", err)
	}

	// User1 sends a friend request to User2
	if err := client1.FriendsClient.SendFriendRequest("user2"); err != nil {
		t.Fatalf("User1 failed to send friend request to user2: %v", err)
	}

	// User2 accepts the friend request
	incomingRequests, err := client2.FriendsClient.GetIncomingFriendRequests()
	if err != nil {
		t.Fatalf("Failed to get incoming friend requests for user2: %v", err)
	}
	if len(incomingRequests) == 0 {
		t.Fatalf("No incoming friend requests found for user2")
	}
	if err := client2.FriendsClient.AcceptFriendRequest(incomingRequests[0].RequestId); err != nil {
		t.Fatalf("Failed to accept friend request: %v", err)
	}

	// Both users should now appear in each other's friend lists
	user1Friends, err := client1.FriendsClient.GetFriendList()
	if err != nil {
		t.Fatalf("Failed to get friend list for user1: %v", err)
	}
	if len(user1Friends) != 1 {
		t.Fatalf("Expected 1 friend for user1, got: %d", len(user1Friends))
	}

	user2Friends, err := client2.FriendsClient.GetFriendList()
	if err != nil {
		t.Fatalf("Failed to get friend list for user2: %v", err)
	}
	if len(user2Friends) != 1 {
		t.Fatalf("Expected 1 friend for user2, got: %d", len(user2Friends))
	}
}
