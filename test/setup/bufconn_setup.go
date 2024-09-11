// test/setup/bufconn_setup.go
package setup

import (
	"context"
	"database/sql"
	"net"
	"testing"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"github.com/johnkhk/cli_chat_app/genproto/auth"
	"github.com/johnkhk/cli_chat_app/genproto/friends"
	"github.com/johnkhk/cli_chat_app/server/app"
)

const BufSize = 1024 * 1024

var lis *bufconn.Listener

// InitBufconn initializes the in-memory gRPC server for testing.
// we dont return the server, just the conn and db and that's good enough for the client to interact with
func InitBufconn(t *testing.T, log *logrus.Logger) (*grpc.ClientConn, *sql.DB) {
	// Load environment variables from .env file
	if err := godotenv.Load("../../.env"); err != nil {
		if t != nil {
			t.Fatalf("Error loading .env file: %v", err)
		} else {
			log.Panicf("Error loading .env file: %v", err)
		}
	}

	lis = bufconn.Listen(BufSize)
	s := grpc.NewServer()

	// Set up the database for testing
	db, err := SetupTestDatabase()
	if err != nil {
		if t != nil {
			t.Fatalf("Failed to set up test database: %v", err)
		} else {
			log.Panicf("Failed to set up test database: %v", err)
		}
	}

	// Initialize the servers with the test database
	authServer := app.NewAuthServer(db, log)
	auth.RegisterAuthServiceServer(s, authServer)

	friendsServer := app.NewFriendsServer(db, log)
	friends.RegisterFriendsServiceServer(s, friendsServer)

	// Start serving the in-memory server
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	// Create a client connection using bufconn
	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		if t != nil {
			t.Fatalf("Failed to dial bufnet: %v", err)
		} else {
			log.Panicf("Failed to dial bufnet: %v", err)
		}
	}
	return conn, db
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}
