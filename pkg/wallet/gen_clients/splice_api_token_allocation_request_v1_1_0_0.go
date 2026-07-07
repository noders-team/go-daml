package gen_clients

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/noders-team/go-daml/pkg/codec"
	"github.com/noders-team/go-daml/pkg/model"
	. "github.com/noders-team/go-daml/pkg/types"
)

var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = fmt.Sprintf
	_ = codec.NewJsonCodec
	_ = model.CreateCommand{}
	_ = NewNumericFromDecimal
)

const packageNameSpliceApiTokenAllocationRequestV1100 = "splice-api-token-allocation-request-v1"

// IAllocationRequest is a DAML interface
type IAllocationRequest interface {
	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// AllocationRequestReject executes the AllocationRequest_Reject choice
	AllocationRequestReject(contractID string, args AllocationRequestReject) *model.ExerciseCommand

	// AllocationRequestWithdraw executes the AllocationRequest_Withdraw choice
	AllocationRequestWithdraw(contractID string, args AllocationRequestWithdraw) *model.ExerciseCommand
}

// AllocationRequestView is a Record type
type AllocationRequestView struct {
	Settlement   SettlementInfo         `json:"settlement"`
	TransferLegs map[string]TransferLeg `json:"transferLegs"`
	Meta         Metadata               `json:"meta"`
}

// ToMap converts AllocationRequestView to a map for DAML arguments
func (t AllocationRequestView) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["settlement"] = t.Settlement

	transferLegs := make(map[string]interface{}, len(t.TransferLegs))
	for key, value := range t.TransferLegs {
		transferLegs[key] = value
	}
	m["transferLegs"] = map[string]interface{}{"_type": "textmap", "value": transferLegs}

	m["meta"] = t.Meta

	return m
}

// MarshalJSON implements custom JSON marshaling for AllocationRequestView using JsonCodec
func (t AllocationRequestView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AllocationRequestView using JsonCodec
func (t *AllocationRequestView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AllocationRequestReject is a Record type
type AllocationRequestReject struct {
	Actor     PARTY     `json:"actor"`
	ExtraArgs ExtraArgs `json:"extraArgs"`
}

// ToMap converts AllocationRequestReject to a map for DAML arguments
func (t AllocationRequestReject) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["actor"] = t.Actor.ToMap()

	m["extraArgs"] = t.ExtraArgs

	return m
}

// MarshalJSON implements custom JSON marshaling for AllocationRequestReject using JsonCodec
func (t AllocationRequestReject) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AllocationRequestReject using JsonCodec
func (t *AllocationRequestReject) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AllocationRequestWithdraw is a Record type
type AllocationRequestWithdraw struct {
	ExtraArgs ExtraArgs `json:"extraArgs"`
}

// ToMap converts AllocationRequestWithdraw to a map for DAML arguments
func (t AllocationRequestWithdraw) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["extraArgs"] = t.ExtraArgs

	return m
}

// MarshalJSON implements custom JSON marshaling for AllocationRequestWithdraw using JsonCodec
func (t AllocationRequestWithdraw) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AllocationRequestWithdraw using JsonCodec
func (t *AllocationRequestWithdraw) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// IAllocationRequestInterfaceID returns the interface ID for the IAllocationRequest interface
func IAllocationRequestInterfaceID(packageID *string) string {
	pkgName := packageNameSpliceApiTokenAllocationRequestV1100
	if packageID != nil {
		pkgName = *packageID
	}
	return fmt.Sprintf("#%s:%s:%s", pkgName, "Splice.Api.Token.AllocationRequestV1", "AllocationRequest")
}
