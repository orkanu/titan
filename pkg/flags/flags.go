package flags

import (
	"flag"
	"log"
	"titan/internal/container"
	"titan/internal/utils"
	"titan/pkg/types"
)

type Flags struct {
	Command    types.Action
	ConfigPath string
}

// ParseFlags will create and parse the CLI flags
func ParseFlags(container *container.Container) error {
	// String that contains the configured configuration path
	var configPath string
	flag.StringVar(&configPath, "c", "./titan.yaml", "path to config file")

	flag.Parse()
	args := flag.Args()
	if len(args) != 0 {
		cmd, args := args[0], args[1:]
		// Parse flags based on command
		switch cmd {
		case "fetch":
			fetch(container, args)
		case "install":
			install(container, args)
		case "build":
			build(container, args)
		case "clean":
			clean(container, args)
		case "all":
			all(container, args)
		case "serve":
			serve(container, args)
		default:
			log.Fatal("Please specify a valid subcommand.")
		}
	} else {
		log.Fatal("Please specify a subcommand.")
	}

	// Validate the path first
	if err := utils.CheckIsFile(configPath); err != nil {
		return err
	}

	container.ConfigData.ConfigFilePath = configPath

	return nil
}

func registerGlobalFlags(fset *flag.FlagSet) {
	flag.VisitAll(func(f *flag.Flag) {
		fset.Var(f.Value, f.Name, f.Usage)
	})
}

func fetch(container *container.Container, args []string) {
	flag := flag.NewFlagSet("fetch", flag.ExitOnError)
	registerGlobalFlags(flag)
	flag.Parse(args)
	container.Command.Action = utils.FETCH
}

func install(container *container.Container, args []string) {
	flag := flag.NewFlagSet("install", flag.ExitOnError)
	registerGlobalFlags(flag)
	flag.Parse(args)
	container.Command.Action = utils.INSTALL
}

func build(container *container.Container, args []string) {
	flag := flag.NewFlagSet("build", flag.ExitOnError)
	registerGlobalFlags(flag)
	flag.Parse(args)
	container.Command.Action = utils.BUILD
}

func clean(container *container.Container, args []string) {
	flag := flag.NewFlagSet("clean", flag.ExitOnError)
	registerGlobalFlags(flag)
	flag.Parse(args)
	container.Command.Action = utils.CLEAN
}

func all(container *container.Container, args []string) {
	flag := flag.NewFlagSet("all", flag.ExitOnError)
	registerGlobalFlags(flag)
	flag.Parse(args)
	container.Command.Action = utils.REPO_ALL
}

func serve(container *container.Container, args []string) {
	flag := flag.NewFlagSet("serve", flag.ExitOnError)
	var profile = flag.String("p", "", "profile to use")
	registerGlobalFlags(flag)
	flag.Parse(args)
	if *profile == "" {
		log.Fatal("missing profile flag in serve command")
	}
	container.Command.Action = utils.PROXY_SERVER
	container.Command.Profile = *profile
}
