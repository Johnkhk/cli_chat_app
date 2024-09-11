// test/setup/bufconn_setup.go
package setup

import (
	"context"
	"database/sql"
	"log"
	"net"
	"testing"
	"time"

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

type TestServerConfig struct {
	DbName               string
	Log                  *logrus.Logger
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

// InitBufconn initializes the in-memory gRPC server for testing.
// we dont return the server, just the conn and db and that's good enough for the client to interact with
func InitBufconn(t *testing.T, serverConfig TestServerConfig) (*grpc.ClientConn, *sql.DB) {
	// Load environment variables from .env file
	if err := godotenv.Load("../../.env"); err != nil {
		if t != nil {
			t.Fatalf("Error loading .env file: %v", err)
		} else {
			serverConfig.Log.Panicf("Error loading .env file: %v", err)
		}
	}

	lis = bufconn.Listen(BufSize)
	s := grpc.NewServer()

	// Set up the database for testing
	db, err := SetupTestDatabase(serverConfig.DbName)
	if err != nil {
		if t != nil {
			t.Fatalf("Failed to set up test database: %v", err)
		} else {
			serverConfig.Log.Panicf("Failed to set up test database: %v", err)
		}
	}

	// Initialize the servers with the test database
	authServer := app.NewAuthServer(db, serverConfig.Log, serverConfig.AccessTokenDuration, serverConfig.RefreshTokenDuration)
	auth.RegisterAuthServiceServer(s, authServer)

	friendsServer := app.NewFriendsServer(db, serverConfig.Log)
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
