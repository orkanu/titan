package utils

import (
	"fmt"
	"os"
	"strings"
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

func getPathWithUserHome(dir string) string {
	if strings.HasPrefix(dir, "~") {
		home, _ := os.UserHomeDir()
		return home + strings.TrimPrefix(dir, "~")
	}
	return dir
}

// CheckIsFile check if the given path is a file
func CheckIsFile(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}
