package logger

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"

	"github.com/johnkhk/cli_chat_app/client/app"
)

func InitLogger() *logrus.Logger {
	log := logrus.New()

	// Set log format
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Optionally set log level
	log.SetLevel(logrus.InfoLevel) // Set to desired level, e.g., Info, Warn, Error

	appDirPath, err := app.GetAppDirPath()
	if err != nil {
		log.Fatalf("Failed to get app directory path: %v", err)
	}

	// Open the log file for writing
	logFilePath := filepath.Join(appDirPath, "debug.log")
	os.MkdirAll(filepath.Dir(logFilePath), 0755)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err == nil {
		log.SetOutput(logFile) // Write logs to debug.log
	} else {
		log.Warn("Failed to log to file, using default stderr")
		log.SetOutput(os.Stderr) // Fallback to stderr if file logging fails
	}

	return log
}
