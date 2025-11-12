package v3

import (
	"errors"
	"fmt"
	"strings"

	daml "github.com/digital-asset/dazl-client/v8/go/api/com/daml/daml_lf_2_1"
	"github.com/noders-team/go-daml/internal/codegen/model"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

const (
	RawTypeTemplate   = "Template"
	RawTypeOptional   = "OPTIONAL"
	RawTypeInterface  = "Interface"
	RawTypeRecord     = "Record"
	RawTypeVariant    = "Variant"
	RawTypeEnum       = "Enum"
	RawTypeContractID = "CONTRACT_ID"
	RawTypeList       = "LIST"
)

type codeGenAst struct {
	payload []byte
}

func NewCodegenAst(payload []byte) *codeGenAst {
	return &codeGenAst{payload: payload}
}

func (c *codeGenAst) isEnumType(typeName string, pkg *daml.Package) bool {
	for _, module := range pkg.Modules {
		for _, dataType := range module.GetDataTypes() {
			if !dataType.Serializable {
				continue
			}

			name := c.getName(pkg, dataType.GetNameInternedDname())
			if name == typeName {
				if _, isEnum := dataType.DataCons.(*daml.DefDataType_Enum); isEnum {
					return true
				}
			}
		}
	}
	return false
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

	damlLf := payloadMapped.GetDamlLf_2()
	if damlLf == nil {
		return nil, errors.New("unsupported daml version")
	}

	for _, module := range damlLf.Modules {
		if len(damlLf.InternedStrings) == 0 {
			continue
		}

		idx := damlLf.InternedDottedNames[module.GetNameInternedDname()].SegmentsInternedStr
		moduleName := damlLf.InternedStrings[idx[len(idx)-1]]

		interfaces, err := c.getInterfaces(damlLf, module, moduleName)
		if err != nil {
			return nil, err
		}
		for key, val := range interfaces {
			interfaceMap[key] = val
		}
	}

	return interfaceMap, nil
}

func (c *codeGenAst) GetTemplateStructs(ifcByModule map[string]model.InterfaceMap) (map[string]*model.TmplStruct, error) {
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

	damlLf := payloadMapped.GetDamlLf_2()
	if damlLf == nil {
		return nil, errors.New("unsupported daml version")
	}

	for _, module := range damlLf.Modules {
		if len(damlLf.InternedStrings) == 0 {
			continue
		}

		idx := damlLf.InternedDottedNames[module.GetNameInternedDname()].SegmentsInternedStr
		moduleName := damlLf.InternedStrings[idx[len(idx)-1]]
		log.Info().Msgf("processing module %s", moduleName)

		dataTypes, err := c.getDataTypes(damlLf, module, moduleName)
		if err != nil {
			return nil, err
		}
		for key, val := range dataTypes {
			structs[key] = val
		}

		templates, err := c.getTemplates(damlLf, module, moduleName, ifcByModule)
		if err != nil {
			return nil, err
		}
		for key, val := range templates {
			structs[key] = val
		}

	}

	return structs, nil
}

func (c *codeGenAst) getName(pkg *daml.Package, id int32) string {
	idx := pkg.InternedDottedNames[id].SegmentsInternedStr
	return pkg.InternedStrings[idx[len(idx)-1]]
}

func (c *codeGenAst) getTemplates(
	pkg *daml.Package, module *daml.Module, moduleName string,
	interfaces map[string]model.InterfaceMap,
) (map[string]*model.TmplStruct, error) {
	structs := make(map[string]*model.TmplStruct, 0)

	for _, template := range module.Templates {
		templateName := c.getName(pkg, template.TyconInternedDname)
		log.Debug().Msgf("processing template: %s", templateName)

		var templateDataType *daml.DefDataType
		for _, dataType := range module.DataTypes {
			dtName := c.getName(pkg, dataType.GetNameInternedDname())
			if dtName == templateName {
				templateDataType = dataType
				break
			}
		}

		if templateDataType == nil {
			log.Debug().Msgf("could not find data type for template: %s", templateName)
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
				isOptional := typeExtracted == RawTypeOptional || strings.HasPrefix(typeExtracted, "*")
				tmplStruct.Fields = append(tmplStruct.Fields, &model.TmplField{
					Name:       fieldExtracted,
					Type:       typeExtracted,
					RawType:    field.String(),
					IsOptional: isOptional,
					IsEnum:     c.isEnumType(typeExtracted, pkg),
				})
			}
		default:
			log.Debug().Msgf("template %s has non-record data type: %T", templateName, v)
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
					log.Debug().Msgf("template %s key field: %s", templateName, keyFieldName)
				}
			}
		}

		if len(template.Implements) > 0 {
			for _, impl := range template.Implements {
				if impl.Interface != nil {
					interfaceName := "I" + c.getName(pkg, impl.Interface.GetNameInternedDname())
					tmplStruct.Implements = append(tmplStruct.Implements, interfaceName)
					ifcModuleName := c.getName(pkg, impl.Interface.Module.ModuleNameInternedDname)
					log.Debug().Msgf("template %s -implements interface: %s location %s", templateName, interfaceName, ifcModuleName)

					if interfaceStruct, exists := interfaces[ifcModuleName][interfaceName]; exists {
						log.Debug().Msgf("found interface %s in map with %d choices", interfaceName, len(interfaceStruct.Choices))
						for _, ifaceChoice := range interfaceStruct.Choices {
							found := false
							for _, tmplChoice := range tmplStruct.Choices {
								if tmplChoice.Name == ifaceChoice.Name {
									found = true
									break
								}
							}
							if !found {
								log.Debug().Msgf("adding interface choice %s to template %s", ifaceChoice.Name, templateName)
								tmplStruct.Choices = append(tmplStruct.Choices, &model.TmplChoice{
									Name:              ifaceChoice.Name,
									ArgType:           ifaceChoice.ArgType,
									ReturnType:        ifaceChoice.ReturnType,
									InterfaceName:     interfaceName,
									InterfaceDAMLName: interfaceStruct.DAMLName,
								})
							}
						}
					}
				}
			}
		}

		structs[templateName] = &tmplStruct
	}

	return structs, nil
}

