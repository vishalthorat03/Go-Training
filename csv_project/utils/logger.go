package utils

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()

	// Open log file for appending
	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Fatalf("Failed to open log file: %v", err)
	}

	// Set the log file as the output
	logger.Out = file

	// Set the log level
	logger.SetLevel(logrus.InfoLevel)

	// Set a custom log formatter
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Optional: Add a hook to write logs to both file and stdout (console)
	logger.AddHook(&logrusHook{file})
}

// Custom hook to duplicate log entries to a file and stdout
type logrusHook struct {
	file *os.File
}

func (h *logrusHook) Levels() []logrus.Level {
	return logrus.AllLevels // Hook for all log levels
}

func (h *logrusHook) Fire(entry *logrus.Entry) error {
	// Get the log message string and handle any potential errors
	logMessage, err := entry.String()
	if err != nil {
		return err
	}

	// Write the log message to the file
	_, err = h.file.WriteString(logMessage)
	if err != nil {
		return err
	}

	// log to stdout here as well
	fmt.Print(logMessage) // Uncomment to log to stdout

	return nil
}

func GetLogger() *logrus.Logger {
	return logger
}
