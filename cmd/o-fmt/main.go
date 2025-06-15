package main

import (
	"fmt"
	"os"

	"github.com/0x726f6f6b6965/openapi-fmt/utils"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const (
	InputFileFlag       = "input"
	InputFileShortFlag  = "i"
	OutputFileFlag      = "output"
	OutputFileShortFlag = "o"
	ExcludesFlag        = "excludes"
	ExcludesShortFlag   = "e"
	PathsFlag           = "paths"
	PathsShortFlag      = "p"

	CmdRemoveExtensions = "rm-exts"
	CmdSplitByPath      = "sp"
)

func main() {
	// This is a command-line tool for formatting OpenAPI documents.
	// It is intended to be run with the `o-fmt` command.
	// The tool will read an OpenAPI document from a specified input file,
	// remove specified extensions, and write the formatted document to an output file.
	// This is write by cobra
	var rootCmd = &cobra.Command{
		Use:   "o-fmt",
		Short: "This is a command-line tool for formatting OpenAPI documents.",
	}

	var (
		input    string
		output   string
		excludes []string
	)

	var removeCmd = &cobra.Command{
		Use:   CmdRemoveExtensions,
		Short: "remove extensions from OpenAPI document",
		Run:   RunRemove,
	}
	removeCmd.Flags().StringVarP(&input, InputFileFlag, InputFileShortFlag, "", "the input OpenAPI file path (must be provided)")
	removeCmd.Flags().StringVarP(&output, OutputFileFlag, OutputFileShortFlag, "", "the output OpenAPI file path (must be provided)")
	removeCmd.Flags().StringSliceVarP(&excludes, ExcludesFlag, ExcludesShortFlag, []string{}, "the fields to exclude from the OpenAPI document (comma-separated)")

	removeCmd.MarkFlagRequired(InputFileFlag)
	removeCmd.MarkFlagRequired(OutputFileFlag)

	var splitCmd = &cobra.Command{
		Use:   CmdSplitByPath,
		Short: "split OpenAPI document by path",
		Run:   RunSplit,
	}
	splitCmd.Flags().StringVarP(&input, InputFileFlag, InputFileShortFlag, "", "the input OpenAPI file path (must be provided)")
	splitCmd.Flags().StringVarP(&output, OutputFileFlag, OutputFileShortFlag, "", "the output OpenAPI file path (must be provided)")
	splitCmd.Flags().StringSliceVarP(&excludes, PathsFlag, PathsShortFlag, []string{}, "the paths to split the OpenAPI document by (comma-separated)")

	splitCmd.MarkFlagRequired(InputFileFlag)
	splitCmd.MarkFlagRequired(OutputFileFlag)
	splitCmd.MarkFlagRequired(PathsFlag)

	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(splitCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func RunRemove(cmd *cobra.Command, args []string) {
	input, err := cmd.Flags().GetString(InputFileFlag)
	if err != nil {
		cmd.PrintErrf("Error getting input file flag: %v\n", err)
		return
	}
	output, err := cmd.Flags().GetString(OutputFileFlag)
	if err != nil {
		cmd.PrintErrf("Error getting output file flag: %v\n", err)
		return
	}
	excludes, err := cmd.Flags().GetStringSlice(ExcludesFlag)
	if err != nil {
		cmd.PrintErrf("Error getting excludes flag: %v\n", err)
		return
	}
	// get input file path
	f, err := os.ReadFile(input)
	if err != nil {
		cmd.PrintErrf("Error reading input file: %s \nerror: %v\n", input, err)
		return
	}
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(f)
	if err != nil {
		cmd.PrintErrf("Error loading OpenAPI document: %v\n", err)
		return
	}
	// get exclude
	keep := map[string]struct{}{}
	for _, exclude := range excludes {
		if exclude == "" {
			continue
		}
		// add the exclude to the keep map
		keep[exclude] = struct{}{}
	}

	// remove extensions
	utils.RemoveExtensions(doc, keep)
	// marshal the document to YAML
	out, err := doc.MarshalYAML()
	if err != nil {
		cmd.PrintErrf("Error marshalling OpenAPI document: %v\n", err)
		return
	}
	data, err := yaml.Marshal(out)
	if err != nil {
		cmd.PrintErrf("Error marshalling OpenAPI document to YAML: %v\n", err)
		return
	}
	// write the formatted document to the output file
	if err := os.WriteFile(output, data, 0644); err != nil {
		cmd.PrintErrf("Error writing output file: %v\n", err)
		return
	}
	cmd.Printf("Formatted OpenAPI document written to %s\n", output)
}

func RunSplit(cmd *cobra.Command, args []string) {
	input, err := cmd.Flags().GetString(InputFileFlag)
	if err != nil {
		cmd.PrintErrf("Error getting input file flag: %v\n", err)
		return
	}
	output, err := cmd.Flags().GetString(OutputFileFlag)
	if err != nil {
		cmd.PrintErrf("Error getting output file flag: %v\n", err)
		return
	}
	// get input file path
	f, err := os.ReadFile(input)
	if err != nil {
		cmd.PrintErrf("Error reading input file: %s \nerror: %v\n", input, err)
		return
	}

	paths, err := cmd.Flags().GetStringSlice(PathsFlag)
	if err != nil {
		cmd.PrintErrf("Error getting paths flag: %v\n", err)
		return
	}
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(f)
	if err != nil {
		cmd.PrintErrf("Error loading OpenAPI document: %v\n", err)
		return
	}
	targets := map[string]struct{}{}
	for _, path := range paths {
		if path == "" {
			continue
		}
		targets[path] = struct{}{}
	}
	// split by path
	splitDoc, err := utils.SplitByPath(doc, targets)
	if err != nil {
		cmd.PrintErrf("Error splitting OpenAPI document by path: %v\n", err)
		return
	}
	// marshal the document to YAML
	out, err := splitDoc.MarshalYAML()
	if err != nil {
		cmd.PrintErrf("Error marshalling OpenAPI document: %v\n", err)
		return
	}
	data, err := yaml.Marshal(out)
	if err != nil {
		cmd.PrintErrf("Error marshalling OpenAPI document to YAML: %v\n", err)
		return
	}
	// write the formatted document to the output file
	if err := os.WriteFile(output, data, 0644); err != nil {
		cmd.PrintErrf("Error writing output file: %v\n", err)
		return
	}
	cmd.Printf("Formatted OpenAPI document written to %s\n", output)
}
