package gen_clients

import (
	"fmt"

	"github.com/noders-team/go-daml/pkg/codec"
	"github.com/noders-team/go-daml/pkg/model"
	. "github.com/noders-team/go-daml/pkg/types"
)

const packageNameSpliceUtilTokenStandardWallet101 = "splice-util-token-standard-wallet"

// BatchMergeUtility is a Template type
type BatchMergeUtility struct {
	Operator PARTY `json:"operator"`
}

// GetTemplateID returns the template ID for this template
func (t BatchMergeUtility) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceUtilTokenStandardWallet101, "Splice.Util.Token.Wallet.MergeDelegation", "BatchMergeUtility")
}

// CreateCommand returns a CreateCommand for this template
func (t BatchMergeUtility) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["operator"] = t.Operator.ToMap()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for BatchMergeUtility using JsonCodec
func (t BatchMergeUtility) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for BatchMergeUtility using JsonCodec
func (t *BatchMergeUtility) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for BatchMergeUtility

// BatchMergeUtilityBatchMerge exercises the BatchMergeUtility_BatchMerge choice on this BatchMergeUtility contract
func (t BatchMergeUtility) BatchMergeUtilityBatchMerge(contractID string, args BatchMergeUtilityBatchMerge) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceUtilTokenStandardWallet101, "Splice.Util.Token.Wallet.MergeDelegation", "BatchMergeUtility"),
		ContractID: contractID,
		Choice:     "BatchMergeUtility_BatchMerge",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this BatchMergeUtility contract
func (t BatchMergeUtility) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceUtilTokenStandardWallet101, "Splice.Util.Token.Wallet.MergeDelegation", "BatchMergeUtility"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// BatchMergeUtilityBatchMerge is a Record type
type BatchMergeUtilityBatchMerge struct {
	MergeCalls []MergeDelegationCall `json:"mergeCalls"`
}

// ToMap converts BatchMergeUtilityBatchMerge to a map for DAML arguments
func (t BatchMergeUtilityBatchMerge) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["mergeCalls"] = t.MergeCalls

	return m
}

// MarshalJSON implements custom JSON marshaling for BatchMergeUtilityBatchMerge using JsonCodec
func (t BatchMergeUtilityBatchMerge) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for BatchMergeUtilityBatchMerge using JsonCodec
func (t *BatchMergeUtilityBatchMerge) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// BatchMergeUtilityBatchMergeResult is a Record type
type BatchMergeUtilityBatchMergeResult struct {
	Results           []MergeDelegationMergeResult `json:"results"`
	OperatorChangeMap GENMAP                       `json:"operatorChangeMap"`
}

// ToMap converts BatchMergeUtilityBatchMergeResult to a map for DAML arguments
func (t BatchMergeUtilityBatchMergeResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["results"] = t.Results

	m["operatorChangeMap"] = map[string]interface{}{"_type": "genmap", "value": t.OperatorChangeMap}

	return m
}

// MarshalJSON implements custom JSON marshaling for BatchMergeUtilityBatchMergeResult using JsonCodec
func (t BatchMergeUtilityBatchMergeResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for BatchMergeUtilityBatchMergeResult using JsonCodec
func (t *BatchMergeUtilityBatchMergeResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// FeaturedAppRightCall is a Record type
type FeaturedAppRightCall struct {
	AppRightCid   CONTRACT_ID            `json:"appRightCid"`
	Beneficiaries []AppRewardBeneficiary `json:"beneficiaries"`
}

// ToMap converts FeaturedAppRightCall to a map for DAML arguments
func (t FeaturedAppRightCall) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["appRightCid"] = t.AppRightCid

	m["beneficiaries"] = t.Beneficiaries

	return m
}

// MarshalJSON implements custom JSON marshaling for FeaturedAppRightCall using JsonCodec
func (t FeaturedAppRightCall) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for FeaturedAppRightCall using JsonCodec
func (t *FeaturedAppRightCall) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// MergeDelegation is a Template type
type MergeDelegation struct {
	Operator PARTY    `json:"operator"`
	Owner    PARTY    `json:"owner"`
	Meta     Metadata `json:"meta"`
}

// GetTemplateID returns the template ID for this template
func (t MergeDelegation) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceUtilTokenStandardWallet101, "Splice.Util.Token.Wallet.MergeDelegation", "MergeDelegation")
}

