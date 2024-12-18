// package main

// import (
// 	"bufio"
// 	"database/sql"
// 	"encoding/json"
// 	"fmt"
// 	"html/template"
// 	"io"
// 	"log"
// 	"net/http"
// 	"os"
// 	"sync"

// 	_ "github.com/lib/pq" // PostgreSQL driver
// )

// type LogEntry struct {
// 	Timestamp string `json:"timestamp"`
// 	Level     string `json:"level"`
// 	Message   string `json:"message"`
// }

// var db *sql.DB
// var mutex sync.Mutex

// func initDB() (*sql.DB, error) {
// 	dsn := "host=postgres user=log_user password=log_password dbname=log_db port=5432 sslmode=disable" // Update to port 5432
// 	db, err := sql.Open("postgres", dsn)
// 	if err != nil {
// 		log.Fatalf("Failed to connect to the database: %v", err)
// 		return nil, err
// 	}

// 	if err := db.Ping(); err != nil {
// 		log.Fatalf("Failed to ping database: %v", err)
// 		return nil, err
// 	}

// 	log.Println("Database connected successfully")
// 	return db, nil
// }

// func uploadHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == "POST" {
// 		file, handler, err := r.FormFile("file")
// 		if err != nil {
// 			http.Error(w, "Failed to upload file", http.StatusBadRequest)
// 			return
// 		}
// 		defer file.Close()

// 		// Save uploaded file
// 		filePath := "./" + handler.Filename
// 		out, err := os.Create(filePath)
// 		if err != nil {
// 			http.Error(w, "Failed to save file", http.StatusInternalServerError)
// 			return
// 		}
// 		defer out.Close()

// 		_, err = io.Copy(out, file)
// 		if err != nil {
// 			http.Error(w, "Failed to write file", http.StatusInternalServerError)
// 			return
// 		}

// 		// Process log file
// 		results := make(chan LogEntry)
// 		go processFile(filePath, results)

// 		for entry := range results {
// 			mutex.Lock()
// 			_, err := db.Exec("INSERT INTO logs (timestamp, level, message) VALUES ($1, $2, $3)", entry.Timestamp, entry.Level, entry.Message)
// 			mutex.Unlock()
// 			if err != nil {
// 				fmt.Println("Error inserting log entry:", err)
// 			}
// 		}

// 		http.Redirect(w, r, "/", http.StatusSeeOther)
// 	} else {
// 		tmpl := template.Must(template.ParseFiles("frontend/index.html"))
// 		tmpl.Execute(w, nil)
// 	}
// }

// func processFile(filepath string, results chan LogEntry) {
// 	file, err := os.Open(filepath)
// 	if err != nil {
// 		fmt.Println("Error opening file:", err)
// 		close(results)
// 		return
// 	}
// 	defer file.Close()

// 	scanner := bufio.NewScanner(file)
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		results <- LogEntry{
// 			Timestamp: "2024-06-28T00:00:00Z", // Mock timestamp
// 			Level:     "INFO",
// 			Message:   line,
// 		}
// 	}
// 	close(results)
// }

// func tableHandler(w http.ResponseWriter, r *http.Request) {
// 	rows, err := db.Query("SELECT timestamp, level, message FROM logs")
// 	if err != nil {
// 		http.Error(w, "Failed to retrieve logs", http.StatusInternalServerError)
// 		return
// 	}
// 	defer rows.Close()

// 	var logs []LogEntry
// 	for rows.Next() {
// 		var log LogEntry
// 		if err := rows.Scan(&log.Timestamp, &log.Level, &log.Message); err != nil {
// 			http.Error(w, "Error reading logs", http.StatusInternalServerError)
// 			return
// 		}
// 		logs = append(logs, log)
// 	}

// 	json.NewEncoder(w).Encode(logs)
// }

// func main() {
// 	initDB()

// 	http.HandleFunc("/", uploadHandler)
// 	http.HandleFunc("/logs", tableHandler)

// 	fmt.Println("Server started on :9090")
// 	http.ListenAndServe(":9090", nil)
// }

package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"sync"

	_ "github.com/lib/pq" // PostgreSQL driver
)

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

