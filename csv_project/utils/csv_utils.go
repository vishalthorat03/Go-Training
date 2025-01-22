package utils

import (
	"bufio"
	"database/sql"
	"io"
	"strconv"
	"strings"

	"github.com/lib/pq"
)

const (
	BatchSize  = 10000
	MaxWorkers = 15
)

func ReadCSV(file io.Reader, rowChannel chan<- []string, pauseProcessing <-chan bool) {
	scanner := bufio.NewScanner(file)
	buffer := make([]byte, 10*1024*1024) // 10MB buffer
	scanner.Buffer(buffer, len(buffer))

	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			logger.Errorf("Error reading CSV header: %v", err)
		}
		return
	}

	var paused bool

	for scanner.Scan() {
		select {
		case paused = <-pauseProcessing:
			if paused {
				logger.Info("ReadCSV: Paused due to resource constraints")
				for paused {
					paused = <-pauseProcessing
				}
				logger.Info("ReadCSV: Resumed processing")
			}
		default:
			row := strings.Split(scanner.Text(), ",")
			rowChannel <- row
		}
	}

	close(rowChannel)

	if err := scanner.Err(); err != nil {
		logger.Errorf("Error reading CSV: %v", err)
	}
}

func ProcessBatches(db *sql.DB, rowChannel <-chan []string, workerID int, pauseProcessing <-chan bool) {
	var rows [][]string
	var paused bool

	for {
		select {
		case paused = <-pauseProcessing:
			if paused {
				logger.Infof("Worker %d: Paused due to resource constraints", workerID)
				rows = nil // Clear buffer to release memory
				for paused {
					paused = <-pauseProcessing
				}
				logger.Infof("Worker %d: Resumed processing", workerID)
			}
		case row, ok := <-rowChannel:
			if !ok {
				if len(rows) > 0 {
					insertBatch(db, rows, workerID)
				}
				return
			}

			rows = append(rows, row)
			if len(rows) >= BatchSize {
				insertBatch(db, rows, workerID)
				rows = nil // Clear memory after batch insert
			}
		}
	}
}

func insertBatch(db *sql.DB, rows [][]string, workerID int) {
	if len(rows) == 0 {
		return
	}

	// Start COPY operation
	txn, err := db.Begin()
	if err != nil {
		logger.Errorf("Worker %d: Failed to begin transaction: %v", workerID, err)
		return
	}
	defer txn.Rollback()

	copyStmt := pq.CopyIn("devices", "devicename", "devicetype", "brand", "model", "os", "osversion", "purchasedate", "warrantyend", "status", "price")
	stmt, err := txn.Prepare(copyStmt)
	if err != nil {
		logger.Errorf("Worker %d: Failed to prepare COPY statement: %v", workerID, err)
		return
	}
	defer stmt.Close()

	for _, row := range rows {
		if len(row) != 11 {
			logger.Warnf("Worker %d: Skipping invalid row: %+v", workerID, row)
			continue
		}

		price, err := strconv.ParseFloat(strings.TrimSpace(row[10]), 64)
		if err != nil {
			logger.Warnf("Worker %d: Invalid price format for row: %+v", workerID, row)
			continue
		}

		_, err = stmt.Exec(row[1], row[2], row[3], row[4], row[5], row[6], row[7], row[8], row[9], price)
		if err != nil {
			logger.Errorf("Worker %d: Failed to execute COPY", workerID)
			return
		}
	}

	if _, err := stmt.Exec(); err != nil {
		logger.Errorf("Worker %d: Failed to finalize COPY: %v", workerID, err)
		return
	}

	if err := txn.Commit(); err != nil {
		logger.Errorf("Worker %d: Failed to commit transaction: %v", workerID, err)
	}
}
