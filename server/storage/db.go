package storage

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql" // Import MySQL driver

	"github.com/johnkhk/cli_chat_app/server/logger"
)

var DB *sql.DB

// InitDB initializes the database connection.
func InitDB() {
	// Get the database URL from environment variables
	dsn := os.Getenv("DATABASE_URL")
	logger.Log.Debug("Database URL:", dsn)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		logger.Log.Fatalf("Failed to open database connection: %v", err)
	}

	// // Defer closing the database connection until the program exits
	// // This will ensure the connection is properly closed
	// // when the main function returns
	// defer func() {
	// 	if err := DB.Close(); err != nil {
	// 		logger.Log.Fatalf("Failed to close database connection: %v", err)
	// 	}
	// }()

	// Test the database connection
	if err := DB.Ping(); err != nil {
		logger.Log.Fatalf("Failed to ping database: %v", err)
	}

	logger.Log.Info("Database connection established successfully")
}
