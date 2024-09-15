// friend_management.go

package pages

import (
	"fmt"
	"strings"
	"time"

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
	statusMessage  string // Message to display
	statusIsError  bool   // True if it's an error message
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

		case "r":
			// Refresh friend list, incoming, and outgoing friend requests
			cmds = append(cmds,
				fetchFriendListCmd(m.rpcClient),
				fetchIncomingFriendRequestsCmd(m.rpcClient),
				fetchOutgoingFriendRequestsCmd(m.rpcClient),
			)
			m.rpcClient.Logger.Info("Refreshing friend list and friend requests")

		default:
			m.rpcClient.Logger.Infoln("Key pressed:", keypress)
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
			m.statusMessage = fmt.Sprintf("Failed to send friend request: %v", msg.Err)
			m.statusIsError = true
		} else {
			m.statusMessage = "Friend request sent successfully."
			m.statusIsError = false
			cmds = append(cmds, fetchOutgoingFriendRequestsCmd(m.rpcClient))
		}
		cmds = append(cmds, clearStatusMessageCmd())

	case AcceptFriendRequestResultMsg:
		if msg.Err != nil {
			m.rpcClient.Logger.Error("Failed to accept friend request:", msg.Err)
			m.statusMessage = fmt.Sprintf("Failed to accept friend request: %v", msg.Err)
			m.statusIsError = true
		} else {
			m.statusMessage = "Friend request accepted."
			m.statusIsError = false
			cmds = append(cmds,
				fetchFriendListCmd(m.rpcClient),
				fetchIncomingFriendRequestsCmd(m.rpcClient),
			)
		}
		cmds = append(cmds, clearStatusMessageCmd())

	case DeclineFriendRequestResultMsg:
		if msg.Err != nil {
			m.rpcClient.Logger.Error("Failed to decline friend request:", msg.Err)
			m.statusMessage = fmt.Sprintf("Failed to decline friend request: %v", msg.Err)
			m.statusIsError = true
		} else {
			m.statusMessage = "Friend request declined."
			m.statusIsError = false
			cmds = append(cmds, fetchIncomingFriendRequestsCmd(m.rpcClient))
		}
		cmds = append(cmds, clearStatusMessageCmd())

	case RemoveFriendResultMsg:
		if msg.Err != nil {
			m.rpcClient.Logger.Error("Failed to remove friend:", msg.Err)
			m.statusMessage = fmt.Sprintf("Failed to remove friend: %v", msg.Err)
			m.statusIsError = true
		} else {
			m.statusMessage = "Friend removed successfully."
			m.statusIsError = false
			cmds = append(cmds, fetchFriendListCmd(m.rpcClient))
		}
		cmds = append(cmds, clearStatusMessageCmd())

	// Clear status message after delay
	case ClearStatusMessageMsg:
		m.statusMessage = ""
		m.statusIsError = false
	}

	// Update the currently active tab's content
	updatedModel, subCmd := m.tabContent[m.activeTab].Update(msg)
	if updatedContent, ok := updatedModel.(tea.Model); ok {
		m.tabContent[m.activeTab] = updatedContent
	}
	cmds = append(cmds, subCmd)

	return m, tea.Batch(cmds...)
}

// Command to clear the status message after a delay
func clearStatusMessageCmd() tea.Cmd {
	return tea.Tick(time.Second*5, func(t time.Time) tea.Msg {
		return ClearStatusMessageMsg{}
	})
}

type ClearStatusMessageMsg struct{}

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
	availableHeight := m.terminalHeight - 10

	// Render the main content
	doc.WriteString(
		windowStyle.
			Width(availableWidth).
			Height(availableHeight).
			Render(m.tabContent[m.activeTab].View()),
	)

	// Render the status message if it exists
	if m.statusMessage != "" {
		var status string
		if m.statusIsError {
			status = errorMsgStyle.Render("\n" + m.statusMessage)
		} else {
			status = successMsgStyle.Render("\n" + m.statusMessage)
		}
		doc.WriteString(status)
	}

	// Render the help message
	help := helpStyle.Render("\nesc/ctrl+c: quit | tab: switch tab | r: refresh")
	doc.WriteString(help)

	return docStyle.Align(lipgloss.Center).
		Width(m.terminalWidth).
		Height(m.terminalHeight).
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
