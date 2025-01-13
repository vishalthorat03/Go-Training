package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

var (
	totalRows, processedRows int64 // Track total and processed rows
)

type Device struct {
	ID           int
	DeviceName   string
	DeviceType   string
	Brand        string
	Model        string
	OS           string
	OSVersion    string
	PurchaseDate string
	WarrantyEnd  string
	Status       string
	Price        float64
}

func main() {
	// Serve static files (e.g., index.html, CSS, JS) from the frontend folder
	fs := http.FileServer(http.Dir("frontend"))
	http.Handle("/", fs) // Serve the index.html and other files in the frontend folder

	// Define the other routes
	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/upload-progress", uploadProgress)
	http.HandleFunc("/showTable", showTable)

	// Start the server
	log.Println("Server is running on port 4041...")
	log.Fatal(http.ListenAndServe(":4041", nil))
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Read file from form data
	file, _, err := r.FormFile("csv")
	if err != nil {
		http.Error(w, "Failed to read uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Start measuring time
	start := time.Now()

	// Database connection
	db, err := sql.Open("postgres", "postgres://postgres:postgres@db:5432/devices?sslmode=disable")
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		log.Fatal(err)
	}
	defer db.Close()

	// Create table if it doesn't exist
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS devices (
		id SERIAL PRIMARY KEY,
		devicename TEXT,
		devicetype TEXT,
		brand TEXT,
		model TEXT,
		os TEXT,
		osversion TEXT,
		purchasedate DATE,
		warrantyend DATE,
		status TEXT,
		price FLOAT
	)`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		http.Error(w, "Failed to create table", http.StatusInternalServerError)
		log.Fatal(err)
	}

	// Setup CSV reader
	reader := csv.NewReader(file)
	_, err = reader.Read() // Skip header
	if err != nil {
		http.Error(w, "Failed to read CSV file", http.StatusInternalServerError)
		log.Fatal(err)
	}

	// Channel for rows and batch size
	const batchSize = 5000 // Adjust batch size for optimal performance
	rowChannel := make(chan [][]string, 100)

	var wg sync.WaitGroup

	// Read CSV in chunks
	wg.Add(1)
	go func() {
		defer wg.Done()
		var rows [][]string
		for {
			row, err := reader.Read()
			if err == io.EOF {
				if len(rows) > 0 {
					rowChannel <- rows
				}
				close(rowChannel)
				break
			}
			if err != nil {
				log.Printf("Failed to read row: %v\n", err)
				continue
			}
			rows = append(rows, row)
			if len(rows) >= batchSize {
				rowChannel <- rows
				rows = nil
			}
		}
	}()

	// Concurrent database writers
	const writerCount = 10 // Number of parallel database writers
	for i := 0; i < writerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for rows := range rowChannel {
				insertBatch(db, rows)
			}
		}()
	}

	// Wait for completion
	wg.Wait()

	// Measure elapsed time
	elapsed := time.Since(start)
	log.Printf("CSV upload completed in %v\n", elapsed)
	w.Write([]byte("Upload completed successfully"))
}

func insertBatch(db *sql.DB, rows [][]string) {
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v\n", err)
		return
	}

	// Build the query for bulk insert
	query := `
		INSERT INTO devices (id, devicename, devicetype, brand, model, os, osversion, purchasedate, warrantyend, status, price)
		VALUES `
	valueStrings := make([]string, 0, len(rows))
	valueArgs := make([]interface{}, 0, len(rows)*11)

	// Construct query dynamically
	for i, row := range rows {
		id, err := strconv.Atoi(row[0])
		if err != nil {
			log.Printf("Invalid id format for row %+v: %v\n", row, err)
			continue
		}

		price, err := strconv.ParseFloat(row[10], 64)
		if err != nil {
			log.Printf("Invalid price format for row %+v: %v\n", row, err)
			continue
		}

		// Append the placeholder for this row
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			i*11+1, i*11+2, i*11+3, i*11+4, i*11+5, i*11+6, i*11+7, i*11+8, i*11+9, i*11+10, i*11+11))

		// Append values for the placeholders
		valueArgs = append(valueArgs, id, row[1], row[2], row[3], row[4], row[5], row[6], row[7], row[8], row[9], price)
	}

	// Combine query and placeholders
	query += strings.Join(valueStrings, ",")
	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Printf("Failed to prepare bulk insert statement: %v\n", err)
		tx.Rollback()
		return
	}
	defer stmt.Close()

	// Execute the bulk insert
	_, err = stmt.Exec(valueArgs...)
	if err != nil {
		log.Printf("Error executing bulk insert: %v\n", err)
		tx.Rollback()
		return
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Printf("Failed to commit transaction: %v\n", err)
	}
}

// Function to monitor upload progress (can be enhanced later with a better tracking mechanism)
func uploadProgress(w http.ResponseWriter, r *http.Request) {
	progress := float64(processedRows) / float64(totalRows) * 100
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"progress": %.2f}`, progress)))
}

