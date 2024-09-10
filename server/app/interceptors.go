package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// UnaryServerInterceptor returns a new unary server interceptor for validating tokens.
func UnaryServerInterceptor(tokenValidator TokenValidator, logger *logrus.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Log the incoming request details for debugging
		logger.Infof("Intercepting request: Method=%s, FullMethod=%s", info.FullMethod, info.FullMethod)

		// Skip authentication for certain methods like Login or Register.
		if isUnauthenticatedMethod(info.FullMethod) {
			logger.Infof("Skipping authentication for method: %s", info.FullMethod)
			return handler(ctx, req)
		}

		// Extract the token from the metadata.
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			logger.Error("Missing metadata in request")
			return nil, fmt.Errorf("missing metadata")
		}

		authorization := md["authorization"]
		if len(authorization) == 0 {
			logger.Error("Missing authorization token in metadata")
			return nil, fmt.Errorf("missing authorization token")
		}

		// Split and validate the token.
		tokenParts := strings.SplitN(authorization[0], " ", 2)
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			logger.Errorf("Invalid authorization token format: %v", authorization)
			return nil, fmt.Errorf("invalid authorization token format")
		}
		token := tokenParts[1]

		// Validate the token using the TokenValidator.
		userID, username, err := tokenValidator.ValidateToken(token)
		if err != nil {
			logger.Errorf("Invalid token: %v", err)
			return nil, fmt.Errorf("invalid token: %w", err)
		}

		// Log the successful validation
		logger.Infof("Successfully validated token for user ID: %s, Username: %s", userID, username)

		// Add the user ID and username to the context.
		ctx = context.WithValue(ctx, "userID", userID)
		ctx = context.WithValue(ctx, "username", username)

		// Continue with the request.
		return handler(ctx, req)
	}
}

// isUnauthenticatedMethod checks if a gRPC method does not require authentication.
func isUnauthenticatedMethod(method string) bool {
	unauthenticatedMethods := []string{
		"/auth.AuthService/LoginUser",
		"/auth.AuthService/RegisterUser",
		"/auth.AuthService/RefreshToken",
	}
	for _, m := range unauthenticatedMethods {
		if method == m {
			return true
		}
	}
	return false
}
