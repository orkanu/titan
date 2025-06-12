package actions

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"titan/internal/utils"
)

// Fetch action
type FetchAction struct {
	name     string
	commands []utils.Command
}

func NewFetchAction() FetchAction {
	return FetchAction{
		name:     "fetch",
		commands: []utils.Command{utils.ALL, utils.FETCH},
	}
}

func (fa FetchAction) Name() string {
	return fa.name
}

func (fa FetchAction) ShouldExecute(command utils.Command) bool {
	return slices.Contains(fa.commands, command)
}

func (fa FetchAction) Execute(repoPath string, projectName string, env []string) error {
	// create temp shell script
	script := fmt.Sprintf(`#!/bin/bash
		set -e
		echo 'FETCH ACTION in %v'
		node -v
		pnpm -v`, projectName)
	// Write script to temp file
	tmpFile, err := utils.CreateTempFile("", "fetch-action-*.sh", script)
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
