package app

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/johnkhk/cli_chat_app/genproto/friends"
	"github.com/johnkhk/cli_chat_app/server/storage"
)

// FriendsServer implements the FriendsService.
type FriendsServer struct {
	friends.UnimplementedFriendManagementServer
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

	// Convert requesterID from string to int
	requesterIDInt, err := strconv.Atoi(requesterID)
	if err != nil {
		return nil, fmt.Errorf("invalid requester ID format: %w", err)
	}

	// Retrieve the requester username from the context
	requesterUsername, ok := ctx.Value("username").(string)
	if !ok || requesterUsername == "" {
		return nil, fmt.Errorf("requester username not found in context")
	}

	// Check if the requester is trying to send a request to themselves
	if requesterUsername == req.RecipientUsername {
		return &friends.SendFriendRequestResponse{
			Status:    friends.FriendRequestStatus_FAILED,
			Message:   "Cannot send a friend request to yourself",
			Timestamp: timestamppb.Now(),
		}, nil
	}

	s.Logger.Infof("Received friend request from user ID: %s (username: %s) to username: %s", requesterID, requesterUsername, req.RecipientUsername)

	// Step 1: Retrieve the recipient's ID from the username
	var recipientID int
	err = s.DB.QueryRow("SELECT id FROM users WHERE username = ?", req.RecipientUsername).Scan(&recipientID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Recipient not found
			return &friends.SendFriendRequestResponse{
				Status:    friends.FriendRequestStatus_FAILED,
				Message:   "Recipient not found",
				Timestamp: timestamppb.Now(),
			}, nil
		}
		s.Logger.Errorf("Error retrieving recipient ID: %v", err)
		return nil, fmt.Errorf("error retrieving recipient ID: %w", err)
	}

	// Step 2: Check if a friend request already exists
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

	if err == nil { // If there is an existing friend request
		s.Logger.Infof("Existing status for friend request: %s", existingStatus)

		statusEnum, ok := friends.FriendRequestStatus_value[existingStatus]
		if !ok {
			// If the status does not match any known value, use UNKNOWN
			statusEnum = storage.StatusUnknownInt32
		}

		// Handle the relevant existing statuses
		switch friends.FriendRequestStatus(statusEnum) {
		case friends.FriendRequestStatus_PENDING:
			return &friends.SendFriendRequestResponse{
				Status:    friends.FriendRequestStatus_FAILED,
				Message:   "A friend request is already pending",
				Timestamp: timestamppb.Now(),
			}, nil
		case friends.FriendRequestStatus_ACCEPTED:
			return &friends.SendFriendRequestResponse{
				Status:    friends.FriendRequestStatus_FAILED,
				Message:   "You are already friends",
				Timestamp: timestamppb.Now(),
			}, nil
		case friends.FriendRequestStatus_DECLINED, friends.FriendRequestStatus_CANCELED:
			// Allow sending the friend request again if it was previously declined or canceled
			_, err = s.DB.Exec(`
				UPDATE friend_requests
				SET status = ?, created_at = NOW()
				WHERE requester_id = ? AND recipient_id = ? AND status IN (?, ?)`,
				storage.StatusPendingStr, requesterIDInt, recipientID, storage.StatusDeclinedStr, storage.StatusCancelledStr)
			if err != nil {
				return nil, fmt.Errorf("error updating friend request to pending: %w", err)
			}
			return &friends.SendFriendRequestResponse{
				Status:    friends.FriendRequestStatus_PENDING,
				Message:   "Friend request sent again successfully",
				Timestamp: timestamppb.Now(),
			}, nil
		}
	}

	// Step 3: Insert a new friend request if no existing request
	_, err = s.DB.Exec(`
		INSERT INTO friend_requests (requester_id, recipient_id, status)
		VALUES (?, ?, ?)`,
		requesterIDInt, recipientID, storage.StatusPendingStr)
	if err != nil {
		return nil, fmt.Errorf("error inserting friend request into database: %w", err)
	}

	return &friends.SendFriendRequestResponse{
		Status:    friends.FriendRequestStatus_PENDING,
		Message:   "Friend request sent successfully",
		Timestamp: timestamppb.Now(),
	}, nil
}

