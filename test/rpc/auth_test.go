package rpc

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"google.golang.org/grpc"

	"github.com/johnkhk/cli_chat_app/client/app"
	"github.com/johnkhk/cli_chat_app/genproto/auth"
	"github.com/johnkhk/cli_chat_app/server/logger"
	"github.com/johnkhk/cli_chat_app/test/setup"
)

var (
	conn       *grpc.ClientConn
	authClient *app.AuthClient
)

func TestMain(m *testing.M) {
	// Initialize the logger
	log := logger.InitLogger()
	var db *sql.DB

	// Initialize Bufconn and test database
	conn, db = setup.InitBufconn(nil, log)
	defer conn.Close()
	defer func() {
		if err := setup.TeardownTestDatabase(db); err != nil {
			log.Errorf("Failed to teardown test database: %v", err)
		}
	}()

	// Initialize the token manager (replace with your actual token manager initialization)
	filePath := filepath.Join(os.Getenv("HOME"), ".test_cli_chat_app", "jwt_tokens") // For Linux/macOS
	tokenManager := app.NewTokenManager(filePath, nil)

	// Initialize the auth client with the shared connection
	authClient = &app.AuthClient{
		Client:       auth.NewAuthServiceClient(conn),
		Logger:       log,
		TokenManager: tokenManager,
	}
	tokenManager.SetClient(authClient.Client)

	// Run all tests
	code := m.Run()
	os.Exit(code)
}

// Test the new user registration
func TestNewUserRegister(t *testing.T) {

	// clear the directory
	authClient.Logger.Infof("Removing test directory")
	if err := os.RemoveAll(filepath.Join(os.Getenv("HOME"), ".test_cli_chat_app")); err != nil {
		t.Fatalf("Failed to remove test directory: %v", err)
	}

	// Test the new user registration
	authClient.Logger.Infof("Testing new user registration")
	err := authClient.RegisterUser("testuser", "testpassword")
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	// Make sure the directory was not created (directory gets created on login)
	authClient.Logger.Infof("Checking test directory")
	if _, err := os.Stat(filepath.Join(os.Getenv("HOME"), ".test_cli_chat_app")); err == nil {
		t.Fatalf("Test directory was created")
	}

}

// Test the login of an unregistered user
func TestUnregisteredUserLogin(t *testing.T) {
	authClient.Logger.Infof("Testing unregistered user login")

	// Attempt to login with an unregistered user
}
