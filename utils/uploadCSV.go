package utils

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func UploadCsvToSpreadSheet(csvName string) error {
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
			Title: "EMAIL LEADS",
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

	count := 0

	for _, row := range records {
		var interfaceRow []interface{}

		for _, cell := range row {
			count += 1
			interfaceRow = append(interfaceRow, cell)
		}

		values = append(values, interfaceRow)
	}

	log.Printf("count: %d", count)

	// Create a new value range for the CSV data.
	valueRange := &sheets.ValueRange{
		Values: values,
	}

	// Update the Google Spreadsheet.
	_, err = srv.Spreadsheets.Values.Update(spreadsheet.SpreadsheetId, "A1", valueRange).ValueInputOption("RAW").Do()
	if err != nil {
		return fmt.Errorf("unable to update spreadsheet: %v", err)
	}

	srv.Spreadsheets.BatchUpdate(spreadsheet.SpreadsheetId, &sheets.BatchUpdateSpreadsheetRequest{
		IncludeSpreadsheetInResponse: true,
		Requests:                     []*sheets.Request{},
		ResponseIncludeGridData:      true,
	})

	return nil
}
