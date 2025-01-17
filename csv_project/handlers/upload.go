package handlers

import (
	"csvproject/utils"
	"database/sql"
	"io"
	"net/http"
	"sync"
)

// HandleUpload processes CSV uploads.
func HandleUpload(
	w http.ResponseWriter,
	r *http.Request,
	connectDB func() (*sql.DB, error),
	ensureTable func(*sql.DB) error,
	readCSV func(io.Reader, chan<- []string, <-chan bool),
	processBatches func(*sql.DB, <-chan []string, int, <-chan bool),
) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	file := r.Body
	defer file.Close()

	db, err := connectDB()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	if err := ensureTable(db); err != nil {
		http.Error(w, "Failed to prepare table", http.StatusInternalServerError)
		return
	}

	rowChannel := make(chan []string, utils.BatchSize*2)
	var wg sync.WaitGroup
	stopMonitor := make(chan struct{})
	pauseProcessing := make(chan bool)

	go utils.MonitorResources(pauseProcessing, stopMonitor)

	wg.Add(1)
	go func() {
		defer wg.Done()
		readCSV(file, rowChannel, pauseProcessing)
		close(rowChannel)
	}()

	for i := 0; i < utils.MaxWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			processBatches(db, rowChannel, workerID, pauseProcessing)
		}(i + 1)
	}

	wg.Wait()
	close(stopMonitor)
	w.Write([]byte("Upload completed successfully"))
}
