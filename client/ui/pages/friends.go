package pages

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/johnkhk/cli_chat_app/client/app"
)

var friendListStyle = lipgloss.NewStyle().
	Align(lipgloss.Left)

type friendItem struct {
	title, desc string
}

func (i friendItem) Title() string       { return i.title }
func (i friendItem) Description() string { return i.desc }
func (i friendItem) FilterValue() string { return i.title }

// Define a message type to handle the response from adding a friend
type friendAddedMsg struct {
	Success bool
	Message string
}

type friendsModel struct {
	list      list.Model
	showInput bool
	textInput textinput.Model
	rpcClient *app.RpcClient
}

func (m friendsModel) Init() tea.Cmd {
	return nil
}

func NewFriendsModel(rpcClient *app.RpcClient) friendsModel {
	items := []list.Item{
		friendItem{title: "Alice", desc: "Alice is a good friend"},
		friendItem{title: "Bob", desc: "Bob is a great friend"},
		friendItem{title: "Charlie", desc: "Charlie is a close friend"},
		// Initial friend items...
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Friends List"
	l.SetStatusBarItemName("Friend", "Friends")
	l.SetShowHelp(false)
	l.SetShowFilter(true)

	ti := textinput.New()
	ti.Placeholder = "Username (esc to cancel)"
	ti.CharLimit = 20
	ti.Width = 30

	return friendsModel{
		list:      l,
		showInput: false,
		textInput: ti,
		rpcClient: rpcClient,
	}
}

func (m friendsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			m.showInput = true
			m.textInput.Focus()
			return m, textinput.Blink
		case "enter":
			if m.showInput {
				// Trigger AddFriend request
				newFriendName := m.textInput.Value()
				if newFriendName != "" {
					m.showInput = false
					cmd = addFriendCmd(m.rpcClient, newFriendName) // Create the tea.Cmd for the request
					m.textInput.SetValue("")
					return m, cmd
				}
				m.showInput = false
				return m, nil
			}
		case "esc":
			if m.showInput {
				m.showInput = false
				return m, nil
			}
		}

	case friendAddedMsg: // Handle the response of the AddFriend request
		if msg.Success {
			m.list.InsertItem(len(m.list.Items()), friendItem{title: msg.Message, desc: fmt.Sprintf("%s is a new friend", msg.Message)})
		} else {
			// Handle the error, e.g., show a message to the user
			fmt.Printf("Failed to add friend: %s\n", msg.Message)
		}
		return m, nil
	}

	if m.showInput {
		m.textInput, cmd = m.textInput.Update(msg)
	} else {
		m.list, cmd = m.list.Update(msg)
	}
	return m, cmd
}

func (m friendsModel) View() string {
	if m.showInput {
		return friendListStyle.Render(m.list.View()) + "\n" + m.textInput.View()
	}
	return friendListStyle.Render(m.list.View()) + "\n[ Press 'a' to Add Friend ]"
}

// Command to add a friend by making a gRPC call
func addFriendCmd(rpcClient *app.RpcClient, username string) tea.Cmd {
	return func() tea.Msg {
		err := rpcClient.FriendsClient.AddFriend(username)
		if err != nil {
			return friendAddedMsg{Success: false, Message: err.Error()}
		}
		return friendAddedMsg{Success: true, Message: username}
	}
}
