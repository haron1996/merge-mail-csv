package utils

import (
	"encoding/csv"
	"io"
	"os"
	"path/filepath"
)

func MergeCSVFiles(inputDir, outputFileName string) error {
	// Create the output CSV file
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		return err
	}

	defer outputFile.Close()

	// Create a CSV writer for the output file
	outputWriter := csv.NewWriter(outputFile)

	defer outputWriter.Flush()

	// Read and merge CSV files from the input directory
	files, err := filepath.Glob(filepath.Join(inputDir, "*.csv"))
	if err != nil {
		return err
	}

	// Initialize a flag to write the CSV header only once
	writeHeader := false

	for _, file := range files {
		// Open the CSV file
		csvFile, err := os.Open(file)
		if err != nil {
			return err
		}

		defer csvFile.Close()

		// Create a CSV reader for the current file
		csvReader := csv.NewReader(csvFile)

		// Read and merge the CSV data
		for {
			record, err := csvReader.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}

			// Write the header to the output file only once
			if !writeHeader {
				err := outputWriter.Write(record)
				if err != nil {
					return err
				}
				writeHeader = true
			}

			// Write the data to the output file
			err = outputWriter.Write(record)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
