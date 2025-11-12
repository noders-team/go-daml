package v2

import (
	"errors"
	"fmt"

	daml "github.com/digital-asset/dazl-client/v8/go/api/com/daml/daml_lf_1_17"
	"github.com/noders-team/go-daml/internal/codegen/model"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

const (
	RawTypeTemplate  = "Template"
	RawTypeOptional  = "OPTIONAL"
	RawTypeInterface = "Interface"
	RawTypeRecord    = "Record"
	RawTypeVariant   = "Variant"
	RawTypeEnum      = "Enum"
)

type codeGenAst struct {
	payload []byte
}

func NewCodegenAst(payload []byte) *codeGenAst {
	return &codeGenAst{payload: payload}
}

func (c *codeGenAst) GetInterfaces() (map[string]*model.TmplStruct, error) {
	interfaceMap := make(map[string]*model.TmplStruct)

	var archive daml.Archive
	err := proto.Unmarshal(c.payload, &archive)
	if err != nil {
		return nil, err
	}

	var payloadMapped daml.ArchivePayload
	err = proto.Unmarshal(archive.Payload, &payloadMapped)
	if err != nil {
		return nil, err
	}

	damlLf1 := payloadMapped.GetDamlLf_1()
	if damlLf1 == nil {
		return nil, errors.New("unsupported daml version")
	}

	for _, module := range damlLf1.Modules {
		if len(damlLf1.InternedStrings) == 0 {
			continue
		}

		idx := damlLf1.InternedDottedNames[module.GetNameInternedDname()].SegmentsInternedStr
		moduleName := damlLf1.InternedStrings[idx[len(idx)-1]]

		interfaces, err := c.getInterfaces(damlLf1, module, moduleName)
		if err != nil {
			return nil, err
		}
		for key, val := range interfaces {
			interfaceMap[key] = val
		}
	}

	return interfaceMap, nil
}

func (c *codeGenAst) GetTemplateStructs(_ map[string]model.InterfaceMap) (map[string]*model.TmplStruct, error) {
	structs := make(map[string]*model.TmplStruct)

	var archive daml.Archive
	err := proto.Unmarshal(c.payload, &archive)
	if err != nil {
		return nil, err
	}

	var payloadMapped daml.ArchivePayload
	err = proto.Unmarshal(archive.Payload, &payloadMapped)
	if err != nil {
		return nil, err
	}

	damlLf1 := payloadMapped.GetDamlLf_1()
	if damlLf1 == nil {
		return nil, errors.New("unsupported daml version")
	}

	for _, module := range damlLf1.Modules {
		if len(damlLf1.InternedStrings) == 0 {
			continue
		}

		idx := damlLf1.InternedDottedNames[module.GetNameInternedDname()].SegmentsInternedStr
		moduleName := damlLf1.InternedStrings[idx[len(idx)-1]]
		log.Info().Msgf("processing module %s", moduleName)

		templates, err := c.getTemplates(damlLf1, module, moduleName)
		if err != nil {
			return nil, err
		}
		for key, val := range templates {
			structs[key] = val
		}

		interfaces, err := c.getInterfaces(damlLf1, module, moduleName)
		if err != nil {
			return nil, err
		}
		for key, val := range interfaces {
			structs[key] = val
		}

		dataTypes, err := c.getDataTypes(damlLf1, module, moduleName)
		if err != nil {
			return nil, err
		}
		for key, val := range dataTypes {
			if _, exists := structs[key]; !exists {
				structs[key] = val
			}
		}
	}

	return structs, nil
}

func (c *codeGenAst) getName(pkg *daml.Package, id int32) string {
	idx := pkg.InternedDottedNames[id].SegmentsInternedStr
	return pkg.InternedStrings[idx[len(idx)-1]]
}

func (c *codeGenAst) getTemplates(pkg *daml.Package, module *daml.Module, moduleName string) (map[string]*model.TmplStruct, error) {
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
			ModuleName: moduleName,
			RawType:    RawTypeTemplate,
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
					Name:       fieldExtracted,
					Type:       typeExtracted,
					RawType:    field.String(),
					IsOptional: typeExtracted == RawTypeOptional,
				})
			}
		default:
			log.Warn().Msgf("template %s has non-record data type: %T", templateName, v)
		}

		choices := c.getChoices(pkg, template.Choices)
		tmplStruct.Choices = append(tmplStruct.Choices, choices...)

		if template.Key != nil {
			keyType := template.Key.GetType().String()
			normalizedKeyType := model.NormalizeDAMLType(keyType)
			log.Debug().Msgf("template %s has key of type: %s (normalized: %s)", templateName, keyType, normalizedKeyType)
			keyFieldNames := c.parseKeyExpression(pkg, template.Key)

			if len(keyFieldNames) > 0 {
				// For now, we support single-field keys
				// TODO: Support composite keys with multiple fields
				keyFieldName := keyFieldNames[0]
				var keyField *model.TmplField
				for _, field := range tmplStruct.Fields {
					if field.Name == keyFieldName {
						keyField = &model.TmplField{
							Name:    field.Name,
							Type:    field.Type,
							RawType: keyType,
						}
						break
					}
				}

				if keyField != nil {
					tmplStruct.Key = keyField
				}
			}
		}

		structs[templateName] = &tmplStruct
	}

	return structs, nil
}

