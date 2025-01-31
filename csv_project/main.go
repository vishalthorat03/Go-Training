package main

import (
	"csv_project/handlers"
	"csv_project/utils"
	"database/sql"
	"io"
	"net/http"
)

// WrapHandleUpload wraps HandleUpload for http.HandleFunc.
func WrapHandleUpload(
	connectDB func() (*sql.DB, error),
	ensureTable func(*sql.DB) error,
	readCSV func(io.Reader, chan<- []string, <-chan bool),
	processBatches func(*sql.DB, <-chan []string, int, <-chan bool),
) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleUpload(w, r, connectDB, ensureTable, readCSV, processBatches)
	}
}

func main() {
	logger := utils.GetLogger()

	logger.Info("Starting server on port 4041...")

	// Serve static files
	fs := http.FileServer(http.Dir("frontend"))
	http.Handle("/", fs)

	// Routes
	http.HandleFunc("/upload", WrapHandleUpload(
		utils.ConnectDatabase,
		utils.EnsureTableExistsAndTruncate,
		utils.ReadCSV,
		utils.ProcessBatches,
	))
	http.HandleFunc("/showTable", handlers.ShowTable)

	// Start the server
	err := http.ListenAndServe(":4041", nil)
	if err != nil {
		// Log the error string
		logger.Fatal("Server failed to start: " + err.Error())
	}
}
