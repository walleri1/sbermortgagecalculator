// Package utils provides utilities for reading application configurations and other additional features.
package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Errors for validation.
var (
	ErrInvalidFileName = errors.New("invalid file name: only 'config.yml' is allowed")
	ErrNotInCurrentDir = errors.New("file must be in the current directory")
)

// Config yaml file.
type Config struct {
	Server struct {
		Port         int           `yaml:"port"`
		ReadTimeout  time.Duration `yaml:"read_timeout"`
		WriteTimeout time.Duration `yaml:"write_timeout"`
		IdleTimeout  time.Duration `yaml:"idle_timeout"`
	} `yaml:"server"`

	Logging struct {
		Level string `yaml:"level"`
		File  string `yaml:"file"`
	} `yaml:"logging"`

	CORS struct {
		AllowedOrigins []string `yaml:"allowed_origins"`
		AllowedMethods []string `yaml:"allowed_methods"`
		AllowedHeaders []string `yaml:"allowed_headers"`
	} `yaml:"cors"`

	Cache struct {
		MaxSize int           `yaml:"max_size"`
		TTL     time.Duration `yaml:"ttl"`
	} `yaml:"cache"`

	Swagger struct {
		Enabled bool   `yaml:"enabled"`
		Path    string `yaml:"path"`
	} `yaml:"swagger"`
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

	// Устанавливаем значения по умолчанию, если они не указаны
	if config.Server.Port == 0 {
		config.Server.Port = 8080
	}
	if config.Server.ReadTimeout == 0 {
		config.Server.ReadTimeout = 10 * time.Second
	}
	if config.Server.WriteTimeout == 0 {
		config.Server.WriteTimeout = 10 * time.Second
	}
	if config.Server.IdleTimeout == 0 {
		config.Server.IdleTimeout = 20 * time.Second
	}
	if config.Logging.Level == "" {
		config.Logging.Level = "info"
	}
	if len(config.CORS.AllowedOrigins) == 0 {
		config.CORS.AllowedOrigins = []string{"*"}
	}
	if len(config.CORS.AllowedMethods) == 0 {
		config.CORS.AllowedMethods = []string{"GET", "POST", "OPTIONS"}
	}
	if len(config.CORS.AllowedHeaders) == 0 {
		config.CORS.AllowedHeaders = []string{"Content-Type", "Authorization"}
	}
	if config.Cache.MaxSize == 0 {
		config.Cache.MaxSize = 1000
	}
	if config.Cache.TTL == 0 {
		config.Cache.TTL = 24 * time.Hour
	}
	if config.Swagger.Path == "" {
		config.Swagger.Path = "/swagger"
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
