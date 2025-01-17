package utils

import (
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func init() {
	// Initialize the logger
	logger = logrus.New()

	// Get the log file path from the environment variable or use a default
	logFilePath := os.Getenv("LOG_FILE_PATH")
	if logFilePath == "" {
		logFilePath = "./app.log" // Default path for non-Docker environments
	}

	// Ensure the directory exists
	dir := filepath.Dir(logFilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		logger.Fatalf("Failed to create log directory: %v", err)
	}

	// Open the log file for appending
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Fatalf("Failed to open log file: %v", err)
	}

	// Use io.MultiWriter to write to both the file and stdout (console)
	logger.Out = io.MultiWriter(file, os.Stdout)

	// Set the log level (can be adjusted to Debug, Error, etc.)
	logger.SetLevel(logrus.InfoLevel)

	// Set a custom log formatter (optional, adds timestamps and log formatting)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

// GetLogger returns the logger instance
func GetLogger() *logrus.Logger {
	return logger
}
