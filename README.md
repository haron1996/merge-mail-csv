# merge-mail-csv
A simple script to merge email csv files, remove duplicates, and upload to Google Sheets using sheets API.

# How to use

- Clone project
- [Read this quickstart to get your Google credentials](https://developers.google.com/sheets/api/quickstart/go)
- Download Google credential to your project root
- Create .env file in your your root
- Group email CSV files in one folder in your local machine
- Add 2 values as follows:
    - folderPath=Path to folder (in your local machine) containing email CVS files you want to match
    - outPutFileName=merged.csv
- Init go mod and run go mod tidy to install required packages
- Run program... Sip your coffee... And now you have a merge list... without duplicates... in your Google Sheet 😊
- If you need help: Don't hesitae to email me at: haronkibetrutoh@gmail.com