func (c *codeGenAst) getInterfaces(pkg *daml.Package, module *daml.Module, moduleName string) (map[string]*model.TmplStruct, error) {
	structs := make(map[string]*model.TmplStruct, 0)

	for _, iface := range module.Interfaces {
		originalName := c.getName(pkg, iface.TyconInternedDname)
		interfaceName := "I" + originalName
		log.Debug().Msgf("processing interface: %s", interfaceName)

		tmplStruct := model.TmplStruct{
			Name:        interfaceName,
			DAMLName:    originalName,
			ModuleName:  moduleName,
			RawType:     RawTypeInterface,
			IsInterface: true, // TODO dont need
			Choices:     make([]*model.TmplChoice, 0),
		}

		choices := c.getChoices(pkg, iface.Choices)
		tmplStruct.Choices = append(tmplStruct.Choices, choices...)

		structs[interfaceName] = &tmplStruct
	}

	return structs, nil
}

func (c *codeGenAst) getChoices(pkg *daml.Package, choices []*daml.TemplateChoice) []*model.TmplChoice {
	res := make([]*model.TmplChoice, 0)
	for _, choice := range choices {
		var choiceName string
		switch v := choice.Name.(type) {
		case *daml.TemplateChoice_NameStr:
			choiceName = v.NameStr
		case *daml.TemplateChoice_NameInternedStr:
			choiceName = pkg.InternedStrings[v.NameInternedStr]
		default:
			log.Debug().Msgf("unknown choice name type: %T", v)
			continue
		}

		choiceStruct := &model.TmplChoice{
			Name: choiceName,
		}

		if argBinder := choice.GetArgBinder(); argBinder != nil && argBinder.Type != nil {
			argType := c.extractType(pkg, argBinder.Type)
			if argType != "UNIT" && argType != "" {
				choiceStruct.ArgType = argType
			}
		}

		if retType := choice.GetRetType(); retType != nil {
			choiceStruct.ReturnType = c.extractType(pkg, retType)
		}
		res = append(res, choiceStruct)
	}

	return res
}

func (c *codeGenAst) getDataTypes(pkg *daml.Package, module *daml.Module, moduleName string) (map[string]*model.TmplStruct, error) {
	structs := make(map[string]*model.TmplStruct, 0)
	for _, dataType := range module.GetDataTypes() {
		if !dataType.Serializable {
			continue
		}

		name := c.getName(pkg, dataType.GetNameInternedDname())
		tmplStruct := model.TmplStruct{
			Name:       name,
			ModuleName: moduleName,
		}

		switch v := dataType.DataCons.(type) {
		case *daml.DefDataType_Record:
			tmplStruct.RawType = RawTypeRecord
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
			tmplStruct.RawType = RawTypeVariant
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
			}
		case *daml.DefDataType_Enum:
			tmplStruct.RawType = RawTypeEnum
			for _, constructorIdx := range v.Enum.ConstructorsInternedStr {
				constructorName := c.getName(pkg, constructorIdx)
				tmplStruct.Fields = append(tmplStruct.Fields, &model.TmplField{
					Name: constructorName,
					Type: "enum",
				})
			}
		case *daml.DefDataType_Interface:
			tmplStruct.RawType = "Interface"
			log.Warn().Msgf("interface not supported %s, ignoring", v.Interface.String())
		default:
			log.Warn().Msgf("unknown data cons type: %T, ignoring", v)
		}
		structs[name] = &tmplStruct
	}

	return structs, nil
}

