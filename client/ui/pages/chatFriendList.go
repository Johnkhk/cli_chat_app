package pages

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/johnkhk/cli_chat_app/client/app"
	"github.com/johnkhk/cli_chat_app/genproto/friends"
)

// ChatFriendListModel manages the friends list within the chat context.
type ChatFriendListModel struct {
	rpcClient *app.RpcClient
	friends   []*friends.Friend // Holds the list of friends
	selected  int               // Currently selected index in the friend list
	loading   bool              // Indicates whether the friend list is being fetched
}

// NewChatFriendListModel initializes the ChatFriendListModel.
func NewChatFriendListModel(rpcClient *app.RpcClient) ChatFriendListModel {
	return ChatFriendListModel{
		rpcClient: rpcClient,
		friends:   []*friends.Friend{},
		selected:  0,    // Default to the first friend
		loading:   true, // Start in a loading state until data is fetched
	}
}

// Update handles key presses to navigate the friend list and processes incoming messages.
func (m ChatFriendListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle navigation with arrow keys
		switch msg.String() {
		case "up":
			if m.selected > 0 {
				m.selected--
			}
		case "down":
			if m.selected < len(m.friends)-1 {
				m.selected++
			}
		case "enter":
			// Trigger a user selection event by pressing "enter".
			if m.selected >= 0 && m.selected < len(m.friends) {
				selectedUserID := m.friends[m.selected].UserId
				selectedUsername := m.friends[m.selected].Username
				m.rpcClient.Logger.Infof("Selected User ID: %d", selectedUserID)
				return m, func() tea.Msg {
					return FriendSelectedMsg{UserID: selectedUserID, Username: selectedUsername}
				}
			}
		}
	case FriendListMsg:
		// Update the friend list once the data is fetched
		if msg.Err == nil {
			m.friends = msg.Friends
			m.loading = false
			if len(m.friends) > 0 {
				defaultUserID := m.friends[m.selected].UserId
				defaultUsername := m.friends[m.selected].Username
				return m, func() tea.Msg {
					return FriendSelectedMsg{UserID: defaultUserID, Username: defaultUsername}
				}
			}
		} else {
			// m.friends = []string{"Failed to load friends."}
			m.rpcClient.Logger.Errorf("Error fetching friend list: %v", msg.Err)
			m.loading = false

		}
		// return m, nil
	}

	return m, nil
}

// View renders the chat friend list with the currently selected friend highlighted.
func (m ChatFriendListModel) View() string {
	if m.loading {
		return "Loading friends..."
	}

	view := ""
	for i, friend := range m.friends {
		cursor := " " // No cursor
		if i == m.selected {
			cursor = ">" // Show cursor on selected item
		}
		view += cursor + " " + friend.Username + "\n"
	}
	return view
}

// Init initializes the ChatFriendListModel with a command to fetch the friend list.
func (m ChatFriendListModel) Init() tea.Cmd {
	return fetchFriendListCmd(m.rpcClient)
}
