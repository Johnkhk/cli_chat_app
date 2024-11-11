package setup

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// SetupTestDatabase initializes the test database.
func SetupTestDatabase(testDBName string) (*sql.DB, error) {
	// Get the Data Source Name (DSN) from environment variables
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("TEST_DATABASE_URL is not set")
	}

	// Connect to the MySQL server (without specifying a database initially)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to MySQL server: %v", err)
	}

	// Check MySQL version (optional but useful for debugging)
	checkMySQLVersion(db)

	// Create the test database if it does not exist
	if _, err := db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", testDBName)); err != nil {
		db.Close()
		return nil, fmt.Errorf("Failed to create test database: %v", err)
	}

	// Close the initial connection
	db.Close()

	// Reconnect to the MySQL server with the specific test database
	testDSN := fmt.Sprintf("%s%s?parseTime=true", dsn, testDBName)

	db, err = sql.Open("mysql", testDSN)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to the test database: %v", err)
	}

	// Drop all tables in the test database
	if err := dropAllTables(db, testDBName); err != nil {
		db.Close()
		return nil, fmt.Errorf("Failed to drop all tables: %v", err)
	}

	// Execute the SQL file to set up tables and other structures
	upSQLPath := getAbsolutePath("db/migrations/02_up.sql")
	if err := runSQLFile(db, upSQLPath); err != nil {
		db.Close()
		return nil, fmt.Errorf("Failed to execute setup SQL file: %v", err)
	}

	return db, nil
}

// TeardownTestDatabase cleans up the test database.
func TeardownTestDatabase(db *sql.DB, testDBName string) error {
	// Optionally run the teardown SQL script if needed
	downSQLPath := getAbsolutePath("db/migrations/01_down.sql")
	if err := runSQLFile(db, downSQLPath); err != nil {
		return fmt.Errorf("Failed to run teardown SQL file: %v", err)
	}

	// Drop the entire test database to ensure a clean state
	if _, err := db.Exec("DROP DATABASE IF EXISTS " + testDBName); err != nil {
		return fmt.Errorf("Failed to drop test database: %v", err)
	}

	// Close the database connection
	if db != nil {
		return db.Close()
	}
	return nil
}

// dropAllTables drops all tables in the current database.
func dropAllTables(db *sql.DB, dbName string) error {
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		return fmt.Errorf("Failed to list tables: %v", err)
	}
	defer rows.Close()

	var tableName string
	for rows.Next() {
		if err := rows.Scan(&tableName); err != nil {
			return fmt.Errorf("Failed to scan table name: %v", err)
		}

		// Explicitly use the database name in the DROP statement
		if _, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s.%s", dbName, tableName)); err != nil {
			return fmt.Errorf("Failed to drop table %s: %v", tableName, err)
		}
	}
	return nil
}

// runSQLFile executes SQL commands from a file, statement by statement.
func runSQLFile(db *sql.DB, filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("Failed to read SQL file: %v", err)
	}

	// Split the SQL file content by semicolons
	statements := strings.Split(string(content), ";")

	for _, statement := range statements {
		trimmedStmt := strings.TrimSpace(statement)
		if trimmedStmt == "" {
			continue // Skip empty statements
		}

		if _, err = db.Exec(trimmedStmt); err != nil {
			return fmt.Errorf("Failed to execute SQL statement: %v\nStatement:\n%s", err, trimmedStmt)
		}
	}
	return nil
}

// checkMySQLVersion logs the MySQL version.
func checkMySQLVersion(db *sql.DB) {
	var version string
	err := db.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil {
		fmt.Printf("Failed to retrieve MySQL version: %v", err)
		return
	}
}

// getAbsolutePath constructs an absolute path for a given relative file path.
func getAbsolutePath(relativePath string) string {
	absPath, err := filepath.Abs("../../" + relativePath)
	if err != nil {
		panic(fmt.Sprintf("Failed to get absolute path for %s: %v", relativePath, err))
	}
	return absPath
}
