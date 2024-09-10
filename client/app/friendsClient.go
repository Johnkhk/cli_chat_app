package app

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/johnkhk/cli_chat_app/genproto/friends"
)

// FriendsClient encapsulates the gRPC client for friend services.
type FriendsClient struct {
	Client friends.FriendsServiceClient
	Logger *logrus.Logger
}

// SendFriendRequest sends a friend request to another user.
func (c *FriendsClient) SendFriendRequest(recipientUsername string) error {
	req := &friends.SendFriendRequestRequest{
		RecipientUsername: recipientUsername,
	}

	resp, err := c.Client.SendFriendRequest(context.Background(), req)
	if err != nil {
		c.Logger.Errorf("Failed to send friend request: %v", err)
		return fmt.Errorf("failed to send friend request: %w", err)
	}

	if resp.Success {
		c.Logger.Infof("Friend request sent successfully: %s", resp.Message)
	} else {
		c.Logger.Infof("Failed to send friend request: %s", resp.Message)
		return fmt.Errorf("failed to send friend request: %s", resp.Message)
	}

	return nil
}

// GetFriendRequests retrieves all pending friend requests for the current user.
func (c *FriendsClient) GetFriendRequests() (*friends.GetFriendRequestsResponse, error) {
	req := &friends.GetFriendRequestsRequest{}

	resp, err := c.Client.GetFriendRequests(context.Background(), req)
	if err != nil {
		c.Logger.Errorf("Failed to get friend requests: %v", err)
		return nil, fmt.Errorf("failed to get friend requests: %w", err)
	}

	c.Logger.Infof("Fetched friend requests: %v", resp.FriendRequests)
	return resp, nil
}

// RespondToFriendRequest responds to a pending friend request.
func (c *FriendsClient) RespondToFriendRequest(requesterID string, accept bool) error {
	req := &friends.RespondToFriendRequestRequest{
		RequesterId: requesterID,
		Accept:      accept,
	}

	resp, err := c.Client.RespondToFriendRequest(context.Background(), req)
	if err != nil {
		c.Logger.Errorf("Failed to respond to friend request: %v", err)
		return fmt.Errorf("failed to respond to friend request: %w", err)
	}

	if resp.Success {
		c.Logger.Infof("Responded to friend request successfully: %s", resp.Message)
	} else {
		c.Logger.Infof("Failed to respond to friend request: %s", resp.Message)
		return fmt.Errorf("failed to respond to friend request: %s", resp.Message)
	}

	return nil
}

// GetFriend retrieves details of a specific friend.
func (c *FriendsClient) GetFriend(friendID string) (*friends.Friend, error) {
	req := &friends.GetFriendRequest{
		FriendId: friendID,
	}

	resp, err := c.Client.GetFriend(context.Background(), req)
	if err != nil {
		c.Logger.Errorf("Failed to get friend details: %v", err)
		return nil, fmt.Errorf("failed to get friend details: %w", err)
	}

	if resp.Success {
		c.Logger.Infof("Fetched friend details: %v", resp.Friend)
		return resp.Friend, nil
	} else {
		c.Logger.Infof("Failed to fetch friend details: %s", resp.Message)
		return nil, fmt.Errorf("failed to fetch friend details: %s", resp.Message)
	}
}

// GetFriends retrieves all friends of the current user.
func (c *FriendsClient) GetFriends() ([]*friends.Friend, error) {
	req := &friends.GetFriendsRequest{}

	resp, err := c.Client.GetFriends(context.Background(), req)
	if err != nil {
		c.Logger.Errorf("Failed to get friends: %v", err)
		return nil, fmt.Errorf("failed to get friends: %w", err)
	}

	c.Logger.Infof("Fetched friends: %v", resp.Friends)
	return resp.Friends, nil
}

// RemoveFriend removes a friend from the current user's friend list.
func (c *FriendsClient) RemoveFriend(friendID string) error {
	req := &friends.RemoveFriendRequest{
		FriendId: friendID,
	}

	resp, err := c.Client.RemoveFriend(context.Background(), req)
	if err != nil {
		c.Logger.Errorf("Failed to remove friend: %v", err)
		return fmt.Errorf("failed to remove friend: %w", err)
	}

	if resp.Success {
		c.Logger.Infof("Friend removed successfully: %s", resp.Message)
	} else {
		c.Logger.Infof("Failed to remove friend: %s", resp.Message)
		return fmt.Errorf("failed to remove friend: %s", resp.Message)
	}

	return nil
}
