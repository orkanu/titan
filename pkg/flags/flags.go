package flags

import (
	"errors"
	"flag"
	"os"
	"titan/internal/utils"
)

type Command struct {
	Runner     func(vars ...any) error
	Subcommand []Command
}

type AppCommands struct {
	commands map[string]Command
}

type AppCommandsOptions struct {
	Commands map[string]Command
}

func NewAppCommands(options *AppCommandsOptions) *AppCommands {
	return &AppCommands{
		commands: options.Commands,
	}
}

func (ac *AppCommands) Run() error {
	// String that contains the configured configuration path
	var configPath string
	flag.StringVar(&configPath, "c", "./titan.yaml", "path to config file")

	// Define subcommands
	fetchCmd := flag.NewFlagSet("fetch", flag.ExitOnError)
	registerGlobalFlags(fetchCmd)
	installCmd := flag.NewFlagSet("install", flag.ExitOnError)
	registerGlobalFlags(installCmd)
	buildCmd := flag.NewFlagSet("build", flag.ExitOnError)
	registerGlobalFlags(buildCmd)
	cleanCmd := flag.NewFlagSet("clean", flag.ExitOnError)
	registerGlobalFlags(cleanCmd)
	allCmd := flag.NewFlagSet("all", flag.ExitOnError)
	registerGlobalFlags(allCmd)
	serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)
	registerGlobalFlags(serveCmd)
	var profile string
	serveCmd.StringVar(&profile, "p", "", "profile to use")

	if len(os.Args) < 2 {
		return errors.New("Please specify a subcommand.")
	}

	runCommand := func(name string, vars ...any) error {
		// Validate the config file path first
		if err := utils.CheckIsFile(vars[0].(string)); err != nil {
			return err
		}
		// Validate profile is present for serve command
		if name == "serve" && vars[1].(string) == "" {
			return errors.New("missing profile")
		}

		if command, ok := ac.commands[name]; ok {
			command.Runner(vars...)
		} else {
			return errors.New("No command runner found for fetch")
		}
		return nil
	}
	// Parse flags based on command
	switch os.Args[1] {
	case "fetch":
		fetchCmd.Parse(os.Args[2:])
		return runCommand("fetch", configPath)
	case "install":
		installCmd.Parse(os.Args[2:])
		return runCommand("install", configPath)
	case "build":
		buildCmd.Parse(os.Args[2:])
		return runCommand("build", configPath)
	case "clean":
		cleanCmd.Parse(os.Args[2:])
		return runCommand("clean", configPath)
	case "all":
		allCmd.Parse(os.Args[2:])
		return runCommand("all", configPath)
	case "serve":
		serveCmd.Parse(os.Args[2:])
		return runCommand("serve", configPath, profile)
	default:
		return errors.New("Please specify a valid subcommand.")
	}
}

func registerGlobalFlags(fset *flag.FlagSet) {
	flag.VisitAll(func(f *flag.Flag) {
		fset.Var(f.Value, f.Name, f.Usage)
	})
}
