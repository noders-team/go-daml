package model

import (
	"strings"
	"time"
)

type InterfaceMap map[string]*TmplStruct

type TmplStruct struct {
	Name        string
	DAMLName    string
	ModuleName  string
	Fields      []*TmplField
	RawType     string
	IsTemplate  bool
	IsInterface bool
	Key         *TmplField
	Choices     []*TmplChoice
	Implements  []string
	Signatories []string
	Observers   []string
	Location    string
}

type TmplField struct {
	Type       string
	Name       string
	RawType    string
	IsOptional bool
	IsEnum     bool
}

type TmplChoice struct {
	Name              string
	ArgType           string
	ReturnType        string
	InterfaceName     string // The Go name of the interface this choice comes from (e.g., "ITransferable")
	InterfaceDAMLName string // The original DAML name of the interface (e.g., "Transferable")
}

type Package struct {
	Name      string
	Version   string
	PackageID string
	Structs   map[string]*TmplStruct
	Metadata  *Metadata
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

var primTypes = []string{
	"PARTY", "TEXT", "INT64", "BOOL", "DECIMAL", "NUMERIC", "DATE",
	"TIMESTAMP", "UNIT", "LIST", "MAP", "OPTIONAL", "CONTRACT_ID",
	"GENMAP", "TEXTMAP", "BIGNUMERIC", "ROUNDING_MODE", "ANY",
}

var tuplePref = []string{
	"TUPLE2[", "TUPLE3[",
	"[]TUPLE2[", "[]TUPLE3[",
	"*TUPLE2[", "*TUPLE3[",
	"*[]TUPLE2[", "*[]TUPLE3[",
}

type typeRule struct {
	match func(string) bool
	norm  func(string) string
}

func constNorm(v string) func(string) string {
	return func(string) string { return v }
}

func asIs(d string) string { return d }

func hasAnyPref(s string, prefs []string) bool {
	for _, p := range prefs {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}

var typeRules = []typeRule{
	{func(d string) bool { return strings.Contains(d, "RelTime") || strings.Contains(d, "RELTIME") }, constNorm("RELTIME")},
	{func(d string) bool {
		return strings.Contains(d, "Set") && !strings.Contains(d, "Settle") && !strings.Contains(d, "Setup")
	}, constNorm("SET")},
	{func(d string) bool { return hasAnyPref(d, tuplePref) }, asIs},
	{func(d string) bool { return strings.Contains(d, "Tuple2") || strings.Contains(d, "TUPLE2") }, constNorm("TUPLE2")},
	{func(d string) bool { return strings.Contains(d, "Tuple3") || strings.Contains(d, "TUPLE3") }, constNorm("TUPLE3")},
	{func(d string) bool { return d == "enum" }, constNorm("string")},
	{func(d string) bool { return d == "19" }, constNorm("GENMAP")},
	{func(d string) bool { return d == "20" }, constNorm("TEXTMAP")},
	{func(d string) bool { return strings.Contains(d, "var:{var_interned_str:") }, constNorm("interface{}")},
	{func(d string) bool { return d == "prim:{}" || d == "{}" }, constNorm("UNIT")},
	{func(d string) bool { return d == "Archive" }, constNorm("UNIT")},
	{func(d string) bool { return strings.HasPrefix(d, "[]") && len(d) > 2 }, asIs},
	{func(d string) bool { return strings.HasPrefix(d, "*") && len(d) > 1 }, asIs},
}

func NormalizeDAMLType(damlType string) string {
	for _, name := range primTypes {
		if damlType == name || strings.Contains(damlType, "prim:"+name) {
			return name
		}
	}
	for _, rule := range typeRules {
		if rule.match(damlType) {
			return rule.norm(damlType)
		}
	}
	return strings.ReplaceAll(damlType, "_", "")
}
