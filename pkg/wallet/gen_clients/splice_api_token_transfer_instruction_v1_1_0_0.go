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

const packageNameSpliceApiTokenTransferInstructionV1100 = "splice-api-token-transfer-instruction-v1"

// ITransferFactory is a DAML interface
type ITransferFactory interface {
	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// TransferFactoryTransfer executes the TransferFactory_Transfer choice
	TransferFactoryTransfer(contractID string, args TransferFactoryTransfer) *model.ExerciseCommand

	// TransferFactoryPublicFetch executes the TransferFactory_PublicFetch choice
	TransferFactoryPublicFetch(contractID string, args TransferFactoryPublicFetch) *model.ExerciseCommand
}

// ITransferInstruction is a DAML interface
type ITransferInstruction interface {
	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// TransferInstructionAccept executes the TransferInstruction_Accept choice
	TransferInstructionAccept(contractID string, args TransferInstructionAccept) *model.ExerciseCommand

	// TransferInstructionReject executes the TransferInstruction_Reject choice
	TransferInstructionReject(contractID string, args TransferInstructionReject) *model.ExerciseCommand

	// TransferInstructionWithdraw executes the TransferInstruction_Withdraw choice
	TransferInstructionWithdraw(contractID string, args TransferInstructionWithdraw) *model.ExerciseCommand

	// TransferInstructionUpdate executes the TransferInstruction_Update choice
	TransferInstructionUpdate(contractID string, args TransferInstructionUpdate) *model.ExerciseCommand
}

// Transfer2 is a Record type
type Transfer2 struct {
	Sender           PARTY         `json:"sender"`
	Receiver         PARTY         `json:"receiver"`
	Amount           NUMERIC       `json:"amount"`
	InstrumentId     InstrumentId  `json:"instrumentId"`
	RequestedAt      TIMESTAMP     `json:"requestedAt"`
	ExecuteBefore    TIMESTAMP     `json:"executeBefore"`
	InputHoldingCids []CONTRACT_ID `json:"inputHoldingCids"`
	Meta             Metadata      `json:"meta"`
}

// ToMap converts Transfer2 to a map for DAML arguments
func (t Transfer2) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["sender"] = t.Sender.ToMap()

	m["receiver"] = t.Receiver.ToMap()

	m["amount"] = (*big.Int)(t.Amount)

	m["instrumentId"] = t.InstrumentId

	m["requestedAt"] = t.RequestedAt

	m["executeBefore"] = t.ExecuteBefore

	m["inputHoldingCids"] = t.InputHoldingCids

	m["meta"] = t.Meta

	return m
}

