package storage

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql" // Import MySQL driver
)

// InitDB initializes the database connection with retry logic.
func InitDB() (*sql.DB, error) {
	// Get the database URL from environment variables
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	// Retry logic
	maxRetries := 10
	var db *sql.DB
	var err error
	for i := 0; i < maxRetries; i++ {
		// Open the database connection
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			return nil, fmt.Errorf("failed to open database connection: %w", err)
		}

		// Test the database connection
		if err = db.Ping(); err == nil {
			// Connection successful, break the retry loop
			fmt.Println("Successfully connected to the database.")
			return db, nil
		}

		// Close the database connection on error before retrying
		db.Close()
		fmt.Printf("Failed to connect to database (attempt %d/%d): %v\n", i+1, maxRetries, err)
		time.Sleep(2 * time.Second) // Wait before retrying
	}

	// If all retries fail, return the last error
	return nil, fmt.Errorf("could not connect to the database after %d attempts: %w", maxRetries, err)
}
