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

const packageNameSpliceAmulet0116 = "splice-amulet"

// Amulet is an enum type
type Amulet string

const (
	AmuletAmulet Amulet = "Amulet"
)

// GetEnumConstructor implements types.ENUM interface
func (e Amulet) GetEnumConstructor() string {
	return string(e)
}

// GetEnumTypeID implements types.ENUM interface
func (e Amulet) GetEnumTypeID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletConfig", "Amulet")
}

// MarshalJSON implements custom JSON marshaling for Amulet using JsonCodec
func (e Amulet) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(e)
}

// UnmarshalJSON implements custom JSON unmarshaling for Amulet using JsonCodec
func (e *Amulet) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, e)
}

// Verify interface implementation
var _ ENUM = Amulet("")

// AmuletAllocation is a Template type
type AmuletAllocation struct {
	LockedAmulet CONTRACT_ID             `json:"lockedAmulet"`
	Allocation   AllocationSpecification `json:"allocation"`
}

// GetTemplateID returns the template ID for this template
func (t AmuletAllocation) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletAllocation", "AmuletAllocation")
}

// CreateCommand returns a CreateCommand for this template
func (t AmuletAllocation) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["lockedAmulet"] = t.LockedAmulet

	args["allocation"] = t.Allocation

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for AmuletAllocation using JsonCodec
func (t AmuletAllocation) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletAllocation using JsonCodec
func (t *AmuletAllocation) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for AmuletAllocation

// Archive exercises the Archive choice on this AmuletAllocation contract
func (t AmuletAllocation) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletAllocation", "AmuletAllocation"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// AllocationWithdraw exercises the Allocation_Withdraw choice on this AmuletAllocation contract via the IAllocation interface
func (t AmuletAllocation) AllocationWithdraw(contractID string, args AllocationWithdraw) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletAllocation", "Allocation"),
		ContractID: contractID,
		Choice:     "Allocation_Withdraw",
		Arguments:  argsToMap(args),
	}
}

// AllocationCancel exercises the Allocation_Cancel choice on this AmuletAllocation contract via the IAllocation interface
func (t AmuletAllocation) AllocationCancel(contractID string, args AllocationCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletAllocation", "Allocation"),
		ContractID: contractID,
		Choice:     "Allocation_Cancel",
		Arguments:  argsToMap(args),
	}
}

// AllocationExecuteTransfer exercises the Allocation_ExecuteTransfer choice on this AmuletAllocation contract via the IAllocation interface
func (t AmuletAllocation) AllocationExecuteTransfer(contractID string, args AllocationExecuteTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletAllocation", "Allocation"),
		ContractID: contractID,
		Choice:     "Allocation_ExecuteTransfer",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for AmuletAllocation

var _ IAllocation = (*AmuletAllocation)(nil)

// AmuletConfig is a Record type
type AmuletConfig struct {
	TransferConfig                  TransferConfig                        `json:"transferConfig"`
	IssuanceCurve                   Schedule                              `json:"issuanceCurve"`
	DecentralizedSynchronizer       AmuletDecentralizedSynchronizerConfig `json:"decentralizedSynchronizer"`
	TickDuration                    RELTIME                               `json:"tickDuration"`
	PackageConfig                   PackageConfig                         `json:"packageConfig"`
	TransferPreapprovalFee          *NUMERIC                              `json:"transferPreapprovalFee"`
	FeaturedAppActivityMarkerAmount *NUMERIC                              `json:"featuredAppActivityMarkerAmount"`
	OptDevelopmentFundManager       *PARTY                                `json:"optDevelopmentFundManager"`
}

// ToMap converts AmuletConfig to a map for DAML arguments
func (t AmuletConfig) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["transferConfig"] = t.TransferConfig

	m["issuanceCurve"] = t.IssuanceCurve

	m["decentralizedSynchronizer"] = t.DecentralizedSynchronizer

	m["tickDuration"] = t.TickDuration

	m["packageConfig"] = t.PackageConfig

	if t.TransferPreapprovalFee != nil {
		m["transferPreapprovalFee"] = map[string]interface{}{
			"_type": "optional",
			"value": (*big.Int)(*t.TransferPreapprovalFee),
		}
	} else {
		m["transferPreapprovalFee"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	if t.FeaturedAppActivityMarkerAmount != nil {
		m["featuredAppActivityMarkerAmount"] = map[string]interface{}{
			"_type": "optional",
			"value": (*big.Int)(*t.FeaturedAppActivityMarkerAmount),
		}
	} else {
		m["featuredAppActivityMarkerAmount"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	if t.OptDevelopmentFundManager != nil {
		m["optDevelopmentFundManager"] = map[string]interface{}{
			"_type": "optional",
			"value": (*t.OptDevelopmentFundManager).ToMap(),
		}
	} else {
		m["optDevelopmentFundManager"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletConfig using JsonCodec
func (t AmuletConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletConfig using JsonCodec
func (t *AmuletConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletCreateSummary is a Record type
type AmuletCreateSummary struct {
	Amulet      interface{} `json:"amulet"`
	AmuletPrice NUMERIC     `json:"amuletPrice"`
	Round       Round       `json:"round"`
}

// ToMap converts AmuletCreateSummary to a map for DAML arguments
func (t AmuletCreateSummary) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["amulet"] = t.Amulet

	m["amuletPrice"] = (*big.Int)(t.AmuletPrice)

	m["round"] = t.Round

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletCreateSummary using JsonCodec
func (t AmuletCreateSummary) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletCreateSummary using JsonCodec
func (t *AmuletCreateSummary) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletDecentralizedSynchronizerConfig is a Record type
type AmuletDecentralizedSynchronizerConfig struct {
	RequiredSynchronizers SET                    `json:"requiredSynchronizers"`
	ActiveSynchronizer    TEXT                   `json:"activeSynchronizer"`
	Fees                  SynchronizerFeesConfig `json:"fees"`
}

// ToMap converts AmuletDecentralizedSynchronizerConfig to a map for DAML arguments
func (t AmuletDecentralizedSynchronizerConfig) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["requiredSynchronizers"] = t.RequiredSynchronizers

	m["activeSynchronizer"] = string(t.ActiveSynchronizer)

	m["fees"] = t.Fees

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletDecentralizedSynchronizerConfig using JsonCodec
func (t AmuletDecentralizedSynchronizerConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletDecentralizedSynchronizerConfig using JsonCodec
func (t *AmuletDecentralizedSynchronizerConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletExpireSummary is a Record type
type AmuletExpireSummary struct {
	Owner                              PARTY   `json:"owner"`
	Round                              Round   `json:"round"`
	ChangeToInitialAmountAsOfRoundZero NUMERIC `json:"changeToInitialAmountAsOfRoundZero"`
	ChangeToHoldingFeesRate            NUMERIC `json:"changeToHoldingFeesRate"`
}

// ToMap converts AmuletExpireSummary to a map for DAML arguments
func (t AmuletExpireSummary) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["owner"] = t.Owner.ToMap()

	m["round"] = t.Round

	m["changeToInitialAmountAsOfRoundZero"] = (*big.Int)(t.ChangeToInitialAmountAsOfRoundZero)

	m["changeToHoldingFeesRate"] = (*big.Int)(t.ChangeToHoldingFeesRate)

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletExpireSummary using JsonCodec
func (t AmuletExpireSummary) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletExpireSummary using JsonCodec
func (t *AmuletExpireSummary) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRules is a Template type
type AmuletRules struct {
	Dso            PARTY    `json:"dso"`
	ConfigSchedule Schedule `json:"configSchedule"`
	IsDevNet       BOOL     `json:"isDevNet"`
}

// GetTemplateID returns the template ID for this template
func (t AmuletRules) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletRules", "AmuletRules")
}

// CreateCommand returns a CreateCommand for this template
func (t AmuletRules) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["dso"] = t.Dso.ToMap()

	args["configSchedule"] = t.ConfigSchedule

	args["isDevNet"] = bool(t.IsDevNet)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for AmuletRules using JsonCodec
func (t AmuletRules) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRules using JsonCodec
func (t *AmuletRules) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for AmuletRules

// AmuletRulesComputeFees exercises the AmuletRules_ComputeFees choice on this AmuletRules contract
func (t AmuletRules) AmuletRulesComputeFees(contractID string, args AmuletRulesComputeFees) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletRules", "AmuletRules"),
		ContractID: contractID,
		Choice:     "AmuletRules_ComputeFees",
		Arguments:  argsToMap(args),
	}
}

// AmuletRulesTransfer exercises the AmuletRules_Transfer choice on this AmuletRules contract
func (t AmuletRules) AmuletRulesTransfer(contractID string, args AmuletRulesTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletRules", "AmuletRules"),
		ContractID: contractID,
		Choice:     "AmuletRules_Transfer",
		Arguments:  argsToMap(args),
	}
}

// AmuletRulesCreateExternalPartySetupProposal exercises the AmuletRules_CreateExternalPartySetupProposal choice on this AmuletRules contract
func (t AmuletRules) AmuletRulesCreateExternalPartySetupProposal(contractID string, args AmuletRulesCreateExternalPartySetupProposal) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletRules", "AmuletRules"),
		ContractID: contractID,
		Choice:     "AmuletRules_CreateExternalPartySetupProposal",
		Arguments:  argsToMap(args),
	}
}

// AmuletRulesCreateTransferPreapproval exercises the AmuletRules_CreateTransferPreapproval choice on this AmuletRules contract
func (t AmuletRules) AmuletRulesCreateTransferPreapproval(contractID string, args AmuletRulesCreateTransferPreapproval) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletRules", "AmuletRules"),
		ContractID: contractID,
		Choice:     "AmuletRules_CreateTransferPreapproval",
		Arguments:  argsToMap(args),
	}
}

// AmuletRulesBuyMemberTraffic exercises the AmuletRules_BuyMemberTraffic choice on this AmuletRules contract
func (t AmuletRules) AmuletRulesBuyMemberTraffic(contractID string, args AmuletRulesBuyMemberTraffic) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletRules", "AmuletRules"),
		ContractID: contractID,
		Choice:     "AmuletRules_BuyMemberTraffic",
		Arguments:  argsToMap(args),
	}
}

// AmuletRulesMergeMemberTrafficContracts exercises the AmuletRules_MergeMemberTrafficContracts choice on this AmuletRules contract
func (t AmuletRules) AmuletRulesMergeMemberTrafficContracts(contractID string, args AmuletRulesMergeMemberTrafficContracts) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletRules", "AmuletRules"),
		ContractID: contractID,
		Choice:     "AmuletRules_MergeMemberTrafficContracts",
		Arguments:  argsToMap(args),
	}
}

// AmuletRulesMint exercises the AmuletRules_Mint choice on this AmuletRules contract
func (t AmuletRules) AmuletRulesMint(contractID string, args AmuletRulesMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletRules", "AmuletRules"),
		ContractID: contractID,
		Choice:     "AmuletRules_Mint",
		Arguments:  argsToMap(args),
	}
}

// AmuletRulesDevNetTap exercises the AmuletRules_DevNet_Tap choice on this AmuletRules contract
func (t AmuletRules) AmuletRulesDevNetTap(contractID string, args AmuletRulesDevNetTap) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletRules", "AmuletRules"),
		ContractID: contractID,
		Choice:     "AmuletRules_DevNet_Tap",
		Arguments:  argsToMap(args),
	}
}

// AmuletRulesDevNetFeatureApp exercises the AmuletRules_DevNet_FeatureApp choice on this AmuletRules contract
func (t AmuletRules) AmuletRulesDevNetFeatureApp(contractID string, args AmuletRulesDevNetFeatureApp) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletRules", "AmuletRules"),
		ContractID: contractID,
		Choice:     "AmuletRules_DevNet_FeatureApp",
		Arguments:  argsToMap(args),
	}
}

// AmuletRulesBootstrapRounds exercises the AmuletRules_Bootstrap_Rounds choice on this AmuletRules contract
func (t AmuletRules) AmuletRulesBootstrapRounds(contractID string, args AmuletRulesBootstrapRounds) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletRules", "AmuletRules"),
		ContractID: contractID,
		Choice:     "AmuletRules_Bootstrap_Rounds",
		Arguments:  argsToMap(args),
	}
}

// AmuletRulesAdvanceOpenMiningRounds exercises the AmuletRules_AdvanceOpenMiningRounds choice on this AmuletRules contract
func (t AmuletRules) AmuletRulesAdvanceOpenMiningRounds(contractID string, args AmuletRulesAdvanceOpenMiningRounds) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletRules", "AmuletRules"),
		ContractID: contractID,
		Choice:     "AmuletRules_AdvanceOpenMiningRounds",
		Arguments:  argsToMap(args),
	}
}

// AmuletRulesMiningRoundStartIssuing exercises the AmuletRules_MiningRound_StartIssuing choice on this AmuletRules contract
func (t AmuletRules) AmuletRulesMiningRoundStartIssuing(contractID string, args AmuletRulesMiningRoundStartIssuing) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletRules", "AmuletRules"),
		ContractID: contractID,
		Choice:     "AmuletRules_MiningRound_StartIssuing",
		Arguments:  argsToMap(args),
	}
}

// AmuletRulesMiningRoundClose exercises the AmuletRules_MiningRound_Close choice on this AmuletRules contract
func (t AmuletRules) AmuletRulesMiningRoundClose(contractID string, args AmuletRulesMiningRoundClose) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletRules", "AmuletRules"),
		ContractID: contractID,
		Choice:     "AmuletRules_MiningRound_Close",
		Arguments:  argsToMap(args),
	}
}

// AmuletRulesMiningRoundArchive exercises the AmuletRules_MiningRound_Archive choice on this AmuletRules contract
func (t AmuletRules) AmuletRulesMiningRoundArchive(contractID string, args AmuletRulesMiningRoundArchive) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletRules", "AmuletRules"),
		ContractID: contractID,
		Choice:     "AmuletRules_MiningRound_Archive",
		Arguments:  argsToMap(args),
	}
}

// AmuletRulesClaimExpiredRewards exercises the AmuletRules_ClaimExpiredRewards choice on this AmuletRules contract
func (t AmuletRules) AmuletRulesClaimExpiredRewards(contractID string, args AmuletRulesClaimExpiredRewards) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletRules", "AmuletRules"),
		ContractID: contractID,
		Choice:     "AmuletRules_ClaimExpiredRewards",
		Arguments:  argsToMap(args),
	}
}

// AmuletRulesMergeUnclaimedRewards exercises the AmuletRules_MergeUnclaimedRewards choice on this AmuletRules contract
func (t AmuletRules) AmuletRulesMergeUnclaimedRewards(contractID string, args AmuletRulesMergeUnclaimedRewards) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletRules", "AmuletRules"),
		ContractID: contractID,
		Choice:     "AmuletRules_MergeUnclaimedRewards",
		Arguments:  argsToMap(args),
	}
}

// AmuletRulesMergeUnclaimedDevelopmentFundCoupons exercises the AmuletRules_MergeUnclaimedDevelopmentFundCoupons choice on this AmuletRules contract
func (t AmuletRules) AmuletRulesMergeUnclaimedDevelopmentFundCoupons(contractID string, args AmuletRulesMergeUnclaimedDevelopmentFundCoupons) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletRules", "AmuletRules"),
		ContractID: contractID,
		Choice:     "AmuletRules_MergeUnclaimedDevelopmentFundCoupons",
		Arguments:  argsToMap(args),
	}
}

// AmuletRulesAllocateDevelopmentFundCoupon exercises the AmuletRules_AllocateDevelopmentFundCoupon choice on this AmuletRules contract
func (t AmuletRules) AmuletRulesAllocateDevelopmentFundCoupon(contractID string, args AmuletRulesAllocateDevelopmentFundCoupon) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletRules", "AmuletRules"),
		ContractID: contractID,
		Choice:     "AmuletRules_AllocateDevelopmentFundCoupon",
		Arguments:  argsToMap(args),
	}
}

// AmuletRulesSetConfig exercises the AmuletRules_SetConfig choice on this AmuletRules contract
func (t AmuletRules) AmuletRulesSetConfig(contractID string, args SET) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletRules", "AmuletRules"),
		ContractID: contractID,
		Choice:     "AmuletRules_SetConfig",
		Arguments:  argsToMap(args),
	}
}

// AmuletRulesConvertFeaturedAppActivityMarkers exercises the AmuletRules_ConvertFeaturedAppActivityMarkers choice on this AmuletRules contract
func (t AmuletRules) AmuletRulesConvertFeaturedAppActivityMarkers(contractID string, args AmuletRulesConvertFeaturedAppActivityMarkers) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletRules", "AmuletRules"),
		ContractID: contractID,
		Choice:     "AmuletRules_ConvertFeaturedAppActivityMarkers",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this AmuletRules contract
func (t AmuletRules) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletRules", "AmuletRules"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// AmuletRulesFetch exercises the AmuletRules_Fetch choice on this AmuletRules contract
func (t AmuletRules) AmuletRulesFetch(contractID string, args AmuletRulesFetch) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletRules", "AmuletRules"),
		ContractID: contractID,
		Choice:     "AmuletRules_Fetch",
		Arguments:  argsToMap(args),
	}
}

// AmuletRulesAddFutureAmuletConfigSchedule exercises the AmuletRules_AddFutureAmuletConfigSchedule choice on this AmuletRules contract
func (t AmuletRules) AmuletRulesAddFutureAmuletConfigSchedule(contractID string, args AmuletRulesAddFutureAmuletConfigSchedule) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletRules", "AmuletRules"),
		ContractID: contractID,
		Choice:     "AmuletRules_AddFutureAmuletConfigSchedule",
		Arguments:  argsToMap(args),
	}
}

// AmuletRulesRemoveFutureAmuletConfigSchedule exercises the AmuletRules_RemoveFutureAmuletConfigSchedule choice on this AmuletRules contract
func (t AmuletRules) AmuletRulesRemoveFutureAmuletConfigSchedule(contractID string, args AmuletRulesRemoveFutureAmuletConfigSchedule) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletRules", "AmuletRules"),
		ContractID: contractID,
		Choice:     "AmuletRules_RemoveFutureAmuletConfigSchedule",
		Arguments:  argsToMap(args),
	}
}

// AmuletRulesUpdateFutureAmuletConfigSchedule exercises the AmuletRules_UpdateFutureAmuletConfigSchedule choice on this AmuletRules contract
func (t AmuletRules) AmuletRulesUpdateFutureAmuletConfigSchedule(contractID string, args AmuletRulesUpdateFutureAmuletConfigSchedule) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletRules", "AmuletRules"),
		ContractID: contractID,
		Choice:     "AmuletRules_UpdateFutureAmuletConfigSchedule",
		Arguments:  argsToMap(args),
	}
}

// AmuletRulesAddFutureAmuletConfigSchedule is a Record type
type AmuletRulesAddFutureAmuletConfigSchedule struct {
	NewScheduleItem TUPLE2[TIMESTAMP, AmuletConfig] `json:"newScheduleItem"`
}

// ToMap converts AmuletRulesAddFutureAmuletConfigSchedule to a map for DAML arguments
func (t AmuletRulesAddFutureAmuletConfigSchedule) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["newScheduleItem"] = t.NewScheduleItem

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesAddFutureAmuletConfigSchedule using JsonCodec
func (t AmuletRulesAddFutureAmuletConfigSchedule) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesAddFutureAmuletConfigSchedule using JsonCodec
func (t *AmuletRulesAddFutureAmuletConfigSchedule) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesAddFutureAmuletConfigScheduleResult is a Record type
type AmuletRulesAddFutureAmuletConfigScheduleResult struct {
	NewAmuletRules CONTRACT_ID `json:"newAmuletRules"`
}

