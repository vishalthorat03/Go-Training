package utils

import (
	"bufio"
	"database/sql"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	BatchSize  = 20000
	MaxWorkers = 30
)

// var logger = GetLogger()

func ReadCSV(file io.Reader, rowChannel chan<- []string, pauseProcessing <-chan bool) {
	scanner := bufio.NewScanner(file)
	buffer := make([]byte, 10*1024*1024) // 10MB buffer
	scanner.Buffer(buffer, len(buffer))

	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			logger.Error("Error reading CSV header:", err)
		}
		return
	}

	for scanner.Scan() {
		select {
		case <-pauseProcessing:
			time.Sleep(10 * time.Second)
		default:
			row := strings.Split(scanner.Text(), ",")
			rowChannel <- row
		}
	}

	if err := scanner.Err(); err != nil {
		logger.Error("Error reading CSV:", err)
	}
}

// Exported ProcessBatches function
func ProcessBatches(db *sql.DB, rowChannel <-chan []string, workerID int, pauseProcessing <-chan bool) {
	var rows [][]string

	for {
		select {
		case <-pauseProcessing:
			log.Printf("Worker %d: Pausing due to resource constraints", workerID)
			time.Sleep(5 * time.Second)
		case row, ok := <-rowChannel:
			if !ok {
				if len(rows) > 0 {
					insertBatch(db, rows, workerID) // Assuming insertBatch is defined somewhere
				}
				return
			}

			rows = append(rows, row)
			if len(rows) >= BatchSize {
				insertBatch(db, rows, workerID)
				rows = nil
			}
		}
	}
}

func insertBatch(db *sql.DB, rows [][]string, workerID int) {
	if len(rows) == 0 {
		return
	}

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		logger.Errorf("Worker %d: Failed to begin transaction: %v", workerID, err)
		return
	}
	defer tx.Rollback() // Ensure the transaction is rolled back on failure

	// Prepare the statement for batch insert
	stmt, err := tx.Prepare(`
        INSERT INTO devices (
            devicename, devicetype, brand, model, os, osversion, purchasedate,
            warrantyend, status, price
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `)
	if err != nil {
		logger.Errorf("Worker %d: Failed to prepare statement: %v", workerID, err)
		return
	}
	defer stmt.Close()

	// Insert each row in the batch using the prepared statement
	for i, row := range rows {
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
			logger.Errorf("Worker %d: Failed to execute batch insert for row %d: %v", workerID, i, err)
			return
		}
	}

	// Commit the transaction after the batch insert
	err = tx.Commit()
	if err != nil {
		logger.Errorf("Worker %d: Failed to commit transaction: %v", workerID, err)
	}
}
