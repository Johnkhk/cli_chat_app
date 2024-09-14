// friend_management.go

package pages

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/johnkhk/cli_chat_app/client/app"
)

type FriendManagementModel struct {
	rpcClient      *app.RpcClient
	terminalWidth  int
	terminalHeight int
	activeTab      int
	tabs           []string
	tabContent     []tea.Model
}

func NewFriendManagementModel(rpcClient *app.RpcClient) FriendManagementModel {
	friendListModel := NewFriendListModel(rpcClient)
	incomingModel := NewIncomingRequestsModel(rpcClient)
	outgoingModel := NewOutgoingRequestsModel(rpcClient)

	return FriendManagementModel{
		rpcClient:  rpcClient,
		tabs:       []string{"Friends", "Incoming", "Outgoing"},
		activeTab:  0,
		tabContent: []tea.Model{friendListModel, incomingModel, outgoingModel},
	}
}

func (m FriendManagementModel) Init() tea.Cmd {
	return tea.Batch(
		fetchFriendListCmd(m.rpcClient),
		fetchIncomingFriendRequestsCmd(m.rpcClient),
		fetchOutgoingFriendRequestsCmd(m.rpcClient),
	)
}

func (m FriendManagementModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.terminalWidth = msg.Width
		m.terminalHeight = msg.Height

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			m.rpcClient.Logger.Info("Exiting the application from main menu")
			return m, tea.Quit

		case "tab":
			m.activeTab = (m.activeTab + 1) % len(m.tabs)

		case "shift+tab":
			m.activeTab = (m.activeTab - 1 + len(m.tabs)) % len(m.tabs)
		default:
			m.rpcClient.Logger.Infoln("Using mah boi", keypress)
		}

	// Data messages: pass to child models
	case FriendListMsg:
		updatedModel, subCmd := m.tabContent[0].Update(msg)
		m.tabContent[0] = updatedModel
		cmds = append(cmds, subCmd)

	case IncomingFriendRequestsMsg:
		updatedModel, subCmd := m.tabContent[1].Update(msg)
		m.tabContent[1] = updatedModel
		cmds = append(cmds, subCmd)

	case OutgoingFriendRequestsMsg:
		updatedModel, subCmd := m.tabContent[2].Update(msg)
		m.tabContent[2] = updatedModel
		cmds = append(cmds, subCmd)

	// Action messages: execute commands
	case SendFriendRequestMsg:
		cmd := sendFriendRequestCmd(m.rpcClient, msg.RecipientUsername)
		cmds = append(cmds, cmd)

	case AcceptFriendRequestMsg:
		cmd := acceptFriendRequestCmd(m.rpcClient, msg.RequestID)
		cmds = append(cmds, cmd)

	case DeclineFriendRequestMsg:
		cmd := declineFriendRequestCmd(m.rpcClient, msg.RequestID)
		cmds = append(cmds, cmd)

	case RemoveFriendMsg:
		cmd := removeFriendCmd(m.rpcClient, msg.FriendID)
		cmds = append(cmds, cmd)

	// Result messages: handle outcomes
	case SendFriendRequestResultMsg:
		if msg.Err != nil {
			m.rpcClient.Logger.Error("Failed to send friend request:", msg.Err)
			// Optionally, notify the child model or display an error
		} else {
			// Refresh outgoing friend requests
			cmds = append(cmds, fetchOutgoingFriendRequestsCmd(m.rpcClient))
		}

	case AcceptFriendRequestResultMsg:
		if msg.Err != nil {
			m.rpcClient.Logger.Error("Failed to accept friend request:", msg.Err)
			// Optionally, notify the child model or display an error
		} else {
			// Refresh friend list and incoming requests
			cmds = append(cmds, fetchFriendListCmd(m.rpcClient))
			cmds = append(cmds, fetchIncomingFriendRequestsCmd(m.rpcClient))
		}

	case DeclineFriendRequestResultMsg:
		if msg.Err != nil {
			m.rpcClient.Logger.Error("Failed to decline friend request:", msg.Err)
			// Optionally, notify the child model or display an error
		} else {
			// Refresh incoming requests
			cmds = append(cmds, fetchIncomingFriendRequestsCmd(m.rpcClient))
		}

	case RemoveFriendResultMsg:
		if msg.Err != nil {
			m.rpcClient.Logger.Error("Failed to remove friend:", msg.Err)
			// Optionally, notify the child model or display an error
		} else {
			// Refresh friend list
			cmds = append(cmds, fetchFriendListCmd(m.rpcClient))
		}
	}
	// Update the currently active tab's content
	updatedModel, subCmd := m.tabContent[m.activeTab].Update(msg)
	if updatedContent, ok := updatedModel.(tea.Model); ok {
		m.tabContent[m.activeTab] = updatedContent
	}
	cmds = append(cmds, subCmd)

	return m, tea.Batch(cmds...)
}

func (m FriendManagementModel) View() string {
	doc := strings.Builder{}

	// Render the tabs
	var renderedTabs []string
	for i, t := range m.tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.tabs)-1, i == m.activeTab
		if isActive {
			style = activeTabStyle
		} else {
			style = inactiveTabStyle
		}

		// Adjust borders for the tabs
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Width(m.terminalWidth/len(m.tabs)-5).Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")

	availableWidth := (lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())
	availableHeight := m.terminalHeight - 5

	doc.WriteString(
		windowStyle.
			Width(availableWidth).
			Height(availableHeight).
			Render(m.tabContent[m.activeTab].View()),
	)

	return docStyle.Align(lipgloss.Center).Width(m.terminalWidth).Height(m.terminalHeight).
		Render(doc.String())
}

// Helper functions for min and max
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
