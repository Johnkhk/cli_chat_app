// incoming_requests_model.go

package pages

import (
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/johnkhk/cli_chat_app/client/app"
	"github.com/johnkhk/cli_chat_app/genproto/friends"
)

type incomingRequestsModel struct {
	incomingRequests      []*friends.FriendRequest
	incFriendRequestTable table.Model
	rpcClient             *app.RpcClient
}

func NewIncomingRequestsModel(rpcClient *app.RpcClient) incomingRequestsModel {
	columns := []table.Column{
		{Title: "Sender", Width: 20},
		{Title: "Status", Width: 20},
	}

	// Initially empty rows
	rows := []table.Row{}

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
			m.rpcClient.Logger.Errorf("Error fetching incoming friend requests: %v", msg.Err)
		} else {
			m.rpcClient.Logger.Infof("Received incoming friend requests: %v", msg.Requests)
			m.incomingRequests = msg.Requests

			// Update the table rows
			rows := make([]table.Row, len(msg.Requests))
			for i, request := range msg.Requests {
				rows[i] = table.Row{request.SenderUsername, request.Status.String()}
			}
			m.incFriendRequestTable.SetRows(rows)
		}

	case AcceptFriendRequestResultMsg:
		if msg.Err != nil {
			m.rpcClient.Logger.Errorf("Failed to accept friend request %d: %v", msg.RequestID, msg.Err)
		} else {
			m.rpcClient.Logger.Infof("Friend request accepted!")
			// Remove the accepted request from the list
			for i, req := range m.incomingRequests {
				if req.RequestId == msg.RequestID {
					m.incomingRequests = append(m.incomingRequests[:i], m.incomingRequests[i+1:]...)
					break
				}
			}
			// Update the table rows
			rows := make([]table.Row, len(m.incomingRequests))
			for i, request := range m.incomingRequests {
				rows[i] = table.Row{request.SenderUsername, request.Status.String()}
			}
			m.incFriendRequestTable.SetRows(rows)
		}

	case DeclineFriendRequestResultMsg:
		if msg.Err != nil {
			m.rpcClient.Logger.Errorf("Failed to decline friend request %d: %v", msg.RequestID, msg.Err)
		} else {
			m.rpcClient.Logger.Infof("Friend request declined.")
			// Remove the declined request from the list
			for i, req := range m.incomingRequests {
				if req.RequestId == msg.RequestID {
					m.incomingRequests = append(m.incomingRequests[:i], m.incomingRequests[i+1:]...)
					break
				}
			}
			// Update the table rows
			rows := make([]table.Row, len(m.incomingRequests))
			for i, request := range m.incomingRequests {
				rows[i] = table.Row{request.SenderUsername, request.Status.String()}
			}
			m.incFriendRequestTable.SetRows(rows)
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			// Accept the selected friend request
			selectedRow := m.incFriendRequestTable.Cursor()
			if selectedRow >= 0 && selectedRow < len(m.incomingRequests) {
				request := m.incomingRequests[selectedRow]
				cmd = func() tea.Msg {
					err := m.rpcClient.FriendsClient.AcceptFriendRequest(request.RequestId)
					return AcceptFriendRequestResultMsg{RequestID: request.RequestId, Err: err}
				}
				return m, cmd
			}
		case "d":
			// Decline the selected friend request
			selectedRow := m.incFriendRequestTable.Cursor()
			if selectedRow >= 0 && selectedRow < len(m.incomingRequests) {
				request := m.incomingRequests[selectedRow]
				cmd = func() tea.Msg {
					err := m.rpcClient.FriendsClient.DeclineFriendRequest(request.RequestId)
					return DeclineFriendRequestResultMsg{RequestID: request.RequestId, Err: err}
				}
				return m, cmd
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		default:
			var tableCmd tea.Cmd
			m.incFriendRequestTable, tableCmd = m.incFriendRequestTable.Update(msg)
			cmd = tea.Batch(cmd, tableCmd)
			return m, cmd
		}

	default:
		// Update the table state with other messages
		var tableCmd tea.Cmd
		m.incFriendRequestTable, tableCmd = m.incFriendRequestTable.Update(msg)
		cmd = tea.Batch(cmd, tableCmd)
	}

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
	view.WriteString("[ ↑/↓: navigate | 'a': Accept | 'd': Decline ]\n")

	return view.String()
}
