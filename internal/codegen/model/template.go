package model

import (
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type TmplStruct struct {
	Name        string
	Fields      []*TmplField
	RawType     string
	IsTemplate  bool
	IsInterface bool
	Key         *TmplField
	Choices     []*TmplChoice
	Implements  []string
	Signatories []string
	Observers   []string
}

type TmplField struct {
	Type       string
	Name       string
	RawType    string
	IsOptional bool
}

type TmplChoice struct {
	Name        string
	ArgType     string
	ReturnType  string
	IsConsuming bool
	Controllers []string
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

func NormalizeDAMLType(damlType string) string {
	switch {
	// Handle both v1/v2 format (prim:TYPE) and v3 format (TYPE)
	case strings.Contains(damlType, "prim:PARTY") || damlType == "PARTY":
		return "PARTY"
	case strings.Contains(damlType, "prim:TEXT") || damlType == "TEXT":
		return "TEXT"
	case strings.Contains(damlType, "prim:INT64") || damlType == "INT64":
		return "INT64"
	case strings.Contains(damlType, "prim:BOOL") || damlType == "BOOL":
		return "BOOL"
	case strings.Contains(damlType, "prim:DECIMAL") || damlType == "DECIMAL":
		return "DECIMAL"
	case strings.Contains(damlType, "prim:NUMERIC") || damlType == "NUMERIC":
		return "NUMERIC"
	case strings.Contains(damlType, "prim:DATE") || damlType == "DATE":
		return "DATE"
	case strings.Contains(damlType, "prim:TIMESTAMP") || damlType == "TIMESTAMP":
		return "TIMESTAMP"
	case strings.Contains(damlType, "prim:UNIT") || damlType == "UNIT":
		return "UNIT"
	case strings.Contains(damlType, "prim:LIST") || damlType == "LIST":
		return "LIST"
	case strings.Contains(damlType, "prim:MAP") || damlType == "MAP":
		return "MAP"
	case strings.Contains(damlType, "prim:OPTIONAL") || damlType == "OPTIONAL":
		return "OPTIONAL"
	case strings.Contains(damlType, "prim:CONTRACT_ID") || damlType == "CONTRACT_ID":
		return "CONTRACT_ID"
	case strings.Contains(damlType, "prim:GENMAP") || damlType == "GENMAP":
		return "GENMAP"
	case strings.Contains(damlType, "prim:TEXTMAP") || damlType == "TEXTMAP":
		return "TEXTMAP"
	case strings.Contains(damlType, "prim:BIGNUMERIC") || damlType == "BIGNUMERIC":
		return "BIGNUMERIC"
	case strings.Contains(damlType, "prim:ROUNDING_MODE") || damlType == "ROUNDING_MODE":
		return "ROUNDING_MODE"
	case strings.Contains(damlType, "prim:ANY") || damlType == "ANY":
		return "ANY"
	case damlType == "enum":
		return "string"
	// Handle numeric builtin IDs from DAML LF 2.1
	case damlType == "19":
		return "GENMAP" // Builtin ID 19 is typically GENMAP
	case damlType == "20":
		return "TEXTMAP"
	// Handle unresolved variable types as interface{} for now
	case strings.Contains(damlType, "var:{var_interned_str:"):
		return "interface{}"
	default:
		log.Warn().Msgf("unknown daml type %s", damlType)
		return damlType
	}
}
