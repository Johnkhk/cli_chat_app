package pages

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/johnkhk/cli_chat_app/client/app"
	"github.com/johnkhk/cli_chat_app/client/ui/ascii"
	// "github.com/johnkhk/cli_chat_app/client/ui/components"
)

type Choice int

const (
	ChoiceRegister Choice = iota
	ChoiceLogin
)

type model struct {
	cursor    Choice
	choices   []string
	selected  Choice
	rpcClient *app.RpcClient
}

// Initialize the model
func NewLandingModel(rpcClient *app.RpcClient) model {
	return model{
		choices:   []string{"Register", "Login"},
		cursor:    ChoiceRegister,
		selected:  -1,
		rpcClient: rpcClient,
	}
}

// Update function handles the state transitions
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		// Move the cursor up
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}

		// Move the cursor down
		case "down":
			if m.cursor < 1 {
				m.cursor++
			}

		// Select the current choice
		case "enter":
			m.selected = m.cursor
			return handleSelection(m) // Call a function to handle the selection

		case "ctrl+c", "q":
			m.rpcClient.Logger.Info("Exiting the application from landing page")
			return m, tea.Quit
		}
	}

	return m, nil
}

// Handle user selection and navigate to the appropriate page
func handleSelection(m model) (tea.Model, tea.Cmd) {
	switch m.selected {
	case ChoiceRegister:
		// Here, integrate your logic for handling registration
		return NewRegisterModel(m.rpcClient), nil // Switch to registration page
	case ChoiceLogin:
		// Here, integrate your logic for handling login
		return NewLoginModel(m.rpcClient), nil // Switch to login page
		// return m, tea.Quit
	}
	return m, nil
}

// View function renders the UI
func (m model) View() string {
	s := logoStyle.Render(ascii.Logo)
	s += "\n\n" // Add some space after the title
	s += titleStyle.Render("Choose an option:")
	s += "\n\n" // Add some space after the title

	for i, choice := range m.choices {
		cursor := " " // No cursor
		if Choice(i) == m.cursor {
			cursor = ">" // Selected cursor
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	// s += "\nPress Enter to select, Up/Down arrows to navigate."
	s += helpStyle.Render("\nenter: select • up/down: navigate • esc/ctrl+c: quit")

	return s
}

func (m model) Init() tea.Cmd {
	return nil
}
