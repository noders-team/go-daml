package {{.Package}}

import (
	"math/big"
	"strings"
	"errors"
)

var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
)

{{$structs := .Structs}}
{{range $structs}}
	type {{capitalise .Name}} struct {
	{{range $field := .Fields}}
	{{capitalise $field.Name}} {{$field.Type}}{{end}}
	}
{{end}}
