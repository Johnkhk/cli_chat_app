package app

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/johnkhk/cli_chat_app/genproto/friends"
)

// FriendsServer implements the FriendsService.
type FriendsServer struct {
	friends.UnimplementedFriendsServiceServer
	DB     *sql.DB
	Logger *logrus.Logger
}

// NewFriendsServer creates a new FriendsServer with the given dependencies.
func NewFriendsServer(db *sql.DB, logger *logrus.Logger) *FriendsServer {
	return &FriendsServer{
		DB:     db,
		Logger: logger,
	}
}

// SendFriendRequest handles sending a friend request.
func (s *FriendsServer) SendFriendRequest(ctx context.Context, req *friends.SendFriendRequestRequest) (*friends.SendFriendRequestResponse, error) {
	// Retrieve the requester ID from the context
	requesterID, ok := ctx.Value("userID").(string)
	if !ok || requesterID == "" {
		return nil, fmt.Errorf("requester ID not found in context")
	}

	// Retrieve the requester username from the context
	requesterUsername, ok := ctx.Value("username").(string)
	if !ok || requesterUsername == "" {
		return nil, fmt.Errorf("requester username not found in context")
	}

	s.Logger.Infof("Received friend request from user ID: %s (username: %s) to username: %s", requesterID, requesterUsername, req.RecipientUsername)

	// Step 1: Retrieve the recipient's ID from the username
	var recipientID int
	err := s.DB.QueryRow("SELECT id FROM users WHERE username = ?", req.RecipientUsername).Scan(&recipientID)
	if err != nil {
		if err == sql.ErrNoRows {
			s.Logger.Warnf("Recipient username does not exist: %s", req.RecipientUsername)
			return &friends.SendFriendRequestResponse{
				Success: false,
				Message: "Recipient username does not exist",
			}, nil
		}
		s.Logger.Errorf("Error retrieving recipient ID: %v", err)
		return nil, fmt.Errorf("error retrieving recipient ID: %w", err)
	}

	// Convert requesterID from string to int
	requesterIDInt, err := strconv.Atoi(requesterID)
	if err != nil {
		return nil, fmt.Errorf("invalid requester ID format: %w", err)
	}

	// Step 2: Check if a friend request is already pending, accepted, or rejected
	var existingStatus string
	err = s.DB.QueryRow(`
		SELECT status 
		FROM friend_requests 
		WHERE (requester_id = ? AND recipient_id = ?) 
		   OR (requester_id = ? AND recipient_id = ?)`,
		requesterIDInt, recipientID, recipientID, requesterIDInt).Scan(&existingStatus)

	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error checking existing friend request: %w", err)
	}

	// If a record is found, handle the different statuses
	if err == nil {
		switch existingStatus {
		case "pending":
			return &friends.SendFriendRequestResponse{
				Success: false,
				Message: "A friend request is already pending",
			}, nil
		case "accepted":
			return &friends.SendFriendRequestResponse{
				Success: false,
				Message: "You are already friends",
			}, nil
		case "rejected":
			// Allow sending the friend request again if it was previously rejected
			_, err = s.DB.Exec(`
				UPDATE friend_requests
				SET status = 'pending', requested_at = NOW()
				WHERE requester_id = ? AND recipient_id = ? AND status = 'rejected'`,
				requesterIDInt, recipientID)
			if err != nil {
				return nil, fmt.Errorf("error updating friend request to pending: %w", err)
			}
			return &friends.SendFriendRequestResponse{
				Success: true,
				Message: "Friend request sent again successfully",
			}, nil
		}
	}

	// Step 3: Insert the new friend request into the database
	_, err = s.DB.Exec(`
		INSERT INTO friend_requests (requester_id, recipient_id, requested_at, status)
		VALUES (?, ?, NOW(), 'pending')`,
		requesterIDInt, recipientID)
	if err != nil {
		return nil, fmt.Errorf("error inserting friend request into database: %w", err)
	}

	return &friends.SendFriendRequestResponse{
		Success: true,
		Message: "Friend request sent successfully",
	}, nil
}

// GetFriendRequests handles fetching pending friend requests.
func (s *FriendsServer) GetFriendRequests(ctx context.Context, req *friends.GetFriendRequestsRequest) (*friends.GetFriendRequestsResponse, error) {
	s.Logger.Info("Fetching pending friend requests")

	// Mock implementation
	return &friends.GetFriendRequestsResponse{
		Success:        true,
		Message:        "Friend requests fetched successfully",
		FriendRequests: []*friends.FriendRequest{}, // Empty list for now
	}, nil
}

