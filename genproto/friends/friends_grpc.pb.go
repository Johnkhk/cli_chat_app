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
	FriendManagement_GetFriendList_FullMethodName             = "/friends.FriendManagement/GetFriendList"
	FriendManagement_GetIncomingFriendRequests_FullMethodName = "/friends.FriendManagement/GetIncomingFriendRequests"
	FriendManagement_GetOutgoingFriendRequests_FullMethodName = "/friends.FriendManagement/GetOutgoingFriendRequests"
	FriendManagement_SendFriendRequest_FullMethodName         = "/friends.FriendManagement/SendFriendRequest"
	FriendManagement_AcceptFriendRequest_FullMethodName       = "/friends.FriendManagement/AcceptFriendRequest"
	FriendManagement_DeclineFriendRequest_FullMethodName      = "/friends.FriendManagement/DeclineFriendRequest"
	FriendManagement_RemoveFriend_FullMethodName              = "/friends.FriendManagement/RemoveFriend"
)

// FriendManagementClient is the client API for FriendManagement service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Service definition for Friend Management
type FriendManagementClient interface {
	GetFriendList(ctx context.Context, in *GetFriendListRequest, opts ...grpc.CallOption) (*GetFriendListResponse, error)
	GetIncomingFriendRequests(ctx context.Context, in *GetIncomingFriendRequestsRequest, opts ...grpc.CallOption) (*GetIncomingFriendRequestsResponse, error)
	GetOutgoingFriendRequests(ctx context.Context, in *GetOutgoingFriendRequestsRequest, opts ...grpc.CallOption) (*GetOutgoingFriendRequestsResponse, error)
	SendFriendRequest(ctx context.Context, in *SendFriendRequestRequest, opts ...grpc.CallOption) (*SendFriendRequestResponse, error)
	AcceptFriendRequest(ctx context.Context, in *AcceptFriendRequestRequest, opts ...grpc.CallOption) (*AcceptFriendRequestResponse, error)
	DeclineFriendRequest(ctx context.Context, in *DeclineFriendRequestRequest, opts ...grpc.CallOption) (*DeclineFriendRequestResponse, error)
	RemoveFriend(ctx context.Context, in *RemoveFriendRequest, opts ...grpc.CallOption) (*RemoveFriendResponse, error)
}

type friendManagementClient struct {
	cc grpc.ClientConnInterface
}

func NewFriendManagementClient(cc grpc.ClientConnInterface) FriendManagementClient {
	return &friendManagementClient{cc}
}

