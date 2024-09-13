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

	// If a record is found, use the map to get the enum value
	statusEnum, ok := storage.StatusMap[existingStatus]
	if !ok {
		// If the status does not match any known value, use UNKNOWN
		statusEnum = friends.FriendRequestStatus_UNKNOWN
	}

	// Handle the relevant existing statuses
	switch statusEnum {
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
	case friends.FriendRequestStatus_DECLINED:
		// Allow sending the friend request again if it was previously rejected
		_, err = s.DB.Exec(`
			UPDATE friend_requests
			SET status = ?, requested_at = NOW()
			WHERE requester_id = ? AND recipient_id = ? AND status = ?`,
			storage.StatusPending, requesterIDInt, recipientID, storage.StatusDeclined)
		if err != nil {
			return nil, fmt.Errorf("error updating friend request to pending: %w", err)
		}
		return &friends.SendFriendRequestResponse{
			Status:    friends.FriendRequestStatus_PENDING,
			Message:   "Friend request sent again successfully",
			Timestamp: timestamppb.Now(),
		}, nil
	}

	// Step 3: Insert the new friend request into the database
	_, err = s.DB.Exec(`
		INSERT INTO friend_requests (requester_id, recipient_id, requested_at, status)
		VALUES (?, ?, NOW(), ?)`,
		requesterIDInt, recipientID, storage.StatusPending)
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
	res, err := s.DB.Exec(`UPDATE friend_requests SET status = ?, updated_at = NOW() WHERE id = ? AND recipient_id = ? AND status = ?`,
		storage.StatusAccepted, req.RequestId, userIDInt, storage.StatusPending)

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

	// Query to get all incoming friend requests for this user
	rows, err := s.DB.Query(`
        SELECT id, requester_id, recipient_id, status, created_at
        FROM friend_requests
        WHERE recipient_id = ? AND status = ?`, userIDInt, storage.StatusPending)
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

		if err := rows.Scan(&friendReq.RequestId, &friendReq.SenderId, &friendReq.RecipientId, &status, &createdAt); err != nil {
			return nil, fmt.Errorf("error scanning friend request row: %w", err)
		}

		// Convert status to enum and time to protobuf timestamp
		friendReq.Status = friends.FriendRequestStatus(friends.FriendRequestStatus_value[status])
		friendReq.CreatedAt = timestamppb.New(createdAt)

		incomingRequests = append(incomingRequests, &friendReq)
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

	// Query to get all outgoing friend requests for this user
	rows, err := s.DB.Query(`
        SELECT id, requester_id, recipient_id, status, created_at
        FROM friend_requests
        WHERE requester_id = ? AND status = ?`, userIDInt, storage.StatusPending)
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

		if err := rows.Scan(&friendReq.RequestId, &friendReq.SenderId, &friendReq.RecipientId, &status, &createdAt); err != nil {
			return nil, fmt.Errorf("error scanning friend request row: %w", err)
		}

		// Convert status to enum and time to protobuf timestamp
		friendReq.Status = friends.FriendRequestStatus(friends.FriendRequestStatus_value[status])
		friendReq.CreatedAt = timestamppb.New(createdAt)

		outgoingRequests = append(outgoingRequests, &friendReq)
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
        SELECT f.friend_id, u.username, u.status, u.last_active_at, f.created_at
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
		var lastActiveAt, addedAt time.Time

		if err := rows.Scan(&friend.UserId, &friend.Username, &friend.Status, &lastActiveAt, &addedAt); err != nil {
			return nil, fmt.Errorf("error scanning friend row: %w", err)
		}

		// Convert times to protobuf timestamps
		friend.LastActiveAt = timestamppb.New(lastActiveAt)
		friend.AddedAt = timestamppb.New(addedAt)

		friendsList = append(friendsList, &friend)
	}

	return &friends.GetFriendListResponse{
		Friends: friendsList,
	}, nil
}
