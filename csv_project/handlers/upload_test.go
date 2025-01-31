package handlers

import (
	"bytes"
	"database/sql"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"

	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// Mock functions
func mockReadCSV(file io.Reader, rowChannel chan<- []string, pauseProcessing <-chan bool) {
	for i := 0; i < 10; i++ {
		rowChannel <- []string{"mock data", strconv.Itoa(i)}
	}
}

func mockProcessBatches(db *sql.DB, rowChannel <-chan []string, workerID int, pauseProcessing <-chan bool) {
	for range rowChannel {
		// Simulate processing
	}
}

func mockConnectDB() (*sql.DB, error) {
	db, _, err := sqlmock.New()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func mockEnsureTable(db *sql.DB) error {
	return nil
}

func mockFailConnectDB() (*sql.DB, error) {
	return nil, errors.New("mock connection error")
}

func mockFailEnsureTable(db *sql.DB) error {
	return errors.New("mock table error")
}

func TestHandleUpload(t *testing.T) {
	tests := []struct {
		name            string
		method          string
		body            io.Reader
		connectDB       func() (*sql.DB, error)
		ensureTable     func(*sql.DB) error
		expectedStatus  int
		expectedMessage string
	}{
		{
			name:            "Invalid Method",
			method:          http.MethodGet,
			body:            nil,
			connectDB:       mockConnectDB,
			ensureTable:     mockEnsureTable,
			expectedStatus:  http.StatusMethodNotAllowed,
			expectedMessage: "Invalid request method",
		},
		{
			name:            "Database Connection Failure",
			method:          http.MethodPost,
			body:            bytes.NewReader([]byte("mock,csv,data")),
			connectDB:       mockFailConnectDB,
			ensureTable:     mockEnsureTable,
			expectedStatus:  http.StatusInternalServerError,
			expectedMessage: "Failed to connect to database",
		},
		{
			name:            "Table Preparation Failure",
			method:          http.MethodPost,
			body:            bytes.NewReader([]byte("mock,csv,data")),
			connectDB:       mockConnectDB,
			ensureTable:     mockFailEnsureTable,
			expectedStatus:  http.StatusInternalServerError,
			expectedMessage: "Failed to prepare table",
		},
		{
			name:            "Successful Upload",
			method:          http.MethodPost,
			body:            bytes.NewReader([]byte("mock,csv,data")),
			connectDB:       mockConnectDB,
			ensureTable:     mockEnsureTable,
			expectedStatus:  http.StatusOK,
			expectedMessage: "Upload completed successfully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/upload", tt.body)
			w := httptest.NewRecorder()

			HandleUpload(
				w,
				req,
				tt.connectDB,
				tt.ensureTable,
				mockReadCSV,
				mockProcessBatches,
			)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, res.StatusCode)
			}

			body, _ := io.ReadAll(res.Body)
			if !bytes.Contains(body, []byte(tt.expectedMessage)) {
				t.Errorf("expected body to contain %q, got %q", tt.expectedMessage, body)
			}
		})
	}
}
