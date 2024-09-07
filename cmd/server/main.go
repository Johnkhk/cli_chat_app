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
	log := logger.InitLogger()

	// Connect to the database
	db, err := storage.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize the database: %v", err)
	}
	defer db.Close() // Ensure the database connection is closed when the application exits

	// Create a new gRPC server with a logging interceptor
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(logger.UnaryInterceptor(log)))
	authServer := app.NewAuthServer(db, log)
	auth.RegisterAuthServiceServer(grpcServer, authServer)

	// Listen specifically on localhost
	listener, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		log.Fatalf("Failed to listen on localhost:%s: %v", port, err)
	}

	log.Infof("gRPC server is listening on localhost:%s", port)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
