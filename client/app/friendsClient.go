package app

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/johnkhk/cli_chat_app/genproto/friends"
)

// FriendsClient encapsulates the gRPC client for friend services.
type FriendsClient struct {
	Client friends.FriendManagementClient
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

	if resp.Status == friends.FriendRequestStatus_PENDING {
		c.Logger.Infof("Friend request sent successfully: %s", resp.Message)
	} else {
		c.Logger.Infof("Failed to send friend request: %s", resp.Message)
		return fmt.Errorf("failed to send friend request: %s", resp.Message)
	}

	return nil
}

// GetFriendList retrieves the list of friends for the current user.
func (c *FriendsClient) GetFriendList() ([]*friends.Friend, error) {
	req := &friends.GetFriendListRequest{}

	resp, err := c.Client.GetFriendList(context.Background(), req)
	if err != nil {
		c.Logger.Errorf("Failed to get friend list: %v", err)
		return nil, fmt.Errorf("failed to get friend list: %w", err)
	}

	c.Logger.Infof("Retrieved %d friends", len(resp.Friends))
	return resp.Friends, nil
}

// GetIncomingFriendRequests retrieves the incoming friend requests for the current user.
func (c *FriendsClient) GetIncomingFriendRequests() ([]*friends.FriendRequest, error) {
	req := &friends.GetIncomingFriendRequestsRequest{}

	resp, err := c.Client.GetIncomingFriendRequests(context.Background(), req)
	if err != nil {
		c.Logger.Errorf("Failed to get incoming friend requests: %v", err)
		return nil, fmt.Errorf("failed to get incoming friend requests: %w", err)
	}

	c.Logger.Infof("Retrieved %d incoming friend requests", len(resp.IncomingRequests))
	return resp.IncomingRequests, nil
}

// GetOutgoingFriendRequests retrieves the outgoing friend requests sent by the current user.
func (c *FriendsClient) GetOutgoingFriendRequests() ([]*friends.FriendRequest, error) {
	req := &friends.GetOutgoingFriendRequestsRequest{}

	resp, err := c.Client.GetOutgoingFriendRequests(context.Background(), req)
	if err != nil {
		c.Logger.Errorf("Failed to get outgoing friend requests: %v", err)
		return nil, fmt.Errorf("failed to get outgoing friend requests: %w", err)
	}

	c.Logger.Infof("Retrieved %d outgoing friend requests", len(resp.OutgoingRequests))
	return resp.OutgoingRequests, nil
}

// AcceptFriendRequest accepts an incoming friend request.
func (c *FriendsClient) AcceptFriendRequest(requestID string) error {
	req := &friends.AcceptFriendRequestRequest{
		RequestId: requestID,
	}

	resp, err := c.Client.AcceptFriendRequest(context.Background(), req)
	if err != nil {
		c.Logger.Errorf("Failed to accept friend request: %v", err)
		return fmt.Errorf("failed to accept friend request: %w", err)
	}

	if resp.Status == friends.FriendRequestStatus_ACCEPTED {
		c.Logger.Infof("Friend request accepted successfully: %s", resp.Message)
	} else {
		c.Logger.Infof("Failed to accept friend request: %s", resp.Message)
		return fmt.Errorf("failed to accept friend request: %s", resp.Message)
	}

	return nil
}
