package utils

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"titan/pkg/types"
)

// CaptureEnvironment sets up NVM and pnpm and captures the resulting environment
func CaptureEnvironment(versions types.Versions) ([]string, error) {
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
func ExecScript(script string, env []string, dir string) error {
	// Write script to temp file
	tmpFile, err := CreateTempFile("", "titan-action-*.sh", script)
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	// Execute the script
	options := NewExecCommandOptions(env, dir, "bash", tmpFile.Name())
	return ExecCommand(options)
}

// ExecCommandOptions holds options for ExecCommand functionality
type ExecCommandOptions struct {
	Env     []string
	Dir     string
	Command string
	Args    []string
}

// NewExecCommandOptions returns an ExecCommandOptions struct
func NewExecCommandOptions(env []string, dir string, command string, args ...string) ExecCommandOptions {
	return ExecCommandOptions{
		Env:     env,
		Dir:     dir,
		Command: command,
		Args:    args,
	}
}

// ExecCommand is a utility function that executes simple shell commands
func ExecCommand(options ExecCommandOptions) error {
	workingDir := getPathWithUserHome(options.Dir)
	cmd := exec.Command(options.Command, options.Args...)
	cmd.Dir = workingDir
	cmd.Env = options.Env
	cmd.Stderr = cmd.Stdout // redirect stderr to stdout
	// using pipes should be quicker than capturing the full output with cmd.Output()
	stdoutPipe, _ := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		return err
	}

	// Stream output directly to stdout
	go func() {
		_, _ = io.Copy(os.Stdout, stdoutPipe)
	}()

	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}
