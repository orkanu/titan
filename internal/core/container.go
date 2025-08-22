package core

import (
	"context"
	"fmt"
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
	// Context keeps hold of the context to allow cancellations
	Context context.Context
	// ConfigData holds the configuration data loaded from an YAML file
	ConfigData Configuration
	// Command holds the data required for the requested command
	Command Command
	// SharedEnvironment
	SharedEnvironment []string
	// CleanUpFuncs is a list of functions to clean up the container
	CleanUpFuncs []func() error
}

type ContainerOptions struct {
	CommandAction types.Action
	Profile       string
	ConfigPath    string
}

// NewContainer retuns a Container
func NewContainer(options ContainerOptions) *Container {
	var cleanUpFuncs []func() error
	// Load configuration
	config, err := config.NewConfig(options.ConfigPath)
	if err != nil {
		utils.PrintlnRed(fmt.Sprintf("Error retrieving configuration: %v", err))
		os.Exit(1)
	}
	// Setup nvm and pnpm to use as environment on other shell executions
	env, err := utils.CaptureEnvironment(config.Versions)
	if err != nil {
		utils.PrintlnRed(fmt.Sprintf("Error setting up shared bash environment: %v", err))
		os.Exit(1)
	}

	// Create error channel. For proxy server we use an unbuffered on or buffered for repository actions
	// var errorChannel chan error
	// if options.CommandAction == utils.PROXY_SERVER {
	// 	errorChannel = make(chan error)
	// 	cleanUpFuncs = addCleanUpFunc(cleanUpFuncs, "close error channel", func() error {
	// 		close(errorChannel)
	// 		return nil
	// 	})
	// } else {
	// 	errorChannel = make(chan error, len(config.Repositories))
	// 	// we do not close the error channel here cause we need to do it after we process actions to read all
	// 	// errors added in it. If we close it there and then via the cleanup function, we get a panic here
	// 	// as the channel was already closed
	// }

	// cleanUpFuncs = addCleanUpFunc(cleanUpFuncs, "sample cleanup name", func() error {
	// 	return nil
	// })

	return &Container{
		Command: Command{
			Action:  options.CommandAction,
			Profile: options.Profile,
		},
		ConfigData: Configuration{
			ConfigFilePath: options.ConfigPath,
			Config:         config,
		},
		SharedEnvironment: env,
		CleanUpFuncs:      cleanUpFuncs,
	}
}

// CleanUp executes all clean up functions available in CleanUpFuncs
func (c Container) CleanUp() {
	for _, clean := range c.CleanUpFuncs {
		if err := clean(); err != nil {
			fmt.Printf("%v\n", err)

			// TODO - I'm doing clean up here. Do I need to exit ungracefuly (i.e. os.Exit(1))
		}
	}
}

func addCleanUpFunc(cleanUpFuncs []func() error, message string, callback func() error) []func() error {
	return append(cleanUpFuncs, func() error {
		fmt.Printf("CleanUp - %v\n", message)
		return callback()
	})
}
