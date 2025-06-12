package actions

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"titan/internal/utils"
)

type InstallAction struct {
	name     string
	commands []utils.Command
}

func NewInstallAction() InstallAction {
	return InstallAction{
		name:     "install",
		commands: []utils.Command{utils.ALL, utils.INSTALL},
	}
}

func (ia InstallAction) Name() string {
	return ia.name
}

func (ia InstallAction) ShouldExecute(command utils.Command) bool {
	return slices.Contains(ia.commands, command)
}

func (ia InstallAction) Execute(repoPath string, projectName string, env []string) error {
	// create temp shell script
	script := fmt.Sprintf(`#!/bin/bash
		set -e
		echo 'INSTALL ACTION in %v'
		cd %v
		pnpm install`, projectName, repoPath)
	// Write script to temp file
	tmpFile, err := utils.CreateTempFile("", "install-action-*.sh", script)
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
