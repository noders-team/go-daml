package {{.Package}}

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"errors"
	"time"
)

var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
)

type PARTY string
type TEXT string
type INT64 int64
type BOOL bool
type DECIMAL *big.Int
type NUMERIC *big.Int
type DATE time.Time
type TIMESTAMP time.Time
type UNIT struct{}
type LIST []string
type MAP map[string]interface{}
type OPTIONAL *interface{}

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
	{{else if eq .RawType "Enum"}}
	// {{capitalise .Name}} is an enum type
	type {{capitalise .Name}} string

	const (
		{{$structName := .Name}}{{range $field := .Fields}}
		{{capitalise $structName}}{{$field.Name}} {{capitalise $structName}} = "{{$field.Name}}"{{end}}
	)
	{{else}}
	// {{capitalise .Name}} is a {{.RawType}} type
	type {{capitalise .Name}} struct {
		{{range $field := .Fields}}
		{{capitalise $field.Name}} {{$field.Type}}{{end}}
	}
	{{end}}
{{end}}
