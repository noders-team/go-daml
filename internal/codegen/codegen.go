package codegen

import (
	"archive/zip"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/noders-team/go-daml/internal/codegen/astgen"
	"github.com/noders-team/go-daml/internal/codegen/model"
	"github.com/rs/zerolog/log"
)

func generateRandomID() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 15)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return string(b), nil
}

func UnzipDar(src string, output *string) (string, error) {
	r, err := zip.OpenReader(src)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := r.Close(); err != nil {
			log.Err(err).Msgf("failed to close zip file")
		}
	}()

	var out string
	if output == nil {
		tmpDir := os.TempDir()
		out = tmpDir
	} else {
		out = *output
	}

	randomID, err := generateRandomID()
	if err != nil {
		return "", fmt.Errorf("failed to generate random ID: %w", err)
	}
	out = filepath.Join(out, randomID)

	os.MkdirAll(out, 0o755)

	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(out, f.Name)

		if !strings.HasPrefix(path, filepath.Clean(out)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, 0o755)
		} else {
			os.MkdirAll(filepath.Dir(path), 0o755)
			outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
			if err != nil {
				return err
			}
			defer func() {
				if err := outFile.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(outFile, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return "", err
		}
	}

	return out, nil
}

func GetManifest(srcPath string) (*model.Manifest, error) {
	manifestPath := strings.Join([]string{srcPath, "META-INF", "MANIFEST.MF"}, "/")
	file, err := os.Open(manifestPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	content := strings.ReplaceAll(string(b), "\n ", "")

	manifest := &model.Manifest{}
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Manifest-Version:") {
			manifest.Version = strings.TrimSpace(strings.TrimPrefix(line, "Manifest-Version:"))
		} else if strings.HasPrefix(line, "Created-By:") {
			manifest.CreatedBy = strings.TrimSpace(strings.TrimPrefix(line, "Created-By:"))
		} else if strings.HasPrefix(line, "Name:") {
			manifest.Name = strings.TrimSpace(strings.TrimPrefix(line, "Name:"))
		} else if strings.HasPrefix(line, "Sdk-Version:") {
			manifest.SdkVersion = strings.TrimSpace(strings.TrimPrefix(line, "Sdk-Version:"))
		} else if strings.HasPrefix(line, "Main-Dalf:") {
			manifest.MainDalf = strings.TrimSpace(strings.TrimPrefix(line, "Main-Dalf:"))
		} else if strings.HasPrefix(line, "Dalfs:") {
			dalfsStr := strings.TrimSpace(strings.TrimPrefix(line, "Dalfs:"))
			if dalfsStr != "" {
				manifest.Dalfs = strings.Split(dalfsStr, ", ")
				for i, dalf := range manifest.Dalfs {
					manifest.Dalfs[i] = strings.TrimSpace(dalf)
				}
			}
		} else if strings.HasPrefix(line, "Format:") {
			manifest.Format = strings.TrimSpace(strings.TrimPrefix(line, "Format:"))
		} else if strings.HasPrefix(line, "Encryption:") {
			manifest.Encryption = strings.TrimSpace(strings.TrimPrefix(line, "Encryption:"))
		}
	}

	if manifest.MainDalf == "" {
		return nil, fmt.Errorf("main-dalf not found in manifest")
	}

	return manifest, nil
}