// MarshalJSON implements custom JSON marshaling for Transfer2 using JsonCodec
func (t Transfer2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for Transfer2 using JsonCodec
func (t *Transfer2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferFactoryView is a Record type
type TransferFactoryView struct {
	Admin PARTY    `json:"admin"`
	Meta  Metadata `json:"meta"`
}

// ToMap converts TransferFactoryView to a map for DAML arguments
func (t TransferFactoryView) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["admin"] = t.Admin.ToMap()

	m["meta"] = t.Meta

	return m
}

// MarshalJSON implements custom JSON marshaling for TransferFactoryView using JsonCodec
func (t TransferFactoryView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferFactoryView using JsonCodec
func (t *TransferFactoryView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferFactoryPublicFetch is a Record type
type TransferFactoryPublicFetch struct {
	ExpectedAdmin PARTY `json:"expectedAdmin"`
	Actor         PARTY `json:"actor"`
}

// ToMap converts TransferFactoryPublicFetch to a map for DAML arguments
func (t TransferFactoryPublicFetch) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["expectedAdmin"] = t.ExpectedAdmin.ToMap()

	m["actor"] = t.Actor.ToMap()

	return m
}

// MarshalJSON implements custom JSON marshaling for TransferFactoryPublicFetch using JsonCodec
func (t TransferFactoryPublicFetch) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferFactoryPublicFetch using JsonCodec
func (t *TransferFactoryPublicFetch) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferFactoryTransfer is a Record type
type TransferFactoryTransfer struct {
	ExpectedAdmin PARTY     `json:"expectedAdmin"`
	Transfer      Transfer2 `json:"transfer"`
	ExtraArgs     ExtraArgs `json:"extraArgs"`
}

// ToMap converts TransferFactoryTransfer to a map for DAML arguments
func (t TransferFactoryTransfer) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["expectedAdmin"] = t.ExpectedAdmin.ToMap()

	m["transfer"] = t.Transfer

	m["extraArgs"] = t.ExtraArgs

	return m
}

// MarshalJSON implements custom JSON marshaling for TransferFactoryTransfer using JsonCodec
func (t TransferFactoryTransfer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferFactoryTransfer using JsonCodec
func (t *TransferFactoryTransfer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferInstructionResult is a Record type
type TransferInstructionResult struct {
	Output           TransferInstructionResultOutput `json:"output"`
	SenderChangeCids []CONTRACT_ID                   `json:"senderChangeCids"`
	Meta             Metadata                        `json:"meta"`
}

// ToMap converts TransferInstructionResult to a map for DAML arguments
func (t TransferInstructionResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["output"] = t.Output

	m["senderChangeCids"] = t.SenderChangeCids

	m["meta"] = t.Meta

	return m
}

// MarshalJSON implements custom JSON marshaling for TransferInstructionResult using JsonCodec
func (t TransferInstructionResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferInstructionResult using JsonCodec
func (t *TransferInstructionResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferInstructionResultCompleted is a Record type
type TransferInstructionResultCompleted struct {
	ReceiverHoldingCids []CONTRACT_ID `json:"receiverHoldingCids"`
}

// ToMap converts TransferInstructionResultCompleted to a map for DAML arguments
func (t TransferInstructionResultCompleted) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["receiverHoldingCids"] = t.ReceiverHoldingCids

	return m
}

// MarshalJSON implements custom JSON marshaling for TransferInstructionResultCompleted using JsonCodec
func (t TransferInstructionResultCompleted) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferInstructionResultCompleted using JsonCodec
func (t *TransferInstructionResultCompleted) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferInstructionResultOutput is a variant/union type
type TransferInstructionResultOutput struct {
	TransferInstructionResultPending   *TransferInstructionResultPending   `json:"TransferInstructionResult_Pending,omitempty"`
	TransferInstructionResultCompleted *TransferInstructionResultCompleted `json:"TransferInstructionResult_Completed,omitempty"`
	TransferInstructionResultFailed    *UNIT                               `json:"TransferInstructionResult_Failed,omitempty"`
}

// MarshalJSON implements custom JSON marshaling for TransferInstructionResultOutput
func (v TransferInstructionResultOutput) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(v)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferInstructionResultOutput
func (v *TransferInstructionResultOutput) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, v)
}

// GetVariantTag implements types.VARIANT interface
func (v TransferInstructionResultOutput) GetVariantTag() string {
	if v.TransferInstructionResultPending != nil {
		return "TransferInstructionResult_Pending"
	}

	if v.TransferInstructionResultCompleted != nil {
		return "TransferInstructionResult_Completed"
	}

	if v.TransferInstructionResultFailed != nil {
		return "TransferInstructionResult_Failed"
	}

	return ""
}

// GetVariantValue implements types.VARIANT interface
func (v TransferInstructionResultOutput) GetVariantValue() interface{} {
	if v.TransferInstructionResultPending != nil {
		return v.TransferInstructionResultPending
	}

	if v.TransferInstructionResultCompleted != nil {
		return v.TransferInstructionResultCompleted
	}

	if v.TransferInstructionResultFailed != nil {
		return v.TransferInstructionResultFailed
	}

	return nil
}

// Verify interface implementation
var _ VARIANT = (*TransferInstructionResultOutput)(nil)

// TransferInstructionResultPending is a Record type
type TransferInstructionResultPending struct {
	TransferInstructionCid CONTRACT_ID `json:"transferInstructionCid"`
}

// ToMap converts TransferInstructionResultPending to a map for DAML arguments
func (t TransferInstructionResultPending) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["transferInstructionCid"] = t.TransferInstructionCid

	return m
}

// MarshalJSON implements custom JSON marshaling for TransferInstructionResultPending using JsonCodec
func (t TransferInstructionResultPending) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferInstructionResultPending using JsonCodec
func (t *TransferInstructionResultPending) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferInstructionStatus is a variant/union type
type TransferInstructionStatus struct {
	TransferPendingReceiverAcceptance *UNIT                            `json:"TransferPendingReceiverAcceptance,omitempty"`
	TransferPendingInternalWorkflow   *TransferPendingInternalWorkflow `json:"TransferPendingInternalWorkflow,omitempty"`
}

// MarshalJSON implements custom JSON marshaling for TransferInstructionStatus
func (v TransferInstructionStatus) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(v)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferInstructionStatus
func (v *TransferInstructionStatus) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, v)
}

// GetVariantTag implements types.VARIANT interface
func (v TransferInstructionStatus) GetVariantTag() string {
	if v.TransferPendingReceiverAcceptance != nil {
		return "TransferPendingReceiverAcceptance"
	}

	if v.TransferPendingInternalWorkflow != nil {
		return "TransferPendingInternalWorkflow"
	}

	return ""
}

// GetVariantValue implements types.VARIANT interface
func (v TransferInstructionStatus) GetVariantValue() interface{} {
	if v.TransferPendingReceiverAcceptance != nil {
		return v.TransferPendingReceiverAcceptance
	}

	if v.TransferPendingInternalWorkflow != nil {
		return v.TransferPendingInternalWorkflow
	}

	return nil
}

