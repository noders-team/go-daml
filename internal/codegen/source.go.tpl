package {{.Package}}

import (
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
type LIST []interface{}
type MAP map[string]interface{}
type OPTIONAL *interface{}

{{$structs := .Structs}}
{{range $structs}}
	type {{capitalise .Name}} struct {
	{{range $field := .Fields}}
	{{capitalise $field.Name}} {{$field.Type}}{{end}}
	}
{{end}}
