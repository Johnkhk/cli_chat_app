package rpc

import (
	"testing"

	"github.com/johnkhk/cli_chat_app/genproto/friends"
	"github.com/johnkhk/cli_chat_app/server/storage"
	"github.com/johnkhk/cli_chat_app/test/setup"
)

// TestSendFriendRequestAndVerifyStatus tests the flow of sending a friend request and verifying its status, including usernames.
func TestSendFriendRequestAndVerifyStatus(t *testing.T) {
	// Initialize resources with two clients for two different users
	rpcClients, _, cleanup, _ := setup.InitializeTestResources(t, nil, 2)
	defer cleanup()

	client1 := rpcClients[0] // Represents User1
	client2 := rpcClients[1] // Represents User2

	// User details
	user1Username := "user1"
	user2Username := "user2"
	password := "password"

	// Register and login two users
	if err := client1.AuthClient.RegisterUser(user1Username, password); err != nil {
		t.Fatalf("Failed to register %s: %v", user1Username, err)
	}
	if err, _ := client1.AuthClient.LoginUser(user1Username, password); err != nil {
		t.Fatalf("Failed to login %s: %v", user1Username, err)
	}
	if err := client2.AuthClient.RegisterUser(user2Username, password); err != nil {
		t.Fatalf("Failed to register %s: %v", user2Username, err)
	}
	if err, _ := client2.AuthClient.LoginUser(user2Username, password); err != nil {
		t.Fatalf("Failed to login %s: %v", user2Username, err)
	}

	// User1 sends a friend request to User2
	if err := client1.FriendsClient.SendFriendRequest(user2Username); err != nil {
		t.Fatalf("%s failed to send friend request to %s: %v", user1Username, user2Username, err)
	}

	// User1 verifies the friend request is in outgoing requests with correct usernames
	outgoingRequests, err := client1.FriendsClient.GetOutgoingFriendRequests()
	if err != nil {
		t.Fatalf("Failed to get outgoing friend requests for %s: %v", user1Username, err)
	}
	if len(outgoingRequests) != 1 {
		t.Fatalf("Expected 1 outgoing friend request for %s, got: %d", user1Username, len(outgoingRequests))
	}

	// Verify sender_username and recipient_username in outgoing request
	outReq := outgoingRequests[0]
	if outReq.SenderUsername != user1Username {
		t.Errorf("Expected sender_username to be %s, got: %s", user1Username, outReq.SenderUsername)
	}
	if outReq.RecipientUsername != user2Username {
		t.Errorf("Expected recipient_username to be %s, got: %s", user2Username, outReq.RecipientUsername)
	}

	// Verify status of the friend request
	if outReq.Status != friends.FriendRequestStatus_PENDING {
		t.Errorf("Expected friend request status to be %s, got: %s", friends.FriendRequestStatus_PENDING, outReq.Status)
	}

	// User2 verifies the friend request is in incoming requests with correct usernames
	incomingRequests, err := client2.FriendsClient.GetIncomingFriendRequests()
	if err != nil {
		t.Fatalf("Failed to get incoming friend requests for %s: %v", user2Username, err)
	}
	if len(incomingRequests) != 1 {
		t.Fatalf("Expected 1 incoming friend request for %s, got: %d", user2Username, len(incomingRequests))
	}

	// Verify sender_username and recipient_username in incoming request
	inReq := incomingRequests[0]
	if inReq.SenderUsername != user1Username {
		t.Errorf("Expected sender_username to be %s, got: %s", user1Username, inReq.SenderUsername)
	}
	if inReq.RecipientUsername != user2Username {
		t.Errorf("Expected recipient_username to be %s, got: %s", user2Username, inReq.RecipientUsername)
	}

	// Verify status of the friend request
	if inReq.Status != friends.FriendRequestStatus_PENDING {
		t.Errorf("Expected friend request status to be %s, got: %s", friends.FriendRequestStatus_PENDING, inReq.Status)
	}
}

