package app

import (
	"fmt"

	"github.com/muesli/termenv" // For enhanced terminal output (colors, styles)
)

const asciiArt = `
 _______           __    _______ __               __
|   |   |.-----.--|  |  |    ___|  |--.-----.----|  |
|       ||  _  |  _  |  |    ___|     |  -__|  __|  |
|___|___||_____|_____|  |_______|__|__|_____|____|__|
`

// DisplayBanner prints the ASCII art banner to the terminal.
func DisplayBanner() {
	// Set up the termenv output for styling.
	p := termenv.ColorProfile()
	banner := termenv.String(asciiArt).Foreground(p.Color("5")).Bold()

	fmt.Println(banner)
	displayInfo()

}

// displayInfo prints additional info to help the user get started.
func displayInfo() {
	fmt.Println("\nWelcome to the Terminal Chat App!")
	fmt.Println("===================================")
	fmt.Println("This app allows you to chat with friends securely in a terminal environment.\n")

	fmt.Println("To get started:")
	fmt.Println("- Type '/register' to create a new account.")
	fmt.Println("- Type '/login' to log in with your existing account.")
	fmt.Println()

	fmt.Println("Available Commands:")
	fmt.Println("- /register    Register a new user account.")
	fmt.Println("- /login       Log in with your username and password.")
	fmt.Println("- /addfriend   Add a friend by their username.")
	fmt.Println("- /senddm      Send a direct message to a friend.")
	fmt.Println("- /help        Display this help information again.")
	fmt.Println("- /quit        Exit the application.\n")

	fmt.Println("Tip: Use '/help' at any time to see this list of commands again.\n")
}
