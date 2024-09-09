package pages

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Style for the friend list panel without border
var friendListStyle = lipgloss.NewStyle().
	// Margin(1, 2). // Add margins to the style without a border
	Align(lipgloss.Left)

// Define a struct to represent a friend as an item in the list
type friendItem struct {
	title, desc string
}

func (i friendItem) Title() string       { return i.title }
func (i friendItem) Description() string { return i.desc }
func (i friendItem) FilterValue() string { return i.title }

// friendsModel structure with state for the friend list
type friendsModel struct {
	list list.Model
}

func (m friendsModel) Init() tea.Cmd {
	return nil
}

// Initialize the friends model with a list of friends
func NewFriendsModel() friendsModel {
	items := []list.Item{
		friendItem{title: "Alice", desc: "Alice is a good friend"},
		friendItem{title: "Bob", desc: "Bob is a great friend"},
		friendItem{title: "Charlie", desc: "Charlie is a close friend"},
		friendItem{title: "Alice", desc: "Alice is a good friend"},
		friendItem{title: "Bob", desc: "Bob is a great friend"},
		friendItem{title: "Charlie", desc: "Charlie is a close friend"},
		friendItem{title: "Alice", desc: "Alice is a good friend"},
		friendItem{title: "Alice", desc: "Alice is a good friend"},
		friendItem{title: "Alice", desc: "Alice is a good friend"},
		friendItem{title: "Alice", desc: "Alice is a good friend"},
		friendItem{title: "Bob", desc: "Bob is a great friend"},
		friendItem{title: "Charlie", desc: "Charlie is a close friend"},
		friendItem{title: "Alice", desc: "Alice is a good friend"},
		friendItem{title: "Bob", desc: "Bob is a great friend"},
		friendItem{title: "Charlie", desc: "Charlie is a close friend"},
		friendItem{title: "Bob", desc: "Bob is a great friend"},
		friendItem{title: "Charlie", desc: "Charlie is a close friend"},
		friendItem{title: "Alice", desc: "Alice is a good friend"},
		friendItem{title: "Bob", desc: "Bob is a great friend"},
		friendItem{title: "Charlie", desc: "Charlie is a close friend"},
		friendItem{title: "Bob", desc: "Bob is a great friend"},
		friendItem{title: "Charlie", desc: "Charlie is a close friend"},
		friendItem{title: "Alice", desc: "Alice is a good friend"},
		friendItem{title: "Bob", desc: "Bob is a great friend"},
		friendItem{title: "Charlie", desc: "Charlie is a close friend"},
		friendItem{title: "Bob", desc: "Bob is a great friend"},
		friendItem{title: "Charlie", desc: "Charlie is a close friend"},
		friendItem{title: "Alice", desc: "Alice is a good friend"},
		friendItem{title: "Bob", desc: "Bob is a great friend"},
		friendItem{title: "Charlie", desc: "Charlie is a close friend"},
	}

	// Create a new list model with the friend items
	l := list.New(items, list.NewDefaultDelegate(), 0, 0) // Initial size set to zero; will adjust dynamically
	l.Title = "Friends List"
	l.SetStatusBarItemName("Friend", "Friends")
	l.SetShowHelp(false)
	// l.SetItems(items)
	l.SetShowFilter(true)

	return friendsModel{list: l}
}

// Update function to handle key inputs for the friend list
func (m friendsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Select friend for DM (trigger some action or callback)
			selectedItem := m.list.SelectedItem()
			if selectedItem != nil {
				friend := selectedItem.(friendItem)
				// Perform action with the selected friend, e.g., open chat
				m.list.NewStatusMessage("Selected friend: " + friend.title)
			}
		}
	}

	// Let the list model handle its own updates
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View function renders the friend list UI
func (m friendsModel) View() string {
	// Apply the style and render the list without an additional border
	return friendListStyle.Render(m.list.View())
	// return m.list.View()
}

// Add a new friend to the friend list
func (m *friendsModel) AddFriend(name string) {
	m.list.InsertItem(len(m.list.Items()), friendItem{title: name})
}

// Remove a friend from the friend list
func (m *friendsModel) RemoveFriend(index int) {
	if index >= 0 && index < len(m.list.Items()) {
		m.list.RemoveItem(index)
	}
}
