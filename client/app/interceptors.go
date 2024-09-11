package app

import (
	"context"
	"fmt"
	"strings"

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
func UnaryInterceptor(tokenManager *TokenManager) grpc.UnaryClientInterceptor {
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
