package app

import (
	"context"
	"database/sql"
	"net"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/johnkhk/cli_chat_app/genproto/auth"
	"github.com/johnkhk/cli_chat_app/genproto/chat"
	"github.com/johnkhk/cli_chat_app/genproto/friends"
)

// RunGRPCServer initializes and runs the gRPC server.
func RunGRPCServer(ctx context.Context, port string, db *sql.DB, log *logrus.Logger) error {
	// Initialize the token validator
	secretKey := os.Getenv("CLI_CHAT_APP_JWT_SECRET_KEY")
	if secretKey == "" {
		log.Fatal("JWT secret key is not set in GRPC server.")
	}
	tokenValidator := NewJWTTokenValidator(secretKey)

	// Create a new gRPC server with the authentication interceptor
	grpcServer := SetupGRPCServer(tokenValidator, log)

	// Register the AuthServer
	authServer := NewAuthServer(db, log, time.Hour, time.Hour*24*7)
	auth.RegisterAuthServiceServer(grpcServer, authServer)

	// Register the FriendsServer
	friendsServer := NewFriendsServer(db, log)
	friends.RegisterFriendManagementServer(grpcServer, friendsServer)

	// Register the ChatServer
	chatServer := NewChatServiceServer(log)
	chat.RegisterChatServiceServer(grpcServer, chatServer)

	// Listen on the specified port
	listener, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		return err
	}

	log.Infof("gRPC server is listening on localhost:%s", port)

	// Start the server in a separate goroutine
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Errorf("Failed to serve gRPC server: %v", err)
		}
	}()

	// Wait for context cancellation for graceful shutdown
	<-ctx.Done()

	log.Info("Received shutdown signal, shutting down gRPC server...")

	// Gracefully stop the gRPC server
	grpcServer.GracefulStop()

	log.Info("gRPC server stopped.")
	return nil
}
