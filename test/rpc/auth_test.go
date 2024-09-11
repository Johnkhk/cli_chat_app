package rpc

import (
	"database/sql"
	"os"
	"testing"

	"google.golang.org/grpc"

	"github.com/johnkhk/cli_chat_app/genproto/auth"
	"github.com/johnkhk/cli_chat_app/server/logger"
	"github.com/johnkhk/cli_chat_app/test/setup"
)

var (
	conn   *grpc.ClientConn
	client auth.AuthServiceClient
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

	// Initialize the gRPC client
	client = auth.NewAuthServiceClient(conn)

	// Run all tests
	code := m.Run()
	os.Exit(code)
}

func TestAuthFunctionality(t *testing.T) {
	// Use the initialized gRPC client to make a call
	// req := &auth.LoginRequest{Username: "testuser", Password: "testpassword"}

	// // Call the Login method on the client
	// resp, err := client.Login(context.Background(), req)
	// if err != nil {
	// 	t.Fatalf("Login failed: %v", err)
	// }

	// // Check the response
	// if resp.GetToken() == "" {
	// 	t.Errorf("Expected a token, got an empty string")
	// }
}