// CreateCommand returns a CreateCommand for this template
func (t MergeDelegation) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["operator"] = t.Operator.ToMap()

	args["owner"] = t.Owner.ToMap()

	args["meta"] = t.Meta

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for MergeDelegation using JsonCodec
func (t MergeDelegation) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for MergeDelegation using JsonCodec
func (t *MergeDelegation) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for MergeDelegation

// MergeDelegationMerge exercises the MergeDelegation_Merge choice on this MergeDelegation contract
func (t MergeDelegation) MergeDelegationMerge(contractID string, args MergeDelegationMerge) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceUtilTokenStandardWallet101, "Splice.Util.Token.Wallet.MergeDelegation", "MergeDelegation"),
		ContractID: contractID,
		Choice:     "MergeDelegation_Merge",
		Arguments:  argsToMap(args),
	}
}

// MergeDelegationReject exercises the MergeDelegation_Reject choice on this MergeDelegation contract
func (t MergeDelegation) MergeDelegationReject(contractID string, args MergeDelegationReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceUtilTokenStandardWallet101, "Splice.Util.Token.Wallet.MergeDelegation", "MergeDelegation"),
		ContractID: contractID,
		Choice:     "MergeDelegation_Reject",
		Arguments:  argsToMap(args),
	}
}

// MergeDelegationWithdraw exercises the MergeDelegation_Withdraw choice on this MergeDelegation contract
func (t MergeDelegation) MergeDelegationWithdraw(contractID string, args MergeDelegationWithdraw) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceUtilTokenStandardWallet101, "Splice.Util.Token.Wallet.MergeDelegation", "MergeDelegation"),
		ContractID: contractID,
		Choice:     "MergeDelegation_Withdraw",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this MergeDelegation contract
func (t MergeDelegation) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceUtilTokenStandardWallet101, "Splice.Util.Token.Wallet.MergeDelegation", "MergeDelegation"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// MergeDelegationCall is a Record type
type MergeDelegationCall struct {
	DelegationCid CONTRACT_ID          `json:"delegationCid"`
	ChoiceArg     MergeDelegationMerge `json:"choiceArg"`
}

// ToMap converts MergeDelegationCall to a map for DAML arguments
func (t MergeDelegationCall) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["delegationCid"] = t.DelegationCid

	m["choiceArg"] = t.ChoiceArg

	return m
}

// MarshalJSON implements custom JSON marshaling for MergeDelegationCall using JsonCodec
func (t MergeDelegationCall) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for MergeDelegationCall using JsonCodec
func (t *MergeDelegationCall) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// MergeDelegationProposal is a Template type
type MergeDelegationProposal struct {
	Delegation MergeDelegation `json:"delegation"`
}

// GetTemplateID returns the template ID for this template
func (t MergeDelegationProposal) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceUtilTokenStandardWallet101, "Splice.Util.Token.Wallet.MergeDelegation", "MergeDelegationProposal")
}

// CreateCommand returns a CreateCommand for this template
func (t MergeDelegationProposal) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["delegation"] = t.Delegation

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for MergeDelegationProposal using JsonCodec
func (t MergeDelegationProposal) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for MergeDelegationProposal using JsonCodec
func (t *MergeDelegationProposal) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for MergeDelegationProposal

// MergeDelegationProposalAccept exercises the MergeDelegationProposal_Accept choice on this MergeDelegationProposal contract
func (t MergeDelegationProposal) MergeDelegationProposalAccept(contractID string, args MergeDelegationProposalAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceUtilTokenStandardWallet101, "Splice.Util.Token.Wallet.MergeDelegation", "MergeDelegationProposal"),
		ContractID: contractID,
		Choice:     "MergeDelegationProposal_Accept",
		Arguments:  argsToMap(args),
	}
}

// MergeDelegationProposalReject exercises the MergeDelegationProposal_Reject choice on this MergeDelegationProposal contract
func (t MergeDelegationProposal) MergeDelegationProposalReject(contractID string, args MergeDelegationProposalReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceUtilTokenStandardWallet101, "Splice.Util.Token.Wallet.MergeDelegation", "MergeDelegationProposal"),
		ContractID: contractID,
		Choice:     "MergeDelegationProposal_Reject",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this MergeDelegationProposal contract
