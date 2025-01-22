package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func NewLogger(logFilePath string) (*logrus.Logger, error) {
	// Create a new logrus logger instance
	log := logrus.New()

	// Get the current working directory (project root)
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %v", err)
	}

	// Ensure the log path is platform-independent (use filepath.Join to handle OS-specific file paths)
	fullLogPath := filepath.Join(wd, logFilePath)

	// Ensure the directory exists before opening the file
	dir := filepath.Dir(fullLogPath) // Extract directory from path
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %v", err)
	}

	// Open or create the log file
	file, err := os.OpenFile(fullLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}

	// Use io.MultiWriter to write to both stdout and the log file
	log.SetOutput(io.MultiWriter(os.Stdout, file))

	// Optionally set log level
	log.SetLevel(logrus.DebugLevel)

	// Set the log format with custom time format (ISO8601-like format)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,                       // Enable full timestamp
		TimestampFormat: "2006-01-02T15:04:05.000Z", // Set the desired timestamp format
		DisableColors:   true,                       // Disable color codes in the logs
	})

	// Log the initialization message
	log.Info("Logger initialized successfully.")

	return log, nil
}

func init() {
	// Initialize the logger with the path to the log file
	logFilePath := "app/app.log" // Path relative to project root
	var err error
	logger, err = NewLogger(logFilePath)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
}

// GetLogger returns the logger instance
func GetLogger() *logrus.Logger {
	return logger
}
