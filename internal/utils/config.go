// Package utils provides utilities for reading application configurations and other additional features.
package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Errors for validation.
var (
	ErrInvalidFileName = errors.New("invalid file name: only 'config.yml' is allowed")
	ErrNotInCurrentDir = errors.New("file must be in the current directory")
)

// Config yaml file.
type Config struct {
	Port int `yaml:"port"`
}

// LoadConfig read config from yml file.
func LoadConfig(filepath string) (*Config, error) {
	if err := validateFilePath(filepath); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(filepath) // #nosec G304
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	return &config, nil
}

func validateFilePath(filePath string) error {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("invalid file path: %w", err)
	}

	if filepath.Base(absPath) != "config.yml" {
		return ErrInvalidFileName
	}

	return nil
}
