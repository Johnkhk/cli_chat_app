package logger

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var Log *logrus.Logger

func InitLogger() {
	Log = logrus.New()
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Optional: Output to stdout instead of stderr (default)
	Log.SetOutput(os.Stdout)
	Log.SetLevel(logrus.InfoLevel)
	// Log.SetLevel(logrus.PanicLevel) // Disable all logs below PanicLevel

}

func UnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	// Ensure logger is initialized before use
	if Log == nil {
		logrus.Fatal("Logger is not initialized.")
	}

	// Log incoming RPC call
	Log.Infof("Incoming RPC call: %s", info.FullMethod)

	// Call the handler to complete the normal execution of the RPC call
	resp, err := handler(ctx, req)
	if err != nil {
		Log.Errorf("Error handling RPC: %v", err)
	} else {
		Log.Infof("RPC call completed successfully: %s", info.FullMethod)
	}
	return resp, err
}