func (c *friendManagementClient) GetFriendList(ctx context.Context, in *GetFriendListRequest, opts ...grpc.CallOption) (*GetFriendListResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetFriendListResponse)
	err := c.cc.Invoke(ctx, FriendManagement_GetFriendList_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *friendManagementClient) GetIncomingFriendRequests(ctx context.Context, in *GetIncomingFriendRequestsRequest, opts ...grpc.CallOption) (*GetIncomingFriendRequestsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetIncomingFriendRequestsResponse)
	err := c.cc.Invoke(ctx, FriendManagement_GetIncomingFriendRequests_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *friendManagementClient) GetOutgoingFriendRequests(ctx context.Context, in *GetOutgoingFriendRequestsRequest, opts ...grpc.CallOption) (*GetOutgoingFriendRequestsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetOutgoingFriendRequestsResponse)
	err := c.cc.Invoke(ctx, FriendManagement_GetOutgoingFriendRequests_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *friendManagementClient) SendFriendRequest(ctx context.Context, in *SendFriendRequestRequest, opts ...grpc.CallOption) (*SendFriendRequestResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SendFriendRequestResponse)
	err := c.cc.Invoke(ctx, FriendManagement_SendFriendRequest_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *friendManagementClient) AcceptFriendRequest(ctx context.Context, in *AcceptFriendRequestRequest, opts ...grpc.CallOption) (*AcceptFriendRequestResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AcceptFriendRequestResponse)
	err := c.cc.Invoke(ctx, FriendManagement_AcceptFriendRequest_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *friendManagementClient) DeclineFriendRequest(ctx context.Context, in *DeclineFriendRequestRequest, opts ...grpc.CallOption) (*DeclineFriendRequestResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeclineFriendRequestResponse)
	err := c.cc.Invoke(ctx, FriendManagement_DeclineFriendRequest_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *friendManagementClient) RemoveFriend(ctx context.Context, in *RemoveFriendRequest, opts ...grpc.CallOption) (*RemoveFriendResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RemoveFriendResponse)
	err := c.cc.Invoke(ctx, FriendManagement_RemoveFriend_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FriendManagementServer is the server API for FriendManagement service.
// All implementations must embed UnimplementedFriendManagementServer
// for forward compatibility.
//
// Service definition for Friend Management
type FriendManagementServer interface {
	GetFriendList(context.Context, *GetFriendListRequest) (*GetFriendListResponse, error)
	GetIncomingFriendRequests(context.Context, *GetIncomingFriendRequestsRequest) (*GetIncomingFriendRequestsResponse, error)
	GetOutgoingFriendRequests(context.Context, *GetOutgoingFriendRequestsRequest) (*GetOutgoingFriendRequestsResponse, error)
	SendFriendRequest(context.Context, *SendFriendRequestRequest) (*SendFriendRequestResponse, error)
	AcceptFriendRequest(context.Context, *AcceptFriendRequestRequest) (*AcceptFriendRequestResponse, error)
	DeclineFriendRequest(context.Context, *DeclineFriendRequestRequest) (*DeclineFriendRequestResponse, error)
	RemoveFriend(context.Context, *RemoveFriendRequest) (*RemoveFriendResponse, error)
	mustEmbedUnimplementedFriendManagementServer()
}

// UnimplementedFriendManagementServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedFriendManagementServer struct{}

func (UnimplementedFriendManagementServer) GetFriendList(context.Context, *GetFriendListRequest) (*GetFriendListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFriendList not implemented")
}
func (UnimplementedFriendManagementServer) GetIncomingFriendRequests(context.Context, *GetIncomingFriendRequestsRequest) (*GetIncomingFriendRequestsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetIncomingFriendRequests not implemented")
}
func (UnimplementedFriendManagementServer) GetOutgoingFriendRequests(context.Context, *GetOutgoingFriendRequestsRequest) (*GetOutgoingFriendRequestsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOutgoingFriendRequests not implemented")
}
func (UnimplementedFriendManagementServer) SendFriendRequest(context.Context, *SendFriendRequestRequest) (*SendFriendRequestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendFriendRequest not implemented")
}
func (UnimplementedFriendManagementServer) AcceptFriendRequest(context.Context, *AcceptFriendRequestRequest) (*AcceptFriendRequestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AcceptFriendRequest not implemented")
}
func (UnimplementedFriendManagementServer) DeclineFriendRequest(context.Context, *DeclineFriendRequestRequest) (*DeclineFriendRequestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeclineFriendRequest not implemented")
}
func (UnimplementedFriendManagementServer) RemoveFriend(context.Context, *RemoveFriendRequest) (*RemoveFriendResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveFriend not implemented")
}
func (UnimplementedFriendManagementServer) mustEmbedUnimplementedFriendManagementServer() {}
func (UnimplementedFriendManagementServer) testEmbeddedByValue()                          {}

// UnsafeFriendManagementServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FriendManagementServer will
// result in compilation errors.
type UnsafeFriendManagementServer interface {
	mustEmbedUnimplementedFriendManagementServer()
}

func RegisterFriendManagementServer(s grpc.ServiceRegistrar, srv FriendManagementServer) {
	// If the following call pancis, it indicates UnimplementedFriendManagementServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&FriendManagement_ServiceDesc, srv)
}

func _FriendManagement_GetFriendList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFriendListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FriendManagementServer).GetFriendList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FriendManagement_GetFriendList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FriendManagementServer).GetFriendList(ctx, req.(*GetFriendListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FriendManagement_GetIncomingFriendRequests_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetIncomingFriendRequestsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FriendManagementServer).GetIncomingFriendRequests(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FriendManagement_GetIncomingFriendRequests_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FriendManagementServer).GetIncomingFriendRequests(ctx, req.(*GetIncomingFriendRequestsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FriendManagement_GetOutgoingFriendRequests_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetOutgoingFriendRequestsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FriendManagementServer).GetOutgoingFriendRequests(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FriendManagement_GetOutgoingFriendRequests_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FriendManagementServer).GetOutgoingFriendRequests(ctx, req.(*GetOutgoingFriendRequestsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FriendManagement_SendFriendRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendFriendRequestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FriendManagementServer).SendFriendRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FriendManagement_SendFriendRequest_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FriendManagementServer).SendFriendRequest(ctx, req.(*SendFriendRequestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FriendManagement_AcceptFriendRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AcceptFriendRequestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FriendManagementServer).AcceptFriendRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FriendManagement_AcceptFriendRequest_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FriendManagementServer).AcceptFriendRequest(ctx, req.(*AcceptFriendRequestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FriendManagement_DeclineFriendRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeclineFriendRequestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FriendManagementServer).DeclineFriendRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FriendManagement_DeclineFriendRequest_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FriendManagementServer).DeclineFriendRequest(ctx, req.(*DeclineFriendRequestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FriendManagement_RemoveFriend_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveFriendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FriendManagementServer).RemoveFriend(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FriendManagement_RemoveFriend_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FriendManagementServer).RemoveFriend(ctx, req.(*RemoveFriendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FriendManagement_ServiceDesc is the grpc.ServiceDesc for FriendManagement service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FriendManagement_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "friends.FriendManagement",
	HandlerType: (*FriendManagementServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetFriendList",
			Handler:    _FriendManagement_GetFriendList_Handler,
		},
		{
			MethodName: "GetIncomingFriendRequests",
			Handler:    _FriendManagement_GetIncomingFriendRequests_Handler,
		},
		{
			MethodName: "GetOutgoingFriendRequests",
			Handler:    _FriendManagement_GetOutgoingFriendRequests_Handler,
		},
		{
			MethodName: "SendFriendRequest",
			Handler:    _FriendManagement_SendFriendRequest_Handler,
		},
		{
			MethodName: "AcceptFriendRequest",
			Handler:    _FriendManagement_AcceptFriendRequest_Handler,
		},
		{
			MethodName: "DeclineFriendRequest",
			Handler:    _FriendManagement_DeclineFriendRequest_Handler,
		},
		{
			MethodName: "RemoveFriend",
			Handler:    _FriendManagement_RemoveFriend_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/friends/friends.proto",
}
