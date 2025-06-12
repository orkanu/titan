package actions

import (
	"fmt"
	"slices"
	"titan/internal/utils"
)

type CleanAction struct {
	name     string
	commands []utils.Command
}

func NewCleanAction() CleanAction {
	return CleanAction{
		name:     "clean",
		commands: []utils.Command{utils.ALL, utils.CLEAN},
	}
}

func (ca CleanAction) Name() string {
	return ca.name
}

func (ca CleanAction) ShouldExecute(command utils.Command) bool {
	return slices.Contains(ca.commands, command)
}

func (ca CleanAction) Execute(repoPath string, projectName string, env []string) error {
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
	fmt.Printf("Action [clean] on project [%v]\n", projectName)
	if err := utils.ExecScript(script, env, repoPath); err != nil {
		return fmt.Errorf("Error executing clean action script: %v", err)
	}
	return nil
}
