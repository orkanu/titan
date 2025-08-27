package actions

import (
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

func (fa FetchAction) Execute(repoActions map[string]types.RepoAction, logger *slog.Logger, repoPath string, projectName string, env []string) error {
	defaultScript := `
		git fetch -p && git pull
		git fetch --tags --force && git fetch --prune --prune-tags
	`
	scriptFromConfig := getScriptFromConfig(fa.name, repoActions, nil, defaultScript, logger)

	return executeScript(fa.name, scriptFromConfig, logger, repoPath, projectName, env)
}
