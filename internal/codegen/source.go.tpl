package {{.Package}}

import (
	"fmt"
	"math/big"
	"strings"
	"errors"
	
	"github.com/noders-team/go-daml/pkg/model"
	. "github.com/noders-team/go-daml/pkg/types"
	"github.com/noders-team/go-daml/pkg/codec"
)

var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
)

const PackageID = "{{.PackageID}}"
const SDKVersion = "{{.SdkVersion}}"

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
}

{{$structs := .Structs}}
{{range $structs}}
{{if .IsInterface}}
// {{capitalise .Name}} is a DAML interface
type {{capitalise .Name}} interface {
	{{range $choice := .Choices}}
	// {{capitalise $choice.Name}} executes the {{$choice.Name}} choice
	{{capitalise $choice.Name}}(contractID string{{if and (ne $choice.ArgType "UNIT") (ne $choice.ArgType "")}}, args {{$choice.ArgType}}{{end}}) *model.ExerciseCommand
	{{end}}
}
{{end}}
{{end}}

func argsToMap(args interface{}) map[string]interface{} {
	if args == nil {
		return map[string]interface{}{}
	}
	
	if m, ok := args.(map[string]interface{}); ok {
		return m
	}
	
	// Check if the type has a toMap method
	type mapper interface {
		toMap() map[string]interface{}
	}
	
	if mapper, ok := args.(mapper); ok {
		return mapper.toMap()
	}
	
	return map[string]interface{}{
		"args": args,
	}
}


