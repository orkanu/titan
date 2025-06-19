package actions

import (
	"slices"
	"titan/internal/utils"
	"titan/pkg/types"
)

type ProxyAction struct {
	name     string
	commands []types.Action
}

func NewProxyAction() ProxyAction {
	return ProxyAction{
		name:     "proxy-server",
		commands: []types.Action{utils.PROXY_SERVER},
	}
}

func (ba ProxyAction) Name() string {
	return ba.name
}

func (ba ProxyAction) ShouldExecute(command types.Action) bool {
	return slices.Contains(ba.commands, command)
}

func (ba ProxyAction) Execute(repoPath string, projectName string, env []string) error {
	// TODO
	return nil
}
