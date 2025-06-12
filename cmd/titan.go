package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"titan/internal/actions"
	"titan/internal/config"
	"titan/internal/utils"
	"titan/pkg/flags"
)

func main() {
	flagsData, err := flags.ParseFlags()
	if err != nil {
		utils.PrintlnRed(fmt.Sprintf("Error parsing flags: %v", err))
		os.Exit(1)
	}
	cfg, err := config.NewConfig(flagsData.ConfigPath)
	if err != nil {
		utils.PrintlnRed(fmt.Sprintf("Error retrieving configuration: %v", err))
		os.Exit(1)
	}

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
		if actionToCheck.ShouldExecute(flagsData.Command) {
			a = append(a, actionToCheck)
		}
	}

	// Setup nvm and pnpm environment
	env, err := utils.CaptureEnvironment(cfg.Versions)
	if err != nil {
		utils.PrintlnRed(fmt.Sprintf("Error setting up shared bash environment: %v", err))
		os.Exit(1)
	}

	// Run actions concurrently for each repo
	var wg sync.WaitGroup
	errorChan := make(chan error, len(cfg.Repositories))
	for _, repository := range cfg.Repositories {
		wg.Add(1)
		go func() {
			defer wg.Done()
			repoPath := repoFullPath(cfg.BasePath, repository)
			// Run actions one after the other. Those should be ordered in the array
			for _, act := range a {
				err := act.Execute(repoPath, repository, env)
				if err != nil {
					errorChan <- err
					// Stop procession further actions
					break
				}
			}
		}()
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(errorChan)

	var errors []error
	for err := range errorChan {
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

func repoFullPath(base string, repo string) string {
	var b strings.Builder
	b.WriteString(base)
	if !strings.HasSuffix(base, "/") {
		b.WriteString("/")
	}
	b.WriteString(repo)
	return b.String()
}
