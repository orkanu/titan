package actions

import (
	"fmt"
	"log/slog"
	"slices"
	"titan/internal/utils"
	"titan/pkg/types"
)

type BuildAction struct {
	name     string
	commands []types.Action
}

func NewBuildAction() BuildAction {
	return BuildAction{
		name:     "build",
		commands: []types.Action{utils.REPO_ALL, utils.BUILD},
	}
}

func (ba BuildAction) Name() string {
	return ba.name
}

func (ba BuildAction) ShouldExecute(command types.Action) bool {
	return slices.Contains(ba.commands, command)
}

func (ba BuildAction) Execute(logger *slog.Logger, repoPath string, projectName string, env []string) error {
	options := utils.NewExecCommandOptions(env, repoPath, "pnpm", "run", "build:local")
	logger.Info("Action [build]", "project", projectName)
	if err := utils.ExecCommand(options); err != nil {
		return fmt.Errorf("Error executing build action script: %v", err)
	}
	return nil
}
