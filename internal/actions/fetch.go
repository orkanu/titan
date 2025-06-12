package actions

import (
	"fmt"
	"slices"
	"titan/internal/utils"
)

// Fetch action
type FetchAction struct {
	name     string
	commands []utils.Command
}

func NewFetchAction() FetchAction {
	return FetchAction{
		name:     "fetch",
		commands: []utils.Command{utils.ALL, utils.FETCH},
	}
}

func (fa FetchAction) Name() string {
	return fa.name
}

func (fa FetchAction) ShouldExecute(command utils.Command) bool {
	return slices.Contains(fa.commands, command)
}

func (fa FetchAction) Execute(repoPath string, env []string) {
	fmt.Printf("Action %v executed on repo %v\n", fa.Name(), repoPath)
}
