package handlers

import (
	"csv_manager/utils"
	"encoding/json"
	"net/http"
)

func ListHandler(w http.ResponseWriter, r *http.Request) {
	data, err := utils.ReadCSV()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(data)
}
