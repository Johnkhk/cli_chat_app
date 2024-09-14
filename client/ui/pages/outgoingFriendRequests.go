// outgoing_requests_model.go

package pages

import (
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

type outgoingRequestsModel struct {
	outgoingRequests  []*friends.FriendRequest
	sentRequestsTable table.Model
	showInput         bool
	textInput         textinput.Model
	rpcClient         *app.RpcClient
}

func NewOutgoingRequestsModel(rpcClient *app.RpcClient) outgoingRequestsModel {
	columns := []table.Column{
		{Title: "Recipient", Width: 20},
		{Title: "Status", Width: 30},
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
			m.rpcClient.Logger.Errorf("Error fetching outgoing friend requests: %v", msg.Err)
		} else {
			m.rpcClient.Logger.Infof("Received outgoing friend requests: %v", msg.Requests)
			m.outgoingRequests = msg.Requests

			// Update the table rows
			rows := make([]table.Row, len(msg.Requests))
			for i, request := range msg.Requests {
				rows[i] = table.Row{request.RecipientId, "Pending"}
			}
			m.sentRequestsTable.SetRows(rows)
		}

	case SendFriendRequestResultMsg:
		if msg.Err != nil {
			m.rpcClient.Logger.Errorf("Failed to send friend request to %s: %v", msg.RecipientUsername, msg.Err)
			// Optionally, display an error message in the UI
		} else {
			// Optionally, handle success (the parent model will refresh the data)
		}

	case tea.KeyMsg:
		if m.showInput {
			var tiCmd tea.Cmd
			m.textInput, tiCmd = m.textInput.Update(msg)
			cmd = tea.Batch(cmd, tiCmd)
			switch msg.String() {
			case "enter":
				newFriendName := m.textInput.Value()
				m.textInput.Blur()
				m.sentRequestsTable.Focus()
				m.showInput = false
				if newFriendName != "" {
					cmd = tea.Batch(cmd, func() tea.Msg {
						return SendFriendRequestMsg{RecipientUsername: newFriendName}
					})
					m.textInput.SetValue("")
					return m, cmd
				}
				return m, cmd
			case "esc":
				m.textInput.Blur()
				m.sentRequestsTable.Focus()
				m.showInput = false
				return m, cmd
			default:
				return m, cmd
			}
		} else {
			switch msg.String() {
			case "a":
				m.showInput = true
				m.textInput.Focus()
				m.sentRequestsTable.Blur()
				return m, textinput.Blink
			case "q", "ctrl+c":
				return m, tea.Quit
			default:
				var tableCmd tea.Cmd
				m.sentRequestsTable, tableCmd = m.sentRequestsTable.Update(msg)
				cmd = tea.Batch(cmd, tableCmd)
				return m, cmd
			}
		}

	default:
		// Update the table or text input with other messages
		if m.showInput {
			var tiCmd tea.Cmd
			m.textInput, tiCmd = m.textInput.Update(msg)
			cmd = tea.Batch(cmd, tiCmd)
		} else {
			var tableCmd tea.Cmd
			m.sentRequestsTable, tableCmd = m.sentRequestsTable.Update(msg)
			cmd = tea.Batch(cmd, tableCmd)
		}
	}

	return m, cmd
}

func (m outgoingRequestsModel) View() string {
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
