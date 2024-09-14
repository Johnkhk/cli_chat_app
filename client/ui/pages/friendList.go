package pages

import (
	"fmt"
	"strings"

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
			m.rpcClient.Logger.Errorf("Error fetching friend list: %v", msg.Err)
		} else {
			m.rpcClient.Logger.Infof("Received friend list: %v", msg.Friends)
			m.friendList = msg.Friends
		}
		// case tea.KeyMsg:
		// 	switch msg.String() {
		// 	case "d":
		// 		selectedFriendID := ... // Logic to get selected friend
		// 		cmd = func() tea.Msg {
		// 			return RemoveFriendMsg{FriendID: selectedFriendID}
		// 		}
		// 		return m, cmd
		// 	}

	}

	return m, nil
}

func (m friendListModel) View() string {
	if len(m.friendList) == 0 {
		return "You have no friends yet."
	}

	var b strings.Builder
	b.WriteString("Friends:\n")
	for _, friend := range m.friendList {
		b.WriteString(fmt.Sprintf("- %s\n", friend.Username)) // Adjust field as necessary
	}
	return b.String()
}

// NewFriendListModel function to create and return a new friend list model
func NewFriendListModel(rpcClient *app.RpcClient) friendListModel {
	return friendListModel{
		friendList: []*friends.Friend{}, // Initialize with an empty friend list
		rpcClient:  rpcClient,           // Set the RPC client reference
	}
}
