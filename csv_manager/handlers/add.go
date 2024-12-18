package handlers

import (
	"csv_manager/utils"
	"encoding/json"
	"net/http"
)

func AddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var newRow []string
	if err := json.NewDecoder(r.Body).Decode(&newRow); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	data, err := utils.ReadCSV()
	if err != nil {
		http.Error(w, "Error reading CSV file", http.StatusInternalServerError)
		return
	}

	data = append(data, newRow)

	// Save updated data to a new file
	if err := utils.SaveNewVersion(data); err != nil {
		http.Error(w, "Error writing to CSV file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Row added successfully"))
}
