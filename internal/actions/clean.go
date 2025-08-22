package actions

import (
	"fmt"
	"log/slog"
	"slices"
	"titan/internal/utils"
	"titan/pkg/types"
)

type CleanAction struct {
	name     string
	commands []types.Action
}

func NewCleanAction() CleanAction {
	return CleanAction{
		name:     "clean",
		commands: []types.Action{utils.REPO_ALL, utils.CLEAN},
	}
}

func (ca CleanAction) Name() string {
	return ca.name
}

func (ca CleanAction) ShouldExecute(command types.Action) bool {
	return slices.Contains(ca.commands, command)
}

func (ca CleanAction) Execute(logger *slog.Logger, repoPath string, projectName string, env []string) error {
	// create temp shell script
	cleanYalc := ""
	if projectName == "cbs-residential-web" {
		cleanYalc = "rm -rf ~/.yalc/packages/@wavelength"
	}
	script := fmt.Sprintf(`#!/bin/bash
		set -e
		find $(pwd) -maxdepth 3 -name "node_modules" -type d -exec rm -rf {} +
		find $(pwd) -maxdepth 3 -name "dist" -type d -exec rm -rf {} +
		%v
		`, cleanYalc)
	logger.Info("Action [clean]", "project", projectName)
	if err := utils.ExecScript(script, env, repoPath); err != nil {
		return fmt.Errorf("Error executing clean action script: %v", err)
	}
	return nil
}
