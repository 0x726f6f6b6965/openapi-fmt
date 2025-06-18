package config_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/0x726f6f6b6965/openapi-fmt/config"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "test-config")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test case 1: Valid config file
	validYAML := `
input:
  path: "input/data"
  format: "json"
output:
  path: "output/data"
  format: "yaml"
rm-exts:
  enable: true
  excludes:
    - ".txt"
    - ".log"
sp:
  enable: false
  paths:
    - "input/data/sp1"
    - "input/data/sp2"
`
	validConfigPath := filepath.Join(tmpDir, "valid_config.yaml")
	if err := os.WriteFile(validConfigPath, []byte(validYAML), 0644); err != nil {
		t.Fatalf("Failed to write valid config file: %v", err)
	}

	expectedConfig := &config.Config{
		Input: config.InputConfig{
			Path:   "input/data",
			Format: "json",
		},
		Output: config.OutputConfig{
			Path:   "output/data",
			Format: "yaml",
		},
		RmExts: config.RmExtsConfig{
			Enable:   true,
			Excludes: []string{".txt", ".log"},
		},
		Sp: config.SpConfig{
			Enable: false,
			Paths:  []string{"input/data/sp1", "input/data/sp2"},
		},
	}

	cfg, err := config.LoadConfig(validConfigPath)
	if err != nil {
		t.Errorf("Expected no error for valid config, got %v", err)
	}
	if !reflect.DeepEqual(cfg, expectedConfig) {
		t.Errorf("Loaded config does not match expected config.\nGot: %+v\nExpected: %+v", cfg, expectedConfig)
	}

	// Test case 2: File not found
	_, err = config.LoadConfig(filepath.Join(tmpDir, "non_existent_config.yaml"))
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}

	// Test case 3: Invalid YAML format
	invalidYAML := `
input:
  path: "input/data"
  format: "json"
output
  path: "output/data"
  format: "yaml"
`
	invalidConfigPath := filepath.Join(tmpDir, "invalid_config.yaml")
	if err := os.WriteFile(invalidConfigPath, []byte(invalidYAML), 0644); err != nil {
		t.Fatalf("Failed to write invalid config file: %v", err)
	}

	_, err = config.LoadConfig(invalidConfigPath)
	if err == nil {
		t.Error("Expected error for invalid YAML format, got nil")
	}
}
