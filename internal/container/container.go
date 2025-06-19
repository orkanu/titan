package container

import (
	"titan/pkg/types"
)

// Command holds the data required for the requested command
type Command struct {
	// Action to execute
	Action types.Action
	// Profile is required for server proxy action
	Profile string
}

type Configuration struct {
	// ConfigFilePath path where the config file is available
	ConfigFilePath string
	// Config holds the configuration data loaded from an YAML file
	Config *types.Config
}

// Container holds data that can be used across the app
type Container struct {
	// ConfigData holds the configuration data loaded from an YAML file
	ConfigData Configuration
	// Command holds the data required for the requested command
	Command Command
}

// NewContainer retuns a Container
func NewContainer() *Container {
	return &Container{
		Command: Command{},
	}
}
