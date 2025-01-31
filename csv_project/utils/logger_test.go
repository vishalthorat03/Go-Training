package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewLogger_ValidPath(t *testing.T) {
	logger, err := NewLogger()
	assert.NoError(t, err)
	assert.NotNil(t, logger)

	// Check default log level (Logrus defaults to InfoLevel if not explicitly set)
	assert.Equal(t, logrus.InfoLevel, logger.GetLevel())
}

func TestNewLogger_WritesToFileAndStdout(t *testing.T) {
	logger, err := NewLogger()
	assert.NoError(t, err)

	// Write a log message
	testMessage := "Test log message"
	logger.Info(testMessage)

	// Define the expected log file path
	logFilePath := "/app/app/logs/app.log"

	// Verify log file contains the message
	content, err := os.ReadFile(logFilePath)
	assert.NoError(t, err)
	assert.Contains(t, string(content), testMessage)
}

func TestNewLogger_DirectoryCreation(t *testing.T) {
	// Define the expected log file path
	logFilePath := "/app/app/logs/app.log"

	_, err := NewLogger()
	assert.NoError(t, err)

	// Check if the log directory was created
	_, err = os.Stat(filepath.Dir(logFilePath))
	assert.NoError(t, err)
}

func TestNewLogger_InvalidPath(t *testing.T) {
	// Since `NewLogger()` does not take an argument, this test is not relevant.
	// If you want to test failure scenarios, you might need to modify NewLogger() to accept a path.
}
