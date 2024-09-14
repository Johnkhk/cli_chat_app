package pages

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/johnkhk/cli_chat_app/client/app"
	"github.com/johnkhk/cli_chat_app/genproto/friends"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type friendAddedMsg struct {
	Success bool
	Message string
}

type outgoingRequestsModel struct {
	outgoingRequests  []*friends.FriendRequest
	sentRequestsTable table.Model
	showInput         bool
	textInput         textinput.Model
	rpcClient         *app.RpcClient
}

func NewOutgoingRequestsModel(rpcClient *app.RpcClient) outgoingRequestsModel {
	columns := []table.Column{
		{Title: "Name", Width: 20},
		{Title: "Status", Width: 30},
	}

	rows := []table.Row{
		{"Alice", "Pending"},
		{"Bob", "Pending"},
		{"Charlie", "Pending"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	ti := textinput.New()
	ti.Placeholder = "Username (esc to cancel)"
	ti.CharLimit = 20
	ti.Width = 30

	return outgoingRequestsModel{
		sentRequestsTable: t,
		showInput:         false,
		textInput:         ti,
		rpcClient:         rpcClient,
	}
}

func (m outgoingRequestsModel) Init() tea.Cmd {
	return nil
}

func (m outgoingRequestsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case OutgoingFriendRequestsMsg:
		if msg.Err != nil {
			// Handle error (e.g., log it or display a message)
			m.rpcClient.Logger.Errorf("Error fetching incoming friend requests: %v", msg.Err)
		} else {
			m.rpcClient.Logger.Infof("Received outgoing friend requests: %v", msg.Requests)
			m.outgoingRequests = msg.Requests
		}
	case tea.KeyMsg:
		if m.showInput {
			m.textInput, cmd = m.textInput.Update(msg)
			if msg.String() == "enter" {
				newFriendName := m.textInput.Value()
				if newFriendName != "" {
					m.showInput = false
					cmd = addFriendCmd(m.rpcClient, newFriendName)
					m.textInput.SetValue("")
					return m, cmd
				}
				m.showInput = false
				return m, nil
			} else if msg.String() == "esc" {
				m.showInput = false
				return m, nil
			}
			return m, cmd
		}

		switch msg.String() {
		case "a":
			m.showInput = true
			m.textInput.Focus()
			return m, textinput.Blink
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Selected: %s", m.sentRequestsTable.SelectedRow()[0]),
			)
		}
	case friendAddedMsg:
		if msg.Success {
			newRow := table.Row{msg.Message, fmt.Sprintf("%s is a new friend", msg.Message)}
			oldRows := m.sentRequestsTable.Rows()
			m.sentRequestsTable.SetRows(append(oldRows, newRow))

		} else {
			// fmt.Printf("Failed to add friend: %s\n", msg.Message)
		}
		return m, nil
	}

	if m.showInput {
		m.textInput, cmd = m.textInput.Update(msg)
	} else {
		m.sentRequestsTable, cmd = m.sentRequestsTable.Update(msg)
	}
	return m, cmd
}

func (m outgoingRequestsModel) View() string {
	// if m.showInput {
	// 	return baseStyle.Render(m.sentRequests.View()) + "\n" + m.textInput.View()
	// }
	// return baseStyle.Render(m.sentRequests.View()) + "\n[ Press 'a' to Add Friend ]"
	var view strings.Builder

	// Title for the sent requests table
	view.WriteString(titleStyle.Render("Sent Friend Requests:"))
	view.WriteString("\n")
	// Render the sent requests table
	view.WriteString(baseStyle.Render(m.sentRequestsTable.View()) + "\n")

	// Show the input if it's active
	if m.showInput {
		view.WriteString(m.textInput.View() + "\n")
	} else {
		view.WriteString("[ Press 'a' to Add Friend ]\n")
	}

	return view.String()
}

func addFriendCmd(rpcClient *app.RpcClient, username string) tea.Cmd {
	return func() tea.Msg {
		err := rpcClient.FriendsClient.SendFriendRequest(username)
		if err != nil {
			return friendAddedMsg{Success: false, Message: err.Error()}
		}
		return friendAddedMsg{Success: true, Message: username}
	}
}
