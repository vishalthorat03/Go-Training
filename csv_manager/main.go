package main

import (
	"bytes"
	"csv_manager/handlers"
	"csv_manager/utils"
	"embed"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

//go:embed templates/*
var templateFiles embed.FS

type PageData struct {
	Message string
}

func main() {
	// Initialize the logger
	utils.SetupLogger()

	// Handlers for API endpoints
	http.HandleFunc("/api/list", handlers.ListHandler)
	http.HandleFunc("/api/query", handlers.QueryHandler)
	http.HandleFunc("/api/sort", handlers.SortHandler)
	http.HandleFunc("/api/add", handlers.AddHandler)
	http.HandleFunc("/api/delete", handlers.DeleteHandler)
	http.HandleFunc("/upload", uploadHandler)

	// Serve the dynamic frontend
	http.HandleFunc("/", frontendHandler)

	// Start the server
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// frontendHandler serves the HTML page with embedded templates
func frontendHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(templateFiles, "templates/index.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	// Data to be passed to the template
	data := PageData{
		Message: "Welcome to the CSV Manager! Use the options below to interact with your CSV.",
	}

	// Render the template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

	// Set content type and write the rendered HTML
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}

// uploadHandler handles CSV file upload
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form (max memory for file)
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get the file from the form
	file, _, err := r.FormFile("csvfile")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Save the uploaded file
	out, err := os.Create("uploaded.csv")
	if err != nil {
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Set the CSV file path
	utils.SetCSVPath("uploaded.csv")

	// Redirect back to home page after upload
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
