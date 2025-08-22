package core

import (
	"log/slog"
	"os"
	"titan/internal/utils"
	"titan/pkg/config"
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
	// Profile is the profile config to use
	Profile types.Profile
}

// Container holds data that can be used across the app
type Container struct {
	// Logger holds the logger instance to use across the app
	Logger *slog.Logger
	// ConfigData holds the configuration data loaded from an YAML file
	ConfigData Configuration
	// Command holds the data required for the requested command
	Command Command
	// SharedEnvironment
	SharedEnvironment []string
}

type ContainerOptions struct {
	Logger        *slog.Logger
	CommandAction types.Action
	Profile       string
	ConfigPath    string
}

// NewContainer retuns a Container
func NewContainer(options ContainerOptions) *Container {
	// Load configuration
	config, err := config.NewConfig(options.ConfigPath)
	if err != nil {
		options.Logger.Error("Error retrieving configuration", "error", err)
		os.Exit(1)
	}
	// Setup nvm and pnpm to use as environment on other shell executions
	env, err := utils.CaptureEnvironment(config.Versions)
	if err != nil {
		options.Logger.Error("Error setting up shared bash environment", "error", err)
		os.Exit(1)
	}

	// cleanUpFuncs = addCleanUpFunc(cleanUpFuncs, "sample cleanup name", func() error {
	// 	return nil
	// })

	return &Container{
		Logger: options.Logger,
		Command: Command{
			Action:  options.CommandAction,
			Profile: options.Profile,
		},
		ConfigData: Configuration{
			ConfigFilePath: options.ConfigPath,
			Config:         config,
		},
		SharedEnvironment: env,
	}
}
