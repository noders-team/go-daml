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

func processDalf(dalfRelPath, unzippedPath, pkgName, sdkVersion, outputDir string, isMainDalf bool, allInterfaces map[string]*model.TmplStruct) error {
	dalfFullPath := filepath.Join(unzippedPath, dalfRelPath)
	dalfContent, err := os.ReadFile(dalfFullPath)
	if err != nil {
		return fmt.Errorf("failed to read dalf file '%s': %w", dalfFullPath, err)
	}

	manifest := &model.Manifest{
		SdkVersion: sdkVersion,
		MainDalf:   dalfRelPath,
	}

	pkg, err := codegen.GetASTWithInterfaces(dalfContent, manifest, allInterfaces)
	if err != nil {
		return fmt.Errorf("failed to generate AST: %w", err)
	}

	code, err := codegen.Bind(pkgName, pkg.PackageID, sdkVersion, pkg.Structs, isMainDalf)
	if err != nil {
		return fmt.Errorf("failed to generate Go code: %w", err)
	}

	baseFileName := getFilenameFromDalf(dalfRelPath)
	outputFile := filepath.Join(outputDir, baseFileName+".go")

	if err := os.WriteFile(outputFile, []byte(code), 0o644); err != nil {
		return fmt.Errorf("failed to write file '%s': %w", outputFile, err)
	}

	log.Info().Msgf("successfully generated: %s", outputFile)
	return nil
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

	log.Info().Msg("first pass: collecting interfaces from all DALFs")
	allInterfaces := make(map[string]*model.TmplStruct)

	dalfFullPath := filepath.Join(unzippedPath, manifest.MainDalf)
	dalfContent, err := os.ReadFile(dalfFullPath)
	if err != nil {
		return fmt.Errorf("failed to read MainDalf '%s': %w", dalfFullPath, err)
	}

	dalfManifest := &model.Manifest{
		SdkVersion: manifest.SdkVersion,
		MainDalf:   manifest.MainDalf,
	}

	interfaces, err := codegen.GetInterfaces(dalfContent, dalfManifest)
	if err != nil {
		log.Warn().Err(err).Msgf("failed to extract interfaces from MainDalf: %s", manifest.MainDalf)
	} else {
		for key, val := range interfaces {
			allInterfaces[key] = val
		}
		log.Info().Msgf("collected %d interfaces from MainDalf", len(interfaces))
	}

	for _, dalf := range manifest.Dalfs {
		if dalf == manifest.MainDalf {
			continue
		}

		dalfLower := strings.ToLower(dalf)
		if strings.Contains(dalfLower, "prim") || strings.Contains(dalfLower, "stdlib") {
			continue
		}

		dalfFullPath := filepath.Join(unzippedPath, dalf)
		dalfContent, err := os.ReadFile(dalfFullPath)
		if err != nil {
			log.Warn().Err(err).Msgf("failed to read dalf '%s': %s", dalf, err)
			continue
		}

		dalfManifest := &model.Manifest{
			SdkVersion: manifest.SdkVersion,
			MainDalf:   dalf,
		}

		interfaces, err := codegen.GetInterfaces(dalfContent, dalfManifest)
		if err != nil {
			log.Warn().Err(err).Msgf("failed to extract interfaces from dalf: %s", dalf)
			continue
		}

		for key, val := range interfaces {
			allInterfaces[key] = val
		}
		log.Info().Msgf("collected %d interfaces from %s", len(interfaces), dalf)
	}

	log.Info().Msgf("total interfaces collected: %d", len(allInterfaces))

	log.Info().Msg("second pass: generating code for all DALFs")
	allStructNames := make(map[string]string)

	log.Info().Msgf("processing MainDalf: %s", manifest.MainDalf)
	err = processDalfWithConflictCheck(manifest.MainDalf, unzippedPath, pkgFile, manifest.SdkVersion, outputDir, true, allInterfaces, allStructNames)
	if err != nil {
		return fmt.Errorf("failed to process MainDalf: %w", err)
	}

	successCount := 1
	failedCount := 0
	skippedCount := 0

	for _, dalf := range manifest.Dalfs {
		if dalf == manifest.MainDalf {
			log.Debug().Msgf("skipping MainDalf (already processed): %s", dalf)
			skippedCount++
			continue
		}

		dalfLower := strings.ToLower(dalf)
		if strings.Contains(dalfLower, "prim") || strings.Contains(dalfLower, "stdlib") {
			log.Debug().Msgf("skipping dalf (prim/stdlib): %s", dalf)
			skippedCount++
			continue
		}

		log.Info().Msgf("processing dependency dalf: %s", dalf)
		err = processDalfWithConflictCheck(dalf, unzippedPath, pkgFile, manifest.SdkVersion, outputDir, false, allInterfaces, allStructNames)
		if err != nil {
			log.Error().Err(err).Msgf("failed to process dalf: %s", dalf)
			failedCount++
			continue
		}
		successCount++
	}

	log.Info().Msgf("code generation summary: %d succeeded, %d failed, %d skipped", successCount, failedCount, skippedCount)

	if successCount == 0 {
		return fmt.Errorf("all dalf files failed to process")
	}

	return nil
}

