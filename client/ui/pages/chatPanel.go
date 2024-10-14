package pages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/johnkhk/cli_chat_app/client/app"
)

type focusState uint

const (
	leftPanel focusState = iota
	rightPanel
)

// Main menu model now has ChatModel as a pointer for in-place updates
type ChatPanelModel struct {
	rpcClient      *app.RpcClient
	terminalWidth  int
	terminalHeight int
	focusState     focusState
	// friendsModel   DummyModel // Replace with actual friends list model.
	friendsModel *ChatFriendListModel // Replace with actual friends list model.
	chatModel    *ChatModel           // Use a pointer to the ChatModel.
}

// Initialize the main menu model
func NewChatPanelModel(rpcClient *app.RpcClient) ChatPanelModel {
	chat := NewChatModel(rpcClient)
	// friend := NewFriendListModel(rpcClient)
	friend := NewChatFriendListModel(rpcClient)
	mm := ChatPanelModel{
		rpcClient: rpcClient,
		// friendsModel: NewDummyModel(), // Replace with actual friends list model.
		friendsModel: &friend, // Replace with actual friends list model.

		chatModel: &chat, // Pass chat as a pointer.
	}
	return mm
}

func (m ChatPanelModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Handle global messages first, like window resizing and quitting
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.terminalWidth = msg.Width
		m.terminalHeight = msg.Height
		// Pass the window size message to both panels
		// m.friendsModel, _ = m.friendsModel.Update(msg)
		// m.chatModel, _ = m.chatModel.Update(msg)

	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			// Switch focus between panels
			if m.focusState == leftPanel {
				m.focusState = rightPanel
			} else {
				m.focusState = leftPanel
			}
		case "ctrl+c", "q":
			return m, tea.Quit

		default:

			// Pass key messages to the active panel only
			switch m.focusState {
			case leftPanel:
				if msg.String() == "f" {
					friendManagementModel := NewFriendManagementModel(m.rpcClient)

					// Initialize the chat panel (e.g., fetch friends) and get the command for it
					chatCmd := friendManagementModel.Init() // Init returns the command for loading friends

					// Update the chat panel to handle the incoming message
					// updatedChatPanelModel, updateCmd := chatPanelModel.Update(msg)
					updatedChatPanelModel, updateCmd := friendManagementModel.Update(tea.WindowSizeMsg{Width: m.terminalWidth, Height: m.terminalHeight})

					// Return only the updated chat model and its commands, ignoring the previous model's cmds
					return updatedChatPanelModel, tea.Batch(chatCmd, updateCmd)
				}
				// Switch to the friend management page

				// var friendCmd tea.Cmd
				// m.friendsModel, friendCmd = m.friendsModel.Update(msg)
				// cmd = tea.Batch(cmd, friendCmd)
				friendModel, friendCmd := m.friendsModel.Update(msg)
				castedFriendModel, ok := friendModel.(ChatFriendListModel)
				if !ok {
					m.rpcClient.Logger.Error("Failed to assert tea.Model to ChatFriendListModel")
					return m, nil
				}
				m.friendsModel = &castedFriendModel
				cmd = tea.Batch(cmd, friendCmd)
			case rightPanel:
				// var chatCmd tea.Cmd
				// m.chatModel, chatCmd = m.chatModel.Update(msg)
				// cmd = tea.Batch(cmd, chatCmd)
				chatModel, chatCmd := m.chatModel.Update(msg)
				castedChatModel, ok := chatModel.(ChatModel) // Assert the tea.Model to ChatModel type
				if !ok {
					// If type assertion fails, return the current state and log an error
					m.rpcClient.Logger.Error("Failed to assert tea.Model to ChatModel")
					return m, nil
				}
				m.chatModel = &castedChatModel // Use the address-of operator to get the pointer
				cmd = tea.Batch(cmd, chatCmd)
			}
		}

	case FriendSelectedMsg:
		// When a friend is selected, set the active user ID in the chat model.
		m.chatModel.SetActiveUser(msg.UserID, msg.Username)
		m.rpcClient.Logger.Infof("Switched to chat with user ID: %d", msg.UserID)
		m.focusState = rightPanel

	default:
		// For other messages, update both models as necessary
		// m.friendsModel, _ = m.friendsModel.Update(msg)
		friendModel, friendCmd := m.friendsModel.Update(msg)
		castedFriendModel, ok := friendModel.(ChatFriendListModel)
		if !ok {
			m.rpcClient.Logger.Error("Failed to assert tea.Model to ChatFriendListModel")
			return m, nil
		}
		m.friendsModel = &castedFriendModel
		cmd = tea.Batch(cmd, friendCmd)
		// m.chatModel, _ = m.chatModel.Update(msg)
		chatModel, chatCmd := m.chatModel.Update(msg)
		castedChatModel, ok := chatModel.(ChatModel) // Assert the tea.Model to ChatModel type
		if !ok {
			// If type assertion fails, return the current state and log an error
			m.rpcClient.Logger.Error("Failed to assert tea.Model to ChatModel")
			return m, nil
		}
		m.chatModel = &castedChatModel // Use the address-of operator to get the pointer
		cmd = tea.Batch(cmd, chatCmd)
	}

	return m, cmd
}

// View function renders the Main Menu UI
func (m ChatPanelModel) View() string {
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

	return finalView + m.renderHelpBar()
}

func (m ChatPanelModel) renderHelpBar() string {
	helpBarStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	// Determine the content based on the focused state
	var helpBarContent string
	if m.focusState == leftPanel {
		helpBarContent = "\nPress Tab to switch panels | esc/ctrl+c: quit | f: friends management"
	} else {
		helpBarContent = "\nPress Tab to switch panels | esc/ctrl+c: quit"
	}

	// Render and return the styled help bar
	return helpBarStyle.Render(helpBarContent)
}

func (m ChatPanelModel) Init() tea.Cmd {
	// Initialize both chatModel and friendsModel and batch their commands together
	return tea.Batch(
		m.chatModel.Init(),
		m.friendsModel.Init(),
	)
}
