package utils

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTempConfigFile(t *testing.T, content string) string {
	t.Helper()
	tmpFile, err := os.CreateTemp("", "config-*.yml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	_, err = tmpFile.Write([]byte(content))
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()
	return tmpFile.Name()
}

func TestLoadConfig_ValidFile(t *testing.T) {
	content := `port: 8080`
	fileName := createTempConfigFile(t, content)
	defer os.Remove(fileName)

	renamedFilePath := filepath.Join(filepath.Dir(fileName), "config.yml")
	err := os.Rename(fileName, renamedFilePath)
	assert.NoError(t, err)
	defer os.Remove(renamedFilePath)

	conf, err := LoadConfig(renamedFilePath)
	assert.NoError(t, err)
	assert.NotNil(t, conf)
	assert.Equal(t, 8080, conf.Port)
}

func TestLoadConfig_InvalidFileName(t *testing.T) {
	content := `port: 8080`
	fileName := createTempConfigFile(t, content)
	defer os.Remove(fileName)

	_, err := LoadConfig(fileName)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidFileName))
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	content := `port: not-a-number`
	fileName := createTempConfigFile(t, content)
	defer os.Remove(fileName)

	renamedFilePath := filepath.Join(filepath.Dir(fileName), "config.yml")
	err := os.Rename(fileName, renamedFilePath)
	assert.NoError(t, err)
	defer os.Remove(renamedFilePath)

	_, err = LoadConfig(renamedFilePath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal yaml")
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := LoadConfig("nonexistent/config.yml")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read file content")
}

func TestValidateFilePath(t *testing.T) {
	err := validateFilePath("config.yml")
	assert.NoError(t, err)

	err = validateFilePath("not_config.yml")
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidFileName, err)

	err = validateFilePath("subdir/config.yml")
	assert.NoError(t, err)
}