// TestAcceptFriendRequestAndVerify tests the flow of accepting a friend request, including verifying usernames.
func TestAcceptFriendRequestAndVerify(t *testing.T) {
	// Initialize resources with two clients for two different users
	rpcClients, _, cleanup, _ := setup.InitializeTestResources(t, nil, 2)
	defer cleanup()

	client1 := rpcClients[0] // Represents User1
	client2 := rpcClients[1] // Represents User2

	// User details
	user1Username := "user1"
	user2Username := "user2"
	password := "password"

	// Register and login two users
	if err := client1.AuthClient.RegisterUser(user1Username, password); err != nil {
		t.Fatalf("Failed to register %s: %v", user1Username, err)
	}
	if err, _ := client1.AuthClient.LoginUser(user1Username, password); err != nil {
		t.Fatalf("Failed to login %s: %v", user1Username, err)
	}
	if err := client2.AuthClient.RegisterUser(user2Username, password); err != nil {
		t.Fatalf("Failed to register %s: %v", user2Username, err)
	}
	if err, _ := client2.AuthClient.LoginUser(user2Username, password); err != nil {
		t.Fatalf("Failed to login %s: %v", user2Username, err)
	}

	// User1 sends a friend request to User2
	if err := client1.FriendsClient.SendFriendRequest(user2Username); err != nil {
		t.Fatalf("%s failed to send friend request to %s: %v", user1Username, user2Username, err)
	}

	// User2 gets the friend request
	incomingRequests, err := client2.FriendsClient.GetIncomingFriendRequests()
	if err != nil {
		t.Fatalf("Failed to get incoming friend requests for %s: %v", user2Username, err)
	}
	if len(incomingRequests) == 0 {
		t.Fatalf("No incoming friend requests found for %s", user2Username)
	}

	// Verify usernames in the incoming request before accepting
	inReq := incomingRequests[0]
	if inReq.SenderUsername != user1Username {
		t.Errorf("Expected sender_username to be %s, got: %s", user1Username, inReq.SenderUsername)
	}
	if inReq.RecipientUsername != user2Username {
		t.Errorf("Expected recipient_username to be %s, got: %s", user2Username, inReq.RecipientUsername)
	}

	// Accept the friend request
	if err := client2.FriendsClient.AcceptFriendRequest(inReq.RequestId); err != nil {
		t.Fatalf("Failed to accept friend request: %v", err)
	}

	// Verify that no incoming or outgoing requests exist after acceptance
	incomingRequests, err = client2.FriendsClient.GetIncomingFriendRequests()
	if err != nil {
		t.Fatalf("Failed to get incoming friend requests after accepting: %v", err)
	}
	if len(incomingRequests) != 0 {
		t.Fatalf("Expected 0 incoming friend requests for %s after acceptance, got: %d", user2Username, len(incomingRequests))
	}

	outgoingRequests, err := client1.FriendsClient.GetOutgoingFriendRequests()
	if err != nil {
		t.Fatalf("Failed to get outgoing friend requests for %s after acceptance: %v", user1Username, err)
	}
	if len(outgoingRequests) != 0 {
		t.Fatalf("Expected 0 outgoing friend requests for %s after acceptance, got: %d", user1Username, len(outgoingRequests))
	}
}

