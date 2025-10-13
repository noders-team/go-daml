package {{.Package}}

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"errors"
	
	"github.com/noders-team/go-daml/pkg/model"
	. "github.com/noders-team/go-daml/pkg/types"
)

var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
)

const PackageID = "{{.PackageID}}"

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
}

func argsToMap(args interface{}) map[string]interface{} {
	if args == nil {
		return map[string]interface{}{}
	}
	
	if m, ok := args.(map[string]interface{}); ok {
		return m
	}
	
	return map[string]interface{}{
		"args": args,
	}
}


{{$structs := .Structs}}
{{range $structs}}
	{{if eq .RawType "Variant"}}
	// {{capitalise .Name}} is a variant/union type
	type {{capitalise .Name}} struct {
		{{range $field := .Fields}}
		{{capitalise $field.Name}} *{{$field.Type}} `json:"{{$field.Name}},omitempty"`{{end}}
	}

	// MarshalJSON implements custom JSON marshaling for {{capitalise .Name}}
	func (v {{capitalise .Name}}) MarshalJSON() ([]byte, error) {
		{{range $field := .Fields}}
		if v.{{capitalise $field.Name}} != nil {
			return json.Marshal(map[string]interface{}{
				"tag":   "{{$field.Name}}",
				"value": v.{{capitalise $field.Name}},
			})
		}
		{{end}}
		return json.Marshal(map[string]interface{}{})
	}

	// UnmarshalJSON implements custom JSON unmarshaling for {{capitalise .Name}}
	func (v *{{capitalise .Name}}) UnmarshalJSON(data []byte) error {
		var tagged struct {
			Tag   string          `json:"tag"`
			Value json.RawMessage `json:"value"`
		}
		
		if err := json.Unmarshal(data, &tagged); err != nil {
			return err
		}
		
		switch tagged.Tag {
		{{range $field := .Fields}}
		case "{{$field.Name}}":
			var value {{$field.Type}}
			if err := json.Unmarshal(tagged.Value, &value); err != nil {
				return err
			}
			v.{{capitalise $field.Name}} = &value
		{{end}}
		default:
			return fmt.Errorf("unknown tag: %s", tagged.Tag)
		}
		
		return nil
	}
	
	// GetVariantTag implements types.VARIANT interface
	func (v {{capitalise .Name}}) GetVariantTag() string {
		{{range $field := .Fields}}
		if v.{{capitalise $field.Name}} != nil {
			return "{{$field.Name}}"
		}
		{{end}}
		return ""
	}
	
	// GetVariantValue implements types.VARIANT interface
	func (v {{capitalise .Name}}) GetVariantValue() interface{} {
		{{range $field := .Fields}}
		if v.{{capitalise $field.Name}} != nil {
			return v.{{capitalise $field.Name}}
		}
		{{end}}
		return nil
	}
	
	// Verify interface implementation
	var _ VARIANT = (*{{capitalise .Name}})(nil)
	{{else if eq .RawType "Enum"}}
	// {{capitalise .Name}} is an enum type
	type {{capitalise .Name}} string

	const (
		{{$structName := .Name}}{{range $field := .Fields}}
		{{capitalise $structName}}{{$field.Name}} {{capitalise $structName}} = "{{$field.Name}}"{{end}}
	)

	// GetEnumConstructor implements types.ENUM interface
	func (e {{capitalise .Name}}) GetEnumConstructor() string {
		return string(e)
	}

	// GetEnumTypeID implements types.ENUM interface
	func (e {{capitalise .Name}}) GetEnumTypeID() string {
		return fmt.Sprintf("%s:%s:%s", PackageID, "{{.ModuleName}}", "{{capitalise .Name}}")
	}

	// Verify interface implementation
	var _ ENUM = {{capitalise .Name}}("")
	{{else}}
	// {{capitalise .Name}} is a {{.RawType}} type
	type {{capitalise .Name}} struct {
		{{range $field := .Fields}}
		{{capitalise $field.Name}} {{$field.Type}} `json:"{{$field.Name}}"`{{end}}
	}
	{{if .IsTemplate}}
	
	// GetTemplateID returns the template ID for this template
	func (t {{capitalise .Name}}) GetTemplateID() string {
		return fmt.Sprintf("%s:%s:%s", PackageID, "{{.ModuleName}}", "{{capitalise .Name}}")
	}
	
	// CreateCommand returns a CreateCommand for this template
	func (t {{capitalise .Name}}) CreateCommand() *model.CreateCommand {
		args := make(map[string]interface{})
		{{range $field := .Fields}}
		{{if $field.IsOptional}}
		if t.{{capitalise $field.Name}} != nil {
			args["{{$field.Name}}"] = map[string]interface{}{
				"_type": "optional",
				"value": {{template "fieldToDAMLValue" $field}},
			}
		} else {
			args["{{$field.Name}}"] = map[string]interface{}{
				"_type": "optional",
			}
		}
		{{else if or $field.IsEnum (eq $field.Type "GENMAP") (eq $field.Type "MAP") (eq $field.Type "LIST") (eq $field.Type "NUMERIC") (eq $field.Type "DECIMAL")}}
		if {{template "fieldIsNotEmpty" $field}} {
			args["{{$field.Name}}"] = {{template "fieldToDAMLValue" $field}}
		}{{else}}
		args["{{$field.Name}}"] = {{template "fieldToDAMLValue" $field}}{{end}}
		{{end}}
		return &model.CreateCommand{
			TemplateID: t.GetTemplateID(),
			Arguments: args,
		}
	}
	{{if .Key}}
	
	// GetKey returns the key for this template as a string
	func (t {{capitalise .Name}}) GetKey() string {
		{{if eq .Key.Type "TEXT"}}
		return string(t.{{capitalise .Key.Name}})
		{{else if eq .Key.Type "PARTY"}}
		return string(t.{{capitalise .Key.Name}})
		{{else if eq .Key.Type "INT64"}}
		return fmt.Sprintf("%d", t.{{capitalise .Key.Name}})
		{{else}}
		return fmt.Sprintf("%v", t.{{capitalise .Key.Name}})
		{{end}}
	}
	{{end}}
	{{end}}
	{{if and .IsTemplate .Choices}}
	{{$templateName := .Name}}
	{{$moduleName := .ModuleName}}
	// Choice methods for {{capitalise .Name}}
	{{range $choice := .Choices}}
	// {{capitalise $choice.Name}} exercises the {{$choice.Name}} choice on this {{capitalise $templateName}} contract
	func (t {{capitalise $templateName}}) {{capitalise $choice.Name}}(contractID string{{if $choice.ArgType}}, args {{$choice.ArgType}}{{end}}) *model.ExerciseCommand {
		return &model.ExerciseCommand{
			TemplateID: fmt.Sprintf("%s:%s:%s", PackageID, "{{$moduleName}}", "{{capitalise $templateName}}"),
			ContractID: contractID,
			Choice: "{{$choice.Name}}",
			{{if $choice.ArgType}}Arguments: argsToMap(args),{{else}}Arguments: map[string]interface{}{},{{end}}
		}
	}
	{{end}}
	{{end}}
	{{end}}
{{end}}