var db *sql.DB
var mutex sync.Mutex

// Initialize database and create table if it doesn't exist
func initDB() (*sql.DB, error) {
	dsn := "host=postgres user=log_user password=log_password dbname=log_db port=5432 sslmode=disable" // Update to port 5432
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
		return nil, err
	}

	log.Println("Database connected successfully")
	return db, nil
}

// Handle file upload and save to disk
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Failed to upload file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Save uploaded file
		filePath := "./" + handler.Filename
		out, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			http.Error(w, "Failed to write file", http.StatusInternalServerError)
			return
		}

		// Process log file in parallel
		results := make(chan LogEntry, 100)
		var wg sync.WaitGroup

		// Get the number of available CPU cores
		numCores := runtime.NumCPU()
		wg.Add(numCores)

		// Split the work across multiple goroutines
		go processFileInChunks(filePath, results, &wg, numCores)

		// Wait for goroutines to finish
		go func() {
			wg.Wait()
			close(results)
		}()

		// Insert logs into the database
		for entry := range results {
			mutex.Lock()
			_, err := db.Exec("INSERT INTO logs (timestamp, level, message) VALUES ($1, $2, $3)", entry.Timestamp, entry.Level, entry.Message)
			mutex.Unlock()
			if err != nil {
				fmt.Println("Error inserting log entry:", err)
			}
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		tmpl := template.Must(template.ParseFiles("frontend/index.html"))
		tmpl.Execute(w, nil)
	}
}

// Process file in chunks with concurrency
func processFileInChunks(filepath string, results chan LogEntry, wg *sync.WaitGroup, numChunks int) {
	// Open the file
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		close(results)
		wg.Done()
		return
	}
	defer file.Close()

	// Determine the size of each chunk
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		close(results)
		wg.Done()
		return
	}
	chunkSize := int(fileInfo.Size()) / numChunks

	// Create a buffered reader for efficient file reading
	reader := bufio.NewReader(file)

	// Read file in chunks
	for i := 0; i < numChunks; i++ {
		go func(i int) {
			defer wg.Done()
			// Each goroutine processes a chunk of the file
			startByte := i * chunkSize
			endByte := startByte + chunkSize
			if i == numChunks-1 {
				endByte = int(fileInfo.Size())
			}

			// Seek to the chunk's starting byte
			_, err := file.Seek(int64(startByte), io.SeekStart)
			if err != nil {
				fmt.Println("Error seeking file:", err)
				return
			}

			// Read lines within the chunk
			scanner := bufio.NewScanner(io.LimitReader(reader, int64(endByte-startByte)))
			for scanner.Scan() {
				line := scanner.Text()
				// Mock timestamp and level for the sake of this example
				results <- LogEntry{
					Timestamp: "2024-06-28T00:00:00Z", // Mock timestamp
					Level:     "INFO",
					Message:   line,
				}
			}
		}(i)
	}
}

// Fetch logs from database
func tableHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT timestamp, level, message FROM logs")
	if err != nil {
		http.Error(w, "Failed to retrieve logs", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var logs []LogEntry
	for rows.Next() {
		var log LogEntry
		if err := rows.Scan(&log.Timestamp, &log.Level, &log.Message); err != nil {
			http.Error(w, "Error reading logs", http.StatusInternalServerError)
			return
		}
		logs = append(logs, log)
	}

	// Check if logs are being fetched correctly
	log.Println("Fetched logs:", logs)

	// Return logs as JSON
	json.NewEncoder(w).Encode(logs)
}

func main() {
	var err error
	db, err = initDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	http.HandleFunc("/", uploadHandler)
	http.HandleFunc("/logs", tableHandler)

	fmt.Println("Server started on :9090")
	http.ListenAndServe(":9090", nil)
}
