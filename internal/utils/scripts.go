package utils

import (
	"os"
	"os/exec"
)

// ExecScript is a utility function that creates a shell script and executes it
func ExecScript(script string, env []string) error {
	// Write script to temp file
	tmpFile, err := CreateTempFile("", "build-action-*.sh", script)
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	// Execute the script
	cmd := exec.Command("bash", tmpFile.Name())
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
