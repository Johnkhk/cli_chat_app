package pages

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Style for the friend list panel
var friendListStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("63")).
	Align(lipgloss.Left)

// friendsModel structure with state for the friend list
type friendsModel struct {
	friends        []string
	selectedFriend int
}

func (m friendsModel) Init() tea.Cmd {
	// Get the friends list
	// m.friends = []string{"Alice", "Bob", "Charlie"}
	return nil
}

// Initialize the friends model
func NewFriendsModel() friendsModel {
	return friendsModel{
		friends:        []string{"Alice", "Bob", "Charlie"},
		selectedFriend: 0, // Default to the first friend
	}
}

// Update function to handle key inputs for the friend list
func (m friendsModel) Update(msg tea.Msg) (friendsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			// Navigate up the friend list
			if m.selectedFriend > 0 {
				m.selectedFriend--
			}

		case "down":
			// Navigate down the friend list
			if m.selectedFriend < len(m.friends)-1 {
				m.selectedFriend++
			}

		case "enter":
			// Select friend for DM (You might trigger a callback or other action here)
			fmt.Printf("Selected friend: %s\n", m.friends[m.selectedFriend])
		}
	}

	return m, nil
}

// View function renders the friend list UI
func (m friendsModel) View(width, height int) string {
	// Render the friend list with selection highlighting
	friendList := "Friends List:\n"
	for i, friend := range m.friends {
		if i == m.selectedFriend {
			friendList += "> " + friend + " (selected)\n"
		} else {
			friendList += "  " + friend + "\n"
		}
	}

	// Apply the style with dynamic width and height
	return friendListStyle.
		Width(width).
		Height(height).
		Render(friendList)
}

// // Add a new friend to the friend list
// func (m *friendsModel) AddFriend(name string) {
// 	m.friends = append(m.friends, name)
// }

// // Remove a friend from the friend list
// func (m *friendsModel) RemoveFriend(index int) {
// 	if index >= 0 && index < len(m.friends) {
// 		m.friends = append(m.friends[:index], m.friends[index+1:]...)
// 		if m.selectedFriend >= len(m.friends) {
// 			m.selectedFriend = len(m.friends) - 1
// 		}
// 	}
// }
