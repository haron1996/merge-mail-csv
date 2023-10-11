package utils

import (
	"io"
	"os"
)

func GetCVSinFolder(folderPath string) error {
	// Open the source folder
	sourceFolder, err := os.Open(folderPath)
	if err != nil {
		return err
	}
	defer sourceFolder.Close()

	// Read the source folder contents
	files, err := sourceFolder.Readdir(-1)
	if err != nil {
		return err
	}

	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Loop through the files
	for _, file := range files {
		if file.Mode().IsRegular() {
			destFilePath := currentDir + "/" + file.Name()

			// Delete the file in the destination directory if it exists
			if _, err := os.Stat(destFilePath); err == nil {
				if err := os.Remove(destFilePath); err != nil {
					return err
				}
			}
			copyEachCSVToCurrentDir(file, folderPath, currentDir)
		}
	}
	return nil
}

func copyEachCSVToCurrentDir(file os.FileInfo, folderPath, currentDir string) error {

	sourceFile, err := os.Open(folderPath + "/" + file.Name())
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFilePath := currentDir + "/" + file.Name()

	destFile, err := os.Create(destFilePath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy the file contents from source to destination
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}
