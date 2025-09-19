package codegen

import "time"

type tmplStruct struct {
	Name   string
	Fields []*tmplField
}

type tmplField struct {
	Type string
	Name string
}

type Package struct {
	Name     string
	Version  string
	Structs  map[string]*tmplStruct
	Metadata *Metadata
}

type Metadata struct {
	Name         string
	Version      string
	Dependencies []string
	LangVersion  string
	CreatedBy    string
	SdkVersion   string
	CreatedAt    *time.Time
}
