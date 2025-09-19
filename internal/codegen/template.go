package codegen

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"strings"
	"text/template"
)

type tmplData struct {
	Package string
	Structs map[string]*tmplStruct
}

//go:embed source.go.tpl
var tmplSource string

func Bind(pkg string, structs map[string]*tmplStruct) (string, error) {
	data := &tmplData{
		Package: pkg,
		Structs: structs,
	}
	buffer := new(bytes.Buffer)

	funcs := map[string]interface{}{
		"capitalise":   capitalize,
		"decapitalise": decapitalize,
	}
	tmpl := template.Must(template.New("").Funcs(funcs).Parse(tmplSource))
	if err := tmpl.Execute(buffer, data); err != nil {
		return "", err
	}
	// Pass the code through gofmt to clean it up
	code, err := format.Source(buffer.Bytes())
	if err != nil {
		return "", fmt.Errorf("%v\n%s", err, buffer)
	}
	return string(code), nil
}

func capitalize(input string) string {
	if len(input) == 0 {
		return input
	}

	if isAllCaps(input) {
		return strings.ToUpper(input[:1]) + strings.ToLower(input[1:])
	}

	if len(input) > 0 && input[0] >= 'A' && input[0] <= 'Z' && !strings.ContainsAny(input, "_- ") {
		return input
	}

	result := toCamelCase(input)
	return strings.ToUpper(result[:1]) + result[1:]
}

func decapitalize(input string) string {
	if len(input) == 0 {
		return input
	}

	if isAllCaps(input) {
		return strings.ToLower(input)
	}

	if len(input) > 0 && input[0] >= 'a' && input[0] <= 'z' && !strings.ContainsAny(input, "_- ") {
		return input
	}

	result := toCamelCase(input)
	return strings.ToLower(result[:1]) + result[1:]
}

func toCamelCase(input string) string {
	if len(input) == 0 {
		return input
	}

	if !strings.ContainsAny(input, "_- ") {
		return input
	}

	words := strings.FieldsFunc(input, func(c rune) bool {
		return c == '_' || c == '-' || c == ' '
	})

	if len(words) == 0 {
		return input
	}

	var result strings.Builder
	for i, word := range words {
		if len(word) == 0 {
			continue
		}
		if i == 0 {
			result.WriteString(strings.ToLower(word[:1]) + strings.ToLower(word[1:]))
		} else {
			result.WriteString(strings.ToUpper(word[:1]) + strings.ToLower(word[1:]))
		}
	}

	return result.String()
}

func isAllCaps(input string) bool {
	if len(input) == 0 {
		return false
	}
	for _, r := range input {
		if r >= 'a' && r <= 'z' {
			return false
		}
	}
	for _, r := range input {
		if r >= 'A' && r <= 'Z' {
			return true
		}
	}
	return false
}
