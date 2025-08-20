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
	"titan/internal/container"
	"titan/internal/proxy"
	"titan/internal/tasks"
	"titan/internal/utils"
	"titan/pkg/flags"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	flagData, err := flags.ParseFlags()
	if err != nil {
		utils.PrintlnRed(fmt.Sprintf("Error parsing flags: %v", err))
		os.Exit(1)
	}

	// Handle SIGINT/SIGTERM
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigCh
		fmt.Printf("Received signal: %v. Initiating shutdown...\n", sig)
		cancel()
	}()

	options := container.ContainerOptions{
		CommandAction: flagData.Command,
		Profile:       flagData.Profile,
		ConfigPath:    flagData.ConfigPath,
	}
	container := container.NewContainer(options)
	defer container.CleanUp()

	if container.Command.Action == utils.PROXY_SERVER {
		processProxyCommand(ctx, container)
	} else {
		processRepositoryCommand(container)
	}
}

func processProxyCommand(ctx context.Context, container *container.Container) {

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

func processRepositoryCommand(container *container.Container) {
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
