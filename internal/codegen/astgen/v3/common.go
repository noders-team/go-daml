package v3

import (
	"errors"
	"fmt"

	daml "github.com/digital-asset/dazl-client/v8/go/api/com/daml/daml_lf_2_1"
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

func (c *codeGenAst) GetTemplateStructs() (map[string]*model.TmplStruct, error) {
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

		// Process templates first (template-centric approach)
		templates, err := c.getTemplates(damlLf, module, moduleName)
		if err != nil {
			return nil, err
		}
		for key, val := range templates {
			structs[key] = val
		}

		// Process interfaces
		interfaces, err := c.getInterfaces(damlLf, module, moduleName)
		if err != nil {
			return nil, err
		}
		for key, val := range interfaces {
			structs[key] = val
		}

		// Process remaining data types that aren't covered by templates/interfaces
		dataTypes, err := c.getDataTypes(damlLf, module, moduleName)
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

		// In DAML LF 2.1, template name is directly in TyconInternedDname
		templateName = c.getName(pkg, template.TyconInternedDname)

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
			ModuleName: moduleName,
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
					Name:       fieldExtracted,
					Type:       typeExtracted,
					RawType:    field.String(),
					IsOptional: typeExtracted == "OPTIONAL",
				})
			}
		default:
			log.Warn().Msgf("template %s has non-record data type: %T", templateName, v)
		}

		for _, choice := range template.Choices {
			// In DAML LF 2.1, choice name is directly in NameInternedStr
			choiceName := pkg.InternedStrings[choice.NameInternedStr]

			choiceStruct := &model.TmplChoice{
				Name:        choiceName,
				IsConsuming: choice.Consuming,
			}

			// Extract argument type if present
			if argBinder := choice.GetArgBinder(); argBinder != nil && argBinder.Type != nil {
				argType := c.extractType(pkg, argBinder.Type)
				// Only set ArgType if it's not a UNIT type
				if argType != "UNIT" && argType != "" {
					choiceStruct.ArgType = argType
				}
			}

			// Extract return type
			if retType := choice.GetRetType(); retType != nil {
				choiceStruct.ReturnType = c.extractType(pkg, retType)
			}

			tmplStruct.Choices = append(tmplStruct.Choices, choiceStruct)
		}

		// Extract key if present
		if template.Key != nil {
			keyType := template.Key.GetType().String()
			normalizedKeyType := model.NormalizeDAMLType(keyType)
			log.Info().Msgf("template %s has key of type: %s (normalized: %s)", templateName, keyType, normalizedKeyType)
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
					log.Info().Msgf("template %s key field: %s", templateName, keyFieldName)
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
		interfaceName := c.getName(pkg, iface.TyconInternedDname)
		log.Info().Msgf("processing interface: %s", interfaceName)

		tmplStruct := model.TmplStruct{
			Name:        interfaceName,
			ModuleName:  moduleName,
			RawType:     "Interface",
			IsInterface: true,
			Choices:     make([]*model.TmplChoice, 0),
		}

		// Extract interface choices
		for _, choice := range iface.Choices {
			// In DAML LF 2.1, choice name is directly in NameInternedStr
			choiceName := pkg.InternedStrings[choice.NameInternedStr]

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
				// For enum constructors, use interned strings directly
				if int(constructorIdx) < len(pkg.InternedStrings) {
					constructorName := pkg.InternedStrings[constructorIdx]
					tmplStruct.Fields = append(tmplStruct.Fields, &model.TmplField{
						Name: constructorName,
						Type: "enum",
					})
					log.Info().Msgf("enum constructor: %s", constructorName)
				}
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

// parseKeyExpression parses the key expression to extract field names used in the key
// In DAML LF 2.1, the key is represented as a general Expr rather than a specialized KeyExpr
func (c *codeGenAst) parseKeyExpression(pkg *daml.Package, key *daml.DefTemplate_DefKey) []string {
	var fieldNames []string

	if key == nil || key.KeyExpr == nil {
		return fieldNames
	}

	// In DAML LF 2.1, we need to parse the general expression
	// This is more complex than 1.17 as it doesn't have specialized KeyExpr types
	fieldNames = c.parseExpressionForFields(pkg, key.KeyExpr)

	if len(fieldNames) == 0 {
		log.Warn().Msg("could not extract fields from key expression")
	}

	return fieldNames
}

// parseExpressionForFields recursively parses an expression to find field references
func (c *codeGenAst) parseExpressionForFields(pkg *daml.Package, expr *daml.Expr) []string {
	var fieldNames []string

	if expr == nil {
		return fieldNames
	}

	// Check the expression type
	switch e := expr.Sum.(type) {
	case *daml.Expr_RecProj_:
		// Record projection (field access)
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
		// Record construction
		if e.RecCon != nil {
			for _, field := range e.RecCon.Fields {
				if field.FieldInternedStr != 0 {
					fieldName := pkg.InternedStrings[field.FieldInternedStr]
					fieldNames = append(fieldNames, fieldName)
				}
			}
		}
	case *daml.Expr_VarInternedStr:
		// Variable reference - might be a field parameter
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
	case *daml.Expr_App_:
		// Function application
		if e.App != nil {
			if e.App.Fun != nil {
				subFields := c.parseExpressionForFields(pkg, e.App.Fun)
				fieldNames = append(fieldNames, subFields...)
			}
			for _, arg := range e.App.Args {
				subFields := c.parseExpressionForFields(pkg, arg)
				fieldNames = append(fieldNames, subFields...)
			}
		}
	default:
		// For other expression types, log for debugging
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
		if prim != nil {
			isConType := prim.GetCon()
			if isConType != nil {
				tyconName := c.getName(pkg, isConType.Tycon.GetNameInternedDname())
				fieldType = tyconName
			} else if builtinType := prim.GetBuiltin(); builtinType != nil {
				// Handle builtin types properly in DAML LF 2.1
				switch builtinType.Builtin.String() {
				case "PARTY":
					fieldType = "PARTY"
				case "TEXT":
					fieldType = "TEXT"
				case "INT64":
					fieldType = "INT64"
				case "BOOL":
					fieldType = "BOOL"
				case "NUMERIC":
					fieldType = "NUMERIC"
				case "DECIMAL":
					fieldType = "DECIMAL"
				case "DATE":
					fieldType = "DATE"
				case "TIMESTAMP":
					fieldType = "TIMESTAMP"
				case "UNIT":
					fieldType = "UNIT"
				case "LIST":
					fieldType = "LIST"
				case "OPTIONAL":
					fieldType = "OPTIONAL"
				default:
					fieldType = builtinType.Builtin.String()
				}
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
			// For variables, we use the interned string directly
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
			} else if builtinType := prim.GetBuiltin(); builtinType != nil {
				// Handle builtin types properly in DAML LF 2.1
				switch builtinType.Builtin.String() {
				case "PARTY":
					fieldType = "PARTY"
				case "TEXT":
					fieldType = "TEXT"
				case "INT64":
					fieldType = "INT64"
				case "BOOL":
					fieldType = "BOOL"
				case "NUMERIC":
					fieldType = "NUMERIC"
				case "DECIMAL":
					fieldType = "DECIMAL"
				case "DATE":
					fieldType = "DATE"
				case "TIMESTAMP":
					fieldType = "TIMESTAMP"
				case "LIST":
					fieldType = "LIST"
				case "OPTIONAL":
					fieldType = "OPTIONAL"
				default:
					fieldType = builtinType.Builtin.String()
				}
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
