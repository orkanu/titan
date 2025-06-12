package actions

import (
	"fmt"
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
	if err := utils.ExecScript(script, env); err != nil {
		return fmt.Errorf("Error executing build action script: %v", err)
	}
	return nil
}
