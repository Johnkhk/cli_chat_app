package pages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/johnkhk/cli_chat_app/client/app"
)

type mainMenuModel struct {
	rpcClient      *app.AuthClient
	terminalWidth  int
	terminalHeight int
	focusState     focusState
	friendsModel   friendsModel
	// chatModel	  chatModel
}

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
		friendsModel: NewFriendsModel(),
		// chatModel: NewChatModel(),
	}
}

// Update function for main menu (currently no state updates needed)
func (m mainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.terminalWidth = msg.Width
		m.terminalHeight = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			// Cycle through the focus states
			switch m.focusState {
			case mainPanel:
				m.focusState = leftPanel
			case leftPanel:
				m.focusState = rightPanel
			case rightPanel:
				m.focusState = mainPanel
			}
		// Quit the program on "ctrl+c" or "q"
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	switch m.focusState {
	case mainPanel:
		// m, _ = m.Update(msg)
		m.rpcClient.Logger.Println("Main panel focused")
	case leftPanel:
		m.friendsModel, _ = m.friendsModel.Update(msg)
	case rightPanel:
		// m.chatModel, cmd = m.chatModel.Update(msg)
		m.rpcClient.Logger.Println("Right panel focused")

	}

	return m, nil
}

// View function renders the Main Menu UI
func (m mainMenuModel) View() string {
	leftPanelContent := "Friends List\n1. Alice\n2. Bob\n3. Charlie"
	rightPanelContent := "Main Chat\nHello, world! Press any key to quit."

	// Define the margin from all edges
	margin := 2

	// Calculate panel widths based on percentage
	leftPanelWidth := int(float64(m.terminalWidth) * 0.20)             // 20% for the left panel
	rightPanelWidth := m.terminalWidth - leftPanelWidth - (margin * 2) // Remaining 80% for the right panel

	// // Render the left panel
	// leftPanel := leftPanelStyle.
	// 	Width(leftPanelWidth).
	// 	Height(m.terminalHeight - (margin * 2)). // Adjust height based on terminal size
	// 	Render(leftPanelContent)

	// // Render the right panel
	// rightPanel := rightPanelStyle.
	// 	Width(rightPanelWidth).
	// 	Height(m.terminalHeight - (margin * 2)). // Adjust height based on terminal size
	// 	Render(rightPanelContent)

	// // Combine both panels side by side
	// finalView := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)

	// return finalView + helpStyle.Render("\nesc/ctrl+c: quit")
	// Determine the border style based on the current focus state
	var leftPanelStyle, rightPanelStyle lipgloss.Style

	switch m.focusState {
	case mainPanel:
		// Both panels have gray borders
		leftPanelStyle = grayBorderStyle
		rightPanelStyle = grayBorderStyle
	case leftPanel:
		// Left panel has a blue border, right panel has a gray border
		leftPanelStyle = blueBorderStyle
		rightPanelStyle = grayBorderStyle
	case rightPanel:
		// Right panel has a blue border, left panel has a gray border
		leftPanelStyle = grayBorderStyle
		rightPanelStyle = blueBorderStyle
	}

	// Render the left panel
	leftPanel := leftPanelStyle.
		Width(leftPanelWidth).
		Height(m.terminalHeight - (margin * 2)). // Adjust height based on terminal size
		Render(leftPanelContent)

	// Render the right panel
	rightPanel := rightPanelStyle.
		Width(rightPanelWidth).
		Height(m.terminalHeight - (margin * 2)). // Adjust height based on terminal size
		Render(rightPanelContent)

	// Combine both panels side by side
	finalView := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)

	// Add a help bar or instructions at the bottom
	helpBar := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("\nPress Tab to switch panels | esc/ctrl+c: quit")

	return finalView + helpBar

}

func (m mainMenuModel) Init() tea.Cmd {
	return nil
}
