package main

import (
	"context"
	"fmt"
	"log"
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
	"titan/pkg/config"
	"titan/pkg/flags"
)

func main() {
	container := container.NewContainer()
	if err := flags.ParseFlags(container); err != nil {
		utils.PrintlnRed(fmt.Sprintf("Error parsing flags: %v", err))
		os.Exit(1)
	}

	if err := config.NewConfig(container); err != nil {
		utils.PrintlnRed(fmt.Sprintf("Error retrieving configuration: %v", err))
		os.Exit(1)
	}

	// Setup nvm and pnpm environment
	env, err := utils.CaptureEnvironment(container.ConfigData.Config.Versions)
	if err != nil {
		utils.PrintlnRed(fmt.Sprintf("Error setting up shared bash environment: %v", err))
		os.Exit(1)
	}
	container.SharedEnvironment = env

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

	// Create a WaitGroup to wait for all workers to finish
	var wg sync.WaitGroup

	// Start tasks
	tasks.StartTasks(container, &wg)
	// Start proxy
	proxy.StartProxy(container, &wg)

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
	wg.Wait()
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
	var wg sync.WaitGroup
	container.ErrorChannel = make(chan error, len(container.ConfigData.Config.Repositories))
	for _, repository := range container.ConfigData.Config.Repositories {
		wg.Add(1)
		go func() {
			defer wg.Done()
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
	wg.Wait()
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
