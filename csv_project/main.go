package main

import (
	"csv_project/handlers"
	"csv_project/utils"
	"net/http"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func init() {
	// Use GetLogger instead of InitializeLogger
	logger = utils.GetLogger()
}

func main() {
	logger.Info("Starting server on port 4041...")

	// Serve static files
	fs := http.FileServer(http.Dir("frontend"))
	http.Handle("/", fs)

	// Routes
	http.HandleFunc("/upload", handlers.HandleUpload)
	http.HandleFunc("/upload-progress", handlers.UploadProgress)
	http.HandleFunc("/showTable", handlers.ShowTable)

	// Start the server
	logger.Fatal(http.ListenAndServe(":4041", nil))
}
