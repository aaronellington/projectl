package configuration

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// Errors
var (
	ErrMissingConfigFile = errors.New("missing config file")
	ErrInvalidConfigFile = errors.New("invalid config file")
)

// NewConfig creates a new config object with the defaults already set
func NewConfig(configFilePath string) (*Config, error) {
	config := &Config{}

	// Open the config file
	configFile, err := os.Open(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrMissingConfigFile, configFilePath)
	}
	defer configFile.Close()

	// Parse the config file
	decoder := json.NewDecoder(configFile)
	if err := decoder.Decode(config); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidConfigFile, configFilePath)
	}

	return config, nil
}

// Config of projectl
type Config struct {
	Gitignore    []string `json:"gitignore"`
	DistedFiles  []string `json:"disted_files"`
	DockerName   string   `json:"docker_name"`
	DockerTarget string   `json:"docker_target"`
	DockerPort   int      `json:"docker_port"`
}
