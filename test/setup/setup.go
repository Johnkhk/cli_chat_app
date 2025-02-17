package setup

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/johnkhk/cli_chat_app/client/app"
)

// NewDefaultTestServerConfig creates a TestServerConfig with default values
func NewDefaultTestServerConfig(t *testing.T) (*TestServerConfig, error) {
	// Load environment variables from .env file
	if _, err := os.Stat("../../.env"); err == nil {
		if err := godotenv.Load("../../.env"); err != nil {
			if t != nil {
				t.Fatalf("Error loading .env file: %v", err)
			} else {
				t.Fatalf("Error loading .env file: %v", err)
			}
		}
	}
	// Get the log directory from the environment variable
	logDir := os.Getenv("TEST_LOG_DIR")
	if logDir == "" {
		// return nil, fmt.Errorf("environment variable TEST_LOG_DIR is not set")
		appDirPath, err := app.GetAppDirPath()
		if err != nil {
			t.Fatalf("Failed to get app directory path: %v", err)
		}
		logDir = filepath.Join(appDirPath, "test_logs")
	}

	// Ensure the directory exists
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %v", err)
	}

	// Construct the log file name using the test name and the current date
	currentDateTime := time.Now().Format("2006-01-02_15-04-05")
	logFileName := fmt.Sprintf("%s_%s.log", t.Name(), currentDateTime)
	logFilePath := filepath.Join(logDir, logFileName)
	t.Logf("Log file path created: %s", logFilePath)

	// Open or create the log file
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}

	// Create a new logger instance
	logger := logrus.New()
	logger.SetOutput(logFile) // Set the logger to write to the file
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		PadLevelText:    true,
		ForceColors:     false, // Disable colors for file output
	})
	logger.SetLevel(logrus.InfoLevel) // Ensure log level is set to INFO or lower

	return &TestServerConfig{
		DbName:               "default_test_db",
		Log:                  logger,
		AccessTokenDuration:  time.Hour,          // Default to 1 hour
		RefreshTokenDuration: time.Hour * 24 * 7, // Default to 7 days
		TimeProvider:         app.RealTimeProvider{},
	}, nil
}

// TestFieldHook is a custom Logrus hook that adds a field to every log entry
type TestFieldHook struct {
	TestName string
}

// Fire adds the test name field to every log entry
func (hook *TestFieldHook) Fire(entry *logrus.Entry) error {
	entry.Message = "[" + hook.TestName + "] " + entry.Message
	return nil
}

// Levels returns the log levels the hook should be applied to
func (hook *TestFieldHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// InitializeTestResources initializes the desired number of RpcClients for testing.
func InitializeTestResources(t *testing.T, serverConfig *TestServerConfig, numClients int) ([]*app.RpcClient, *sql.DB, func(), *ServerStruct) {

	// Use default server configuration if none provided
	if serverConfig == nil {
		var err error
		serverConfig, err = NewDefaultTestServerConfig(t)
		if err != nil {
			t.Fatalf("Failed to create test server config: %v", err)
		}
	}

	// Add test-specific details to the logger
	serverConfig.Log.AddHook(&TestFieldHook{TestName: t.Name()})

	// Generate a unique database name for each test if not set in serverConfig
	if serverConfig.DbName == "default_test_db" {
		hashedName := fmt.Sprintf("%x", md5.Sum([]byte(t.Name()))) // Create a hash of the test name
		serverConfig.DbName = fmt.Sprintf("db_%s_%d", hashedName, time.Now().UnixNano())

	}

	// Initialize the in-memory gRPC server and test database
	db, server := InitTestServer(t, *serverConfig)

	// Slice to hold the RpcClients
	var rpcClients []*app.RpcClient

	// Loop to initialize the desired number of RpcClients
	for i := 0; i < numClients; i++ {
		// Create a unique directory for each client
		appDir := filepath.Join(os.TempDir(), fmt.Sprintf(".test_cli_chat_app_%s_client_%d", t.Name(), i))
		tokenDir := filepath.Join(appDir, "jwt_tokens")

		// Create the TokenManager for the client
		tokenManager := app.NewTokenManager(tokenDir, nil)
		tokenManager.TimeProvider = serverConfig.TimeProvider

		conn := CreateTestClientConn(t, app.UnaryInterceptor(tokenManager, serverConfig.Log), app.StreamInterceptor(tokenManager, serverConfig.Log))

		// Include custom RPC client configuration
		rpcClientConfig := app.RpcClientConfig{
			ServerAddress: "", // Not needed since we pass the connection directly
			Logger:        serverConfig.Log,
			AppDirPath:    appDir,
			Conn:          conn,
			TokenManager:  tokenManager,
		}
		// Initialize the gRPC client using RpcClient
		rpcClient, err := app.NewRpcClient(rpcClientConfig)
		tokenManager.SetClient(rpcClient.AuthClient)
		if err != nil {
			t.Fatalf("Failed to initialize RPC clients: %v", err)
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
	return rpcClients, db, cleanup, server
}