// ToMap converts AmuletRulesAddFutureAmuletConfigScheduleResult to a map for DAML arguments
func (t AmuletRulesAddFutureAmuletConfigScheduleResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["newAmuletRules"] = t.NewAmuletRules

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesAddFutureAmuletConfigScheduleResult using JsonCodec
func (t AmuletRulesAddFutureAmuletConfigScheduleResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesAddFutureAmuletConfigScheduleResult using JsonCodec
func (t *AmuletRulesAddFutureAmuletConfigScheduleResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesAdvanceOpenMiningRounds is a Record type
type AmuletRulesAdvanceOpenMiningRounds struct {
	AmuletPrice       NUMERIC     `json:"amuletPrice"`
	RoundToArchiveCid CONTRACT_ID `json:"roundToArchiveCid"`
	MiddleRoundCid    CONTRACT_ID `json:"middleRoundCid"`
	LatestRoundCid    CONTRACT_ID `json:"latestRoundCid"`
}

// ToMap converts AmuletRulesAdvanceOpenMiningRounds to a map for DAML arguments
func (t AmuletRulesAdvanceOpenMiningRounds) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["amuletPrice"] = (*big.Int)(t.AmuletPrice)

	m["roundToArchiveCid"] = t.RoundToArchiveCid

	m["middleRoundCid"] = t.MiddleRoundCid

	m["latestRoundCid"] = t.LatestRoundCid

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesAdvanceOpenMiningRounds using JsonCodec
func (t AmuletRulesAdvanceOpenMiningRounds) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesAdvanceOpenMiningRounds using JsonCodec
func (t *AmuletRulesAdvanceOpenMiningRounds) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesAdvanceOpenMiningRoundsResult is a Record type
type AmuletRulesAdvanceOpenMiningRoundsResult struct {
	SummarizingRoundCid CONTRACT_ID `json:"summarizingRoundCid"`
	OpenRoundCid        CONTRACT_ID `json:"openRoundCid"`
}

// ToMap converts AmuletRulesAdvanceOpenMiningRoundsResult to a map for DAML arguments
func (t AmuletRulesAdvanceOpenMiningRoundsResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["summarizingRoundCid"] = t.SummarizingRoundCid

	m["openRoundCid"] = t.OpenRoundCid

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesAdvanceOpenMiningRoundsResult using JsonCodec
func (t AmuletRulesAdvanceOpenMiningRoundsResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesAdvanceOpenMiningRoundsResult using JsonCodec
func (t *AmuletRulesAdvanceOpenMiningRoundsResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesAllocateDevelopmentFundCoupon is a Record type
type AmuletRulesAllocateDevelopmentFundCoupon struct {
	UnclaimedDevelopmentFundCouponCids []CONTRACT_ID `json:"unclaimedDevelopmentFundCouponCids"`
	Beneficiary                        PARTY         `json:"beneficiary"`
	Amount                             NUMERIC       `json:"amount"`
	ExpiresAt                          TIMESTAMP     `json:"expiresAt"`
	Reason                             TEXT          `json:"reason"`
	FundManager                        PARTY         `json:"fundManager"`
}

// ToMap converts AmuletRulesAllocateDevelopmentFundCoupon to a map for DAML arguments
func (t AmuletRulesAllocateDevelopmentFundCoupon) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["unclaimedDevelopmentFundCouponCids"] = t.UnclaimedDevelopmentFundCouponCids

	m["beneficiary"] = t.Beneficiary.ToMap()

	m["amount"] = (*big.Int)(t.Amount)

	m["expiresAt"] = t.ExpiresAt

	m["reason"] = string(t.Reason)

	m["fundManager"] = t.FundManager.ToMap()

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesAllocateDevelopmentFundCoupon using JsonCodec
func (t AmuletRulesAllocateDevelopmentFundCoupon) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesAllocateDevelopmentFundCoupon using JsonCodec
func (t *AmuletRulesAllocateDevelopmentFundCoupon) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesAllocateDevelopmentFundCouponResult is a Record type
type AmuletRulesAllocateDevelopmentFundCouponResult struct {
	DevelopmentFundCouponCid             CONTRACT_ID  `json:"developmentFundCouponCid"`
	OptUnclaimedDevelopmentFundCouponCid *CONTRACT_ID `json:"optUnclaimedDevelopmentFundCouponCid"`
}

// ToMap converts AmuletRulesAllocateDevelopmentFundCouponResult to a map for DAML arguments
func (t AmuletRulesAllocateDevelopmentFundCouponResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["developmentFundCouponCid"] = t.DevelopmentFundCouponCid

	if t.OptUnclaimedDevelopmentFundCouponCid != nil {
		m["optUnclaimedDevelopmentFundCouponCid"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.OptUnclaimedDevelopmentFundCouponCid,
		}
	} else {
		m["optUnclaimedDevelopmentFundCouponCid"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesAllocateDevelopmentFundCouponResult using JsonCodec
func (t AmuletRulesAllocateDevelopmentFundCouponResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesAllocateDevelopmentFundCouponResult using JsonCodec
func (t *AmuletRulesAllocateDevelopmentFundCouponResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesBootstrapRounds is a Record type
type AmuletRulesBootstrapRounds struct {
	AmuletPrice    NUMERIC `json:"amuletPrice"`
	Round0Duration RELTIME `json:"round0Duration"`
	InitialRound   *INT64  `json:"initialRound"`
}

// ToMap converts AmuletRulesBootstrapRounds to a map for DAML arguments
func (t AmuletRulesBootstrapRounds) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["amuletPrice"] = (*big.Int)(t.AmuletPrice)

	m["round0Duration"] = t.Round0Duration

	if t.InitialRound != nil {
		m["initialRound"] = map[string]interface{}{
			"_type": "optional",
			"value": int64(*t.InitialRound),
		}
	} else {
		m["initialRound"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesBootstrapRounds using JsonCodec
func (t AmuletRulesBootstrapRounds) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesBootstrapRounds using JsonCodec
func (t *AmuletRulesBootstrapRounds) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesBootstrapRoundsResult is a Record type
type AmuletRulesBootstrapRoundsResult struct {
	OpenMiningRoundCid CONTRACT_ID `json:"openMiningRoundCid"`
	InitialRound       *Round      `json:"initialRound"`
}

// ToMap converts AmuletRulesBootstrapRoundsResult to a map for DAML arguments
func (t AmuletRulesBootstrapRoundsResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["openMiningRoundCid"] = t.OpenMiningRoundCid

	if t.InitialRound != nil {
		m["initialRound"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.InitialRound,
		}
	} else {
		m["initialRound"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesBootstrapRoundsResult using JsonCodec
func (t AmuletRulesBootstrapRoundsResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesBootstrapRoundsResult using JsonCodec
func (t *AmuletRulesBootstrapRoundsResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesBuyMemberTraffic is a Record type
type AmuletRulesBuyMemberTraffic struct {
	Inputs         []TransferInput `json:"inputs"`
	Context        TransferContext `json:"context"`
	Provider       PARTY           `json:"provider"`
	MemberId       TEXT            `json:"memberId"`
	SynchronizerId TEXT            `json:"synchronizerId"`
	MigrationId    INT64           `json:"migrationId"`
	TrafficAmount  INT64           `json:"trafficAmount"`
	ExpectedDso    *PARTY          `json:"expectedDso"`
}

// ToMap converts AmuletRulesBuyMemberTraffic to a map for DAML arguments
func (t AmuletRulesBuyMemberTraffic) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["inputs"] = t.Inputs

	m["context"] = t.Context

	m["provider"] = t.Provider.ToMap()

	m["memberId"] = string(t.MemberId)

	m["synchronizerId"] = string(t.SynchronizerId)

	m["migrationId"] = int64(t.MigrationId)

	m["trafficAmount"] = int64(t.TrafficAmount)

	if t.ExpectedDso != nil {
		m["expectedDso"] = map[string]interface{}{
			"_type": "optional",
			"value": (*t.ExpectedDso).ToMap(),
		}
	} else {
		m["expectedDso"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesBuyMemberTraffic using JsonCodec
func (t AmuletRulesBuyMemberTraffic) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesBuyMemberTraffic using JsonCodec
func (t *AmuletRulesBuyMemberTraffic) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesBuyMemberTrafficResult is a Record type
type AmuletRulesBuyMemberTrafficResult struct {
	Round              Round           `json:"round"`
	Summary            TransferSummary `json:"summary"`
	AmuletPaid         NUMERIC         `json:"amuletPaid"`
	PurchasedTraffic   CONTRACT_ID     `json:"purchasedTraffic"`
	SenderChangeAmulet *CONTRACT_ID    `json:"senderChangeAmulet"`
	Meta               *Metadata       `json:"meta"`
}

// ToMap converts AmuletRulesBuyMemberTrafficResult to a map for DAML arguments
func (t AmuletRulesBuyMemberTrafficResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["round"] = t.Round

	m["summary"] = t.Summary

	m["amuletPaid"] = (*big.Int)(t.AmuletPaid)

	m["purchasedTraffic"] = t.PurchasedTraffic

	if t.SenderChangeAmulet != nil {
		m["senderChangeAmulet"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.SenderChangeAmulet,
		}
	} else {
		m["senderChangeAmulet"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	if t.Meta != nil {
		m["meta"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.Meta,
		}
	} else {
		m["meta"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesBuyMemberTrafficResult using JsonCodec
func (t AmuletRulesBuyMemberTrafficResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesBuyMemberTrafficResult using JsonCodec
func (t *AmuletRulesBuyMemberTrafficResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesClaimExpiredRewards is a Record type
type AmuletRulesClaimExpiredRewards struct {
	ClosedRoundCid                         CONTRACT_ID    `json:"closedRoundCid"`
	ValidatorRewardCouponCids              []CONTRACT_ID  `json:"validatorRewardCouponCids"`
	AppCouponCids                          []CONTRACT_ID  `json:"appCouponCids"`
	SvRewardCouponCids                     []CONTRACT_ID  `json:"svRewardCouponCids"`
	OptValidatorFaucetCouponCids           *[]CONTRACT_ID `json:"optValidatorFaucetCouponCids"`
	OptValidatorLivenessActivityRecordCids *[]CONTRACT_ID `json:"optValidatorLivenessActivityRecordCids"`
}

// ToMap converts AmuletRulesClaimExpiredRewards to a map for DAML arguments
func (t AmuletRulesClaimExpiredRewards) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["closedRoundCid"] = t.ClosedRoundCid

	m["validatorRewardCouponCids"] = t.ValidatorRewardCouponCids

	m["appCouponCids"] = t.AppCouponCids

	m["svRewardCouponCids"] = t.SvRewardCouponCids

	if t.OptValidatorFaucetCouponCids != nil {
		m["optValidatorFaucetCouponCids"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.OptValidatorFaucetCouponCids,
		}
	} else {
		m["optValidatorFaucetCouponCids"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	if t.OptValidatorLivenessActivityRecordCids != nil {
		m["optValidatorLivenessActivityRecordCids"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.OptValidatorLivenessActivityRecordCids,
		}
	} else {
		m["optValidatorLivenessActivityRecordCids"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesClaimExpiredRewards using JsonCodec
func (t AmuletRulesClaimExpiredRewards) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesClaimExpiredRewards using JsonCodec
func (t *AmuletRulesClaimExpiredRewards) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesClaimExpiredRewardsResult is a Record type
type AmuletRulesClaimExpiredRewardsResult struct {
	UnclaimedRewardCid *CONTRACT_ID `json:"unclaimedRewardCid"`
}

// ToMap converts AmuletRulesClaimExpiredRewardsResult to a map for DAML arguments
func (t AmuletRulesClaimExpiredRewardsResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	if t.UnclaimedRewardCid != nil {
		m["unclaimedRewardCid"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.UnclaimedRewardCid,
		}
	} else {
		m["unclaimedRewardCid"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesClaimExpiredRewardsResult using JsonCodec
func (t AmuletRulesClaimExpiredRewardsResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesClaimExpiredRewardsResult using JsonCodec
func (t *AmuletRulesClaimExpiredRewardsResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesComputeFees is a Record type
type AmuletRulesComputeFees struct {
	Context     TransferContext  `json:"context"`
	Sender      PARTY            `json:"sender"`
	Outputs     []TransferOutput `json:"outputs"`
	ExpectedDso *PARTY           `json:"expectedDso"`
}

// ToMap converts AmuletRulesComputeFees to a map for DAML arguments
func (t AmuletRulesComputeFees) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["context"] = t.Context

	m["sender"] = t.Sender.ToMap()

	m["outputs"] = t.Outputs

	if t.ExpectedDso != nil {
		m["expectedDso"] = map[string]interface{}{
			"_type": "optional",
			"value": (*t.ExpectedDso).ToMap(),
		}
	} else {
		m["expectedDso"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesComputeFees using JsonCodec
func (t AmuletRulesComputeFees) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesComputeFees using JsonCodec
func (t *AmuletRulesComputeFees) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesComputeFeesResult is a Record type
type AmuletRulesComputeFeesResult struct {
	Fees []NUMERIC `json:"fees"`
}

// ToMap converts AmuletRulesComputeFeesResult to a map for DAML arguments
func (t AmuletRulesComputeFeesResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["fees"] = t.Fees

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesComputeFeesResult using JsonCodec
func (t AmuletRulesComputeFeesResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesComputeFeesResult using JsonCodec
func (t *AmuletRulesComputeFeesResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesConvertFeaturedAppActivityMarkers is a Record type
type AmuletRulesConvertFeaturedAppActivityMarkers struct {
	MarkerCids         []CONTRACT_ID `json:"markerCids"`
	OpenMiningRoundCid CONTRACT_ID   `json:"openMiningRoundCid"`
	Observers          *[]PARTY      `json:"observers"`
}

// ToMap converts AmuletRulesConvertFeaturedAppActivityMarkers to a map for DAML arguments
func (t AmuletRulesConvertFeaturedAppActivityMarkers) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["markerCids"] = t.MarkerCids

	m["openMiningRoundCid"] = t.OpenMiningRoundCid

	if t.Observers != nil {
		m["observers"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.Observers,
		}
	} else {
		m["observers"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesConvertFeaturedAppActivityMarkers using JsonCodec
func (t AmuletRulesConvertFeaturedAppActivityMarkers) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesConvertFeaturedAppActivityMarkers using JsonCodec
func (t *AmuletRulesConvertFeaturedAppActivityMarkers) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesConvertFeaturedAppActivityMarkersResult is a Record type
type AmuletRulesConvertFeaturedAppActivityMarkersResult struct {
	AppRewardCouponCids []CONTRACT_ID `json:"appRewardCouponCids"`
}

// ToMap converts AmuletRulesConvertFeaturedAppActivityMarkersResult to a map for DAML arguments
func (t AmuletRulesConvertFeaturedAppActivityMarkersResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["appRewardCouponCids"] = t.AppRewardCouponCids

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesConvertFeaturedAppActivityMarkersResult using JsonCodec
func (t AmuletRulesConvertFeaturedAppActivityMarkersResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesConvertFeaturedAppActivityMarkersResult using JsonCodec
func (t *AmuletRulesConvertFeaturedAppActivityMarkersResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesCreateExternalPartySetupProposal is a Record type
type AmuletRulesCreateExternalPartySetupProposal struct {
	Context              PaymentTransferContext `json:"context"`
	Inputs               []TransferInput        `json:"inputs"`
	User                 PARTY                  `json:"user"`
	Validator            PARTY                  `json:"validator"`
	PreapprovalExpiresAt TIMESTAMP              `json:"preapprovalExpiresAt"`
	ExpectedDso          *PARTY                 `json:"expectedDso"`
}

// ToMap converts AmuletRulesCreateExternalPartySetupProposal to a map for DAML arguments
func (t AmuletRulesCreateExternalPartySetupProposal) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["context"] = t.Context

	m["inputs"] = t.Inputs

	m["user"] = t.User.ToMap()

	m["validator"] = t.Validator.ToMap()

	m["preapprovalExpiresAt"] = t.PreapprovalExpiresAt

	if t.ExpectedDso != nil {
		m["expectedDso"] = map[string]interface{}{
			"_type": "optional",
			"value": (*t.ExpectedDso).ToMap(),
		}
	} else {
		m["expectedDso"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesCreateExternalPartySetupProposal using JsonCodec
func (t AmuletRulesCreateExternalPartySetupProposal) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesCreateExternalPartySetupProposal using JsonCodec
func (t *AmuletRulesCreateExternalPartySetupProposal) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesCreateExternalPartySetupProposalResult is a Record type
type AmuletRulesCreateExternalPartySetupProposalResult struct {
	ProposalCid    CONTRACT_ID    `json:"proposalCid"`
	User           PARTY          `json:"user"`
	Validator      PARTY          `json:"validator"`
	TransferResult TransferResult `json:"transferResult"`
	AmuletPaid     NUMERIC        `json:"amuletPaid"`
	Meta           *Metadata      `json:"meta"`
}

// ToMap converts AmuletRulesCreateExternalPartySetupProposalResult to a map for DAML arguments
func (t AmuletRulesCreateExternalPartySetupProposalResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["proposalCid"] = t.ProposalCid

	m["user"] = t.User.ToMap()

	m["validator"] = t.Validator.ToMap()

	m["transferResult"] = t.TransferResult

	m["amuletPaid"] = (*big.Int)(t.AmuletPaid)

	if t.Meta != nil {
		m["meta"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.Meta,
		}
	} else {
		m["meta"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesCreateExternalPartySetupProposalResult using JsonCodec
func (t AmuletRulesCreateExternalPartySetupProposalResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesCreateExternalPartySetupProposalResult using JsonCodec
func (t *AmuletRulesCreateExternalPartySetupProposalResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesCreateTransferPreapproval is a Record type
type AmuletRulesCreateTransferPreapproval struct {
	Context     PaymentTransferContext `json:"context"`
	Inputs      []TransferInput        `json:"inputs"`
	Receiver    PARTY                  `json:"receiver"`
	Provider    PARTY                  `json:"provider"`
	ExpiresAt   TIMESTAMP              `json:"expiresAt"`
	ExpectedDso *PARTY                 `json:"expectedDso"`
}

// ToMap converts AmuletRulesCreateTransferPreapproval to a map for DAML arguments
func (t AmuletRulesCreateTransferPreapproval) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["context"] = t.Context

	m["inputs"] = t.Inputs

	m["receiver"] = t.Receiver.ToMap()

	m["provider"] = t.Provider.ToMap()

	m["expiresAt"] = t.ExpiresAt

	if t.ExpectedDso != nil {
		m["expectedDso"] = map[string]interface{}{
			"_type": "optional",
			"value": (*t.ExpectedDso).ToMap(),
		}
	} else {
		m["expectedDso"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesCreateTransferPreapproval using JsonCodec
func (t AmuletRulesCreateTransferPreapproval) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesCreateTransferPreapproval using JsonCodec
func (t *AmuletRulesCreateTransferPreapproval) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesCreateTransferPreapprovalResult is a Record type
type AmuletRulesCreateTransferPreapprovalResult struct {
	TransferPreapprovalCid CONTRACT_ID    `json:"transferPreapprovalCid"`
	TransferResult         TransferResult `json:"transferResult"`
	AmuletPaid             NUMERIC        `json:"amuletPaid"`
	Meta                   *Metadata      `json:"meta"`
}

// ToMap converts AmuletRulesCreateTransferPreapprovalResult to a map for DAML arguments
func (t AmuletRulesCreateTransferPreapprovalResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["transferPreapprovalCid"] = t.TransferPreapprovalCid

	m["transferResult"] = t.TransferResult

	m["amuletPaid"] = (*big.Int)(t.AmuletPaid)

	if t.Meta != nil {
		m["meta"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.Meta,
		}
	} else {
		m["meta"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesCreateTransferPreapprovalResult using JsonCodec
func (t AmuletRulesCreateTransferPreapprovalResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesCreateTransferPreapprovalResult using JsonCodec
func (t *AmuletRulesCreateTransferPreapprovalResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesDevNetFeatureApp is a Record type
type AmuletRulesDevNetFeatureApp struct {
	Provider PARTY `json:"provider"`
}

// ToMap converts AmuletRulesDevNetFeatureApp to a map for DAML arguments
func (t AmuletRulesDevNetFeatureApp) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["provider"] = t.Provider.ToMap()

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesDevNetFeatureApp using JsonCodec
func (t AmuletRulesDevNetFeatureApp) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesDevNetFeatureApp using JsonCodec
func (t *AmuletRulesDevNetFeatureApp) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesDevNetFeatureAppResult is a Record type
type AmuletRulesDevNetFeatureAppResult struct {
	FeaturedAppRightCid CONTRACT_ID `json:"featuredAppRightCid"`
}

// ToMap converts AmuletRulesDevNetFeatureAppResult to a map for DAML arguments
func (t AmuletRulesDevNetFeatureAppResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["featuredAppRightCid"] = t.FeaturedAppRightCid

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesDevNetFeatureAppResult using JsonCodec
func (t AmuletRulesDevNetFeatureAppResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesDevNetFeatureAppResult using JsonCodec
func (t *AmuletRulesDevNetFeatureAppResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesDevNetTap is a Record type
type AmuletRulesDevNetTap struct {
	Receiver  PARTY       `json:"receiver"`
	Amount    NUMERIC     `json:"amount"`
	OpenRound CONTRACT_ID `json:"openRound"`
}

// ToMap converts AmuletRulesDevNetTap to a map for DAML arguments
func (t AmuletRulesDevNetTap) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["receiver"] = t.Receiver.ToMap()

	m["amount"] = (*big.Int)(t.Amount)

	m["openRound"] = t.OpenRound

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesDevNetTap using JsonCodec
func (t AmuletRulesDevNetTap) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesDevNetTap using JsonCodec
func (t *AmuletRulesDevNetTap) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesDevNetTapResult is a Record type
type AmuletRulesDevNetTapResult struct {
	AmuletSum AmuletCreateSummary `json:"amuletSum"`
	Meta      *Metadata           `json:"meta"`
}

// ToMap converts AmuletRulesDevNetTapResult to a map for DAML arguments
func (t AmuletRulesDevNetTapResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["amuletSum"] = t.AmuletSum

	if t.Meta != nil {
		m["meta"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.Meta,
		}
	} else {
		m["meta"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesDevNetTapResult using JsonCodec
func (t AmuletRulesDevNetTapResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesDevNetTapResult using JsonCodec
func (t *AmuletRulesDevNetTapResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesFetch is a Record type
type AmuletRulesFetch struct {
	P PARTY `json:"p"`
}

// ToMap converts AmuletRulesFetch to a map for DAML arguments
func (t AmuletRulesFetch) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["p"] = t.P.ToMap()

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesFetch using JsonCodec
func (t AmuletRulesFetch) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesFetch using JsonCodec
func (t *AmuletRulesFetch) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesMergeMemberTrafficContracts is a Record type
type AmuletRulesMergeMemberTrafficContracts struct {
	TrafficCids []CONTRACT_ID `json:"trafficCids"`
}

// ToMap converts AmuletRulesMergeMemberTrafficContracts to a map for DAML arguments
func (t AmuletRulesMergeMemberTrafficContracts) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["trafficCids"] = t.TrafficCids

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesMergeMemberTrafficContracts using JsonCodec
func (t AmuletRulesMergeMemberTrafficContracts) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesMergeMemberTrafficContracts using JsonCodec
func (t *AmuletRulesMergeMemberTrafficContracts) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesMergeMemberTrafficContractsResult is a Record type
type AmuletRulesMergeMemberTrafficContractsResult struct {
	MergedTrafficCid CONTRACT_ID `json:"mergedTrafficCid"`
}

// ToMap converts AmuletRulesMergeMemberTrafficContractsResult to a map for DAML arguments
func (t AmuletRulesMergeMemberTrafficContractsResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["mergedTrafficCid"] = t.MergedTrafficCid

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesMergeMemberTrafficContractsResult using JsonCodec
func (t AmuletRulesMergeMemberTrafficContractsResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesMergeMemberTrafficContractsResult using JsonCodec
func (t *AmuletRulesMergeMemberTrafficContractsResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesMergeUnclaimedDevelopmentFundCoupons is a Record type
type AmuletRulesMergeUnclaimedDevelopmentFundCoupons struct {
	UnclaimedDevelopmentFundCouponCids []CONTRACT_ID `json:"unclaimedDevelopmentFundCouponCids"`
}

// ToMap converts AmuletRulesMergeUnclaimedDevelopmentFundCoupons to a map for DAML arguments
func (t AmuletRulesMergeUnclaimedDevelopmentFundCoupons) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["unclaimedDevelopmentFundCouponCids"] = t.UnclaimedDevelopmentFundCouponCids

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesMergeUnclaimedDevelopmentFundCoupons using JsonCodec
func (t AmuletRulesMergeUnclaimedDevelopmentFundCoupons) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesMergeUnclaimedDevelopmentFundCoupons using JsonCodec
func (t *AmuletRulesMergeUnclaimedDevelopmentFundCoupons) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesMergeUnclaimedDevelopmentFundCouponsResult is a Record type
type AmuletRulesMergeUnclaimedDevelopmentFundCouponsResult struct {
	UnclaimedDevelopmentFundCouponCid CONTRACT_ID `json:"unclaimedDevelopmentFundCouponCid"`
}

// ToMap converts AmuletRulesMergeUnclaimedDevelopmentFundCouponsResult to a map for DAML arguments
func (t AmuletRulesMergeUnclaimedDevelopmentFundCouponsResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["unclaimedDevelopmentFundCouponCid"] = t.UnclaimedDevelopmentFundCouponCid

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesMergeUnclaimedDevelopmentFundCouponsResult using JsonCodec
func (t AmuletRulesMergeUnclaimedDevelopmentFundCouponsResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesMergeUnclaimedDevelopmentFundCouponsResult using JsonCodec
func (t *AmuletRulesMergeUnclaimedDevelopmentFundCouponsResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesMergeUnclaimedRewards is a Record type
type AmuletRulesMergeUnclaimedRewards struct {
	UnclaimedRewardCids []CONTRACT_ID `json:"unclaimedRewardCids"`
}

// ToMap converts AmuletRulesMergeUnclaimedRewards to a map for DAML arguments
func (t AmuletRulesMergeUnclaimedRewards) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["unclaimedRewardCids"] = t.UnclaimedRewardCids

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesMergeUnclaimedRewards using JsonCodec
func (t AmuletRulesMergeUnclaimedRewards) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesMergeUnclaimedRewards using JsonCodec
func (t *AmuletRulesMergeUnclaimedRewards) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesMergeUnclaimedRewardsResult is a Record type
type AmuletRulesMergeUnclaimedRewardsResult struct {
	UnclaimedRewardCid CONTRACT_ID `json:"unclaimedRewardCid"`
}

// ToMap converts AmuletRulesMergeUnclaimedRewardsResult to a map for DAML arguments
func (t AmuletRulesMergeUnclaimedRewardsResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["unclaimedRewardCid"] = t.UnclaimedRewardCid

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesMergeUnclaimedRewardsResult using JsonCodec
func (t AmuletRulesMergeUnclaimedRewardsResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesMergeUnclaimedRewardsResult using JsonCodec
func (t *AmuletRulesMergeUnclaimedRewardsResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesMiningRoundArchive is a Record type
type AmuletRulesMiningRoundArchive struct {
	ClosedRoundCid CONTRACT_ID `json:"closedRoundCid"`
}

// ToMap converts AmuletRulesMiningRoundArchive to a map for DAML arguments
func (t AmuletRulesMiningRoundArchive) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["closedRoundCid"] = t.ClosedRoundCid

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesMiningRoundArchive using JsonCodec
func (t AmuletRulesMiningRoundArchive) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesMiningRoundArchive using JsonCodec
func (t *AmuletRulesMiningRoundArchive) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesMiningRoundArchiveResult is an enum type
type AmuletRulesMiningRoundArchiveResult string

const (
	AmuletRulesMiningRoundArchiveResultAmuletRules_MiningRound_ArchiveResult AmuletRulesMiningRoundArchiveResult = "AmuletRules_MiningRound_ArchiveResult"
)

// GetEnumConstructor implements types.ENUM interface
func (e AmuletRulesMiningRoundArchiveResult) GetEnumConstructor() string {
	return string(e)
}

// GetEnumTypeID implements types.ENUM interface
func (e AmuletRulesMiningRoundArchiveResult) GetEnumTypeID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletRules", "AmuletRulesMiningRoundArchiveResult")
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesMiningRoundArchiveResult using JsonCodec
func (e AmuletRulesMiningRoundArchiveResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(e)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesMiningRoundArchiveResult using JsonCodec
func (e *AmuletRulesMiningRoundArchiveResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, e)
}

// Verify interface implementation
var _ ENUM = AmuletRulesMiningRoundArchiveResult("")

// AmuletRulesMiningRoundClose is a Record type
type AmuletRulesMiningRoundClose struct {
	IssuingRoundCid CONTRACT_ID `json:"issuingRoundCid"`
}

// ToMap converts AmuletRulesMiningRoundClose to a map for DAML arguments
func (t AmuletRulesMiningRoundClose) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["issuingRoundCid"] = t.IssuingRoundCid

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesMiningRoundClose using JsonCodec
func (t AmuletRulesMiningRoundClose) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesMiningRoundClose using JsonCodec
func (t *AmuletRulesMiningRoundClose) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesMiningRoundCloseResult is a Record type
type AmuletRulesMiningRoundCloseResult struct {
	ClosedRoundCid CONTRACT_ID `json:"closedRoundCid"`
}

// ToMap converts AmuletRulesMiningRoundCloseResult to a map for DAML arguments
func (t AmuletRulesMiningRoundCloseResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["closedRoundCid"] = t.ClosedRoundCid

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesMiningRoundCloseResult using JsonCodec
func (t AmuletRulesMiningRoundCloseResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesMiningRoundCloseResult using JsonCodec
func (t *AmuletRulesMiningRoundCloseResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesMiningRoundStartIssuing is a Record type
type AmuletRulesMiningRoundStartIssuing struct {
	MiningRoundCid CONTRACT_ID            `json:"miningRoundCid"`
	Summary        OpenMiningRoundSummary `json:"summary"`
}

// ToMap converts AmuletRulesMiningRoundStartIssuing to a map for DAML arguments
func (t AmuletRulesMiningRoundStartIssuing) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["miningRoundCid"] = t.MiningRoundCid

	m["summary"] = t.Summary

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesMiningRoundStartIssuing using JsonCodec
func (t AmuletRulesMiningRoundStartIssuing) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesMiningRoundStartIssuing using JsonCodec
func (t *AmuletRulesMiningRoundStartIssuing) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesMiningRoundStartIssuingResult is a Record type
type AmuletRulesMiningRoundStartIssuingResult struct {
	IssuingRoundCid                   CONTRACT_ID  `json:"issuingRoundCid"`
	UnclaimedDevelopmentFundCouponCid *CONTRACT_ID `json:"unclaimedDevelopmentFundCouponCid"`
}

// ToMap converts AmuletRulesMiningRoundStartIssuingResult to a map for DAML arguments
func (t AmuletRulesMiningRoundStartIssuingResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["issuingRoundCid"] = t.IssuingRoundCid

	if t.UnclaimedDevelopmentFundCouponCid != nil {
		m["unclaimedDevelopmentFundCouponCid"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.UnclaimedDevelopmentFundCouponCid,
		}
	} else {
		m["unclaimedDevelopmentFundCouponCid"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesMiningRoundStartIssuingResult using JsonCodec
func (t AmuletRulesMiningRoundStartIssuingResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesMiningRoundStartIssuingResult using JsonCodec
func (t *AmuletRulesMiningRoundStartIssuingResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesMint is a Record type
type AmuletRulesMint struct {
	Receiver  PARTY       `json:"receiver"`
	Amount    NUMERIC     `json:"amount"`
	OpenRound CONTRACT_ID `json:"openRound"`
}

// ToMap converts AmuletRulesMint to a map for DAML arguments
func (t AmuletRulesMint) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["receiver"] = t.Receiver.ToMap()

	m["amount"] = (*big.Int)(t.Amount)

	m["openRound"] = t.OpenRound

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesMint using JsonCodec
func (t AmuletRulesMint) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesMint using JsonCodec
func (t *AmuletRulesMint) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesMintResult is a Record type
type AmuletRulesMintResult struct {
	AmuletSum AmuletCreateSummary `json:"amuletSum"`
}

// ToMap converts AmuletRulesMintResult to a map for DAML arguments
func (t AmuletRulesMintResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["amuletSum"] = t.AmuletSum

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesMintResult using JsonCodec
func (t AmuletRulesMintResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesMintResult using JsonCodec
func (t *AmuletRulesMintResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesRemoveFutureAmuletConfigSchedule is a Record type
type AmuletRulesRemoveFutureAmuletConfigSchedule struct {
	ScheduleTime TIMESTAMP `json:"scheduleTime"`
}

// ToMap converts AmuletRulesRemoveFutureAmuletConfigSchedule to a map for DAML arguments
func (t AmuletRulesRemoveFutureAmuletConfigSchedule) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["scheduleTime"] = t.ScheduleTime

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesRemoveFutureAmuletConfigSchedule using JsonCodec
func (t AmuletRulesRemoveFutureAmuletConfigSchedule) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesRemoveFutureAmuletConfigSchedule using JsonCodec
func (t *AmuletRulesRemoveFutureAmuletConfigSchedule) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesRemoveFutureAmuletConfigScheduleResult is a Record type
type AmuletRulesRemoveFutureAmuletConfigScheduleResult struct {
	NewAmuletRules CONTRACT_ID `json:"newAmuletRules"`
}

// ToMap converts AmuletRulesRemoveFutureAmuletConfigScheduleResult to a map for DAML arguments
func (t AmuletRulesRemoveFutureAmuletConfigScheduleResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["newAmuletRules"] = t.NewAmuletRules

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesRemoveFutureAmuletConfigScheduleResult using JsonCodec
func (t AmuletRulesRemoveFutureAmuletConfigScheduleResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesRemoveFutureAmuletConfigScheduleResult using JsonCodec
func (t *AmuletRulesRemoveFutureAmuletConfigScheduleResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesSetConfig is a Record type
type AmuletRulesSetConfig struct {
	NewConfig  AmuletConfig `json:"newConfig"`
	BaseConfig AmuletConfig `json:"baseConfig"`
}

// ToMap converts AmuletRulesSetConfig to a map for DAML arguments
func (t AmuletRulesSetConfig) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["newConfig"] = t.NewConfig

	m["baseConfig"] = t.BaseConfig

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesSetConfig using JsonCodec
func (t AmuletRulesSetConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesSetConfig using JsonCodec
func (t *AmuletRulesSetConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesSetConfigResult is a Record type
type AmuletRulesSetConfigResult struct {
	NewAmuletRules CONTRACT_ID `json:"newAmuletRules"`
}

// ToMap converts AmuletRulesSetConfigResult to a map for DAML arguments
func (t AmuletRulesSetConfigResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["newAmuletRules"] = t.NewAmuletRules

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesSetConfigResult using JsonCodec
func (t AmuletRulesSetConfigResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesSetConfigResult using JsonCodec
func (t *AmuletRulesSetConfigResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesTransfer is a Record type
type AmuletRulesTransfer struct {
	Transfer    Transfer        `json:"transfer"`
	Context     TransferContext `json:"context"`
	ExpectedDso *PARTY          `json:"expectedDso"`
}

// ToMap converts AmuletRulesTransfer to a map for DAML arguments
func (t AmuletRulesTransfer) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["transfer"] = t.Transfer

	m["context"] = t.Context

	if t.ExpectedDso != nil {
		m["expectedDso"] = map[string]interface{}{
			"_type": "optional",
			"value": (*t.ExpectedDso).ToMap(),
		}
	} else {
		m["expectedDso"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesTransfer using JsonCodec
func (t AmuletRulesTransfer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesTransfer using JsonCodec
func (t *AmuletRulesTransfer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesUpdateFutureAmuletConfigSchedule is a Record type
type AmuletRulesUpdateFutureAmuletConfigSchedule struct {
	ScheduleItem TUPLE2[TIMESTAMP, AmuletConfig] `json:"scheduleItem"`
}

// ToMap converts AmuletRulesUpdateFutureAmuletConfigSchedule to a map for DAML arguments
func (t AmuletRulesUpdateFutureAmuletConfigSchedule) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["scheduleItem"] = t.ScheduleItem

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesUpdateFutureAmuletConfigSchedule using JsonCodec
func (t AmuletRulesUpdateFutureAmuletConfigSchedule) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesUpdateFutureAmuletConfigSchedule using JsonCodec
func (t *AmuletRulesUpdateFutureAmuletConfigSchedule) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletRulesUpdateFutureAmuletConfigScheduleResult is a Record type
type AmuletRulesUpdateFutureAmuletConfigScheduleResult struct {
	NewAmuletRules CONTRACT_ID `json:"newAmuletRules"`
}

// ToMap converts AmuletRulesUpdateFutureAmuletConfigScheduleResult to a map for DAML arguments
func (t AmuletRulesUpdateFutureAmuletConfigScheduleResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["newAmuletRules"] = t.NewAmuletRules

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletRulesUpdateFutureAmuletConfigScheduleResult using JsonCodec
func (t AmuletRulesUpdateFutureAmuletConfigScheduleResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletRulesUpdateFutureAmuletConfigScheduleResult using JsonCodec
func (t *AmuletRulesUpdateFutureAmuletConfigScheduleResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletTransferInstruction is a Template type
type AmuletTransferInstruction struct {
	LockedAmulet CONTRACT_ID `json:"lockedAmulet"`
	Transfer     Transfer    `json:"transfer"`
}

// GetTemplateID returns the template ID for this template
func (t AmuletTransferInstruction) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletTransferInstruction", "AmuletTransferInstruction")
}

// CreateCommand returns a CreateCommand for this template
func (t AmuletTransferInstruction) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["lockedAmulet"] = t.LockedAmulet

	args["transfer"] = t.Transfer

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for AmuletTransferInstruction using JsonCodec
func (t AmuletTransferInstruction) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletTransferInstruction using JsonCodec
func (t *AmuletTransferInstruction) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for AmuletTransferInstruction

// Archive exercises the Archive choice on this AmuletTransferInstruction contract
func (t AmuletTransferInstruction) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletTransferInstruction", "AmuletTransferInstruction"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// TransferInstructionAccept exercises the TransferInstruction_Accept choice on this AmuletTransferInstruction contract via the ITransferInstruction interface
func (t AmuletTransferInstruction) TransferInstructionAccept(contractID string, args TransferInstructionAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletTransferInstruction", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "TransferInstruction_Accept",
		Arguments:  argsToMap(args),
	}
}

// TransferInstructionReject exercises the TransferInstruction_Reject choice on this AmuletTransferInstruction contract via the ITransferInstruction interface
func (t AmuletTransferInstruction) TransferInstructionReject(contractID string, args TransferInstructionReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletTransferInstruction", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "TransferInstruction_Reject",
		Arguments:  argsToMap(args),
	}
}

// TransferInstructionWithdraw exercises the TransferInstruction_Withdraw choice on this AmuletTransferInstruction contract via the ITransferInstruction interface
func (t AmuletTransferInstruction) TransferInstructionWithdraw(contractID string, args TransferInstructionWithdraw) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletTransferInstruction", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "TransferInstruction_Withdraw",
		Arguments:  argsToMap(args),
	}
}

// TransferInstructionUpdate exercises the TransferInstruction_Update choice on this AmuletTransferInstruction contract via the ITransferInstruction interface
func (t AmuletTransferInstruction) TransferInstructionUpdate(contractID string, args TransferInstructionUpdate) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletTransferInstruction", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "TransferInstruction_Update",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for AmuletTransferInstruction

var _ ITransferInstruction = (*AmuletTransferInstruction)(nil)

// AmuletExpire is a Record type
type AmuletExpire struct {
	RoundCid CONTRACT_ID `json:"roundCid"`
}

// ToMap converts AmuletExpire to a map for DAML arguments
func (t AmuletExpire) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["roundCid"] = t.RoundCid

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletExpire using JsonCodec
func (t AmuletExpire) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletExpire using JsonCodec
func (t *AmuletExpire) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AmuletExpireResult is a Record type
type AmuletExpireResult struct {
	ExpireSum AmuletExpireSummary `json:"expireSum"`
	Meta      *Metadata           `json:"meta"`
}

// ToMap converts AmuletExpireResult to a map for DAML arguments
func (t AmuletExpireResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["expireSum"] = t.ExpireSum

	if t.Meta != nil {
		m["meta"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.Meta,
		}
	} else {
		m["meta"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for AmuletExpireResult using JsonCodec
func (t AmuletExpireResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AmuletExpireResult using JsonCodec
func (t *AmuletExpireResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AppRewardCoupon is a Template type
type AppRewardCoupon struct {
	Dso         PARTY   `json:"dso"`
	Provider    PARTY   `json:"provider"`
	Featured    BOOL    `json:"featured"`
	Amount      NUMERIC `json:"amount"`
	Round       Round   `json:"round"`
	Beneficiary *PARTY  `json:"beneficiary"`
}

// GetTemplateID returns the template ID for this template
func (t AppRewardCoupon) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "AppRewardCoupon")
}

// CreateCommand returns a CreateCommand for this template
func (t AppRewardCoupon) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["dso"] = t.Dso.ToMap()

	args["provider"] = t.Provider.ToMap()

	args["featured"] = bool(t.Featured)

	if t.Amount != nil {
		args["amount"] = (*big.Int)(t.Amount)
	}

	args["round"] = t.Round

	if t.Beneficiary != nil {
		args["beneficiary"] = map[string]interface{}{
			"_type": "optional",
			"value": (*t.Beneficiary).ToMap(),
		}
	} else {
		args["beneficiary"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for AppRewardCoupon using JsonCodec
func (t AppRewardCoupon) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AppRewardCoupon using JsonCodec
func (t *AppRewardCoupon) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for AppRewardCoupon

// AppRewardCouponDsoExpire exercises the AppRewardCoupon_DsoExpire choice on this AppRewardCoupon contract
func (t AppRewardCoupon) AppRewardCouponDsoExpire(contractID string, args AppRewardCouponDsoExpire) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "AppRewardCoupon"),
		ContractID: contractID,
		Choice:     "AppRewardCoupon_DsoExpire",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this AppRewardCoupon contract
func (t AppRewardCoupon) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "AppRewardCoupon"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// AppRewardCouponDsoExpire is a Record type
type AppRewardCouponDsoExpire struct {
	ClosedRoundCid CONTRACT_ID `json:"closedRoundCid"`
}

// ToMap converts AppRewardCouponDsoExpire to a map for DAML arguments
func (t AppRewardCouponDsoExpire) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["closedRoundCid"] = t.ClosedRoundCid

	return m
}

// MarshalJSON implements custom JSON marshaling for AppRewardCouponDsoExpire using JsonCodec
func (t AppRewardCouponDsoExpire) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AppRewardCouponDsoExpire using JsonCodec
func (t *AppRewardCouponDsoExpire) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AppRewardCouponDsoExpireResult is a Record type
type AppRewardCouponDsoExpireResult struct {
	Featured BOOL    `json:"featured"`
	Amount   NUMERIC `json:"amount"`
}

// ToMap converts AppRewardCouponDsoExpireResult to a map for DAML arguments
func (t AppRewardCouponDsoExpireResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["featured"] = bool(t.Featured)

	m["amount"] = (*big.Int)(t.Amount)

	return m
}

// MarshalJSON implements custom JSON marshaling for AppRewardCouponDsoExpireResult using JsonCodec
func (t AppRewardCouponDsoExpireResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AppRewardCouponDsoExpireResult using JsonCodec
func (t *AppRewardCouponDsoExpireResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AppTransferContext is a Record type
type AppTransferContext struct {
	AmuletRules      CONTRACT_ID  `json:"amuletRules"`
	OpenMiningRound  CONTRACT_ID  `json:"openMiningRound"`
	FeaturedAppRight *CONTRACT_ID `json:"featuredAppRight"`
}

// ToMap converts AppTransferContext to a map for DAML arguments
func (t AppTransferContext) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["amuletRules"] = t.AmuletRules

	m["openMiningRound"] = t.OpenMiningRound

	if t.FeaturedAppRight != nil {
		m["featuredAppRight"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.FeaturedAppRight,
		}
	} else {
		m["featuredAppRight"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for AppTransferContext using JsonCodec
func (t AppTransferContext) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AppTransferContext using JsonCodec
func (t *AppTransferContext) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// BalanceChange is a Record type
type BalanceChange struct {
	ChangeToInitialAmountAsOfRoundZero NUMERIC `json:"changeToInitialAmountAsOfRoundZero"`
	ChangeToHoldingFeesRate            NUMERIC `json:"changeToHoldingFeesRate"`
}

// ToMap converts BalanceChange to a map for DAML arguments
func (t BalanceChange) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["changeToInitialAmountAsOfRoundZero"] = (*big.Int)(t.ChangeToInitialAmountAsOfRoundZero)

	m["changeToHoldingFeesRate"] = (*big.Int)(t.ChangeToHoldingFeesRate)

	return m
}

// MarshalJSON implements custom JSON marshaling for BalanceChange using JsonCodec
func (t BalanceChange) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for BalanceChange using JsonCodec
func (t *BalanceChange) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// BaseRateTrafficLimits is a Record type
type BaseRateTrafficLimits struct {
	BurstAmount INT64   `json:"burstAmount"`
	BurstWindow RELTIME `json:"burstWindow"`
}

// ToMap converts BaseRateTrafficLimits to a map for DAML arguments
func (t BaseRateTrafficLimits) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["burstAmount"] = int64(t.BurstAmount)

	m["burstWindow"] = t.BurstWindow

	return m
}

// MarshalJSON implements custom JSON marshaling for BaseRateTrafficLimits using JsonCodec
func (t BaseRateTrafficLimits) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for BaseRateTrafficLimits using JsonCodec
func (t *BaseRateTrafficLimits) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// BoundedSet is a variant/union type
type BoundedSet struct {
	Singleton     *interface{} `json:"Singleton,omitempty"`
	AfterMaxBound *UNIT        `json:"AfterMaxBound,omitempty"`
}

// MarshalJSON implements custom JSON marshaling for BoundedSet
func (v BoundedSet) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(v)
}

// UnmarshalJSON implements custom JSON unmarshaling for BoundedSet
func (v *BoundedSet) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, v)
}

// GetVariantTag implements types.VARIANT interface
func (v BoundedSet) GetVariantTag() string {
	if v.Singleton != nil {
		return "Singleton"
	}

	if v.AfterMaxBound != nil {
		return "AfterMaxBound"
	}

	return ""
}

// GetVariantValue implements types.VARIANT interface
func (v BoundedSet) GetVariantValue() interface{} {
	if v.Singleton != nil {
		return v.Singleton
	}

	if v.AfterMaxBound != nil {
		return v.AfterMaxBound
	}

	return nil
}

// Verify interface implementation
var _ VARIANT = (*BoundedSet)(nil)

// ClosedMiningRound is a Template type
type ClosedMiningRound struct {
	Dso                                  PARTY    `json:"dso"`
	Round                                Round    `json:"round"`
	IssuancePerValidatorRewardCoupon     NUMERIC  `json:"issuancePerValidatorRewardCoupon"`
	IssuancePerFeaturedAppRewardCoupon   NUMERIC  `json:"issuancePerFeaturedAppRewardCoupon"`
	IssuancePerUnfeaturedAppRewardCoupon NUMERIC  `json:"issuancePerUnfeaturedAppRewardCoupon"`
	IssuancePerSvRewardCoupon            NUMERIC  `json:"issuancePerSvRewardCoupon"`
	OptIssuancePerValidatorFaucetCoupon  *NUMERIC `json:"optIssuancePerValidatorFaucetCoupon"`
}

// GetTemplateID returns the template ID for this template
func (t ClosedMiningRound) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Round", "ClosedMiningRound")
}

// CreateCommand returns a CreateCommand for this template
func (t ClosedMiningRound) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["dso"] = t.Dso.ToMap()

	args["round"] = t.Round

	if t.IssuancePerValidatorRewardCoupon != nil {
		args["issuancePerValidatorRewardCoupon"] = (*big.Int)(t.IssuancePerValidatorRewardCoupon)
	}

	if t.IssuancePerFeaturedAppRewardCoupon != nil {
		args["issuancePerFeaturedAppRewardCoupon"] = (*big.Int)(t.IssuancePerFeaturedAppRewardCoupon)
	}

	if t.IssuancePerUnfeaturedAppRewardCoupon != nil {
		args["issuancePerUnfeaturedAppRewardCoupon"] = (*big.Int)(t.IssuancePerUnfeaturedAppRewardCoupon)
	}

	if t.IssuancePerSvRewardCoupon != nil {
		args["issuancePerSvRewardCoupon"] = (*big.Int)(t.IssuancePerSvRewardCoupon)
	}

	if t.OptIssuancePerValidatorFaucetCoupon != nil {
		args["optIssuancePerValidatorFaucetCoupon"] = map[string]interface{}{
			"_type": "optional",
			"value": (*big.Int)(*t.OptIssuancePerValidatorFaucetCoupon),
		}
	} else {
		args["optIssuancePerValidatorFaucetCoupon"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for ClosedMiningRound using JsonCodec
func (t ClosedMiningRound) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ClosedMiningRound using JsonCodec
func (t *ClosedMiningRound) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for ClosedMiningRound

// Archive exercises the Archive choice on this ClosedMiningRound contract
func (t ClosedMiningRound) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Round", "ClosedMiningRound"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// CreatedAmulet is a variant/union type
type CreatedAmulet struct {
	TransferResultAmulet       *CONTRACT_ID      `json:"TransferResultAmulet,omitempty"`
	TransferResultLockedAmulet *CONTRACT_ID      `json:"TransferResultLockedAmulet,omitempty"`
	ExtCreatedAmulet           *ExtCreatedAmulet `json:"ExtCreatedAmulet,omitempty"`
}

// MarshalJSON implements custom JSON marshaling for CreatedAmulet
func (v CreatedAmulet) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(v)
}

// UnmarshalJSON implements custom JSON unmarshaling for CreatedAmulet
func (v *CreatedAmulet) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, v)
}

// GetVariantTag implements types.VARIANT interface
func (v CreatedAmulet) GetVariantTag() string {
	if v.TransferResultAmulet != nil {
		return "TransferResultAmulet"
	}

	if v.TransferResultLockedAmulet != nil {
		return "TransferResultLockedAmulet"
	}

	if v.ExtCreatedAmulet != nil {
		return "ExtCreatedAmulet"
	}

	return ""
}

// GetVariantValue implements types.VARIANT interface
func (v CreatedAmulet) GetVariantValue() interface{} {
	if v.TransferResultAmulet != nil {
		return v.TransferResultAmulet
	}

	if v.TransferResultLockedAmulet != nil {
		return v.TransferResultLockedAmulet
	}

	if v.ExtCreatedAmulet != nil {
		return v.ExtCreatedAmulet
	}

	return nil
}

// Verify interface implementation
var _ VARIANT = (*CreatedAmulet)(nil)

// DevelopmentFundCoupon is a Template type
type DevelopmentFundCoupon struct {
	Dso         PARTY     `json:"dso"`
	Beneficiary PARTY     `json:"beneficiary"`
	FundManager PARTY     `json:"fundManager"`
	Amount      NUMERIC   `json:"amount"`
	ExpiresAt   TIMESTAMP `json:"expiresAt"`
	Reason      TEXT      `json:"reason"`
}

// GetTemplateID returns the template ID for this template
func (t DevelopmentFundCoupon) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "DevelopmentFundCoupon")
}

// CreateCommand returns a CreateCommand for this template
func (t DevelopmentFundCoupon) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["dso"] = t.Dso.ToMap()

	args["beneficiary"] = t.Beneficiary.ToMap()

	args["fundManager"] = t.FundManager.ToMap()

	if t.Amount != nil {
		args["amount"] = (*big.Int)(t.Amount)
	}

	args["expiresAt"] = t.ExpiresAt

	args["reason"] = string(t.Reason)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for DevelopmentFundCoupon using JsonCodec
func (t DevelopmentFundCoupon) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for DevelopmentFundCoupon using JsonCodec
func (t *DevelopmentFundCoupon) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for DevelopmentFundCoupon

// DevelopmentFundCouponWithdraw exercises the DevelopmentFundCoupon_Withdraw choice on this DevelopmentFundCoupon contract
func (t DevelopmentFundCoupon) DevelopmentFundCouponWithdraw(contractID string, args DevelopmentFundCouponWithdraw) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "DevelopmentFundCoupon"),
		ContractID: contractID,
		Choice:     "DevelopmentFundCoupon_Withdraw",
		Arguments:  argsToMap(args),
	}
}

// DevelopmentFundCouponReject exercises the DevelopmentFundCoupon_Reject choice on this DevelopmentFundCoupon contract
func (t DevelopmentFundCoupon) DevelopmentFundCouponReject(contractID string, args DevelopmentFundCouponReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "DevelopmentFundCoupon"),
		ContractID: contractID,
		Choice:     "DevelopmentFundCoupon_Reject",
		Arguments:  argsToMap(args),
	}
}

// DevelopmentFundCouponDsoExpire exercises the DevelopmentFundCoupon_DsoExpire choice on this DevelopmentFundCoupon contract
func (t DevelopmentFundCoupon) DevelopmentFundCouponDsoExpire(contractID string, args DevelopmentFundCouponDsoExpire) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "DevelopmentFundCoupon"),
		ContractID: contractID,
		Choice:     "DevelopmentFundCoupon_DsoExpire",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this DevelopmentFundCoupon contract
func (t DevelopmentFundCoupon) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "DevelopmentFundCoupon"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// DevelopmentFundCouponDsoExpire is a Record type
type DevelopmentFundCouponDsoExpire struct{}

// ToMap converts DevelopmentFundCouponDsoExpire to a map for DAML arguments
func (t DevelopmentFundCouponDsoExpire) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	return m
}

// MarshalJSON implements custom JSON marshaling for DevelopmentFundCouponDsoExpire using JsonCodec
func (t DevelopmentFundCouponDsoExpire) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for DevelopmentFundCouponDsoExpire using JsonCodec
func (t *DevelopmentFundCouponDsoExpire) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// DevelopmentFundCouponDsoExpireResult is a Record type
type DevelopmentFundCouponDsoExpireResult struct {
	UnclaimedDevelopmentFundCouponCid CONTRACT_ID `json:"unclaimedDevelopmentFundCouponCid"`
}

// ToMap converts DevelopmentFundCouponDsoExpireResult to a map for DAML arguments
func (t DevelopmentFundCouponDsoExpireResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["unclaimedDevelopmentFundCouponCid"] = t.UnclaimedDevelopmentFundCouponCid

	return m
}

// MarshalJSON implements custom JSON marshaling for DevelopmentFundCouponDsoExpireResult using JsonCodec
func (t DevelopmentFundCouponDsoExpireResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for DevelopmentFundCouponDsoExpireResult using JsonCodec
func (t *DevelopmentFundCouponDsoExpireResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// DevelopmentFundCouponReject is a Record type
type DevelopmentFundCouponReject struct {
	Reason TEXT `json:"reason"`
}

// ToMap converts DevelopmentFundCouponReject to a map for DAML arguments
func (t DevelopmentFundCouponReject) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["reason"] = string(t.Reason)

	return m
}

// MarshalJSON implements custom JSON marshaling for DevelopmentFundCouponReject using JsonCodec
func (t DevelopmentFundCouponReject) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for DevelopmentFundCouponReject using JsonCodec
func (t *DevelopmentFundCouponReject) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// DevelopmentFundCouponRejectResult is a Record type
type DevelopmentFundCouponRejectResult struct {
	UnclaimedDevelopmentFundCouponCid CONTRACT_ID `json:"unclaimedDevelopmentFundCouponCid"`
}

// ToMap converts DevelopmentFundCouponRejectResult to a map for DAML arguments
func (t DevelopmentFundCouponRejectResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["unclaimedDevelopmentFundCouponCid"] = t.UnclaimedDevelopmentFundCouponCid

	return m
}

// MarshalJSON implements custom JSON marshaling for DevelopmentFundCouponRejectResult using JsonCodec
func (t DevelopmentFundCouponRejectResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for DevelopmentFundCouponRejectResult using JsonCodec
func (t *DevelopmentFundCouponRejectResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// DevelopmentFundCouponWithdraw is a Record type
type DevelopmentFundCouponWithdraw struct {
	Reason TEXT `json:"reason"`
}

// ToMap converts DevelopmentFundCouponWithdraw to a map for DAML arguments
func (t DevelopmentFundCouponWithdraw) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["reason"] = string(t.Reason)

	return m
}

// MarshalJSON implements custom JSON marshaling for DevelopmentFundCouponWithdraw using JsonCodec
func (t DevelopmentFundCouponWithdraw) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for DevelopmentFundCouponWithdraw using JsonCodec
func (t *DevelopmentFundCouponWithdraw) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// DevelopmentFundCouponWithdrawResult is a Record type
type DevelopmentFundCouponWithdrawResult struct {
	UnclaimedDevelopmentFundCouponCid CONTRACT_ID `json:"unclaimedDevelopmentFundCouponCid"`
}

// ToMap converts DevelopmentFundCouponWithdrawResult to a map for DAML arguments
func (t DevelopmentFundCouponWithdrawResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["unclaimedDevelopmentFundCouponCid"] = t.UnclaimedDevelopmentFundCouponCid

	return m
}

// MarshalJSON implements custom JSON marshaling for DevelopmentFundCouponWithdrawResult using JsonCodec
func (t DevelopmentFundCouponWithdrawResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for DevelopmentFundCouponWithdrawResult using JsonCodec
func (t *DevelopmentFundCouponWithdrawResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ExpiringAmount is a Record type
type ExpiringAmount struct {
	InitialAmount NUMERIC      `json:"initialAmount"`
	CreatedAt     Round        `json:"createdAt"`
	RatePerRound  RatePerRound `json:"ratePerRound"`
}

// ToMap converts ExpiringAmount to a map for DAML arguments
func (t ExpiringAmount) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["initialAmount"] = (*big.Int)(t.InitialAmount)

	m["createdAt"] = t.CreatedAt

	m["ratePerRound"] = t.RatePerRound

	return m
}

// MarshalJSON implements custom JSON marshaling for ExpiringAmount using JsonCodec
func (t ExpiringAmount) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ExpiringAmount using JsonCodec
func (t *ExpiringAmount) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ExtCreatedAmulet is a Record type
type ExtCreatedAmulet struct {
	DummyUnitField UNIT `json:"dummyUnitField"`
}

// ToMap converts ExtCreatedAmulet to a map for DAML arguments
func (t ExtCreatedAmulet) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["dummyUnitField"] = map[string]interface{}{"_type": "unit"}

	return m
}

// MarshalJSON implements custom JSON marshaling for ExtCreatedAmulet using JsonCodec
func (t ExtCreatedAmulet) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ExtCreatedAmulet using JsonCodec
func (t *ExtCreatedAmulet) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ExtInvalidTransferReason is a Record type
type ExtInvalidTransferReason struct {
	DummyUnitField UNIT `json:"dummyUnitField"`
}

// ToMap converts ExtInvalidTransferReason to a map for DAML arguments
func (t ExtInvalidTransferReason) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["dummyUnitField"] = map[string]interface{}{"_type": "unit"}

	return m
}

// MarshalJSON implements custom JSON marshaling for ExtInvalidTransferReason using JsonCodec
func (t ExtInvalidTransferReason) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ExtInvalidTransferReason using JsonCodec
func (t *ExtInvalidTransferReason) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ExtTransferInput is a Record type
type ExtTransferInput struct {
	DummyUnitField                UNIT         `json:"dummyUnitField"`
	OptInputValidatorFaucetCoupon *CONTRACT_ID `json:"optInputValidatorFaucetCoupon"`
}

// ToMap converts ExtTransferInput to a map for DAML arguments
func (t ExtTransferInput) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["dummyUnitField"] = map[string]interface{}{"_type": "unit"}

	if t.OptInputValidatorFaucetCoupon != nil {
		m["optInputValidatorFaucetCoupon"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.OptInputValidatorFaucetCoupon,
		}
	} else {
		m["optInputValidatorFaucetCoupon"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for ExtTransferInput using JsonCodec
func (t ExtTransferInput) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ExtTransferInput using JsonCodec
func (t *ExtTransferInput) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ExternalPartyAmuletRules is a Template type
type ExternalPartyAmuletRules struct {
	Dso PARTY `json:"dso"`
}

// GetTemplateID returns the template ID for this template
func (t ExternalPartyAmuletRules) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.ExternalPartyAmuletRules", "ExternalPartyAmuletRules")
}

// CreateCommand returns a CreateCommand for this template
func (t ExternalPartyAmuletRules) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["dso"] = t.Dso.ToMap()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for ExternalPartyAmuletRules using JsonCodec
func (t ExternalPartyAmuletRules) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ExternalPartyAmuletRules using JsonCodec
func (t *ExternalPartyAmuletRules) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for ExternalPartyAmuletRules

// ExternalPartyAmuletRulesCreateTransferCommand exercises the ExternalPartyAmuletRules_CreateTransferCommand choice on this ExternalPartyAmuletRules contract
func (t ExternalPartyAmuletRules) ExternalPartyAmuletRulesCreateTransferCommand(contractID string, args ExternalPartyAmuletRulesCreateTransferCommand) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.ExternalPartyAmuletRules", "ExternalPartyAmuletRules"),
		ContractID: contractID,
		Choice:     "ExternalPartyAmuletRules_CreateTransferCommand",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this ExternalPartyAmuletRules contract
func (t ExternalPartyAmuletRules) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.ExternalPartyAmuletRules", "ExternalPartyAmuletRules"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// TransferFactoryTransfer exercises the TransferFactory_Transfer choice on this ExternalPartyAmuletRules contract via the ITransferFactory interface
func (t ExternalPartyAmuletRules) TransferFactoryTransfer(contractID string, args TransferFactoryTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.ExternalPartyAmuletRules", "TransferFactory"),
		ContractID: contractID,
		Choice:     "TransferFactory_Transfer",
		Arguments:  argsToMap(args),
	}
}

// TransferFactoryPublicFetch exercises the TransferFactory_PublicFetch choice on this ExternalPartyAmuletRules contract via the ITransferFactory interface
func (t ExternalPartyAmuletRules) TransferFactoryPublicFetch(contractID string, args TransferFactoryPublicFetch) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.ExternalPartyAmuletRules", "TransferFactory"),
		ContractID: contractID,
		Choice:     "TransferFactory_PublicFetch",
		Arguments:  argsToMap(args),
	}
}

// AllocationFactoryAllocate exercises the AllocationFactory_Allocate choice on this ExternalPartyAmuletRules contract via the IAllocationFactory interface
func (t ExternalPartyAmuletRules) AllocationFactoryAllocate(contractID string, args AllocationFactoryAllocate) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.ExternalPartyAmuletRules", "AllocationFactory"),
		ContractID: contractID,
		Choice:     "AllocationFactory_Allocate",
		Arguments:  argsToMap(args),
	}
}

// AllocationFactoryPublicFetch exercises the AllocationFactory_PublicFetch choice on this ExternalPartyAmuletRules contract via the IAllocationFactory interface
func (t ExternalPartyAmuletRules) AllocationFactoryPublicFetch(contractID string, args AllocationFactoryPublicFetch) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.ExternalPartyAmuletRules", "AllocationFactory"),
		ContractID: contractID,
		Choice:     "AllocationFactory_PublicFetch",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for ExternalPartyAmuletRules

var _ ITransferFactory = (*ExternalPartyAmuletRules)(nil)

var _ IAllocationFactory = (*ExternalPartyAmuletRules)(nil)

// ExternalPartyAmuletRulesCreateTransferCommand is a Record type
type ExternalPartyAmuletRulesCreateTransferCommand struct {
	Sender      PARTY     `json:"sender"`
	Receiver    PARTY     `json:"receiver"`
	Delegate    PARTY     `json:"delegate"`
	Amount      NUMERIC   `json:"amount"`
	ExpiresAt   TIMESTAMP `json:"expiresAt"`
	Nonce       INT64     `json:"nonce"`
	Description *TEXT     `json:"description"`
	ExpectedDso *PARTY    `json:"expectedDso"`
}

// ToMap converts ExternalPartyAmuletRulesCreateTransferCommand to a map for DAML arguments
func (t ExternalPartyAmuletRulesCreateTransferCommand) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["sender"] = t.Sender.ToMap()

	m["receiver"] = t.Receiver.ToMap()

	m["delegate"] = t.Delegate.ToMap()

	m["amount"] = (*big.Int)(t.Amount)

	m["expiresAt"] = t.ExpiresAt

	m["nonce"] = int64(t.Nonce)

	if t.Description != nil {
		m["description"] = map[string]interface{}{
			"_type": "optional",
			"value": string(*t.Description),
		}
	} else {
		m["description"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	if t.ExpectedDso != nil {
		m["expectedDso"] = map[string]interface{}{
			"_type": "optional",
			"value": (*t.ExpectedDso).ToMap(),
		}
	} else {
		m["expectedDso"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for ExternalPartyAmuletRulesCreateTransferCommand using JsonCodec
func (t ExternalPartyAmuletRulesCreateTransferCommand) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ExternalPartyAmuletRulesCreateTransferCommand using JsonCodec
func (t *ExternalPartyAmuletRulesCreateTransferCommand) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ExternalPartyAmuletRulesCreateTransferCommandResult is a Record type
type ExternalPartyAmuletRulesCreateTransferCommandResult struct {
	TransferCommandCid CONTRACT_ID `json:"transferCommandCid"`
}

// ToMap converts ExternalPartyAmuletRulesCreateTransferCommandResult to a map for DAML arguments
func (t ExternalPartyAmuletRulesCreateTransferCommandResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["transferCommandCid"] = t.TransferCommandCid

	return m
}

// MarshalJSON implements custom JSON marshaling for ExternalPartyAmuletRulesCreateTransferCommandResult using JsonCodec
func (t ExternalPartyAmuletRulesCreateTransferCommandResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ExternalPartyAmuletRulesCreateTransferCommandResult using JsonCodec
func (t *ExternalPartyAmuletRulesCreateTransferCommandResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ExternalPartySetupProposal is a Template type
type ExternalPartySetupProposal struct {
	Validator            PARTY     `json:"validator"`
	User                 PARTY     `json:"user"`
	Dso                  PARTY     `json:"dso"`
	CreatedAt            TIMESTAMP `json:"createdAt"`
	PreapprovalExpiresAt TIMESTAMP `json:"preapprovalExpiresAt"`
}

// GetTemplateID returns the template ID for this template
func (t ExternalPartySetupProposal) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletRules", "ExternalPartySetupProposal")
}

// CreateCommand returns a CreateCommand for this template
func (t ExternalPartySetupProposal) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["validator"] = t.Validator.ToMap()

	args["user"] = t.User.ToMap()

	args["dso"] = t.Dso.ToMap()

	args["createdAt"] = t.CreatedAt

	args["preapprovalExpiresAt"] = t.PreapprovalExpiresAt

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for ExternalPartySetupProposal using JsonCodec
func (t ExternalPartySetupProposal) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ExternalPartySetupProposal using JsonCodec
func (t *ExternalPartySetupProposal) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for ExternalPartySetupProposal

// ExternalPartySetupProposalAccept exercises the ExternalPartySetupProposal_Accept choice on this ExternalPartySetupProposal contract
func (t ExternalPartySetupProposal) ExternalPartySetupProposalAccept(contractID string, args ExternalPartySetupProposalAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletRules", "ExternalPartySetupProposal"),
		ContractID: contractID,
		Choice:     "ExternalPartySetupProposal_Accept",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this ExternalPartySetupProposal contract
func (t ExternalPartySetupProposal) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletRules", "ExternalPartySetupProposal"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// ExternalPartySetupProposalReject exercises the ExternalPartySetupProposal_Reject choice on this ExternalPartySetupProposal contract
func (t ExternalPartySetupProposal) ExternalPartySetupProposalReject(contractID string, args ExternalPartySetupProposalReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletRules", "ExternalPartySetupProposal"),
		ContractID: contractID,
		Choice:     "ExternalPartySetupProposal_Reject",
		Arguments:  argsToMap(args),
	}
}

// ExternalPartySetupProposalWithdraw exercises the ExternalPartySetupProposal_Withdraw choice on this ExternalPartySetupProposal contract
func (t ExternalPartySetupProposal) ExternalPartySetupProposalWithdraw(contractID string, args ExternalPartySetupProposalWithdraw) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletRules", "ExternalPartySetupProposal"),
		ContractID: contractID,
		Choice:     "ExternalPartySetupProposal_Withdraw",
		Arguments:  argsToMap(args),
	}
}

// ExternalPartySetupProposalAccept is a Record type
type ExternalPartySetupProposalAccept struct{}

// ToMap converts ExternalPartySetupProposalAccept to a map for DAML arguments
func (t ExternalPartySetupProposalAccept) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	return m
}

// MarshalJSON implements custom JSON marshaling for ExternalPartySetupProposalAccept using JsonCodec
func (t ExternalPartySetupProposalAccept) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ExternalPartySetupProposalAccept using JsonCodec
func (t *ExternalPartySetupProposalAccept) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ExternalPartySetupProposalAcceptResult is a Record type
type ExternalPartySetupProposalAcceptResult struct {
	ValidatorRightCid      CONTRACT_ID `json:"validatorRightCid"`
	TransferPreapprovalCid CONTRACT_ID `json:"transferPreapprovalCid"`
}

// ToMap converts ExternalPartySetupProposalAcceptResult to a map for DAML arguments
func (t ExternalPartySetupProposalAcceptResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["validatorRightCid"] = t.ValidatorRightCid

	m["transferPreapprovalCid"] = t.TransferPreapprovalCid

	return m
}

// MarshalJSON implements custom JSON marshaling for ExternalPartySetupProposalAcceptResult using JsonCodec
func (t ExternalPartySetupProposalAcceptResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ExternalPartySetupProposalAcceptResult using JsonCodec
func (t *ExternalPartySetupProposalAcceptResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ExternalPartySetupProposalReject is a Record type
type ExternalPartySetupProposalReject struct {
	Reason TEXT `json:"reason"`
}

// ToMap converts ExternalPartySetupProposalReject to a map for DAML arguments
func (t ExternalPartySetupProposalReject) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["reason"] = string(t.Reason)

	return m
}

// MarshalJSON implements custom JSON marshaling for ExternalPartySetupProposalReject using JsonCodec
func (t ExternalPartySetupProposalReject) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ExternalPartySetupProposalReject using JsonCodec
func (t *ExternalPartySetupProposalReject) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ExternalPartySetupProposalRejectResult is a Record type
type ExternalPartySetupProposalRejectResult struct {
	DummyArg UNIT `json:"dummyArg"`
}

// ToMap converts ExternalPartySetupProposalRejectResult to a map for DAML arguments
func (t ExternalPartySetupProposalRejectResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["dummyArg"] = map[string]interface{}{"_type": "unit"}

	return m
}

// MarshalJSON implements custom JSON marshaling for ExternalPartySetupProposalRejectResult using JsonCodec
func (t ExternalPartySetupProposalRejectResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ExternalPartySetupProposalRejectResult using JsonCodec
func (t *ExternalPartySetupProposalRejectResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ExternalPartySetupProposalWithdraw is a Record type
type ExternalPartySetupProposalWithdraw struct {
	Reason TEXT `json:"reason"`
}

// ToMap converts ExternalPartySetupProposalWithdraw to a map for DAML arguments
func (t ExternalPartySetupProposalWithdraw) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["reason"] = string(t.Reason)

	return m
}

// MarshalJSON implements custom JSON marshaling for ExternalPartySetupProposalWithdraw using JsonCodec
func (t ExternalPartySetupProposalWithdraw) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ExternalPartySetupProposalWithdraw using JsonCodec
func (t *ExternalPartySetupProposalWithdraw) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ExternalPartySetupProposalWithdrawResult is a Record type
type ExternalPartySetupProposalWithdrawResult struct {
	DummyArg UNIT `json:"dummyArg"`
}

// ToMap converts ExternalPartySetupProposalWithdrawResult to a map for DAML arguments
func (t ExternalPartySetupProposalWithdrawResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["dummyArg"] = map[string]interface{}{"_type": "unit"}

	return m
}

// MarshalJSON implements custom JSON marshaling for ExternalPartySetupProposalWithdrawResult using JsonCodec
func (t ExternalPartySetupProposalWithdrawResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ExternalPartySetupProposalWithdrawResult using JsonCodec
func (t *ExternalPartySetupProposalWithdrawResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// FaucetState is a Record type
type FaucetState struct {
	FirstReceivedFor Round `json:"firstReceivedFor"`
	LastReceivedFor  Round `json:"lastReceivedFor"`
	NumCouponsMissed INT64 `json:"numCouponsMissed"`
}

// ToMap converts FaucetState to a map for DAML arguments
func (t FaucetState) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["firstReceivedFor"] = t.FirstReceivedFor

	m["lastReceivedFor"] = t.LastReceivedFor

	m["numCouponsMissed"] = int64(t.NumCouponsMissed)

	return m
}

// MarshalJSON implements custom JSON marshaling for FaucetState using JsonCodec
func (t FaucetState) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for FaucetState using JsonCodec
func (t *FaucetState) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// FeaturedAppActivityMarker is a Template type
type FeaturedAppActivityMarker struct {
	Dso         PARTY   `json:"dso"`
	Provider    PARTY   `json:"provider"`
	Beneficiary PARTY   `json:"beneficiary"`
	Weight      NUMERIC `json:"weight"`
}

// GetTemplateID returns the template ID for this template
func (t FeaturedAppActivityMarker) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "FeaturedAppActivityMarker")
}

// CreateCommand returns a CreateCommand for this template
func (t FeaturedAppActivityMarker) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["dso"] = t.Dso.ToMap()

	args["provider"] = t.Provider.ToMap()

	args["beneficiary"] = t.Beneficiary.ToMap()

	if t.Weight != nil {
		args["weight"] = (*big.Int)(t.Weight)
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for FeaturedAppActivityMarker using JsonCodec
func (t FeaturedAppActivityMarker) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for FeaturedAppActivityMarker using JsonCodec
func (t *FeaturedAppActivityMarker) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for FeaturedAppActivityMarker

// Archive exercises the Archive choice on this FeaturedAppActivityMarker contract
func (t FeaturedAppActivityMarker) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "FeaturedAppActivityMarker"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// Verify interface implementations for FeaturedAppActivityMarker

var _ IFeaturedAppActivityMarker = (*FeaturedAppActivityMarker)(nil)

var _ IFeaturedAppActivityMarker = (*FeaturedAppActivityMarker)(nil)

// FeaturedAppRight is a Template type
type FeaturedAppRight struct {
	Dso      PARTY `json:"dso"`
	Provider PARTY `json:"provider"`
}

// GetTemplateID returns the template ID for this template
func (t FeaturedAppRight) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "FeaturedAppRight")
}

// CreateCommand returns a CreateCommand for this template
func (t FeaturedAppRight) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["dso"] = t.Dso.ToMap()

	args["provider"] = t.Provider.ToMap()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for FeaturedAppRight using JsonCodec
func (t FeaturedAppRight) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for FeaturedAppRight using JsonCodec
func (t *FeaturedAppRight) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for FeaturedAppRight

// FeaturedAppRightWithdraw exercises the FeaturedAppRight_Withdraw choice on this FeaturedAppRight contract
func (t FeaturedAppRight) FeaturedAppRightWithdraw(contractID string, args FeaturedAppRightWithdraw) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "FeaturedAppRight"),
		ContractID: contractID,
		Choice:     "FeaturedAppRight_Withdraw",
		Arguments:  argsToMap(args),
	}
}

// FeaturedAppRightCancel exercises the FeaturedAppRight_Cancel choice on this FeaturedAppRight contract
func (t FeaturedAppRight) FeaturedAppRightCancel(contractID string, args FeaturedAppRightCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "FeaturedAppRight"),
		ContractID: contractID,
		Choice:     "FeaturedAppRight_Cancel",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this FeaturedAppRight contract
func (t FeaturedAppRight) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "FeaturedAppRight"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// FeaturedAppRightCreateActivityMarker exercises the FeaturedAppRight_CreateActivityMarker choice on this FeaturedAppRight contract via the IFeaturedAppRight interface
func (t FeaturedAppRight) FeaturedAppRightCreateActivityMarker(contractID string, args FeaturedAppRightCreateActivityMarker) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "FeaturedAppRight"),
		ContractID: contractID,
		Choice:     "FeaturedAppRight_CreateActivityMarker",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for FeaturedAppRight

var _ IFeaturedAppRight = (*FeaturedAppRight)(nil)

var _ IFeaturedAppRight = (*FeaturedAppRight)(nil)

// FeaturedAppRightCancel is a Record type
type FeaturedAppRightCancel struct{}

// ToMap converts FeaturedAppRightCancel to a map for DAML arguments
func (t FeaturedAppRightCancel) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	return m
}

// MarshalJSON implements custom JSON marshaling for FeaturedAppRightCancel using JsonCodec
func (t FeaturedAppRightCancel) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for FeaturedAppRightCancel using JsonCodec
func (t *FeaturedAppRightCancel) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// FeaturedAppRightCancelResult is an enum type
type FeaturedAppRightCancelResult string

const (
	FeaturedAppRightCancelResultFeaturedAppRight_CancelResult FeaturedAppRightCancelResult = "FeaturedAppRight_CancelResult"
)

// GetEnumConstructor implements types.ENUM interface
func (e FeaturedAppRightCancelResult) GetEnumConstructor() string {
	return string(e)
}

// GetEnumTypeID implements types.ENUM interface
func (e FeaturedAppRightCancelResult) GetEnumTypeID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "FeaturedAppRightCancelResult")
}

// MarshalJSON implements custom JSON marshaling for FeaturedAppRightCancelResult using JsonCodec
func (e FeaturedAppRightCancelResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(e)
}

// UnmarshalJSON implements custom JSON unmarshaling for FeaturedAppRightCancelResult using JsonCodec
func (e *FeaturedAppRightCancelResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, e)
}

// Verify interface implementation
var _ ENUM = FeaturedAppRightCancelResult("")

// FeaturedAppRightWithdraw is a Record type
type FeaturedAppRightWithdraw struct {
	Reason TEXT `json:"reason"`
}

// ToMap converts FeaturedAppRightWithdraw to a map for DAML arguments
func (t FeaturedAppRightWithdraw) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["reason"] = string(t.Reason)

	return m
}

// MarshalJSON implements custom JSON marshaling for FeaturedAppRightWithdraw using JsonCodec
func (t FeaturedAppRightWithdraw) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for FeaturedAppRightWithdraw using JsonCodec
func (t *FeaturedAppRightWithdraw) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// FeaturedAppRightWithdrawResult is an enum type
type FeaturedAppRightWithdrawResult string

const (
	FeaturedAppRightWithdrawResultFeaturedAppRight_WithdrawResult FeaturedAppRightWithdrawResult = "FeaturedAppRight_WithdrawResult"
)

// GetEnumConstructor implements types.ENUM interface
func (e FeaturedAppRightWithdrawResult) GetEnumConstructor() string {
	return string(e)
}

// GetEnumTypeID implements types.ENUM interface
func (e FeaturedAppRightWithdrawResult) GetEnumTypeID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "FeaturedAppRightWithdrawResult")
}

// MarshalJSON implements custom JSON marshaling for FeaturedAppRightWithdrawResult using JsonCodec
func (e FeaturedAppRightWithdrawResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(e)
}

// UnmarshalJSON implements custom JSON unmarshaling for FeaturedAppRightWithdrawResult using JsonCodec
func (e *FeaturedAppRightWithdrawResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, e)
}

// Verify interface implementation
var _ ENUM = FeaturedAppRightWithdrawResult("")

// FixedFee is a Record type
type FixedFee struct {
	Fee NUMERIC `json:"fee"`
}

// ToMap converts FixedFee to a map for DAML arguments
func (t FixedFee) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["fee"] = (*big.Int)(t.Fee)

	return m
}

// MarshalJSON implements custom JSON marshaling for FixedFee using JsonCodec
func (t FixedFee) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for FixedFee using JsonCodec
func (t *FixedFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ForDso is a Record type
type ForDso struct {
	Dso PARTY `json:"dso"`
}

// ToMap converts ForDso to a map for DAML arguments
func (t ForDso) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["dso"] = t.Dso.ToMap()

	return m
}

// MarshalJSON implements custom JSON marshaling for ForDso using JsonCodec
func (t ForDso) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ForDso using JsonCodec
func (t *ForDso) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ForMemberTraffic is a Record type
type ForMemberTraffic struct {
	Dso            PARTY `json:"dso"`
	MemberId       TEXT  `json:"memberId"`
	SynchronizerId TEXT  `json:"synchronizerId"`
	MigrationId    INT64 `json:"migrationId"`
}

// ToMap converts ForMemberTraffic to a map for DAML arguments
func (t ForMemberTraffic) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["dso"] = t.Dso.ToMap()

	m["memberId"] = string(t.MemberId)

	m["synchronizerId"] = string(t.SynchronizerId)

	m["migrationId"] = int64(t.MigrationId)

	return m
}

// MarshalJSON implements custom JSON marshaling for ForMemberTraffic using JsonCodec
func (t ForMemberTraffic) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ForMemberTraffic using JsonCodec
func (t *ForMemberTraffic) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ForOwner is a Record type
type ForOwner struct {
	Dso   PARTY `json:"dso"`
	Owner PARTY `json:"owner"`
}

// ToMap converts ForOwner to a map for DAML arguments
func (t ForOwner) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["dso"] = t.Dso.ToMap()

	m["owner"] = t.Owner.ToMap()

	return m
}

// MarshalJSON implements custom JSON marshaling for ForOwner using JsonCodec
func (t ForOwner) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ForOwner using JsonCodec
func (t *ForOwner) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ForRound is a Record type
type ForRound struct {
	Dso   PARTY `json:"dso"`
	Round Round `json:"round"`
}

// ToMap converts ForRound to a map for DAML arguments
func (t ForRound) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["dso"] = t.Dso.ToMap()

	m["round"] = t.Round

	return m
}

// MarshalJSON implements custom JSON marshaling for ForRound using JsonCodec
func (t ForRound) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ForRound using JsonCodec
func (t *ForRound) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ITRInsufficientFunds is a Record type
type ITRInsufficientFunds struct {
	MissingAmount NUMERIC `json:"missingAmount"`
}

// ToMap converts ITRInsufficientFunds to a map for DAML arguments
func (t ITRInsufficientFunds) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["missingAmount"] = (*big.Int)(t.MissingAmount)

	return m
}

// MarshalJSON implements custom JSON marshaling for ITRInsufficientFunds using JsonCodec
func (t ITRInsufficientFunds) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ITRInsufficientFunds using JsonCodec
func (t *ITRInsufficientFunds) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ITRInsufficientTopupAmount is a Record type
type ITRInsufficientTopupAmount struct {
	RequestedTopupAmount INT64 `json:"requestedTopupAmount"`
	MinTopupAmount       INT64 `json:"minTopupAmount"`
}

// ToMap converts ITRInsufficientTopupAmount to a map for DAML arguments
func (t ITRInsufficientTopupAmount) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["requestedTopupAmount"] = int64(t.RequestedTopupAmount)

	m["minTopupAmount"] = int64(t.MinTopupAmount)

	return m
}

// MarshalJSON implements custom JSON marshaling for ITRInsufficientTopupAmount using JsonCodec
func (t ITRInsufficientTopupAmount) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ITRInsufficientTopupAmount using JsonCodec
func (t *ITRInsufficientTopupAmount) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ITROther is a Record type
type ITROther struct {
	Description TEXT `json:"description"`
}

// ToMap converts ITROther to a map for DAML arguments
func (t ITROther) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["description"] = string(t.Description)

	return m
}

// MarshalJSON implements custom JSON marshaling for ITROther using JsonCodec
func (t ITROther) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ITROther using JsonCodec
func (t *ITROther) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ITRUnknownSynchronizer is a Record type
type ITRUnknownSynchronizer struct {
	SynchronizerId TEXT `json:"synchronizerId"`
}

// ToMap converts ITRUnknownSynchronizer to a map for DAML arguments
func (t ITRUnknownSynchronizer) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["synchronizerId"] = string(t.SynchronizerId)

	return m
}

// MarshalJSON implements custom JSON marshaling for ITRUnknownSynchronizer using JsonCodec
func (t ITRUnknownSynchronizer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ITRUnknownSynchronizer using JsonCodec
func (t *ITRUnknownSynchronizer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// InvalidTransfer is a Record type
type InvalidTransfer struct {
	Reason InvalidTransferReason `json:"reason"`
}

// ToMap converts InvalidTransfer to a map for DAML arguments
func (t InvalidTransfer) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["reason"] = t.Reason

	return m
}

// MarshalJSON implements custom JSON marshaling for InvalidTransfer using JsonCodec
func (t InvalidTransfer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for InvalidTransfer using JsonCodec
func (t *InvalidTransfer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// InvalidTransferReason is a variant/union type
type InvalidTransferReason struct {
	ITRInsufficientFunds       *ITRInsufficientFunds       `json:"ITR_InsufficientFunds,omitempty"`
	ITRUnknownSynchronizer     *ITRUnknownSynchronizer     `json:"ITR_UnknownSynchronizer,omitempty"`
	ITRInsufficientTopupAmount *ITRInsufficientTopupAmount `json:"ITR_InsufficientTopupAmount,omitempty"`
	ITROther                   *ITROther                   `json:"ITR_Other,omitempty"`
	ExtInvalidTransferReason   *ExtInvalidTransferReason   `json:"ExtInvalidTransferReason,omitempty"`
}

// MarshalJSON implements custom JSON marshaling for InvalidTransferReason
func (v InvalidTransferReason) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(v)
}

// UnmarshalJSON implements custom JSON unmarshaling for InvalidTransferReason
func (v *InvalidTransferReason) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, v)
}

// GetVariantTag implements types.VARIANT interface
func (v InvalidTransferReason) GetVariantTag() string {
	if v.ITRInsufficientFunds != nil {
		return "ITR_InsufficientFunds"
	}

	if v.ITRUnknownSynchronizer != nil {
		return "ITR_UnknownSynchronizer"
	}

	if v.ITRInsufficientTopupAmount != nil {
		return "ITR_InsufficientTopupAmount"
	}

	if v.ITROther != nil {
		return "ITR_Other"
	}

	if v.ExtInvalidTransferReason != nil {
		return "ExtInvalidTransferReason"
	}

	return ""
}

// GetVariantValue implements types.VARIANT interface
func (v InvalidTransferReason) GetVariantValue() interface{} {
	if v.ITRInsufficientFunds != nil {
		return v.ITRInsufficientFunds
	}

	if v.ITRUnknownSynchronizer != nil {
		return v.ITRUnknownSynchronizer
	}

	if v.ITRInsufficientTopupAmount != nil {
		return v.ITRInsufficientTopupAmount
	}

	if v.ITROther != nil {
		return v.ITROther
	}

	if v.ExtInvalidTransferReason != nil {
		return v.ExtInvalidTransferReason
	}

	return nil
}

// Verify interface implementation
var _ VARIANT = (*InvalidTransferReason)(nil)

// IssuanceConfig is a Record type
type IssuanceConfig struct {
	AmuletToIssuePerYear         NUMERIC  `json:"amuletToIssuePerYear"`
	ValidatorRewardPercentage    NUMERIC  `json:"validatorRewardPercentage"`
	AppRewardPercentage          NUMERIC  `json:"appRewardPercentage"`
	ValidatorRewardCap           NUMERIC  `json:"validatorRewardCap"`
	FeaturedAppRewardCap         NUMERIC  `json:"featuredAppRewardCap"`
	UnfeaturedAppRewardCap       NUMERIC  `json:"unfeaturedAppRewardCap"`
	OptValidatorFaucetCap        *NUMERIC `json:"optValidatorFaucetCap"`
	OptDevelopmentFundPercentage *NUMERIC `json:"optDevelopmentFundPercentage"`
}

// ToMap converts IssuanceConfig to a map for DAML arguments
func (t IssuanceConfig) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["amuletToIssuePerYear"] = (*big.Int)(t.AmuletToIssuePerYear)

	m["validatorRewardPercentage"] = (*big.Int)(t.ValidatorRewardPercentage)

	m["appRewardPercentage"] = (*big.Int)(t.AppRewardPercentage)

	m["validatorRewardCap"] = (*big.Int)(t.ValidatorRewardCap)

	m["featuredAppRewardCap"] = (*big.Int)(t.FeaturedAppRewardCap)

	m["unfeaturedAppRewardCap"] = (*big.Int)(t.UnfeaturedAppRewardCap)

	if t.OptValidatorFaucetCap != nil {
		m["optValidatorFaucetCap"] = map[string]interface{}{
			"_type": "optional",
			"value": (*big.Int)(*t.OptValidatorFaucetCap),
		}
	} else {
		m["optValidatorFaucetCap"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	if t.OptDevelopmentFundPercentage != nil {
		m["optDevelopmentFundPercentage"] = map[string]interface{}{
			"_type": "optional",
			"value": (*big.Int)(*t.OptDevelopmentFundPercentage),
		}
	} else {
		m["optDevelopmentFundPercentage"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for IssuanceConfig using JsonCodec
func (t IssuanceConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for IssuanceConfig using JsonCodec
func (t *IssuanceConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// IssuanceTranche is a Record type
type IssuanceTranche struct {
	RewardsToIssue    NUMERIC `json:"rewardsToIssue"`
	IssuancePerCoupon NUMERIC `json:"issuancePerCoupon"`
	UnclaimedRewards  NUMERIC `json:"unclaimedRewards"`
}

// ToMap converts IssuanceTranche to a map for DAML arguments
func (t IssuanceTranche) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["rewardsToIssue"] = (*big.Int)(t.RewardsToIssue)

	m["issuancePerCoupon"] = (*big.Int)(t.IssuancePerCoupon)

	m["unclaimedRewards"] = (*big.Int)(t.UnclaimedRewards)

	return m
}

// MarshalJSON implements custom JSON marshaling for IssuanceTranche using JsonCodec
func (t IssuanceTranche) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for IssuanceTranche using JsonCodec
func (t *IssuanceTranche) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// IssuingMiningRound is a Template type
type IssuingMiningRound struct {
	Dso                                  PARTY     `json:"dso"`
	Round                                Round     `json:"round"`
	IssuancePerValidatorRewardCoupon     NUMERIC   `json:"issuancePerValidatorRewardCoupon"`
	IssuancePerFeaturedAppRewardCoupon   NUMERIC   `json:"issuancePerFeaturedAppRewardCoupon"`
	IssuancePerUnfeaturedAppRewardCoupon NUMERIC   `json:"issuancePerUnfeaturedAppRewardCoupon"`
	IssuancePerSvRewardCoupon            NUMERIC   `json:"issuancePerSvRewardCoupon"`
	OpensAt                              TIMESTAMP `json:"opensAt"`
	TargetClosesAt                       TIMESTAMP `json:"targetClosesAt"`
	OptIssuancePerValidatorFaucetCoupon  *NUMERIC  `json:"optIssuancePerValidatorFaucetCoupon"`
}

// GetTemplateID returns the template ID for this template
func (t IssuingMiningRound) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Round", "IssuingMiningRound")
}

// CreateCommand returns a CreateCommand for this template
func (t IssuingMiningRound) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["dso"] = t.Dso.ToMap()

	args["round"] = t.Round

	if t.IssuancePerValidatorRewardCoupon != nil {
		args["issuancePerValidatorRewardCoupon"] = (*big.Int)(t.IssuancePerValidatorRewardCoupon)
	}

	if t.IssuancePerFeaturedAppRewardCoupon != nil {
		args["issuancePerFeaturedAppRewardCoupon"] = (*big.Int)(t.IssuancePerFeaturedAppRewardCoupon)
	}

	if t.IssuancePerUnfeaturedAppRewardCoupon != nil {
		args["issuancePerUnfeaturedAppRewardCoupon"] = (*big.Int)(t.IssuancePerUnfeaturedAppRewardCoupon)
	}

	if t.IssuancePerSvRewardCoupon != nil {
		args["issuancePerSvRewardCoupon"] = (*big.Int)(t.IssuancePerSvRewardCoupon)
	}

	args["opensAt"] = t.OpensAt

	args["targetClosesAt"] = t.TargetClosesAt

	if t.OptIssuancePerValidatorFaucetCoupon != nil {
		args["optIssuancePerValidatorFaucetCoupon"] = map[string]interface{}{
			"_type": "optional",
			"value": (*big.Int)(*t.OptIssuancePerValidatorFaucetCoupon),
		}
	} else {
		args["optIssuancePerValidatorFaucetCoupon"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for IssuingMiningRound using JsonCodec
func (t IssuingMiningRound) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for IssuingMiningRound using JsonCodec
func (t *IssuingMiningRound) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for IssuingMiningRound

// Archive exercises the Archive choice on this IssuingMiningRound contract
func (t IssuingMiningRound) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Round", "IssuingMiningRound"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// IssuingRoundParameters is a Record type
type IssuingRoundParameters struct {
	IssuancePerValidatorRewardCoupon     NUMERIC  `json:"issuancePerValidatorRewardCoupon"`
	IssuancePerFeaturedAppRewardCoupon   NUMERIC  `json:"issuancePerFeaturedAppRewardCoupon"`
	IssuancePerUnfeaturedAppRewardCoupon NUMERIC  `json:"issuancePerUnfeaturedAppRewardCoupon"`
	IssuancePerSvRewardCoupon            NUMERIC  `json:"issuancePerSvRewardCoupon"`
	UnclaimedAppRewards                  NUMERIC  `json:"unclaimedAppRewards"`
	UnclaimedValidatorRewards            NUMERIC  `json:"unclaimedValidatorRewards"`
	UnclaimedSvRewards                   NUMERIC  `json:"unclaimedSvRewards"`
	IssuancePerValidatorFaucetCoupon     NUMERIC  `json:"issuancePerValidatorFaucetCoupon"`
	OptAmuletsToIssueToDevelopmentFund   *NUMERIC `json:"optAmuletsToIssueToDevelopmentFund"`
}

// ToMap converts IssuingRoundParameters to a map for DAML arguments
func (t IssuingRoundParameters) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["issuancePerValidatorRewardCoupon"] = (*big.Int)(t.IssuancePerValidatorRewardCoupon)

	m["issuancePerFeaturedAppRewardCoupon"] = (*big.Int)(t.IssuancePerFeaturedAppRewardCoupon)

	m["issuancePerUnfeaturedAppRewardCoupon"] = (*big.Int)(t.IssuancePerUnfeaturedAppRewardCoupon)

	m["issuancePerSvRewardCoupon"] = (*big.Int)(t.IssuancePerSvRewardCoupon)

	m["unclaimedAppRewards"] = (*big.Int)(t.UnclaimedAppRewards)

	m["unclaimedValidatorRewards"] = (*big.Int)(t.UnclaimedValidatorRewards)

	m["unclaimedSvRewards"] = (*big.Int)(t.UnclaimedSvRewards)

	m["issuancePerValidatorFaucetCoupon"] = (*big.Int)(t.IssuancePerValidatorFaucetCoupon)

	if t.OptAmuletsToIssueToDevelopmentFund != nil {
		m["optAmuletsToIssueToDevelopmentFund"] = map[string]interface{}{
			"_type": "optional",
			"value": (*big.Int)(*t.OptAmuletsToIssueToDevelopmentFund),
		}
	} else {
		m["optAmuletsToIssueToDevelopmentFund"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for IssuingRoundParameters using JsonCodec
func (t IssuingRoundParameters) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for IssuingRoundParameters using JsonCodec
func (t *IssuingRoundParameters) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// LockedAmulet is a Template type
type LockedAmulet struct {
	Amulet Amulet   `json:"amulet"`
	Lock   TimeLock `json:"lock"`
}

// GetTemplateID returns the template ID for this template
func (t LockedAmulet) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "LockedAmulet")
}

// CreateCommand returns a CreateCommand for this template
func (t LockedAmulet) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	if t.Amulet != "" {
		args["amulet"] = t.Amulet
	}

	args["lock"] = t.Lock

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for LockedAmulet using JsonCodec
func (t LockedAmulet) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for LockedAmulet using JsonCodec
func (t *LockedAmulet) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for LockedAmulet

// LockedAmuletUnlock exercises the LockedAmulet_Unlock choice on this LockedAmulet contract
func (t LockedAmulet) LockedAmuletUnlock(contractID string, args LockedAmuletUnlock) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "LockedAmulet"),
		ContractID: contractID,
		Choice:     "LockedAmulet_Unlock",
		Arguments:  argsToMap(args),
	}
}

// LockedAmuletOwnerExpireLock exercises the LockedAmulet_OwnerExpireLock choice on this LockedAmulet contract
func (t LockedAmulet) LockedAmuletOwnerExpireLock(contractID string, args LockedAmuletOwnerExpireLock) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "LockedAmulet"),
		ContractID: contractID,
		Choice:     "LockedAmulet_OwnerExpireLock",
		Arguments:  argsToMap(args),
	}
}

// LockedAmuletExpireAmulet exercises the LockedAmulet_ExpireAmulet choice on this LockedAmulet contract
func (t LockedAmulet) LockedAmuletExpireAmulet(contractID string, args LockedAmuletExpireAmulet) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "LockedAmulet"),
		ContractID: contractID,
		Choice:     "LockedAmulet_ExpireAmulet",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this LockedAmulet contract
func (t LockedAmulet) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "LockedAmulet"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// Verify interface implementations for LockedAmulet

var _ IHolding = (*LockedAmulet)(nil)

// LockedAmuletExpireAmulet is a Record type
type LockedAmuletExpireAmulet struct {
	RoundCid CONTRACT_ID `json:"roundCid"`
}

// ToMap converts LockedAmuletExpireAmulet to a map for DAML arguments
func (t LockedAmuletExpireAmulet) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["roundCid"] = t.RoundCid

	return m
}

// MarshalJSON implements custom JSON marshaling for LockedAmuletExpireAmulet using JsonCodec
func (t LockedAmuletExpireAmulet) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for LockedAmuletExpireAmulet using JsonCodec
func (t *LockedAmuletExpireAmulet) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// LockedAmuletExpireAmuletResult is a Record type
type LockedAmuletExpireAmuletResult struct {
	ExpireSum AmuletExpireSummary `json:"expireSum"`
	Meta      *Metadata           `json:"meta"`
}

// ToMap converts LockedAmuletExpireAmuletResult to a map for DAML arguments
func (t LockedAmuletExpireAmuletResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["expireSum"] = t.ExpireSum

	if t.Meta != nil {
		m["meta"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.Meta,
		}
	} else {
		m["meta"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for LockedAmuletExpireAmuletResult using JsonCodec
func (t LockedAmuletExpireAmuletResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for LockedAmuletExpireAmuletResult using JsonCodec
func (t *LockedAmuletExpireAmuletResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// LockedAmuletOwnerExpireLock is a Record type
type LockedAmuletOwnerExpireLock struct {
	OpenRoundCid CONTRACT_ID `json:"openRoundCid"`
}

// ToMap converts LockedAmuletOwnerExpireLock to a map for DAML arguments
func (t LockedAmuletOwnerExpireLock) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["openRoundCid"] = t.OpenRoundCid

	return m
}

// MarshalJSON implements custom JSON marshaling for LockedAmuletOwnerExpireLock using JsonCodec
func (t LockedAmuletOwnerExpireLock) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for LockedAmuletOwnerExpireLock using JsonCodec
func (t *LockedAmuletOwnerExpireLock) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// LockedAmuletOwnerExpireLockResult is a Record type
type LockedAmuletOwnerExpireLockResult struct {
	AmuletSum AmuletCreateSummary `json:"amuletSum"`
	Meta      *Metadata           `json:"meta"`
}

// ToMap converts LockedAmuletOwnerExpireLockResult to a map for DAML arguments
func (t LockedAmuletOwnerExpireLockResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["amuletSum"] = t.AmuletSum

	if t.Meta != nil {
		m["meta"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.Meta,
		}
	} else {
		m["meta"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for LockedAmuletOwnerExpireLockResult using JsonCodec
func (t LockedAmuletOwnerExpireLockResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for LockedAmuletOwnerExpireLockResult using JsonCodec
func (t *LockedAmuletOwnerExpireLockResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// LockedAmuletUnlock is a Record type
type LockedAmuletUnlock struct {
	OpenRoundCid CONTRACT_ID `json:"openRoundCid"`
}

// ToMap converts LockedAmuletUnlock to a map for DAML arguments
func (t LockedAmuletUnlock) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["openRoundCid"] = t.OpenRoundCid

	return m
}

// MarshalJSON implements custom JSON marshaling for LockedAmuletUnlock using JsonCodec
func (t LockedAmuletUnlock) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for LockedAmuletUnlock using JsonCodec
func (t *LockedAmuletUnlock) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// LockedAmuletUnlockResult is a Record type
type LockedAmuletUnlockResult struct {
	AmuletSum AmuletCreateSummary `json:"amuletSum"`
	Meta      *Metadata           `json:"meta"`
}

// ToMap converts LockedAmuletUnlockResult to a map for DAML arguments
func (t LockedAmuletUnlockResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["amuletSum"] = t.AmuletSum

	if t.Meta != nil {
		m["meta"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.Meta,
		}
	} else {
		m["meta"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for LockedAmuletUnlockResult using JsonCodec
func (t LockedAmuletUnlockResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for LockedAmuletUnlockResult using JsonCodec
func (t *LockedAmuletUnlockResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// MemberTraffic is a Template type
type MemberTraffic struct {
	Dso            PARTY   `json:"dso"`
	MemberId       TEXT    `json:"memberId"`
	SynchronizerId TEXT    `json:"synchronizerId"`
	MigrationId    INT64   `json:"migrationId"`
	TotalPurchased INT64   `json:"totalPurchased"`
	NumPurchases   INT64   `json:"numPurchases"`
	AmuletSpent    NUMERIC `json:"amuletSpent"`
	UsdSpent       NUMERIC `json:"usdSpent"`
}

// GetTemplateID returns the template ID for this template
func (t MemberTraffic) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.DecentralizedSynchronizer", "MemberTraffic")
}

// CreateCommand returns a CreateCommand for this template
func (t MemberTraffic) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["dso"] = t.Dso.ToMap()

	args["memberId"] = string(t.MemberId)

	args["synchronizerId"] = string(t.SynchronizerId)

	args["migrationId"] = int64(t.MigrationId)

	args["totalPurchased"] = int64(t.TotalPurchased)

	args["numPurchases"] = int64(t.NumPurchases)

	if t.AmuletSpent != nil {
		args["amuletSpent"] = (*big.Int)(t.AmuletSpent)
	}

	if t.UsdSpent != nil {
		args["usdSpent"] = (*big.Int)(t.UsdSpent)
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for MemberTraffic using JsonCodec
func (t MemberTraffic) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for MemberTraffic using JsonCodec
func (t *MemberTraffic) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for MemberTraffic

// Archive exercises the Archive choice on this MemberTraffic contract
func (t MemberTraffic) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.DecentralizedSynchronizer", "MemberTraffic"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// OpenMiningRound is a Template type
type OpenMiningRound struct {
	Dso               PARTY          `json:"dso"`
	Round             Round          `json:"round"`
	AmuletPrice       NUMERIC        `json:"amuletPrice"`
	OpensAt           TIMESTAMP      `json:"opensAt"`
	TargetClosesAt    TIMESTAMP      `json:"targetClosesAt"`
	IssuingFor        RELTIME        `json:"issuingFor"`
	TransferConfigUsd TransferConfig `json:"transferConfigUsd"`
	IssuanceConfig    IssuanceConfig `json:"issuanceConfig"`
	TickDuration      RELTIME        `json:"tickDuration"`
}

// GetTemplateID returns the template ID for this template
func (t OpenMiningRound) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Round", "OpenMiningRound")
}

// CreateCommand returns a CreateCommand for this template
func (t OpenMiningRound) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["dso"] = t.Dso.ToMap()

	args["round"] = t.Round

	if t.AmuletPrice != nil {
		args["amuletPrice"] = (*big.Int)(t.AmuletPrice)
	}

	args["opensAt"] = t.OpensAt

	args["targetClosesAt"] = t.TargetClosesAt

	args["issuingFor"] = t.IssuingFor

	args["transferConfigUsd"] = t.TransferConfigUsd

	args["issuanceConfig"] = t.IssuanceConfig

	args["tickDuration"] = t.TickDuration

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for OpenMiningRound using JsonCodec
func (t OpenMiningRound) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for OpenMiningRound using JsonCodec
func (t *OpenMiningRound) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for OpenMiningRound

// Archive exercises the Archive choice on this OpenMiningRound contract
func (t OpenMiningRound) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Round", "OpenMiningRound"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// OpenMiningRoundFetch exercises the OpenMiningRound_Fetch choice on this OpenMiningRound contract
func (t OpenMiningRound) OpenMiningRoundFetch(contractID string, args OpenMiningRoundFetch) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Round", "OpenMiningRound"),
		ContractID: contractID,
		Choice:     "OpenMiningRound_Fetch",
		Arguments:  argsToMap(args),
	}
}

// OpenMiningRoundSummary is a Record type
type OpenMiningRoundSummary struct {
	TotalValidatorRewardCoupons     NUMERIC `json:"totalValidatorRewardCoupons"`
	TotalFeaturedAppRewardCoupons   NUMERIC `json:"totalFeaturedAppRewardCoupons"`
	TotalUnfeaturedAppRewardCoupons NUMERIC `json:"totalUnfeaturedAppRewardCoupons"`
	TotalSvRewardWeight             INT64   `json:"totalSvRewardWeight"`
	OptTotalValidatorFaucetCoupons  *INT64  `json:"optTotalValidatorFaucetCoupons"`
}

// ToMap converts OpenMiningRoundSummary to a map for DAML arguments
func (t OpenMiningRoundSummary) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["totalValidatorRewardCoupons"] = (*big.Int)(t.TotalValidatorRewardCoupons)

	m["totalFeaturedAppRewardCoupons"] = (*big.Int)(t.TotalFeaturedAppRewardCoupons)

	m["totalUnfeaturedAppRewardCoupons"] = (*big.Int)(t.TotalUnfeaturedAppRewardCoupons)

	m["totalSvRewardWeight"] = int64(t.TotalSvRewardWeight)

	if t.OptTotalValidatorFaucetCoupons != nil {
		m["optTotalValidatorFaucetCoupons"] = map[string]interface{}{
			"_type": "optional",
			"value": int64(*t.OptTotalValidatorFaucetCoupons),
		}
	} else {
		m["optTotalValidatorFaucetCoupons"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for OpenMiningRoundSummary using JsonCodec
func (t OpenMiningRoundSummary) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for OpenMiningRoundSummary using JsonCodec
func (t *OpenMiningRoundSummary) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// OpenMiningRoundFetch is a Record type
type OpenMiningRoundFetch struct {
	P PARTY `json:"p"`
}

// ToMap converts OpenMiningRoundFetch to a map for DAML arguments
func (t OpenMiningRoundFetch) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["p"] = t.P.ToMap()

	return m
}

// MarshalJSON implements custom JSON marshaling for OpenMiningRoundFetch using JsonCodec
func (t OpenMiningRoundFetch) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for OpenMiningRoundFetch using JsonCodec
func (t *OpenMiningRoundFetch) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// PackageConfig is a Record type
type PackageConfig struct {
	Amulet             TEXT `json:"amulet"`
	AmuletNameService  TEXT `json:"amuletNameService"`
	DsoGovernance      TEXT `json:"dsoGovernance"`
	ValidatorLifecycle TEXT `json:"validatorLifecycle"`
	Wallet             TEXT `json:"wallet"`
	WalletPayments     TEXT `json:"walletPayments"`
}

// ToMap converts PackageConfig to a map for DAML arguments
func (t PackageConfig) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["amulet"] = string(t.Amulet)

	m["amuletNameService"] = string(t.AmuletNameService)

	m["dsoGovernance"] = string(t.DsoGovernance)

	m["validatorLifecycle"] = string(t.ValidatorLifecycle)

	m["wallet"] = string(t.Wallet)

	m["walletPayments"] = string(t.WalletPayments)

	return m
}

// MarshalJSON implements custom JSON marshaling for PackageConfig using JsonCodec
func (t PackageConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for PackageConfig using JsonCodec
func (t *PackageConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// PaymentTransferContext is a Record type
type PaymentTransferContext struct {
	AmuletRules CONTRACT_ID     `json:"amuletRules"`
	Context     TransferContext `json:"context"`
}

// ToMap converts PaymentTransferContext to a map for DAML arguments
func (t PaymentTransferContext) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["amuletRules"] = t.AmuletRules

	m["context"] = t.Context

	return m
}

// MarshalJSON implements custom JSON marshaling for PaymentTransferContext using JsonCodec
func (t PaymentTransferContext) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for PaymentTransferContext using JsonCodec
func (t *PaymentTransferContext) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// PreprocessedTransferOutput is a Record type
type PreprocessedTransferOutput struct {
	Owner     PARTY     `json:"owner"`
	OutputFee NUMERIC   `json:"outputFee"`
	Amount    NUMERIC   `json:"amount"`
	Lock      *TimeLock `json:"lock"`
}

// ToMap converts PreprocessedTransferOutput to a map for DAML arguments
func (t PreprocessedTransferOutput) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["owner"] = t.Owner.ToMap()

	m["outputFee"] = (*big.Int)(t.OutputFee)

	m["amount"] = (*big.Int)(t.Amount)

	if t.Lock != nil {
		m["lock"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.Lock,
		}
	} else {
		m["lock"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for PreprocessedTransferOutput using JsonCodec
func (t PreprocessedTransferOutput) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for PreprocessedTransferOutput using JsonCodec
func (t *PreprocessedTransferOutput) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// RatePerDay is a Record type
type RatePerDay struct {
	Rate NUMERIC `json:"rate"`
}

// ToMap converts RatePerDay to a map for DAML arguments
func (t RatePerDay) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["rate"] = (*big.Int)(t.Rate)

	return m
}

// MarshalJSON implements custom JSON marshaling for RatePerDay using JsonCodec
func (t RatePerDay) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for RatePerDay using JsonCodec
func (t *RatePerDay) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// RatePerRound is a Record type
type RatePerRound struct {
	Rate NUMERIC `json:"rate"`
}

// ToMap converts RatePerRound to a map for DAML arguments
func (t RatePerRound) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["rate"] = (*big.Int)(t.Rate)

	return m
}

// MarshalJSON implements custom JSON marshaling for RatePerRound using JsonCodec
func (t RatePerRound) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for RatePerRound using JsonCodec
func (t *RatePerRound) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// RelRound is a Record type
type RelRound struct {
	Diff INT64 `json:"diff"`
}

// ToMap converts RelRound to a map for DAML arguments
func (t RelRound) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["diff"] = int64(t.Diff)

	return m
}

// MarshalJSON implements custom JSON marshaling for RelRound using JsonCodec
func (t RelRound) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for RelRound using JsonCodec
func (t *RelRound) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// RewardsIssuanceConfig is a Record type
type RewardsIssuanceConfig struct {
	IssueAppRewards       BOOL `json:"issueAppRewards"`
	IssueValidatorRewards BOOL `json:"issueValidatorRewards"`
}

// ToMap converts RewardsIssuanceConfig to a map for DAML arguments
func (t RewardsIssuanceConfig) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["issueAppRewards"] = bool(t.IssueAppRewards)

	m["issueValidatorRewards"] = bool(t.IssueValidatorRewards)

	return m
}

// MarshalJSON implements custom JSON marshaling for RewardsIssuanceConfig using JsonCodec
func (t RewardsIssuanceConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for RewardsIssuanceConfig using JsonCodec
func (t *RewardsIssuanceConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Round is a Record type
type Round struct {
	Number INT64 `json:"number"`
}

// ToMap converts Round to a map for DAML arguments
func (t Round) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["number"] = int64(t.Number)

	return m
}

// MarshalJSON implements custom JSON marshaling for Round using JsonCodec
func (t Round) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for Round using JsonCodec
func (t *Round) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Schedule is a Record type
type Schedule struct {
	InitialValue interface{}                        `json:"initialValue"`
	FutureValues []TUPLE2[interface{}, interface{}] `json:"futureValues"`
}

// ToMap converts Schedule to a map for DAML arguments
func (t Schedule) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["initialValue"] = t.InitialValue

	m["futureValues"] = t.FutureValues

	return m
}

// MarshalJSON implements custom JSON marshaling for Schedule using JsonCodec
func (t Schedule) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for Schedule using JsonCodec
func (t *Schedule) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// SteppedRate is a Record type
type SteppedRate struct {
	InitialRate NUMERIC                    `json:"initialRate"`
	Steps       []TUPLE2[NUMERIC, NUMERIC] `json:"steps"`
}

// ToMap converts SteppedRate to a map for DAML arguments
func (t SteppedRate) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["initialRate"] = (*big.Int)(t.InitialRate)

	m["steps"] = t.Steps

	return m
}

// MarshalJSON implements custom JSON marshaling for SteppedRate using JsonCodec
func (t SteppedRate) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for SteppedRate using JsonCodec
func (t *SteppedRate) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// SummarizingMiningRound is a Template type
type SummarizingMiningRound struct {
	Dso            PARTY          `json:"dso"`
	Round          Round          `json:"round"`
	AmuletPrice    NUMERIC        `json:"amuletPrice"`
	IssuanceConfig IssuanceConfig `json:"issuanceConfig"`
	TickDuration   RELTIME        `json:"tickDuration"`
}

// GetTemplateID returns the template ID for this template
func (t SummarizingMiningRound) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Round", "SummarizingMiningRound")
}

// CreateCommand returns a CreateCommand for this template
func (t SummarizingMiningRound) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["dso"] = t.Dso.ToMap()

	args["round"] = t.Round

	if t.AmuletPrice != nil {
		args["amuletPrice"] = (*big.Int)(t.AmuletPrice)
	}

	args["issuanceConfig"] = t.IssuanceConfig

	args["tickDuration"] = t.TickDuration

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for SummarizingMiningRound using JsonCodec
func (t SummarizingMiningRound) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for SummarizingMiningRound using JsonCodec
func (t *SummarizingMiningRound) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for SummarizingMiningRound

// Archive exercises the Archive choice on this SummarizingMiningRound contract
func (t SummarizingMiningRound) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Round", "SummarizingMiningRound"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// SvRewardCoupon is a Template type
type SvRewardCoupon struct {
	Dso         PARTY `json:"dso"`
	Sv          PARTY `json:"sv"`
	Beneficiary PARTY `json:"beneficiary"`
	Round       Round `json:"round"`
	Weight      INT64 `json:"weight"`
}

// GetTemplateID returns the template ID for this template
func (t SvRewardCoupon) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "SvRewardCoupon")
}

// CreateCommand returns a CreateCommand for this template
func (t SvRewardCoupon) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["dso"] = t.Dso.ToMap()

	args["sv"] = t.Sv.ToMap()

	args["beneficiary"] = t.Beneficiary.ToMap()

	args["round"] = t.Round

	args["weight"] = int64(t.Weight)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for SvRewardCoupon using JsonCodec
func (t SvRewardCoupon) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for SvRewardCoupon using JsonCodec
func (t *SvRewardCoupon) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for SvRewardCoupon

// SvRewardCouponDsoExpire exercises the SvRewardCoupon_DsoExpire choice on this SvRewardCoupon contract
func (t SvRewardCoupon) SvRewardCouponDsoExpire(contractID string, args SvRewardCouponDsoExpire) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "SvRewardCoupon"),
		ContractID: contractID,
		Choice:     "SvRewardCoupon_DsoExpire",
		Arguments:  argsToMap(args),
	}
}

// SvRewardCouponArchiveAsBeneficiary exercises the SvRewardCoupon_ArchiveAsBeneficiary choice on this SvRewardCoupon contract
func (t SvRewardCoupon) SvRewardCouponArchiveAsBeneficiary(contractID string, args SvRewardCouponArchiveAsBeneficiary) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "SvRewardCoupon"),
		ContractID: contractID,
		Choice:     "SvRewardCoupon_ArchiveAsBeneficiary",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this SvRewardCoupon contract
func (t SvRewardCoupon) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "SvRewardCoupon"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// SvRewardCouponArchiveAsBeneficiary is a Record type
type SvRewardCouponArchiveAsBeneficiary struct{}

// ToMap converts SvRewardCouponArchiveAsBeneficiary to a map for DAML arguments
func (t SvRewardCouponArchiveAsBeneficiary) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	return m
}

// MarshalJSON implements custom JSON marshaling for SvRewardCouponArchiveAsBeneficiary using JsonCodec
func (t SvRewardCouponArchiveAsBeneficiary) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for SvRewardCouponArchiveAsBeneficiary using JsonCodec
func (t *SvRewardCouponArchiveAsBeneficiary) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// SvRewardCouponArchiveAsBeneficiaryResult is an enum type
type SvRewardCouponArchiveAsBeneficiaryResult string

const (
	SvRewardCouponArchiveAsBeneficiaryResultSvRewardCoupon_ArchiveAsBeneficiaryResult SvRewardCouponArchiveAsBeneficiaryResult = "SvRewardCoupon_ArchiveAsBeneficiaryResult"
)

// GetEnumConstructor implements types.ENUM interface
func (e SvRewardCouponArchiveAsBeneficiaryResult) GetEnumConstructor() string {
	return string(e)
}

// GetEnumTypeID implements types.ENUM interface
func (e SvRewardCouponArchiveAsBeneficiaryResult) GetEnumTypeID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "SvRewardCouponArchiveAsBeneficiaryResult")
}

// MarshalJSON implements custom JSON marshaling for SvRewardCouponArchiveAsBeneficiaryResult using JsonCodec
func (e SvRewardCouponArchiveAsBeneficiaryResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(e)
}

// UnmarshalJSON implements custom JSON unmarshaling for SvRewardCouponArchiveAsBeneficiaryResult using JsonCodec
func (e *SvRewardCouponArchiveAsBeneficiaryResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, e)
}

// Verify interface implementation
var _ ENUM = SvRewardCouponArchiveAsBeneficiaryResult("")

// SvRewardCouponDsoExpire is a Record type
type SvRewardCouponDsoExpire struct {
	ClosedRoundCid CONTRACT_ID `json:"closedRoundCid"`
}

// ToMap converts SvRewardCouponDsoExpire to a map for DAML arguments
func (t SvRewardCouponDsoExpire) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["closedRoundCid"] = t.ClosedRoundCid

	return m
}

// MarshalJSON implements custom JSON marshaling for SvRewardCouponDsoExpire using JsonCodec
func (t SvRewardCouponDsoExpire) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for SvRewardCouponDsoExpire using JsonCodec
func (t *SvRewardCouponDsoExpire) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// SvRewardCouponDsoExpireResult is a Record type
type SvRewardCouponDsoExpireResult struct {
	Weight INT64 `json:"weight"`
}

// ToMap converts SvRewardCouponDsoExpireResult to a map for DAML arguments
func (t SvRewardCouponDsoExpireResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["weight"] = int64(t.Weight)

	return m
}

// MarshalJSON implements custom JSON marshaling for SvRewardCouponDsoExpireResult using JsonCodec
func (t SvRewardCouponDsoExpireResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for SvRewardCouponDsoExpireResult using JsonCodec
func (t *SvRewardCouponDsoExpireResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// SynchronizerFeesConfig is a Record type
type SynchronizerFeesConfig struct {
	BaseRateTrafficLimits    BaseRateTrafficLimits `json:"baseRateTrafficLimits"`
	ExtraTrafficPrice        NUMERIC               `json:"extraTrafficPrice"`
	ReadVsWriteScalingFactor INT64                 `json:"readVsWriteScalingFactor"`
	MinTopupAmount           INT64                 `json:"minTopupAmount"`
}

// ToMap converts SynchronizerFeesConfig to a map for DAML arguments
func (t SynchronizerFeesConfig) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["baseRateTrafficLimits"] = t.BaseRateTrafficLimits

	m["extraTrafficPrice"] = (*big.Int)(t.ExtraTrafficPrice)

	m["readVsWriteScalingFactor"] = int64(t.ReadVsWriteScalingFactor)

	m["minTopupAmount"] = int64(t.MinTopupAmount)

	return m
}

// MarshalJSON implements custom JSON marshaling for SynchronizerFeesConfig using JsonCodec
func (t SynchronizerFeesConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for SynchronizerFeesConfig using JsonCodec
func (t *SynchronizerFeesConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TimeLock is a Record type
type TimeLock struct {
	Holders    []PARTY   `json:"holders"`
	ExpiresAt  TIMESTAMP `json:"expiresAt"`
	OptContext *TEXT     `json:"optContext"`
}

// ToMap converts TimeLock to a map for DAML arguments
func (t TimeLock) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["holders"] = t.Holders

	m["expiresAt"] = t.ExpiresAt

	if t.OptContext != nil {
		m["optContext"] = map[string]interface{}{
			"_type": "optional",
			"value": string(*t.OptContext),
		}
	} else {
		m["optContext"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for TimeLock using JsonCodec
func (t TimeLock) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TimeLock using JsonCodec
func (t *TimeLock) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Transfer is a Record type
type Transfer struct {
	Sender        PARTY                   `json:"sender"`
	Provider      PARTY                   `json:"provider"`
	Inputs        []TransferInput         `json:"inputs"`
	Outputs       []TransferOutput        `json:"outputs"`
	Beneficiaries *[]AppRewardBeneficiary `json:"beneficiaries"`
}

// ToMap converts Transfer to a map for DAML arguments
func (t Transfer) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["sender"] = t.Sender.ToMap()

	m["provider"] = t.Provider.ToMap()

	m["inputs"] = t.Inputs

	m["outputs"] = t.Outputs

	if t.Beneficiaries != nil {
		m["beneficiaries"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.Beneficiaries,
		}
	} else {
		m["beneficiaries"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for Transfer using JsonCodec
func (t Transfer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for Transfer using JsonCodec
func (t *Transfer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferCommand is a Template type
type TransferCommand struct {
	Dso         PARTY     `json:"dso"`
	Sender      PARTY     `json:"sender"`
	Receiver    PARTY     `json:"receiver"`
	Delegate    PARTY     `json:"delegate"`
	Amount      NUMERIC   `json:"amount"`
	ExpiresAt   TIMESTAMP `json:"expiresAt"`
	Nonce       INT64     `json:"nonce"`
	Description *TEXT     `json:"description"`
}

// GetTemplateID returns the template ID for this template
func (t TransferCommand) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.ExternalPartyAmuletRules", "TransferCommand")
}

// CreateCommand returns a CreateCommand for this template
func (t TransferCommand) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["dso"] = t.Dso.ToMap()

	args["sender"] = t.Sender.ToMap()

	args["receiver"] = t.Receiver.ToMap()

	args["delegate"] = t.Delegate.ToMap()

	if t.Amount != nil {
		args["amount"] = (*big.Int)(t.Amount)
	}

	args["expiresAt"] = t.ExpiresAt

	args["nonce"] = int64(t.Nonce)

	if t.Description != nil {
		args["description"] = map[string]interface{}{
			"_type": "optional",
			"value": string(*t.Description),
		}
	} else {
		args["description"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for TransferCommand using JsonCodec
func (t TransferCommand) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferCommand using JsonCodec
func (t *TransferCommand) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for TransferCommand

// TransferCommandExpire exercises the TransferCommand_Expire choice on this TransferCommand contract
func (t TransferCommand) TransferCommandExpire(contractID string, args TransferCommandExpire) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.ExternalPartyAmuletRules", "TransferCommand"),
		ContractID: contractID,
		Choice:     "TransferCommand_Expire",
		Arguments:  argsToMap(args),
	}
}

// TransferCommandSend exercises the TransferCommand_Send choice on this TransferCommand contract
func (t TransferCommand) TransferCommandSend(contractID string, args TransferCommandSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.ExternalPartyAmuletRules", "TransferCommand"),
		ContractID: contractID,
		Choice:     "TransferCommand_Send",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this TransferCommand contract
func (t TransferCommand) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.ExternalPartyAmuletRules", "TransferCommand"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// TransferCommandWithdraw exercises the TransferCommand_Withdraw choice on this TransferCommand contract
func (t TransferCommand) TransferCommandWithdraw(contractID string, args TransferCommandWithdraw) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.ExternalPartyAmuletRules", "TransferCommand"),
		ContractID: contractID,
		Choice:     "TransferCommand_Withdraw",
		Arguments:  argsToMap(args),
	}
}

// TransferCommandCounter is a Template type
type TransferCommandCounter struct {
	Dso       PARTY `json:"dso"`
	Sender    PARTY `json:"sender"`
	NextNonce INT64 `json:"nextNonce"`
}

// GetTemplateID returns the template ID for this template
func (t TransferCommandCounter) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.ExternalPartyAmuletRules", "TransferCommandCounter")
}

// CreateCommand returns a CreateCommand for this template
func (t TransferCommandCounter) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["dso"] = t.Dso.ToMap()

	args["sender"] = t.Sender.ToMap()

	args["nextNonce"] = int64(t.NextNonce)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for TransferCommandCounter using JsonCodec
func (t TransferCommandCounter) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferCommandCounter using JsonCodec
func (t *TransferCommandCounter) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for TransferCommandCounter

// Archive exercises the Archive choice on this TransferCommandCounter contract
func (t TransferCommandCounter) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.ExternalPartyAmuletRules", "TransferCommandCounter"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// TransferCommandResult is a variant/union type
type TransferCommandResult struct {
	TransferCommandResultFailure *TransferCommandResultFailure `json:"TransferCommandResultFailure,omitempty"`
	TransferCommandResultSuccess *TransferCommandResultSuccess `json:"TransferCommandResultSuccess,omitempty"`
}

// MarshalJSON implements custom JSON marshaling for TransferCommandResult
func (v TransferCommandResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(v)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferCommandResult
func (v *TransferCommandResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, v)
}

// GetVariantTag implements types.VARIANT interface
func (v TransferCommandResult) GetVariantTag() string {
	if v.TransferCommandResultFailure != nil {
		return "TransferCommandResultFailure"
	}

	if v.TransferCommandResultSuccess != nil {
		return "TransferCommandResultSuccess"
	}

	return ""
}

// GetVariantValue implements types.VARIANT interface
func (v TransferCommandResult) GetVariantValue() interface{} {
	if v.TransferCommandResultFailure != nil {
		return v.TransferCommandResultFailure
	}

	if v.TransferCommandResultSuccess != nil {
		return v.TransferCommandResultSuccess
	}

	return nil
}

// Verify interface implementation
var _ VARIANT = (*TransferCommandResult)(nil)

// TransferCommandResultFailure is a Record type
type TransferCommandResultFailure struct {
	Reason InvalidTransferReason `json:"reason"`
}

// ToMap converts TransferCommandResultFailure to a map for DAML arguments
func (t TransferCommandResultFailure) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["reason"] = t.Reason

	return m
}

// MarshalJSON implements custom JSON marshaling for TransferCommandResultFailure using JsonCodec
func (t TransferCommandResultFailure) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferCommandResultFailure using JsonCodec
func (t *TransferCommandResultFailure) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferCommandResultSuccess is a Record type
type TransferCommandResultSuccess struct {
	Result TransferResult `json:"result"`
}

// ToMap converts TransferCommandResultSuccess to a map for DAML arguments
func (t TransferCommandResultSuccess) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["result"] = t.Result

	return m
}

// MarshalJSON implements custom JSON marshaling for TransferCommandResultSuccess using JsonCodec
func (t TransferCommandResultSuccess) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferCommandResultSuccess using JsonCodec
func (t *TransferCommandResultSuccess) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferCommandExpire is a Record type
type TransferCommandExpire struct {
	P PARTY `json:"p"`
}

// ToMap converts TransferCommandExpire to a map for DAML arguments
func (t TransferCommandExpire) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["p"] = t.P.ToMap()

	return m
}

// MarshalJSON implements custom JSON marshaling for TransferCommandExpire using JsonCodec
func (t TransferCommandExpire) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferCommandExpire using JsonCodec
func (t *TransferCommandExpire) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferCommandExpireResult is a Record type
type TransferCommandExpireResult struct {
	Sender PARTY `json:"sender"`
	Nonce  INT64 `json:"nonce"`
}

// ToMap converts TransferCommandExpireResult to a map for DAML arguments
func (t TransferCommandExpireResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["sender"] = t.Sender.ToMap()

	m["nonce"] = int64(t.Nonce)

	return m
}

// MarshalJSON implements custom JSON marshaling for TransferCommandExpireResult using JsonCodec
func (t TransferCommandExpireResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferCommandExpireResult using JsonCodec
func (t *TransferCommandExpireResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferCommandSend is a Record type
type TransferCommandSend struct {
	Context                 PaymentTransferContext `json:"context"`
	Inputs                  []TransferInput        `json:"inputs"`
	TransferPreapprovalCidO *CONTRACT_ID           `json:"transferPreapprovalCidO"`
	TransferCounterCid      CONTRACT_ID            `json:"transferCounterCid"`
}

// ToMap converts TransferCommandSend to a map for DAML arguments
func (t TransferCommandSend) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["context"] = t.Context

	m["inputs"] = t.Inputs

	if t.TransferPreapprovalCidO != nil {
		m["transferPreapprovalCidO"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.TransferPreapprovalCidO,
		}
	} else {
		m["transferPreapprovalCidO"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	m["transferCounterCid"] = t.TransferCounterCid

	return m
}

// MarshalJSON implements custom JSON marshaling for TransferCommandSend using JsonCodec
func (t TransferCommandSend) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferCommandSend using JsonCodec
func (t *TransferCommandSend) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferCommandSendResult is a Record type
type TransferCommandSendResult struct {
	Result TransferCommandResult `json:"result"`
	Sender PARTY                 `json:"sender"`
	Nonce  INT64                 `json:"nonce"`
}

// ToMap converts TransferCommandSendResult to a map for DAML arguments
func (t TransferCommandSendResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["result"] = t.Result

	m["sender"] = t.Sender.ToMap()

	m["nonce"] = int64(t.Nonce)

	return m
}

// MarshalJSON implements custom JSON marshaling for TransferCommandSendResult using JsonCodec
func (t TransferCommandSendResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferCommandSendResult using JsonCodec
func (t *TransferCommandSendResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferCommandWithdraw is a Record type
type TransferCommandWithdraw struct{}

// ToMap converts TransferCommandWithdraw to a map for DAML arguments
func (t TransferCommandWithdraw) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	return m
}

// MarshalJSON implements custom JSON marshaling for TransferCommandWithdraw using JsonCodec
func (t TransferCommandWithdraw) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferCommandWithdraw using JsonCodec
func (t *TransferCommandWithdraw) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferCommandWithdrawResult is a Record type
type TransferCommandWithdrawResult struct {
	Sender PARTY `json:"sender"`
	Nonce  INT64 `json:"nonce"`
}

// ToMap converts TransferCommandWithdrawResult to a map for DAML arguments
func (t TransferCommandWithdrawResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["sender"] = t.Sender.ToMap()

	m["nonce"] = int64(t.Nonce)

	return m
}

// MarshalJSON implements custom JSON marshaling for TransferCommandWithdrawResult using JsonCodec
func (t TransferCommandWithdrawResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferCommandWithdrawResult using JsonCodec
func (t *TransferCommandWithdrawResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferConfig is a Record type
type TransferConfig struct {
	CreateFee                    FixedFee     `json:"createFee"`
	HoldingFee                   RatePerRound `json:"holdingFee"`
	TransferFee                  SteppedRate  `json:"transferFee"`
	LockHolderFee                FixedFee     `json:"lockHolderFee"`
	ExtraFeaturedAppRewardAmount NUMERIC      `json:"extraFeaturedAppRewardAmount"`
	MaxNumInputs                 INT64        `json:"maxNumInputs"`
	MaxNumOutputs                INT64        `json:"maxNumOutputs"`
	MaxNumLockHolders            INT64        `json:"maxNumLockHolders"`
}

// ToMap converts TransferConfig to a map for DAML arguments
func (t TransferConfig) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["createFee"] = t.CreateFee

	m["holdingFee"] = t.HoldingFee

	m["transferFee"] = t.TransferFee

	m["lockHolderFee"] = t.LockHolderFee

	m["extraFeaturedAppRewardAmount"] = (*big.Int)(t.ExtraFeaturedAppRewardAmount)

	m["maxNumInputs"] = int64(t.MaxNumInputs)

	m["maxNumOutputs"] = int64(t.MaxNumOutputs)

	m["maxNumLockHolders"] = int64(t.MaxNumLockHolders)

	return m
}

// MarshalJSON implements custom JSON marshaling for TransferConfig using JsonCodec
func (t TransferConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferConfig using JsonCodec
func (t *TransferConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferContext is a Record type
type TransferContext struct {
	OpenMiningRound     CONTRACT_ID  `json:"openMiningRound"`
	IssuingMiningRounds GENMAP       `json:"issuingMiningRounds"`
	ValidatorRights     GENMAP       `json:"validatorRights"`
	FeaturedAppRight    *CONTRACT_ID `json:"featuredAppRight"`
}

// ToMap converts TransferContext to a map for DAML arguments
func (t TransferContext) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["openMiningRound"] = t.OpenMiningRound

	m["issuingMiningRounds"] = map[string]interface{}{"_type": "genmap", "value": t.IssuingMiningRounds}

	m["validatorRights"] = map[string]interface{}{"_type": "genmap", "value": t.ValidatorRights}

	if t.FeaturedAppRight != nil {
		m["featuredAppRight"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.FeaturedAppRight,
		}
	} else {
		m["featuredAppRight"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for TransferContext using JsonCodec
func (t TransferContext) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferContext using JsonCodec
func (t *TransferContext) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferContextSummary is a Record type
type TransferContextSummary struct {
	FeaturedAppProvider *PARTY          `json:"featuredAppProvider"`
	Config              TransferConfig  `json:"config"`
	OpenRound           OpenMiningRound `json:"openRound"`
	IssuingMiningRounds GENMAP          `json:"issuingMiningRounds"`
	ValidatorRights     GENMAP          `json:"validatorRights"`
}

// ToMap converts TransferContextSummary to a map for DAML arguments
func (t TransferContextSummary) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	if t.FeaturedAppProvider != nil {
		m["featuredAppProvider"] = map[string]interface{}{
			"_type": "optional",
			"value": (*t.FeaturedAppProvider).ToMap(),
		}
	} else {
		m["featuredAppProvider"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	m["config"] = t.Config

	m["openRound"] = t.OpenRound

	m["issuingMiningRounds"] = map[string]interface{}{"_type": "genmap", "value": t.IssuingMiningRounds}

	m["validatorRights"] = map[string]interface{}{"_type": "genmap", "value": t.ValidatorRights}

	return m
}

// MarshalJSON implements custom JSON marshaling for TransferContextSummary using JsonCodec
func (t TransferContextSummary) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferContextSummary using JsonCodec
func (t *TransferContextSummary) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferInput is a variant/union type
type TransferInput struct {
	InputAppRewardCoupon                 *CONTRACT_ID      `json:"InputAppRewardCoupon,omitempty"`
	InputValidatorRewardCoupon           *CONTRACT_ID      `json:"InputValidatorRewardCoupon,omitempty"`
	InputSvRewardCoupon                  *CONTRACT_ID      `json:"InputSvRewardCoupon,omitempty"`
	InputAmulet                          *CONTRACT_ID      `json:"InputAmulet,omitempty"`
	ExtTransferInput                     *ExtTransferInput `json:"ExtTransferInput,omitempty"`
	InputValidatorLivenessActivityRecord *CONTRACT_ID      `json:"InputValidatorLivenessActivityRecord,omitempty"`
	InputUnclaimedActivityRecord         *CONTRACT_ID      `json:"InputUnclaimedActivityRecord,omitempty"`
	InputDevelopmentFundCoupon           *CONTRACT_ID      `json:"InputDevelopmentFundCoupon,omitempty"`
}

// MarshalJSON implements custom JSON marshaling for TransferInput
func (v TransferInput) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(v)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferInput
func (v *TransferInput) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, v)
}

// GetVariantTag implements types.VARIANT interface
func (v TransferInput) GetVariantTag() string {
	if v.InputAppRewardCoupon != nil {
		return "InputAppRewardCoupon"
	}

	if v.InputValidatorRewardCoupon != nil {
		return "InputValidatorRewardCoupon"
	}

	if v.InputSvRewardCoupon != nil {
		return "InputSvRewardCoupon"
	}

	if v.InputAmulet != nil {
		return "InputAmulet"
	}

	if v.ExtTransferInput != nil {
		return "ExtTransferInput"
	}

	if v.InputValidatorLivenessActivityRecord != nil {
		return "InputValidatorLivenessActivityRecord"
	}

	if v.InputUnclaimedActivityRecord != nil {
		return "InputUnclaimedActivityRecord"
	}

	if v.InputDevelopmentFundCoupon != nil {
		return "InputDevelopmentFundCoupon"
	}

	return ""
}

// GetVariantValue implements types.VARIANT interface
func (v TransferInput) GetVariantValue() interface{} {
	if v.InputAppRewardCoupon != nil {
		return v.InputAppRewardCoupon
	}

	if v.InputValidatorRewardCoupon != nil {
		return v.InputValidatorRewardCoupon
	}

	if v.InputSvRewardCoupon != nil {
		return v.InputSvRewardCoupon
	}

	if v.InputAmulet != nil {
		return v.InputAmulet
	}

	if v.ExtTransferInput != nil {
		return v.ExtTransferInput
	}

	if v.InputValidatorLivenessActivityRecord != nil {
		return v.InputValidatorLivenessActivityRecord
	}

	if v.InputUnclaimedActivityRecord != nil {
		return v.InputUnclaimedActivityRecord
	}

	if v.InputDevelopmentFundCoupon != nil {
		return v.InputDevelopmentFundCoupon
	}

	return nil
}

// Verify interface implementation
var _ VARIANT = (*TransferInput)(nil)

// TransferInputsSummary is a Record type
type TransferInputsSummary struct {
	TotalAmuletAmount                  NUMERIC  `json:"totalAmuletAmount"`
	TotalAppRewardAmount               NUMERIC  `json:"totalAppRewardAmount"`
	TotalValidatorRewardAmount         NUMERIC  `json:"totalValidatorRewardAmount"`
	TotalValidatorFaucetAmount         NUMERIC  `json:"totalValidatorFaucetAmount"`
	TotalSvRewardAmount                NUMERIC  `json:"totalSvRewardAmount"`
	TotalHoldingFees                   NUMERIC  `json:"totalHoldingFees"`
	AmountArchivedAsOfRoundZero        NUMERIC  `json:"amountArchivedAsOfRoundZero"`
	ChangeToHoldingFeesRate            NUMERIC  `json:"changeToHoldingFeesRate"`
	TotalUnclaimedActivityRecordAmount *NUMERIC `json:"totalUnclaimedActivityRecordAmount"`
	TotalDevelopmentFundAmount         *NUMERIC `json:"totalDevelopmentFundAmount"`
}

// ToMap converts TransferInputsSummary to a map for DAML arguments
func (t TransferInputsSummary) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["totalAmuletAmount"] = (*big.Int)(t.TotalAmuletAmount)

	m["totalAppRewardAmount"] = (*big.Int)(t.TotalAppRewardAmount)

	m["totalValidatorRewardAmount"] = (*big.Int)(t.TotalValidatorRewardAmount)

	m["totalValidatorFaucetAmount"] = (*big.Int)(t.TotalValidatorFaucetAmount)

	m["totalSvRewardAmount"] = (*big.Int)(t.TotalSvRewardAmount)

	m["totalHoldingFees"] = (*big.Int)(t.TotalHoldingFees)

	m["amountArchivedAsOfRoundZero"] = (*big.Int)(t.AmountArchivedAsOfRoundZero)

	m["changeToHoldingFeesRate"] = (*big.Int)(t.ChangeToHoldingFeesRate)

	if t.TotalUnclaimedActivityRecordAmount != nil {
		m["totalUnclaimedActivityRecordAmount"] = map[string]interface{}{
			"_type": "optional",
			"value": (*big.Int)(*t.TotalUnclaimedActivityRecordAmount),
		}
	} else {
		m["totalUnclaimedActivityRecordAmount"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	if t.TotalDevelopmentFundAmount != nil {
		m["totalDevelopmentFundAmount"] = map[string]interface{}{
			"_type": "optional",
			"value": (*big.Int)(*t.TotalDevelopmentFundAmount),
		}
	} else {
		m["totalDevelopmentFundAmount"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for TransferInputsSummary using JsonCodec
func (t TransferInputsSummary) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferInputsSummary using JsonCodec
func (t *TransferInputsSummary) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferOutput is a Record type
type TransferOutput struct {
	Receiver         PARTY     `json:"receiver"`
	ReceiverFeeRatio NUMERIC   `json:"receiverFeeRatio"`
	Amount           NUMERIC   `json:"amount"`
	Lock             *TimeLock `json:"lock"`
}

// ToMap converts TransferOutput to a map for DAML arguments
func (t TransferOutput) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["receiver"] = t.Receiver.ToMap()

	m["receiverFeeRatio"] = (*big.Int)(t.ReceiverFeeRatio)

	m["amount"] = (*big.Int)(t.Amount)

	if t.Lock != nil {
		m["lock"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.Lock,
		}
	} else {
		m["lock"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for TransferOutput using JsonCodec
func (t TransferOutput) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferOutput using JsonCodec
func (t *TransferOutput) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferPreapproval is a Template type
type TransferPreapproval struct {
	Dso           PARTY     `json:"dso"`
	Receiver      PARTY     `json:"receiver"`
	Provider      PARTY     `json:"provider"`
	ValidFrom     TIMESTAMP `json:"validFrom"`
	LastRenewedAt TIMESTAMP `json:"lastRenewedAt"`
	ExpiresAt     TIMESTAMP `json:"expiresAt"`
}

// GetTemplateID returns the template ID for this template
func (t TransferPreapproval) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletRules", "TransferPreapproval")
}

// CreateCommand returns a CreateCommand for this template
func (t TransferPreapproval) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["dso"] = t.Dso.ToMap()

	args["receiver"] = t.Receiver.ToMap()

	args["provider"] = t.Provider.ToMap()

	args["validFrom"] = t.ValidFrom

	args["lastRenewedAt"] = t.LastRenewedAt

	args["expiresAt"] = t.ExpiresAt

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for TransferPreapproval using JsonCodec
func (t TransferPreapproval) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferPreapproval using JsonCodec
func (t *TransferPreapproval) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for TransferPreapproval

// TransferPreapprovalRenew exercises the TransferPreapproval_Renew choice on this TransferPreapproval contract
func (t TransferPreapproval) TransferPreapprovalRenew(contractID string, args TransferPreapprovalRenew) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletRules", "TransferPreapproval"),
		ContractID: contractID,
		Choice:     "TransferPreapproval_Renew",
		Arguments:  argsToMap(args),
	}
}

// TransferPreapprovalSend exercises the TransferPreapproval_Send choice on this TransferPreapproval contract
func (t TransferPreapproval) TransferPreapprovalSend(contractID string, args TransferPreapprovalSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletRules", "TransferPreapproval"),
		ContractID: contractID,
		Choice:     "TransferPreapproval_Send",
		Arguments:  argsToMap(args),
	}
}

// TransferPreapprovalExpire exercises the TransferPreapproval_Expire choice on this TransferPreapproval contract
func (t TransferPreapproval) TransferPreapprovalExpire(contractID string, args TransferPreapprovalExpire) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletRules", "TransferPreapproval"),
		ContractID: contractID,
		Choice:     "TransferPreapproval_Expire",
		Arguments:  argsToMap(args),
	}
}

// TransferPreapprovalCancel exercises the TransferPreapproval_Cancel choice on this TransferPreapproval contract
func (t TransferPreapproval) TransferPreapprovalCancel(contractID string, args TransferPreapprovalCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletRules", "TransferPreapproval"),
		ContractID: contractID,
		Choice:     "TransferPreapproval_Cancel",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this TransferPreapproval contract
func (t TransferPreapproval) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletRules", "TransferPreapproval"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// TransferPreapprovalFetch exercises the TransferPreapproval_Fetch choice on this TransferPreapproval contract
func (t TransferPreapproval) TransferPreapprovalFetch(contractID string, args TransferPreapprovalFetch) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletRules", "TransferPreapproval"),
		ContractID: contractID,
		Choice:     "TransferPreapproval_Fetch",
		Arguments:  argsToMap(args),
	}
}

// TransferPreapprovalCancel is a Record type
type TransferPreapprovalCancel struct {
	P PARTY `json:"p"`
}

// ToMap converts TransferPreapprovalCancel to a map for DAML arguments
func (t TransferPreapprovalCancel) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["p"] = t.P.ToMap()

	return m
}

// MarshalJSON implements custom JSON marshaling for TransferPreapprovalCancel using JsonCodec
func (t TransferPreapprovalCancel) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferPreapprovalCancel using JsonCodec
func (t *TransferPreapprovalCancel) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferPreapprovalCancelResult is an enum type
type TransferPreapprovalCancelResult string

const (
	TransferPreapprovalCancelResultTransferPreapproval_CancelResult TransferPreapprovalCancelResult = "TransferPreapproval_CancelResult"
)

// GetEnumConstructor implements types.ENUM interface
func (e TransferPreapprovalCancelResult) GetEnumConstructor() string {
	return string(e)
}

// GetEnumTypeID implements types.ENUM interface
func (e TransferPreapprovalCancelResult) GetEnumTypeID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletRules", "TransferPreapprovalCancelResult")
}

// MarshalJSON implements custom JSON marshaling for TransferPreapprovalCancelResult using JsonCodec
func (e TransferPreapprovalCancelResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(e)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferPreapprovalCancelResult using JsonCodec
func (e *TransferPreapprovalCancelResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, e)
}

// Verify interface implementation
var _ ENUM = TransferPreapprovalCancelResult("")

// TransferPreapprovalExpire is a Record type
type TransferPreapprovalExpire struct{}

// ToMap converts TransferPreapprovalExpire to a map for DAML arguments
func (t TransferPreapprovalExpire) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	return m
}

// MarshalJSON implements custom JSON marshaling for TransferPreapprovalExpire using JsonCodec
func (t TransferPreapprovalExpire) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferPreapprovalExpire using JsonCodec
func (t *TransferPreapprovalExpire) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferPreapprovalExpireResult is a Record type
type TransferPreapprovalExpireResult struct{}

// ToMap converts TransferPreapprovalExpireResult to a map for DAML arguments
func (t TransferPreapprovalExpireResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	return m
}

// MarshalJSON implements custom JSON marshaling for TransferPreapprovalExpireResult using JsonCodec
func (t TransferPreapprovalExpireResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferPreapprovalExpireResult using JsonCodec
func (t *TransferPreapprovalExpireResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferPreapprovalFetch is a Record type
type TransferPreapprovalFetch struct {
	P PARTY `json:"p"`
}

// ToMap converts TransferPreapprovalFetch to a map for DAML arguments
func (t TransferPreapprovalFetch) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["p"] = t.P.ToMap()

	return m
}

// MarshalJSON implements custom JSON marshaling for TransferPreapprovalFetch using JsonCodec
func (t TransferPreapprovalFetch) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferPreapprovalFetch using JsonCodec
func (t *TransferPreapprovalFetch) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferPreapprovalRenew is a Record type
type TransferPreapprovalRenew struct {
	Context      PaymentTransferContext `json:"context"`
	Inputs       []TransferInput        `json:"inputs"`
	NewExpiresAt TIMESTAMP              `json:"newExpiresAt"`
}

// ToMap converts TransferPreapprovalRenew to a map for DAML arguments
func (t TransferPreapprovalRenew) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["context"] = t.Context

	m["inputs"] = t.Inputs

	m["newExpiresAt"] = t.NewExpiresAt

	return m
}

// MarshalJSON implements custom JSON marshaling for TransferPreapprovalRenew using JsonCodec
func (t TransferPreapprovalRenew) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferPreapprovalRenew using JsonCodec
func (t *TransferPreapprovalRenew) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferPreapprovalRenewResult is a Record type
type TransferPreapprovalRenewResult struct {
	TransferPreapprovalCid CONTRACT_ID    `json:"transferPreapprovalCid"`
	TransferResult         TransferResult `json:"transferResult"`
	Receiver               PARTY          `json:"receiver"`
	Provider               PARTY          `json:"provider"`
	AmuletPaid             NUMERIC        `json:"amuletPaid"`
	Meta                   *Metadata      `json:"meta"`
}

// ToMap converts TransferPreapprovalRenewResult to a map for DAML arguments
func (t TransferPreapprovalRenewResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["transferPreapprovalCid"] = t.TransferPreapprovalCid

	m["transferResult"] = t.TransferResult

	m["receiver"] = t.Receiver.ToMap()

	m["provider"] = t.Provider.ToMap()

	m["amuletPaid"] = (*big.Int)(t.AmuletPaid)

	if t.Meta != nil {
		m["meta"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.Meta,
		}
	} else {
		m["meta"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for TransferPreapprovalRenewResult using JsonCodec
func (t TransferPreapprovalRenewResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferPreapprovalRenewResult using JsonCodec
func (t *TransferPreapprovalRenewResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferPreapprovalSend is a Record type
type TransferPreapprovalSend struct {
	Context     PaymentTransferContext `json:"context"`
	Inputs      []TransferInput        `json:"inputs"`
	Amount      NUMERIC                `json:"amount"`
	Sender      PARTY                  `json:"sender"`
	Description *TEXT                  `json:"description"`
}

// ToMap converts TransferPreapprovalSend to a map for DAML arguments
func (t TransferPreapprovalSend) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["context"] = t.Context

	m["inputs"] = t.Inputs

	m["amount"] = (*big.Int)(t.Amount)

	m["sender"] = t.Sender.ToMap()

	if t.Description != nil {
		m["description"] = map[string]interface{}{
			"_type": "optional",
			"value": string(*t.Description),
		}
	} else {
		m["description"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for TransferPreapprovalSend using JsonCodec
func (t TransferPreapprovalSend) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferPreapprovalSend using JsonCodec
func (t *TransferPreapprovalSend) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferPreapprovalSendResult is a Record type
type TransferPreapprovalSendResult struct {
	Result TransferResult `json:"result"`
	Meta   *Metadata      `json:"meta"`
}

// ToMap converts TransferPreapprovalSendResult to a map for DAML arguments
func (t TransferPreapprovalSendResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["result"] = t.Result

	if t.Meta != nil {
		m["meta"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.Meta,
		}
	} else {
		m["meta"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for TransferPreapprovalSendResult using JsonCodec
func (t TransferPreapprovalSendResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferPreapprovalSendResult using JsonCodec
func (t *TransferPreapprovalSendResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferResult is a Record type
type TransferResult struct {
	Round              Round           `json:"round"`
	Summary            TransferSummary `json:"summary"`
	CreatedAmulets     []CreatedAmulet `json:"createdAmulets"`
	SenderChangeAmulet *CONTRACT_ID    `json:"senderChangeAmulet"`
	Meta               *Metadata       `json:"meta"`
}

// ToMap converts TransferResult to a map for DAML arguments
func (t TransferResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["round"] = t.Round

	m["summary"] = t.Summary

	m["createdAmulets"] = t.CreatedAmulets

	if t.SenderChangeAmulet != nil {
		m["senderChangeAmulet"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.SenderChangeAmulet,
		}
	} else {
		m["senderChangeAmulet"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	if t.Meta != nil {
		m["meta"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.Meta,
		}
	} else {
		m["meta"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for TransferResult using JsonCodec
func (t TransferResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferResult using JsonCodec
func (t *TransferResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferSummary is a Record type
type TransferSummary struct {
	InputAppRewardAmount               NUMERIC   `json:"inputAppRewardAmount"`
	InputValidatorRewardAmount         NUMERIC   `json:"inputValidatorRewardAmount"`
	InputSvRewardAmount                NUMERIC   `json:"inputSvRewardAmount"`
	InputAmuletAmount                  NUMERIC   `json:"inputAmuletAmount"`
	BalanceChanges                     GENMAP    `json:"balanceChanges"`
	HoldingFees                        NUMERIC   `json:"holdingFees"`
	OutputFees                         []NUMERIC `json:"outputFees"`
	SenderChangeFee                    NUMERIC   `json:"senderChangeFee"`
	SenderChangeAmount                 NUMERIC   `json:"senderChangeAmount"`
	AmuletPrice                        NUMERIC   `json:"amuletPrice"`
	InputValidatorFaucetAmount         *NUMERIC  `json:"inputValidatorFaucetAmount"`
	InputUnclaimedActivityRecordAmount *NUMERIC  `json:"inputUnclaimedActivityRecordAmount"`
	InputDevelopmentFundAmount         *NUMERIC  `json:"inputDevelopmentFundAmount"`
}

// ToMap converts TransferSummary to a map for DAML arguments
func (t TransferSummary) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["inputAppRewardAmount"] = (*big.Int)(t.InputAppRewardAmount)

	m["inputValidatorRewardAmount"] = (*big.Int)(t.InputValidatorRewardAmount)

	m["inputSvRewardAmount"] = (*big.Int)(t.InputSvRewardAmount)

	m["inputAmuletAmount"] = (*big.Int)(t.InputAmuletAmount)

	m["balanceChanges"] = map[string]interface{}{"_type": "genmap", "value": t.BalanceChanges}

	m["holdingFees"] = (*big.Int)(t.HoldingFees)

	m["outputFees"] = t.OutputFees

	m["senderChangeFee"] = (*big.Int)(t.SenderChangeFee)

	m["senderChangeAmount"] = (*big.Int)(t.SenderChangeAmount)

	m["amuletPrice"] = (*big.Int)(t.AmuletPrice)

	if t.InputValidatorFaucetAmount != nil {
		m["inputValidatorFaucetAmount"] = map[string]interface{}{
			"_type": "optional",
			"value": (*big.Int)(*t.InputValidatorFaucetAmount),
		}
	} else {
		m["inputValidatorFaucetAmount"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	if t.InputUnclaimedActivityRecordAmount != nil {
		m["inputUnclaimedActivityRecordAmount"] = map[string]interface{}{
			"_type": "optional",
			"value": (*big.Int)(*t.InputUnclaimedActivityRecordAmount),
		}
	} else {
		m["inputUnclaimedActivityRecordAmount"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	if t.InputDevelopmentFundAmount != nil {
		m["inputDevelopmentFundAmount"] = map[string]interface{}{
			"_type": "optional",
			"value": (*big.Int)(*t.InputDevelopmentFundAmount),
		}
	} else {
		m["inputDevelopmentFundAmount"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for TransferSummary using JsonCodec
func (t TransferSummary) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferSummary using JsonCodec
func (t *TransferSummary) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TwoStepTransfer is a Record type
type TwoStepTransfer struct {
	Dso                    PARTY     `json:"dso"`
	Sender                 PARTY     `json:"sender"`
	Receiver               PARTY     `json:"receiver"`
	Amount                 NUMERIC   `json:"amount"`
	LockContext            TEXT      `json:"lockContext"`
	TransferBefore         TIMESTAMP `json:"transferBefore"`
	TransferBeforeDeadline TEXT      `json:"transferBeforeDeadline"`
	Provider               PARTY     `json:"provider"`
	AllowFeaturing         BOOL      `json:"allowFeaturing"`
}

// ToMap converts TwoStepTransfer to a map for DAML arguments
func (t TwoStepTransfer) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["dso"] = t.Dso.ToMap()

	m["sender"] = t.Sender.ToMap()

	m["receiver"] = t.Receiver.ToMap()

	m["amount"] = (*big.Int)(t.Amount)

	m["lockContext"] = string(t.LockContext)

	m["transferBefore"] = t.TransferBefore

	m["transferBeforeDeadline"] = string(t.TransferBeforeDeadline)

	m["provider"] = t.Provider.ToMap()

	m["allowFeaturing"] = bool(t.AllowFeaturing)

	return m
}

// MarshalJSON implements custom JSON marshaling for TwoStepTransfer using JsonCodec
func (t TwoStepTransfer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TwoStepTransfer using JsonCodec
func (t *TwoStepTransfer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TxKind is an enum type
type TxKind string

const (
	TxKindTxKind_Transfer   TxKind = "TxKind_Transfer"
	TxKindTxKind_Unlock     TxKind = "TxKind_Unlock"
	TxKindTxKind_MergeSplit TxKind = "TxKind_MergeSplit"
	TxKindTxKind_Burn       TxKind = "TxKind_Burn"
	TxKindTxKind_Mint       TxKind = "TxKind_Mint"
	TxKindTxKind_ExpireDust TxKind = "TxKind_ExpireDust"
)

// GetEnumConstructor implements types.ENUM interface
func (e TxKind) GetEnumConstructor() string {
	return string(e)
}

// GetEnumTypeID implements types.ENUM interface
func (e TxKind) GetEnumTypeID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet.TokenApiUtils", "TxKind")
}

// MarshalJSON implements custom JSON marshaling for TxKind using JsonCodec
func (e TxKind) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(e)
}

// UnmarshalJSON implements custom JSON unmarshaling for TxKind using JsonCodec
func (e *TxKind) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, e)
}

// Verify interface implementation
var _ ENUM = TxKind("")

// USD is an enum type
type USD string

const (
	USDUSD USD = "USD"
)

// GetEnumConstructor implements types.ENUM interface
func (e USD) GetEnumConstructor() string {
	return string(e)
}

// GetEnumTypeID implements types.ENUM interface
func (e USD) GetEnumTypeID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.AmuletConfig", "USD")
}

// MarshalJSON implements custom JSON marshaling for USD using JsonCodec
func (e USD) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(e)
}

// UnmarshalJSON implements custom JSON unmarshaling for USD using JsonCodec
func (e *USD) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, e)
}

// Verify interface implementation
var _ ENUM = USD("")

// UnclaimedActivityRecord is a Template type
type UnclaimedActivityRecord struct {
	Dso         PARTY     `json:"dso"`
	Beneficiary PARTY     `json:"beneficiary"`
	Amount      NUMERIC   `json:"amount"`
	Reason      TEXT      `json:"reason"`
	ExpiresAt   TIMESTAMP `json:"expiresAt"`
}

// GetTemplateID returns the template ID for this template
func (t UnclaimedActivityRecord) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "UnclaimedActivityRecord")
}

// CreateCommand returns a CreateCommand for this template
func (t UnclaimedActivityRecord) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["dso"] = t.Dso.ToMap()

	args["beneficiary"] = t.Beneficiary.ToMap()

	if t.Amount != nil {
		args["amount"] = (*big.Int)(t.Amount)
	}

	args["reason"] = string(t.Reason)

	args["expiresAt"] = t.ExpiresAt

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for UnclaimedActivityRecord using JsonCodec
func (t UnclaimedActivityRecord) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for UnclaimedActivityRecord using JsonCodec
func (t *UnclaimedActivityRecord) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for UnclaimedActivityRecord

// UnclaimedActivityRecordDsoExpire exercises the UnclaimedActivityRecord_DsoExpire choice on this UnclaimedActivityRecord contract
func (t UnclaimedActivityRecord) UnclaimedActivityRecordDsoExpire(contractID string, args UnclaimedActivityRecordDsoExpire) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "UnclaimedActivityRecord"),
		ContractID: contractID,
		Choice:     "UnclaimedActivityRecord_DsoExpire",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this UnclaimedActivityRecord contract
func (t UnclaimedActivityRecord) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "UnclaimedActivityRecord"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// UnclaimedActivityRecordArchiveAsBeneficiaryResult is an enum type
type UnclaimedActivityRecordArchiveAsBeneficiaryResult string

const (
	UnclaimedActivityRecordArchiveAsBeneficiaryResultUnclaimedActivityRecord_ArchiveAsBeneficiaryResult UnclaimedActivityRecordArchiveAsBeneficiaryResult = "UnclaimedActivityRecord_ArchiveAsBeneficiaryResult"
)

// GetEnumConstructor implements types.ENUM interface
func (e UnclaimedActivityRecordArchiveAsBeneficiaryResult) GetEnumConstructor() string {
	return string(e)
}

// GetEnumTypeID implements types.ENUM interface
func (e UnclaimedActivityRecordArchiveAsBeneficiaryResult) GetEnumTypeID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "UnclaimedActivityRecordArchiveAsBeneficiaryResult")
}

// MarshalJSON implements custom JSON marshaling for UnclaimedActivityRecordArchiveAsBeneficiaryResult using JsonCodec
func (e UnclaimedActivityRecordArchiveAsBeneficiaryResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(e)
}

// UnmarshalJSON implements custom JSON unmarshaling for UnclaimedActivityRecordArchiveAsBeneficiaryResult using JsonCodec
func (e *UnclaimedActivityRecordArchiveAsBeneficiaryResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, e)
}

// Verify interface implementation
var _ ENUM = UnclaimedActivityRecordArchiveAsBeneficiaryResult("")

// UnclaimedActivityRecordDsoExpire is a Record type
type UnclaimedActivityRecordDsoExpire struct{}

// ToMap converts UnclaimedActivityRecordDsoExpire to a map for DAML arguments
func (t UnclaimedActivityRecordDsoExpire) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	return m
}

// MarshalJSON implements custom JSON marshaling for UnclaimedActivityRecordDsoExpire using JsonCodec
func (t UnclaimedActivityRecordDsoExpire) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for UnclaimedActivityRecordDsoExpire using JsonCodec
func (t *UnclaimedActivityRecordDsoExpire) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// UnclaimedActivityRecordDsoExpireResult is a Record type
type UnclaimedActivityRecordDsoExpireResult struct {
	UnclaimedRewardCid CONTRACT_ID `json:"unclaimedRewardCid"`
}

// ToMap converts UnclaimedActivityRecordDsoExpireResult to a map for DAML arguments
func (t UnclaimedActivityRecordDsoExpireResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["unclaimedRewardCid"] = t.UnclaimedRewardCid

	return m
}

// MarshalJSON implements custom JSON marshaling for UnclaimedActivityRecordDsoExpireResult using JsonCodec
func (t UnclaimedActivityRecordDsoExpireResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for UnclaimedActivityRecordDsoExpireResult using JsonCodec
func (t *UnclaimedActivityRecordDsoExpireResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// UnclaimedDevelopmentFundCoupon is a Template type
type UnclaimedDevelopmentFundCoupon struct {
	Dso    PARTY   `json:"dso"`
	Amount NUMERIC `json:"amount"`
}

// GetTemplateID returns the template ID for this template
func (t UnclaimedDevelopmentFundCoupon) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "UnclaimedDevelopmentFundCoupon")
}

// CreateCommand returns a CreateCommand for this template
func (t UnclaimedDevelopmentFundCoupon) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["dso"] = t.Dso.ToMap()

	if t.Amount != nil {
		args["amount"] = (*big.Int)(t.Amount)
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for UnclaimedDevelopmentFundCoupon using JsonCodec
func (t UnclaimedDevelopmentFundCoupon) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for UnclaimedDevelopmentFundCoupon using JsonCodec
func (t *UnclaimedDevelopmentFundCoupon) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for UnclaimedDevelopmentFundCoupon

// Archive exercises the Archive choice on this UnclaimedDevelopmentFundCoupon contract
func (t UnclaimedDevelopmentFundCoupon) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "UnclaimedDevelopmentFundCoupon"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// UnclaimedReward is a Template type
type UnclaimedReward struct {
	Dso    PARTY   `json:"dso"`
	Amount NUMERIC `json:"amount"`
}

// GetTemplateID returns the template ID for this template
func (t UnclaimedReward) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "UnclaimedReward")
}

// CreateCommand returns a CreateCommand for this template
func (t UnclaimedReward) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["dso"] = t.Dso.ToMap()

	if t.Amount != nil {
		args["amount"] = (*big.Int)(t.Amount)
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for UnclaimedReward using JsonCodec
func (t UnclaimedReward) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for UnclaimedReward using JsonCodec
func (t *UnclaimedReward) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for UnclaimedReward

// Archive exercises the Archive choice on this UnclaimedReward contract
func (t UnclaimedReward) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "UnclaimedReward"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// ValidatorFaucetCoupon is a Template type
type ValidatorFaucetCoupon struct {
	Dso       PARTY `json:"dso"`
	Validator PARTY `json:"validator"`
	Round     Round `json:"round"`
}

// GetTemplateID returns the template ID for this template
func (t ValidatorFaucetCoupon) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.ValidatorLicense", "ValidatorFaucetCoupon")
}

// CreateCommand returns a CreateCommand for this template
func (t ValidatorFaucetCoupon) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["dso"] = t.Dso.ToMap()

	args["validator"] = t.Validator.ToMap()

	args["round"] = t.Round

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for ValidatorFaucetCoupon using JsonCodec
func (t ValidatorFaucetCoupon) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ValidatorFaucetCoupon using JsonCodec
func (t *ValidatorFaucetCoupon) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for ValidatorFaucetCoupon

// ValidatorFaucetCouponDsoExpire exercises the ValidatorFaucetCoupon_DsoExpire choice on this ValidatorFaucetCoupon contract
func (t ValidatorFaucetCoupon) ValidatorFaucetCouponDsoExpire(contractID string, args ValidatorFaucetCouponDsoExpire) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.ValidatorLicense", "ValidatorFaucetCoupon"),
		ContractID: contractID,
		Choice:     "ValidatorFaucetCoupon_DsoExpire",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this ValidatorFaucetCoupon contract
func (t ValidatorFaucetCoupon) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.ValidatorLicense", "ValidatorFaucetCoupon"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// ValidatorFaucetCouponDsoExpire is a Record type
type ValidatorFaucetCouponDsoExpire struct {
	ClosedRoundCid CONTRACT_ID `json:"closedRoundCid"`
}

// ToMap converts ValidatorFaucetCouponDsoExpire to a map for DAML arguments
func (t ValidatorFaucetCouponDsoExpire) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["closedRoundCid"] = t.ClosedRoundCid

	return m
}

// MarshalJSON implements custom JSON marshaling for ValidatorFaucetCouponDsoExpire using JsonCodec
func (t ValidatorFaucetCouponDsoExpire) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ValidatorFaucetCouponDsoExpire using JsonCodec
func (t *ValidatorFaucetCouponDsoExpire) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ValidatorFaucetCouponDsoExpireResult is an enum type
type ValidatorFaucetCouponDsoExpireResult string

const (
	ValidatorFaucetCouponDsoExpireResultValidatorFaucetCoupon_DsoExpireResult ValidatorFaucetCouponDsoExpireResult = "ValidatorFaucetCoupon_DsoExpireResult"
)

// GetEnumConstructor implements types.ENUM interface
func (e ValidatorFaucetCouponDsoExpireResult) GetEnumConstructor() string {
	return string(e)
}

// GetEnumTypeID implements types.ENUM interface
func (e ValidatorFaucetCouponDsoExpireResult) GetEnumTypeID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.ValidatorLicense", "ValidatorFaucetCouponDsoExpireResult")
}

// MarshalJSON implements custom JSON marshaling for ValidatorFaucetCouponDsoExpireResult using JsonCodec
func (e ValidatorFaucetCouponDsoExpireResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(e)
}

// UnmarshalJSON implements custom JSON unmarshaling for ValidatorFaucetCouponDsoExpireResult using JsonCodec
func (e *ValidatorFaucetCouponDsoExpireResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, e)
}

// Verify interface implementation
var _ ENUM = ValidatorFaucetCouponDsoExpireResult("")

// ValidatorLicense is a Template type
type ValidatorLicense struct {
	Validator    PARTY                     `json:"validator"`
	Sponsor      PARTY                     `json:"sponsor"`
	Dso          PARTY                     `json:"dso"`
	FaucetState  *FaucetState              `json:"faucetState"`
	Metadata     *ValidatorLicenseMetadata `json:"metadata"`
	LastActiveAt *TIMESTAMP                `json:"lastActiveAt"`
}

// GetTemplateID returns the template ID for this template
func (t ValidatorLicense) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.ValidatorLicense", "ValidatorLicense")
}

// CreateCommand returns a CreateCommand for this template
func (t ValidatorLicense) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["validator"] = t.Validator.ToMap()

	args["sponsor"] = t.Sponsor.ToMap()

	args["dso"] = t.Dso.ToMap()

	if t.FaucetState != nil {
		args["faucetState"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.FaucetState,
		}
	} else {
		args["faucetState"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	if t.Metadata != nil {
		args["metadata"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.Metadata,
		}
	} else {
		args["metadata"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	if t.LastActiveAt != nil {
		args["lastActiveAt"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.LastActiveAt,
		}
	} else {
		args["lastActiveAt"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for ValidatorLicense using JsonCodec
func (t ValidatorLicense) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ValidatorLicense using JsonCodec
func (t *ValidatorLicense) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for ValidatorLicense

// ValidatorLicenseReceiveFaucetCoupon exercises the ValidatorLicense_ReceiveFaucetCoupon choice on this ValidatorLicense contract
func (t ValidatorLicense) ValidatorLicenseReceiveFaucetCoupon(contractID string, args ValidatorLicenseReceiveFaucetCoupon) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.ValidatorLicense", "ValidatorLicense"),
		ContractID: contractID,
		Choice:     "ValidatorLicense_ReceiveFaucetCoupon",
		Arguments:  argsToMap(args),
	}
}

// ValidatorLicenseRecordValidatorLivenessActivity exercises the ValidatorLicense_RecordValidatorLivenessActivity choice on this ValidatorLicense contract
func (t ValidatorLicense) ValidatorLicenseRecordValidatorLivenessActivity(contractID string, args ValidatorLicenseRecordValidatorLivenessActivity) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.ValidatorLicense", "ValidatorLicense"),
		ContractID: contractID,
		Choice:     "ValidatorLicense_RecordValidatorLivenessActivity",
		Arguments:  argsToMap(args),
	}
}

// ValidatorLicenseWithdraw exercises the ValidatorLicense_Withdraw choice on this ValidatorLicense contract
func (t ValidatorLicense) ValidatorLicenseWithdraw(contractID string, args ValidatorLicenseWithdraw) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.ValidatorLicense", "ValidatorLicense"),
		ContractID: contractID,
		Choice:     "ValidatorLicense_Withdraw",
		Arguments:  argsToMap(args),
	}
}

// ValidatorLicenseCancel exercises the ValidatorLicense_Cancel choice on this ValidatorLicense contract
func (t ValidatorLicense) ValidatorLicenseCancel(contractID string, args ValidatorLicenseCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.ValidatorLicense", "ValidatorLicense"),
		ContractID: contractID,
		Choice:     "ValidatorLicense_Cancel",
		Arguments:  argsToMap(args),
	}
}

// ValidatorLicenseUpdateMetadata exercises the ValidatorLicense_UpdateMetadata choice on this ValidatorLicense contract
func (t ValidatorLicense) ValidatorLicenseUpdateMetadata(contractID string, args ValidatorLicenseUpdateMetadata) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.ValidatorLicense", "ValidatorLicense"),
		ContractID: contractID,
		Choice:     "ValidatorLicense_UpdateMetadata",
		Arguments:  argsToMap(args),
	}
}

// ValidatorLicenseReportActive exercises the ValidatorLicense_ReportActive choice on this ValidatorLicense contract
func (t ValidatorLicense) ValidatorLicenseReportActive(contractID string, args ValidatorLicenseReportActive) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.ValidatorLicense", "ValidatorLicense"),
		ContractID: contractID,
		Choice:     "ValidatorLicense_ReportActive",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this ValidatorLicense contract
func (t ValidatorLicense) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.ValidatorLicense", "ValidatorLicense"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// ValidatorLicenseMetadata is a Record type
type ValidatorLicenseMetadata struct {
	LastUpdatedAt TIMESTAMP `json:"lastUpdatedAt"`
	Version       TEXT      `json:"version"`
	ContactPoint  TEXT      `json:"contactPoint"`
}

// ToMap converts ValidatorLicenseMetadata to a map for DAML arguments
func (t ValidatorLicenseMetadata) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["lastUpdatedAt"] = t.LastUpdatedAt

	m["version"] = string(t.Version)

	m["contactPoint"] = string(t.ContactPoint)

	return m
}

// MarshalJSON implements custom JSON marshaling for ValidatorLicenseMetadata using JsonCodec
func (t ValidatorLicenseMetadata) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ValidatorLicenseMetadata using JsonCodec
func (t *ValidatorLicenseMetadata) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ValidatorLicenseCancel is a Record type
type ValidatorLicenseCancel struct {
	Reason TEXT `json:"reason"`
}

// ToMap converts ValidatorLicenseCancel to a map for DAML arguments
func (t ValidatorLicenseCancel) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["reason"] = string(t.Reason)

	return m
}

// MarshalJSON implements custom JSON marshaling for ValidatorLicenseCancel using JsonCodec
func (t ValidatorLicenseCancel) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ValidatorLicenseCancel using JsonCodec
func (t *ValidatorLicenseCancel) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ValidatorLicenseCancelResult is an enum type
type ValidatorLicenseCancelResult string

const (
	ValidatorLicenseCancelResultValidatorLicense_CancelResult ValidatorLicenseCancelResult = "ValidatorLicense_CancelResult"
)

// GetEnumConstructor implements types.ENUM interface
func (e ValidatorLicenseCancelResult) GetEnumConstructor() string {
	return string(e)
}

// GetEnumTypeID implements types.ENUM interface
func (e ValidatorLicenseCancelResult) GetEnumTypeID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.ValidatorLicense", "ValidatorLicenseCancelResult")
}

// MarshalJSON implements custom JSON marshaling for ValidatorLicenseCancelResult using JsonCodec
func (e ValidatorLicenseCancelResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(e)
}

// UnmarshalJSON implements custom JSON unmarshaling for ValidatorLicenseCancelResult using JsonCodec
func (e *ValidatorLicenseCancelResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, e)
}

// Verify interface implementation
var _ ENUM = ValidatorLicenseCancelResult("")

// ValidatorLicenseReceiveFaucetCoupon is a Record type
type ValidatorLicenseReceiveFaucetCoupon struct {
	OpenRoundCid CONTRACT_ID `json:"openRoundCid"`
}

// ToMap converts ValidatorLicenseReceiveFaucetCoupon to a map for DAML arguments
func (t ValidatorLicenseReceiveFaucetCoupon) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["openRoundCid"] = t.OpenRoundCid

	return m
}

// MarshalJSON implements custom JSON marshaling for ValidatorLicenseReceiveFaucetCoupon using JsonCodec
func (t ValidatorLicenseReceiveFaucetCoupon) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ValidatorLicenseReceiveFaucetCoupon using JsonCodec
func (t *ValidatorLicenseReceiveFaucetCoupon) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ValidatorLicenseReceiveFaucetCouponResult is a Record type
type ValidatorLicenseReceiveFaucetCouponResult struct {
	LicenseCid CONTRACT_ID `json:"licenseCid"`
	CouponCid  CONTRACT_ID `json:"couponCid"`
}

// ToMap converts ValidatorLicenseReceiveFaucetCouponResult to a map for DAML arguments
func (t ValidatorLicenseReceiveFaucetCouponResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["licenseCid"] = t.LicenseCid

	m["couponCid"] = t.CouponCid

	return m
}

// MarshalJSON implements custom JSON marshaling for ValidatorLicenseReceiveFaucetCouponResult using JsonCodec
func (t ValidatorLicenseReceiveFaucetCouponResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ValidatorLicenseReceiveFaucetCouponResult using JsonCodec
func (t *ValidatorLicenseReceiveFaucetCouponResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ValidatorLicenseRecordValidatorLivenessActivity is a Record type
type ValidatorLicenseRecordValidatorLivenessActivity struct {
	OpenRoundCid CONTRACT_ID `json:"openRoundCid"`
}

// ToMap converts ValidatorLicenseRecordValidatorLivenessActivity to a map for DAML arguments
func (t ValidatorLicenseRecordValidatorLivenessActivity) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["openRoundCid"] = t.OpenRoundCid

	return m
}

// MarshalJSON implements custom JSON marshaling for ValidatorLicenseRecordValidatorLivenessActivity using JsonCodec
func (t ValidatorLicenseRecordValidatorLivenessActivity) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ValidatorLicenseRecordValidatorLivenessActivity using JsonCodec
func (t *ValidatorLicenseRecordValidatorLivenessActivity) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ValidatorLicenseRecordValidatorLivenessActivityResult is a Record type
type ValidatorLicenseRecordValidatorLivenessActivityResult struct {
	LicenseCid CONTRACT_ID `json:"licenseCid"`
	CouponCid  CONTRACT_ID `json:"couponCid"`
}

// ToMap converts ValidatorLicenseRecordValidatorLivenessActivityResult to a map for DAML arguments
func (t ValidatorLicenseRecordValidatorLivenessActivityResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["licenseCid"] = t.LicenseCid

	m["couponCid"] = t.CouponCid

	return m
}

// MarshalJSON implements custom JSON marshaling for ValidatorLicenseRecordValidatorLivenessActivityResult using JsonCodec
func (t ValidatorLicenseRecordValidatorLivenessActivityResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ValidatorLicenseRecordValidatorLivenessActivityResult using JsonCodec
func (t *ValidatorLicenseRecordValidatorLivenessActivityResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ValidatorLicenseReportActive is a Record type
type ValidatorLicenseReportActive struct{}

// ToMap converts ValidatorLicenseReportActive to a map for DAML arguments
func (t ValidatorLicenseReportActive) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	return m
}

// MarshalJSON implements custom JSON marshaling for ValidatorLicenseReportActive using JsonCodec
func (t ValidatorLicenseReportActive) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ValidatorLicenseReportActive using JsonCodec
func (t *ValidatorLicenseReportActive) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ValidatorLicenseReportActiveResult is a Record type
type ValidatorLicenseReportActiveResult struct {
	LicenseCid CONTRACT_ID `json:"licenseCid"`
}

// ToMap converts ValidatorLicenseReportActiveResult to a map for DAML arguments
func (t ValidatorLicenseReportActiveResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["licenseCid"] = t.LicenseCid

	return m
}

// MarshalJSON implements custom JSON marshaling for ValidatorLicenseReportActiveResult using JsonCodec
func (t ValidatorLicenseReportActiveResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ValidatorLicenseReportActiveResult using JsonCodec
func (t *ValidatorLicenseReportActiveResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ValidatorLicenseUpdateMetadata is a Record type
type ValidatorLicenseUpdateMetadata struct {
	Version      TEXT `json:"version"`
	ContactPoint TEXT `json:"contactPoint"`
}

// ToMap converts ValidatorLicenseUpdateMetadata to a map for DAML arguments
func (t ValidatorLicenseUpdateMetadata) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["version"] = string(t.Version)

	m["contactPoint"] = string(t.ContactPoint)

	return m
}

// MarshalJSON implements custom JSON marshaling for ValidatorLicenseUpdateMetadata using JsonCodec
func (t ValidatorLicenseUpdateMetadata) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ValidatorLicenseUpdateMetadata using JsonCodec
func (t *ValidatorLicenseUpdateMetadata) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ValidatorLicenseUpdateMetadataResult is a Record type
type ValidatorLicenseUpdateMetadataResult struct {
	LicenseCid CONTRACT_ID `json:"licenseCid"`
}

// ToMap converts ValidatorLicenseUpdateMetadataResult to a map for DAML arguments
func (t ValidatorLicenseUpdateMetadataResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["licenseCid"] = t.LicenseCid

	return m
}

// MarshalJSON implements custom JSON marshaling for ValidatorLicenseUpdateMetadataResult using JsonCodec
func (t ValidatorLicenseUpdateMetadataResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ValidatorLicenseUpdateMetadataResult using JsonCodec
func (t *ValidatorLicenseUpdateMetadataResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ValidatorLicenseWithdraw is a Record type
type ValidatorLicenseWithdraw struct {
	Reason TEXT `json:"reason"`
}

// ToMap converts ValidatorLicenseWithdraw to a map for DAML arguments
func (t ValidatorLicenseWithdraw) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["reason"] = string(t.Reason)

	return m
}

// MarshalJSON implements custom JSON marshaling for ValidatorLicenseWithdraw using JsonCodec
func (t ValidatorLicenseWithdraw) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ValidatorLicenseWithdraw using JsonCodec
func (t *ValidatorLicenseWithdraw) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ValidatorLicenseWithdrawResult is an enum type
type ValidatorLicenseWithdrawResult string

const (
	ValidatorLicenseWithdrawResultValidatorLicense_WithdrawResult ValidatorLicenseWithdrawResult = "ValidatorLicense_WithdrawResult"
)

// GetEnumConstructor implements types.ENUM interface
func (e ValidatorLicenseWithdrawResult) GetEnumConstructor() string {
	return string(e)
}

// GetEnumTypeID implements types.ENUM interface
func (e ValidatorLicenseWithdrawResult) GetEnumTypeID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.ValidatorLicense", "ValidatorLicenseWithdrawResult")
}

// MarshalJSON implements custom JSON marshaling for ValidatorLicenseWithdrawResult using JsonCodec
func (e ValidatorLicenseWithdrawResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(e)
}

// UnmarshalJSON implements custom JSON unmarshaling for ValidatorLicenseWithdrawResult using JsonCodec
func (e *ValidatorLicenseWithdrawResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, e)
}

// Verify interface implementation
var _ ENUM = ValidatorLicenseWithdrawResult("")

// ValidatorLivenessActivityRecord is a Template type
type ValidatorLivenessActivityRecord struct {
	Dso       PARTY `json:"dso"`
	Validator PARTY `json:"validator"`
	Round     Round `json:"round"`
}

// GetTemplateID returns the template ID for this template
func (t ValidatorLivenessActivityRecord) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.ValidatorLicense", "ValidatorLivenessActivityRecord")
}

// CreateCommand returns a CreateCommand for this template
func (t ValidatorLivenessActivityRecord) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["dso"] = t.Dso.ToMap()

	args["validator"] = t.Validator.ToMap()

	args["round"] = t.Round

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for ValidatorLivenessActivityRecord using JsonCodec
func (t ValidatorLivenessActivityRecord) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ValidatorLivenessActivityRecord using JsonCodec
func (t *ValidatorLivenessActivityRecord) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for ValidatorLivenessActivityRecord

// Archive exercises the Archive choice on this ValidatorLivenessActivityRecord contract
func (t ValidatorLivenessActivityRecord) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.ValidatorLicense", "ValidatorLivenessActivityRecord"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// ValidatorLivenessActivityRecordDsoExpire exercises the ValidatorLivenessActivityRecord_DsoExpire choice on this ValidatorLivenessActivityRecord contract
func (t ValidatorLivenessActivityRecord) ValidatorLivenessActivityRecordDsoExpire(contractID string, args ValidatorLivenessActivityRecordDsoExpire) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.ValidatorLicense", "ValidatorLivenessActivityRecord"),
		ContractID: contractID,
		Choice:     "ValidatorLivenessActivityRecord_DsoExpire",
		Arguments:  argsToMap(args),
	}
}

// ValidatorLivenessActivityRecordDsoExpire is a Record type
type ValidatorLivenessActivityRecordDsoExpire struct {
	ClosedRoundCid CONTRACT_ID `json:"closedRoundCid"`
}

// ToMap converts ValidatorLivenessActivityRecordDsoExpire to a map for DAML arguments
func (t ValidatorLivenessActivityRecordDsoExpire) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["closedRoundCid"] = t.ClosedRoundCid

	return m
}

// MarshalJSON implements custom JSON marshaling for ValidatorLivenessActivityRecordDsoExpire using JsonCodec
func (t ValidatorLivenessActivityRecordDsoExpire) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ValidatorLivenessActivityRecordDsoExpire using JsonCodec
func (t *ValidatorLivenessActivityRecordDsoExpire) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ValidatorLivenessActivityRecordDsoExpireResult is an enum type
type ValidatorLivenessActivityRecordDsoExpireResult string

const (
	ValidatorLivenessActivityRecordDsoExpireResultValidatorLivenessActivityRecord_DsoExpireResult ValidatorLivenessActivityRecordDsoExpireResult = "ValidatorLivenessActivityRecord_DsoExpireResult"
)

// GetEnumConstructor implements types.ENUM interface
func (e ValidatorLivenessActivityRecordDsoExpireResult) GetEnumConstructor() string {
	return string(e)
}

// GetEnumTypeID implements types.ENUM interface
func (e ValidatorLivenessActivityRecordDsoExpireResult) GetEnumTypeID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.ValidatorLicense", "ValidatorLivenessActivityRecordDsoExpireResult")
}

// MarshalJSON implements custom JSON marshaling for ValidatorLivenessActivityRecordDsoExpireResult using JsonCodec
func (e ValidatorLivenessActivityRecordDsoExpireResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(e)
}

// UnmarshalJSON implements custom JSON unmarshaling for ValidatorLivenessActivityRecordDsoExpireResult using JsonCodec
func (e *ValidatorLivenessActivityRecordDsoExpireResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, e)
}

// Verify interface implementation
var _ ENUM = ValidatorLivenessActivityRecordDsoExpireResult("")

// ValidatorRewardCoupon is a Template type
type ValidatorRewardCoupon struct {
	Dso    PARTY   `json:"dso"`
	User   PARTY   `json:"user"`
	Amount NUMERIC `json:"amount"`
	Round  Round   `json:"round"`
}

// GetTemplateID returns the template ID for this template
func (t ValidatorRewardCoupon) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "ValidatorRewardCoupon")
}

// CreateCommand returns a CreateCommand for this template
func (t ValidatorRewardCoupon) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["dso"] = t.Dso.ToMap()

	args["user"] = t.User.ToMap()

	if t.Amount != nil {
		args["amount"] = (*big.Int)(t.Amount)
	}

	args["round"] = t.Round

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for ValidatorRewardCoupon using JsonCodec
func (t ValidatorRewardCoupon) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ValidatorRewardCoupon using JsonCodec
func (t *ValidatorRewardCoupon) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for ValidatorRewardCoupon

// ValidatorRewardCouponDsoExpire exercises the ValidatorRewardCoupon_DsoExpire choice on this ValidatorRewardCoupon contract
func (t ValidatorRewardCoupon) ValidatorRewardCouponDsoExpire(contractID string, args ValidatorRewardCouponDsoExpire) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "ValidatorRewardCoupon"),
		ContractID: contractID,
		Choice:     "ValidatorRewardCoupon_DsoExpire",
		Arguments:  argsToMap(args),
	}
}

// ValidatorRewardCouponArchiveAsValidator exercises the ValidatorRewardCoupon_ArchiveAsValidator choice on this ValidatorRewardCoupon contract
func (t ValidatorRewardCoupon) ValidatorRewardCouponArchiveAsValidator(contractID string, args ValidatorRewardCouponArchiveAsValidator) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "ValidatorRewardCoupon"),
		ContractID: contractID,
		Choice:     "ValidatorRewardCoupon_ArchiveAsValidator",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this ValidatorRewardCoupon contract
func (t ValidatorRewardCoupon) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "ValidatorRewardCoupon"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// ValidatorRewardCouponArchiveAsValidator is a Record type
type ValidatorRewardCouponArchiveAsValidator struct {
	Validator PARTY       `json:"validator"`
	RightCid  CONTRACT_ID `json:"rightCid"`
}

// ToMap converts ValidatorRewardCouponArchiveAsValidator to a map for DAML arguments
func (t ValidatorRewardCouponArchiveAsValidator) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["validator"] = t.Validator.ToMap()

	m["rightCid"] = t.RightCid

	return m
}

// MarshalJSON implements custom JSON marshaling for ValidatorRewardCouponArchiveAsValidator using JsonCodec
func (t ValidatorRewardCouponArchiveAsValidator) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ValidatorRewardCouponArchiveAsValidator using JsonCodec
func (t *ValidatorRewardCouponArchiveAsValidator) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ValidatorRewardCouponArchiveAsValidatorResult is a Record type
type ValidatorRewardCouponArchiveAsValidatorResult struct{}

// ToMap converts ValidatorRewardCouponArchiveAsValidatorResult to a map for DAML arguments
func (t ValidatorRewardCouponArchiveAsValidatorResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	return m
}

// MarshalJSON implements custom JSON marshaling for ValidatorRewardCouponArchiveAsValidatorResult using JsonCodec
func (t ValidatorRewardCouponArchiveAsValidatorResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ValidatorRewardCouponArchiveAsValidatorResult using JsonCodec
func (t *ValidatorRewardCouponArchiveAsValidatorResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ValidatorRewardCouponDsoExpire is a Record type
type ValidatorRewardCouponDsoExpire struct {
	ClosedRoundCid CONTRACT_ID `json:"closedRoundCid"`
}

// ToMap converts ValidatorRewardCouponDsoExpire to a map for DAML arguments
func (t ValidatorRewardCouponDsoExpire) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["closedRoundCid"] = t.ClosedRoundCid

	return m
}

// MarshalJSON implements custom JSON marshaling for ValidatorRewardCouponDsoExpire using JsonCodec
func (t ValidatorRewardCouponDsoExpire) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ValidatorRewardCouponDsoExpire using JsonCodec
func (t *ValidatorRewardCouponDsoExpire) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ValidatorRewardCouponDsoExpireResult is a Record type
type ValidatorRewardCouponDsoExpireResult struct {
	Amount NUMERIC `json:"amount"`
}

// ToMap converts ValidatorRewardCouponDsoExpireResult to a map for DAML arguments
func (t ValidatorRewardCouponDsoExpireResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["amount"] = (*big.Int)(t.Amount)

	return m
}

// MarshalJSON implements custom JSON marshaling for ValidatorRewardCouponDsoExpireResult using JsonCodec
func (t ValidatorRewardCouponDsoExpireResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ValidatorRewardCouponDsoExpireResult using JsonCodec
func (t *ValidatorRewardCouponDsoExpireResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ValidatorRight is a Template type
type ValidatorRight struct {
	Dso       PARTY `json:"dso"`
	User      PARTY `json:"user"`
	Validator PARTY `json:"validator"`
}

// GetTemplateID returns the template ID for this template
func (t ValidatorRight) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "ValidatorRight")
}

// CreateCommand returns a CreateCommand for this template
func (t ValidatorRight) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["dso"] = t.Dso.ToMap()

	args["user"] = t.User.ToMap()

	args["validator"] = t.Validator.ToMap()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for ValidatorRight using JsonCodec
func (t ValidatorRight) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ValidatorRight using JsonCodec
func (t *ValidatorRight) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for ValidatorRight

// ValidatorRightArchiveAsValidator exercises the ValidatorRight_ArchiveAsValidator choice on this ValidatorRight contract
func (t ValidatorRight) ValidatorRightArchiveAsValidator(contractID string, args ValidatorRightArchiveAsValidator) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "ValidatorRight"),
		ContractID: contractID,
		Choice:     "ValidatorRight_ArchiveAsValidator",
		Arguments:  argsToMap(args),
	}
}

// ValidatorRightArchiveAsUser exercises the ValidatorRight_ArchiveAsUser choice on this ValidatorRight contract
func (t ValidatorRight) ValidatorRightArchiveAsUser(contractID string, args ValidatorRightArchiveAsUser) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "ValidatorRight"),
		ContractID: contractID,
		Choice:     "ValidatorRight_ArchiveAsUser",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this ValidatorRight contract
func (t ValidatorRight) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "ValidatorRight"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// ValidatorRightArchiveAsUser is a Record type
type ValidatorRightArchiveAsUser struct{}

// ToMap converts ValidatorRightArchiveAsUser to a map for DAML arguments
func (t ValidatorRightArchiveAsUser) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	return m
}

// MarshalJSON implements custom JSON marshaling for ValidatorRightArchiveAsUser using JsonCodec
func (t ValidatorRightArchiveAsUser) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ValidatorRightArchiveAsUser using JsonCodec
func (t *ValidatorRightArchiveAsUser) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ValidatorRightArchiveAsUserResult is an enum type
type ValidatorRightArchiveAsUserResult string

const (
	ValidatorRightArchiveAsUserResultValidatorRight_ArchiveAsUserResult ValidatorRightArchiveAsUserResult = "ValidatorRight_ArchiveAsUserResult"
)

// GetEnumConstructor implements types.ENUM interface
func (e ValidatorRightArchiveAsUserResult) GetEnumConstructor() string {
	return string(e)
}

// GetEnumTypeID implements types.ENUM interface
func (e ValidatorRightArchiveAsUserResult) GetEnumTypeID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "ValidatorRightArchiveAsUserResult")
}

// MarshalJSON implements custom JSON marshaling for ValidatorRightArchiveAsUserResult using JsonCodec
func (e ValidatorRightArchiveAsUserResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(e)
}

// UnmarshalJSON implements custom JSON unmarshaling for ValidatorRightArchiveAsUserResult using JsonCodec
func (e *ValidatorRightArchiveAsUserResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, e)
}

// Verify interface implementation
var _ ENUM = ValidatorRightArchiveAsUserResult("")

// ValidatorRightArchiveAsValidator is a Record type
type ValidatorRightArchiveAsValidator struct{}

// ToMap converts ValidatorRightArchiveAsValidator to a map for DAML arguments
func (t ValidatorRightArchiveAsValidator) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	return m
}

// MarshalJSON implements custom JSON marshaling for ValidatorRightArchiveAsValidator using JsonCodec
func (t ValidatorRightArchiveAsValidator) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ValidatorRightArchiveAsValidator using JsonCodec
func (t *ValidatorRightArchiveAsValidator) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ValidatorRightArchiveAsValidatorResult is an enum type
type ValidatorRightArchiveAsValidatorResult string

const (
	ValidatorRightArchiveAsValidatorResultValidatorRight_ArchiveAsValidatorResult ValidatorRightArchiveAsValidatorResult = "ValidatorRight_ArchiveAsValidatorResult"
)

// GetEnumConstructor implements types.ENUM interface
func (e ValidatorRightArchiveAsValidatorResult) GetEnumConstructor() string {
	return string(e)
}

// GetEnumTypeID implements types.ENUM interface
func (e ValidatorRightArchiveAsValidatorResult) GetEnumTypeID() string {
	return fmt.Sprintf("#%s:%s:%s", packageNameSpliceAmulet0116, "Splice.Amulet", "ValidatorRightArchiveAsValidatorResult")
}

// MarshalJSON implements custom JSON marshaling for ValidatorRightArchiveAsValidatorResult using JsonCodec
func (e ValidatorRightArchiveAsValidatorResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(e)
}

// UnmarshalJSON implements custom JSON unmarshaling for ValidatorRightArchiveAsValidatorResult using JsonCodec
func (e *ValidatorRightArchiveAsValidatorResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, e)
}

// Verify interface implementation
var _ ENUM = ValidatorRightArchiveAsValidatorResult("")

func argsToMap(args interface{}) map[string]interface{} {
	if args == nil {
		return map[string]interface{}{}
	}

	if m, ok := args.(map[string]interface{}); ok {
		return m
	}

	// Check if the type has a ToMap method
	type mapper interface {
		ToMap() map[string]interface{}
	}

	if mapper, ok := args.(mapper); ok {
		return mapper.ToMap()
	}

	return map[string]interface{}{
		"args": args,
	}
}
