package pages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/johnkhk/cli_chat_app/client/app"
)

type mainMenuModel struct {
	rpcClient      *app.RpcClient
	terminalWidth  int
	terminalHeight int
	focusState     focusState
	friendsModel   DummyModel // Replace with actual friends list model.
	chatModel      ChatModel  // The chat model is now part of the main menu model.
}

type focusState uint

const (
	mainPanel focusState = iota
	leftPanel
	rightPanel
)

// Initialize the main menu model
func NewMainMenuModel(rpcClient *app.RpcClient) mainMenuModel {
	return mainMenuModel{
		rpcClient:    rpcClient,
		friendsModel: NewDummyModel(),         // Replace with actual friends list model.
		chatModel:    NewChatModel(rpcClient), // Initialize the chat model.
	}
}

// Update function for main menu
func (m mainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

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
		// Update the main panel (no additional logic here)
	case leftPanel:
		// Update the left panel (friends list, etc.)
		// m.friendsModel, _ = m.friendsModel.Update(msg)
	case rightPanel:
		// Update the chat model when the right panel is focused
		// var chatCmd tea.Cmd
		// m.chatModel, chatCmd = m.chatModel.Update(msg)
		// cmd = tea.Batch(cmd, chatCmd)
		updatedChatModel, subCmd := m.chatModel.Update(msg)
		cmd = tea.Batch(cmd, subCmd)
		m.chatModel = updatedChatModel.(ChatModel)
	}

	return m, cmd
}

// View function renders the Main Menu UI
func (m mainMenuModel) View() string {
	leftPanelContent := "Friends List\n1. Alice\n2. Bob\n3. Charlie"
	// Right panel will now display the chat content.
	rightPanelContent := m.chatModel.View()

	// Define the margin from all edges
	margin := 2

	// Calculate panel widths based on percentage
	leftPanelWidth := int(float64(m.terminalWidth) * 0.20)             // 20% for the left panel
	rightPanelWidth := m.terminalWidth - leftPanelWidth - (margin * 2) // Remaining 80% for the right panel

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

	// Render the left panel
	leftPanel := leftPanelStyle.
		Width(leftPanelWidth).
		Height(m.terminalHeight - (margin * 2)). // Adjust height based on terminal size
		Render(leftPanelContent)

	// Render the right panel with the chat model content
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
