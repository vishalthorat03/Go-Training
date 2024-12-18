package utils

import "net/http"

// Operation defines the interface for all CSV operations
type Operation interface {
	Perform(w http.ResponseWriter, r *http.Request) error
}
