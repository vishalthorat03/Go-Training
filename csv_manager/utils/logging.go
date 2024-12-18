package utils

import (
	"log"
	"os"
)

var Logger *log.Logger

// SetupLogger sets up logging to a file
func SetupLogger() {
	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to set up logger: %v", err)
	}
	Logger = log.New(file, "INFO: ", log.LstdFlags)
}
