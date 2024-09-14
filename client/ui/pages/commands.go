// commands.go

package pages

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/johnkhk/cli_chat_app/client/app"
)

func fetchFriendListCmd(rpcClient *app.RpcClient) tea.Cmd {
	return func() tea.Msg {
		friends, err := rpcClient.FriendsClient.GetFriendList()
		return FriendListMsg{Friends: friends, Err: err}
	}
}

func fetchIncomingFriendRequestsCmd(rpcClient *app.RpcClient) tea.Cmd {
	return func() tea.Msg {
		requests, err := rpcClient.FriendsClient.GetIncomingFriendRequests()
		return IncomingFriendRequestsMsg{Requests: requests, Err: err}
	}
}

func fetchOutgoingFriendRequestsCmd(rpcClient *app.RpcClient) tea.Cmd {
	return func() tea.Msg {
		requests, err := rpcClient.FriendsClient.GetOutgoingFriendRequests()
		return OutgoingFriendRequestsMsg{Requests: requests, Err: err}
	}
}
