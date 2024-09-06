package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLogger() {
	Log = logrus.New()

	// Set log output to standard output (console)
	Log.SetOutput(os.Stdout)

	// Set log format
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Optionally set log level
	Log.SetLevel(logrus.InfoLevel)
	// Log.SetLevel(logrus.PanicLevel)

	// Optionally log to file
	// file, err := os.OpenFile("client.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	// if err != nil {
	// 	Log.Fatalf("Failed to open log file: %v", err)
	// }
	// Log.SetOutput(file)
}
