package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"titan/internal/actions"
	"titan/internal/core"
	"titan/internal/proxy"
	"titan/internal/tasks"
	"titan/internal/utils"
	"titan/pkg/flags"
	"titan/pkg/types"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repoRunner := func(action types.Action) func(vars ...any) error {
		return func(vars ...any) error {
			options := core.ContainerOptions{
				CommandAction: action,
				ConfigPath:    vars[0].(string),
			}
			container := core.NewContainer(options)
			defer container.CleanUp()

			processRepositoryCommand(container)
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
						CommandAction: utils.PROXY_SERVER,
						Profile:       vars[1].(string),
						ConfigPath:    vars[0].(string),
					}
					container := core.NewContainer(options)
					defer container.CleanUp()

					processProxyCommand(ctx, container)
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
		fmt.Printf("Received signal: %v. Initiating shutdown...\n", sig)
		cancel()
	}()

	appComands := flags.NewAppCommands(&commandOptions)
	err := appComands.Run()
	if err != nil {
		utils.PrintlnRed(fmt.Sprintf("Error parsing flags: %v", err))
		os.Exit(1)
	}
}

func processProxyCommand(ctx context.Context, container *core.Container) {

	// Create unbuffered error channel for proxy server and tasks
	errorChannel := make(chan error)

	profileData, found := container.ConfigData.Config.Server.Profiles[container.Command.Profile]
	if !found {
		// TODO - need to deal with this potential error properly
		fmt.Printf("profile %+v not found in config\n", container.Command.Profile)
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
		fmt.Printf("fatal error: %v\n", err)
	case <-ctx.Done():
		fmt.Println("context canceled, shutting down")
	}
	fmt.Println("All workers have stopped.")
}

func processRepositoryCommand(container *core.Container) {
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
	var a []actions.Action
	for _, actionToCheck := range availableActions {
		if actionToCheck.ShouldExecute(container.Command.Action) {
			a = append(a, actionToCheck)
		}
	}

	// Run actions concurrently for each repo
	// container.ErrorChannel = make(chan error, len(container.ConfigData.Config.Repositories))
	for _, repository := range container.ConfigData.Config.Repositories {
		wg.Go(func() {
			repoName := repoName(repository)
			// Run actions one after the other. Those should be ordered in the array
			for _, act := range a {
				err := act.Execute(repository, repoName, container.SharedEnvironment)
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
		fmt.Println("WaitGroup finished!")
	case err := <-errorChannel:
		fmt.Printf("Error encountered: %v. Initiating shutdown...\n", err)
		// cancel() // Cancel the context to stop the workers
	}

	close(errorChannel)

	var errors []error
	for err := range errorChannel {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		utils.PrintlnRed("Some actions failed:")
		for _, err := range errors {
			utils.PrintlnRed(fmt.Sprintf("  - %v", err))
		}
		os.Exit(1)
	} else {
		fmt.Println("All actions completed successfully")
	}
}
func repoName(repository string) string {
	split := strings.Split(repository, "/")
	return split[len(split)-1]
}
