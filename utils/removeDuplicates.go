package utils

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func RemoveDuplicates(csvName string) error {
	file, err := os.Open(csvName)
	if err != nil {
		return fmt.Errorf("error opening csv: %v", err)
	}
	defer file.Close()

	// Read the CSV data
	reader := csv.NewReader(file)

	// Allow variable number of fields
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("error reading csv: %v", err)
	}

	// Create a map to store unique email addresses.
	uniqueEmails := make(map[string]bool)

	// Iterate through the CSV data and remove duplicates.
	var uniqueRecords [][]string
	for _, record := range records {
		if len(record) > 0 {
			email := record[0] // Assuming the email is in the first column
			if !uniqueEmails[email] {
				uniqueEmails[email] = true
				uniqueRecords = append(uniqueRecords, record)
			}
		}
	}

	// Close the file before truncating.
	file.Close()

	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	filePath := currentDir + "/" + file.Name()

	// Open the file in write mode, truncating the content.
	file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error reopening the file: %v", err)
	}
	defer file.Close()

	// Write the unique email addresses back to the input CSV file.
	writer := csv.NewWriter(file)
	for _, record := range uniqueRecords {
		err := writer.Write(record)
		if err != nil {
			return fmt.Errorf("error writing to the input file: %v", err)
		}
	}
	writer.Flush()

	if err := writer.Error(); err != nil {
		return fmt.Errorf("error flushing the writer: %v", err)
	}

	return nil
}