// Verify interface implementation
var _ VARIANT = (*TransferInstructionStatus)(nil)

// TransferInstructionView is a Record type
type TransferInstructionView struct {
	OriginalInstructionCid *CONTRACT_ID              `json:"originalInstructionCid"`
	Transfer               Transfer2                 `json:"transfer"`
	Status                 TransferInstructionStatus `json:"status"`
	Meta                   Metadata                  `json:"meta"`
}

// ToMap converts TransferInstructionView to a map for DAML arguments
func (t TransferInstructionView) ToMap() map[string]interface{} {
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

	m["transfer"] = t.Transfer

	m["status"] = t.Status

	m["meta"] = t.Meta

	return m
}

// MarshalJSON implements custom JSON marshaling for TransferInstructionView using JsonCodec
func (t TransferInstructionView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferInstructionView using JsonCodec
func (t *TransferInstructionView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferInstructionAccept is a Record type
type TransferInstructionAccept struct {
	ExtraArgs ExtraArgs `json:"extraArgs"`
}

// ToMap converts TransferInstructionAccept to a map for DAML arguments
func (t TransferInstructionAccept) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["extraArgs"] = t.ExtraArgs

	return m
}

// MarshalJSON implements custom JSON marshaling for TransferInstructionAccept using JsonCodec
func (t TransferInstructionAccept) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferInstructionAccept using JsonCodec
func (t *TransferInstructionAccept) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferInstructionReject is a Record type
type TransferInstructionReject struct {
	ExtraArgs ExtraArgs `json:"extraArgs"`
}

// ToMap converts TransferInstructionReject to a map for DAML arguments
func (t TransferInstructionReject) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["extraArgs"] = t.ExtraArgs

	return m
}

// MarshalJSON implements custom JSON marshaling for TransferInstructionReject using JsonCodec
func (t TransferInstructionReject) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferInstructionReject using JsonCodec
func (t *TransferInstructionReject) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferInstructionUpdate is a Record type
type TransferInstructionUpdate struct {
	ExtraActors []PARTY   `json:"extraActors"`
	ExtraArgs   ExtraArgs `json:"extraArgs"`
}

// ToMap converts TransferInstructionUpdate to a map for DAML arguments
func (t TransferInstructionUpdate) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["extraActors"] = t.ExtraActors

	m["extraArgs"] = t.ExtraArgs

	return m
}

// MarshalJSON implements custom JSON marshaling for TransferInstructionUpdate using JsonCodec
func (t TransferInstructionUpdate) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferInstructionUpdate using JsonCodec
func (t *TransferInstructionUpdate) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferInstructionWithdraw is a Record type
type TransferInstructionWithdraw struct {
	ExtraArgs ExtraArgs `json:"extraArgs"`
}

// ToMap converts TransferInstructionWithdraw to a map for DAML arguments
func (t TransferInstructionWithdraw) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["extraArgs"] = t.ExtraArgs

	return m
}

// MarshalJSON implements custom JSON marshaling for TransferInstructionWithdraw using JsonCodec
func (t TransferInstructionWithdraw) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferInstructionWithdraw using JsonCodec
func (t *TransferInstructionWithdraw) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferPendingInternalWorkflow is a Record type
type TransferPendingInternalWorkflow struct {
	PendingActions GENMAP `json:"pendingActions"`
}

// ToMap converts TransferPendingInternalWorkflow to a map for DAML arguments
func (t TransferPendingInternalWorkflow) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["pendingActions"] = map[string]interface{}{"_type": "genmap", "value": t.PendingActions}

	return m
}

// MarshalJSON implements custom JSON marshaling for TransferPendingInternalWorkflow using JsonCodec
func (t TransferPendingInternalWorkflow) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferPendingInternalWorkflow using JsonCodec
func (t *TransferPendingInternalWorkflow) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ITransferFactoryInterfaceID returns the interface ID for the ITransferFactory interface
func ITransferFactoryInterfaceID(packageID *string) string {
	pkgName := packageNameSpliceApiTokenTransferInstructionV1100
	if packageID != nil {
		pkgName = *packageID
	}
	return fmt.Sprintf("#%s:%s:%s", pkgName, "Splice.Api.Token.TransferInstructionV1", "TransferFactory")
}

// ITransferInstructionInterfaceID returns the interface ID for the ITransferInstruction interface
func ITransferInstructionInterfaceID(packageID *string) string {
	pkgName := packageNameSpliceApiTokenTransferInstructionV1100
	if packageID != nil {
		pkgName = *packageID
	}
	return fmt.Sprintf("#%s:%s:%s", pkgName, "Splice.Api.Token.TransferInstructionV1", "TransferInstruction")
}
