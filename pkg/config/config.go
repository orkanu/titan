package config

import (
	"os"
	"titan/pkg/types"

	"gopkg.in/yaml.v3"
)

// NewConfig returns a new decoded Config struct
func NewConfig(configFilePath string) (*types.Config, error) {
	// Create config structure
	config := &types.Config{}

	// Open config file
	file, err := os.Open(configFilePath)
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

	// container.ConfigData.Config = config
	return config, nil
}
