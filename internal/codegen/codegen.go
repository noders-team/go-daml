package codegen

import (
	"archive/zip"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/proto"

	"github.com/digital-asset/dazl-client/v8/go/api/com/daml/daml_lf_1_17"
	"github.com/digital-asset/dazl-client/v8/go/api/com/daml/daml_lf_2_1"
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

	if output == nil {
		tmpDir := os.TempDir()
		output = &tmpDir
	}
	randomID, err := generateRandomID()
	if err != nil {
		return "", fmt.Errorf("failed to generate random ID: %w", err)
	}
	*output = filepath.Join(*output, randomID)

	os.MkdirAll(*output, 0o755)

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

		path := filepath.Join(*output, f.Name)

		if !strings.HasPrefix(path, filepath.Clean(*output)+string(os.PathSeparator)) {
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

	return *output, nil
}

func GetManifest(srcPath string) (*Manifest, error) {
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

	manifest := &Manifest{}
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

func GetAST(payload []byte, manifest *Manifest) (*Package, error) {
	structs := make(map[string]*tmplStruct)

	if strings.HasPrefix(manifest.SdkVersion, "1.") {
		var archive daml_lf_1_17.Archive
		err := proto.Unmarshal(payload, &archive)
		if err != nil {
			return nil, err
		}

		var payloadMapped daml_lf_1_17.ArchivePayload
		err = proto.Unmarshal(archive.Payload, &payloadMapped)
		if err != nil {
			return nil, err
		}

		damlLf1 := payloadMapped.GetDamlLf_1()
		for _, module := range damlLf1.Modules {
			for _, dataType := range module.DataTypes {
				name, err := extractDataTypeName(dataType.Name, damlLf1.InternedStrings)
				if err != nil {
					return nil, err
				}
				tmplStruct := tmplStruct{
					Name: name,
				}

				switch v := dataType.DataCons.(type) {
				case *daml_lf_1_17.DefDataType_Record:
					for _, field := range v.Record.Fields {
						fieldExtracted, typeExtracted, err := extractField(field, damlLf1.InternedStrings, damlLf1.InternedTypes)
						if err != nil {
							return nil, err
						}
						tmplStruct.Fields = append(tmplStruct.Fields, &tmplField{
							Name: fieldExtracted,
							Type: typeExtracted,
						})
					}
				case *daml_lf_1_17.DefDataType_Variant:
					for _, field := range v.Variant.Fields {
						fieldExtracted, typeExtracted, err := extractField(field, damlLf1.InternedStrings, damlLf1.InternedTypes)
						if err != nil {
							return nil, err
						}
						tmplStruct.Fields = append(tmplStruct.Fields, &tmplField{
							Name: fieldExtracted,
							Type: typeExtracted,
						})
						log.Info().Msgf("variant constructor: %s, type: %s", fieldExtracted, typeExtracted)
					}
				case *daml_lf_1_17.DefDataType_Enum:
					for _, constructorStr := range v.Enum.ConstructorsStr {
						tmplStruct.Fields = append(tmplStruct.Fields, &tmplField{
							Name: constructorStr,
							Type: "enum",
						})
						log.Info().Msgf("enum constructor: %s", constructorStr)
					}
					for _, constructorIdx := range v.Enum.ConstructorsInternedStr {
						if int(constructorIdx) >= len(damlLf1.InternedStrings) {
							return nil, fmt.Errorf("interned enum constructor index out of bounds: %d", constructorIdx)
						}
						constructorName := damlLf1.InternedStrings[constructorIdx]
						tmplStruct.Fields = append(tmplStruct.Fields, &tmplField{
							Name: constructorName,
							Type: "enum",
						})
						log.Info().Msgf("enum constructor: %s", constructorName)
					}
				default:
					log.Warn().Msgf("unknown data cons type: %T", v)
				}
				structs[name] = &tmplStruct
			}
		}
	} else {
		var archive daml_lf_2_1.Archive
		err := proto.Unmarshal(payload, &archive)
		if err != nil {
			return nil, err
		}

		var payloadMapped daml_lf_2_1.ArchivePayload
		err = proto.Unmarshal(archive.Payload, &payloadMapped)
		if err != nil {
			return nil, err
		}

		damlLf1 := payloadMapped.GetDamlLf_2()
		for _, module := range damlLf1.Modules {
			for _, dataType := range module.DataTypes {
				log.Info().Msgf("data type: %+v", dataType)
			}
			for _, template := range module.Templates {
				log.Info().Msgf("template: %+v", template)
			}
			for _, inter := range module.Interfaces {
				log.Info().Msgf("interface: %+v", inter)
			}
		}
	}

	return &Package{
		Structs: structs,
	}, nil
}

func extractDataTypeName(dataTypeName interface{}, internedStrings []string) (string, error) {
	switch v := dataTypeName.(type) {
	case *daml_lf_1_17.DefDataType_NameInternedDname:
		if int(v.NameInternedDname) >= len(internedStrings) {
			return "", fmt.Errorf("interned string index out of bounds: %d", v.NameInternedDname)
		}
		return internedStrings[v.NameInternedDname], nil
	case *daml_lf_1_17.DefDataType_NameDname:
		return v.NameDname.String(), nil
	default:
		return "", fmt.Errorf("unknown name type: %T", dataTypeName)
	}
}

func extractField(field *daml_lf_1_17.FieldWithType, internedStrings []string, internedTypes []*daml_lf_1_17.Type) (string, string, error) {
	if field == nil {
		return "", "", fmt.Errorf("field is nil")
	}

	var fieldName string
	switch v := field.Field.(type) {
	case *daml_lf_1_17.FieldWithType_FieldStr:
		fieldName = v.FieldStr
	case *daml_lf_1_17.FieldWithType_FieldInternedStr:
		if int(v.FieldInternedStr) >= len(internedStrings) {
			return "", "", fmt.Errorf("interned string index out of bounds: %d", v.FieldInternedStr)
		}
		fieldName = internedStrings[v.FieldInternedStr]
	default:
		return "", "", fmt.Errorf("unknown field type: %T", field.Field)
	}

	var fieldType string
	if field.Type == nil {
		return fieldName, "", fmt.Errorf("field type is nil")
	}

	//	*Type_Var_
	//	*Type_Con_
	//	*Type_Prim_
	//	*Type_Forall_
	//	*Type_Struct_
	//	*Type_Nat
	//	*Type_Syn_
	//	*Type_Interned
	switch v := field.Type.Sum.(type) {
	case *daml_lf_1_17.Type_Interned:
		if int(v.Interned) >= len(internedTypes) {
			return fieldName, "", fmt.Errorf("interned type index out of bounds: %d", v.Interned)
		}
		if prim := internedTypes[v.Interned].GetPrim(); prim != nil {
			fieldType = prim.String()
		} else {
			fieldType = "complex_interned_type"
		}
	case *daml_lf_1_17.Type_Prim_:
		fieldType = v.Prim.String()
	case *daml_lf_1_17.Type_Con_:
		if v.Con.Tycon != nil {
			switch {
			case v.Con.Tycon.GetNameInternedDname() != 0:
				if int(v.Con.Tycon.GetNameInternedDname()) >= len(internedStrings) {
					return fieldName, "", fmt.Errorf("interned tycon index out of bounds: %d", v.Con.Tycon.GetNameInternedDname())
				}
				fieldType = internedStrings[v.Con.Tycon.GetNameInternedDname()]
			case v.Con.Tycon.GetNameDname() != nil:
				fieldType = v.Con.Tycon.GetNameDname().String()
			default:
				fieldType = "unknown_con_type"
			}
		} else {
			fieldType = "con_without_tycon"
		}
	case *daml_lf_1_17.Type_Var_:
		switch {
		case v.Var.GetVarInternedStr() != 0:
			if int(v.Var.GetVarInternedStr()) >= len(internedStrings) {
				return fieldName, "", fmt.Errorf("interned var index out of bounds: %d", v.Var.GetVarInternedStr())
			}
			fieldType = internedStrings[v.Var.GetVarInternedStr()]
		case v.Var.GetVarStr() != "":
			fieldType = v.Var.GetVarStr()
		default:
			fieldType = "unnamed_var"
		}
	case *daml_lf_1_17.Type_Forall_:
		fieldType = fmt.Sprintf("forall[%d_vars]", len(v.Forall.Vars))
	case *daml_lf_1_17.Type_Struct_:
		fieldType = fmt.Sprintf("struct[%d_fields]", len(v.Struct.Fields))
	case *daml_lf_1_17.Type_Nat:
		fieldType = fmt.Sprintf("nat_%d", v.Nat)
	case *daml_lf_1_17.Type_Syn_:
		if v.Syn.Tysyn != nil {
			switch {
			case v.Syn.Tysyn.GetNameInternedDname() != 0:
				if int(v.Syn.Tysyn.GetNameInternedDname()) >= len(internedStrings) {
					return fieldName, "", fmt.Errorf("interned tysyn index out of bounds: %d", v.Syn.Tysyn.GetNameInternedDname())
				}
				fieldType = fmt.Sprintf("syn_%s", internedStrings[v.Syn.Tysyn.GetNameInternedDname()])
			case v.Syn.Tysyn.GetNameDname() != nil:
				fieldType = fmt.Sprintf("syn_%s", v.Syn.Tysyn.GetNameDname().String())
			default:
				fieldType = "syn_unknown"
			}
		} else {
			fieldType = "syn_without_name"
		}
	default:
		return fieldName, "", fmt.Errorf("unsupported type sum: %T", field.Type.Sum)
	}

	return fieldName, normalizeDAMLType(fieldType), nil
}

func normalizeDAMLType(damlType string) string {
	switch {
	case strings.HasPrefix(damlType, "prim:PARTY"):
		return "PARTY"
	case strings.HasPrefix(damlType, "prim:TEXT"):
		return "TEXT"
	case strings.HasPrefix(damlType, "prim:INT64"):
		return "INT64"
	case strings.HasPrefix(damlType, "prim:BOOL"):
		return "BOOL"
	case strings.HasPrefix(damlType, "prim:DECIMAL"):
		return "DECIMAL"
	case strings.HasPrefix(damlType, "prim:NUMERIC"):
		return "NUMERIC"
	case strings.HasPrefix(damlType, "prim:DATE"):
		return "DATE"
	case strings.HasPrefix(damlType, "prim:TIMESTAMP"):
		return "TIMESTAMP"
	case strings.HasPrefix(damlType, "prim:UNIT"):
		return "UNIT"
	case strings.HasPrefix(damlType, "prim:LIST"):
		return "LIST"
	case strings.HasPrefix(damlType, "prim:MAP"):
		return "MAP"
	case strings.HasPrefix(damlType, "prim:OPTIONAL"):
		return "OPTIONAL"
	case damlType == "enum":
		return "string"
	default:
		return "interface{}"
	}
}
