package actions

import (
	"log/slog"
	"titan/pkg/types"
)

// Action defines what an action is
type Action interface {
	// Name gives the name of the action
	Name() string
	// ShouldExecute evaluates is the action has to be executed based on the command requested
	ShouldExecute(command types.Action) bool
	// Execute executes the action
	Execute(logger *slog.Logger, repoPath string, projectName string, env []string) error
}
