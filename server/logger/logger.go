package logger

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// InitLogger initializes a new logger instance and returns it.
func InitLogger() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Optional: Output to stdout instead of stderr (default)
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)
	// log.SetLevel(logrus.PanicLevel) // Disable all logs below PanicLevel

	return log
}

// UnaryInterceptor is a gRPC interceptor that logs incoming RPC calls.
func UnaryInterceptor(
	log *logrus.Logger, // Pass logger as a parameter
) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Log incoming RPC call
		log.Infof("Incoming RPC call: %s", info.FullMethod)

		// Call the handler to complete the normal execution of the RPC call
		resp, err := handler(ctx, req)
		if err != nil {
			log.Errorf("Error handling RPC: %v", err)
		} else {
			log.Infof("RPC call completed successfully: %s", info.FullMethod)
		}
		return resp, err
	}
}
