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
	case strings.Contains(damlType, "prim:PARTY"):
		return "PARTY"
	case strings.Contains(damlType, "prim:TEXT"):
		return "TEXT"
	case strings.Contains(damlType, "prim:INT64"):
		return "INT64"
	case strings.Contains(damlType, "prim:BOOL"):
		return "BOOL"
	case strings.Contains(damlType, "prim:DECIMAL"):
		return "DECIMAL"
	case strings.Contains(damlType, "prim:NUMERIC"):
		return "NUMERIC"
	case strings.Contains(damlType, "prim:DATE"):
		return "DATE"
	case strings.Contains(damlType, "prim:TIMESTAMP"):
		return "TIMESTAMP"
	case strings.Contains(damlType, "prim:UNIT"):
		return "UNIT"
	case strings.Contains(damlType, "prim:LIST"):
		return "LIST"
	case strings.Contains(damlType, "prim:MAP"):
		return "MAP"
	case strings.Contains(damlType, "prim:OPTIONAL"):
		return "OPTIONAL"
	case strings.Contains(damlType, "prim:CONTRACT_ID"):
		return "CONTRACT_ID"
	case strings.Contains(damlType, "prim:GENMAP"):
		return "GENMAP"
	case strings.Contains(damlType, "prim:TEXTMAP"):
		return "TEXTMAP"
	case strings.Contains(damlType, "prim:BIGNUMERIC"):
		return "BIGNUMERIC"
	case strings.Contains(damlType, "prim:ROUNDING_MODE"):
		return "ROUNDING_MODE"
	case strings.Contains(damlType, "prim:ANY"):
		return "ANY"
	case damlType == "enum":
		return "string"
	default:
		log.Warn().Msgf("unknown daml type %s", damlType)
		return damlType
	}
}