func CodegenDalfs(dalfToProcess []string, unzippedPath string, pkgFile string, dalfManifest *model.Manifest) (map[string]string, error) {
	ifcByModule := make(map[string]model.InterfaceMap)
	result := make(map[string]string)

	for _, dalf := range dalfToProcess {
		dalfFullPath := filepath.Join(unzippedPath, dalf)
		dalfContent, err := os.ReadFile(dalfFullPath)
		if err != nil {
			log.Warn().Err(err).Msgf("failed to read dalf '%s': %s", dalf, err)
			continue
		}

		interfaces, err := GetInterfaces(dalfContent, dalfManifest)
		if err != nil {
			log.Warn().Err(err).Msgf("failed to extract interfaces from dalf: %s", dalf)
			continue
		}

		for key, val := range interfaces {
			equalNames := 0
			for _, ifcName := range ifcByModule {
				for ifcKey := range ifcName {
					res, found := strings.CutPrefix(ifcKey, key)
					_, atoiErr := strconv.Atoi(res)
					if found && (res == "" || atoiErr == nil) {
						equalNames++
					}
				}
			}
			if equalNames > 0 {
				equalNames++
				val.Name = fmt.Sprintf("%s%d", key, equalNames)
			}

			m, ok := ifcByModule[val.ModuleName]
			if !ok {
				m = make(model.InterfaceMap)
				ifcByModule[val.ModuleName] = m
			}
			m[val.Name] = val
		}
	}

	allStructNames := make(map[string]int)

	for _, dalf := range dalfToProcess {
		dalfFullPath := filepath.Join(unzippedPath, dalf)
		dalfContent, err := os.ReadFile(dalfFullPath)
		if err != nil {
			log.Warn().Err(err).Msgf("failed to read dalf '%s': %s", dalf, err)
			continue
		}

		pkg, err := GetAST(dalfContent, dalfManifest, ifcByModule)
		if err != nil {
			return nil, fmt.Errorf("failed to generate AST: %w", err)
		}

		currentModules := make(map[string]bool)
		for _, structDef := range pkg.Structs {
			if structDef.ModuleName != "" {
				currentModules[structDef.ModuleName] = true
			}
		}

		log.Info().Msgf("adding interfaces for dalf %s from modules: %v", dalf, currentModules)
		for moduleName := range currentModules {
			if ifcMap, exists := ifcByModule[moduleName]; exists {
				for key, val := range ifcMap {
					log.Debug().Msgf("adding interface %s from module %s to output", key, moduleName)
					pkg.Structs[key] = val
				}
			}
		}

		renamedStructs := make(map[string]*model.TmplStruct)

		for structName, structDef := range pkg.Structs {
			if structDef.IsInterface {
				continue
			}

			equalNames := 0
			for existingName := range allStructNames {
				res, found := strings.CutPrefix(existingName, structName)
				_, atoiErr := strconv.Atoi(res)
				if found && (res == "" || atoiErr == nil) {
					equalNames++
				}
			}

			if equalNames > 0 {
				equalNames++
				originalName := structName
				newName := fmt.Sprintf("%s%d", structName, equalNames)
				structDef.Name = newName

				delete(pkg.Structs, originalName)
				pkg.Structs[newName] = structDef
				renamedStructs[originalName] = structDef
				allStructNames[newName] = equalNames
			} else {
				allStructNames[structName] = 0
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

		code, err := Bind(pkgFile, pkg.Name, dalfManifest.SdkVersion, pkg.Structs, dalf == dalfManifest.MainDalf)
		if err != nil {
			return nil, fmt.Errorf("failed to generate Go code: %w", err)
		}

		result[dalf] = code
	}
	return result, nil
}

func GetInterfaces(payload []byte, manifest *model.Manifest) (map[string]*model.TmplStruct, error) {
	var version string
	if strings.HasPrefix(manifest.SdkVersion, astgen.V3) {
		version = astgen.V3
	} else if strings.HasPrefix(manifest.SdkVersion, astgen.V2) || strings.HasPrefix(manifest.SdkVersion, astgen.V1) {
		version = astgen.V2
	} else {
		return nil, fmt.Errorf("unsupported sdk version %s", manifest.SdkVersion)
	}

	gen, err := astgen.GetAstGenFromVersion(payload, version)
	if err != nil {
		return nil, err
	}

	return gen.GetInterfaces()
}

func GetAST(payload []byte, manifest *model.Manifest, ifcByModule map[string]model.InterfaceMap) (*model.Package, error) {
	var version string
	if strings.HasPrefix(manifest.SdkVersion, astgen.V3) {
		version = astgen.V3
	} else if strings.HasPrefix(manifest.SdkVersion, astgen.V2) || strings.HasPrefix(manifest.SdkVersion, astgen.V1) {
		version = astgen.V2
	} else {
		return nil, fmt.Errorf("unsupported sdk version %s", manifest.SdkVersion)
	}

	gen, err := astgen.GetAstGenFromVersion(payload, version)
	if err != nil {
		return nil, err
	}
	var structs map[string]*model.TmplStruct
	structs, err = gen.GetTemplateStructs(ifcByModule)
	if err != nil {
		return nil, err
	}

	packageID := getPackageID(manifest.MainDalf)
	if packageID == "" {
		return nil, fmt.Errorf("could not extract package ID from MainDalf: %s", manifest.MainDalf)
	}

	packageName := manifest.Name
	if packageName == "" {
		packageName = packageID
	}

	return &model.Package{
		Name:      packageName,
		PackageID: packageID,
		Structs:   structs,
	}, nil
}

func getPackageID(mainDalf string) string {
	parts := strings.Split(mainDalf, "/")
	filename := strings.TrimSuffix(parts[len(parts)-1], ".dalf")

	lastHyphen := strings.LastIndex(filename, "-")
	if lastHyphen != -1 && lastHyphen < len(filename)-1 {
		return filename[lastHyphen+1:]
	}

	return ""
}
