package pages

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/johnkhk/cli_chat_app/client/app"
	"github.com/johnkhk/cli_chat_app/genproto/friends"
)

// friendListModel represents the model for the friend list page
type friendListModel struct {
	friendList         []*friends.Friend // List of friends
	rpcClient          *app.RpcClient    // Reference to the RPC client
	cursor             int               // Cursor position in the list
	removeConfirmation bool              // Indicates if we're in the remove confirmation state
}

// Init initializes the model (no initialization needed here)
func (m friendListModel) Init() tea.Cmd {
	return nil
}

// Update handles incoming messages and updates the model's state
func (m friendListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Handle friend list message
	case FriendListMsg:
		if msg.Err != nil {
			m.rpcClient.Logger.Errorf("Error fetching friend list: %v", msg.Err)
		} else {
			m.rpcClient.Logger.Infof("Received friend list: %v", msg.Friends)
			m.friendList = msg.Friends
		}

	// Handle key presses
	case tea.KeyMsg:
		if m.removeConfirmation {
			// Handle confirmation inputs
			switch msg.String() {
			case "y", "Y":
				// Confirm removal
				friendID := m.friendList[m.cursor].UserId
				return m, removeFriendCmd(m.rpcClient, friendID)
			case "n", "N":
				// Cancel removal
				m.removeConfirmation = false
			}
		} else {
			switch msg.String() {
			case "up", "k":
				// Move cursor up
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				// Move cursor down
				if m.cursor < len(m.friendList)-1 {
					m.cursor++
				}
			case "r":
				// Initiate remove confirmation
				if len(m.friendList) > 0 {
					m.removeConfirmation = true
				}
			case "ctrl+c", "q":
				// Quit the application
				return m, tea.Quit
			}
		}

	// Handle remove friend result
	case RemoveFriendResultMsg:
		if msg.Err != nil {
			m.rpcClient.Logger.Errorf("Error removing friend: %v", msg.Err)
		} else {
			// Remove friend from the list
			for i, f := range m.friendList {
				if f.UserId == msg.FriendID {
					m.friendList = append(m.friendList[:i], m.friendList[i+1:]...)
					break
				}
			}
			// Adjust cursor if necessary
			if m.cursor >= len(m.friendList) && m.cursor > 0 {
				m.cursor--
			}
			m.removeConfirmation = false
		}
	}

	return m, nil
}

// View renders the UI
func (m friendListModel) View() string {
	if len(m.friendList) == 0 {
		return "You have no friends yet."
	}

	var b strings.Builder

	if m.removeConfirmation {
		// Display confirmation prompt
		friend := m.friendList[m.cursor]
		b.WriteString(fmt.Sprintf("Are you sure you want to remove %s? (y/n)\n", friend.Username))
	} else {
		// Display friend list with cursor
		b.WriteString("Friends:\n\n")
		for i, friend := range m.friendList {
			cursor := " " // No cursor
			if m.cursor == i {
				cursor = ">" // Cursor
			}
			b.WriteString(fmt.Sprintf("%s %s\n", cursor, friend.Username))
		}
		b.WriteString("\nUse ↑/↓ to navigate. Press 'r' to remove the selected friend.")
	}

	return b.String()
}

// NewFriendListModel creates and returns a new friend list model
func NewFriendListModel(rpcClient *app.RpcClient) friendListModel {
	return friendListModel{
		friendList: []*friends.Friend{}, // Initialize with an empty friend list
		rpcClient:  rpcClient,           // Set the RPC client reference
		cursor:     0,                   // Start cursor at the top
	}
}

// // removeFriendCmd creates a command to remove a friend
// func removeFriendCmd(rpcClient *app.RpcClient, friendID int32) tea.Cmd {
// 	return func() tea.Msg {
// 		err := rpcClient.FriendsClient.RemoveFriend(friendID)
// 		return RemoveFriendResultMsg{
// 			FriendID: friendID,
// 			Err:      err,
// 		}
// 	}
// }
