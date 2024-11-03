package test

import (
	"testing"
	"time"

	"github.com/johnkhk/cli_chat_app/client/app"
)

type MockTimeProvider struct {
	CurrentTime time.Time
}

func (mtp *MockTimeProvider) Now() time.Time {
	return mtp.CurrentTime
}

func (mtp *MockTimeProvider) Advance(d time.Duration) {
	mtp.CurrentTime = mtp.CurrentTime.Add(d)
}

// Helper function for registering and logging in a user
func RegisterAndLoginUser(t *testing.T, client *app.RpcClient, username string) {
	if err := client.AuthClient.RegisterUser(username, "password"); err != nil {
		t.Fatalf("Failed to register %s: %v", username, err)
	}
	if err, _ := client.AuthClient.LoginUser(username, "password"); err != nil {
		t.Fatalf("Failed to login %s: %v", username, err)
	}
}

// Helper function for waiting for the welcome message
func WaitForWelcomeMessage(t *testing.T, client *app.RpcClient, username string) {
	select {
	case msg := <-client.ChatClient.MessageChannel:
		t.Logf("%s received welcome message: %v", username, msg)
	case <-time.After(3 * time.Second):
		t.Fatalf("%s did not receive the welcome message within timeout period", username)
	}
}
