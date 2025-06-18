package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

// runTestMain is a helper function to execute RunE with specified configurations.
// It resets global flags before each run.
func runTestMain(t *testing.T, configFilePathVal, inputPathVal, outputPathVal, outputFmtVal string, excludesVal []string, pathsVal []string, rmEnableVal bool) error {
	// Reset global flags to defaults or empty states
	configFile = ""
	inputPath = ""
	outputPath = ""
	outputFmt = "yaml" // Default value in main.go
	excludesSlice = nil
	pathsSlice = nil
	rmEnable = false

	// Set global flags based on test case parameters
	if configFilePathVal != "" {
		configFile = configFilePathVal
	}
	if inputPathVal != "" {
		inputPath = inputPathVal
	}
	if outputPathVal != "" {
		outputPath = outputPathVal
	}
	if outputFmtVal != "" {
		outputFmt = outputFmtVal
	}
	if excludesVal != nil {
		excludesSlice = make([]string, len(excludesVal))
		copy(excludesSlice, excludesVal) // Make a defensive copy
	}
	if pathsVal != nil {
		pathsSlice = make([]string, len(pathsVal))
		copy(pathsSlice, pathsVal) // Make a defensive copy
	}
	if rmEnableVal {
		rmEnable = rmEnableVal
	}

	// Diagnostic prints
	// t.Logf("Debug: In runTestMain, before RunE:")
	// t.Logf("Debug: configFile = %q", configFile)
	// t.Logf("Debug: inputPath = %q", inputPath)
	// t.Logf("Debug: outputPath = %q", outputPath)
	// t.Logf("Debug: outputFmt = %q", outputFmt)
	// t.Logf("Debug: excludesSlice = %v", excludesSlice)
	// t.Logf("Debug: pathsSlice = %v", pathsSlice)
	// t.Logf("Debug: rmEnable = %t", rmEnable)

	// Specific check for TestRunE_RemoveExtensions to ensure helper is working as expected for that case
	if filepath.Base(inputPathVal) == "input_ext.yaml" {
		// This assertion is in the helper, so it runs for every test that might use "input_ext.yaml".
		// However, we're interested in the state specifically for TestRunE_RemoveExtensions's call.
		expectedExcludes := []string{"x-go-type"} // Corrected to x-go-type
		currentExcludes := excludesSlice // The global that was just set
		assert.Equal(t, expectedExcludes, currentExcludes, "HELPER CHECK: excludesSlice not set correctly before RunE for input_ext.yaml")

		// For TestRunE_RemoveExtensions, rmEnableVal is now true.
		// The global rmEnable will also be true.
		// The logic within RunE might further modify rmEnable based on excludesSlice,
		// but at this stage (before RunE), it reflects rmEnableVal.
		if filepath.Base(inputPathVal) == "input_ext.yaml" { // Only for this specific test
			assert.True(t, rmEnableVal, "HELPER CHECK: rmEnableVal for input_ext.yaml should be true as per test setup")
			assert.True(t, rmEnable, "HELPER CHECK: global rmEnable should be true before RunE for input_ext.yaml if rmEnableVal was true")
		}
	}


	// We pass nil for cmd and empty slice for args as RunE doesn't use them directly,
	// but relies on the global flag variables.
	return RunE(nil, []string{})
}

const simpleOpenAPIYAML = `
openapi: 3.0.0
info:
  title: Simple Test API
  version: 1.0.0
paths:
  /hello:
    get:
      summary: Says hello
      responses:
        '200':
          description: OK
`

