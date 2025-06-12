package actions

import (
	"fmt"
	"os"
	"os/exec"
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

func (ca CleanAction) Execute(repoPath string, env []string) {
	// create temp shell script
	script := `#!/bin/bash
set -e
echo 'CLEAN ACTION'
node -v
pnpm -v
`
	// Write script to temp file
	tmpFile, err := utils.CreateTempFile("", "clean-action-*.sh", script)
	if err != nil {
		// TODO need to propagate error. Perhaps via channels?
		panic(err)
	}
	defer os.Remove(tmpFile.Name())

	// Execute the script
	cmd := exec.Command("bash", tmpFile.Name())
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		// TODO need to propagate error. Perhaps via channels?
		panic(err)
	}
	fmt.Printf("Action %v executed on repo %v\n", ca.Name(), repoPath)
}
