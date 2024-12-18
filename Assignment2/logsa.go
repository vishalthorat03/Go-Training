package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
)

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

var logPattern = regexp.MustCompile(`^(?P<timestamp>\S+) \[(?P<level>[Aa-Z]+)\] (?P<message>.+)$`)

func parseLine(line string) (LogEntry, bool) {
	matches := logPattern.FindStringSubmatch(line)
	if matches == nil {
		return LogEntry{}, false
	}
	return LogEntry{
		Timestamp: matches[1],
		Level:     matches[2],
		Message:   matches[3],
	}, true
}

func processChunk(lines []string, results chan LogEntry, seen *sync.Map) {
	for _, line := range lines {
		if logEntry, valid := parseLine(line); valid {
			// Check for duplicates using a concurrent map
			entryKey := fmt.Sprintf("%s|%s|%s", logEntry.Timestamp, logEntry.Level, logEntry.Message)
			if _, exists := seen.LoadOrStore(entryKey, struct{}{}); !exists {
				results <- logEntry
			}
		}
	}
}

func readChunks(filepath string, chunkSize int) ([][]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var chunks [][]string
	scanner := bufio.NewScanner(file)
	var chunk []string
	for scanner.Scan() {
		chunk = append(chunk, scanner.Text())
		if len(chunk) >= chunkSize {
			chunks = append(chunks, chunk)
			chunk = nil
		}
	}
	if len(chunk) > 0 {
		chunks = append(chunks, chunk)
	}
	return chunks, scanner.Err()
}

func main() {
	// Ask for the file name
	fmt.Print("Enter the log file name: ")
	var filepath string
	fmt.Scanln(&filepath)

	outputFile := strings.Replace(filepath, ".log", ".json", 1)
	chunkSize := 1000

	// Read file in chunks
	chunks, err := readChunks(filepath, chunkSize)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	results := make(chan LogEntry, chunkSize)
	var wg sync.WaitGroup
	var seen sync.Map

	// Process chunks concurrently
	for _, chunk := range chunks {
		wg.Add(1)
		go func(chunk []string) {
			defer wg.Done()
			processChunk(chunk, results, &seen)
		}(chunk)
	}

	// Close results channel after all processing is done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	var logEntries []LogEntry
	for entry := range results {
		logEntries = append(logEntries, entry)
	}

	// Save results to JSON
	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(logEntries); err != nil {
		fmt.Printf("Error encoding JSON: %v\n", err)
		return
	}

	// Print a sample of the JSON
	sampleSize := 5
	if len(logEntries) < sampleSize {
		sampleSize = len(logEntries)
	}

	sample := logEntries[:sampleSize]
	prettySample, err := json.MarshalIndent(sample, "", "  ")
	if err != nil {
		fmt.Printf("Error creating JSON sample: %v\n", err)
		return
	}

	fmt.Println("Sample JSON Output:")
	fmt.Println(string(prettySample))

	fmt.Printf("Processed results saved to %s\n", outputFile)
}
