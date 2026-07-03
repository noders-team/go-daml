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

const packageNameSpliceApiTokenAllocationInstructionV1100 = "splice-api-token-allocation-instruction-v1"

// IAllocationFactory is a DAML interface
type IAllocationFactory interface {
	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// AllocationFactoryAllocate executes the AllocationFactory_Allocate choice
	AllocationFactoryAllocate(contractID string, args AllocationFactoryAllocate) *model.ExerciseCommand

	// AllocationFactoryPublicFetch executes the AllocationFactory_PublicFetch choice
	AllocationFactoryPublicFetch(contractID string, args AllocationFactoryPublicFetch) *model.ExerciseCommand
}

// IAllocationInstruction is a DAML interface
type IAllocationInstruction interface {
	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// AllocationInstructionWithdraw executes the AllocationInstruction_Withdraw choice
	AllocationInstructionWithdraw(contractID string, args AllocationInstructionWithdraw) *model.ExerciseCommand

	// AllocationInstructionUpdate executes the AllocationInstruction_Update choice
	AllocationInstructionUpdate(contractID string, args AllocationInstructionUpdate) *model.ExerciseCommand
}

// AllocationFactoryView is a Record type
type AllocationFactoryView struct {
	Admin PARTY    `json:"admin"`
	Meta  Metadata `json:"meta"`
}

// ToMap converts AllocationFactoryView to a map for DAML arguments
func (t AllocationFactoryView) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["admin"] = t.Admin.ToMap()

	m["meta"] = t.Meta

	return m
}

// MarshalJSON implements custom JSON marshaling for AllocationFactoryView using JsonCodec
func (t AllocationFactoryView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AllocationFactoryView using JsonCodec
func (t *AllocationFactoryView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AllocationFactoryAllocate is a Record type
type AllocationFactoryAllocate struct {
	ExpectedAdmin    PARTY                   `json:"expectedAdmin"`
	Allocation       AllocationSpecification `json:"allocation"`
	RequestedAt      TIMESTAMP               `json:"requestedAt"`
	InputHoldingCids []CONTRACT_ID           `json:"inputHoldingCids"`
	ExtraArgs        ExtraArgs               `json:"extraArgs"`
}

// ToMap converts AllocationFactoryAllocate to a map for DAML arguments
func (t AllocationFactoryAllocate) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["expectedAdmin"] = t.ExpectedAdmin.ToMap()

	m["allocation"] = t.Allocation

	m["requestedAt"] = t.RequestedAt

	m["inputHoldingCids"] = t.InputHoldingCids

	m["extraArgs"] = t.ExtraArgs

	return m
}

// MarshalJSON implements custom JSON marshaling for AllocationFactoryAllocate using JsonCodec
func (t AllocationFactoryAllocate) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AllocationFactoryAllocate using JsonCodec
func (t *AllocationFactoryAllocate) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AllocationFactoryPublicFetch is a Record type
type AllocationFactoryPublicFetch struct {
	ExpectedAdmin PARTY `json:"expectedAdmin"`
	Actor         PARTY `json:"actor"`
}

// ToMap converts AllocationFactoryPublicFetch to a map for DAML arguments
func (t AllocationFactoryPublicFetch) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["expectedAdmin"] = t.ExpectedAdmin.ToMap()

	m["actor"] = t.Actor.ToMap()

	return m
}

// MarshalJSON implements custom JSON marshaling for AllocationFactoryPublicFetch using JsonCodec
func (t AllocationFactoryPublicFetch) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AllocationFactoryPublicFetch using JsonCodec
func (t *AllocationFactoryPublicFetch) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AllocationInstructionResult is a Record type
type AllocationInstructionResult struct {
	Output           AllocationInstructionResultOutput `json:"output"`
	SenderChangeCids []CONTRACT_ID                     `json:"senderChangeCids"`
	Meta             Metadata                          `json:"meta"`
}

// ToMap converts AllocationInstructionResult to a map for DAML arguments
func (t AllocationInstructionResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["output"] = t.Output

	m["senderChangeCids"] = t.SenderChangeCids

	m["meta"] = t.Meta

	return m
}

// MarshalJSON implements custom JSON marshaling for AllocationInstructionResult using JsonCodec
func (t AllocationInstructionResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AllocationInstructionResult using JsonCodec
func (t *AllocationInstructionResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AllocationInstructionResultCompleted is a Record type
type AllocationInstructionResultCompleted struct {
	AllocationCid CONTRACT_ID `json:"allocationCid"`
}

// ToMap converts AllocationInstructionResultCompleted to a map for DAML arguments
func (t AllocationInstructionResultCompleted) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["allocationCid"] = t.AllocationCid

	return m
}

// MarshalJSON implements custom JSON marshaling for AllocationInstructionResultCompleted using JsonCodec
func (t AllocationInstructionResultCompleted) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AllocationInstructionResultCompleted using JsonCodec
func (t *AllocationInstructionResultCompleted) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AllocationInstructionResultOutput is a variant/union type
type AllocationInstructionResultOutput struct {
	AllocationInstructionResultPending   *AllocationInstructionResultPending   `json:"AllocationInstructionResult_Pending,omitempty"`
	AllocationInstructionResultCompleted *AllocationInstructionResultCompleted `json:"AllocationInstructionResult_Completed,omitempty"`
	AllocationInstructionResultFailed    *UNIT                                 `json:"AllocationInstructionResult_Failed,omitempty"`
}

// MarshalJSON implements custom JSON marshaling for AllocationInstructionResultOutput
func (v AllocationInstructionResultOutput) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(v)
}

