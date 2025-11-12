package astgen

import (
	"fmt"

	v2 "github.com/noders-team/go-daml/internal/codegen/astgen/v2"
	v3 "github.com/noders-team/go-daml/internal/codegen/astgen/v3"
	"github.com/noders-team/go-daml/internal/codegen/model"
)

const (
	V1 = "1."
	V2 = "2."
	V3 = "3."
)

type AstGen interface {
	GetInterfaces() (map[string]*model.TmplStruct, error)
	GetTemplateStructs(ifcByModule map[string]model.InterfaceMap) (map[string]*model.TmplStruct, error)
}

func GetAstGenFromVersion(payload []byte, ver string) (AstGen, error) {
	switch ver {
	case V2:
		return v2.NewCodegenAst(payload), nil
	case V3:
		return v3.NewCodegenAst(payload), nil
	default:
		return nil, fmt.Errorf("none supported version")
	}
}
