package handlers

import (
	"csv_manager/utils"
	"encoding/json"
	"net/http"
)

func QueryHandler(w http.ResponseWriter, r *http.Request) {
	level := r.URL.Query().Get("level")
	if level == "" {
		http.Error(w, "Missing 'level' query parameter", http.StatusBadRequest)
		return
	}

	data, err := utils.ReadCSV()
	if err != nil {
		http.Error(w, "Error reading CSV file", http.StatusInternalServerError)
		return
	}

	var result [][]string
	for _, row := range data {
		if len(row) > 3 && row[3] == level { // Match the Criticality column
			result = append(result, row)
		}
	}

	if len(result) == 0 {
		http.Error(w, "No matching rows found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(result)
}
