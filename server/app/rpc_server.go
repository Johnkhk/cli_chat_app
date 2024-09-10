package app

import (
	"database/sql"
	"net"
	"os"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/johnkhk/cli_chat_app/genproto/auth"
	"github.com/johnkhk/cli_chat_app/genproto/friends"
)

// RunGRPCServer initializes and runs the gRPC server.
func RunGRPCServer(port string, db *sql.DB, log *logrus.Logger) error {
	// Initialize the token validator
	secretKey := os.Getenv("CLI_CHAT_APP_JWT_SECRET_KEY")
	if secretKey == "" {
		log.Fatal("JWT secret key is not set.")
	}
	tokenValidator := NewJWTTokenValidator(secretKey)

	// Create a new gRPC server with the authentication interceptor
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(UnaryServerInterceptor(tokenValidator, log)),
	)

	// Register the AuthServer
	authServer := NewAuthServer(db, log)
	auth.RegisterAuthServiceServer(grpcServer, authServer)

	// Register the FriendsServer
	friendsServer := NewFriendsServer(db, log)
	friends.RegisterFriendsServiceServer(grpcServer, friendsServer)

	// Listen on the specified port
	listener, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		return err
	}

	log.Infof("gRPC server is listening on localhost:%s", port)

	// Start serving
	return grpcServer.Serve(listener)
}
