package handlers

import (
	"bytes"
	"database/sql"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/mock"
)

// Mock database utilities
type MockDB struct {
	mock.Mock
}

func mockReadCSV(file io.Reader, rowChannel chan<- []string, pauseProcessing <-chan bool) {
	// Simulate reading rows and sending them to the channel
	for i := 0; i < 10; i++ {
		rowChannel <- []string{"mock data", strconv.Itoa(i)}
	}
}

func mockProcessBatches(db *sql.DB, rowChannel <-chan []string, workerID int, pauseProcessing <-chan bool) {
	// Simulate batch processing
	for row := range rowChannel {
		log.Printf("Worker %d processing row: %v", workerID, row)
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
	return nil // Simulate successful table preparation
}

func TestHandleUpload(t *testing.T) {
	reqBody := bytes.NewReader([]byte("mock,csv,data"))
	req := httptest.NewRequest(http.MethodPost, "/upload", reqBody)
	w := httptest.NewRecorder()

	HandleUpload(
		w,
		req,
		mockConnectDB,      // Use mock database connection
		mockEnsureTable,    // Mock table setup
		mockReadCSV,        // Mock CSV reading
		mockProcessBatches, // Mock batch processing
	)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}
}
