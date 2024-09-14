package pages

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/johnkhk/cli_chat_app/client/app"
	"github.com/johnkhk/cli_chat_app/genproto/friends"
)

// Define the friendListModel struct
type friendListModel struct {
	friendList []*friends.Friend // List of friends
	rpcClient  *app.RpcClient    // Reference to the RPC client
}

// Init function to initialize the model (can be empty if no initialization is needed)
func (m friendListModel) Init() tea.Cmd {
	return nil
}

// Update function to handle messages and update the model's state
func (m friendListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case FriendListMsg:
		if msg.Err != nil {
			// Handle error (e.g., log it or display a message)
			m.rpcClient.Logger.Errorf("Error fetching friend list: %v", msg.Err)
		} else {
			m.rpcClient.Logger.Infof("Received friend list: %v", msg.Friends)
			m.friendList = msg.Friends
		}
	}
	return m, nil
}

// View function to render the model
func (m friendListModel) View() string {
	// Convert the friend list to a single string with each friend on a new line
	// return fmt.Sprintf("Friends:\n%s", strings.Join(m.friendList, "\n"))
	return "Friend List View" // Placeholder for actual rendering logic
}

// NewFriendListModel function to create and return a new friend list model
func NewFriendListModel(rpcClient *app.RpcClient) friendListModel {
	return friendListModel{
		friendList: []*friends.Friend{}, // Initialize with an empty friend list
		rpcClient:  rpcClient,           // Set the RPC client reference
	}
}
