package app

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/johnkhk/cli_chat_app/client/e2ee/store"
	"github.com/johnkhk/cli_chat_app/genproto/auth"
	"github.com/johnkhk/cli_chat_app/genproto/chat"
	"github.com/johnkhk/cli_chat_app/genproto/friends"
)

// RpcClient manages multiple gRPC clients for different services.
type RpcClient struct {
	AuthClient      *AuthClient
	FriendsClient   *FriendsClient
	ChatClient      *ChatClient
	Conn            *grpc.ClientConn
	Logger          *logrus.Logger
	AppDirPath      string
	Store           *store.SQLiteStore
	CurrentUserID   uint32
	CurrentDeviceID uint32
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
	config.Logger.Infof("Creating App directory at: %s", dir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		logger.Errorf("Failed to create directory %s: %v", dir, err)
		return nil, err
	}

	// Create a new Store instance
	sqlitePath := filepath.Join(config.AppDirPath, "store.db")
	logger.Info("Creating SQLite store at: ", sqlitePath)
	sqliteStore, err := store.NewSQLiteStore(sqlitePath)
	if err != nil {
		logger.Errorf("Failed to create SQLite store: %v", err)
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
		// Your unary and stream interceptor functions
		unaryInterceptor := UnaryInterceptor(tokenManager, logger)
		streamInterceptor := StreamInterceptor(tokenManager, logger)
		conn, err = grpc.Dial(
			config.ServerAddress,
			grpc.WithTransportCredentials(insecure.NewCredentials()), // Using insecure credentials for local development/testing
			grpc.WithChainUnaryInterceptor(unaryInterceptor),         // Add the unary interceptor
			grpc.WithChainStreamInterceptor(streamInterceptor),       // Add the stream interceptor
		)
		if err != nil {
			logger.Errorf("Failed to connect to server: %v", err)
			return nil, err
		}
	}

	// Create the RpcClient instance
	rpcClient := &RpcClient{
		Logger: logger,
		Conn:   conn,
		Store:  sqliteStore,
	}

	// Initialize individual clients and set the ParentClient
	authClient := &AuthClient{
		Client:       auth.NewAuthServiceClient(conn),
		Logger:       logger,
		TokenManager: tokenManager,
		AppDirPath:   config.AppDirPath,
		SqliteStore:  sqliteStore,
		ParentClient: rpcClient, // Set reference to the parent RpcClient
	}

	chatClient := &ChatClient{
		Client:         chat.NewChatServiceClient(conn),
		AuthClient:     authClient,
		Store:          sqliteStore,
		Logger:         logger,
		MessageChannel: make(chan *chat.MessageResponse, 10), // Initialize the channel with a buffer size of 10
	}

	friendsClient := &FriendsClient{
		Client: friends.NewFriendManagementClient(conn),
		Logger: logger,
	}

	// Set clients in RpcClient
	rpcClient.AuthClient = authClient
	rpcClient.ChatClient = chatClient
	rpcClient.FriendsClient = friendsClient

	// Set the AuthService client in the TokenManager
	tokenManager.SetClient(authClient)

	return rpcClient, nil
}

// CloseConnections closes the shared gRPC connection.
func (r *RpcClient) CloseConnections() {
	if err := r.Conn.Close(); err != nil {
		r.Logger.Errorf("Failed to close the connection: %v", err)
	}

	if err := r.AuthClient.LogoutUser(); err != nil {
		r.Logger.Errorf("Failed to log out user: %v", err)
	}
}

func (r *RpcClient) GetAppDirPath() string {
	return r.AppDirPath
}
