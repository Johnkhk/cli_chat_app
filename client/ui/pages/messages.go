// messages.go

package pages

import "github.com/johnkhk/cli_chat_app/genproto/friends"

// Data Messages (used to pass data to child models)
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

// Action Messages (sent from child models to parent model to request an action)
type SendFriendRequestMsg struct {
	RecipientUsername string
}

type AcceptFriendRequestMsg struct {
	RequestID string
}

type DeclineFriendRequestMsg struct {
	RequestID string
}

type RemoveFriendMsg struct {
	FriendID int32
}

// Result Messages (returned by commands after action execution)
type SendFriendRequestResultMsg struct {
	RecipientUsername string
	Err               error
}

type AcceptFriendRequestResultMsg struct {
	RequestID string
	Err       error
}

type DeclineFriendRequestResultMsg struct {
	RequestID string
	Err       error
}

type RemoveFriendResultMsg struct {
	FriendID int32
	Err      error
}