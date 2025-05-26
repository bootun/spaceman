package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Agents Agents `yaml:"agents"`
}

type Agents struct {
	Spaceman Agent `yaml:"spaceman"`
}

type Agent struct {
	ModelID string `yaml:"model_id"`
	BaseURL string `yaml:"base_url"`
	Token   string `yaml:"token"`
}

func LoadConfig(filePath string) (*Config, error) {
	if filePath == "" {
		filePath = "config.yml"
	}

	// If not an absolute path, resolve to absolute
	if !filepath.IsAbs(filePath) {
		// Get the current working directory
		cwd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get current working directory: %w", err)
		}
		filePath = filepath.Join(cwd, filePath)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}
