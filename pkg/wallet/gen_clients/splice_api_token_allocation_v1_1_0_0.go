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

const packageNameSpliceApiTokenAllocationV1100 = "splice-api-token-allocation-v1"

// IAllocation is a DAML interface
type IAllocation interface {
	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// AllocationWithdraw executes the Allocation_Withdraw choice
	AllocationWithdraw(contractID string, args AllocationWithdraw) *model.ExerciseCommand

	// AllocationCancel executes the Allocation_Cancel choice
	AllocationCancel(contractID string, args AllocationCancel) *model.ExerciseCommand

	// AllocationExecuteTransfer executes the Allocation_ExecuteTransfer choice
	AllocationExecuteTransfer(contractID string, args AllocationExecuteTransfer) *model.ExerciseCommand
}

// AllocationSpecification is a Record type
type AllocationSpecification struct {
	Settlement    SettlementInfo `json:"settlement"`
	TransferLegId TEXT           `json:"transferLegId"`
	TransferLeg   TransferLeg    `json:"transferLeg"`
}

// ToMap converts AllocationSpecification to a map for DAML arguments
func (t AllocationSpecification) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["settlement"] = t.Settlement

	m["transferLegId"] = string(t.TransferLegId)

	m["transferLeg"] = t.TransferLeg

	return m
}

// MarshalJSON implements custom JSON marshaling for AllocationSpecification using JsonCodec
func (t AllocationSpecification) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AllocationSpecification using JsonCodec
func (t *AllocationSpecification) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AllocationView is a Record type
type AllocationView struct {
	Allocation  AllocationSpecification `json:"allocation"`
	HoldingCids []CONTRACT_ID           `json:"holdingCids"`
	Meta        Metadata                `json:"meta"`
}

// ToMap converts AllocationView to a map for DAML arguments
func (t AllocationView) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["allocation"] = t.Allocation

	m["holdingCids"] = t.HoldingCids

	m["meta"] = t.Meta

	return m
}

