package flags

import (
	"flag"
	"log"
	"titan/internal/utils"
	"titan/pkg/types"
)

type FlagData struct {
	Command    types.Action
	Profile    string
	ConfigPath string
}

// ParseFlags will create and parse the CLI flags
func ParseFlags() (*FlagData, error) {
	flagData := &FlagData{}
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
			command := fetch(args)
			flagData.Command = command
		case "install":
			command := install(args)
			flagData.Command = command
		case "build":
			command := build(args)
			flagData.Command = command
		case "clean":
			command := clean(args)
			flagData.Command = command
		case "all":
			command := all(args)
			flagData.Command = command
		case "serve":
			command, profile := serve(args)
			flagData.Command = command
			flagData.Profile = profile
		default:
			log.Fatal("Please specify a valid subcommand.")
		}
	} else {
		log.Fatal("Please specify a subcommand.")
	}

	// Validate the path first
	if err := utils.CheckIsFile(configPath); err != nil {
		return nil, err
	}

	flagData.ConfigPath = configPath

	return flagData, nil
}

func registerGlobalFlags(fset *flag.FlagSet) {
	flag.VisitAll(func(f *flag.Flag) {
		fset.Var(f.Value, f.Name, f.Usage)
	})
}

func fetch(args []string) types.Action {
	flag := flag.NewFlagSet("fetch", flag.ExitOnError)
	registerGlobalFlags(flag)
	flag.Parse(args)
	return utils.FETCH
}

func install(args []string) types.Action {
	flag := flag.NewFlagSet("install", flag.ExitOnError)
	registerGlobalFlags(flag)
	flag.Parse(args)
	return utils.INSTALL
}

func build(args []string) types.Action {
	flag := flag.NewFlagSet("build", flag.ExitOnError)
	registerGlobalFlags(flag)
	flag.Parse(args)
	return utils.BUILD
}

func clean(args []string) types.Action {
	flag := flag.NewFlagSet("clean", flag.ExitOnError)
	registerGlobalFlags(flag)
	flag.Parse(args)
	return utils.CLEAN
}

func all(args []string) types.Action {
	flag := flag.NewFlagSet("all", flag.ExitOnError)
	registerGlobalFlags(flag)
	flag.Parse(args)
	return utils.REPO_ALL
}

func serve(args []string) (types.Action, string) {
	flag := flag.NewFlagSet("serve", flag.ExitOnError)
	var profile = flag.String("p", "", "profile to use")
	registerGlobalFlags(flag)
	flag.Parse(args)
	if *profile == "" {
		log.Fatal("missing profile flag in serve command")
	}
	return utils.PROXY_SERVER, *profile
}
