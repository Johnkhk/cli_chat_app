// db_setup.go
package setup

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// SetupTestDatabase initializes the test database.
func SetupTestDatabase() {
	var err error
	testDBName := "test_db"

	// Get the Data Source Name (DSN) from environment variables
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		panic("DATABASE_URL is not set")
	}

	// Connect to the MySQL server (without specifying a database)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to MySQL server: %v", err))
	}

	// Check MySQL version (optional but useful for debugging)
	checkMySQLVersion()

	// Create the test database if it does not exist
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", testDBName))
	if err != nil {
		panic(fmt.Sprintf("Failed to create test database: %v", err))
	}

	// Explicitly set the database to be used in this connection
	_, err = db.Exec(fmt.Sprintf("USE %s", testDBName))
	if err != nil {
		panic(fmt.Sprintf("Failed to switch to test database: %v", err))
	}

	// Drop all tables in the test database before running the setup script
	dropAllTables(testDBName)

	// Execute the SQL file to set up tables and other structures
	upSQLPath := getAbsolutePath("db/migrations/up.sql")
	runSQLFile(upSQLPath)
}

// TeardownTestDatabase cleans up the test database.
func TeardownTestDatabase() {
	// Optionally run the teardown SQL script if needed
	downSQLPath := getAbsolutePath("db/migrations/down.sql")
	runSQLFile(downSQLPath)

	// Drop the entire test database to ensure a clean state
	_, err := db.Exec("DROP DATABASE IF EXISTS test_db")
	if err != nil {
		panic(fmt.Sprintf("Failed to drop test database: %v", err))
	}

	// Close the database connection
	if db != nil {
		db.Close()
	}
}

// dropAllTables drops all tables in the current database.
func dropAllTables(dbName string) {
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		panic(fmt.Sprintf("Failed to list tables: %v", err))
	}
	defer rows.Close()

	var tableName string
	for rows.Next() {
		if err := rows.Scan(&tableName); err != nil {
			panic(fmt.Sprintf("Failed to scan table name: %v", err))
		}

		// Explicitly use the database name in the DROP statement
		_, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s.%s", dbName, tableName))
		if err != nil {
			panic(fmt.Sprintf("Failed to drop table %s: %v", tableName, err))
		}
	}
}

// runSQLFile executes SQL commands from a file, statement by statement.
func runSQLFile(filePath string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		panic(fmt.Sprintf("Failed to read SQL file: %v", err))
	}

	// // Log the content of the SQL file for debugging
	// fmt.Printf("Executing SQL file: %s\nContent:\n%s\n", filePath, string(content))

	// Split the SQL file content by semicolons
	statements := strings.Split(string(content), ";")

	for _, statement := range statements {
		trimmedStmt := strings.TrimSpace(statement)
		if trimmedStmt == "" {
			continue // Skip empty statements
		}

		_, err = db.Exec(trimmedStmt)
		if err != nil {
			panic(fmt.Sprintf("Failed to execute SQL statement: %v\nStatement:\n%s", err, trimmedStmt))
		}
	}
}

// checkMySQLVersion logs the MySQL version.
func checkMySQLVersion() {
	var version string
	err := db.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil {
		panic(fmt.Sprintf("Failed to retrieve MySQL version: %v", err))
	}
	fmt.Printf("Connected to MySQL version: %s\n", version)
}

// getAbsolutePath constructs an absolute path for a given relative file path.
func getAbsolutePath(relativePath string) string {
	absPath, err := filepath.Abs("../../" + relativePath)
	if err != nil {
		panic(fmt.Sprintf("Failed to get absolute path for %s: %v", relativePath, err))
	}
	return absPath
}