// MarshalJSON implements custom JSON marshaling for AllocationView using JsonCodec
func (t AllocationView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AllocationView using JsonCodec
func (t *AllocationView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AllocationCancel is a Record type
type AllocationCancel struct {
	ExtraArgs ExtraArgs `json:"extraArgs"`
}

// ToMap converts AllocationCancel to a map for DAML arguments
func (t AllocationCancel) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["extraArgs"] = t.ExtraArgs

	return m
}

// MarshalJSON implements custom JSON marshaling for AllocationCancel using JsonCodec
func (t AllocationCancel) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AllocationCancel using JsonCodec
func (t *AllocationCancel) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AllocationCancelResult is a Record type
type AllocationCancelResult struct {
	SenderHoldingCids []CONTRACT_ID `json:"senderHoldingCids"`
	Meta              Metadata      `json:"meta"`
}

// ToMap converts AllocationCancelResult to a map for DAML arguments
func (t AllocationCancelResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["senderHoldingCids"] = t.SenderHoldingCids

	m["meta"] = t.Meta

	return m
}

// MarshalJSON implements custom JSON marshaling for AllocationCancelResult using JsonCodec
func (t AllocationCancelResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AllocationCancelResult using JsonCodec
func (t *AllocationCancelResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AllocationExecuteTransfer is a Record type
type AllocationExecuteTransfer struct {
	ExtraArgs ExtraArgs `json:"extraArgs"`
}

// ToMap converts AllocationExecuteTransfer to a map for DAML arguments
func (t AllocationExecuteTransfer) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["extraArgs"] = t.ExtraArgs

	return m
}

// MarshalJSON implements custom JSON marshaling for AllocationExecuteTransfer using JsonCodec
func (t AllocationExecuteTransfer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AllocationExecuteTransfer using JsonCodec
func (t *AllocationExecuteTransfer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AllocationExecuteTransferResult is a Record type
type AllocationExecuteTransferResult struct {
	SenderHoldingCids   []CONTRACT_ID `json:"senderHoldingCids"`
	ReceiverHoldingCids []CONTRACT_ID `json:"receiverHoldingCids"`
	Meta                Metadata      `json:"meta"`
}

// ToMap converts AllocationExecuteTransferResult to a map for DAML arguments
func (t AllocationExecuteTransferResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["senderHoldingCids"] = t.SenderHoldingCids

	m["receiverHoldingCids"] = t.ReceiverHoldingCids

	m["meta"] = t.Meta

	return m
}

// MarshalJSON implements custom JSON marshaling for AllocationExecuteTransferResult using JsonCodec
func (t AllocationExecuteTransferResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AllocationExecuteTransferResult using JsonCodec
func (t *AllocationExecuteTransferResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AllocationWithdraw is a Record type
type AllocationWithdraw struct {
	ExtraArgs ExtraArgs `json:"extraArgs"`
}

// ToMap converts AllocationWithdraw to a map for DAML arguments
func (t AllocationWithdraw) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["extraArgs"] = t.ExtraArgs

	return m
}

// MarshalJSON implements custom JSON marshaling for AllocationWithdraw using JsonCodec
func (t AllocationWithdraw) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AllocationWithdraw using JsonCodec
func (t *AllocationWithdraw) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AllocationWithdrawResult is a Record type
type AllocationWithdrawResult struct {
	SenderHoldingCids []CONTRACT_ID `json:"senderHoldingCids"`
	Meta              Metadata      `json:"meta"`
}

// ToMap converts AllocationWithdrawResult to a map for DAML arguments
func (t AllocationWithdrawResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["senderHoldingCids"] = t.SenderHoldingCids

	m["meta"] = t.Meta

	return m
}

// MarshalJSON implements custom JSON marshaling for AllocationWithdrawResult using JsonCodec
func (t AllocationWithdrawResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AllocationWithdrawResult using JsonCodec
func (t *AllocationWithdrawResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Reference is a Record type
type Reference struct {
	Id  TEXT         `json:"id"`
	Cid *CONTRACT_ID `json:"cid"`
}

// ToMap converts Reference to a map for DAML arguments
func (t Reference) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["id"] = string(t.Id)

	if t.Cid != nil {
		m["cid"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.Cid,
		}
	} else {
		m["cid"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for Reference using JsonCodec
func (t Reference) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for Reference using JsonCodec
func (t *Reference) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// SettlementInfo is a Record type
type SettlementInfo struct {
	Executor       PARTY     `json:"executor"`
	SettlementRef  Reference `json:"settlementRef"`
	RequestedAt    TIMESTAMP `json:"requestedAt"`
	AllocateBefore TIMESTAMP `json:"allocateBefore"`
	SettleBefore   TIMESTAMP `json:"settleBefore"`
	Meta           Metadata  `json:"meta"`
}

// ToMap converts SettlementInfo to a map for DAML arguments
func (t SettlementInfo) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["executor"] = t.Executor.ToMap()

	m["settlementRef"] = t.SettlementRef

	m["requestedAt"] = t.RequestedAt

	m["allocateBefore"] = t.AllocateBefore

	m["settleBefore"] = t.SettleBefore

	m["meta"] = t.Meta

	return m
}

// MarshalJSON implements custom JSON marshaling for SettlementInfo using JsonCodec
func (t SettlementInfo) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for SettlementInfo using JsonCodec
func (t *SettlementInfo) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferLeg is a Record type
type TransferLeg struct {
	Sender       PARTY        `json:"sender"`
	Receiver     PARTY        `json:"receiver"`
	Amount       NUMERIC      `json:"amount"`
	InstrumentId InstrumentId `json:"instrumentId"`
	Meta         Metadata     `json:"meta"`
}

// ToMap converts TransferLeg to a map for DAML arguments
func (t TransferLeg) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["sender"] = t.Sender.ToMap()

	m["receiver"] = t.Receiver.ToMap()

	m["amount"] = (*big.Int)(t.Amount)

	m["instrumentId"] = t.InstrumentId

	m["meta"] = t.Meta

	return m
}

// MarshalJSON implements custom JSON marshaling for TransferLeg using JsonCodec
func (t TransferLeg) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferLeg using JsonCodec
func (t *TransferLeg) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// IAllocationInterfaceID returns the interface ID for the IAllocation interface
func IAllocationInterfaceID(packageID *string) string {
	pkgName := packageNameSpliceApiTokenAllocationV1100
	if packageID != nil {
		pkgName = *packageID
	}
	return fmt.Sprintf("#%s:%s:%s", pkgName, "Splice.Api.Token.AllocationV1", "Allocation")
}