func (c *codeGenAst) getChoices(pkg *daml.Package, choices []*daml.TemplateChoice) []*model.TmplChoice {
	res := make([]*model.TmplChoice, 0)
	for _, choice := range choices {
		choiceName := pkg.InternedStrings[choice.NameInternedStr]
		choiceStruct := &model.TmplChoice{
			Name: choiceName,
		}

		// Extract argument type if present
		if argBinder := choice.GetArgBinder(); argBinder != nil && argBinder.Type != nil {
			argType := c.extractType(pkg, argBinder.Type)
			// Only set ArgType if it's not a UNIT type
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

func (c *codeGenAst) getInterfaces(pkg *daml.Package, module *daml.Module, moduleName string) (map[string]*model.TmplStruct, error) {
	structs := make(map[string]*model.TmplStruct, 0)

	for _, iface := range module.Interfaces {
		originalName := c.getName(pkg, iface.TyconInternedDname)
		interfaceName := "I" + originalName
		location := c.getName(pkg, iface.Location.Module.GetModuleNameInternedDname())
		log.Debug().Msgf("processing interface: %s, original name %s location %s", interfaceName, originalName, location)

		tmplStruct := model.TmplStruct{
			Name:        interfaceName,
			DAMLName:    originalName,
			ModuleName:  moduleName,
			RawType:     RawTypeInterface,
			IsInterface: true, // TODO dont need as we have RawTypeInterface
			Choices:     make([]*model.TmplChoice, 0),
			Location:    location,
		}
		choices := c.getChoices(pkg, iface.Choices)
		tmplStruct.Choices = append(tmplStruct.Choices, choices...)

		// TODO: Process interface view if needed
		// iface.View contains the view type information

		structs[interfaceName] = &tmplStruct
	}

	return structs, nil
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
				isOptional := typeExtracted == RawTypeOptional || strings.HasPrefix(typeExtracted, "*")
				tmplStruct.Fields = append(tmplStruct.Fields, &model.TmplField{
					Name:       fieldExtracted,
					Type:       typeExtracted,
					RawType:    field.String(),
					IsOptional: isOptional,
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
				if int(constructorIdx) < len(pkg.InternedStrings) {
					constructorName := pkg.InternedStrings[constructorIdx]
					tmplStruct.Fields = append(tmplStruct.Fields, &model.TmplField{
						Name: constructorName,
						Type: "enum",
					})
				}
			}
		case *daml.DefDataType_Interface:
			tmplStruct.RawType = RawTypeInterface
			log.Warn().Msgf("interface not supported %s", v.Interface.String())
		default:
			log.Warn().Msgf("unknown data cons type: %T", v)
		}
		structs[name] = &tmplStruct
	}

	return structs, nil
}

func (c *codeGenAst) parseKeyExpression(pkg *daml.Package, key *daml.DefTemplate_DefKey) []string {
	var fieldNames []string
	if key == nil || key.KeyExpr == nil {
		return fieldNames
	}
	fieldNames = c.parseExpressionForFields(pkg, key.KeyExpr)

	if len(fieldNames) == 0 {
		log.Warn().Msg("could not extract fields from key expression")
	}

	return fieldNames
}

func (c *codeGenAst) parseExpressionForFields(pkg *daml.Package, expr *daml.Expr) []string {
	var fieldNames []string

	if expr == nil {
		return fieldNames
	}

	switch e := expr.Sum.(type) {
	case *daml.Expr_RecProj_:
		if e.RecProj != nil {
			if e.RecProj.FieldInternedStr != 0 {
				fieldName := pkg.InternedStrings[e.RecProj.FieldInternedStr]
				fieldNames = append(fieldNames, fieldName)
			}
			// Also check if the record being projected has more fields
			if e.RecProj.Record != nil {
				subFields := c.parseExpressionForFields(pkg, e.RecProj.Record)
				fieldNames = append(fieldNames, subFields...)
			}
		}
	case *daml.Expr_RecCon_:
		if e.RecCon != nil {
			for _, field := range e.RecCon.Fields {
				if field.FieldInternedStr != 0 {
					fieldName := pkg.InternedStrings[field.FieldInternedStr]
					fieldNames = append(fieldNames, fieldName)
				}
			}
		}
	case *daml.Expr_VarInternedStr:
		if e.VarInternedStr != 0 {
			varName := pkg.InternedStrings[e.VarInternedStr]
			// In template keys, the template parameter is often referenced
			// We'll include variable names as they might represent fields
			fieldNames = append(fieldNames, varName)
		}
	case *daml.Expr_Builtin:
		// Builtin function - might have arguments with field references
		// In DAML LF 2.1, builtins are handled differently
		// For now, we don't extract fields from builtins
	default:
		log.Debug().Msgf("unhandled expression type in key parsing: %T", e)
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

		// TODO: add other types here
		isConType := prim.GetCon()
		if isConType != nil {
			tyconName := c.getName(pkg, isConType.Tycon.GetNameInternedDname())
			fieldType = tyconName
		} else if builtinType := prim.GetBuiltin(); builtinType != nil {
			fieldType = c.handleBuiltinType(pkg, builtinType)
		} else {
			fieldType = prim.String()
		}
	case *daml.Type_Con_:
		if v.Con.Tycon != nil {
			fieldType = c.getName(pkg, v.Con.Tycon.GetNameInternedDname())
		} else {
			fieldType = "con_without_tycon"
		}
	case *daml.Type_Var_:
		// For variables, we use the interned string directly
		if int(v.Var.GetVarInternedStr()) < len(pkg.InternedStrings) {
			fieldType = pkg.InternedStrings[v.Var.GetVarInternedStr()]
		} else {
			fieldType = "unknown_var"
		}
	case *daml.Type_Syn_:
		if v.Syn.Tysyn != nil {
			fieldType = fmt.Sprintf("syn_%s", c.getName(pkg, v.Syn.Tysyn.GetNameInternedDname()))
		} else {
			fieldType = "syn_without_name"
		}
	default:
		fieldType = fmt.Sprintf("unknown_type_%T", typ.Sum)
	}

	return model.NormalizeDAMLType(fieldType)
}

func (c *codeGenAst) handleBuiltinType(pkg *daml.Package, builtinType *daml.Type_Builtin) string {
	builtinName := builtinType.Builtin.String()

	switch builtinType.Builtin {
	case daml.BuiltinType_LIST:
		if len(builtinType.Args) > 0 {
			elementType := c.extractType(pkg, builtinType.Args[0])
			normalizedElementType := model.NormalizeDAMLType(elementType)
			return "[]" + normalizedElementType
		}
		return RawTypeList // fallback to generic LIST
	case daml.BuiltinType_OPTIONAL:
		if len(builtinType.Args) > 0 {
			elementType := c.extractType(pkg, builtinType.Args[0])
			normalizedElementType := model.NormalizeDAMLType(elementType)
			return "*" + normalizedElementType
		}
		return RawTypeOptional // fallback to generic OPTIONAL
	case daml.BuiltinType_CONTRACT_ID:
		return RawTypeContractID
	default:
		return builtinName
	}
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
			} else if builtinType := prim.GetBuiltin(); builtinType != nil {
				fieldType = c.handleBuiltinType(pkg, builtinType)
			} else {
				fieldType = prim.String()
			}
		} else {
			fieldType = "complex_interned_type"
		}
	case *daml.Type_Con_:
		if v.Con.Tycon != nil {
			fieldType = c.getName(pkg, v.Con.Tycon.GetNameInternedDname())
		} else {
			fieldType = "con_without_tycon"
		}
	case *daml.Type_Var_:
		switch {
		case v.Var.GetVarInternedStr() != 0:
			// For variables, we use the interned string directly, not getName which expects DottedName
			if int(v.Var.GetVarInternedStr()) < len(pkg.InternedStrings) {
				fieldType = pkg.InternedStrings[v.Var.GetVarInternedStr()]
			} else {
				fieldType = "unknown_var"
			}
		default:
			fieldType = "unnamed_var"
		}
	case *daml.Type_Syn_:
		if v.Syn.Tysyn != nil {
			fieldType = fmt.Sprintf("syn_%s", c.getName(pkg, v.Syn.Tysyn.GetNameInternedDname()))
		} else {
			fieldType = "syn_without_name"
		}
	default:
		return fieldName, "", fmt.Errorf("unsupported type sum: %T", field.Type.Sum)
	}

	return fieldName, model.NormalizeDAMLType(fieldType), nil
}
