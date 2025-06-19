package actions

import (
	"fmt"
	"slices"
	"titan/internal/utils"
	"titan/pkg/types"
)

type InstallAction struct {
	name     string
	commands []types.Action
}

func NewInstallAction() InstallAction {
	return InstallAction{
		name:     "install",
		commands: []types.Action{utils.ALL, utils.INSTALL},
	}
}

func (ia InstallAction) Name() string {
	return ia.name
}

func (ia InstallAction) ShouldExecute(command types.Action) bool {
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