// AcceptFriendRequest handles accepting a friend request.
func (s *FriendsServer) AcceptFriendRequest(ctx context.Context, req *friends.AcceptFriendRequestRequest) (*friends.AcceptFriendRequestResponse, error) {
	// Retrieve the user ID from the context
	userID, ok := ctx.Value("userID").(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("user ID not found in context")
	}

	// Convert userID from string to int
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	// Step 1: Update the friend request status to "ACCEPTED" if it exists and is pending
	res, err := s.DB.Exec(`UPDATE friend_requests SET status = ?, response_at = NOW() WHERE id = ? AND recipient_id = ? AND status = ?`,
		storage.StatusAcceptedStr, req.RequestId, userIDInt, storage.StatusPendingStr)

	if err != nil {
		return nil, fmt.Errorf("error updating friend request status to accepted: %w", err)
	}

	// Step 2: Check if exactly one row was affected by the update
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("error checking affected rows: %w", err)
	}

	if rowsAffected != 1 {
		// No rows were affected, indicating that the request does not exist or is not pending
		return &friends.AcceptFriendRequestResponse{
			Status:    friends.FriendRequestStatus_FAILED,
			Message:   "Friend request does not exist or is not pending",
			Timestamp: timestamppb.Now(),
		}, nil
	}

	// Step 3: Retrieve the requester ID from the friend request
	var requesterID int
	err = s.DB.QueryRow(`SELECT requester_id FROM friend_requests WHERE id = ?`, req.RequestId).Scan(&requesterID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving requester ID: %w", err)
	}

	// Step 4: Insert the new friendship into the friends table
	_, err = s.DB.Exec(`
		INSERT INTO friends (user_id, friend_id, created_at) VALUES (?, ?, NOW()), (?, ?, NOW())`,
		userIDInt, requesterID, requesterID, userIDInt)

	if err != nil {
		return nil, fmt.Errorf("error inserting into friends table: %w", err)
	}

	// Step 5: Return a successful response
	return &friends.AcceptFriendRequestResponse{
		Status:    friends.FriendRequestStatus_ACCEPTED,
		Message:   "Friend request accepted successfully",
		Timestamp: timestamppb.Now(),
	}, nil
}

// GetIncomingFriendRequests retrieves the incoming friend requests for the user.
func (s *FriendsServer) GetIncomingFriendRequests(ctx context.Context, req *friends.GetIncomingFriendRequestsRequest) (*friends.GetIncomingFriendRequestsResponse, error) {
	// Retrieve the user ID from the context (e.g., extracted from the token)
	userID, ok := ctx.Value("userID").(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("user ID not found in context")
	}

	// Convert userID from string to int
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	// Query to get all incoming friend requests for this user, including usernames
	rows, err := s.DB.Query(`
        SELECT fr.id, fr.requester_id, fr.recipient_id, fr.status, fr.created_at,
               u_sender.username AS sender_username,
               u_recipient.username AS recipient_username
        FROM friend_requests fr
        JOIN users u_sender ON fr.requester_id = u_sender.id
        JOIN users u_recipient ON fr.recipient_id = u_recipient.id
        WHERE fr.recipient_id = ? AND fr.status = ?`, userIDInt, storage.StatusPendingStr)
	if err != nil {
		return nil, fmt.Errorf("error fetching incoming friend requests: %w", err)
	}
	defer rows.Close()

	// Prepare the response
	var incomingRequests []*friends.FriendRequest
	for rows.Next() {
		var friendReq friends.FriendRequest
		var createdAt time.Time
		var status string
		var senderUsername string
		var recipientUsername string

		if err := rows.Scan(
			&friendReq.RequestId,
			&friendReq.SenderId,
			&friendReq.RecipientId,
			&status,
			&createdAt,
			&senderUsername,
			&recipientUsername,
		); err != nil {
			return nil, fmt.Errorf("error scanning friend request row: %w", err)
		}

		// Convert status to enum and time to protobuf timestamp
		friendReq.Status = friends.FriendRequestStatus(friends.FriendRequestStatus_value[status])
		friendReq.CreatedAt = timestamppb.New(createdAt)

		// Set the usernames
		friendReq.SenderUsername = senderUsername
		friendReq.RecipientUsername = recipientUsername

		incomingRequests = append(incomingRequests, &friendReq)
	}

	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over incoming friend requests: %w", err)
	}

	return &friends.GetIncomingFriendRequestsResponse{
		IncomingRequests: incomingRequests,
	}, nil
}