func (t MergeDelegationProposal) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceUtilTokenStandardWallet101, "Splice.Util.Token.Wallet.MergeDelegation", "MergeDelegationProposal"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// MergeDelegationProposalWithdraw exercises the MergeDelegationProposal_Withdraw choice on this MergeDelegationProposal contract
func (t MergeDelegationProposal) MergeDelegationProposalWithdraw(contractID string, args MergeDelegationProposalWithdraw) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceUtilTokenStandardWallet101, "Splice.Util.Token.Wallet.MergeDelegation", "MergeDelegationProposal"),
		ContractID: contractID,
		Choice:     "MergeDelegationProposal_Withdraw",
		Arguments:  argsToMap(args),
	}
}

// MergeDelegationProposalAccept is a Record type
type MergeDelegationProposalAccept struct{}

// ToMap converts MergeDelegationProposalAccept to a map for DAML arguments
func (t MergeDelegationProposalAccept) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	return m
}

// MarshalJSON implements custom JSON marshaling for MergeDelegationProposalAccept using JsonCodec
func (t MergeDelegationProposalAccept) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for MergeDelegationProposalAccept using JsonCodec
func (t *MergeDelegationProposalAccept) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// MergeDelegationProposalAcceptResult is a Record type
type MergeDelegationProposalAcceptResult struct {
	MergeDelegationCid CONTRACT_ID `json:"mergeDelegationCid"`
}

// ToMap converts MergeDelegationProposalAcceptResult to a map for DAML arguments
func (t MergeDelegationProposalAcceptResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["mergeDelegationCid"] = t.MergeDelegationCid

	return m
}

