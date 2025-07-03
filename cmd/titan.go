package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"titan/internal/actions"
	"titan/internal/container"
	"titan/internal/proxy"
	"titan/internal/tasks"
	"titan/internal/utils"
	"titan/pkg/flags"
)

func main() {
	flagData, err := flags.ParseFlags()
	if err != nil {
		utils.PrintlnRed(fmt.Sprintf("Error parsing flags: %v", err))
		os.Exit(1)
	}

	options := container.ContainerOptions{
		CommandAction: flagData.Command,
		Profile:       flagData.Profile,
		ConfigPath:    flagData.ConfigPath,
	}
	container := container.NewContainer(options)
	defer container.CleanUp()

	if container.Command.Action == utils.PROXY_SERVER {
		processProxyCommand(container)
	} else {
		processRepositoryCommand(container)
	}
}

func processProxyCommand(container *container.Container) {

	profileData, found := container.ConfigData.Config.Server.Profiles[container.Command.Profile]
	if !found {
		// TODO - need to deal with this potential error properly
		fmt.Printf("profile %+v not found in config\n", container.Command.Profile)
		os.Exit(1)
	}
	container.ConfigData.Profile = profileData

	// Create a context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	container.Context = ctx

	// Channel to capture OS signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Channel to collect errors from workers
	container.ErrorChannel = make(chan error)

	// Start tasks
	tasks.StartTasks(container)
	// Start proxy
	proxy.StartProxy(container)

	// Select loop to wait for signals or errors
	select {
	case sig := <-sigCh:
		log.Printf("Received signal: %v. Initiating shutdown...\n", sig)
		cancel() // Cancel the context to stop the workers
	case err := <-container.ErrorChannel:
		log.Printf("Error encountered: %v. Initiating shutdown...\n", err)
		cancel() // Cancel the context to stop the workers
	}

	// Wait for all workers to finish
	container.WaitGroup.Wait()
	fmt.Println("All workers have stopped.")
}

func processRepositoryCommand(container *container.Container) {
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
	container.ErrorChannel = make(chan error, len(container.ConfigData.Config.Repositories))
	for _, repository := range container.ConfigData.Config.Repositories {
		container.WaitGroup.Add(1)
		go func() {
			defer container.WaitGroup.Done()
			repoName := repoName(repository)
			// Run actions one after the other. Those should be ordered in the array
			for _, act := range a {
				err := act.Execute(repository, repoName, container.SharedEnvironment)
				if err != nil {
					container.ErrorChannel <- err
					// Stop procession further actions
					break
				}
			}
		}()
	}

	// Wait for all goroutines to complete
	container.WaitGroup.Wait()
	close(container.ErrorChannel)

	var errors []error
	for err := range container.ErrorChannel {
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
