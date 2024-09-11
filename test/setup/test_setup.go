// test/setup/test_setup.go
package setup

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

// TestMainSetup sets up the environment for all tests.
func TestMainSetup(m *testing.M) {
	// Load environment variables from .env.test file
	if err := godotenv.Load("../../.env.test"); err != nil {
		log.Fatalf("Error loading .env.test file: %v", err)
	}

	// Set the ENV_PATH variable to point to the .env.test file
	absPath, err := filepath.Abs("../../.env.test")
	if err != nil {
		log.Fatalf("Error getting absolute path to .env.test: %v", err)
	}
	os.Setenv("ENV_PATH", absPath)

	// Run database setup code
	SetupTestDatabase()

	// Ensure that the database teardown runs even if something goes wrong
	defer TeardownTestDatabase()

	// Start the gRPC server
	serverCmd := exec.Command("go", "run", "../../cmd/server/main.go")
	serverCmd.Stdout = os.Stdout
	serverCmd.Stderr = os.Stderr

	// Start the server process
	if err := serverCmd.Start(); err != nil {
		panic("Failed to start server: " + err.Error())
	}

	// Ensure the server is stopped after tests
	defer func() {
		// Attempt a graceful shutdown by sending SIGINT
		if err := serverCmd.Process.Signal(syscall.SIGINT); err != nil {
			log.Fatalf("Failed to send SIGINT to server process: %v", err)
		}

		// Wait for the server process to exit or timeout
		done := make(chan error, 1)
		go func() {
			done <- serverCmd.Wait()
		}()

		select {
		case err := <-done:
			if err != nil {
				log.Fatalf("Failed to wait for server process to exit: %v", err)
			}
			log.Println("Server process stopped.")
		case <-time.After(5 * time.Second):
			log.Println("Timeout waiting for server to stop. Killing the process.")
			if err := serverCmd.Process.Kill(); err != nil {
				log.Fatalf("Failed to kill server process: %v", err)
			}
		}
	}()

	// Wait for the server to start
	time.Sleep(2 * time.Second)

	// Run all tests
	code := m.Run()

	// Exit with the test status code
	os.Exit(code)
}
