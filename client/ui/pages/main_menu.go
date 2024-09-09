package pages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/johnkhk/cli_chat_app/client/app"
)

// mainMenuModel structure with terminal size fields and state
type mainMenuModel struct {
	rpcClient      *app.AuthClient
	terminalWidth  int
	terminalHeight int
	focusState     focusState
	friendsModel   friendsModel
}

// focusState enum to define focus states
type focusState uint

const (
	mainPanel focusState = iota
	leftPanel
	rightPanel
)

// Initialize the main menu model
func NewMainMenuModel(rpcClient *app.AuthClient) mainMenuModel {
	return mainMenuModel{
		rpcClient:    rpcClient,
		friendsModel: NewFriendsModel(rpcClient),
	}
}

// Update function for main menu to handle key inputs and window resizing
func (m mainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.terminalWidth = msg.Width
		m.terminalHeight = msg.Height
		// m.friendsModel.list.SetSize(msg.Width/5, msg.Height/2)
		leftPanelWidth := int(float64(msg.Width) * 0.30) // 20% of the total width
		// m.friendsModel.list.SetSize(leftPanelWidth-10, msg.Height-4) // Adjust the height as needed
		m.friendsModel.list.SetSize(leftPanelWidth-10, int(float64(msg.Height)*0.8)) // Adjust the height as needed

	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			switch m.focusState {
			case mainPanel:
				m.focusState = leftPanel
			case leftPanel:
				m.focusState = rightPanel
			case rightPanel:
				m.focusState = mainPanel
			}
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	if m.focusState == leftPanel {
		updatedModel, subCmd := m.friendsModel.Update(msg)
		m.friendsModel = updatedModel.(friendsModel) // Type assertion to friendsModel
		cmd = tea.Batch(cmd, subCmd)
	}

	// Return the updated mainMenuModel and any command to be executed
	return m, cmd
}

// View function renders the Main Menu UI with a dynamic split border layout
func (m mainMenuModel) View() string {
	// Define the margin from all edges
	margin := 2

	// Calculate panel widths based on percentage
	leftPanelWidth := int(float64(m.terminalWidth) * 0.30)
	rightPanelWidth := m.terminalWidth - leftPanelWidth - (margin * 2)

	// Determine the border style based on the current focus state
	var leftPanelStyle, rightPanelStyle lipgloss.Style

	switch m.focusState {
	case mainPanel:
		leftPanelStyle = grayBorderStyle
		rightPanelStyle = grayBorderStyle
	case leftPanel:
		leftPanelStyle = blueBorderStyle
		rightPanelStyle = grayBorderStyle
	case rightPanel:
		leftPanelStyle = grayBorderStyle
		rightPanelStyle = blueBorderStyle
	}

	_ = leftPanelStyle
	// Render the left panel (friend list) without nested border
	leftPanel := leftPanelStyle.
		Width(leftPanelWidth).
		Height(m.terminalHeight - (margin * 2)).
		Render(m.friendsModel.View())
	// leftPanel := m.friendsModel.View()

	// Render the right panel (chat view)
	rightPanelContent := "Chat View:\nHello, world! Press any key to quit."
	rightPanel := rightPanelStyle.
		Width(rightPanelWidth).
		Height(m.terminalHeight - (margin * 2)).
		Render(rightPanelContent)

	// Combine both panels side by side
	finalView := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)

	// Add a help bar or instructions at the bottom
	helpBar := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("\nPress Tab to switch panels | esc/ctrl+c: quit")

	return finalView + helpBar

}

// Init function initializes the main menu model
func (m mainMenuModel) Init() tea.Cmd {
	return nil
}
