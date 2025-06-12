package actions

import (
	"fmt"
	"slices"
	"titan/internal/utils"
)

type BuildAction struct {
	name     string
	commands []utils.Command
}

func NewBuildAction() BuildAction {
	return BuildAction{
		name:     "build",
		commands: []utils.Command{utils.ALL, utils.BUILD},
	}
}

func (ba BuildAction) Name() string {
	return ba.name
}

func (ba BuildAction) ShouldExecute(command utils.Command) bool {
	return slices.Contains(ba.commands, command)
}

func (ba BuildAction) Execute(repoPath string, env []string) {
	fmt.Printf("Action %v executed on repo %v\n", ba.Name(), repoPath)
}