func TestRunE_SuccessYAML(t *testing.T) {
	tempDir := t.TempDir()
	inputFilePath := filepath.Join(tempDir, "input.yaml")
	outputFilePath := filepath.Join(tempDir, "output.yaml")

	err := os.WriteFile(inputFilePath, []byte(simpleOpenAPIYAML), 0644)
	assert.NoError(t, err, "Failed to write temp input file")

	runErr := runTestMain(t,
		"",            // configFile
		inputFilePath, // inputPath
		outputFilePath, // outputPath
		"yaml",        // outputFmt
		nil,           // excludesSlice
		nil,           // pathsSlice
		false,         // rmEnable
	)
	assert.NoError(t, runErr, "RunE returned an error")

	outputData, err := os.ReadFile(outputFilePath)
	assert.NoError(t, err, "Failed to read output file")
	assert.True(t, len(outputData) > 0, "Output file is empty")

	// Basic check: does it look like YAML?
	var yamlData map[string]interface{}
	err = yaml.Unmarshal(outputData, &yamlData)
	assert.NoError(t, err, "Output file is not valid YAML")
	assert.Equal(t, "Simple Test API", yamlData["info"].(map[string]interface{})["title"], "Unexpected title in output YAML")
}

func TestRunE_SuccessJSON(t *testing.T) {
	tempDir := t.TempDir()
	inputFilePath := filepath.Join(tempDir, "input.yaml") // Input can still be YAML
	outputFilePath := filepath.Join(tempDir, "output.json")

	err := os.WriteFile(inputFilePath, []byte(simpleOpenAPIYAML), 0644)
	assert.NoError(t, err, "Failed to write temp input file")

	runErr := runTestMain(t,
		"",             // configFile
		inputFilePath,  // inputPath
		outputFilePath, // outputPath
		"json",         // outputFmt
		nil,            // excludesSlice
		nil,            // pathsSlice
		false,          // rmEnable
	)
	assert.NoError(t, runErr, "RunE returned an error")

	outputData, err := os.ReadFile(outputFilePath)
	assert.NoError(t, err, "Failed to read output file")
	assert.True(t, len(outputData) > 0, "Output file is empty")

	// Basic check: does it look like JSON?
	// For simplicity, we'll just check if it starts with { and ends with }
	// A more robust check would be json.Unmarshal
	assert.True(t, strings.HasPrefix(string(outputData), "{"), "Output is not valid JSON (missing '{')")
	assert.True(t, strings.HasSuffix(strings.TrimSpace(string(outputData)), "}"), "Output is not valid JSON (missing '}')")
}

const simpleOpenAPIForConfigTest = `
openapi: 3.0.0
info:
  title: Config Test API
  version: 1.0.0
paths:
  /test:
    get:
      summary: Test endpoint for config
      x-remove-me: "should be gone if rmEnable"
      x-go-type: "should stay if excluded" # Changed x-keep-me to x-go-type
      responses:
        '200':
          description: OK
`
const simpleOpenAPIForPathSplit = `
openapi: 3.0.0
info:
  title: Path Split Test API
  version: 1.0.0
x-another-ext: "should be removed by config rm-exts" # Doc level
paths:
  /api/v1/users:
    get:
      summary: Users endpoint
      x-config-keep: "config should keep this"
      x-flag-keep: "flag might try to keep this"
      responses:
        '200':
          description: OK
  /api/v1/orders:
    get:
      summary: Orders endpoint
      x-another-ext: "on order path, should be removed by config rm-exts"
      responses:
        '200':
          description: OK
`