func (c *codeGenAst) parseKeyExpression(pkg *daml.Package, key *daml.DefTemplate_DefKey) []string {
	var fieldNames []string

	if key == nil {
		return fieldNames
	}

	if key.GetComplexKey() != nil {
		log.Debug().Msg("complex key expressions are not fully supported yet")
		return fieldNames
	}

	if key.GetKey() == nil {
		log.Debug().Msgf("key expression is nil, ignoring")
		return fieldNames
	}

	keyExpr := key.GetKey()
	if keyExpr.GetProjections() != nil {
		projections := keyExpr.GetProjections()
		for _, proj := range projections.Projections {
			if proj.GetFieldInternedStr() != 0 {
				fieldName := pkg.InternedStrings[proj.GetFieldInternedStr()]
				fieldNames = append(fieldNames, fieldName)
			} else if proj.GetFieldStr() != "" {
				fieldNames = append(fieldNames, proj.GetFieldStr())
			}
		}
	} else if keyExpr.GetRecord() != nil {
		record := keyExpr.GetRecord()
		for _, field := range record.Fields {
			if field.GetFieldInternedStr() != 0 {
				fieldName := pkg.InternedStrings[field.GetFieldInternedStr()]
				fieldNames = append(fieldNames, fieldName)
			} else if field.GetFieldStr() != "" {
				fieldNames = append(fieldNames, field.GetFieldStr())
			}
		}
	}

	return fieldNames
}

func (c *codeGenAst) extractType(pkg *daml.Package, typ *daml.Type) string {
	if typ == nil {
		return ""
	}

	var fieldType string
	switch v := typ.Sum.(type) {
	case *daml.Type_Interned:
		prim := pkg.InternedTypes[v.Interned]
		if prim == nil {
			return "unknown_interned_type"
		}
		isConType := prim.GetCon()
		if isConType != nil {
			tyconName := c.getName(pkg, isConType.Tycon.GetNameInternedDname())
			fieldType = tyconName
		} else {
			fieldType = prim.String()
		}
	case *daml.Type_Con_:
		if v.Con.Tycon == nil {
			return "unknown_vcon_tycon"
		}
		switch {
		case v.Con.Tycon.GetNameInternedDname() != 0:
			fieldType = c.getName(pkg, v.Con.Tycon.GetNameInternedDname())
		default:
			fieldType = "unknown_con_type"
		}
	case *daml.Type_Var_:
		switch {
		case v.Var.GetVarInternedStr() != 0:
			// For v2, variables might need different handling
			if int(v.Var.GetVarInternedStr()) < len(pkg.InternedStrings) {
				fieldType = pkg.InternedStrings[v.Var.GetVarInternedStr()]
			} else {
				fieldType = "unknown_var"
			}
		default:
			fieldType = "unnamed_var"
		}
	case *daml.Type_Syn_:
		if v.Syn.Tysyn == nil {
			return "unknown_syn_without_name"
		}
		switch {
		case v.Syn.Tysyn.GetNameInternedDname() != 0:
			fieldType = fmt.Sprintf("syn_%s", c.getName(pkg, v.Syn.Tysyn.GetNameInternedDname()))
		default:
			fieldType = "syn_unknown"
		}
	default:
		fieldType = fmt.Sprintf("unknown_type_%T", typ.Sum)
	}

	return model.NormalizeDAMLType(fieldType)
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
