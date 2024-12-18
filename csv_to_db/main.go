// package main

// import (
// 	"encoding/csv"
// 	"fmt"
// 	"html/template"
// 	"io"
// 	"log"
// 	"net/http"
// 	"os"
// 	"strconv"

// 	"github.com/joho/godotenv"
// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// )

// // Define the struct corresponding to the CSV structure
// type Update struct {
// 	ID                    uint   `gorm:"primaryKey"`
// 	FxiletID              string `gorm:"type:varchar(255);not null"`
// 	Name                  string `gorm:"type:text;not null"`
// 	Criticality           string `gorm:"type:varchar(50);not null"`
// 	RelevantComputerCount int    `gorm:"type:int;not null"`
// }

// var db *gorm.DB

// func main() {
// 	// Load environment variables from .env file
// 	err := godotenv.Load(".env")
// 	if err != nil {
// 		log.Fatalf("Error loading .env file: %v", err)
// 	}

// 	// Get PostgreSQL connection details
// 	dbUser := os.Getenv("DB_USER")
// 	dbPassword := os.Getenv("DB_PASSWORD")
// 	dbHost := os.Getenv("DB_HOST")
// 	dbName := os.Getenv("DB_NAME")
// 	dbPort := os.Getenv("DB_PORT")

// 	// Build the connection string
// 	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Calcutta",
// 		dbHost, dbUser, dbPassword, dbName, dbPort)

// 	// Open a connection to the database
// 	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		log.Fatalf("Error opening connection to database: %v", err)
// 	}

// 	// Migrate the schema (create table based on struct)
// 	err = db.AutoMigrate(&Update{})
// 	if err != nil {
// 		log.Fatalf("Error migrating schema: %v", err)
// 	}

// 	// Start the HTTP server
// 	http.HandleFunc("/", uploadHandler)
// 	http.HandleFunc("/view", viewHandler)
// 	log.Println("Starting server on http://localhost:8080...")
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

// // uploadHandler serves the file upload form and handles file uploads
// func uploadHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodGet {
// 		// Display the upload form
// 		tmpl := `
// 		<!DOCTYPE html>
// 		<html>
// 		<head>
// 			<title>Upload CSV</title>
// 		</head>
// 		<body>
// 			<h1>Upload CSV File</h1>
// 			<form enctype="multipart/form-data" action="/" method="post">
// 				<input type="file" name="file" accept=".csv" required>
// 				<button type="submit">Upload</button>
// 			</form>
// 		</body>
// 		</html>`
// 		w.Write([]byte(tmpl))
// 		return
// 	}

// 	if r.Method == http.MethodPost {
// 		// Handle file upload
// 		file, _, err := r.FormFile("file")
// 		if err != nil {
// 			http.Error(w, "Failed to read file", http.StatusBadRequest)
// 			return
// 		}
// 		defer file.Close()

// 		// Parse the CSV file and store it in the database
// 		reader := csv.NewReader(file)
// 		for {
// 			record, err := reader.Read()
// 			if err == io.EOF {
// 				break
// 			}
// 			if err != nil {
// 				http.Error(w, "Failed to parse CSV file", http.StatusInternalServerError)
// 				return
// 			}

// 			// Skip header row if present
// 			if record[0] == "iteID" {
// 				continue
// 			}

// 			// Parse fields
// 			fxiletID := record[1]
// 			name := record[2]
// 			criticality := record[3]
// 			relevantComputerCount, _ := strconv.Atoi(record[4])

// 			// Insert into the database
// 			update := Update{
// 				FxiletID:              fxiletID,
// 				Name:                  name,
// 				Criticality:           criticality,
// 				RelevantComputerCount: relevantComputerCount,
// 			}
// 			db.Create(&update)
// 		}

// 		// Redirect to the view page
// 		http.Redirect(w, r, "/view", http.StatusSeeOther)
// 	}
// }

// // viewHandler displays the records from the database with filters and pagination
// func viewHandler(w http.ResponseWriter, r *http.Request) {
// 	// Parse query parameters
// 	criticality := r.URL.Query().Get("criticality")
// 	pageParam := r.URL.Query().Get("page")
// 	page, err := strconv.Atoi(pageParam)
// 	if err != nil || page < 1 {
// 		page = 1
// 	}

// 	// Pagination logic
// 	pageSize := 100
// 	offset := (page - 1) * pageSize

// 	// Query the database
// 	var updates []Update
// 	query := db.Offset(offset).Limit(pageSize)
// 	if criticality != "" {
// 		query = query.Where("criticality = ?", criticality)
// 	}
// 	result := query.Find(&updates)
// 	if result.Error != nil {
// 		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
// 		return
// 	}

// 	// Render the data as an HTML table
// 	tmpl := `
// 	<!DOCTYPE html>
// 	<html>
// 	<head>
// 		<title>View Updates</title>
// 	</head>
// 	<body>
// 		<h1>Updates Table</h1>
// 		<div>
// 			<a href="/view?criticality=Low&page=1"><button>Low</button></a>
// 			<a href="/view?criticality=Moderate&page=1"><button>Moderate</button></a>
// 			<a href="/view?criticality=Critical&page=1"><button>Critical</button></a>
// 		</div>
// 		<table border="1">
// 			<tr>
// 				<th>ID</th>
// 				<th>FxiletID</th>
// 				<th>Name</th>
// 				<th>Criticality</th>
// 				<th>Relevant Computer Count</th>
// 			</tr>
// 			{{range .Updates}}
// 			<tr>
// 				<td>{{.ID}}</td>
// 				<td>{{.FxiletID}}</td>
// 				<td>{{.Name}}</td>
// 				<td>{{.Criticality}}</td>
// 				<td>{{.RelevantComputerCount}}</td>
// 			</tr>
// 			{{end}}
// 		</table>
// 		<div>
// 			{{if .PrevPage}}
// 			<a href="/view?criticality={{.Criticality}}&page={{.PrevPage}}"><button>Previous</button></a>
// 			{{end}}
// 			{{if .NextPage}}
// 			<a href="/view?criticality={{.Criticality}}&page={{.NextPage}}"><button>Next</button></a>
// 			{{end}}
// 		</div>
// 	</body>
// 	</html>`

// 	// Prepare template data
// 	var nextPage, prevPage int
// 	if len(updates) == pageSize {
// 		nextPage = page + 1
// 	}
// 	if page > 1 {
// 		prevPage = page - 1
// 	}

// 	data := struct {
// 		Updates     []Update
// 		Criticality string
// 		NextPage    int
// 		PrevPage    int
// 	}{
// 		Updates:     updates,
// 		Criticality: criticality,
// 		NextPage:    nextPage,
// 		PrevPage:    prevPage,
// 	}

// 	// Parse and execute template
// 	t := template.Must(template.New("view").Parse(tmpl))
// 	t.Execute(w, data)
// }

package main

import (
	"csv_to_db/handlers"
	"csv_to_db/models"
	"csv_to_db/utils"
	"log"
	"net/http"
)

func main() {
	// Initialize environment variables
	utils.LoadEnv()

	// Initialize the database
	models.InitDB()

	// Route handlers
	http.HandleFunc("/", handlers.UploadHandler)
	http.HandleFunc("/view", handlers.ViewHandler)

	// Start the HTTP server
	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

//http://localhost:5050/
//http://localhost:8080/
