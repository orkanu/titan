package container

import (
	"titan/internal/utils"
	"titan/pkg/config"
)

// Command holds the data required for the requested command
type Command struct {
	// Action to execute
	Action utils.Action
	// Profile is required for server proxy action
	Profile string
}

// Container holds data that can be used across the app
type Container struct {
	// Config holds the configuration data loaded from an YAML file
	Config *config.Config
	// Command holds the data required for the requested command
	Command Command
}

// NewContainer retuns a Container
func NewContainer() *Container {
	return &Container{
		Command: Command{},
	}
}
