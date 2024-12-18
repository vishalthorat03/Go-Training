package utils

import (
	"encoding/csv"
	"errors"
	"os"
	"sync"
)

var (
	entries []map[string]string
	headers []string
	mutex   sync.Mutex
)

// LoadCSV loads the CSV file
func LoadCSV(filename string) error {
	mutex.Lock()
	defer mutex.Unlock()

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return err
	}

	headers = rows[0]
	entries = nil
	for _, row := range rows[1:] {
		entry := make(map[string]string)
		for i, value := range row {
			entry[headers[i]] = value
		}
		entries = append(entries, entry)
	}
	return nil
}

// SaveEntries saves the entries to a CSV file
func SaveEntries(filename string) error {
	mutex.Lock()
	defer mutex.Unlock()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(headers); err != nil {
		return err
	}
	for _, entry := range entries {
		row := make([]string, len(headers))
		for i, header := range headers {
			row[i] = entry[header]
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	return nil
}

// GetEntries retrieves all entries
func GetEntries() ([]map[string]string, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if len(entries) == 0 {
		return nil, errors.New("no entries found")
	}
	return entries, nil
}

// AddEntry adds a new entry
func AddEntry(entry map[string]string) error {
	mutex.Lock()
	defer mutex.Unlock()

	entries = append(entries, entry)
	return nil
}

// DeleteEntry deletes an entry based on a key-value pair
func DeleteEntry(key, value string) error {
	mutex.Lock()
	defer mutex.Unlock()

	for i, entry := range entries {
		if entry[key] == value {
			entries = append(entries[:i], entries[i+1:]...)
			return nil
		}
	}
	return errors.New("entry not found")
}

// QueryEntries filters entries based on a column and value
func QueryEntries(column, value string) ([]map[string]string, error) {
	mutex.Lock()
	defer mutex.Unlock()

	var result []map[string]string
	for _, entry := range entries {
		if entry[column] == value {
			result = append(result, entry)
		}
	}
	if len(result) == 0 {
		return nil, errors.New("no matching entries found")
	}
	return result, nil
}

// SortEntries sorts entries by a column
func SortEntries(column string) error {
	mutex.Lock()
	defer mutex.Unlock()

	for _, header := range headers {
		if header == column {
			// Perform sorting using goroutines
			// Placeholder for implementation
			return nil
		}
	}
	return errors.New("column not found")
}
