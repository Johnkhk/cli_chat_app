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
		{Title: "Request ID", Width: 20},
		{Title: "Sender", Width: 30},
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
				rows[i] = table.Row{request.RequestId, request.SenderId}
			}
			m.incFriendRequestTable.SetRows(rows)
		}

	case AcceptFriendRequestResultMsg:
		if msg.Err != nil {
			m.rpcClient.Logger.Errorf("Failed to process friend request %s: %v", msg.RequestID, msg.Err)
			// Optionally, display an error message in the UI
		} else {
			// Optionally, handle success (the parent model will refresh the data)
		}

	case DeclineFriendRequestResultMsg:
		if msg.Err != nil {
			m.rpcClient.Logger.Errorf("Failed to decline friend request %s: %v", msg.RequestID, msg.Err)
			// Optionally, display an error message in the UI
		} else {
			// Optionally, handle success (the parent model will refresh the data)
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			// Accept the selected friend request
			selectedRow := m.incFriendRequestTable.SelectedRow()
			if len(selectedRow) > 0 {
				requestID := selectedRow[0]
				cmd = func() tea.Msg {
					return AcceptFriendRequestMsg{RequestID: requestID}
				}
				return m, cmd
			}
		case "d":
			// Decline the selected friend request
			selectedRow := m.incFriendRequestTable.SelectedRow()
			if len(selectedRow) > 0 {
				requestID := selectedRow[0]
				cmd = func() tea.Msg {
					return DeclineFriendRequestMsg{RequestID: requestID}
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
	view.WriteString("[ Press 'a' to Accept, 'd' to Decline, 'q' to Quit ]\n")

	return view.String()
}
