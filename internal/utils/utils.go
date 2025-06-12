package utils

import (
	"fmt"
	"os"
)

type Command string

const (
	FETCH   Command = "fetch"
	CLEAN   Command = "clean"
	INSTALL Command = "install"
	BUILD   Command = "build"
	ALL     Command = "all"
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

// CheckIfError should be used to naively panics if an error is not nil.
// TODO should it accept a channel and send the error
func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}
