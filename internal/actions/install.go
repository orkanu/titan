package actions

import (
	"fmt"
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
	options := utils.NewExecCommandOptions(env, repoPath, "pnpm", "install", "--frozen-lockfile", "--prefer-offline")
	fmt.Printf("Action [install] on project [%v]\n", projectName)
	if err := utils.ExecCommand(options); err != nil {
		return fmt.Errorf("Error executing install action script: %v", err)
	}
	return nil
}
