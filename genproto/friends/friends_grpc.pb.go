// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.3
// source: proto/friends/friends.proto

package friends

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	FriendsService_SendFriendRequest_FullMethodName      = "/friends.FriendsService/SendFriendRequest"
	FriendsService_GetFriendRequests_FullMethodName      = "/friends.FriendsService/GetFriendRequests"
	FriendsService_RespondToFriendRequest_FullMethodName = "/friends.FriendsService/RespondToFriendRequest"
	FriendsService_GetFriend_FullMethodName              = "/friends.FriendsService/GetFriend"
	FriendsService_GetFriends_FullMethodName             = "/friends.FriendsService/GetFriends"
	FriendsService_RemoveFriend_FullMethodName           = "/friends.FriendsService/RemoveFriend"
)

// FriendsServiceClient is the client API for FriendsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Define the friends service.
type FriendsServiceClient interface {
	SendFriendRequest(ctx context.Context, in *SendFriendRequestRequest, opts ...grpc.CallOption) (*SendFriendRequestResponse, error)
	GetFriendRequests(ctx context.Context, in *GetFriendRequestsRequest, opts ...grpc.CallOption) (*GetFriendRequestsResponse, error)
	RespondToFriendRequest(ctx context.Context, in *RespondToFriendRequestRequest, opts ...grpc.CallOption) (*RespondToFriendRequestResponse, error)
	GetFriend(ctx context.Context, in *GetFriendRequest, opts ...grpc.CallOption) (*GetFriendResponse, error)
	GetFriends(ctx context.Context, in *GetFriendsRequest, opts ...grpc.CallOption) (*GetFriendsResponse, error)
	RemoveFriend(ctx context.Context, in *RemoveFriendRequest, opts ...grpc.CallOption) (*RemoveFriendResponse, error)
}

type friendsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFriendsServiceClient(cc grpc.ClientConnInterface) FriendsServiceClient {
	return &friendsServiceClient{cc}
}

