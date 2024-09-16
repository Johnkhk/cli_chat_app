package app

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/johnkhk/cli_chat_app/genproto/auth"
	"github.com/johnkhk/cli_chat_app/genproto/chat"
	"github.com/johnkhk/cli_chat_app/genproto/friends"
)

// RpcClient manages multiple gRPC clients for different services.
type RpcClient struct {
	AuthClient    *AuthClient
	FriendsClient *FriendsClient
	ChatClient    *ChatClient
	Conn          *grpc.ClientConn
	Logger        *logrus.Logger
	AppDirPath    string
}

type RpcClientConfig struct {
	Conn          *grpc.ClientConn
	ServerAddress string
	Logger        *logrus.Logger
	AppDirPath    string
	TokenManager  *TokenManager
}

// NewRpcClient initializes all service clients with a shared gRPC connection.
func NewRpcClient(config RpcClientConfig) (*RpcClient, error) {
	logger := config.Logger

	// Create the application directory if it doesn't exist
	dir := filepath.Dir(config.AppDirPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		logger.Errorf("Failed to create directory %s: %v", dir, err)
		return nil, err
	}

	// Create a TokenManager instance
	tokenManager := config.TokenManager
	if tokenManager == nil {
		tokenManager = NewTokenManager(filepath.Join(config.AppDirPath, "jwt_tokens"), nil) // Will set the client later
	}

	// Establish a single gRPC connection to the server
	conn := config.Conn
	if conn == nil {
		var err error
		conn, err = grpc.Dial(config.ServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(UnaryInterceptor(tokenManager, logger))) // Add the interceptor here

		if err != nil {
			logger.Errorf("Failed to connect to server: %v", err)
			return nil, err
		}
	}

	// Initialize individual clients with the shared connection
	authClient := &AuthClient{
		Client:       auth.NewAuthServiceClient(conn),
		Logger:       logger,
		TokenManager: tokenManager,
		AppDirPath:   config.AppDirPath,
	}
	tokenManager.SetClient(authClient.Client) // Set the client in the TokenManager

	friendsClient := &FriendsClient{
		Client: friends.NewFriendManagementClient(conn),
		Logger: logger,
	}

	chatClient := &ChatClient{
		Client: chat.NewChatServiceClient(conn),
		Logger: logger,
	}

	return &RpcClient{
		AuthClient:    authClient,
		FriendsClient: friendsClient,
		ChatClient:    chatClient,
		Conn:          conn,
		Logger:        logger,
	}, nil
}

// CloseConnections closes the shared gRPC connection.
func (r *RpcClient) CloseConnections() {
	if err := r.Conn.Close(); err != nil {
		r.Logger.Errorf("Failed to close the connection: %v", err)
	}
}

func (r *RpcClient) GetAppDirPath() string {
	return r.AppDirPath
}
