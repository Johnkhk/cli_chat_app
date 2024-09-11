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
)

// NewDefaultTestServerConfig creates a TestServerConfig with default values
func NewDefaultTestServerConfig() *TestServerConfig {
	return &TestServerConfig{
		DbName:               "default_test_db",
		Log:                  logrus.New(),
		AccessTokenDuration:  time.Hour,          // Default to 1 hour
		RefreshTokenDuration: time.Hour * 24 * 7, // Default to 7 days
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

// Helper function to initialize test resources with a test-specific logger
func InitializeTestResources(t *testing.T, serverConfig *TestServerConfig) (*app.AuthClient, *sql.DB, func()) {
	// Use default server configuration if none provided
	if serverConfig == nil {
		serverConfig = NewDefaultTestServerConfig()
	}

	// Set up logger with test-specific details
	serverConfig.Log.SetReportCaller(true)
	serverConfig.Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,                  // Include full timestamp with date
		TimestampFormat: "2006-01-02 15:04:05", // Custom date format
		PadLevelText:    true,                  // Align log level text for better readability
		ForceColors:     true,                  // Force colors even if the output is not a terminal
	})

	// Add a custom hook to add the test name to every log entry
	serverConfig.Log.AddHook(&TestFieldHook{TestName: t.Name()})

	// Generate a unique path for JWT tokens using the test name
	tokenDir := filepath.Join(os.TempDir(), fmt.Sprintf(".test_cli_chat_app_%s", t.Name()))
	filePath := filepath.Join(tokenDir, "jwt_tokens")

	// Initialize the token manager with the unique path
	tokenManager := app.NewTokenManager(filePath, nil)

	// Generate a unique database name for each test if not set in serverConfig
	if serverConfig.DbName == "default_test_db" {
		serverConfig.DbName = fmt.Sprintf("test_db_%s", t.Name())
	}

	// Initialize Bufconn and test database with the unique database name
	conn, db := InitBufconn(t, *serverConfig)

	// Initialize the auth client with the shared connection
	authClient := &app.AuthClient{
		Client:       auth.NewAuthServiceClient(conn),
		Logger:       serverConfig.Log, // Use the test-specific logger
		TokenManager: tokenManager,
	}
	tokenManager.SetClient(authClient.Client)

	// Define a cleanup function to close the connection, teardown the database, and remove the token directory
	cleanup := func() {
		conn.Close()
		if err := TeardownTestDatabase(db, serverConfig.DbName); err != nil {
			serverConfig.Log.Errorf("Failed to teardown test database: %v", err)
		}

		// Remove the JWT token directory created by the token manager
		if err := os.RemoveAll(tokenDir); err != nil {
			serverConfig.Log.Errorf("Failed to remove JWT token directory: %v", err)
		}
	}

	// Return initialized resources and cleanup function
	return authClient, db, cleanup
}
