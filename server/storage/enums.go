package storage

import "github.com/johnkhk/cli_chat_app/genproto/friends"

// Status constants for friend request statuses stored in the database
const (
	StatusPending  = "pending"
	StatusAccepted = "accepted"
	StatusDeclined = "declined"
	StatusCanceled = "canceled"
)

// Map to convert string statuses from the database to Protobuf enum values
var StatusMap = map[string]friends.FriendRequestStatus{
	StatusPending:  friends.FriendRequestStatus_PENDING,
	StatusAccepted: friends.FriendRequestStatus_ACCEPTED,
	StatusDeclined: friends.FriendRequestStatus_DECLINED,
	StatusCanceled: friends.FriendRequestStatus_CANCELED,
}
