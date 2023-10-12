package utils

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func LimitColumsToOneAndUploadCSV(csvName string) error {
	ctx := context.Background()

	b, err := os.ReadFile("credentials.json")
	if err != nil {
		return fmt.Errorf("unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, sheets.SpreadsheetsScope)
	if err != nil {
		return fmt.Errorf("unable to parse client secret file to config: %v", err)
	}

	client := getClient(config)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("unable to retrieve Sheets client: %v", err)
	}

	// Create a new spreadsheet.
	spreadsheet := &sheets.Spreadsheet{
		Properties: &sheets.SpreadsheetProperties{
			Title: "EMAILS",
		},
	}

	spreadsheet, err = srv.Spreadsheets.Create(spreadsheet).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("could not create spreadsheet: %v", err)
	}

	csvFile, err := os.Open(csvName)
	if err != nil {
		return fmt.Errorf("error opening CSV file: %v", err)
	}

	defer csvFile.Close()

	// Read the CSV data.
	reader := csv.NewReader(csvFile)

	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading CSV file: %v", err)
	}

	// Convert [][]string to [][]interface{}.
	var values [][]interface{}

	for _, row := range records {
		var interfaceRow []interface{}

		for _, cell := range row {
			interfaceRow = append(interfaceRow, cell)
		}

		values = append(values, interfaceRow)
	}

	sheet := spreadsheet.Sheets[0]

	r := &sheets.Request{DeleteDimension: &sheets.DeleteDimensionRequest{Range: &sheets.DimensionRange{
		Dimension:  "COLUMNS",
		SheetId:    sheet.Properties.SheetId,
		StartIndex: 1,
	}}}

	_, err = srv.Spreadsheets.BatchUpdate(spreadsheet.SpreadsheetId, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{r},
	}).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("could not batch update sheet: %v", err)
	}

	// Create a new value range for the CSV data.
	valueRange := &sheets.ValueRange{
		Values: values,
	}

	// Update the Google Spreadsheet.
	_, err = srv.Spreadsheets.Values.Update(spreadsheet.SpreadsheetId, "A1", valueRange).ValueInputOption("RAW").Do()
	if err != nil {
		return fmt.Errorf("unable to update spreadsheet: %v", err)
	}

	return nil
}