// RespondToFriendRequest handles accepting or rejecting a friend request.
// Takes in the requester ID and a boolean indicating whether to accept the request.
func (s *FriendsServer) RespondToFriendRequest(ctx context.Context, req *friends.RespondToFriendRequestRequest) (*friends.RespondToFriendRequestResponse, error) {
	// Retrieve the recipient's user ID from the context
	recipientID, ok := ctx.Value("userID").(string)
	if !ok || recipientID == "" {
		return nil, fmt.Errorf("recipient ID not found in context")
	}

	// Convert recipient ID from string to int
	recipientIDInt, err := strconv.Atoi(recipientID)
	if err != nil {
		return nil, fmt.Errorf("invalid recipient ID format: %w", err)
	}

	// Convert requester ID from string to int
	requesterIDInt, err := strconv.Atoi(req.RequesterId)
	if err != nil {
		return nil, fmt.Errorf("invalid requester ID format: %w", err)
	}

	// Check if there is a pending friend request from requester to recipient
	var status string
	err = s.DB.QueryRow(`
		SELECT status 
		FROM friend_requests 
		WHERE requester_id = ? AND recipient_id = ? AND status = 'pending'`,
		requesterIDInt, recipientIDInt).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return &friends.RespondToFriendRequestResponse{
				Success: false,
				Message: "No pending friend request found",
			}, nil
		}
		return nil, fmt.Errorf("error checking pending friend request: %w", err)
	}

	if req.Accept {
		// Accepting the friend request
		tx, err := s.DB.Begin()
		if err != nil {
			return nil, fmt.Errorf("failed to begin transaction: %w", err)
		}

		// Step 1: Update the friend request status to 'accepted'
		_, err = tx.Exec(`
			UPDATE friend_requests 
			SET status = 'accepted' 
			WHERE requester_id = ? AND recipient_id = ? AND status = 'pending'`,
			requesterIDInt, recipientIDInt)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to update friend request status: %w", err)
		}

		// Step 2: Insert the friendship into the friends table
		_, err = tx.Exec(`
			INSERT INTO friends (user_id, friend_id, added_at) 
			VALUES (?, ?, NOW()), (?, ?, NOW())`,
			requesterIDInt, recipientIDInt, recipientIDInt, requesterIDInt)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to insert into friends table: %w", err)
		}

		// Commit the transaction
		if err := tx.Commit(); err != nil {
			return nil, fmt.Errorf("failed to commit transaction: %w", err)
		}

		return &friends.RespondToFriendRequestResponse{
			Success: true,
			Message: "Friend request accepted successfully",
		}, nil
	} else {
		// Rejecting the friend request
		_, err := s.DB.Exec(`
			UPDATE friend_requests 
			SET status = 'rejected' 
			WHERE requester_id = ? AND recipient_id = ? AND status = 'pending'`,
			requesterIDInt, recipientIDInt)
		if err != nil {
			return nil, fmt.Errorf("failed to update friend request status to rejected: %w", err)
		}

		return &friends.RespondToFriendRequestResponse{
			Success: true,
			Message: "Friend request rejected successfully",
		}, nil
	}
}

// GetFriend handles fetching a specific friend's details.
func (s *FriendsServer) GetFriend(ctx context.Context, req *friends.GetFriendRequest) (*friends.GetFriendResponse, error) {
	s.Logger.Infof("Fetching friend details for friend ID: %s", req.FriendId)

	// Mock implementation
	return &friends.GetFriendResponse{
		Success: true,
		Message: "Friend details fetched successfully",
		Friend: &friends.Friend{
			Id:       req.FriendId,
			Username: "MockFriend",
			Status:   "offline",
			AddedAt:  "2024-01-01T00:00:00Z",
		},
	}, nil
}

// GetFriends handles fetching all friends for the user.
func (s *FriendsServer) GetFriends(ctx context.Context, req *friends.GetFriendsRequest) (*friends.GetFriendsResponse, error) {
	s.Logger.Info("Fetching all friends")

	// Mock implementation
	return &friends.GetFriendsResponse{
		Success: true,
		Message: "Friends fetched successfully",
		Friends: []*friends.Friend{}, // Empty list for now
	}, nil
}

// RemoveFriend handles removing a friend.
func (s *FriendsServer) RemoveFriend(ctx context.Context, req *friends.RemoveFriendRequest) (*friends.RemoveFriendResponse, error) {
	s.Logger.Infof("Removing friend with ID: %s", req.FriendId)

	// Mock implementation
	return &friends.RemoveFriendResponse{Success: true, Message: "Friend removed successfully"}, nil
}