// GetOutgoingFriendRequests retrieves the outgoing friend requests sent by the user.
func (s *FriendsServer) GetOutgoingFriendRequests(ctx context.Context, req *friends.GetOutgoingFriendRequestsRequest) (*friends.GetOutgoingFriendRequestsResponse, error) {
	// Retrieve the user ID from the context (e.g., extracted from the token)
	userID, ok := ctx.Value("userID").(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("user ID not found in context")
	}

	// Convert userID from string to int
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	// Query to get all outgoing friend requests for this user, including usernames
	rows, err := s.DB.Query(`
        SELECT fr.id, fr.requester_id, fr.recipient_id, fr.status, fr.created_at,
               u_sender.username AS sender_username,
               u_recipient.username AS recipient_username
        FROM friend_requests fr
        JOIN users u_sender ON fr.requester_id = u_sender.id
        JOIN users u_recipient ON fr.recipient_id = u_recipient.id
        WHERE fr.requester_id = ? AND fr.status = ?`, userIDInt, storage.StatusPendingStr)
	if err != nil {
		return nil, fmt.Errorf("error fetching outgoing friend requests: %w", err)
	}
	defer rows.Close()

	// Prepare the response
	var outgoingRequests []*friends.FriendRequest
	for rows.Next() {
		var friendReq friends.FriendRequest
		var createdAt time.Time
		var status string
		var senderUsername string
		var recipientUsername string

		if err := rows.Scan(
			&friendReq.RequestId,
			&friendReq.SenderId,
			&friendReq.RecipientId,
			&status,
			&createdAt,
			&senderUsername,
			&recipientUsername,
		); err != nil {
			return nil, fmt.Errorf("error scanning friend request row: %w", err)
		}

		// Convert status to enum and time to protobuf timestamp
		friendReq.Status = friends.FriendRequestStatus(friends.FriendRequestStatus_value[status])
		friendReq.CreatedAt = timestamppb.New(createdAt)

		// Set the usernames
		friendReq.SenderUsername = senderUsername
		friendReq.RecipientUsername = recipientUsername

		outgoingRequests = append(outgoingRequests, &friendReq)
	}

	// Check for errors after iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over outgoing friend requests: %w", err)
	}

	return &friends.GetOutgoingFriendRequestsResponse{
		OutgoingRequests: outgoingRequests,
	}, nil
}

// GetFriendList retrieves the list of friends for the user.
func (s *FriendsServer) GetFriendList(ctx context.Context, req *friends.GetFriendListRequest) (*friends.GetFriendListResponse, error) {
	// Retrieve the user ID from the context (e.g., extracted from the token)
	userID, ok := ctx.Value("userID").(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("user ID not found in context")
	}

	// Convert userID from string to int
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	// Query to get all friends for this user
	rows, err := s.DB.Query(`
        SELECT f.friend_id, u.username, f.created_at
        FROM friends f
        JOIN users u ON f.friend_id = u.id
        WHERE f.user_id = ?`, userIDInt)
	if err != nil {
		return nil, fmt.Errorf("error fetching friend list: %w", err)
	}
	defer rows.Close()

	// Prepare the response
	var friendsList []*friends.Friend
	for rows.Next() {
		var friend friends.Friend
		var addedAt time.Time

		// Scan the required fields
		if err := rows.Scan(&friend.UserId, &friend.Username, &addedAt); err != nil {
			return nil, fmt.Errorf("error scanning friend row: %w", err)
		}

		// Convert `added_at` to protobuf timestamp
		friend.AddedAt = timestamppb.New(addedAt)

		friendsList = append(friendsList, &friend)
	}

	return &friends.GetFriendListResponse{
		Friends: friendsList,
	}, nil
}