func TestRunE_ConfigTakesPrecedence(t *testing.T) {
	tempDir := t.TempDir()
	cfgFilePath := filepath.Join(tempDir, "config.yaml")
	cfgInputPath := filepath.Join(tempDir, "input_from_cfg.yaml")
	cfgOutputPath := filepath.Join(tempDir, "output_from_cfg.json") // Note: .json for format test

	// Create dummy input file that config will point to
	err := os.WriteFile(cfgInputPath, []byte(simpleOpenAPIForConfigTest), 0644)
	assert.NoError(t, err)

	// Create config file
	configContent := `
input:
  path: ` + cfgInputPath + `
output:
  path: ` + cfgOutputPath + `
  format: json
`
	err = os.WriteFile(cfgFilePath, []byte(configContent), 0644)
	assert.NoError(t, err)

	// These flag values should be ignored in favor of config values
	flagInputPath := filepath.Join(tempDir, "input_from_flag.yaml")
	flagOutputPath := filepath.Join(tempDir, "output_from_flag.yaml")

	runErr := runTestMain(t,
		cfgFilePath,    // configFile
		flagInputPath,  // inputPath (should be overridden by config)
		flagOutputPath, // outputPath (should be overridden by config)
		"yaml",         // outputFmt (should be overridden by config)
		nil,
		nil,
		false,
	)
	assert.NoError(t, runErr, "RunE returned an error")

	// Check that output was created at the path specified in config, not flag
	_, err = os.Stat(flagOutputPath)
	assert.True(t, os.IsNotExist(err), "Output file was created at flag path, but should be at config path")

	outputData, err := os.ReadFile(cfgOutputPath)
	assert.NoError(t, err, "Failed to read output file from config path")
	assert.True(t, len(outputData) > 0, "Output file is empty")
	assert.True(t, strings.HasPrefix(string(outputData), "{"), "Output is not JSON as specified in config")
}

func TestRunE_FlagsOverrideMissingConfigFields(t *testing.T) {
	tempDir := t.TempDir()
	cfgFilePath := filepath.Join(tempDir, "config.yaml")
	cfgInputPath := filepath.Join(tempDir, "input_from_cfg.yaml")
	cfgOutputPath := filepath.Join(tempDir, "output_from_cfg.yaml") // Config will specify this

	err := os.WriteFile(cfgInputPath, []byte(simpleOpenAPIForConfigTest), 0644)
	assert.NoError(t, err)

	// Config file missing output.format
	configContent := `
input:
  path: ` + cfgInputPath + `
output:
  path: ` + cfgOutputPath + `
`
	err = os.WriteFile(cfgFilePath, []byte(configContent), 0644)
	assert.NoError(t, err)

	runErr := runTestMain(t,
		cfgFilePath,
		"",             // inputPath (use from config)
		"",             // outputPath (use from config)
		"json",         // outputFmt (this flag should be used as config doesn't specify it)
		nil,
		nil,
		false,
	)
	assert.NoError(t, runErr, "RunE returned an error")

	outputData, err := os.ReadFile(cfgOutputPath)
	assert.NoError(t, err, "Failed to read output file from config path")
	assert.True(t, strings.HasPrefix(string(outputData), "{"), "Output is not JSON as specified by flag")
}

func TestRunE_ConfigAndFlagsCombined(t *testing.T) {
	tempDir := t.TempDir()
	cfgFilePath := filepath.Join(tempDir, "config.yaml")
	cfgInputPath := filepath.Join(tempDir, "input_from_cfg.yaml")
	flagOutputPath := filepath.Join(tempDir, "output_from_flag.yaml") // Output path from flag

	// Create input file (will be used by config)
	err := os.WriteFile(cfgInputPath, []byte(simpleOpenAPIForPathSplit), 0644) // Use path split version for this test
	assert.NoError(t, err)

	// Config specifies input, rm-exts, and some excludes
	// Flags will specify output path and paths to split
	configContent := `
input:
  path: ` + cfgInputPath + `
rm-exts:
  enable: true
  excludes:
    - "x-config-keep" # Config will try to keep this
`
	err = os.WriteFile(cfgFilePath, []byte(configContent), 0644)
	assert.NoError(t, err)

	runErr := runTestMain(t,
		cfgFilePath,
		"",             // inputPath (from config)
		flagOutputPath, // outputPath (from flag)
		"yaml",         // outputFmt (default, not in config or flag)
		[]string{"x-flag-keep"}, // excludesSlice from flag (will be overridden by config's excludes)
		[]string{"/api/v1/users"}, // pathsSlice (from flag)
		false,          // rmEnable (config rm-exts.enable = true will take precedence)
	)
	assert.NoError(t, runErr, "RunE returned an error")

	outputData, err := os.ReadFile(flagOutputPath)
	assert.NoError(t, err, "Failed to read output file from flag path")
	outputStr := string(outputData)

	// Check path splitting
	var yamlData map[string]interface{}
	err = yaml.Unmarshal(outputData, &yamlData)
	assert.NoError(t, err, "Output is not valid YAML")
	pathsMap := yamlData["paths"].(map[string]interface{})
	assert.Contains(t, pathsMap, "/api/v1/users", "Path /api/v1/users missing")
	assert.NotContains(t, pathsMap, "/api/v1/orders", "Path /api/v1/orders should be absent")

	// Check extension removal based on config
	// Note: RunE's current logic is that if config.RmExts.Excludes is present, it overrides flag's excludesSlice.
	assert.Contains(t, outputStr, "x-config-keep", "Extension 'x-config-keep' should be present due to config")
	assert.NotContains(t, outputStr, "x-flag-keep", "Extension 'x-flag-keep' should be removed (config excludes take precedence)")
	assert.NotContains(t, outputStr, "x-another-ext", "Extension 'x-another-ext' should be removed")
}

