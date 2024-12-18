package handlers

import (
	"csv_manager/utils"
	"net/http"
)

// SaveHandler saves entries to a CSV
func SaveHandler(w http.ResponseWriter, r *http.Request) {
	if err := utils.SaveEntries("data/uploaded.csv"); err != nil {
		http.Error(w, "Failed to save entries: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
