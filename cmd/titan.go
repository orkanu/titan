package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
	"titan/internal/actions"
	"titan/internal/core"
	"titan/internal/proxy"
	"titan/internal/tasks"
	"titan/internal/utils"
	"titan/pkg/flags"
	"titan/pkg/types"

	"github.com/lmittmann/tint"
)

func main() {
	w := os.Stderr

	// logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	// slog.SetDefault(logger)
	logger := slog.New(tint.NewHandler(w, nil))
	// Set global logger with custom options
	slog.SetDefault(slog.New(
		tint.NewHandler(w, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
		}),
	))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repoRunner := func(action types.Action) func(vars ...any) error {
		return func(vars ...any) error {
			options := core.ContainerOptions{
				Logger:        logger,
				CommandAction: action,
				ConfigPath:    vars[0].(string),
			}
			container := core.NewContainer(options)

			processCommand(container)
			return nil
		}
	}
	commandOptions := flags.AppCommandsOptions{
		Commands: map[string]flags.Command{
			"fetch":   {Runner: repoRunner(utils.FETCH)},
			"install": {Runner: repoRunner(utils.INSTALL)},
			"build":   {Runner: repoRunner(utils.BUILD)},
			"clean":   {Runner: repoRunner(utils.CLEAN)},
			"all":     {Runner: repoRunner(utils.REPO_ALL)},
			"serve": {
				Runner: func(vars ...any) error {
					options := core.ContainerOptions{
						Logger:        logger,
						CommandAction: utils.PROXY_SERVER,
						Profile:       vars[1].(string),
						ConfigPath:    vars[0].(string),
					}
					container := core.NewContainer(options)

					processProxy(ctx, container)
					return nil
				},
			},
			"help": {
				Runner: func(_ ...any) error {
					utils.PrintlnWhite("TITAN - Wee CLI app that allows perform some operations against a project as well as start a proxy server")
					utils.PrintlnWhite("to run a bunch of services under the same host")
					utils.PrintlnBlack("")
					utils.PrintlnGreen("Usage:")
					utils.PrintlnGreen("   fetch   - performs a git fetch on the configured project/s")
					utils.PrintlnGreen("   install - performs a pnpm install on the configured project/s")
					utils.PrintlnGreen("   build   - performs a pnpm run build:local on the configured project/s")
					utils.PrintlnGreen("   clean   - performs a clean up of the node_modules and dist folders on the configured project/s")
					utils.PrintlnGreen("   all     - performs all of the above")
					utils.PrintlnGreen("   serve   - starts a proxy server based on configuration. NOTE: required flag \"-p\" to specify a profile to use")
					utils.PrintlnBlack("")
					utils.PrintlnCyan("To run any of the comands, it requires a configuration file (default \"titan.yaml\" in the same place where the")
					utils.PrintlnCyan("binary is run). Using the -c flag, we can specify a different config file location")

					return nil
				},
			},
		},
	}

	// Handle SIGINT/SIGTERM
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigCh
		logger.Info("inititation shutdown - received signal", "signal", sig)
		cancel()
	}()

	appComands := flags.NewAppCommands(&commandOptions)
	err := appComands.Run()
	if err != nil {
		logger.Error("failed parsing flags", "error", err)
		os.Exit(1)
	}
}

func processProxy(ctx context.Context, container *core.Container) {

	// Create unbuffered error channel for proxy server and tasks
	errorChannel := make(chan error)

	profileData, found := container.ConfigData.Config.Server.Profiles[container.Command.Profile]
	if !found {
		// TODO - need to deal with this potential error properly
		container.Logger.Info("profile not found in config", "profile", container.Command.Profile)
		os.Exit(1)
	}
	container.ConfigData.Profile = profileData

	// Start tasks
	tasks.StartTasks(errorChannel, container)
	// Start proxy
	proxy.StartProxy(errorChannel, container)

	// Wait for error or shutdown
	select {
	case err := <-errorChannel:
		container.Logger.Error("fatal error", "error", err)
	case <-ctx.Done():
		container.Logger.Info("context canceled, shutting down")
	}
	container.Logger.Info("all workers have stopped")
}

func processCommand(container *core.Container) {
	// Create a WaitGroup to wait for all workers to finish
	var wg sync.WaitGroup

	// Channel to do not get stuck in the select below if the error or signal channels do not receive anything
	// but the repository commands are finished
	waitCh := make(chan struct{})

	// Create buffered error channel for repository actions
	errorChannel := make(chan error, len(container.ConfigData.Config.Repositories))

	// Slice with all the available actions
	availableActions := []actions.Action{
		actions.NewFetchAction(),
		actions.NewCleanAction(),
		actions.NewInstallAction(),
		actions.NewBuildAction(),
	}

	// Get only actions required based on command passed to the Titan
	var actionList []actions.Action
	for _, actionToCheck := range availableActions {
		if actionToCheck.ShouldExecute(container.Command.Action) {
			actionList = append(actionList, actionToCheck)
		}
	}

	// Run actions concurrently for each repo
	// container.ErrorChannel = make(chan error, len(container.ConfigData.Config.Repositories))
	for _, repository := range container.ConfigData.Config.Repositories {
		wg.Go(func() {
			repoName := repoName(repository)
			repoActions := container.ConfigData.Config.RepoActions
			sharedEnv := container.SharedEnvironment
			// Run actions one after the other. Those should be ordered in the array
			for _, action := range actionList {
				err := action.Execute(repoActions, container.Logger, repository, repoName, sharedEnv)
				if err != nil {
					errorChannel <- err
					// Stop procession further actions
					break
				}
			}
		})
	}

	// Wait for all workers to finish
	wg.Wait()
	close(waitCh)

	// Select loop to wait for signals or errors
	select {
	case <-waitCh:
		container.Logger.Debug("waitGroup finished")
	case err := <-errorChannel:
		container.Logger.Error("fatal error", "error", err)
		// cancel() // Cancel the context to stop the workers
	}

	close(errorChannel)

	var errors []error
	for err := range errorChannel {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		container.Logger.Error("some actions failed:")
		for _, err := range errors {
			container.Logger.Error(fmt.Sprintf("  - %v", err))
		}
		os.Exit(1)
	} else {
		container.Logger.Debug("all actions completed")
	}
}
func repoName(repository string) string {
	split := strings.Split(repository, "/")
	return split[len(split)-1]
}
