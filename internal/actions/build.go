package actions

import (
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

func (ba BuildAction) Execute(repoActions map[string]types.RepoAction, logger *slog.Logger, repoPath string, projectName string, env []string) error {
	defaultScript := "pnpm run build:local"
	scriptFromConfig := getScriptFromConfig(ba.name, repoActions, nil, defaultScript, logger)

	return executeScript(ba.name, scriptFromConfig, logger, repoPath, projectName, env)
}
