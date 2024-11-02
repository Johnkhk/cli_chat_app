package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// publicMethods lists the gRPC methods that do NOT require authentication.
var publicMethods = []string{
	"/auth.AuthService/LoginUser",
	"/auth.AuthService/RegisterUser",
	"/auth.AuthService/RefreshToken",
}

// UnaryInterceptor returns a gRPC interceptor that adds the authorization token to each request,
// except for the methods listed in AuthMethods.
func UnaryInterceptor(tokenManager *TokenManager, logger *logrus.Logger) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		// Check if the method requires authentication
		if !isPublicMethod(method) {
			// Retrieve the access token
			// This also attempts to refresh the token if it has expired
			token, err := tokenManager.GetAccessToken()
			if err != nil {
				return fmt.Errorf("failed to get access token: %v", err)
			}
			logger.Info("Adding authorization token to request")
			// Create a new context with the authorization metadata
			md := metadata.Pairs("authorization", "Bearer "+token)
			ctx = metadata.NewOutgoingContext(ctx, md)
		}

		// Invoke the RPC call with the (potentially) modified context
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// isPublicMethod checks if a gRPC method is an authentication method.
func isPublicMethod(method string) bool {
	for _, m := range publicMethods {
		if strings.HasPrefix(method, m) {
			return true
		}
	}
	return false
}

// StreamInterceptor returns a gRPC stream interceptor that adds the authorization token to each stream request,
// except for the methods listed in publicMethods.
func StreamInterceptor(tokenManager *TokenManager, logger *logrus.Logger) grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		// Check if the method requires authentication
		if !isPublicMethod(method) {
			// Retrieve the access token
			token, err := tokenManager.GetAccessToken()
			if err != nil {
				return nil, fmt.Errorf("failed to get access token: %v", err)
			}
			logger.Info("Adding authorization token to stream")
			// Create a new context with the authorization metadata
			md := metadata.Pairs("authorization", "Bearer "+token)
			ctx = metadata.NewOutgoingContext(ctx, md)
		}

		// Invoke the RPC call with the (potentially) modified context
		return streamer(ctx, desc, cc, method, opts...)
	}
}
