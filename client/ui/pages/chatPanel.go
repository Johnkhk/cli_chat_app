package pages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/johnkhk/cli_chat_app/client/app"
)

type focusState uint

const (
	mainPanel focusState = iota
	leftPanel
	rightPanel
)

// Main menu model now has ChatModel as a pointer for in-place updates
type mainMenuModel struct {
	rpcClient      *app.RpcClient
	terminalWidth  int
	terminalHeight int
	focusState     focusState
	// friendsModel   DummyModel // Replace with actual friends list model.
	friendsModel *ChatFriendListModel // Replace with actual friends list model.
	chatModel    *ChatModel           // Use a pointer to the ChatModel.
}

// Initialize the main menu model
func NewMainMenuModel(rpcClient *app.RpcClient) mainMenuModel {
	chat := NewChatModel(rpcClient)
	// friend := NewFriendListModel(rpcClient)
	friend := NewChatFriendListModel(rpcClient)
	mm := mainMenuModel{
		rpcClient: rpcClient,
		// friendsModel: NewDummyModel(), // Replace with actual friends list model.
		friendsModel: &friend, // Replace with actual friends list model.

		chatModel: &chat, // Pass chat as a pointer.
	}
	return mm
}

// Update function for main menu
func (m mainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Handle global messages first, like window resizing and quitting
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
	case FriendSelectedMsg:
		// When a friend is selected, set the active user ID in the chat model.
		m.chatModel.SetActiveUser(msg.UserID, msg.Username)
		m.rpcClient.Logger.Infof("Switched to chat with user ID: %d", msg.UserID)
		return m, nil
	case FriendListMsg:
		// When friends list is updated, update the friend list model.
		friendModel, friendCmd := m.friendsModel.Update(msg)
		castedFriendModel, ok := friendModel.(ChatFriendListModel)
		if !ok {
			m.rpcClient.Logger.Error("Failed to assert tea.Model to ChatFriendListModel")
			return m, nil
		}
		m.friendsModel = &castedFriendModel
		cmd = tea.Batch(cmd, friendCmd)
	}

	// Always update the chat model, regardless of the focus state
	// First, assert it to the concrete type `ChatModel`, then convert it to a pointer
	chatModel, chatCmd := m.chatModel.Update(msg)
	castedChatModel, ok := chatModel.(ChatModel) // Assert the tea.Model to ChatModel type
	if !ok {
		// If type assertion fails, return the current state and log an error
		m.rpcClient.Logger.Error("Failed to assert tea.Model to ChatModel")
		return m, nil
	}
	m.chatModel = &castedChatModel // Use the address-of operator to get the pointer
	cmd = tea.Batch(cmd, chatCmd)

	// Always update the friends model
	// m.friendsModel, friendCmd = m.friendsModel.Update(msg)
	// castedFriendModel, ok := m.friendsModel.(FriendListModel)
	// if !ok {
	// 	m.rpcClient.Logger.Error("Failed to assert tea.Model to friendListModel")
	// 	return m, nil
	// }
	// m.friendsModel = &castedFriendModel
	// cmd = tea.Batch(cmd, friendCmd)
	// Update friends model and handle pointer reference
	friendModel, friendCmd := m.friendsModel.Update(msg)
	castedFriendModel, ok := friendModel.(ChatFriendListModel)
	if !ok {
		m.rpcClient.Logger.Error("Failed to assert tea.Model to FriendListModel")
		return m, nil
	}
	m.friendsModel = &castedFriendModel
	cmd = tea.Batch(cmd, friendCmd)

	// Update other models based on focus state if necessary
	switch m.focusState {
	case leftPanel:
		// Update the left panel (e.g., friends list)
		// m.friendsModel, _ = m.friendsModel.Update(msg)
	case rightPanel:
		// Any specific logic for right panel can be handled here.
	}

	return m, cmd
}

// View function renders the Main Menu UI
func (m mainMenuModel) View() string {
	// leftPanelContent := "Friends List\n1. Alice\n2. Bob\n3. Charlie"
	leftPanelContent := m.friendsModel.View()
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

// // Init function for main menu model
// func (m mainMenuModel) Init() tea.Cmd {
// 	return m.chatModel.Init() // Initialize the chat model with a command.

// }

func (m mainMenuModel) Init() tea.Cmd {
	// Initialize both chatModel and friendsModel and batch their commands together
	return tea.Batch(
		m.chatModel.Init(),
		m.friendsModel.Init(),
	)
}