// DeclineFriendRequest handles declining a friend request.
func (s *FriendsServer) DeclineFriendRequest(ctx context.Context, req *friends.DeclineFriendRequestRequest) (*friends.DeclineFriendRequestResponse, error) {
	// Retrieve the user ID from the context
	userID, ok := ctx.Value("userID").(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("user ID not found in context")
	}

	// Convert userID from string to int
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	// Step 1: Update the friend request status to "DECLINED" if it exists and is pending
	res, err := s.DB.Exec(`UPDATE friend_requests SET status = ?, response_at = NOW() WHERE id = ? AND recipient_id = ? AND status = ?`,
		storage.StatusDeclinedStr, req.RequestId, userIDInt, storage.StatusPendingStr)

	if err != nil {
		return nil, fmt.Errorf("error updating friend request status to declined: %w", err)
	}

	// Step 2: Check if exactly one row was affected by the update
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("error checking affected rows: %w", err)
	}

	if rowsAffected != 1 {
		// No rows were affected, indicating that the request does not exist or is not pending
		return &friends.DeclineFriendRequestResponse{
			Status:    friends.FriendRequestStatus_FAILED,
			Message:   "Friend request does not exist or is not pending",
			Timestamp: timestamppb.Now(),
		}, nil
	}

	// Step 3: Return a successful response
	return &friends.DeclineFriendRequestResponse{
		Status:    friends.FriendRequestStatus_DECLINED,
		Message:   "Friend request declined successfully",
		Timestamp: timestamppb.Now(),
	}, nil
}

// RemoveFriend handles removing a friend.
func (s *FriendsServer) RemoveFriend(ctx context.Context, req *friends.RemoveFriendRequest) (*friends.RemoveFriendResponse, error) {
	// Retrieve the user ID from the context
	userID, ok := ctx.Value("userID").(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("user ID not found in context")
	}

	// Convert userID from string to int
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	// Step 1: Remove the friendship from the friends table
	res, err := s.DB.Exec(`DELETE FROM friends WHERE (user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)`,
		userIDInt, req.FriendId, req.FriendId, userIDInt)
	if err != nil {
		return nil, fmt.Errorf("error removing friend from friends table: %w", err)
	}

	// Step 2: Check if any rows were affected by the deletion
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("error checking affected rows: %w", err)
	}

	if rowsAffected == 0 {
		// No rows were affected, indicating that the friendship does not exist
		return &friends.RemoveFriendResponse{
			Success:   false,
			Message:   "Friend does not exist or has already been removed",
			Timestamp: timestamppb.Now(),
		}, nil
	}

	// Step 3: Update the friend request status to "CANCELLED" if it exists
	_, err = s.DB.Exec(`UPDATE friend_requests SET status = ? WHERE (requester_id = ? AND recipient_id = ?) OR (requester_id = ? AND recipient_id = ?)`,
		storage.StatusCancelledStr, userIDInt, req.FriendId, req.FriendId, userIDInt)
	if err != nil {
		return nil, fmt.Errorf("error updating friend request status to cancelled: %w", err)
	}

	// Step 4: Return a successful response
	return &friends.RemoveFriendResponse{
		Success:   true,
		Message:   "Friend removed successfully",
		Timestamp: timestamppb.Now(),
	}, nil
}
