package v3

import (
	"errors"

	daml "github.com/digital-asset/dazl-client/v8/go/api/com/daml/daml_lf_2_1"
	"github.com/noders-team/go-daml/internal/codegen/model"
	"google.golang.org/protobuf/proto"
)

type codeGenAst struct {
	payload []byte
}

func NewCodegenAst(payload []byte) *codeGenAst {
	return &codeGenAst{payload: payload}
}

func (c *codeGenAst) GetTemplateStructs() (string, map[string]*model.TmplStruct, error) {
	structs := make(map[string]*model.TmplStruct)

	var archive daml.Archive
	err := proto.Unmarshal(c.payload, &archive)
	if err != nil {
		return "", nil, err
	}

	var payloadMapped daml.ArchivePayload
	err = proto.Unmarshal(archive.Payload, &payloadMapped)
	if err != nil {
		return "", nil, err
	}

	damlLf := payloadMapped.GetDamlLf_2()
	if damlLf == nil {
		return "", nil, errors.New("unsupported daml version")
	}

	return archive.Hash, structs, nil
}