func (c *friendsServiceClient) SendFriendRequest(ctx context.Context, in *SendFriendRequestRequest, opts ...grpc.CallOption) (*SendFriendRequestResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SendFriendRequestResponse)
	err := c.cc.Invoke(ctx, FriendsService_SendFriendRequest_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *friendsServiceClient) GetFriendRequests(ctx context.Context, in *GetFriendRequestsRequest, opts ...grpc.CallOption) (*GetFriendRequestsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetFriendRequestsResponse)
	err := c.cc.Invoke(ctx, FriendsService_GetFriendRequests_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *friendsServiceClient) RespondToFriendRequest(ctx context.Context, in *RespondToFriendRequestRequest, opts ...grpc.CallOption) (*RespondToFriendRequestResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RespondToFriendRequestResponse)
	err := c.cc.Invoke(ctx, FriendsService_RespondToFriendRequest_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *friendsServiceClient) GetFriend(ctx context.Context, in *GetFriendRequest, opts ...grpc.CallOption) (*GetFriendResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetFriendResponse)
	err := c.cc.Invoke(ctx, FriendsService_GetFriend_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *friendsServiceClient) GetFriends(ctx context.Context, in *GetFriendsRequest, opts ...grpc.CallOption) (*GetFriendsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetFriendsResponse)
	err := c.cc.Invoke(ctx, FriendsService_GetFriends_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *friendsServiceClient) RemoveFriend(ctx context.Context, in *RemoveFriendRequest, opts ...grpc.CallOption) (*RemoveFriendResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RemoveFriendResponse)
	err := c.cc.Invoke(ctx, FriendsService_RemoveFriend_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FriendsServiceServer is the server API for FriendsService service.
// All implementations must embed UnimplementedFriendsServiceServer
// for forward compatibility.
//
// Define the friends service.
type FriendsServiceServer interface {
	SendFriendRequest(context.Context, *SendFriendRequestRequest) (*SendFriendRequestResponse, error)
	GetFriendRequests(context.Context, *GetFriendRequestsRequest) (*GetFriendRequestsResponse, error)
	RespondToFriendRequest(context.Context, *RespondToFriendRequestRequest) (*RespondToFriendRequestResponse, error)
	GetFriend(context.Context, *GetFriendRequest) (*GetFriendResponse, error)
	GetFriends(context.Context, *GetFriendsRequest) (*GetFriendsResponse, error)
	RemoveFriend(context.Context, *RemoveFriendRequest) (*RemoveFriendResponse, error)
	mustEmbedUnimplementedFriendsServiceServer()
}

// UnimplementedFriendsServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedFriendsServiceServer struct{}

func (UnimplementedFriendsServiceServer) SendFriendRequest(context.Context, *SendFriendRequestRequest) (*SendFriendRequestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendFriendRequest not implemented")
}
func (UnimplementedFriendsServiceServer) GetFriendRequests(context.Context, *GetFriendRequestsRequest) (*GetFriendRequestsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFriendRequests not implemented")
}
func (UnimplementedFriendsServiceServer) RespondToFriendRequest(context.Context, *RespondToFriendRequestRequest) (*RespondToFriendRequestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RespondToFriendRequest not implemented")
}
func (UnimplementedFriendsServiceServer) GetFriend(context.Context, *GetFriendRequest) (*GetFriendResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFriend not implemented")
}
func (UnimplementedFriendsServiceServer) GetFriends(context.Context, *GetFriendsRequest) (*GetFriendsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFriends not implemented")
}
func (UnimplementedFriendsServiceServer) RemoveFriend(context.Context, *RemoveFriendRequest) (*RemoveFriendResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveFriend not implemented")
}
func (UnimplementedFriendsServiceServer) mustEmbedUnimplementedFriendsServiceServer() {}
func (UnimplementedFriendsServiceServer) testEmbeddedByValue()                        {}

// UnsafeFriendsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FriendsServiceServer will
// result in compilation errors.
type UnsafeFriendsServiceServer interface {
	mustEmbedUnimplementedFriendsServiceServer()
}

func RegisterFriendsServiceServer(s grpc.ServiceRegistrar, srv FriendsServiceServer) {
	// If the following call pancis, it indicates UnimplementedFriendsServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&FriendsService_ServiceDesc, srv)
}

func _FriendsService_SendFriendRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendFriendRequestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FriendsServiceServer).SendFriendRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FriendsService_SendFriendRequest_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FriendsServiceServer).SendFriendRequest(ctx, req.(*SendFriendRequestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FriendsService_GetFriendRequests_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFriendRequestsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FriendsServiceServer).GetFriendRequests(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FriendsService_GetFriendRequests_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FriendsServiceServer).GetFriendRequests(ctx, req.(*GetFriendRequestsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FriendsService_RespondToFriendRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RespondToFriendRequestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FriendsServiceServer).RespondToFriendRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FriendsService_RespondToFriendRequest_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FriendsServiceServer).RespondToFriendRequest(ctx, req.(*RespondToFriendRequestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FriendsService_GetFriend_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFriendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FriendsServiceServer).GetFriend(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FriendsService_GetFriend_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FriendsServiceServer).GetFriend(ctx, req.(*GetFriendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FriendsService_GetFriends_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFriendsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FriendsServiceServer).GetFriends(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FriendsService_GetFriends_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FriendsServiceServer).GetFriends(ctx, req.(*GetFriendsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FriendsService_RemoveFriend_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveFriendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FriendsServiceServer).RemoveFriend(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FriendsService_RemoveFriend_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FriendsServiceServer).RemoveFriend(ctx, req.(*RemoveFriendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FriendsService_ServiceDesc is the grpc.ServiceDesc for FriendsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FriendsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "friends.FriendsService",
	HandlerType: (*FriendsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendFriendRequest",
			Handler:    _FriendsService_SendFriendRequest_Handler,
		},
		{
			MethodName: "GetFriendRequests",
			Handler:    _FriendsService_GetFriendRequests_Handler,
		},
		{
			MethodName: "RespondToFriendRequest",
			Handler:    _FriendsService_RespondToFriendRequest_Handler,
		},
		{
			MethodName: "GetFriend",
			Handler:    _FriendsService_GetFriend_Handler,
		},
		{
			MethodName: "GetFriends",
			Handler:    _FriendsService_GetFriends_Handler,
		},
		{
			MethodName: "RemoveFriend",
			Handler:    _FriendsService_RemoveFriend_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/friends/friends.proto",
}