// TestGetFriendListAfterAcceptingRequest tests retrieving the friend list after a request is accepted.
func TestGetFriendListAfterAcceptingRequest(t *testing.T) {
	// t.Parallel()

	// Initialize resources with two clients for two different users
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

// TestDeclineFriendRequestAndVerify tests the flow of declining a friend request.
func TestDeclineFriendRequestAndVerify(t *testing.T) {
	// t.Parallel()

	// Initialize resources with two clients for two different users
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

	// User1 sends a friend request to User2
	if err := client1.FriendsClient.SendFriendRequest("user2"); err != nil {
		t.Fatalf("User1 failed to send friend request to user2: %v", err)
	}

	// User2 declines the friend request
	incomingRequests, err := client2.FriendsClient.GetIncomingFriendRequests()
	if err != nil {
		t.Fatalf("Failed to get incoming friend requests for user2: %v", err)
	}
	if len(incomingRequests) == 0 {
		t.Fatalf("No incoming friend requests found for user2")
	}
	if err := client2.FriendsClient.DeclineFriendRequest(incomingRequests[0].RequestId); err != nil {
		t.Fatalf("Failed to decline friend request: %v", err)
	}

	// Verify that the friend request is removed from the incoming requests for User2
	incomingRequests, err = client2.FriendsClient.GetIncomingFriendRequests()
	if err != nil {
		t.Fatalf("Failed to get incoming friend requests after declining: %v", err)
	}
	if len(incomingRequests) != 0 {
		t.Fatalf("Expected 0 incoming friend requests for user2 after declining, got: %d", len(incomingRequests))
	}

	// Verify that the friend request is also removed from the outgoing requests for User1
	outgoingRequests, err := client1.FriendsClient.GetOutgoingFriendRequests()
	if err != nil {
		t.Fatalf("Failed to get outgoing friend requests for user1 after declining: %v", err)
	}
	if len(outgoingRequests) != 0 {
		t.Fatalf("Expected 0 outgoing friend requests for user1 after declining, got: %d", len(outgoingRequests))
	}
}

// TestRemoveFriendAndVerify tests the flow of removing a friend.
func TestRemoveFriendAndVerify(t *testing.T) {
	// t.Parallel()

	// Initialize resources with two clients for two different users
	rpcClients, db, cleanup, _ := setup.InitializeTestResources(t, nil, 2)
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

	// User1 sends a friend request to User2
	if err := client1.FriendsClient.SendFriendRequest("user2"); err != nil {
		t.Fatalf("User1 failed to send friend request to user2: %v", err)
	}

	// Make sure status is pending in User2's incoming requests
	incomingRequests, err := client2.FriendsClient.GetIncomingFriendRequests()
	if err != nil {
		t.Fatalf("Failed to get incoming friend requests for user2: %v", err)
	}
	if len(incomingRequests) == 0 {
		t.Fatalf("No incoming friend requests found for user2")
	}
	if incomingRequests[0].Status != friends.FriendRequestStatus_PENDING {
		t.Fatalf("Expected friend request status to be PENDING, got: %s", incomingRequests[0].Status)
	}

	// Make sure status is pending in db
	var status string
	err = db.QueryRow("SELECT status FROM friend_requests WHERE recipient_id = ? AND requester_id = ?", 2, 1).Scan(&status)
	if err != nil {
		t.Fatalf("Failed to query status from friend_requests table: %v", err)
	}
	if status != storage.StatusPendingStr {
		t.Fatalf("Expected friend request status to be PENDING in db, got: %s", status)
	}

	// Accept the friend request
	if err := client2.FriendsClient.AcceptFriendRequest(incomingRequests[0].RequestId); err != nil {
		t.Fatalf("Failed to accept friend request: %v", err)
	}

	// Make sure status is accepted in db
	err = db.QueryRow("SELECT status FROM friend_requests WHERE recipient_id = ? AND requester_id = ?", 2, 1).Scan(&status)
	if err != nil {
		t.Fatalf("Failed to query status from friend_requests table: %v", err)
	}
	if status != storage.StatusAcceptedStr {
		t.Fatalf("Expected friend request status to be ACCEPTED in db, got: %s", status)
	}

	// Both users verify they are friends
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

	// User1 removes User2 from the friend list
	if err := client1.FriendsClient.RemoveFriend(user1Friends[0].UserId); err != nil {
		t.Fatalf("Failed to remove friend: %v", err)
	}

	// Make sure status is CANCELLED in db
	err = db.QueryRow("SELECT status FROM friend_requests WHERE recipient_id = ? AND requester_id = ?", 2, 1).Scan(&status)
	if err != nil {
		t.Fatalf("Failed to query status from friend_requests table: %v", err)
	}
	if status != storage.StatusCancelledStr {
		t.Fatalf("Expected friend request status to be CANCELLED in db, got: %s", status)
	}

	// Verify that User2 is no longer in User1's friend list
	user1Friends, err = client1.FriendsClient.GetFriendList()
	if err != nil {
		t.Fatalf("Failed to get friend list for user1 after removing friend: %v", err)
	}
	if len(user1Friends) != 0 {
		t.Fatalf("Expected 0 friends for user1 after removing, got: %d", len(user1Friends))
	}

	// Verify that User1 is no longer in User2's friend list
	user2Friends, err = client2.FriendsClient.GetFriendList()
	if err != nil {
		t.Fatalf("Failed to get friend list for user2 after being removed: %v", err)
	}
	if len(user2Friends) != 0 {
		t.Fatalf("Expected 0 friends for user2 after being removed, got: %d", len(user2Friends))
	}

	// User1 Sends another friend request to User2
	if err := client1.FriendsClient.SendFriendRequest("user2"); err != nil {
		t.Fatalf("User1 failed to send friend request to user2: %v", err)
	}
	if err != nil {
		t.Fatalf("Failed to send friend request a second time: %v", err)
	}

	// Make sure status is pending in User2's incoming requests
	incomingRequests, err = client2.FriendsClient.GetIncomingFriendRequests()
	if err != nil {
		t.Fatalf("Failed to get incoming friend requests for user2: %v", err)
	}
	if len(incomingRequests) == 0 {
		t.Fatalf("No incoming friend requests found for user2")
	}
	if incomingRequests[0].Status != friends.FriendRequestStatus_PENDING {
		t.Fatalf("Expected friend request status to be PENDING, got: %s", incomingRequests[0].Status)
	}

	// Make sure status is pending in db
	err = db.QueryRow("SELECT status FROM friend_requests WHERE recipient_id = ? AND requester_id = ?", 2, 1).Scan(&status)
	if err != nil {
		t.Fatalf("Failed to query status from friend_requests table: %v", err)
	}
	if status != storage.StatusPendingStr {
		t.Fatalf("Expected friend request status to be PENDING in db, got: %s", status)
	}
}
