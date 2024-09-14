// messages.go

package pages

import "github.com/johnkhk/cli_chat_app/genproto/friends"

type FriendListMsg struct {
	Friends []*friends.Friend // Actual Friend type from proto
	Err     error
}

type IncomingFriendRequestsMsg struct {
	Requests []*friends.FriendRequest // Actual FriendRequest type from proto
	Err      error
}

type OutgoingFriendRequestsMsg struct {
	Requests []*friends.FriendRequest // Actual FriendRequest type from proto
	Err      error
}