func TestRunE_RemoveExtensions(t *testing.T) {
	tempDir := t.TempDir()
	inputFilePath := filepath.Join(tempDir, "input_ext.yaml")
	outputFilePath := filepath.Join(tempDir, "output_ext.yaml")

	// Using simpleOpenAPIForConfigTest as it has x-remove-me and x-keep-me
	err := os.WriteFile(inputFilePath, []byte(simpleOpenAPIForConfigTest), 0644)
	assert.NoError(t, err)

	runErr := runTestMain(t,
		"",             // configFile
		inputFilePath,  // inputPath
		outputFilePath, // outputPath
		"yaml",         // outputFmt
		[]string{"x-go-type"}, // Changed excludesSlice to x-go-type
		nil,            // pathsSlice
		true,           // rmEnable is explicitly true from flag
	)
	assert.NoError(t, runErr, "RunE returned an error")

	outputData, err := os.ReadFile(outputFilePath)
	assert.NoError(t, err, "Failed to read output file")

	outputStr := string(outputData)
	assert.Contains(t, outputStr, "x-go-type", "Excluded extension x-go-type should be present")
	assert.NotContains(t, outputStr, "x-remove-me", "Extension x-remove-me should be removed")
}

func TestRunE_SplitPath(t *testing.T) {
	tempDir := t.TempDir()
	inputFilePath := filepath.Join(tempDir, "input_paths.yaml")
	outputFilePath := filepath.Join(tempDir, "output_paths.yaml")

	err := os.WriteFile(inputFilePath, []byte(simpleOpenAPIForPathSplit), 0644)
	assert.NoError(t, err)

	runErr := runTestMain(t,
		"",             // configFile
		inputFilePath,  // inputPath
		outputFilePath, // outputPath
		"yaml",         // outputFmt
		nil,            // excludesSlice
		[]string{"/api/v1/orders"}, // pathsSlice
		false,          // rmEnable
	)
	assert.NoError(t, runErr, "RunE returned an error")

	outputData, err := os.ReadFile(outputFilePath)
	assert.NoError(t, err, "Failed to read output file")

	var yamlData map[string]interface{}
	err = yaml.Unmarshal(outputData, &yamlData)
	assert.NoError(t, err, "Output is not valid YAML")

	pathsMap := yamlData["paths"].(map[string]interface{})
	assert.Contains(t, pathsMap, "/api/v1/orders", "Path /api/v1/orders missing")
	assert.NotContains(t, pathsMap, "/api/v1/users", "Path /api/v1/users should be absent")
}

// Error Condition Test Cases

func TestRunE_ErrorNoInputPath(t *testing.T) {
	tempDir := t.TempDir()
	outputFilePath := filepath.Join(tempDir, "output.yaml")

	err := runTestMain(t, "", "", outputFilePath, "yaml", nil, nil, false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "input file path must be provided")
}

