package actions

import (
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

func (ca CleanAction) Execute(options *ExecOptions) error {

	defaultScript := `
		find $(pwd) -maxdepth 3 -name "node_modules" -type d -exec rm -rf {} +
        find $(pwd) -maxdepth 3 -name "dist" -type d -exec rm -rf {} +
	`
	ctx := map[string]any{
		"projectName": options.projectName,
	}
	scriptFromConfig := getScriptFromConfig(ca.name, options.repoAction, ctx, defaultScript, options.logger)

	return executeScript(ca.name, scriptFromConfig, options.logger, options.repoPath, options.projectName, options.env)
}
