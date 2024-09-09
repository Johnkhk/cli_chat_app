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

type friendsModel struct {
	list      list.Model
	showInput bool
	textInput textinput.Model
	rpcClient *app.AuthClient
}

func (m friendsModel) Init() tea.Cmd {
	return nil
}

func NewFriendsModel(rpcClient *app.AuthClient) friendsModel {
	items := []list.Item{
		friendItem{title: "Alice", desc: "Alice is a good friend"},
		friendItem{title: "Bob", desc: "Bob is a great friend"},
		friendItem{title: "Charlie", desc: "Charlie is a close friend"},
		friendItem{title: "Alice", desc: "Alice is a good friend"},
		friendItem{title: "Alice", desc: "Alice is a good friend"},
		friendItem{title: "Alice", desc: "Alice is a good friend"},
		friendItem{title: "Alice", desc: "Alice is a good friend"},
		friendItem{title: "Alice", desc: "Alice is a good friend"},
		friendItem{title: "Alice", desc: "Alice is a good friend"},
		friendItem{title: "Alice", desc: "Alice is a good friend"},
		friendItem{title: "Alice", desc: "Alice is a good friend"},
		friendItem{title: "Alice", desc: "Alice is a good friend"},
		friendItem{title: "Alice", desc: "Alice is a good friend"},
		friendItem{title: "Alice", desc: "Alice is a good friend"},
		friendItem{title: "Alice", desc: "Alice is a good friend"},
		friendItem{title: "Bob", desc: "Bob is a great friend"},
		friendItem{title: "Charlie", desc: "Charlie is a close friend"},
		friendItem{title: "Bob", desc: "Bob is a great friend"},
		friendItem{title: "Charlie", desc: "Charlie is a close friend"},
		friendItem{title: "Bob", desc: "Bob is a great friend"},
		friendItem{title: "Charlie", desc: "Charlie is a close friend"},
		friendItem{title: "Bob", desc: "Bob is a great friend"},
		friendItem{title: "Charlie", desc: "Charlie is a close friend"},
		friendItem{title: "Bob", desc: "Bob is a great friend"},
		friendItem{title: "Charlie", desc: "Charlie is a close friend"},
		friendItem{title: "Bob", desc: "Bob is a great friend"},
		friendItem{title: "Charlie", desc: "Charlie is a close friend"},
		friendItem{title: "Bob", desc: "Bob is a great friend"},
		friendItem{title: "Charlie", desc: "Charlie is a close friend"},
		friendItem{title: "Bob", desc: "Bob is a great friend"},
		friendItem{title: "Charlie", desc: "Charlie is a close friend"},
		friendItem{title: "Bob", desc: "Bob is a great friend"},
		friendItem{title: "Charlie", desc: "Charlie is a close friend"},
		friendItem{title: "Bob", desc: "Bob is a great friend"},
		friendItem{title: "Charlie", desc: "Charlie is a close friend"},
		friendItem{title: "Bob", desc: "Bob is a great friend"},
		friendItem{title: "Charlie", desc: "Charlie is a close friend"},
		friendItem{title: "Bob", desc: "Bob is a great friend"},
		friendItem{title: "Charlie", desc: "Charlie is a close friend"},
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
				// Add friend logic
				newFriendName := m.textInput.Value()
				if newFriendName != "" {
					m.AddFriend(newFriendName)
					m.textInput.SetValue("")
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

func (m *friendsModel) AddFriend(name string) {
	m.list.InsertItem(len(m.list.Items()), friendItem{title: name, desc: fmt.Sprintf("%s is a new friend", name)})
}

func (m *friendsModel) RemoveFriend(index int) {
	if index >= 0 && index < len(m.list.Items()) {
		m.list.RemoveItem(index)
	}
}
