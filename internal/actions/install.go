package actions

import (
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

func (ia InstallAction) Execute(options *ExecOptions) error {
	defaultScript := "pnpm install --frozen-lockfile --prefer-offline"
	scriptFromConfig := getScriptFromConfig(ia.name, options.repoAction, nil, defaultScript, options.logger)

	return executeScript(ia.name, scriptFromConfig, options.logger, options.repoPath, options.projectName, options.env)
}
