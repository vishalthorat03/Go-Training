package main

import (
	"bytes"
	"csv_project/handlers"
	"database/sql"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// Mock database utilities
func mockConnectDB() (*sql.DB, error) {
	db, _, err := sqlmock.New() // Creates a new mock DB
	if err != nil {
		return nil, err
	}
	return db, nil
}

func mockEnsureTable(db *sql.DB) error {
	return nil // Simulating no errors during table preparation
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
		_ = row // Avoid unused variable error
		// You can log or process each row if needed
	}
}

func TestMain_RootRoute(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "frontend")
	}))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestMain_UploadRoute(t *testing.T) {
	// Wrap the handler with mock dependencies
	server := httptest.NewServer(WrapHandleUpload(
		mockConnectDB,      // Mocking the DB connection
		mockEnsureTable,    // Mocking table setup
		mockReadCSV,        // Mocking CSV reading
		mockProcessBatches, // Mocking batch processing
	))
	defer server.Close()

	// Create a mock request body for the /upload endpoint
	reqBody := bytes.NewReader([]byte("mock,csv,data"))
	resp, err := http.Post(server.URL+"/upload", "application/json", reqBody)
	if err != nil {
		t.Fatalf("Failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestMain_UploadProgress(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handlers.UploadProgress))
	defer server.Close()

	resp, err := http.Get(server.URL + "/upload-progress")
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", resp.Header.Get("Content-Type"))
	}
}
