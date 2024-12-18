package handlers

import (
	"csv_to_db/models"
	"encoding/csv"
	"html/template"
	"io"
	"net/http"
	"strconv"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Display the upload form
		tmpl := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Upload CSV</title>
		</head>
		<body>
			<h1>Upload CSV File</h1>
			<form enctype="multipart/form-data" action="/" method="post">
				<input type="file" name="file" accept=".csv" required>
				<button type="submit">Upload</button>
			</form>
		</body>
		</html>`
		w.Write([]byte(tmpl))
		return
	}

	if r.Method == http.MethodPost {
		// Handle file upload
		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Failed to read file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Parse the CSV file and store it in the database
		reader := csv.NewReader(file)
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				http.Error(w, "Failed to parse CSV file", http.StatusInternalServerError)
				return
			}

			// Skip header row if present
			if record[0] == "iteID" {
				continue
			}

			// Parse fields
			fxiletID := record[1]
			name := record[2]
			criticality := record[3]
			relevantComputerCount, _ := strconv.Atoi(record[4])

			// Insert into the database
			update := models.Update{
				FxiletID:              fxiletID,
				Name:                  name,
				Criticality:           criticality,
				RelevantComputerCount: relevantComputerCount,
			}
			models.DB.Create(&update)
		}

		// Redirect to the view page
		http.Redirect(w, r, "/view", http.StatusSeeOther)
	}
}

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	criticality := r.URL.Query().Get("criticality")
	pageParam := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 1 {
		page = 1
	}

	// Pagination logic
	pageSize := 100
	offset := (page - 1) * pageSize

	// Query the database
	var updates []models.Update
	query := models.DB.Offset(offset).Limit(pageSize)
	if criticality != "" {
		query = query.Where("criticality = ?", criticality)
	}
	result := query.Find(&updates)
	if result.Error != nil {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

	// Render the data as an HTML table
	tmpl := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>View Updates</title>
	</head>
	<body>
		<h1>Updates Table</h1>
		<div>
			<a href="/view?criticality=Low&page=1"><button>Low</button></a>
			<a href="/view?criticality=Moderate&page=1"><button>Moderate</button></a>
			<a href="/view?criticality=Critical&page=1"><button>Critical</button></a>
		</div>
		<table border="1">
			<tr>
				<th>ID</th>
				<th>FxiletID</th>
				<th>Name</th>
				<th>Criticality</th>
				<th>Relevant Computer Count</th>
			</tr>
			{{range .Updates}}
			<tr>
				<td>{{.ID}}</td>
				<td>{{.FxiletID}}</td>
				<td>{{.Name}}</td>
				<td>{{.Criticality}}</td>
				<td>{{.RelevantComputerCount}}</td>
			</tr>
			{{end}}
		</table>
		<div>
			{{if .PrevPage}}
			<a href="/view?criticality={{.Criticality}}&page={{.PrevPage}}"><button>Previous</button></a>
			{{end}}
			{{if .NextPage}}
			<a href="/view?criticality={{.Criticality}}&page={{.NextPage}}"><button>Next</button></a>
			{{end}}
		</div>
	</body>
	</html>`

	// Prepare template data
	var nextPage, prevPage int
	if len(updates) == pageSize {
		nextPage = page + 1
	}
	if page > 1 {
		prevPage = page - 1
	}

	data := struct {
		Updates     []models.Update
		Criticality string
		NextPage    int
		PrevPage    int
	}{
		Updates:     updates,
		Criticality: criticality,
		NextPage:    nextPage,
		PrevPage:    prevPage,
	}

	// Parse and execute template
	t := template.Must(template.New("view").Parse(tmpl))
	t.Execute(w, data)
}
