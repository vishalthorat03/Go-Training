package handlers

import (
	"csv_manager/utils"
	"fmt"
	"io"
	"net/http"
	"os"
)

// ServeHomepage serves the homepage
func ServeHomepage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		html := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>CSV Manager</title>
			<script>
				// Generic function to perform an operation
				function performOperation(endpoint) {
					fetch(endpoint)
						.then(response => response.text())
						.then(data => {
							document.getElementById("result").innerHTML = data;
						})
						.catch(error => {
							alert("An error occurred: " + error);
						});
				}
			</script>
		</head>
		<body>
			<h1>Welcome to CSV Manager</h1>
			<p>Upload a CSV file to begin.</p>
			<form action="/upload" method="post" enctype="multipart/form-data">
				<label for="file">Choose a CSV file:</label>
				<input type="file" name="file" id="file" accept=".csv" required>
				<button type="submit">Upload</button>
			</form>
			<div id="operations" style="display:none;">
				<h2>Choose an Operation:</h2>
				<button onclick="performOperation('/list')">Show List</button>
				<button onclick="performOperation('/query')">Query Data</button>
				<button onclick="performOperation('/add')">Add Entry</button>
				<button onclick="performOperation('/delete')">Delete Entry</button>
			</div>
			<div id="result"></div>
		</body>
		</html>
		`
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, html)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleFileUpload handles CSV file uploads
func HandleFileUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Failed to read uploaded file: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Save the file
		outputFile, err := os.Create("data/uploaded.csv")
		if err != nil {
			http.Error(w, "Failed to save file: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer outputFile.Close()

		if _, err := io.Copy(outputFile, file); err != nil {
			http.Error(w, "Failed to save file: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Load the CSV into memory
		if err := utils.LoadCSV("data/uploaded.csv"); err != nil {
			http.Error(w, "Failed to load CSV: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `
		<!DOCTYPE html>
		<html>
		<head><title>Upload Successful</title></head>
		<body>
			<h1>File Uploaded Successfully!</h1>
			<div id="operations">
				<h2>Choose an Operation:</h2>
				<button onclick="performOperation('/list')">Show List</button>
				<button onclick="performOperation('/query')">Query Data</button>
				<button onclick="performOperation('/add')">Add Entry</button>
				<button onclick="performOperation('/delete')">Delete Entry</button>
			</div>
			<div id="result"></div>
			<script>
				function performOperation(endpoint) {
					fetch(endpoint)
						.then(response => response.text())
						.then(data => {
							document.getElementById("result").innerHTML = data;
						})
						.catch(error => {
							alert("An error occurred: " + error);
						});
				}
			</script>
		</body>
		</html>
		`)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
