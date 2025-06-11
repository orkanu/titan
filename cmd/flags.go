package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type Flags struct {
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
	// String that contains the configured configuration path
	var configPath string
	// Subcommands https://gobyexample.com/command-line-subcommands
	flag.StringVar(&configPath, "c", "./titan.yaml", "path to config file")

	// Actually parse the flags
	flag.Parse()

	// Validate the path first
	if err := ValidateConfigPath(configPath); err != nil {
		return &Flags{}, err
	}

	// Return the configuration path
	return &Flags{ConfigPath: configPath}, nil
}
