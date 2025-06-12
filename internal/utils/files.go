package utils

import (
	"os"
)

// CreateTempFile creates a temporary file for the given directory, name and contents
func CreateTempFile(dir string, namePattern string, fileContents string) (*os.File, error) {
	// Write script to temp file
	tmpFile, err := os.CreateTemp(dir, namePattern)
	if err != nil {

		return nil, err
	}

	if _, err := tmpFile.WriteString(fileContents); err != nil {
		return nil, err
	}
	tmpFile.Close()

	// Make file executable
	if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
		return nil, err
	}
	return tmpFile, nil
}