// UnmarshalJSON implements custom JSON unmarshaling for AllocationInstructionResultOutput
func (v *AllocationInstructionResultOutput) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, v)
}

// GetVariantTag implements types.VARIANT interface
func (v AllocationInstructionResultOutput) GetVariantTag() string {
	if v.AllocationInstructionResultPending != nil {
		return "AllocationInstructionResult_Pending"
	}

	if v.AllocationInstructionResultCompleted != nil {
		return "AllocationInstructionResult_Completed"
	}

	if v.AllocationInstructionResultFailed != nil {
		return "AllocationInstructionResult_Failed"
	}

	return ""
}

// GetVariantValue implements types.VARIANT interface
func (v AllocationInstructionResultOutput) GetVariantValue() interface{} {
	if v.AllocationInstructionResultPending != nil {
		return v.AllocationInstructionResultPending
	}

	if v.AllocationInstructionResultCompleted != nil {
		return v.AllocationInstructionResultCompleted
	}

	if v.AllocationInstructionResultFailed != nil {
		return v.AllocationInstructionResultFailed
	}

	return nil
}

// Verify interface implementation
var _ VARIANT = (*AllocationInstructionResultOutput)(nil)

// AllocationInstructionResultPending is a Record type
type AllocationInstructionResultPending struct {
	AllocationInstructionCid CONTRACT_ID `json:"allocationInstructionCid"`
}

// ToMap converts AllocationInstructionResultPending to a map for DAML arguments
func (t AllocationInstructionResultPending) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["allocationInstructionCid"] = t.AllocationInstructionCid

	return m
}

// MarshalJSON implements custom JSON marshaling for AllocationInstructionResultPending using JsonCodec
func (t AllocationInstructionResultPending) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AllocationInstructionResultPending using JsonCodec
func (t *AllocationInstructionResultPending) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AllocationInstructionView is a Record type
type AllocationInstructionView struct {
	OriginalInstructionCid *CONTRACT_ID            `json:"originalInstructionCid"`
	Allocation             AllocationSpecification `json:"allocation"`
	PendingActions         GENMAP                  `json:"pendingActions"`
	RequestedAt            TIMESTAMP               `json:"requestedAt"`
	InputHoldingCids       []CONTRACT_ID           `json:"inputHoldingCids"`
	Meta                   Metadata                `json:"meta"`
}

// ToMap converts AllocationInstructionView to a map for DAML arguments
func (t AllocationInstructionView) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	if t.OriginalInstructionCid != nil {
		m["originalInstructionCid"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.OriginalInstructionCid,
		}
	} else {
		m["originalInstructionCid"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	m["allocation"] = t.Allocation

	m["pendingActions"] = map[string]interface{}{"_type": "genmap", "value": t.PendingActions}

	m["requestedAt"] = t.RequestedAt

	m["inputHoldingCids"] = t.InputHoldingCids

	m["meta"] = t.Meta

	return m
}

// MarshalJSON implements custom JSON marshaling for AllocationInstructionView using JsonCodec
func (t AllocationInstructionView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AllocationInstructionView using JsonCodec
func (t *AllocationInstructionView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AllocationInstructionUpdate is a Record type
type AllocationInstructionUpdate struct {
	ExtraActors []PARTY   `json:"extraActors"`
	ExtraArgs   ExtraArgs `json:"extraArgs"`
}

// ToMap converts AllocationInstructionUpdate to a map for DAML arguments
func (t AllocationInstructionUpdate) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["extraActors"] = t.ExtraActors

	m["extraArgs"] = t.ExtraArgs

	return m
}

// MarshalJSON implements custom JSON marshaling for AllocationInstructionUpdate using JsonCodec
func (t AllocationInstructionUpdate) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AllocationInstructionUpdate using JsonCodec
func (t *AllocationInstructionUpdate) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AllocationInstructionWithdraw is a Record type
type AllocationInstructionWithdraw struct {
	ExtraArgs ExtraArgs `json:"extraArgs"`
}

// ToMap converts AllocationInstructionWithdraw to a map for DAML arguments
func (t AllocationInstructionWithdraw) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["extraArgs"] = t.ExtraArgs

	return m
}

// MarshalJSON implements custom JSON marshaling for AllocationInstructionWithdraw using JsonCodec
func (t AllocationInstructionWithdraw) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AllocationInstructionWithdraw using JsonCodec
func (t *AllocationInstructionWithdraw) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// IAllocationFactoryInterfaceID returns the interface ID for the IAllocationFactory interface
func IAllocationFactoryInterfaceID(packageID *string) string {
	pkgName := packageNameSpliceApiTokenAllocationInstructionV1100
	if packageID != nil {
		pkgName = *packageID
	}
	return fmt.Sprintf("#%s:%s:%s", pkgName, "Splice.Api.Token.AllocationInstructionV1", "AllocationFactory")
}

// IAllocationInstructionInterfaceID returns the interface ID for the IAllocationInstruction interface
func IAllocationInstructionInterfaceID(packageID *string) string {
	pkgName := packageNameSpliceApiTokenAllocationInstructionV1100
	if packageID != nil {
		pkgName = *packageID
	}
	return fmt.Sprintf("#%s:%s:%s", pkgName, "Splice.Api.Token.AllocationInstructionV1", "AllocationInstruction")
}
