package handlers

import (
	"csv_manager/utils"
	"encoding/json"
	"net/http"
	"sort"
)

func SortHandler(w http.ResponseWriter, r *http.Request) {
	order := r.URL.Query().Get("order")
	data, _ := utils.ReadCSV()

	sort.Slice(data, func(i, j int) bool {
		if order == "asc" {
			return data[i][2] < data[j][2] // Assuming level is in column 3
		}
		return data[i][2] > data[j][2]
	})

	json.NewEncoder(w).Encode(data)
}
