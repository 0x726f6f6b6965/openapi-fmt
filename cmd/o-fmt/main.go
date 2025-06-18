package main

import (
	"fmt"
	"os"

	"github.com/0x726f6f6b6965/openapi-fmt/config"
	"github.com/0x726f6f6b6965/openapi-fmt/utils"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const (
	InputFileFlag         = "input"
	InputFileShortFlag    = "i"
	OutputFileFlag        = "output"
	OutputFileShortFlag   = "o"
	OutputFormatFlag      = "output-format"
	OutputFormatShortFlag = "f"
	ExcludesFlag          = "excludes"
	ExcludesShortFlag     = "e"
	PathsFlag             = "paths"
	PathsShortFlag        = "p"
	ConfigFlag            = "config"
	ConfigShortFlag       = "c"
	RmExtsFlag            = "remove-exts"
	RmExtsShortFlag       = "r"
)

var (
	configFile    string
	inputPath     string
	outputPath    string
	outputFmt     string
	excludesSlice []string
	pathsSlice    []string
	rmEnable      bool
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "o-fmt",
		Short: "This is a command-line tool for formatting OpenAPI documents.",
		RunE:  RunE,
	}
	rootCmd.PersistentFlags().StringVarP(&configFile, ConfigFlag, ConfigShortFlag, "", "path to the config file (e.g. config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&inputPath, InputFileFlag, InputFileShortFlag, "", "path to the input OpenAPI file")
	rootCmd.PersistentFlags().StringVarP(&outputPath, OutputFileFlag, OutputFileShortFlag, "", "path to the output OpenAPI file")
	rootCmd.PersistentFlags().StringVarP(&outputFmt, OutputFormatFlag, OutputFormatShortFlag, "yaml", "format of the output file (yaml or json)")
	rootCmd.PersistentFlags().StringSliceVarP(&excludesSlice, ExcludesFlag, ExcludesShortFlag, []string{}, "extensions to exclude from the output file")
	rootCmd.PersistentFlags().StringSliceVarP(&pathsSlice, PathsFlag, PathsShortFlag, []string{}, "paths to split the OpenAPI document")
	rootCmd.PersistentFlags().BoolVarP(&rmEnable, RmExtsFlag, RmExtsShortFlag, false, "enable removing extensions from the OpenAPI document")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func RunE(cmd *cobra.Command, args []string) error {

	var cfg *config.Config
	if configFile != "" {
		var err error
		cfg, err = config.LoadConfig(configFile)
		if err != nil {
			return fmt.Errorf("Error loading config file '%s': %w", configFile, err)
		}

		if cfg.Input.Path != "" {
			inputPath = cfg.Input.Path
		}
		if cfg.Output.Path != "" {
			outputPath = cfg.Output.Path
		}
		if cfg.Output.Format != "" {
			outputFmt = cfg.Output.Format
		}
		if cfg.RmExts.Enable {
			rmEnable = cfg.RmExts.Enable
			if len(cfg.RmExts.Excludes) > 0 {
				excludesSlice = cfg.RmExts.Excludes
			}
		}
		if cfg.Sp.Enable && len(cfg.Sp.Paths) > 0 {
			pathsSlice = cfg.Sp.Paths
		}
	}

	if inputPath == "" {
		return fmt.Errorf("Error: input file path must be provided via flag or config file")
	}
	if outputPath == "" {
		return fmt.Errorf("Error: output file path must be provided via flag or config file")
	}
	if outputFmt != "yaml" && outputFmt != "json" {
		return fmt.Errorf("Error: output format must be either 'yaml' or 'json'")
	}
	if len(excludesSlice) != 0 {
		rmEnable = true
	}

	var (
		source *openapi3.T
		err    error
	)

	f, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("Error reading input file '%s': %w", inputPath, err)
	}
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(f)
	if err != nil {
		return fmt.Errorf("Error loading OpenAPI document from '%s': %w", inputPath, err)
	}

	if len(pathsSlice) > 0 {
		targets := map[string]struct{}{}
		for _, path := range pathsSlice {
			if path == "" {
				continue
			}
			targets[path] = struct{}{}
		}

		source, err = utils.SplitByPath(doc, targets)
		if err != nil {
			return fmt.Errorf("Error splitting OpenAPI document by path: %w", err)
		}
	} else {
		source = doc
	}

	if rmEnable {
		// remove extensions
		keep := map[string]struct{}{}
		for _, exclude := range excludesSlice {
			if exclude == "" {
				continue
			}
			keep[exclude] = struct{}{}
		}
		utils.RemoveExtensions(source, keep)
	}
	var data []byte
	if outputFmt == "yaml" {
		out, err := source.MarshalYAML()
		if err != nil {
			return fmt.Errorf("Error marshalling OpenAPI document: %w", err)
		}
		data, err = yaml.Marshal(out)
		if err != nil {
			return fmt.Errorf("Error marshalling OpenAPI document to YAML: %w", err)
		}
	} else if outputFmt == "json" {
		data, err = source.MarshalJSON()
		if err != nil {
			return fmt.Errorf("Error marshalling OpenAPI document to JSON: %w", err)
		}
	}
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("Error writing output file '%s': %w", outputPath, err)
	}

	return nil // Success
}
