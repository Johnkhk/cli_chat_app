package app

import (
	"context"
	"database/sql"

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

// AddFriend handles requests to add a new friend.
func (s *FriendsServer) AddFriend(ctx context.Context, req *friends.AddFriendRequest) (*friends.AddFriendResponse, error) {
	s.Logger.Infof("Adding friend: %s", req.Username)

	// Mock AddFriend logic...
	return &friends.AddFriendResponse{Success: true, Message: "Friend added successfully"}, nil
}

// GetRecentMessages handles requests to fetch recent messages.
func (s *FriendsServer) GetRecentMessages(ctx context.Context, req *friends.RecentMessagesRequest) (*friends.RecentMessagesResponse, error) {
	s.Logger.Infof("Fetching recent messages for user_id: %s, friend_id: %s", req.UserId, req.FriendId)

	// Mock data for recent messages
	messages := []*friends.Message{
		{
			Id:         1,
			SenderId:   req.UserId,
			ReceiverId: req.FriendId,
			Content:    "Hello! This is a mock message.",
			CreatedAt:  "2024-09-09T10:00:00Z",
		},
		{
			Id:         2,
			SenderId:   req.FriendId,
			ReceiverId: req.UserId,
			Content:    "Hi! Another mock message.",
			CreatedAt:  "2024-09-09T10:05:00Z",
		},
	}

	return &friends.RecentMessagesResponse{
		Messages: messages,
		Success:  true,
		Message:  "Fetched recent messages successfully",
	}, nil
}
