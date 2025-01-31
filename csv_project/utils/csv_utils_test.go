package utils

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func mockReadCSV(_ *strings.Reader, rowChannel chan<- []string, _ <-chan bool) {
	for i := 0; i < 5; i++ {
		rowChannel <- []string{"mock device", strconv.Itoa(i)} // Corrected conversion
	}
	close(rowChannel)
}

func mockProcessBatches(_ *sql.DB, rowChannel <-chan []string, _ int, _ <-chan bool) {
	for range rowChannel {
		// Simulate batch processing without interacting with DB.
	}
}

// Mock helpers for DB interactions
func mockEnsureTable(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS devices (id SERIAL PRIMARY KEY, devicename TEXT, price NUMERIC)")
	return err
}

func mockFailEnsureTable(db *sql.DB) error {
	return errors.New("failed to ensure table")
}

func TestReadCSVAndProcessBatches(t *testing.T) {
	tests := []struct {
		name            string
		mockDB          func() (*sql.DB, sqlmock.Sqlmock, error)
		ensureTable     func(*sql.DB) error
		expectedError   string
		expectTableCall bool
		readCSV         func(*strings.Reader, chan<- []string, <-chan bool)
		processBatches  func(*sql.DB, <-chan []string, int, <-chan bool)
	}{
		{
			name: "Successful batch processing",
			mockDB: func() (*sql.DB, sqlmock.Sqlmock, error) {
				db, mock, err := sqlmock.New()
				if err != nil {
					return nil, nil, err
				}
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS devices").WillReturnResult(sqlmock.NewResult(0, 1))
				return db, mock, nil
			},
			ensureTable:     mockEnsureTable,
			expectedError:   "",
			expectTableCall: true,
			readCSV:         mockReadCSV,
			processBatches:  mockProcessBatches,
		},
		{
			name: "Failed to ensure table",
			mockDB: func() (*sql.DB, sqlmock.Sqlmock, error) {
				db, _, err := sqlmock.New()
				return db, nil, err
			},
			ensureTable:     mockFailEnsureTable,
			expectedError:   "failed to ensure table",
			expectTableCall: false,
			readCSV:         mockReadCSV,
			processBatches:  mockProcessBatches,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock database connection
			db, mock, err := tt.mockDB()
			if tt.expectedError != "" && err != nil {
				assert.EqualError(t, err, tt.expectedError)
				return
			}
			assert.NoError(t, err)
			if db != nil {
				defer db.Close()
			}

			rowChannel := make(chan []string, 5)
			pauseProcessing := make(chan bool)

			// Ensure table setup
			err = tt.ensureTable(db)
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
				return
			}
			assert.NoError(t, err)

			// Mock CSV reading
			go tt.readCSV(strings.NewReader("mock,csv,data"), rowChannel, pauseProcessing)

			// Process batches
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				tt.processBatches(db, rowChannel, 1, pauseProcessing)
			}()
			wg.Wait()

			// Verify expectations
			if mock != nil {
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			}
		})
	}
}
