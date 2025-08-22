package actions

import (
	"fmt"
	"log/slog"
	"slices"
	"titan/internal/utils"
	"titan/pkg/types"
)

// Fetch action
type FetchAction struct {
	name     string
	commands []types.Action
}

func NewFetchAction() FetchAction {
	return FetchAction{
		name:     "fetch",
		commands: []types.Action{utils.REPO_ALL, utils.FETCH},
	}
}

func (fa FetchAction) Name() string {
	return fa.name
}

func (fa FetchAction) ShouldExecute(command types.Action) bool {
	return slices.Contains(fa.commands, command)
}

func (fa FetchAction) Execute(logger *slog.Logger, repoPath string, projectName string, env []string) error {
	// create temp shell script
	script := `#!/bin/bash
		set -e
		git fetch -p && git pull
  		git fetch --tags --force && git fetch --prune --prune-tags`
	logger.Info("Action [fetch]", "project", projectName)
	if err := utils.ExecScript(script, env, repoPath); err != nil {
		return fmt.Errorf("Error executing fetch action script: %v", err)
	}
	return nil
}
