package handlers

import (
	"fmt"
	"net/http"
)

var (
	TotalRows     int64
	ProcessedRows int64
)

func UploadProgress(w http.ResponseWriter, r *http.Request) {
	progress := float64(ProcessedRows) / float64(TotalRows) * 100
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"progress": %.2f}`, progress)))
}
