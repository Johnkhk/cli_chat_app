package app

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/johnkhk/cli_chat_app/genproto/friends"
)

// FriendsClient encapsulates the gRPC client for friend services.
type FriendsClient struct {
	Client     friends.FriendsServiceClient
	Connection *grpc.ClientConn
	Logger     *logrus.Logger
}

// AddFriend sends a request to add a friend by username.
func (c *FriendsClient) AddFriend(username string) error {
	req := &friends.AddFriendRequest{
		Username: username,
	}

	resp, err := c.Client.AddFriend(context.Background(), req)
	if err != nil {
		c.Logger.Errorf("Failed to add friend: %v", err)
		return fmt.Errorf("Failed to add friend: %v", err)
	}

	if resp.Success {
		c.Logger.Infof("Friend added successfully: %s", resp.Message)
		return nil
	} else {
		c.Logger.Infof("Failed to add friend: %s", resp.Message)
		return fmt.Errorf("Failed to add friend: %s", resp.Message)
	}
}
