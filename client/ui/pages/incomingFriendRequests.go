package pages

import (
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/johnkhk/cli_chat_app/client/app"
	"github.com/johnkhk/cli_chat_app/genproto/friends"
)

type requestActionMsg struct {
	Success bool
	Message string
	Action  string // "accepted" or "declined"
}

type incomingRequestsModel struct {
	incomingRequests      []*friends.FriendRequest
	incFriendRequestTable table.Model
	rpcClient             *app.RpcClient
}

func NewIncomingRequestsModel(rpcClient *app.RpcClient) incomingRequestsModel {
	columns := []table.Column{
		{Title: "Name", Width: 20},
		{Title: "Action", Width: 30},
	}

	rows := []table.Row{
		{"David", "Awaiting Response"},
		{"Eve", "Awaiting Response"},
		{"Frank", "Awaiting Response"},
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

	return incomingRequestsModel{
		incFriendRequestTable: t,
		rpcClient:             rpcClient,
	}
}

func (m incomingRequestsModel) Init() tea.Cmd {
	return nil
}

func (m incomingRequestsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case IncomingFriendRequestsMsg:
		if msg.Err != nil {
			// Handle error (e.g., log it or display a message)
		} else {
			Rows := make([]table.Row, len(msg.Requests))
			for i, request := range msg.Requests {
				Rows[i] = table.Row{request.SenderId, "Awaiting Response"}
			}
			m.incFriendRequestTable.SetRows(Rows)
			if msg.Err != nil {
				// Handle error (e.g., log it or display a message)
				m.rpcClient.Logger.Errorf("Error fetching incoming friend requests: %v", msg.Err)
			} else {
				m.rpcClient.Logger.Infof("Received incoming friend requests: %v", msg.Requests)
				m.incomingRequests = msg.Requests
			}
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			// Accept the selected friend request
			// selectedName := m.receivedRequests.SelectedRow()[0]
			// cmd = handleFriendRequestCmd(m.rpcClient, selectedName, "accept")
			return m, cmd
		case "d":
			// Decline the selected friend request
			// selectedName := m.receivedRequests.SelectedRow()[0]
			// cmd = handleFriendRequestCmd(m.rpcClient, selectedName, "decline")
			return m, cmd
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case requestActionMsg:
		if msg.Success {
			// Remove the friend request from the list upon acceptance or decline
			oldRows := m.incFriendRequestTable.Rows()
			var newRows []table.Row
			for _, row := range oldRows {
				if row[0] != msg.Message {
					newRows = append(newRows, row)
				}
			}
			m.incFriendRequestTable.SetRows(newRows)
		}
		return m, nil
	}

	// Update the table state
	m.incFriendRequestTable, cmd = m.incFriendRequestTable.Update(msg)
	return m, cmd
}

func (m incomingRequestsModel) View() string {
	var view strings.Builder

	// Title for the received requests table
	view.WriteString(titleStyle.Render("Incoming Friend Requests:"))
	view.WriteString("\n")
	// Render the received requests table
	view.WriteString(baseStyle.Render(m.incFriendRequestTable.View()) + "\n")
	// Show instructions
	view.WriteString("[ Press 'a' to Accept, 'd' to Decline, 'q' to Quit ]\n")

	return view.String()
}

func handleFriendRequestCmd(rpcClient *app.RpcClient, username, action string) tea.Cmd {
	return func() tea.Msg {
		var err error
		if action == "accept" {
			// err = rpcClient.FriendsClient.AcceptFriendRequest(username)
		} else if action == "decline" {
			// err = rpcClient.FriendsClient.DeclineFriendRequest(username)
		}
		if err != nil {
			return requestActionMsg{Success: false, Message: err.Error(), Action: action}
		}
		return requestActionMsg{Success: true, Message: username, Action: action}
	}
}
