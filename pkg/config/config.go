package config

import (
	"os"
	"titan/internal/container"
	"titan/pkg/types"

	"gopkg.in/yaml.v3"
)

// NewConfig returns a new decoded Config struct
func NewConfig(container *container.Container) error {
	// Create config structure
	config := &types.Config{}

	// Open config file
	file, err := os.Open(container.ConfigData.ConfigFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return err
	}

	container.ConfigData.Config = config
	return nil
}
