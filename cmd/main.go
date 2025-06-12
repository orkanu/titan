package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"titan/internal/actions"
	"titan/internal/utils"
)

func main() {
	flagsData, err := ParseFlags()
	if err != nil {
		utils.PrintlnRed(fmt.Sprintf("Error parsing flags: %v", err))
		os.Exit(1)
	}
	cfg, err := NewConfig(flagsData.ConfigPath)
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
	env, err := captureEnvironment(cfg.Versions)
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

	// Get subcomand from flags:
	// 	(fetch) should fetch and pull from remote for each project
	// 	(clean) should remove node_modules, dist and yalc folders on each of the different projects
	// 	(install) should install deps for each of the different projects
	// 	(build) should build each of the different projects
	// 	(all) should fetch, clean, install and build each of the different projects
	// If no subcomand specified, assume it is (all)
	//
	// For install & build we need to do the following:
	// 	Check/install node (use version from config)
	// 	Check/install pnpm (use version from config)
	//
	// Should use goroutines to run everything in parallel
	// Should use prefix when logging to console, for instance:
	// 	(FETCH): blah blah blah
	// 	(INSTALL): blah blah blah more
	// Evaluate using files for logging. Each action should use its own file
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

// captureEnvironment sets up NVM and pnpm and captures the resulting environment
func captureEnvironment(versions Versions) ([]string, error) {
	// Source nvm and install the desired nvm and pnpm versions
	scriptContents := fmt.Sprintf(`
		source ~/.nvm/nvm.sh &&
		nvm install %v &&
		nvm use %v &&
		npm i -g pnpm@%v &&
		env
		`, versions.Node, versions.Node, versions.PNPM)
	setupCmd := exec.Command("bash", "-c", scriptContents)
	output, err := setupCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to setup environment: %w", err)
	}
	// Parese the environment output
	envLines := strings.Split(string(output), "\n")
	var env []string
	for _, line := range envLines {
		line = strings.TrimSpace(line)
		if line != "" && strings.Contains(line, "=") {
			env = append(env, line)
		}
	}

	return env, nil
}
