package setup

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/johnkhk/cli_chat_app/client/app"
	"github.com/johnkhk/cli_chat_app/genproto/auth"
	"github.com/johnkhk/cli_chat_app/genproto/friends"
)

// NewDefaultTestServerConfig creates a TestServerConfig with default values
func NewDefaultTestServerConfig() *TestServerConfig {
	return &TestServerConfig{
		DbName:               "default_test_db",
		Log:                  logrus.New(),
		AccessTokenDuration:  time.Hour,          // Default to 1 hour
		RefreshTokenDuration: time.Hour * 24 * 7, // Default to 7 days
		TimeProvider:         app.RealTimeProvider{},
	}
}

// TestFieldHook is a custom Logrus hook that adds a field to every log entry
type TestFieldHook struct {
	TestName string
}

// Fire adds the test name field to every log entry
func (hook *TestFieldHook) Fire(entry *logrus.Entry) error {
	entry.Data["test"] = hook.TestName
	return nil
}

// Levels returns the log levels the hook should be applied to
func (hook *TestFieldHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// InitializeTestResources initializes the desired number of RpcClients for testing.
func InitializeTestResources(t *testing.T, serverConfig *TestServerConfig, numClients int) ([]*app.RpcClient, *sql.DB, func()) {
	// Use default server configuration if none provided
	if serverConfig == nil {
		serverConfig = NewDefaultTestServerConfig()
	}

	// Set up logger with test-specific details
	serverConfig.Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		PadLevelText:    true,
		ForceColors:     true,
	})
	serverConfig.Log.AddHook(&TestFieldHook{TestName: t.Name()})

	// Generate a unique database name for each test if not set in serverConfig
	if serverConfig.DbName == "default_test_db" {
		serverConfig.DbName = fmt.Sprintf("test_db_%s", t.Name())
	}

	// Initialize the in-memory gRPC server and test database
	db := InitTestServer(t, *serverConfig)

	// Slice to hold the RpcClients
	var rpcClients []*app.RpcClient

	// Loop to initialize the desired number of RpcClients
	for i := 0; i < numClients; i++ {

		// Create a unique path for JWT tokens for each client
		tokenDir := filepath.Join(os.TempDir(), fmt.Sprintf(".test_cli_chat_app_%s_client_%d", t.Name(), i))
		filePath := filepath.Join(tokenDir, "jwt_tokens")

		// Initialize the token manager with the unique path
		tokenManager := app.NewTokenManager(filePath, nil)
		tokenManager.TimeProvider = serverConfig.TimeProvider

		// Each client establishes its own connection to the server
		conn := CreateTestClientConn(t, app.UnaryInterceptor(tokenManager, serverConfig.Log))
		// Initialize the auth client
		authClient := &app.AuthClient{
			Client:       auth.NewAuthServiceClient(conn),
			Logger:       serverConfig.Log,
			TokenManager: tokenManager,
		}
		tokenManager.SetClient(authClient.Client)

		// Initialize the friends client
		friendsClient := &app.FriendsClient{
			Client: friends.NewFriendManagementClient(conn),
			Logger: serverConfig.Log,
		}

		// Create the RpcClient and add it to the slice
		rpcClient := &app.RpcClient{
			AuthClient:    authClient,
			FriendsClient: friendsClient,
			Conn:          conn,
			Logger:        serverConfig.Log,
		}
		rpcClients = append(rpcClients, rpcClient)
	}

	// Define a cleanup function to close all client connections, teardown the database, and remove token directories
	cleanup := func() {
		for _, rpcClient := range rpcClients {
			rpcClient.Conn.Close() // Close each client connection
		}
		if err := TeardownTestDatabase(db, serverConfig.DbName); err != nil {
			serverConfig.Log.Errorf("Failed to teardown test database: %v", err)
		}

		// Remove the JWT token directories created for each client
		for i := 0; i < numClients; i++ {
			tokenDir := filepath.Join(os.TempDir(), fmt.Sprintf(".test_cli_chat_app_%s_client_%d", t.Name(), i))
			if err := os.RemoveAll(tokenDir); err != nil {
				serverConfig.Log.Errorf("Failed to remove JWT token directory for client %d: %v", i, err)
			}
		}
	}

	// Return the initialized RpcClients, the database handle, and the cleanup function
	return rpcClients, db, cleanup
}
