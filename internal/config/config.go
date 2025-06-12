package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Versions struct {
	// Node version to install/use in system via NVM
	Node string `yaml:"node"`
	// PNOM version to install
	PNPM string `yaml:"pnpm"`
}

// Config struct for titan
type Config struct {
	Versions Versions `yaml:"versions"`
	// Base path where the repositories are located
	BasePath string `yaml:"base_path"`

	// List of respositories
	Repositories map[string]string `yaml:"repositories"`
}

// NewConfig returns a new decoded Config struct
func NewConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
