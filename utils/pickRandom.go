package utils

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func PickRandom(csvName, filePath string, count int) error {
	records, err := readCSVFile(csvName)
	if err != nil {
		return fmt.Errorf("could not pick random values: %v", err)
	}

	randomEmails := getRandomValues(records, count)

	if err := writeCSVFile(filePath, randomEmails); err != nil {
		return fmt.Errorf("could not write random values to os: %v", err)
	}

	return nil
}

func readCSVFile(csvName string) ([][]string, error) {
	file, err := os.Open(csvName)
	if err != nil {
		return nil, fmt.Errorf("could not open csv file: %v", err)
	}

	defer file.Close()

	// Read the CSV data
	reader := csv.NewReader(file)

	// Allow variable number of fields
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("could not read csv file: %v", err)
	}

	return records, nil
}

func getRandomValues(data [][]string, count int) []string {
	rand.NewSource(time.Now().UnixNano())

	randomValues := make([]string, 0, count)

	for i := 0; i < count; i++ {
		randomIndex := rand.Intn(len(data))
		randomValues = append(randomValues, data[randomIndex][0]) // Assuming the first column contains the values you want to pick
	}

	return randomValues
}

func writeCSVFile(filePath string, data []string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range data {
		err := writer.Write([]string{value})
		if err != nil {
			return err
		}
	}

	return nil
}