func TestRunE_ErrorNoOutputPath(t *testing.T) {
	tempDir := t.TempDir()
	inputFilePath := filepath.Join(tempDir, "input.yaml")
	err := os.WriteFile(inputFilePath, []byte(simpleOpenAPIYAML), 0644)
	assert.NoError(t, err)

	runErr := runTestMain(t, "", inputFilePath, "", "yaml", nil, nil, false)
	assert.Error(t, runErr)
	assert.Contains(t, runErr.Error(), "output file path must be provided")
}

func TestRunE_ErrorInvalidOutputFormat(t *testing.T) {
	tempDir := t.TempDir()
	inputFilePath := filepath.Join(tempDir, "input.yaml")
	outputFilePath := filepath.Join(tempDir, "output.yaml")
	err := os.WriteFile(inputFilePath, []byte(simpleOpenAPIYAML), 0644)
	assert.NoError(t, err)

	runErr := runTestMain(t, "", inputFilePath, outputFilePath, "xml", nil, nil, false) // xml is invalid
	assert.Error(t, runErr)
	assert.Contains(t, runErr.Error(), "output format must be either 'yaml' or 'json'")
}

func TestRunE_ErrorInputFileNonExistent(t *testing.T) {
	tempDir := t.TempDir()
	inputFilePath := filepath.Join(tempDir, "nonexistent_input.yaml")
	outputFilePath := filepath.Join(tempDir, "output.yaml")

	err := runTestMain(t, "", inputFilePath, outputFilePath, "yaml", nil, nil, false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Error reading input file")
}

func TestRunE_ErrorConfigFileNonExistent(t *testing.T) {
	tempDir := t.TempDir()
	// Don't need input/output for this, as config load is early
	err := runTestMain(t, filepath.Join(tempDir, "nonexistent_config.yaml"), "", "", "yaml", nil, nil, false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Error loading config file")
}

const malformedYAML = `
openapi: 3.0.0
info:
  title: Malformed API
  version: 1.0.0
paths: /hello # This should be a map, not a string
`

func TestRunE_ErrorInvalidInputOpenAPI(t *testing.T) {
	tempDir := t.TempDir()
	inputFilePath := filepath.Join(tempDir, "malformed_input.yaml")
	outputFilePath := filepath.Join(tempDir, "output.yaml")

	err := os.WriteFile(inputFilePath, []byte(malformedYAML), 0644)
	assert.NoError(t, err)

	runErr := runTestMain(t, "", inputFilePath, outputFilePath, "yaml", nil, nil, false)
	assert.Error(t, runErr)
	assert.Contains(t, runErr.Error(), "Error loading OpenAPI document")
}

func TestRunE_ErrorSplitPathNotFound(t *testing.T) {
	tempDir := t.TempDir()
	inputFilePath := filepath.Join(tempDir, "input_paths.yaml")
	outputFilePath := filepath.Join(tempDir, "output_paths.yaml")

	err := os.WriteFile(inputFilePath, []byte(simpleOpenAPIForPathSplit), 0644) // Contains /api/v1/users and /api/v1/orders
	assert.NoError(t, err)

	runErr := runTestMain(t,
		"",             // configFile
		inputFilePath,  // inputPath
		outputFilePath, // outputPath
		"yaml",         // outputFmt
		nil,            // excludesSlice
		[]string{"/api/v1/nonexistentpath"}, // pathsSlice with a path not in the doc
		false,          // rmEnable
	)
	assert.Error(t, runErr)
	assert.Contains(t, runErr.Error(), "Error splitting OpenAPI document by path") // This is generic, check for utils.ErrOpenAPIPathNotFound if possible
	// To check for specific error type, you might need to adjust RunE or use errors.Is with a sentinel error from utils.
	// For now, checking the message is a good start.
}
