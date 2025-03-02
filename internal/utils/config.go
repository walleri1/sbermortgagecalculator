// The utils package provides utilities for reading application configurations and other additional features.
package utils

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// struct Config yaml file
type Config struct {
	Port int `yaml:"port"`
}

// LoadConfig read config from yml file
func LoadConfig(filepath string) (*Config, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("is not read file content: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("is not unmarshal yaml: %w", err)
	}

	return &config, nil
}
