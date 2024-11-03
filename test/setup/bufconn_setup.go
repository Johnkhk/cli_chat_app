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
	"github.com/johnkhk/cli_chat_app/genproto/chat"
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

type ServerStruct struct {
	AuthServer    *app.AuthServer
	FriendsServer *app.FriendsServer
	ChatServer    *app.ChatServiceServer
}

// InitTestServer initializes the in-memory gRPC server and the test database.
func InitTestServer(t *testing.T, serverConfig TestServerConfig) (*sql.DB, *ServerStruct) {

	lis = bufconn.Listen(BufSize)
	secretKey := os.Getenv("CLI_CHAT_APP_JWT_SECRET_KEY")
	if secretKey == "" {
		log.Fatal("JWT secret key is not set.")
	}
	tokenValidator := app.NewJWTTokenValidator(secretKey)

	// Create a new gRPC server with the authentication interceptor
	// s := grpc.NewServer(
	// 	grpc.UnaryInterceptor(app.UnaryServerInterceptor(tokenValidator, serverConfig.Log)),
	// )
	s := app.SetupGRPCServer(tokenValidator, serverConfig.Log)

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

	chatServer := app.NewChatServiceServer(serverConfig.Log)
	chat.RegisterChatServiceServer(s, chatServer)

	// serverStruct
	serverStruct := &ServerStruct{
		AuthServer:    authServer,
		FriendsServer: friendsServer,
		ChatServer:    chatServer,
	}

	// Start serving the in-memory server
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	return db, serverStruct
}

// CreateTestClientConn creates a new client connection using bufconn.
func CreateTestClientConn(t *testing.T, unaryInterceptor grpc.UnaryClientInterceptor, streamInterceptor grpc.StreamClientInterceptor) *grpc.ClientConn {
	// Create a client connection using bufconn
	conn, err := grpc.DialContext(
		context.Background(),
		"bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(unaryInterceptor),   // Add the unary interceptor
		grpc.WithStreamInterceptor(streamInterceptor), // Add the stream interceptor
	)
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
