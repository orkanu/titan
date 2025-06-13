package flags

import (
	"flag"
	"fmt"
	"log"
	"os"
	"titan/internal/utils"
)

type Flags struct {
	Command    utils.Command
	ConfigPath string
}

func checkIsFile(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

// ValidateConfigPath just makes sure, that the path provided is a file,
// that can be read
func ValidateConfigPath(path string) error {
	err := checkIsFile(path)
	if err != nil {
		log.Printf("%v\n", err)
		log.Println("Trying default config location ~/.config/titan/titan.yaml")
		err = checkIsFile("~/.config/titan/titan.yaml")
		if err != nil {
			return err
		}
	}

	return nil
}

// ParseFlags will create and parse the CLI flags
// For now it returns the path to config file to be used
func ParseFlags() (*Flags, error) {
	flagsData := &Flags{}
	// String that contains the configured configuration path
	var configPath string
	flag.StringVar(&configPath, "c", "./titan.yaml", "path to config file")

	// (fetch) subcommand
	fetchCmd := flag.NewFlagSet("fetch", flag.ExitOnError)
	fetchCmd.StringVar(&configPath, "c", "./titan.yaml", "path to config file")
	// (install) subcommand
	installCmd := flag.NewFlagSet("install", flag.ExitOnError)
	installCmd.StringVar(&configPath, "c", "./titan.yaml", "path to config file")
	// (build) subcommand
	buildCmd := flag.NewFlagSet("build", flag.ExitOnError)
	buildCmd.StringVar(&configPath, "c", "./titan.yaml", "path to config file")
	// (clean) subcommand
	cleanCmd := flag.NewFlagSet("clean", flag.ExitOnError)
	cleanCmd.StringVar(&configPath, "c", "./titan.yaml", "path to config file")
	// (all) subcommand
	allCmd := flag.NewFlagSet("all", flag.ExitOnError)
	allCmd.StringVar(&configPath, "c", "./titan.yaml", "path to config file")

	if len(os.Args) > 1 {
		// Parse flags based on command
		switch os.Args[1] {
		case "fetch":
			fetchCmd.Parse(os.Args[2:])
			flagsData.Command = utils.FETCH
		case "install":
			installCmd.Parse(os.Args[2:])
			flagsData.Command = utils.INSTALL
		case "build":
			buildCmd.Parse(os.Args[2:])
			flagsData.Command = utils.BUILD
		case "clean":
			cleanCmd.Parse(os.Args[2:])
			flagsData.Command = utils.CLEAN
		case "all":
			allCmd.Parse(os.Args[2:])
			flagsData.Command = utils.ALL
		default:
			flag.Usage()
			// TODO how to unclude all flag set command options?
			os.Exit(0)
		}
	} else {
		flag.Usage()
		// TODO how to unclude all flag set command options?
		os.Exit(0)
	}

	// Validate the path first
	if err := ValidateConfigPath(configPath); err != nil {
		return &Flags{}, err
	}
	flagsData.ConfigPath = configPath

	// Return the configuration path
	return flagsData, nil
}
