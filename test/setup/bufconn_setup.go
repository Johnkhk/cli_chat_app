// test/setup/bufconn_setup.go
package setup

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	client "github.com/johnkhk/cli_chat_app/client/app"
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
	TimeProvider         client.TimeProvider
}

// InitTestServer initializes the in-memory gRPC server and the test database.
func InitTestServer(t *testing.T, serverConfig TestServerConfig) *sql.DB {

	lis = bufconn.Listen(BufSize)
	secretKey := os.Getenv("CLI_CHAT_APP_JWT_SECRET_KEY")
	if secretKey == "" {
		log.Fatal("JWT secret key is not set.")
	}
	tokenValidator := app.NewJWTTokenValidator(secretKey)

	// Create a new gRPC server with the authentication interceptor
	s := grpc.NewServer(
		grpc.UnaryInterceptor(app.UnaryServerInterceptor(tokenValidator, serverConfig.Log)),
	)

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
	friends.RegisterFriendManagementServer(s, friendsServer)

	// Start serving the in-memory server
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	return db
}

// CreateTestClientConn creates a new client connection using bufconn.
func CreateTestClientConn(t *testing.T, interceptor grpc.UnaryClientInterceptor) *grpc.ClientConn {
	// Create a client connection using bufconn
	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure(), grpc.WithUnaryInterceptor(interceptor))
	if err != nil {
		if t != nil {
			t.Fatalf("Failed to dial bufnet: %v", err)
		} else {
			log.Panicf("Failed to dial bufnet: %v", err)
		}
	}
	return conn
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}
