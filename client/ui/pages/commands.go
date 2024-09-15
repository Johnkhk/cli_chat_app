// commands.go

package pages

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/johnkhk/cli_chat_app/client/app"
)

// fetchFriendListCmd fetches the friend list for the current user.
func fetchFriendListCmd(rpcClient *app.RpcClient) tea.Cmd {
	return func() tea.Msg {
		friends, err := rpcClient.FriendsClient.GetFriendList()
		return FriendListMsg{Friends: friends, Err: err}
	}
}

// fetchIncomingFriendRequestsCmd fetches the incoming friend requests for the current user.
func fetchIncomingFriendRequestsCmd(rpcClient *app.RpcClient) tea.Cmd {
	return func() tea.Msg {
		requests, err := rpcClient.FriendsClient.GetIncomingFriendRequests()
		return IncomingFriendRequestsMsg{Requests: requests, Err: err}
	}
}

// fetchOutgoingFriendRequestsCmd fetches the outgoing friend requests sent by the current user.
func fetchOutgoingFriendRequestsCmd(rpcClient *app.RpcClient) tea.Cmd {
	return func() tea.Msg {
		requests, err := rpcClient.FriendsClient.GetOutgoingFriendRequests()
		return OutgoingFriendRequestsMsg{Requests: requests, Err: err}
	}
}

// sendFriendRequestCmd sends a friend request to another user and returns a result message.
func sendFriendRequestCmd(rpcClient *app.RpcClient, recipientUsername string) tea.Cmd {
	return func() tea.Msg {
		err := rpcClient.FriendsClient.SendFriendRequest(recipientUsername)
		return SendFriendRequestResultMsg{RecipientUsername: recipientUsername, Err: err}
	}
}

// acceptFriendRequestCmd accepts a friend request and returns a result message.
func acceptFriendRequestCmd(rpcClient *app.RpcClient, requestID int32) tea.Cmd {
	return func() tea.Msg {
		err := rpcClient.FriendsClient.AcceptFriendRequest(requestID)
		return AcceptFriendRequestResultMsg{RequestID: requestID, Err: err}
	}
}

// declineFriendRequestCmd declines a friend request and returns a result message.
func declineFriendRequestCmd(rpcClient *app.RpcClient, requestID int32) tea.Cmd {
	return func() tea.Msg {
		err := rpcClient.FriendsClient.DeclineFriendRequest(requestID)
		return DeclineFriendRequestResultMsg{RequestID: requestID, Err: err}
	}
}

// removeFriendCmd removes a friend from the user's friend list and returns a result message.
func removeFriendCmd(rpcClient *app.RpcClient, friendID int32) tea.Cmd {
	return func() tea.Msg {
		err := rpcClient.FriendsClient.RemoveFriend(friendID)
		return RemoveFriendResultMsg{FriendID: friendID, Err: err}
	}
}
