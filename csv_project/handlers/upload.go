package handlers

import (
	"net/http"
	"sync"

	"csvproject/utils"
)

// var logger = utils.GetLogger()

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	file := r.Body
	defer file.Close()

	db, err := utils.ConnectDatabase()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	if err := utils.EnsureTableExistsAndTruncate(db); err != nil {
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
		utils.ReadCSV(file, rowChannel, pauseProcessing)
		close(rowChannel)
	}()

	for i := 0; i < utils.MaxWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			utils.ProcessBatches(db, rowChannel, workerID, pauseProcessing)
		}(i + 1)
	}

	wg.Wait()
	close(stopMonitor)
	w.Write([]byte("Upload completed successfully"))
}
