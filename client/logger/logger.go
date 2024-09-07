package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// InitLogger initializes a new logger instance.
func InitLogger() *logrus.Logger {
	log := logrus.New()

	// Set log output to standard output (console)
	log.SetOutput(os.Stdout)

	// Set log format
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Optionally set log level
	log.SetLevel(logrus.InfoLevel)
	// log.SetLevel(logrus.PanicLevel)

	// Optionally log to file
	// file, err := os.OpenFile("client.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	// if err != nil {
	// 	log.Fatalf("Failed to open log file: %v", err)
	// }
	// log.SetOutput(file)

	return log
}
