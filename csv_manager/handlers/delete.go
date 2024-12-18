package handlers

import (
	"csv_manager/utils"
	"net/http"
	"strconv"
)

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	indexParam := r.URL.Query().Get("index")
	index, err := strconv.Atoi(indexParam)
	if err != nil || index < 0 {
		http.Error(w, "Invalid index", http.StatusBadRequest)
		return
	}

	data, err := utils.ReadCSV()
	if err != nil {
		http.Error(w, "Error reading CSV file", http.StatusInternalServerError)
		return
	}

	if index >= len(data) {
		http.Error(w, "Index out of range", http.StatusBadRequest)
		return
	}

	data = append(data[:index], data[index+1:]...)

	// Save updated data to a new file
	if err := utils.SaveNewVersion(data); err != nil {
		http.Error(w, "Error writing to CSV file", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Row deleted successfully"))
}
