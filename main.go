package main

import (
	"log"

	"github.com/kwandapchumba/merge-mail-csv/utils"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Printf("could not load config: %v", err)
		return
	}

	folderPath := config.FolderPath         // Replace with the source folder path
	outputFileName := config.OutPutFileName // Name of the output merged CSV file

	if err := utils.MergeCSVFiles(folderPath, outputFileName); err != nil {
		log.Printf("could not copy csv to current dir: %v", err)
	} else {
		log.Println("CSV files merged successfully")
	}

	if err := utils.RemoveDuplicates(outputFileName); err != nil {
		log.Printf("could not remove duplicates: %v", err)
	} else {
		log.Println("Duplicates removed from CSV successfully")
	}

	// update merged csv to google sheets without limiting columns
	// if err := utils.UploadCsvToSpreadSheet(outputFileName); err != nil {
	// 	log.Printf("could not upload CSV to spreadsheet: %v", err)
	// } else {
	// 	log.Println("CSV data uploaded to Google Spreadsheet.")
	// }

	// upload merged csv to google sheets and limit columns to 1
	if err := utils.LimitColumsToOneAndUploadCSV(outputFileName); err != nil {
		log.Printf("could not upload CSV to spreadsheet: %v", err)
	} else {
		log.Println("CSV data uploaded to Google Spreadsheet.")
	}
}
