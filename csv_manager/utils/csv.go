package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

var csvFilePath string

func SetCSVPath(filePath string) {
	csvFilePath = filePath
}

func ReadCSV() ([][]string, error) {
	file, err := os.Open(csvFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	return reader.ReadAll()
}

func SaveNewVersion(data [][]string) error {
	originalFile := csvFilePath
	newFile := fmt.Sprintf("%s_%s.csv", originalFile[:len(originalFile)-4], time.Now().Format("20060102_150405"))

	file, err := os.Create(newFile)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	return writer.WriteAll(data)
}
