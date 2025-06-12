package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"titan/internal/config"
)

// captureEnvironment sets up NVM and pnpm and captures the resulting environment
func CaptureEnvironment(versions config.Versions) ([]string, error) {
	// Source nvm and install the desired nvm and pnpm versions
	scriptContents := fmt.Sprintf(`
		source ~/.nvm/nvm.sh &&
		nvm install %v &&
		nvm use %v &&
		npm i -g pnpm@%v &&
		env
		`, versions.Node, versions.Node, versions.PNPM)
	setupCmd := exec.Command("bash", "-c", scriptContents)
	output, err := setupCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to setup environment: %w", err)
	}
	// Parese the environment output
	envLines := strings.Split(string(output), "\n")
	var env []string
	for _, line := range envLines {
		line = strings.TrimSpace(line)
		if line != "" && strings.Contains(line, "=") {
			env = append(env, line)
		}
	}

	return env, nil
}

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
