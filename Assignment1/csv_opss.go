package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Entry represents a record in the CSV file.
type Entry struct {
	SiteID                int
	FixletID              int
	Name                  string
	Criticality           string
	RelevantComputerCount int
}

// ReadCSV reads the CSV file and returns a slice of Entry.
func ReadCSV(filename string) ([]Entry, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var entries []Entry

	// Read and skip the header
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	// Read each record
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		siteID, _ := strconv.Atoi(record[0])
		fixletID, _ := strconv.Atoi(record[1])
		computers, _ := strconv.Atoi(record[4])

		entry := Entry{
			SiteID:                siteID,
			FixletID:              fixletID,
			Name:                  record[2],
			Criticality:           record[3],
			RelevantComputerCount: computers,
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

// WriteCSV writes the entries back to the CSV file.
func WriteCSV(filename string, entries []Entry) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header
	writer.Write([]string{"SiteID", "FixletID", "Name", "Criticality", "RelevantComputerCount"})

	// Write the records
	for _, entry := range entries {
		writer.Write([]string{
			strconv.Itoa(entry.SiteID),
			strconv.Itoa(entry.FixletID),
			entry.Name,
			entry.Criticality,
			strconv.Itoa(entry.RelevantComputerCount),
		})
	}
	return nil
}

// ListEntries prints all entries.
func ListEntries(entries []Entry) {
	fmt.Println("\nEntries:")
	for _, entry := range entries {
		fmt.Printf("SiteID: %d, FixletID: %d, Name: %s, Criticality: %s, Computers: %d\n",
			entry.SiteID, entry.FixletID, entry.Name, entry.Criticality, entry.RelevantComputerCount)
	}
}

// QueryEntry searches for entries by a keyword.
func QueryEntry(entries []Entry, keyword string) {
	fmt.Println("\nSearch Results:")
	keyword = strings.ToLower(keyword)
	found := false
	for _, entry := range entries {
		if strings.Contains(strings.ToLower(entry.Name), keyword) ||
			strings.Contains(strings.ToLower(entry.Criticality), keyword) {
			fmt.Printf("SiteID: %d, FixletID: %d, Name: %s, Criticality: %s, Computers: %d\n",
				entry.SiteID, entry.FixletID, entry.Name, entry.Criticality, entry.RelevantComputerCount)
			found = true
		}
	}
	if !found {
		fmt.Println("No matching entries found.")
	}
}

// AddEntry adds a new entry to the list.
func AddEntry(entries []Entry, entry Entry) []Entry {
	return append(entries, entry)
}

// DeleteEntry removes an entry by FixletID.
func DeleteEntry(entries []Entry, fixletID int) []Entry {
	for i, entry := range entries {
		if entry.FixletID == fixletID {
			fmt.Println("Entry deleted successfully.")
			return append(entries[:i], entries[i+1:]...)
		}
	}
	fmt.Println("No entry found with the given FixletID.")
	return entries
}

// Main menu for user interaction.
func main() {
	var filename string
	fmt.Print("Enter the CSV file name: ")
	fmt.Scanln(&filename)

	entries, err := ReadCSV(filename)
	if err != nil {
		fmt.Println("Error reading CSV file:", err)
		return
	}

	for {
		fmt.Println("\nOptions\nquery\nadd\ndelete\nexit")
		fmt.Print("Enter your choice: ")
		var choice string
		fmt.Scanln(&choice)

		switch choice {
		case "query":
			var keyword string
			fmt.Print("Enter name or criticality to search: ")
			fmt.Scanln(&keyword)
			QueryEntry(entries, keyword)
		case "add":
			var siteID, fixletID, computers int
			var name, criticality string
			fmt.Print("Enter SiteID: ")
			fmt.Scanln(&siteID)
			fmt.Print("Enter FixletID: ")
			fmt.Scanln(&fixletID)
			fmt.Print("Enter Name: ")
			fmt.Scanln(&name)
			fmt.Print("Enter Criticality: ")
			fmt.Scanln(&criticality)
			fmt.Print("Enter RelevantComputerCount: ")
			fmt.Scanln(&computers)

			entry := Entry{
				SiteID:                siteID,
				FixletID:              fixletID,
				Name:                  name,
				Criticality:           criticality,
				RelevantComputerCount: computers,
			}
			entries = AddEntry(entries, entry)
			WriteCSV(filename, entries)
			fmt.Println("Entry added successfully.")
		case "delete":
			var fixletID int
			fmt.Print("Enter FixletID to delete: ")
			fmt.Scanln(&fixletID)
			entries = DeleteEntry(entries, fixletID)
			WriteCSV(filename, entries)
			fmt.Printf("\nEntry has been deleted %d", fixletID)
		case "exit":
			fmt.Println("<<Application Closed>>")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}
