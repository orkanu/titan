package actions

import (
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

func (ca CleanAction) Execute(repoActions map[string]types.RepoAction, logger *slog.Logger, repoPath string, projectName string, env []string) error {

	defaultScript := `
		find $(pwd) -maxdepth 3 -name "node_modules" -type d -exec rm -rf {} +
        find $(pwd) -maxdepth 3 -name "dist" -type d -exec rm -rf {} +
	`
	ctx := map[string]any{
		"projectName": projectName,
	}
	scriptFromConfig := getScriptFromConfig(ca.name, repoActions, ctx, defaultScript, logger)

	return executeScript(ca.name, scriptFromConfig, logger, repoPath, projectName, env)
}
