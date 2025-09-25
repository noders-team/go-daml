package v2

import (
	"errors"
	"fmt"

	daml "github.com/digital-asset/dazl-client/v8/go/api/com/daml/daml_lf_1_17"
	"github.com/noders-team/go-daml/internal/codegen/model"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

type codeGenAst struct {
	payload []byte
}

func NewCodegenAst(payload []byte) *codeGenAst {
	return &codeGenAst{payload: payload}
}

func (c *codeGenAst) GetTemplateStructs() (string, map[string]*model.TmplStruct, error) {
	structs := make(map[string]*model.TmplStruct)

	var archive daml.Archive
	err := proto.Unmarshal(c.payload, &archive)
	if err != nil {
		return "", nil, err
	}

	var payloadMapped daml.ArchivePayload
	err = proto.Unmarshal(archive.Payload, &payloadMapped)
	if err != nil {
		return "", nil, err
	}

	damlLf1 := payloadMapped.GetDamlLf_1()
	if damlLf1 == nil {
		return "", nil, errors.New("unsupported daml version")
	}

	for _, module := range damlLf1.Modules {
		if len(damlLf1.InternedStrings) == 0 {
			continue
		}

		idx := damlLf1.InternedDottedNames[module.GetNameInternedDname()].SegmentsInternedStr
		moduleName := damlLf1.InternedStrings[idx[len(idx)-1]]
		log.Info().Msgf("processing module %s", moduleName)

		// Process templates first (template-centric approach)
		templates, err := c.getTemplates(damlLf1, module)
		if err != nil {
			return "", nil, err
		}
		for key, val := range templates {
			structs[key] = val
		}

		// Process interfaces
		interfaces, err := c.getInterfaces(damlLf1, module)
		if err != nil {
			return "", nil, err
		}
		for key, val := range interfaces {
			structs[key] = val
		}

		// Process remaining data types that aren't covered by templates/interfaces
		dataTypes, err := c.getDataTypes(damlLf1, module)
		if err != nil {
			return "", nil, err
		}
		for key, val := range dataTypes {
			// Only add if not already processed as part of templates/interfaces
			if _, exists := structs[key]; !exists {
				structs[key] = val
			}
		}
	}

	return archive.Hash, structs, nil
}

func (c *codeGenAst) getName(pkg *daml.Package, id int32) string {
	idx := pkg.InternedDottedNames[id].SegmentsInternedStr
	return pkg.InternedStrings[idx[len(idx)-1]]
}

func (c *codeGenAst) getTemplates(pkg *daml.Package, module *daml.Module) (map[string]*model.TmplStruct, error) {
	structs := make(map[string]*model.TmplStruct, 0)

	for _, template := range module.Templates {
		var templateName string

		switch v := template.Tycon.(type) {
		case *daml.DefTemplate_TyconDname:
			templateName = v.TyconDname.String()
		case *daml.DefTemplate_TyconInternedDname:
			templateName = c.getName(pkg, v.TyconInternedDname)
		default:
			log.Warn().Msgf("unknown template tycon type: %T", v)
			continue
		}

		log.Info().Msgf("processing template: %s", templateName)

		var templateDataType *daml.DefDataType
		for _, dataType := range module.DataTypes {
			dtName := c.getName(pkg, dataType.GetNameInternedDname())
			if dtName == templateName {
				templateDataType = dataType
				break
			}
		}

		if templateDataType == nil {
			log.Warn().Msgf("could not find data type for template: %s", templateName)
			continue
		}

		tmplStruct := model.TmplStruct{
			Name:       templateName,
			RawType:    "Template",
			IsTemplate: true,
			Choices:    make([]*model.TmplChoice, 0),
		}

		switch v := templateDataType.DataCons.(type) {
		case *daml.DefDataType_Record:
			for _, field := range v.Record.Fields {
				fieldExtracted, typeExtracted, err := c.extractField(pkg, field)
				if err != nil {
					return nil, err
				}
				tmplStruct.Fields = append(tmplStruct.Fields, &model.TmplField{
					Name:    fieldExtracted,
					Type:    typeExtracted,
					RawType: field.String(),
				})
			}
		default:
			log.Warn().Msgf("template %s has non-record data type: %T", templateName, v)
		}

		for _, choice := range template.Choices {
			var choiceName string
			switch v := choice.Name.(type) {
			case *daml.TemplateChoice_NameStr:
				choiceName = v.NameStr
			case *daml.TemplateChoice_NameInternedStr:
				choiceName = pkg.InternedStrings[v.NameInternedStr]
			default:
				log.Warn().Msgf("unknown choice name type: %T", v)
				continue
			}

			choiceStruct := &model.TmplChoice{
				Name:        choiceName,
				IsConsuming: choice.Consuming,
				ArgType:     choice.GetArgBinder().GetType().String(),
				ReturnType:  choice.GetRetType().String(),
			}
			tmplStruct.Choices = append(tmplStruct.Choices, choiceStruct)
		}

		// Extract key if present
		if template.Key != nil {
			keyType := template.Key.GetType().String()
			normalizedKeyType := model.NormalizeDAMLType(keyType)
			log.Info().Msgf("template %s has key of type: %s (normalized: %s)", templateName, keyType, normalizedKeyType)

			// For now, we assume the key is based on a single field that matches the key type
			// In a real implementation, you'd need to parse the key expression to determine which fields are used
			var keyFieldName string

			// Try to find a field that matches the key type - this is a simplified approach
			for _, field := range tmplStruct.Fields {
				if field.Type == normalizedKeyType {
					keyFieldName = field.Name
					break
				}
			}

			// If we can't match by type, we'll use the first field as a fallback
			if keyFieldName == "" && len(tmplStruct.Fields) > 0 {
				keyFieldName = tmplStruct.Fields[0].Name
				log.Warn().Msgf("could not match key type %s for template %s, using first field %s", normalizedKeyType, templateName, keyFieldName)
			}

			tmplStruct.Key = &model.TmplField{
				Name:    keyFieldName,
				Type:    normalizedKeyType,
				RawType: keyType,
			}
		} else {
			log.Info().Msgf("template %s has no key", templateName)
		}

		structs[templateName] = &tmplStruct
	}

	return structs, nil
}

func (c *codeGenAst) getInterfaces(pkg *daml.Package, module *daml.Module) (map[string]*model.TmplStruct, error) {
	structs := make(map[string]*model.TmplStruct, 0)

	for _, iface := range module.Interfaces {
		interfaceName := c.getName(pkg, iface.TyconInternedDname)
		log.Info().Msgf("processing interface: %s", interfaceName)

		tmplStruct := model.TmplStruct{
			Name:        interfaceName,
			RawType:     "Interface",
			IsInterface: true,
			Choices:     make([]*model.TmplChoice, 0),
		}

		// Extract interface choices
		for _, choice := range iface.Choices {
			var choiceName string
			switch v := choice.Name.(type) {
			case *daml.TemplateChoice_NameStr:
				choiceName = v.NameStr
			case *daml.TemplateChoice_NameInternedStr:
				choiceName = pkg.InternedStrings[v.NameInternedStr]
			default:
				log.Warn().Msgf("unknown choice name type: %T", v)
				continue
			}

			choiceStruct := &model.TmplChoice{
				Name:        choiceName,
				IsConsuming: choice.Consuming,
				ArgType:     choice.GetArgBinder().GetType().String(),
				ReturnType:  choice.GetRetType().String(),
			}
			tmplStruct.Choices = append(tmplStruct.Choices, choiceStruct)
		}

		// TODO: Process interface view if needed
		// iface.View contains the view type information

		structs[interfaceName] = &tmplStruct
	}

	return structs, nil
}

func (c *codeGenAst) getDataTypes(pkg *daml.Package, module *daml.Module) (map[string]*model.TmplStruct, error) {
	structs := make(map[string]*model.TmplStruct, 0)
	for _, dataType := range module.GetDataTypes() {
		if !dataType.Serializable {
			continue
		}

		name := c.getName(pkg, dataType.GetNameInternedDname())
		tmplStruct := model.TmplStruct{
			Name: name,
		}

		switch v := dataType.DataCons.(type) {
		case *daml.DefDataType_Record:
			tmplStruct.RawType = "Record"
			for _, field := range v.Record.Fields {
				fieldExtracted, typeExtracted, err := c.extractField(pkg, field)
				if err != nil {
					return nil, err
				}
				tmplStruct.Fields = append(tmplStruct.Fields, &model.TmplField{
					Name:    fieldExtracted,
					Type:    typeExtracted,
					RawType: field.String(),
				})
			}
		case *daml.DefDataType_Variant:
			tmplStruct.RawType = "Variant"
			for _, field := range v.Variant.Fields {
				fieldExtracted, typeExtracted, err := c.extractField(pkg, field)
				if err != nil {
					return nil, err
				}
				tmplStruct.Fields = append(tmplStruct.Fields, &model.TmplField{
					Name:       fieldExtracted,
					Type:       typeExtracted,
					RawType:    field.String(),
					IsOptional: true,
				})
				log.Info().Msgf("variant constructor: %s, type: %s", fieldExtracted, typeExtracted)
			}
		case *daml.DefDataType_Enum:
			tmplStruct.RawType = "Enum"
			for _, constructorIdx := range v.Enum.ConstructorsInternedStr {
				constructorName := c.getName(pkg, constructorIdx)
				tmplStruct.Fields = append(tmplStruct.Fields, &model.TmplField{
					Name: constructorName,
					Type: "enum",
				})
				log.Info().Msgf("enum constructor: %s", constructorName)
			}
		case *daml.DefDataType_Interface:
			tmplStruct.RawType = "Interface"
			log.Warn().Msgf("interface not supported %s", v.Interface.String())
		default:
			log.Warn().Msgf("unknown data cons type: %T", v)
		}
		structs[name] = &tmplStruct
	}

	return structs, nil
}

func (c *codeGenAst) extractField(pkg *daml.Package, field *daml.FieldWithType) (string, string, error) {
	if field == nil {
		return "", "", fmt.Errorf("field is nil")
	}

	internedStrIdx := field.GetFieldInternedStr()
	if int(internedStrIdx) >= len(pkg.InternedStrings) {
		return "", "", fmt.Errorf("invalid interned string index for field name: %d", internedStrIdx)
	}
	fieldName := pkg.InternedStrings[internedStrIdx]
	if field.Type == nil {
		return fieldName, "", fmt.Errorf("field type is nil")
	}

	//	*Type_Var_
	//	*Type_Con_
	//	*Type_Builtin_
	//	*Type_Syn_
	//	*Type_Interned
	var fieldType string
	switch v := field.Type.Sum.(type) {
	case *daml.Type_Interned:
		prim := pkg.InternedTypes[v.Interned]
		if prim != nil {
			isConType := prim.GetCon()
			if isConType != nil {
				tyconName := c.getName(pkg, isConType.Tycon.GetNameInternedDname())
				fieldType = tyconName
			} else {
				fieldType = prim.String()
			}
		} else {
			fieldType = "complex_interned_type"
		}
	case *daml.Type_Con_:
		if v.Con.Tycon != nil {
			switch {
			case v.Con.Tycon.GetNameInternedDname() != 0:
				fieldType = c.getName(pkg, v.Con.Tycon.GetNameInternedDname())
			default:
				fieldType = "unknown_con_type"
			}
		} else {
			fieldType = "con_without_tycon"
		}
	case *daml.Type_Var_:
		switch {
		case v.Var.GetVarInternedStr() != 0:
			fieldType = c.getName(pkg, v.Var.GetVarInternedStr())
		default:
			fieldType = "unnamed_var"
		}
	case *daml.Type_Syn_:
		if v.Syn.Tysyn != nil {
			switch {
			case v.Syn.Tysyn.GetNameInternedDname() != 0:
				fieldType = fmt.Sprintf("syn_%s", c.getName(pkg, v.Syn.Tysyn.GetNameInternedDname()))
			default:
				fieldType = "syn_unknown"
			}
		} else {
			fieldType = "syn_without_name"
		}
	default:
		return fieldName, "", fmt.Errorf("unsupported type sum: %T", field.Type.Sum)
	}

	return fieldName, model.NormalizeDAMLType(fieldType), nil
}
