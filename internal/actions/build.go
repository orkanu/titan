package actions

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"titan/internal/utils"
)

type BuildAction struct {
	name     string
	commands []utils.Command
}

func NewBuildAction() BuildAction {
	return BuildAction{
		name:     "build",
		commands: []utils.Command{utils.ALL, utils.BUILD},
	}
}

func (ba BuildAction) Name() string {
	return ba.name
}

func (ba BuildAction) ShouldExecute(command utils.Command) bool {
	return slices.Contains(ba.commands, command)
}

func (ba BuildAction) Execute(repoPath string, projectName string, env []string) error {
	// create temp shell script
	script := fmt.Sprintf(`#!/bin/bash
		set -e
		echo 'BUILD ACTION in %v'
		cd %v
		pnpm run build:local`, projectName, repoPath)
	// Write script to temp file
	tmpFile, err := utils.CreateTempFile("", "build-action-*.sh", script)
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
