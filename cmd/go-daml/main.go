package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/noders-team/go-daml/internal/codegen"
	"github.com/noders-team/go-daml/internal/codegen/model"
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

func removePackageID(filename string) string {
	lastHyphen := strings.LastIndex(filename, "-")
	if lastHyphen == -1 {
		return filename
	}

	potentialHash := filename[lastHyphen+1:]
	if len(potentialHash) == 64 {
		allHex := true
		for _, ch := range potentialHash {
			if !((ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')) {
				allHex = false
				break
			}
		}
		if allHex {
			return filename[:lastHyphen]
		}
	}

	return filename
}

func getFilenameFromDalf(dalfRelPath string) string {
	parts := strings.Split(dalfRelPath, "/")
	var baseFileName string
	if len(parts) > 1 {
		dalfFileName := parts[len(parts)-1]
		baseFileName = strings.TrimSuffix(dalfFileName, ".dalf")
	} else {
		baseFileName = strings.TrimSuffix(dalfRelPath, ".dalf")
	}

	baseFileName = removePackageID(baseFileName)
	sanitizedFileName := strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(baseFileName), ".", "_"), "-", "_")
	return sanitizedFileName
}

func runCodeGen(dar, outputDir, pkgFile string, debugMode bool) error {
	if debugMode {
		log.Info().Msg("debug mode enabled")
	}

	unzippedPath, err := codegen.UnzipDar(dar, nil)
	if err != nil {
		return fmt.Errorf("failed to unzip dar file '%s': %w", dar, err)
	}
	defer os.RemoveAll(unzippedPath)

	manifest, err := codegen.GetManifest(unzippedPath)
	if err != nil {
		return fmt.Errorf("failed to get manifest from '%s': %w", unzippedPath, err)
	}

	err = os.MkdirAll(outputDir, 0o755)
	if err != nil {
		return fmt.Errorf("failed to create output directory '%s': %w", outputDir, err)
	}

	dalfs := make([]string, 0)
	for _, dalf := range manifest.Dalfs {
		if dalf == manifest.MainDalf {
			continue
		}

		dalfLower := strings.ToLower(dalf)
		if strings.Contains(dalfLower, "prim") || strings.Contains(dalfLower, "stdlib") {
			continue
		}

		dalfs = append(dalfs, dalf)
	}

	dalfManifest := &model.Manifest{
		SdkVersion: manifest.SdkVersion,
		MainDalf:   manifest.MainDalf,
		Dalfs:      dalfs,
	}

	dalfToProcess := make([]string, 0)
	dalfToProcess = append(dalfToProcess, manifest.MainDalf)
	dalfToProcess = append(dalfToProcess, dalfs...)

	result, err := codegen.CodegenDalfs(dalfToProcess, unzippedPath, pkgFile, dalfManifest)
	if err != nil {
		return err
	}

	for dalf, code := range result {
		baseFileName := getFilenameFromDalf(dalf)
		outputFile := filepath.Join(outputDir, baseFileName+".go")

		if err := os.WriteFile(outputFile, []byte(code), 0o644); err != nil {
			return fmt.Errorf("failed to write file '%s': %w", outputFile, err)
		}

		log.Info().Msgf("successfully generated: %s", outputFile)
	}

	return nil
}
