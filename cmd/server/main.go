package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	"github.com/johnkhk/cli_chat_app/genproto/auth"
	"github.com/johnkhk/cli_chat_app/server/app"
	"github.com/johnkhk/cli_chat_app/server/logger"
	"github.com/johnkhk/cli_chat_app/server/storage" // Import the storage package
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Check the port value
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("PORT environment variable is not set.")
	}

	// Initialize the logger
	logger.InitLogger()

	// Connect to the database
	storage.InitDB()

	// Create a new gRPC server
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(logger.UnaryInterceptor))
	authServer := &app.AuthServer{DB: storage.DB}
	auth.RegisterAuthServiceServer(grpcServer, authServer)

	// Listen specifically on localhost
	listener, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		logger.Log.Fatalf("Failed to listen on localhost:%s: %v", port, err)
	}

	logger.Log.Infof("gRPC server is listening on localhost:%s", port)

	if err := grpcServer.Serve(listener); err != nil {
		logger.Log.Fatalf("Failed to serve: %v", err)
	}
}