{{$structs := .Structs}}
{{range $structs}}
	{{if not .IsInterface}}
	{{if eq .RawType "Variant"}}
	// {{capitalise .Name}} is a variant/union type
	type {{capitalise .Name}} struct {
		{{range $field := .Fields}}
		{{capitalise $field.Name}} *{{$field.Type}} `json:"{{$field.Name}},omitempty"`{{end}}
	}

	// MarshalJSON implements custom JSON marshaling for {{capitalise .Name}}
	func (v {{capitalise .Name}}) MarshalJSON() ([]byte, error) {
		jsonCodec := codec.NewJsonCodec()
		return jsonCodec.Marshall(v)
	}

	// UnmarshalJSON implements custom JSON unmarshaling for {{capitalise .Name}}
	func (v *{{capitalise .Name}}) UnmarshalJSON(data []byte) error {
		jsonCodec := codec.NewJsonCodec()
		return jsonCodec.Unmarshall(data, v)
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

	// MarshalJSON implements custom JSON marshaling for {{capitalise .Name}} using JsonCodec
	func (e {{capitalise .Name}}) MarshalJSON() ([]byte, error) {
		jsonCodec := codec.NewJsonCodec()
		return jsonCodec.Marshall(e)
	}

	// UnmarshalJSON implements custom JSON unmarshaling for {{capitalise .Name}} using JsonCodec
	func (e *{{capitalise .Name}}) UnmarshalJSON(data []byte) error {
		jsonCodec := codec.NewJsonCodec()
		return jsonCodec.Unmarshall(data, e)
	}

	// Verify interface implementation
	var _ ENUM = {{capitalise .Name}}("")
	{{else}}
	// {{capitalise .Name}} is a {{.RawType}} type
	type {{capitalise .Name}} struct {
		{{range $field := .Fields}}
		{{capitalise $field.Name}} {{$field.Type}} `json:"{{$field.Name}}"`{{end}}
	}
	{{if and (eq .RawType "Record") (not .IsTemplate) (not .IsInterface)}}

	// toMap converts {{capitalise .Name}} to a map for DAML arguments
	func (t {{capitalise .Name}}) toMap() map[string]interface{} {
		return map[string]interface{}{
			{{range $field := .Fields}}
			"{{$field.Name}}": {{template "fieldToDAMLValue" $field}},{{end}}
		}
	}

	// MarshalJSON implements custom JSON marshaling for {{capitalise .Name}} using JsonCodec
	func (t {{capitalise .Name}}) MarshalJSON() ([]byte, error) {
		jsonCodec := codec.NewJsonCodec()
		return jsonCodec.Marshall(t)
	}

	// UnmarshalJSON implements custom JSON unmarshaling for {{capitalise .Name}} using JsonCodec
	func (t *{{capitalise .Name}}) UnmarshalJSON(data []byte) error {
		jsonCodec := codec.NewJsonCodec()
		return jsonCodec.Unmarshall(data, t)
	}
	{{end}}
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

	// MarshalJSON implements custom JSON marshaling for {{capitalise .Name}} using JsonCodec
	func (t {{capitalise .Name}}) MarshalJSON() ([]byte, error) {
		jsonCodec := codec.NewJsonCodec()
		return jsonCodec.Marshall(t)
	}

	// UnmarshalJSON implements custom JSON unmarshaling for {{capitalise .Name}} using JsonCodec
	func (t *{{capitalise .Name}}) UnmarshalJSON(data []byte) error {
		jsonCodec := codec.NewJsonCodec()
		return jsonCodec.Unmarshall(data, t)
	}
	{{end}}
	{{if and .IsTemplate .Choices}}
	{{$templateName := .Name}}
	{{$moduleName := .ModuleName}}
	// Choice methods for {{capitalise .Name}}
	{{range $choice := .Choices}}
	// {{capitalise $choice.Name}} exercises the {{$choice.Name}} choice on this {{capitalise $templateName}} contract{{if ne $choice.InterfaceName ""}} via the {{capitalise $choice.InterfaceName}} interface{{end}}
	func (t {{capitalise $templateName}}) {{capitalise $choice.Name}}(contractID string{{if and (ne $choice.ArgType "UNIT") (ne $choice.ArgType "")}}, args {{$choice.ArgType}}{{end}}) *model.ExerciseCommand {
		return &model.ExerciseCommand{
			{{if ne $choice.InterfaceName ""}}TemplateID: fmt.Sprintf("%s:%s:%s", PackageID, "{{$moduleName}}", "{{capitalise $choice.InterfaceName}}"),{{else}}TemplateID: fmt.Sprintf("%s:%s:%s", PackageID, "{{$moduleName}}", "{{capitalise $templateName}}"),{{end}}
			ContractID: contractID,
			Choice: "{{$choice.Name}}",
			{{if and (ne $choice.ArgType "UNIT") (ne $choice.ArgType "")}}Arguments: argsToMap(args),{{else}}Arguments: map[string]interface{}{},{{end}}
		}
	}
	{{end}}
	{{end}}
	{{if and .IsTemplate .Implements}}
	{{$templateName := .Name}}
	// Verify interface implementations for {{capitalise .Name}}
	{{range $interface := .Implements}}
	var _ {{capitalise $interface}} = (*{{capitalise $templateName}})(nil)
	{{end}}
	{{end}}
	{{end}}
	{{end}}
{{end}}

{{define "fieldToDAMLValue"}}{{if .IsOptional}}{{$baseType := stringsTrimPrefix .Type "*"}}{{if eq $baseType "INT64"}}int64(*t.{{capitalise .Name}}){{else if eq $baseType "TEXT"}}string(*t.{{capitalise .Name}}){{else if eq $baseType "BOOL"}}bool(*t.{{capitalise .Name}}){{else if eq $baseType "PARTY"}}(*t.{{capitalise .Name}}).ToMap(){{else if eq $baseType "NUMERIC"}}(*big.Int)(*t.{{capitalise .Name}}){{else if eq $baseType "DECIMAL"}}(*big.Int)(*t.{{capitalise .Name}}){{else if eq $baseType "DATE"}}*t.{{capitalise .Name}}{{else if eq $baseType "TIMESTAMP"}}*t.{{capitalise .Name}}{{else if eq $baseType "UNIT"}}map[string]interface{}{"_type": "unit"}{{else}}*t.{{capitalise .Name}}{{end}}{{else if eq .Type "PARTY"}}t.{{capitalise .Name}}.ToMap(){{else if eq .Type "TEXT"}}string(t.{{capitalise .Name}}){{else if eq .Type "INT64"}}int64(t.{{capitalise .Name}}){{else if eq .Type "BOOL"}}bool(t.{{capitalise .Name}}){{else if eq .Type "NUMERIC"}}(*big.Int)(t.{{capitalise .Name}}){{else if eq .Type "DECIMAL"}}(*big.Int)(t.{{capitalise .Name}}){{else if eq .Type "DATE"}}t.{{capitalise .Name}}{{else if eq .Type "TIMESTAMP"}}t.{{capitalise .Name}}{{else if eq .Type "UNIT"}}map[string]interface{}{"_type": "unit"}{{else if eq .Type "LIST"}}t.{{capitalise .Name}}{{else if eq .Type "GENMAP"}}map[string]interface{}{"_type": "genmap", "value": t.{{capitalise .Name}}}{{else if eq .Type "MAP"}}t.{{capitalise .Name}}{{else if eq .Type "OPTIONAL"}}t.{{capitalise .Name}}{{else if eq .Type "string"}}string(t.{{capitalise .Name}}){{else}}t.{{capitalise .Name}}{{end}}{{end}}

{{define "fieldIsNotEmpty"}}{{if eq .Type "PARTY"}}t.{{capitalise .Name}} != ""{{else if eq .Type "TEXT"}}t.{{capitalise .Name}} != ""{{else if eq .Type "INT64"}}t.{{capitalise .Name}} != 0{{else if eq .Type "BOOL"}}true{{else if eq .Type "NUMERIC"}}t.{{capitalise .Name}} != nil{{else if eq .Type "DECIMAL"}}t.{{capitalise .Name}} != nil{{else if eq .Type "DATE"}}!t.{{capitalise .Name}}.IsZero(){{else if eq .Type "TIMESTAMP"}}!t.{{capitalise .Name}}.IsZero(){{else if eq .Type "LIST"}}len(t.{{capitalise .Name}}) > 0{{else if eq .Type "GENMAP"}}t.{{capitalise .Name}} != nil && len(t.{{capitalise .Name}}) > 0{{else if eq .Type "MAP"}}t.{{capitalise .Name}} != nil && len(t.{{capitalise .Name}}) > 0{{else if eq .Type "OPTIONAL"}}t.{{capitalise .Name}} != nil{{else if .IsOptional}}t.{{capitalise .Name}} != nil{{else if .IsEnum}}t.{{capitalise .Name}} != ""{{else}}t.{{capitalise .Name}} != nil{{end}}{{end}}
