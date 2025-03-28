package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
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

func ShowTable(w http.ResponseWriter, r *http.Request) {
	// Connect to the database
	db, err := sql.Open("postgres", "postgres://postgres:password@db:5432/csvdb?sslmode=disable")
	if err != nil {
		log.Println("Failed to connect to the database:", err)
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
	sortOrder := r.URL.Query().Get("sortOrder")

	if page == "" {
		page = "1"
	}

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}

	rowsPerPage := 1000
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

	// Handling sorting based on query parameters
	if sortColumn == "purchasedate" || sortColumn == "warrantyend" || sortColumn == "price" {
		if sortOrder == "asc" || sortOrder == "desc" {
			query += fmt.Sprintf(" ORDER BY %s %s", sortColumn, sortOrder)
		}
	}
	// Apply LIMIT and OFFSET for pagination
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", rowsPerPage, offset)

	// Fetch the devices for the current page
	rows, err := db.Query(query)
	if err != nil {
		log.Println("Failed to fetch data:", err)
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var devices []Device
	for rows.Next() {
		var device Device
		err := rows.Scan(&device.ID, &device.DeviceName, &device.DeviceType, &device.Brand, &device.Model, &device.OS, &device.OSVersion, &device.PurchaseDate, &device.WarrantyEnd, &device.Status, &device.Price)
		if err != nil {
			log.Println("Failed to read data:", err)
			http.Error(w, "Failed to read data", http.StatusInternalServerError)
			return
		}
		devices = append(devices, device)
	}

	// Query for the total row count using the reltuples estimate from pg_class
	countQuery := "SELECT reltuples::BIGINT AS estimate FROM pg_class WHERE relname = 'devices'"

	var totalRows int64
	err = db.QueryRow(countQuery).Scan(&totalRows)
	if err != nil {
		log.Println("Failed to fetch the filtered row count:", err)
		http.Error(w, "Failed to fetch the filtered row count", http.StatusInternalServerError)
		return
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalRows) / float64(rowsPerPage)))

	// Render the HTML with pagination
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

	// Display the devices
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

	// Pagination controls
	fmt.Fprintf(w, `
			</table>
			<p>%d devices found</p>
			<div class="pagination">
	`, totalRows)

	// Pagination links
	if pageNum > 1 {
		fmt.Fprintf(w, `<a href="?deviceName=%s&deviceType=%s&brand=%s&os=%s&status=%s&page=1" class="button">First</a>`, deviceName, deviceType, brand, os, status)
		fmt.Fprintf(w, `<a href="?deviceName=%s&deviceType=%s&brand=%s&os=%s&status=%s&page=%d" class="button">Previous</a>`, deviceName, deviceType, brand, os, status, pageNum-1)
	}
	if pageNum < totalPages {
		fmt.Fprintf(w, `<a href="?deviceName=%s&deviceType=%s&brand=%s&os=%s&status=%s&page=%d" class="button">Next</a>`, deviceName, deviceType, brand, os, status, pageNum+1)
		fmt.Fprintf(w, `<a href="?deviceName=%s&deviceType=%s&brand=%s&os=%s&status=%s&page=%d" class="button">Last</a>`, deviceName, deviceType, brand, os, status, totalPages)
	}
	fmt.Fprintf(w, `</div></div></body></html>`)
}