func getDALFPrefix(dalfRelPath string) string {
	parts := strings.Split(dalfRelPath, "/")
	var dalfFileName string
	if len(parts) > 1 {
		dalfFileName = parts[len(parts)-1]
	} else {
		dalfFileName = dalfRelPath
	}
	dalfFileName = strings.TrimSuffix(dalfFileName, ".dalf")
	dalfFileName = removePackageID(dalfFileName)
	parts = strings.FieldsFunc(dalfFileName, func(r rune) bool {
		return r == '.' || r == '_' || r == '-' || r == ' '
	})

	var result strings.Builder
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}
		result.WriteString(strings.ToUpper(part[:1]) + part[1:])
	}

	return result.String()
}

func processDalfWithConflictCheck(dalfRelPath, unzippedPath, pkgName, sdkVersion, outputDir string, isMainDalf bool, allInterfaces map[string]*model.TmplStruct, allStructNames map[string]string) error {
	dalfFullPath := filepath.Join(unzippedPath, dalfRelPath)
	dalfContent, err := os.ReadFile(dalfFullPath)
	if err != nil {
		return fmt.Errorf("failed to read dalf file '%s': %w", dalfFullPath, err)
	}

	manifest := &model.Manifest{
		SdkVersion: sdkVersion,
		MainDalf:   dalfRelPath,
	}

	pkg, err := codegen.GetASTWithInterfaces(dalfContent, manifest, allInterfaces)
	if err != nil {
		return fmt.Errorf("failed to generate AST: %w", err)
	}

	dalfPrefix := getDALFPrefix(dalfRelPath)
	renamedStructs := make(map[string]*model.TmplStruct)
	structsToProcess := make(map[string]*model.TmplStruct)

	for structName, structDef := range pkg.Structs {
		if !structDef.IsInterface {
			structsToProcess[structName] = structDef
		}
	}

	for structName, structDef := range structsToProcess {
		if existingDalf, exists := allStructNames[structName]; exists {
			newName := dalfPrefix + structName
			log.Info().Msgf("conflict detected: struct '%s' already exists in %s, renaming to %s", structName, existingDalf, newName)
			structDef.Name = newName
			renamedStructs[structName] = structDef
			allStructNames[newName] = dalfRelPath
			delete(pkg.Structs, structName)
			pkg.Structs[newName] = structDef
		} else {
			allStructNames[structName] = dalfRelPath
		}
	}

	for _, structDef := range pkg.Structs {
		for _, field := range structDef.Fields {
			if renamed, exists := renamedStructs[field.Type]; exists {
				field.Type = renamed.Name
			}
			trimmedType := strings.TrimPrefix(field.Type, "*")
			trimmedType = strings.TrimPrefix(trimmedType, "[]")
			if renamed, exists := renamedStructs[trimmedType]; exists {
				field.Type = strings.Replace(field.Type, trimmedType, renamed.Name, 1)
			}
		}

		for _, choice := range structDef.Choices {
			if renamed, exists := renamedStructs[choice.ArgType]; exists {
				choice.ArgType = renamed.Name
			}
			if renamed, exists := renamedStructs[choice.ReturnType]; exists {
				choice.ReturnType = renamed.Name
			}
		}
	}

	code, err := codegen.Bind(pkgName, pkg.PackageID, sdkVersion, pkg.Structs, isMainDalf)
	if err != nil {
		return fmt.Errorf("failed to generate Go code: %w", err)
	}

	baseFileName := getFilenameFromDalf(dalfRelPath)
	outputFile := filepath.Join(outputDir, baseFileName+".go")

	if err := os.WriteFile(outputFile, []byte(code), 0o644); err != nil {
		return fmt.Errorf("failed to write file '%s': %w", outputFile, err)
	}

	log.Info().Msgf("successfully generated: %s", outputFile)
	return nil
}
