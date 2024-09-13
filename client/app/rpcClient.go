package app

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/johnkhk/cli_chat_app/genproto/auth"
	"github.com/johnkhk/cli_chat_app/genproto/friends"
)

// RpcClient manages multiple gRPC clients for different services.
type RpcClient struct {
	AuthClient    *AuthClient
	FriendsClient *FriendsClient
	conn          *grpc.ClientConn
	Logger        *logrus.Logger
}

// NewRpcClient initializes all service clients with a shared gRPC connection.
func NewRpcClient(serverAddress string, logger *logrus.Logger, tokenManager *TokenManager) (*RpcClient, error) {
	// Establish a single gRPC connection to the server
	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(UnaryInterceptor(tokenManager))) // Add the interceptor here

	if err != nil {
		logger.Errorf("Failed to connect to server: %v", err)
		return nil, err
	}

	// Initialize individual clients with the shared connection
	authClient := &AuthClient{
		Client:       auth.NewAuthServiceClient(conn),
		Logger:       logger,
		TokenManager: tokenManager,
	}

	friendsClient := &FriendsClient{
		Client: friends.NewFriendManagementClient(conn),
		Logger: logger,
	}

	return &RpcClient{
		AuthClient:    authClient,
		FriendsClient: friendsClient,
		conn:          conn,
		Logger:        logger,
	}, nil
}

// CloseConnections closes the shared gRPC connection.
func (r *RpcClient) CloseConnections() {
	if err := r.conn.Close(); err != nil {
		r.Logger.Errorf("Failed to close the connection: %v", err)
	}
}