func showTable(w http.ResponseWriter, r *http.Request) {
	// Connect to the database
	db, err := sql.Open("postgres", "postgres://postgres:postgres@db:5432/devices?sslmode=disable")
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Get filters, pagination, and sorting parameters from query string
	deviceType := r.URL.Query().Get("deviceType")
	deviceName := r.URL.Query().Get("deviceName")
	brand := r.URL.Query().Get("brand")
	os := r.URL.Query().Get("os")
	status := r.URL.Query().Get("status")
	page := r.URL.Query().Get("page")
	sortColumn := r.URL.Query().Get("sortColumn")
	sortOrder := r.URL.Query().Get("sortOrder") // "asc" or "desc"

	if page == "" {
		page = "1"
	}

	pageNum, _ := strconv.Atoi(page)
	if pageNum <= 0 {
		pageNum = 1
	}

	rowsPerPage := 100
	offset := (pageNum - 1) * rowsPerPage

	// Build SQL query with filters and sorting
	query := "SELECT id, devicename, devicetype, brand, model, os, osversion, purchasedate, warrantyend, status, price FROM devices WHERE 1=1"
	if deviceType != "" {
		query += fmt.Sprintf(" AND devicetype = '%s'", deviceType)
	}
	if deviceName != "" {
		query += fmt.Sprintf(" AND devicename = '%s'", deviceName)
	}
	if brand != "" {
		query += fmt.Sprintf(" AND brand = '%s'", brand)
	}
	if os != "" {
		query += fmt.Sprintf(" AND os = '%s'", os)
	}
	if status != "" {
		query += fmt.Sprintf(" AND status = '%s'", status)
	}
	if sortColumn == "purchasedate" || sortColumn == "warrantyend" || sortColumn == "price" {
		if sortOrder == "asc" || sortOrder == "desc" {
			query += fmt.Sprintf(" ORDER BY %s %s", sortColumn, sortOrder)
		}
	}
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", rowsPerPage, offset)

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Fetch data into a slice
	var devices []Device
	for rows.Next() {
		var device Device
		err := rows.Scan(&device.ID, &device.DeviceName, &device.DeviceType, &device.Brand, &device.Model, &device.OS, &device.OSVersion, &device.PurchaseDate, &device.WarrantyEnd, &device.Status, &device.Price)
		if err != nil {
			http.Error(w, "Failed to read data", http.StatusInternalServerError)
			return
		}
		devices = append(devices, device)
	}

	// Get total row count for pagination
	countQuery := "SELECT COUNT(*) FROM devices WHERE 1=1"
	if deviceType != "" {
		countQuery += fmt.Sprintf(" AND devicetype = '%s'", deviceType)
	}
	if deviceName != "" {
		countQuery += fmt.Sprintf(" AND devicename = '%s'", deviceName)
	}
	if brand != "" {
		countQuery += fmt.Sprintf(" AND brand = '%s'", brand)
	}
	if os != "" {
		countQuery += fmt.Sprintf(" AND os = '%s'", os)
	}
	if status != "" {
		countQuery += fmt.Sprintf(" AND status = '%s'", status)
	}

	var totalRows int
	err = db.QueryRow(countQuery).Scan(&totalRows)
	if err != nil {
		http.Error(w, "Failed to get total row count", http.StatusInternalServerError)
		return
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(rowsPerPage)))

	// Render HTML with embedded CSS
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Device Table</title>
		<style>
			body {
				font-family: Arial, sans-serif;
				background-color: #f4f4f9;
				margin: 0;
				padding: 0;
			}
			.container {
				width: 90%%;
				margin: 20px auto;
				background-color: #fff;
				padding: 20px;
				border-radius: 8px;
				box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
				text-align: center;
			}
			h1 {
				text-align: center;
			}
			table {
				width: 100%%;
				border-collapse: collapse;
				margin-top: 20px;
			}
			th, td {
				padding: 10px;
				text-align: center;
				border-bottom: 1px solid #ddd;
			}
			th {
				background-color: #f2f2f2;
			}
			th a {
				text-decoration: none;
				color: #333;
			}
			th a:hover {
				text-decoration: underline;
			}
			.button {
				display: inline-block;
				padding: 10px 15px;
				background-color: #4CAF50;
				color: #fff;
				border: none;
				border-radius: 5px;
				text-decoration: none;
				cursor: pointer;
				margin: 5px;
			}
			.button:hover {
				background-color: #45a049;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h1>Device Table</h1>
			<table>
				<tr>
					<th>ID</th>
					<th>Device Name</th>
					<th>Type</th>
					<th>Brand</th>
					<th>Model</th>
					<th>OS</th>
					<th>OS Version</th>
					<th><a href="?sortColumn=purchasedate&sortOrder=asc">Purchase Date ▲</a> | <a href="?sortColumn=purchasedate&sortOrder=desc">▼</a></th>
					<th><a href="?sortColumn=warrantyend&sortOrder=asc">Warranty End ▲</a> | <a href="?sortColumn=warrantyend&sortOrder=desc">▼</a></th>
					<th>Status</th>
					<th><a href="?sortColumn=price&sortOrder=asc">Price ▲</a> | <a href="?sortColumn=price&sortOrder=desc">▼</a></th>
				</tr>
	`)

	for _, device := range devices {
		fmt.Fprintf(w, `
			<tr>
				<td>%d</td>
				<td>%s</td>
				<td>%s</td>
				<td>%s</td>
				<td>%s</td>
				<td>%s</td>
				<td>%s</td>
				<td>%s</td>
				<td>%s</td>
				<td>%s</td>
				<td>%.2f</td>
			</tr>
		`, device.ID, device.DeviceName, device.DeviceType, device.Brand, device.Model, device.OS, device.OSVersion, device.PurchaseDate, device.WarrantyEnd, device.Status, device.Price)
	}

	// Pagination controls and record count
	fmt.Fprintf(w, `
			</table>
			<p>%d records found</p>
			<div>
	`, totalRows)

	if pageNum > 1 {
		fmt.Fprintf(w, `<a href="?page=1" class="button">First</a> `)
		fmt.Fprintf(w, `<a href="?page=%d" class="button">Prev</a> `, pageNum-1)
	}
	if pageNum < totalPages {
		fmt.Fprintf(w, `<a href="?page=%d" class="button">Next</a> `, pageNum+1)
		fmt.Fprintf(w, `<a href="?page=%d" class="button">Last</a> `, totalPages)
	}

	fmt.Fprintf(w, `
			</div>
		</div>
	</body>
	</html>
	`)
}
