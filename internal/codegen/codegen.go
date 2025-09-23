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

	daml "github.com/digital-asset/dazl-client/v8/go/api/com/daml/daml_lf_1_17"
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

type codeGenAst struct {
	Package *daml.Package
}

func NewCodegenAst(pkg *daml.Package) *codeGenAst {
	return &codeGenAst{Package: pkg}
}

func (c *codeGenAst) getName(id int32) string {
	idx := c.Package.InternedDottedNames[id].SegmentsInternedStr
	if len(idx) > 1 {
		log.Info().Msgf("more then 1")
	}
	return c.Package.InternedStrings[idx[len(idx)-1]]
}

func (c *codeGenAst) getDataTypes(module *daml.Module) (map[string]*tmplStruct, error) {
	structs := make(map[string]*tmplStruct, 0)
	for _, dataType := range module.GetDataTypes() {
		if !dataType.Serializable {
			continue
		}

		name := c.getName(dataType.GetNameInternedDname())
		tmplStruct := tmplStruct{
			Name: name,
		}

		switch v := dataType.DataCons.(type) {
		case *daml.DefDataType_Record:
			for _, field := range v.Record.Fields {
				fieldExtracted, typeExtracted, err := c.extractField(field)
				if err != nil {
					return nil, err
				}
				tmplStruct.Fields = append(tmplStruct.Fields, &tmplField{
					Name:    fieldExtracted,
					Type:    typeExtracted,
					RawType: field.String(),
				})
			}
		case *daml.DefDataType_Variant:
			for _, field := range v.Variant.Fields {
				fieldExtracted, typeExtracted, err := c.extractField(field)
				if err != nil {
					return nil, err
				}
				tmplStruct.Fields = append(tmplStruct.Fields, &tmplField{
					Name:    fieldExtracted,
					Type:    typeExtracted,
					RawType: field.String(),
				})
				log.Info().Msgf("variant constructor: %s, type: %s", fieldExtracted, typeExtracted)
			}
		case *daml.DefDataType_Enum:
			for _, constructorIdx := range v.Enum.ConstructorsInternedStr {
				constructorName := c.getName(constructorIdx)
				tmplStruct.Fields = append(tmplStruct.Fields, &tmplField{
					Name: constructorName,
					Type: "enum",
				})
				log.Info().Msgf("enum constructor: %s", constructorName)
			}
		case *daml.DefDataType_Interface:
			log.Warn().Msgf("interface not supported %s", v.Interface.String())
		default:
			log.Warn().Msgf("unknown data cons type: %T", v)
		}
		structs[name] = &tmplStruct
	}

	return structs, nil
}

