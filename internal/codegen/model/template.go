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
	case strings.Contains(damlType, "RelTime") || strings.Contains(damlType, "RELTIME"):
		return "RELTIME"
	case strings.Contains(damlType, "Set") && !strings.Contains(damlType, "Settle") && !strings.Contains(damlType, "Setup"):
		return "SET"
	case strings.Contains(damlType, "Tuple2") || strings.Contains(damlType, "TUPLE2"):
		return "TUPLE2"
	case damlType == "enum":
		return "string"
	case damlType == "19":
		return "GENMAP"
	case damlType == "20":
		return "TEXTMAP"
	case strings.Contains(damlType, "var:{var_interned_str:"):
		return "interface{}"
	case damlType == "prim:{}" || damlType == "{}":
		return "UNIT"
	case damlType == "Archive":
		return "UNIT"
	case strings.HasPrefix(damlType, "[]") && len(damlType) > 2:
		return damlType
	case strings.HasPrefix(damlType, "*") && len(damlType) > 1:
		return damlType
	default:
		// log.Warn().Msgf("unknown daml type %s, using as-is", damlType)
		// Remove underscores from type names to match struct naming convention
		// e.g., "TransferInstruction_Accept" -> "TransferInstructionAccept"
		return strings.ReplaceAll(damlType, "_", "")
	}
}
