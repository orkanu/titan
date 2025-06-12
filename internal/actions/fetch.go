package actions

import (
	"fmt"
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
		cd %v
		git fetch -p && git pull
  		git fetch --tags --force && git fetch --prune --prune-tags`, projectName, repoPath)
	if err := utils.ExecScript(script, env); err != nil {
		return fmt.Errorf("Error executing fetch action script: %v", err)
	}
	return nil
}
