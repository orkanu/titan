package actions

import (
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

func (ia InstallAction) Execute(repoActions map[string]types.RepoAction, logger *slog.Logger, repoPath string, projectName string, env []string) error {
	defaultScript := "pnpm install --frozen-lockfile --prefer-offline"
	scriptFromConfig := getScriptFromConfig(ia.name, repoActions, nil, defaultScript, logger)

	return executeScript(ia.name, scriptFromConfig, logger, repoPath, projectName, env)
}
