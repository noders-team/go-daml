package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/noders-team/go-daml/internal/codegen"
	"github.com/spf13/cobra"
)

var (
	dar    string
	output string
	debug  bool
	pkg    string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "godaml --dar <path> --output <dir> --go_package <name> [--debug]",
		Short: "Go DAML codegen tool",
		Long: `A command-line interface tool for generating Go code from DAML (.dar) files.

This tool extracts DAML definitions from .dar archives and generates corresponding Go structs and types.`,
		Example: `  godaml --dar ./test.dar --output ./generated --go_package main
  godaml --dar /path/to/contracts.dar --output ./src/daml --go_package contracts --debug`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if dar == "" {
				return fmt.Errorf("--dar parameter is required")
			}
			if output == "" {
				return fmt.Errorf("--output parameter is required")
			}
			if pkg == "" {
				return fmt.Errorf("--go_package parameter is required")
			}

			return runCodeGen(dar, output, pkg, debug)
		},
	}

	rootCmd.Flags().StringVar(&dar, "dar", "", "path to the DAR file (required)")
	rootCmd.Flags().StringVar(&output, "output", "", "output directory where generated Go files will be saved (required)")
	rootCmd.Flags().StringVar(&pkg, "go_package", "", "Go package name for generated code (required)")
	rootCmd.Flags().BoolVar(&debug, "debug", false, "enable debug logging")

	rootCmd.MarkFlagRequired("dar")
	rootCmd.MarkFlagRequired("output")
	rootCmd.MarkFlagRequired("go_package")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runCodeGen(dar, outputDir, pkgFile string, debugMode bool) error {
	if debugMode {
		log.Info().Msg("debug mode enabled")
	}

	unzippedPath, err := codegen.UnzipDar(dar, nil)
	if err != nil {
		return fmt.Errorf("failed to unzip dar file '%s': %w", dar, err)
	}
	defer os.RemoveAll(unzippedPath) // Clean up temporary files

	manifest, err := codegen.GetManifest(unzippedPath)
	if err != nil {
		return fmt.Errorf("failed to get manifest from '%s': %w", unzippedPath, err)
	}

	dalfFullPath := filepath.Join(unzippedPath, manifest.MainDalf)
	dalfContent, err := os.ReadFile(dalfFullPath)
	if err != nil {
		return fmt.Errorf("failed to read dalf file '%s': %w", dalfFullPath, err)
	}

	pkg, err := codegen.GetAST(dalfContent, manifest)
	if err != nil {
		return fmt.Errorf("failed to generate AST: %w", err)
	}

	res, err := codegen.Bind(pkgFile, pkg.PackageID, manifest.SdkVersion, pkg.Structs)
	if err != nil {
		return fmt.Errorf("failed to generate Go code: %w", err)
	}

	err = os.MkdirAll(outputDir, 0o755)
	if err != nil {
		return fmt.Errorf("failed to create output directory '%s': %w", outputDir, err)
	}

	fileName := filepath.Join(outputDir, strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(manifest.Name), ".", "_"), "-", "_")+".go")
	err = os.WriteFile(fileName, []byte(res), 0o644)
	if err != nil {
		return fmt.Errorf("failed to save generated file '%s': %w", fileName, err)
	}

	log.Info().Msgf("successfully generated Go code: %s", fileName)
	return nil
}
