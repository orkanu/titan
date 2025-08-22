package actions

import (
	"fmt"
	"log/slog"
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
		commands: []types.Action{utils.REPO_ALL, utils.INSTALL},
	}
}

func (ia InstallAction) Name() string {
	return ia.name
}

func (ia InstallAction) ShouldExecute(command types.Action) bool {
	return slices.Contains(ia.commands, command)
}

func (ia InstallAction) Execute(logger *slog.Logger, repoPath string, projectName string, env []string) error {
	options := utils.NewExecCommandOptions(env, repoPath, "pnpm", "install", "--frozen-lockfile", "--prefer-offline")
	logger.Info("Action [install]", "project", projectName)
	if err := utils.ExecCommand(options); err != nil {
		return fmt.Errorf("Error executing install action script: %v", err)
	}
	return nil
}