// MarshalJSON implements custom JSON marshaling for MergeDelegationProposalAcceptResult using JsonCodec
func (t MergeDelegationProposalAcceptResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for MergeDelegationProposalAcceptResult using JsonCodec
func (t *MergeDelegationProposalAcceptResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// MergeDelegationProposalReject is a Record type
type MergeDelegationProposalReject struct{}

// ToMap converts MergeDelegationProposalReject to a map for DAML arguments
func (t MergeDelegationProposalReject) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	return m
}

// MarshalJSON implements custom JSON marshaling for MergeDelegationProposalReject using JsonCodec
func (t MergeDelegationProposalReject) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for MergeDelegationProposalReject using JsonCodec
func (t *MergeDelegationProposalReject) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// MergeDelegationProposalRejectResult is a Record type
type MergeDelegationProposalRejectResult struct {
	Result UNIT `json:"result"`
}

// ToMap converts MergeDelegationProposalRejectResult to a map for DAML arguments
func (t MergeDelegationProposalRejectResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["result"] = map[string]interface{}{"_type": "unit"}

	return m
}

// MarshalJSON implements custom JSON marshaling for MergeDelegationProposalRejectResult using JsonCodec
func (t MergeDelegationProposalRejectResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for MergeDelegationProposalRejectResult using JsonCodec
func (t *MergeDelegationProposalRejectResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// MergeDelegationProposalWithdraw is a Record type
type MergeDelegationProposalWithdraw struct{}

// ToMap converts MergeDelegationProposalWithdraw to a map for DAML arguments
func (t MergeDelegationProposalWithdraw) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	return m
}

// MarshalJSON implements custom JSON marshaling for MergeDelegationProposalWithdraw using JsonCodec
func (t MergeDelegationProposalWithdraw) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for MergeDelegationProposalWithdraw using JsonCodec
func (t *MergeDelegationProposalWithdraw) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// MergeDelegationProposalWithdrawResult is a Record type
type MergeDelegationProposalWithdrawResult struct {
	Result UNIT `json:"result"`
}

// ToMap converts MergeDelegationProposalWithdrawResult to a map for DAML arguments
func (t MergeDelegationProposalWithdrawResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["result"] = map[string]interface{}{"_type": "unit"}

	return m
}

// MarshalJSON implements custom JSON marshaling for MergeDelegationProposalWithdrawResult using JsonCodec
func (t MergeDelegationProposalWithdrawResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for MergeDelegationProposalWithdrawResult using JsonCodec
func (t *MergeDelegationProposalWithdrawResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// MergeDelegationMerge is a Record type
type MergeDelegationMerge struct {
	OptMergeTransfer    *TransferCall         `json:"optMergeTransfer"`
	OptExtraTransfer    *TransferCall         `json:"optExtraTransfer"`
	OptFeaturedAppRight *FeaturedAppRightCall `json:"optFeaturedAppRight"`
}

// ToMap converts MergeDelegationMerge to a map for DAML arguments
func (t MergeDelegationMerge) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	if t.OptMergeTransfer != nil {
		m["optMergeTransfer"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.OptMergeTransfer,
		}
	} else {
		m["optMergeTransfer"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	if t.OptExtraTransfer != nil {
		m["optExtraTransfer"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.OptExtraTransfer,
		}
	} else {
		m["optExtraTransfer"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	if t.OptFeaturedAppRight != nil {
		m["optFeaturedAppRight"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.OptFeaturedAppRight,
		}
	} else {
		m["optFeaturedAppRight"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for MergeDelegationMerge using JsonCodec
func (t MergeDelegationMerge) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for MergeDelegationMerge using JsonCodec
func (t *MergeDelegationMerge) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// MergeDelegationMergeResult is a Record type
type MergeDelegationMergeResult struct {
	OptMergeTransferResult *TransferInstructionResult `json:"optMergeTransferResult"`
	OptExtraTransferResult *TransferInstructionResult `json:"optExtraTransferResult"`
}

// ToMap converts MergeDelegationMergeResult to a map for DAML arguments
func (t MergeDelegationMergeResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	if t.OptMergeTransferResult != nil {
		m["optMergeTransferResult"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.OptMergeTransferResult,
		}
	} else {
		m["optMergeTransferResult"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	if t.OptExtraTransferResult != nil {
		m["optExtraTransferResult"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.OptExtraTransferResult,
		}
	} else {
		m["optExtraTransferResult"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for MergeDelegationMergeResult using JsonCodec
func (t MergeDelegationMergeResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for MergeDelegationMergeResult using JsonCodec
func (t *MergeDelegationMergeResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// MergeDelegationReject is a Record type
type MergeDelegationReject struct{}

// ToMap converts MergeDelegationReject to a map for DAML arguments
func (t MergeDelegationReject) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	return m
}

// MarshalJSON implements custom JSON marshaling for MergeDelegationReject using JsonCodec
func (t MergeDelegationReject) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for MergeDelegationReject using JsonCodec
func (t *MergeDelegationReject) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// MergeDelegationRejectResult is a Record type
type MergeDelegationRejectResult struct {
	Result UNIT `json:"result"`
}

// ToMap converts MergeDelegationRejectResult to a map for DAML arguments
func (t MergeDelegationRejectResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["result"] = map[string]interface{}{"_type": "unit"}

	return m
}

// MarshalJSON implements custom JSON marshaling for MergeDelegationRejectResult using JsonCodec
func (t MergeDelegationRejectResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for MergeDelegationRejectResult using JsonCodec
func (t *MergeDelegationRejectResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// MergeDelegationWithdraw is a Record type
type MergeDelegationWithdraw struct{}

// ToMap converts MergeDelegationWithdraw to a map for DAML arguments
func (t MergeDelegationWithdraw) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	return m
}

// MarshalJSON implements custom JSON marshaling for MergeDelegationWithdraw using JsonCodec
func (t MergeDelegationWithdraw) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for MergeDelegationWithdraw using JsonCodec
func (t *MergeDelegationWithdraw) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// MergeDelegationWithdrawResult is a Record type
type MergeDelegationWithdrawResult struct {
	Result UNIT `json:"result"`
}

// ToMap converts MergeDelegationWithdrawResult to a map for DAML arguments
func (t MergeDelegationWithdrawResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["result"] = map[string]interface{}{"_type": "unit"}

	return m
}

// MarshalJSON implements custom JSON marshaling for MergeDelegationWithdrawResult using JsonCodec
func (t MergeDelegationWithdrawResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for MergeDelegationWithdrawResult using JsonCodec
func (t *MergeDelegationWithdrawResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferCall is a Record type
type TransferCall struct {
	FactoryCid CONTRACT_ID             `json:"factoryCid"`
	ChoiceArg  TransferFactoryTransfer `json:"choiceArg"`
}

// ToMap converts TransferCall to a map for DAML arguments
func (t TransferCall) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["factoryCid"] = t.FactoryCid

	m["choiceArg"] = t.ChoiceArg

	return m
}

// MarshalJSON implements custom JSON marshaling for TransferCall using JsonCodec
func (t TransferCall) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferCall using JsonCodec
func (t *TransferCall) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}