{{define "fieldToDAMLValue"}}{{if .IsOptional}}{{$baseType := stringsTrimPrefix .Type "*"}}{{if eq $baseType "INT64"}}int64(*t.{{capitalise .Name}}){{else if eq $baseType "TEXT"}}string(*t.{{capitalise .Name}}){{else if eq $baseType "BOOL"}}bool(*t.{{capitalise .Name}}){{else if eq $baseType "PARTY"}}map[string]interface{}{"_type": "party", "value": string(*t.{{capitalise .Name}})}{{else if eq $baseType "NUMERIC"}}(*big.Int)(*t.{{capitalise .Name}}){{else if eq $baseType "DECIMAL"}}(*big.Int)(*t.{{capitalise .Name}}){{else if eq $baseType "DATE"}}*t.{{capitalise .Name}}{{else if eq $baseType "TIMESTAMP"}}*t.{{capitalise .Name}}{{else if eq $baseType "UNIT"}}map[string]interface{}{"_type": "unit"}{{else}}*t.{{capitalise .Name}}{{end}}{{else if eq .Type "PARTY"}}map[string]interface{}{"_type": "party", "value": string(t.{{capitalise .Name}})}{{else if eq .Type "TEXT"}}string(t.{{capitalise .Name}}){{else if eq .Type "INT64"}}int64(t.{{capitalise .Name}}){{else if eq .Type "BOOL"}}bool(t.{{capitalise .Name}}){{else if eq .Type "NUMERIC"}}(*big.Int)(t.{{capitalise .Name}}){{else if eq .Type "DECIMAL"}}(*big.Int)(t.{{capitalise .Name}}){{else if eq .Type "DATE"}}t.{{capitalise .Name}}{{else if eq .Type "TIMESTAMP"}}t.{{capitalise .Name}}{{else if eq .Type "UNIT"}}map[string]interface{}{"_type": "unit"}{{else if eq .Type "LIST"}}t.{{capitalise .Name}}{{else if eq .Type "GENMAP"}}map[string]interface{}{"_type": "genmap", "value": t.{{capitalise .Name}}}{{else if eq .Type "MAP"}}t.{{capitalise .Name}}{{else if eq .Type "OPTIONAL"}}t.{{capitalise .Name}}{{else if eq .Type "string"}}string(t.{{capitalise .Name}}){{else}}t.{{capitalise .Name}}{{end}}{{end}}

{{define "fieldIsNotEmpty"}}{{if eq .Type "PARTY"}}t.{{capitalise .Name}} != ""{{else if eq .Type "TEXT"}}t.{{capitalise .Name}} != ""{{else if eq .Type "INT64"}}t.{{capitalise .Name}} != 0{{else if eq .Type "BOOL"}}true{{else if eq .Type "NUMERIC"}}t.{{capitalise .Name}} != nil{{else if eq .Type "DECIMAL"}}t.{{capitalise .Name}} != nil{{else if eq .Type "DATE"}}!t.{{capitalise .Name}}.IsZero(){{else if eq .Type "TIMESTAMP"}}!t.{{capitalise .Name}}.IsZero(){{else if eq .Type "LIST"}}len(t.{{capitalise .Name}}) > 0{{else if eq .Type "GENMAP"}}t.{{capitalise .Name}} != nil && len(t.{{capitalise .Name}}) > 0{{else if eq .Type "MAP"}}t.{{capitalise .Name}} != nil && len(t.{{capitalise .Name}}) > 0{{else if eq .Type "OPTIONAL"}}t.{{capitalise .Name}} != nil{{else if .IsOptional}}t.{{capitalise .Name}} != nil{{else if .IsEnum}}t.{{capitalise .Name}} != ""{{else}}t.{{capitalise .Name}} != nil{{end}}{{end}}
