package storage

import "github.com/johnkhk/cli_chat_app/genproto/friends"

var (
	StatusPendingStr  = friends.FriendRequestStatus_name[int32(friends.FriendRequestStatus_PENDING)]
	StatusAcceptedStr = friends.FriendRequestStatus_name[int32(friends.FriendRequestStatus_ACCEPTED)]
	StatusDeclinedStr = friends.FriendRequestStatus_name[int32(friends.FriendRequestStatus_DECLINED)]
	StatusCanceledStr = friends.FriendRequestStatus_name[int32(friends.FriendRequestStatus_CANCELED)]
	StatusUnknownStr  = friends.FriendRequestStatus_name[int32(friends.FriendRequestStatus_UNKNOWN)]
	StatusFailedStr   = friends.FriendRequestStatus_name[int32(friends.FriendRequestStatus_FAILED)]
)

const (
	StatusPendingInt32  = int32(friends.FriendRequestStatus_PENDING)
	StatusAcceptedInt32 = int32(friends.FriendRequestStatus_ACCEPTED)
	StatusDeclinedInt32 = int32(friends.FriendRequestStatus_DECLINED)
	StatusCanceledInt32 = int32(friends.FriendRequestStatus_CANCELED)
	StatusUnknownInt32  = int32(friends.FriendRequestStatus_UNKNOWN)
	StatusFailedInt32   = int32(friends.FriendRequestStatus_FAILED)
)