func (c *codeGenAst) extractField(field *daml.FieldWithType) (string, string, error) {
	if field == nil {
		return "", "", fmt.Errorf("field is nil")
	}

	fieldName := c.getName(field.GetFieldInternedStr())
	if field.Type == nil {
		return fieldName, "", fmt.Errorf("field type is nil")
	}

	//	*Type_Var_
	//	*Type_Con_
	//	*Type_Builtin_
	//	*Type_Forall_
	//	*Type_Struct_
	//	*Type_Nat
	//	*Type_Syn_
	//	*Type_Interned
	var fieldType string
	switch v := field.Type.Sum.(type) {
	case *daml.Type_Interned:
		prim := c.Package.InternedTypes[v.Interned]
		if prim != nil {
			fieldType = prim.String()
			break
		}

		isConType := c.Package.InternedTypes[v.Interned].GetCon()
		if isConType != nil {
			tyconName := c.getName(isConType.Tycon.GetNameInternedDname())
			fieldType = tyconName
		} else {
			fieldType = "complex_interned_type"
		}
	case *daml.Type_Con_:
		if v.Con.Tycon != nil {
			switch {
			case v.Con.Tycon.GetNameInternedDname() != 0:
				fieldType = c.getName(v.Con.Tycon.GetNameInternedDname())
			default:
				fieldType = "unknown_con_type"
			}
		} else {
			fieldType = "con_without_tycon"
		}
	case *daml.Type_Var_:
		switch {
		case v.Var.GetVarInternedStr() != 0:
			fieldType = c.getName(v.Var.GetVarInternedStr())
		default:
			fieldType = "unnamed_var"
		}
	case *daml.Type_Forall_:
		fieldType = fmt.Sprintf("forall[%d_vars]", len(v.Forall.Vars))
	case *daml.Type_Struct_:
		fieldType = fmt.Sprintf("struct[%d_fields]", len(v.Struct.Fields))
	case *daml.Type_Nat:
		fieldType = fmt.Sprintf("nat_%d", v.Nat)
	case *daml.Type_Syn_:
		if v.Syn.Tysyn != nil {
			switch {
			case v.Syn.Tysyn.GetNameInternedDname() != 0:
				fieldType = fmt.Sprintf("syn_%s", c.getName(v.Syn.Tysyn.GetNameInternedDname()))
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

func GetAST(payload []byte, manifest *Manifest) (*Package, error) {
	structs := make(map[string]*tmplStruct)

	var archive daml.Archive
	err := proto.Unmarshal(payload, &archive)
	if err != nil {
		return nil, err
	}

	var payloadMapped daml.ArchivePayload
	err = proto.Unmarshal(archive.Payload, &payloadMapped)
	if err != nil {
		return nil, err
	}

	damlLf1 := payloadMapped.GetDamlLf_1()
	codeGen := NewCodegenAst(damlLf1)
	for _, module := range damlLf1.Modules {
		if len(damlLf1.InternedStrings) == 0 {
			continue
		}

		idx := damlLf1.InternedDottedNames[module.GetNameInternedDname()].SegmentsInternedStr
		moduleName := damlLf1.InternedStrings[idx[len(idx)-1]]
		log.Info().Msgf("processing module %s", moduleName)
		for _, dataType := range module.GetDataTypes() {
			if !dataType.Serializable {
				log.Warn().Msgf("skipping non-serializable data type in module %s", moduleName)
				continue
			}

			dt, err := codeGen.getDataTypes(module)
			if err != nil {
				return nil, err
			}
			for key, val := range dt {
				structs[key] = val
			}
		}

		for _, itrfc := range module.Interfaces {
			log.Info().Msgf("interface: %+v", itrfc)
		}

		/*
			for _, tmplt := range module.Templates {
				switch v := tmplt.Tycon.(type) {
				case *daml.DefTemplate_TyconDname:
					log.Info().Msgf("v: %+v", v)
				case *daml.DefTemplate_TyconInternedDname:
					idx := damlLf1.InternedDottedNames[v.TyconInternedDname].SegmentsInternedStr
					tmplName := damlLf1.InternedStrings[idx[len(idx)-1]]
					log.Info().Msgf("template name: %+s", tmplName)
				}
			}*/
	}

	return &Package{
		Structs: structs,
	}, nil
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
	case strings.HasPrefix(damlType, "prim:CONTRACT_ID"):
		return "CONTRACT_ID"
	case strings.HasPrefix(damlType, "prim:GENMAP"):
		return "GENMAP"
	case strings.HasPrefix(damlType, "prim:TEXTMAP"):
		return "TEXTMAP"
	case strings.HasPrefix(damlType, "prim:BIGNUMERIC"):
		return "BIGNUMERIC"
	case strings.HasPrefix(damlType, "prim:ROUNDING_MODE"):
		return "ROUNDING_MODE"
	case strings.HasPrefix(damlType, "prim:ANY"):
		return "ANY"
	case strings.HasPrefix(damlType, "prim:TYPE_REP"):
		return "TYPE_REP"
	case strings.HasPrefix(damlType, "prim:ARROW"):
		return "ARROW"
	case strings.HasPrefix(damlType, "prim:UPDATE"):
		return "UPDATE"
	case strings.HasPrefix(damlType, "prim:SCENARIO"):
		return "SCENARIO"
	case damlType == "enum":
		return "string"
	default:
		log.Warn().Msgf("unknown daml type %s", damlType)
		return damlType
	}
}
