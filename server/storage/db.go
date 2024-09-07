package storage

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql" // Import MySQL driver
)

// InitDB initializes the database connection.
func InitDB() (*sql.DB, error) {
	// Get the database URL from environment variables
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	// Open the database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the database connection
	if err := db.Ping(); err != nil {
		db.Close() // Close the database connection on error
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
