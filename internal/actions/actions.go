package actions

import (
	"titan/internal/utils"
)

// Action defines what an action is
type Action interface {
	// Name gives the name of the action
	Name() string
	// ShouldExecute evaluates is the action has to be executed based on the command requested
	ShouldExecute(command utils.Command) bool
	// Execute executes the action
	Execute(repoPath string, projectName string, env []string) error
}
