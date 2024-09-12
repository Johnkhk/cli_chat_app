package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sirupsen/logrus"

	"github.com/johnkhk/cli_chat_app/client/app"
	"github.com/johnkhk/cli_chat_app/client/ui/pages"
)

// Run the appropriate Bubble Tea program based on the login status
func RunUIBasedOnAuthStatus(isLoggedIn bool, log *logrus.Logger, rpcClient *app.RpcClient) {
	if isLoggedIn {
		// log.Info("User automatically logged in with stored tokens.")
		runTeaProgram(pages.NewFriendManagementModel(rpcClient)) // Start the main menu if auto-login succeeds
	} else {
		// log.Info("Automatic login failed or no valid token found.")
		runTeaProgram(pages.NewLandingModel(rpcClient)) // Start the landing page if auto-login fails
	}
}

// Function to run the Bubble Tea program
func runTeaProgram(m tea.Model) {
	p := tea.NewProgram(m, tea.WithAltScreen())

	if err := p.Start(); err != nil {
		fmt.Printf("Error starting Bubble Tea program: %v\n", err)
		os.Exit(1)
	}
}
