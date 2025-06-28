package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Input  InputConfig  `yaml:"input"`
	Output OutputConfig `yaml:"output"`
	RmExts RmExtsConfig `yaml:"rm-exts"`
	Sp     SpConfig     `yaml:"sp"`
}

type InputConfig struct {
	Path   string `yaml:"path"`
	Format string `yaml:"format"`
}

type OutputConfig struct {
	Path   string `yaml:"path"`
	Format string `yaml:"format"`
}

type RmExtsConfig struct {
	Enable   bool     `yaml:"enable"`
	Excludes []string `yaml:"excludes"`
}

type SpConfig struct {
	Enable    bool       `yaml:"enable"`
	Endpoints []Endpoint `yaml:"endpoints"`
}

type Endpoint struct {
	Path    string   `yaml:"path"`
	Methods []string `yaml:"methods"`
}

// LoadConfig loads and parses the config.yaml file
// The function takes the file path as input and returns a Config struct or an error.
func LoadConfig[V any](path string) (*V, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config V
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config data: %w", err)
	}

	return &config, nil
}
