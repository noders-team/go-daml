package wallet

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"github.com/noders-team/go-daml/pkg/client"
	"github.com/noders-team/go-daml/pkg/model"
	damlModel "github.com/noders-team/go-daml/pkg/model"
	"github.com/noders-team/go-daml/pkg/types"
	gen "github.com/noders-team/go-daml/pkg/wallet/gen_clients"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

const (
	ALLOCATION_FACTORY_INTERFACE_ID          = "#splice-api-token-allocation-instruction-v1:Splice.Api.Token.AllocationInstructionV1:AllocationFactory"
	ALLOCATION_INSTRUCTION_INTERFACE_ID      = "#splice-api-token-allocation-instruction-v1:Splice.Api.Token.AllocationInstructionV1:AllocationInstruction"
	ALLOCATION_REQUEST_INTERFACE_ID          = "#splice-api-token-allocation-request-v1:Splice.Api.Token.AllocationRequestV1:AllocationRequest"
	ALLOCATION_INTERFACE_ID                  = "#splice-api-token-allocation-v1:Splice.Api.Token.AllocationV1:Allocation"
	HOLDING_INTERFACE_ID                     = "#splice-api-token-holding-v1:Splice.Api.Token.HoldingV1:Holding"
	METADATA_INTERFACE_ID                    = "#splice-api-token-metadata-v1:Splice.Api.Token.MetadataV1:AnyContract"
	TRANSFER_FACTORY_INTERFACE_ID            = "#splice-api-token-transfer-instruction-v1:Splice.Api.Token.TransferInstructionV1:TransferFactory"
	TRANSFER_INSTRUCTION_INTERFACE_ID        = "#splice-api-token-transfer-instruction-v1:Splice.Api.Token.TransferInstructionV1:TransferInstruction"
	FEATURED_APP_DELEGATE_PROXY_INTERFACE_ID = "#splice-util-featured-app-proxies:Splice.Util.FeaturedApp.DelegateProxy:DelegateProxy"
	MERGE_DELEGATION_PROPOSAL_TEMPLATE_ID    = "#splice-util-token-standard-wallet:Splice.Util.Token.Wallet.MergeDelegation:MergeDelegationProposal"
	MERGE_DELEGATION_TEMPLATE_ID             = "#splice-util-token-standard-wallet:Splice.Util.Token.Wallet.MergeDelegation:MergeDelegation"
	MERGE_DELEGATION_BATCH_MERGE_UTILITY     = "#splice-util-token-standard-wallet:Splice.Util.Token.Wallet.MergeDelegation:BatchMergeUtility"
)

type TokenStandardController interface {
	ListHoldingUtxos(ctx context.Context, includeLocked bool, limit int) ([]*model.HoldingUTXO, error)
	ListHoldingUtxo(ctx context.Context,
		contractId string,
	) (*model.HoldingUTXO, error)
	GetInputHoldingsCids(
		ctx context.Context,
		sender PartyID,
		instrumentAdmin string,
		instrumentID string,
		amount *decimal.Decimal,
	) ([]string, error)
	GetInputHoldingsCidsForAmount(amount decimal.Decimal, holdings []*model.HoldingUTXO) ([]string, error)
	CreateTransfer(
		ctx context.Context,
		sender PartyID,
		receiver PartyID,
		amount decimal.Decimal,
		inputUtxos []string,
		expiryDate *time.Time,
	) (*model.CommandRequest, error)
	CreateTransferInstruction(ctx context.Context, cid string, choice string) (*model.CommandRequest, error)
	GetBalance(ctx context.Context) (decimal.Decimal, error)
	FetchPendingTransferInstructionView(ctx context.Context) ([]*model.TransferInstruction, error)
}

type tokenStandardController struct {
	damlClient                 *client.DamlBindingClient
	scanProxy                  *client.ScanProxyClient
	userID                     string
	partyID                    atomic.Value
	synchronizerID             atomic.Value
	transferFactoryRegistryUrl atomic.Value
	amuletRulesContractID      atomic.Value
	amuletRulesTemplateID      atomic.Value
	openMiningRoundContractID  atomic.Value
	logger                     zerolog.Logger
}

func NewTokenStandardController(userID string, damlClient *client.DamlBindingClient, sp *client.ScanProxyClient) (TokenStandardController, error) {
	logger := log.Logger.With().
		Str("component", "token-standard-controller").
		Str("userID", userID).
		Logger()

	return &tokenStandardController{
		damlClient: damlClient,
		userID:     userID,
		logger:     logger,
		scanProxy:  sp,
	}, nil
}

func (t *tokenStandardController) SetPartyID(partyID PartyID) {
	t.partyID.Store(partyID)
}

func (t *tokenStandardController) SetSynchronizerID(synchronizerID PartyID) {
	t.synchronizerID.Store(synchronizerID)
}

func (t *tokenStandardController) GetPartyID() (PartyID, error) {
	v := t.partyID.Load()
	if v == nil {
		return "", fmt.Errorf("partyID not set")
	}
	return v.(PartyID), nil
}

func (t *tokenStandardController) GetSynchronizerID() (PartyID, error) {
	v := t.synchronizerID.Load()
	if v == nil {
		return "", fmt.Errorf("synchronizerID not set")
	}
	return v.(PartyID), nil
}

func (t *tokenStandardController) SetTransferFactoryRegistryUrl(url string) {
	t.transferFactoryRegistryUrl.Store(url)
}

func (t *tokenStandardController) GetTransferFactoryRegistryUrl() (string, error) {
	v := t.transferFactoryRegistryUrl.Load()
	if v == nil {
		return "", fmt.Errorf("transferFactoryRegistryUrl not set")
	}
	return v.(string), nil
}

func (t *tokenStandardController) SetAmuletRulesContractID(id string) {
	t.amuletRulesContractID.Store(id)
}

func (t *tokenStandardController) GetAmuletRulesContractID() string {
	v := t.amuletRulesContractID.Load()
	if v == nil {
		return ""
	}
	return v.(string)
}

func (t *tokenStandardController) SetAmuletRulesTemplateID(id string) {
	t.amuletRulesTemplateID.Store(id)
}

func (t *tokenStandardController) GetAmuletRulesTemplateID() string {
	v := t.amuletRulesTemplateID.Load()
	if v == nil {
		return ""
	}
	return v.(string)
}

func (t *tokenStandardController) SetOpenMiningRoundContractID(id string) {
	t.openMiningRoundContractID.Store(id)
}

func (t *tokenStandardController) GetOpenMiningRoundContractID() string {
	v := t.openMiningRoundContractID.Load()
	if v == nil {
		return ""
	}
	return v.(string)
}

func (t *tokenStandardController) Transfer(ctx context.Context, receiver PartyID, amount decimal.Decimal) (*TransferResponse, error) {
	partyID, err := t.GetPartyID()
	if err != nil {
		return nil, err
	}

	syncID, err := t.GetSynchronizerID()
	if err != nil {
		return nil, err
	}

	transferCmd := &damlModel.Command{
		Command: damlModel.ExerciseCommand{
			TemplateID: "#splice-amulet:Splice.Amulet:Amulet",
			Choice:     "Transfer",
			Arguments: map[string]interface{}{
				"newOwner": string(receiver),
				"amount":   amount.String(),
			},
		},
	}

	prepareReq := &damlModel.PrepareSubmissionRequest{
		UserID:         t.userID,
		CommandID:      fmt.Sprintf("transfer-%d", time.Now().UnixNano()),
		Commands:       []*damlModel.Command{transferCmd},
		ActAs:          []string{string(partyID)},
		ReadAs:         []string{},
		SynchronizerID: string(syncID),
	}

	_, err = t.damlClient.InteractiveSubmissionService.PrepareSubmission(ctx, prepareReq)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare transfer: %w", err)
	}

	t.logger.Info().
		Str("receiver", string(receiver)).
		Str("amount", amount.String()).
		Msg("Transfer prepared (signature required)")

	return &TransferResponse{
		SubmissionID: prepareReq.CommandID,
		Amount:       amount,
		Receiver:     receiver,
	}, nil
}

func (t *tokenStandardController) Lock(
	ctx context.Context,
	instrumentAdmin, instrumentID string,
	amount decimal.Decimal, expiresAt time.Time,
) (*model.CommandRequest, error) {
	partyID, err := t.GetPartyID()
	if err != nil {
		return nil, err
	}

	inputCids, err := t.GetInputHoldingsCids(ctx, partyID, instrumentAdmin, instrumentID, &amount)
	if err != nil {
		return nil, err
	}

	amuletRules, round, dsoParty, disclosed, err := t.amuletTransferContext(ctx)
	if err != nil {
		return nil, err
	}

	outputLock := &gen.TimeLock{
		Holders:   []types.PARTY{types.PARTY(partyID)},
		ExpiresAt: types.TIMESTAMP(expiresAt),
	}

	cmd := t.buildAmuletTransferCommand(partyID, partyID, amount, inputCids, amuletRules, round, dsoParty, outputLock)

	t.logger.Info().
		Str("amount", amount.String()).
		Time("expiresAt", expiresAt).
		Msg("lock amulet operation prepared")

	return &model.CommandRequest{
		Command:            cmd,
		DisclosedContracts: disclosed,
	}, nil
}

// TODO????
func (t *tokenStandardController) Unlock(ctx context.Context, lockContractID string) error {
	_, err := t.GetPartyID()
	if err != nil {
		return err
	}

	t.logger.Info().
		Str("lockContractID", lockContractID).
		Msg("Unlock amulet operation")

	return nil
}

// TODO????
func (t *tokenStandardController) Burn(ctx context.Context, amount decimal.Decimal) error {
	_, err := t.GetPartyID()
	if err != nil {
		return err
	}

	t.logger.Info().
		Str("amount", amount.String()).
		Msg("Burn amulet operation")

	return nil
}

func (t *tokenStandardController) GetBalance(ctx context.Context) (decimal.Decimal, error) {
	partyID, err := t.GetPartyID()
	if err != nil {
		return decimal.Zero, err
	}

	contracts, err := client.NewContractQuery[gen.HoldingView](t.damlClient).
		FindContractsByInterface(ctx, string(partyID), gen.IHoldingInterfaceID(nil))
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to query holdings: %w", err)
	}

	balance := decimal.Zero
	now := time.Now()
	for _, c := range contracts {
		if t.isHoldingLocked(c.Data, now) {
			continue
		}
		balance = balance.Add(t.numericToDecimal(c.Data.Amount))
	}

	t.logger.Debug().
		Str("partyID", string(partyID)).
		Int("holdings", len(contracts)).
		Str("balance", balance.String()).
		Msg("balance retrieved via Holding interface query")

	return balance, nil
}

func (t *tokenStandardController) ListContractsByInterface(ctx context.Context, interfaceID string) ([]*damlModel.CreatedEvent, error) {
	partyID, err := t.GetPartyID()
	if err != nil {
		return nil, err
	}

	filterByParty := map[string]*damlModel.Filters{
		string(partyID): {
			Inclusive: &damlModel.InclusiveFilters{
				InterfaceFilters: []*damlModel.InterfaceFilter{
					{
						InterfaceID:             interfaceID,
						IncludeCreatedEventBlob: true,
					},
				},
			},
		},
	}

	req := &damlModel.GetActiveContractsRequest{
		EventFormat: &damlModel.EventFormat{
			Verbose:        true,
			FiltersByParty: filterByParty,
		},
	}

	stream, errChan := t.damlClient.StateService.GetActiveContracts(ctx, req)

	var contracts []*damlModel.CreatedEvent

	for {
		select {
		case resp, ok := <-stream:
			if !ok {
				return contracts, nil
			}
			if entry, ok := resp.ContractEntry.(*damlModel.ActiveContractEntry); ok {
				if entry.ActiveContract != nil && entry.ActiveContract.CreatedEvent != nil {
					contracts = append(contracts, entry.ActiveContract.CreatedEvent)
				}
			}
		case err := <-errChan:
			if err != nil {
				return nil, fmt.Errorf("failed to list contracts by interface: %w", err)
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func (t *tokenStandardController) numericToDecimal(n types.NUMERIC) decimal.Decimal {
	if n == nil {
		return decimal.Zero
	}
	return decimal.NewFromBigInt((*big.Int)(n), -10)
}

func (t *tokenStandardController) lockToMap(l *gen.Lock) map[string]interface{} {
	if l == nil {
		return nil
	}
	m := map[string]interface{}{}
	if l.ExpiresAt != nil {
		m["expiresAt"] = time.Time(*l.ExpiresAt)
	}
	return m
}

func (t *tokenStandardController) ListHoldingUtxos(ctx context.Context, includeLocked bool, limit int) ([]*model.HoldingUTXO, error) {
	partyID, err := t.GetPartyID()
	if err != nil {
		return nil, err
	}

	contracts, err := client.NewContractQuery[gen.HoldingView](t.damlClient).
		FindContractsByInterface(ctx, string(partyID), gen.IHoldingInterfaceID(nil))
	if err != nil {
		return nil, fmt.Errorf("failed to query holdings: %w", err)
	}

	now := time.Now()
	result := make([]*model.HoldingUTXO, 0, len(contracts))
	for _, c := range contracts {
		if !includeLocked && t.isHoldingLocked(c.Data, now) {
			continue
		}
		result = append(result, &model.HoldingUTXO{
			ContractID:      c.ContractID,
			Amount:          t.numericToDecimal(c.Data.Amount),
			InstrumentID:    string(c.Data.InstrumentId.Id),
			InstrumentAdmin: string(c.Data.InstrumentId.Admin),
			Owner:           string(c.Data.Owner),
			Lock:            t.lockToMap(c.Data.Lock),
		})
		if limit > 0 && len(result) >= limit {
			break
		}
	}

	t.logger.Debug().
		Str("partyID", string(partyID)).
		Int("holdings", len(result)).
		Msg("listed holding utxos via Holding interface query")

	return result, nil
}

func (t *tokenStandardController) FetchPendingTransferInstructionView(ctx context.Context) ([]*model.TransferInstruction, error) {
	partyID, err := t.GetPartyID()
	if err != nil {
		return nil, err
	}

	contracts, err := client.NewContractQuery[gen.TransferInstructionView](t.damlClient).
		FindContractsByInterface(ctx, string(partyID), gen.ITransferInstructionInterfaceID(nil))
	if err != nil {
		return nil, fmt.Errorf("failed to query transfer instructions: %w", err)
	}

	instructions := make([]*model.TransferInstruction, 0, len(contracts))
	for _, contract := range contracts {
		view := contract.Data
		if view.Status.TransferPendingReceiverAcceptance == nil {
			continue
		}

		instructions = append(instructions, &model.TransferInstruction{
			ContractID:       contract.ContractID,
			CreatedEventBlob: contract.CreatedEventBlob,
			Sender:           string(view.Transfer.Sender),
			Receiver:         string(view.Transfer.Receiver),
			Amount:           t.numericToDecimal(view.Transfer.Amount),
			Memo:             view.Transfer.Meta.Values[MemoKey],
		})
	}

	return instructions, nil
}

type InputAmuletVariant struct {
	ContractID string
}

func (v InputAmuletVariant) GetVariantTag() string {
	return "InputAmulet"
}

func (v InputAmuletVariant) GetVariantValue() interface{} {
	return types.CONTRACT_ID(v.ContractID)
}

func decimalToNumeric(d decimal.Decimal) types.NUMERIC {
	scaled := d.Mul(decimal.NewFromInt(10000000000))
	return types.NUMERIC(scaled.BigInt())
}

func (t *tokenStandardController) CreateTransfer(
	ctx context.Context,
	sender PartyID,
	receiver PartyID,
	amount decimal.Decimal,
	inputUtxos []string,
	expiryDate *time.Time,
) (*model.CommandRequest, error) {
	if len(inputUtxos) == 0 {
		return nil, fmt.Errorf("no utxos available for transfer")
	}

	amuletRules, round, dsoParty, disclosed, err := t.amuletTransferContext(ctx)
	if err != nil {
		return nil, err
	}

	// AmuletRules_Transfer has no executeBefore field; an expiry is expressed by
	// time-locking the transferred output to the receiver until expiryDate.
	var outputLock *gen.TimeLock
	if expiryDate != nil {
		outputLock = &gen.TimeLock{
			Holders:   []types.PARTY{types.PARTY(receiver)},
			ExpiresAt: types.TIMESTAMP(*expiryDate),
		}
	}

	cmd := t.buildAmuletTransferCommand(sender, receiver, amount, inputUtxos, amuletRules, round, dsoParty, outputLock)

	return &model.CommandRequest{
		Command:            cmd,
		DisclosedContracts: disclosed,
	}, nil
}

// amuletRulesContract fetches the current AmuletRules contract from the
// scan-proxy (contract id + template id + disclosed blob).
func (t *tokenStandardController) amuletRulesContract(ctx context.Context) (*client.ScanContract, error) {
	if t.scanProxy == nil {
		return nil, fmt.Errorf("scan-proxy client not configured; call SetScanProxyClient")
	}
	rules, err := t.scanProxy.GetAmuletRules(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get AmuletRules from scan-proxy: %w", err)
	}
	return rules, nil
}

// amuletTransferContext fetches the AmuletRules and active OpenMiningRound
// contracts and the DSO party from the scan-proxy, returning both contracts as
// disclosed contracts for an AmuletRules_Transfer submission.
func (t *tokenStandardController) amuletTransferContext(ctx context.Context) (amuletRules, round *client.ScanContract, dsoParty string, disclosed []*damlModel.DisclosedContract, err error) {
	if t.scanProxy == nil {
		return nil, nil, "", nil, fmt.Errorf("scan-proxy client not configured; call SetScanProxyClient")
	}

	amuletRules, err = t.scanProxy.GetAmuletRules(ctx)
	if err != nil {
		return nil, nil, "", nil, fmt.Errorf("failed to get AmuletRules from scan-proxy: %w", err)
	}
	round, err = t.scanProxy.GetActiveOpenMiningRound(ctx)
	if err != nil {
		return nil, nil, "", nil, fmt.Errorf("failed to get open mining round from scan-proxy: %w", err)
	}
	if round == nil {
		return nil, nil, "", nil, fmt.Errorf("no active open mining round available")
	}
	dsoParty, err = t.scanProxy.GetDSOPartyID(ctx)
	if err != nil {
		return nil, nil, "", nil, fmt.Errorf("failed to get DSO party from scan-proxy: %w", err)
	}

	disclosed = make([]*damlModel.DisclosedContract, 0, 2)
	if dc := amuletRules.ToDisclosed(); dc != nil {
		disclosed = append(disclosed, dc)
	}
	if dc := round.ToDisclosed(); dc != nil {
		disclosed = append(disclosed, dc)
	}
	return amuletRules, round, dsoParty, disclosed, nil
}

func (t *tokenStandardController) buildAmuletTransferCommand(
	sender, receiver PartyID,
	amount decimal.Decimal,
	inputCids []string,
	amuletRules, round *client.ScanContract,
	dsoParty string,
	outputLock *gen.TimeLock,
) *damlModel.Command {
	inputs := make([]gen.TransferInput, 0, len(inputCids))
	for _, utxo := range inputCids {
		cid := types.CONTRACT_ID(utxo)
		inputs = append(inputs, gen.TransferInput{InputAmulet: &cid})
	}

	expectedDso := types.PARTY(dsoParty)
	choiceArgs := gen.AmuletRulesTransfer{
		Transfer: gen.Transfer{
			Sender:   types.PARTY(sender),
			Provider: types.PARTY(sender),
			Inputs:   inputs,
			Outputs: []gen.TransferOutput{{
				Receiver:         types.PARTY(receiver),
				ReceiverFeeRatio: types.NewNumericFromDecimal(decimal.Zero),
				Amount:           types.NewNumericFromDecimal(amount),
				Lock:             outputLock,
			}},
		},
		Context: gen.TransferContext{
			OpenMiningRound:     types.CONTRACT_ID(round.ContractID),
			IssuingMiningRounds: types.GENMAP{},
			ValidatorRights:     types.GENMAP{},
		},
		ExpectedDso: &expectedDso,
	}

	exercise := gen.AmuletRules{}.AmuletRulesTransfer(amuletRules.ContractID, choiceArgs)
	exercise.TemplateID = amuletRules.TemplateID
	return &damlModel.Command{Command: exercise}
}

type CreateTapResult struct {
	Command            *damlModel.Command
	DisclosedContracts []*damlModel.DisclosedContract
}

func (t *tokenStandardController) CreateTap(
	ctx context.Context,
	receiver PartyID,
	amount decimal.Decimal,
	instrumentAdmin string,
	instrumentID string,
) (*CreateTapResult, error) {
	amuletRulesTemplateID, amuletRulesContractID, err := t.findAmuletRulesContract(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find AmuletRules contract: %w", err)
	}

	openMiningRoundContractID := t.GetOpenMiningRoundContractID()
	if openMiningRoundContractID == "" {
		openMiningRoundContractID = os.Getenv("OPEN_MINING_ROUND_CONTRACT_ID")
	}
	if openMiningRoundContractID == "" {
		return nil, fmt.Errorf("openMiningRoundContractID not set - OpenMiningRound needs to be bootstrapped first")
	}

	tapCmd := &damlModel.Command{
		Command: &damlModel.ExerciseCommand{
			TemplateID: amuletRulesTemplateID,
			ContractID: amuletRulesContractID,
			Choice:     "AmuletRules_DevNet_Tap",
			Arguments: map[string]interface{}{
				"receiver":  types.PARTY(string(receiver)),
				"amount":    decimalToNumeric(amount),
				"openRound": types.CONTRACT_ID(openMiningRoundContractID),
			},
		},
	}

	return &CreateTapResult{
		Command:            tapCmd,
		DisclosedContracts: []*damlModel.DisclosedContract{},
	}, nil
}

func (t *tokenStandardController) ExerciseTransferInstructionChoice(
	ctx context.Context,
	transferInstructionCid string,
	choice string,
) (*model.CommandRequest, error) {
	exerciseCmd := &damlModel.Command{
		Command: &damlModel.ExerciseCommand{
			ContractID: transferInstructionCid,
			TemplateID: TRANSFER_INSTRUCTION_INTERFACE_ID,
			Choice:     choice,
			Arguments:  map[string]interface{}{},
		},
	}

	return &model.CommandRequest{
		Command:            exerciseCmd,
		DisclosedContracts: []*damlModel.DisclosedContract{},
	}, nil
}

func (t *tokenStandardController) ListHoldingUtxo(ctx context.Context,
	contractId string,
) (*model.HoldingUTXO, error) {
	utxos, err := t.ListHoldingUtxos(ctx, true, 0)
	if err != nil {
		return nil, err
	}

	for _, utxo := range utxos {
		if utxo.ContractID == contractId {
			return utxo, nil
		}
	}

	return nil, fmt.Errorf("holding with contractId %s not found", contractId)
}

func (t *tokenStandardController) MergeHoldingUtxos(ctx context.Context,
	nodeLimit int, partyID PartyID,
	inputUtxos []*model.HoldingUTXO,
) (*MergeUtxosResult, error) {
	if partyID == "" {
		var err error
		partyID, err = t.GetPartyID()
		if err != nil {
			return nil, err
		}
	}

	var utxos []*model.HoldingUTXO
	if len(inputUtxos) > 0 {
		utxos = inputUtxos
	} else {
		var err error
		utxos, err = t.ListHoldingUtxos(ctx, false, nodeLimit)
		if err != nil {
			return nil, err
		}
	}

	// Fetch AmuletRules / OpenMiningRound / DSO once and reuse across every batch
	// (mirrors ambo's mergeInstrument, which resolves the app context once).
	amuletRules, round, dsoParty, disclosed, err := t.amuletTransferContext(ctx)
	if err != nil {
		return nil, err
	}

	utxosByInstrument := make(map[string][]*model.HoldingUTXO)
	for _, utxo := range utxos {
		key := utxo.InstrumentID + "::" + utxo.InstrumentAdmin
		utxosByInstrument[key] = append(utxosByInstrument[key], utxo)
	}

	var allCommands []*damlModel.Command
	transferInputUtxoLimit := 100

	for _, group := range utxosByInstrument {
		transfers := (len(group) + transferInputUtxoLimit - 1) / transferInputUtxoLimit

		for i := 0; i < transfers; i++ {
			start := i * transferInputUtxoLimit
			end := start + transferInputUtxoLimit
			if end > len(group) {
				end = len(group)
			}

			var totalAmount decimal.Decimal
			var inputCids []string
			for _, utxo := range group[start:end] {
				totalAmount = totalAmount.Add(utxo.Amount)
				inputCids = append(inputCids, utxo.ContractID)
			}

			// A self-transfer of all inputs into a single output merges them.
			cmd := t.buildAmuletTransferCommand(partyID, partyID, totalAmount, inputCids, amuletRules, round, dsoParty, nil)
			allCommands = append(allCommands, cmd)
		}
	}

	return &MergeUtxosResult{
		Commands:           allCommands,
		DisclosedContracts: disclosed,
	}, nil
}

func (t *tokenStandardController) FetchPendingAllocationInstructionView(ctx context.Context) ([]*model.AllocationInstruction, error) {
	partyID, err := t.GetPartyID()
	if err != nil {
		return nil, err
	}

	contracts, err := client.NewContractQuery[gen.AllocationInstructionView](t.damlClient).
		FindContractsByInterface(ctx, string(partyID), gen.IAllocationInstructionInterfaceID(nil))
	if err != nil {
		return nil, fmt.Errorf("failed to query allocation instructions: %w", err)
	}

	instructions := make([]*model.AllocationInstruction, 0, len(contracts))
	for _, contract := range contracts {
		view := contract.Data
		instructions = append(instructions, &model.AllocationInstruction{
			ContractID:       contract.ContractID,
			Provider:         string(view.Allocation.Settlement.Executor),
			Specification:    view.Allocation.ToMap(),
			CreatedEventBlob: contract.CreatedEventBlob,
		})
	}

	return instructions, nil
}

func (t *tokenStandardController) FetchPendingAllocationRequestView(ctx context.Context) ([]*model.AllocationRequest, error) {
	partyID, err := t.GetPartyID()
	if err != nil {
		return nil, err
	}

	contracts, err := client.NewContractQuery[gen.AllocationRequestView](t.damlClient).
		FindContractsByInterface(ctx, string(partyID), gen.IAllocationRequestInterfaceID(nil))
	if err != nil {
		return nil, fmt.Errorf("failed to query allocation requests: %w", err)
	}

	requests := make([]*model.AllocationRequest, 0, len(contracts))
	for _, contract := range contracts {
		view := contract.Data
		requests = append(requests, &model.AllocationRequest{
			ContractID:       contract.ContractID,
			Requester:        string(view.Settlement.Executor),
			Specification:    view.ToMap(),
			CreatedEventBlob: contract.CreatedEventBlob,
		})
	}

	return requests, nil
}

func (t *tokenStandardController) FetchPendingAllocationView(ctx context.Context) ([]*Allocation, error) {
	contracts, err := t.ListContractsByInterface(ctx, ALLOCATION_INTERFACE_ID)
	if err != nil {
		return nil, err
	}

	var allocations []*Allocation
	for _, contract := range contracts {
		args, ok := contract.CreateArguments.(map[string]interface{})
		if !ok {
			continue
		}

		allocation := &Allocation{
			ContractID:       contract.ContractID,
			CreatedEventBlob: contract.CreatedEventBlob,
		}

		if provider, ok := args["provider"].(string); ok {
			allocation.Provider = provider
		}
		if receiver, ok := args["receiver"].(string); ok {
			allocation.Receiver = receiver
		}
		if amountVal, ok := args["amount"]; ok {
			if amountStr, ok := amountVal.(string); ok {
				allocation.Amount, _ = decimal.NewFromString(amountStr)
			}
		}

		allocations = append(allocations, allocation)
	}

	return allocations, nil
}

func (t *tokenStandardController) CreateAllocationInstruction(
	ctx context.Context,
	allocationSpecification map[string]interface{},
	expectedAdmin string,
	inputUtxos []string,
	requestedAt string,
) (*model.CommandRequest, error) {
	cmd := &damlModel.Command{
		Command: &damlModel.ExerciseCommand{
			TemplateID: ALLOCATION_FACTORY_INTERFACE_ID,
			Choice:     "CreateAllocationInstruction",
			Arguments: map[string]interface{}{
				"specification": allocationSpecification,
				"expectedAdmin": expectedAdmin,
				"inputs":        inputUtxos,
				"requestedAt":   requestedAt,
			},
		},
	}

	return &model.CommandRequest{
		Command:            cmd,
		DisclosedContracts: []*damlModel.DisclosedContract{},
	}, nil
}

func (t *tokenStandardController) ExerciseAllocationChoice(
	_ context.Context,
	allocationCid string,
	choice string,
) (*model.CommandRequest, error) {
	cmd := &damlModel.Command{
		Command: &damlModel.ExerciseCommand{
			ContractID: allocationCid,
			TemplateID: ALLOCATION_INTERFACE_ID,
			Choice:     choice,
			Arguments:  map[string]interface{}{},
		},
	}

	return &model.CommandRequest{
		Command:            cmd,
		DisclosedContracts: []*damlModel.DisclosedContract{},
	}, nil
}

func (t *tokenStandardController) ExerciseAllocationInstructionChoice(
	ctx context.Context,
	allocationInstructionCid string,
	choice string,
) (*model.CommandRequest, error) {
	cmd := &damlModel.Command{
		Command: &damlModel.ExerciseCommand{
			ContractID: allocationInstructionCid,
			TemplateID: ALLOCATION_INSTRUCTION_INTERFACE_ID,
			Choice:     choice,
			Arguments:  map[string]interface{}{},
		},
	}

	return &model.CommandRequest{
		Command:            cmd,
		DisclosedContracts: []*damlModel.DisclosedContract{},
	}, nil
}

func (t *tokenStandardController) ExerciseAllocationRequestChoice(
	ctx context.Context,
	allocationRequestCid string,
	choice string,
	actor PartyID,
) (*model.CommandRequest, error) {
	cmd := &damlModel.Command{
		Command: &damlModel.ExerciseCommand{
			ContractID: allocationRequestCid,
			TemplateID: ALLOCATION_REQUEST_INTERFACE_ID,
			Choice:     choice,
			Arguments: map[string]interface{}{
				"actor": string(actor),
			},
		},
	}

	return &model.CommandRequest{
		Command:            cmd,
		DisclosedContracts: []*damlModel.DisclosedContract{},
	}, nil
}

func (t *tokenStandardController) CreateTransferUsingDelegateProxy(
	ctx context.Context,
	proxyCid string,
	featuredAppRightCid string,
	sender PartyID,
	receiver PartyID,
	amount decimal.Decimal,
	instrumentID string,
	instrumentAdmin string,
	beneficiaries []map[string]interface{},
	inputUtxos []string,
	memo string,
) (*model.CommandRequest, error) {
	cmd := &damlModel.Command{
		Command: &damlModel.ExerciseCommand{
			ContractID: proxyCid,
			TemplateID: FEATURED_APP_DELEGATE_PROXY_INTERFACE_ID,
			Choice:     "CreateTransfer",
			Arguments: map[string]interface{}{
				"featuredAppRightCid": featuredAppRightCid,
				"sender":              string(sender),
				"receiver":            string(receiver),
				"amount":              amount.String(),
				"instrumentId":        instrumentID,
				"instrumentAdmin":     instrumentAdmin,
				"beneficiaries":       beneficiaries,
				"inputs":              inputUtxos,
				"memo":                memo,
			},
		},
	}

	return &model.CommandRequest{
		Command:            cmd,
		DisclosedContracts: []*damlModel.DisclosedContract{},
	}, nil
}

func (t *tokenStandardController) ExerciseTransferInstructionChoiceWithDelegate(
	ctx context.Context,
	transferInstructionCid string,
	choice string,
	proxyCid string,
	featuredAppRightCid string,
	beneficiaries []map[string]interface{},
) (*model.CommandRequest, error) {
	cmd := &damlModel.Command{
		Command: &damlModel.ExerciseCommand{
			ContractID: proxyCid,
			TemplateID: FEATURED_APP_DELEGATE_PROXY_INTERFACE_ID,
			Choice:     "ExerciseTransferInstructionChoice",
			Arguments: map[string]interface{}{
				"transferInstructionCid": transferInstructionCid,
				"choice":                 choice,
				"featuredAppRightCid":    featuredAppRightCid,
				"beneficiaries":          beneficiaries,
			},
		},
	}

	return &model.CommandRequest{
		Command:            cmd,
		DisclosedContracts: []*damlModel.DisclosedContract{},
	}, nil
}

type TransferPreapproval struct {
	ReceiverID PartyID
	ExpiresAt  time.Time
	Dso        PartyID
	ContractID string
	TemplateID string
}

// Deprecated: requires the scan-proxy HTTP API; not reachable through the ledger client.
func (t *tokenStandardController) GetTransferPreApprovalByParty(ctx context.Context, receiverID PartyID, instrumentID string) (*TransferPreapproval, error) {
	t.logger.Info().Str("receiverId", string(receiverID)).Str("instrumentId", instrumentID).Msg("Getting transfer preapproval by party")
	return nil, fmt.Errorf("scan proxy API call not implemented - requires HTTP client")
}

func (t *tokenStandardController) CreateCancelTransferPreapproval(
	ctx context.Context,
	contractID string,
	templateID string,
	actor PartyID,
) (*model.CommandRequest, error) {
	cmd := &damlModel.Command{
		Command: &damlModel.ExerciseCommand{
			ContractID: contractID,
			TemplateID: templateID,
			Choice:     "TransferPreapproval_Cancel",
			Arguments: map[string]interface{}{
				"actor": string(actor),
			},
		},
	}

	return &model.CommandRequest{
		Command:            cmd,
		DisclosedContracts: []*damlModel.DisclosedContract{},
	}, nil
}

func (t *tokenStandardController) CreateRenewTransferPreapproval(
	ctx context.Context,
	contractID string,
	templateID string,
	provider PartyID,
	newExpiresAt *time.Time,
	inputUtxos []string,
) (*model.CommandRequest, error) {
	syncID, err := t.GetSynchronizerID()
	if err != nil {
		return nil, err
	}

	args := map[string]interface{}{
		"provider":       string(provider),
		"synchronizerId": string(syncID),
	}

	if newExpiresAt != nil {
		args["newExpiresAt"] = newExpiresAt.Format(time.RFC3339)
	}

	if len(inputUtxos) > 0 {
		args["inputUtxos"] = inputUtxos
	}

	cmd := &damlModel.Command{
		Command: &damlModel.ExerciseCommand{
			ContractID: contractID,
			TemplateID: templateID,
			Choice:     "TransferPreapproval_Renew",
			Arguments:  args,
		},
	}

	return &model.CommandRequest{
		Command:            cmd,
		DisclosedContracts: []*damlModel.DisclosedContract{},
	}, nil
}

func (t *tokenStandardController) WaitForPreapprovalFromScanProxy(
	ctx context.Context,
	receiverID PartyID,
	instrumentID string,
	oldCid string,
	expectGone bool,
	intervalMs int,
	timeoutMs int,
) (*TransferPreapproval, error) {
	if intervalMs == 0 {
		intervalMs = 15000
	}
	if timeoutMs == 0 {
		timeoutMs = 300000
	}

	deadline := time.Now().Add(time.Duration(timeoutMs) * time.Millisecond)
	interval := time.Duration(intervalMs) * time.Millisecond

	for attempt := 1; time.Now().Before(deadline); attempt++ {
		preapproval, err := t.GetTransferPreApprovalByParty(ctx, receiverID, instrumentID)

		if expectGone {
			if preapproval == nil || err != nil {
				t.logger.Info().Int("attempt", attempt).Msg("Preapproval is no longer visible")
				return nil, nil
			}
			t.logger.Info().Int("attempt", attempt).Str("seenCid", preapproval.ContractID).Msg("Preapproval still visible - polling again")
		} else if preapproval != nil {
			if oldCid == "" {
				t.logger.Info().Int("attempt", attempt).Str("seenCid", preapproval.ContractID).Msg("Preapproval is visible")
				return preapproval, nil
			}
			if preapproval.ContractID != oldCid {
				t.logger.Info().Int("attempt", attempt).Str("oldCid", oldCid).Str("newCid", preapproval.ContractID).Msg("Preapproval CID changed")
				return preapproval, nil
			}
			t.logger.Info().Int("attempt", attempt).Str("seenCid", preapproval.ContractID).Str("oldCid", oldCid).Msg("Preapproval visible but CID unchanged - polling again")
		} else {
			t.logger.Info().Int("attempt", attempt).Msg("Preapproval not visible yet - polling again")
		}

		select {
		case <-time.After(interval):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	waitingFor := "for preapproval to appear"
	if expectGone {
		waitingFor = "for preapproval to disappear"
	} else if oldCid != "" {
		waitingFor = fmt.Sprintf("for preapproval CID to change from %s", oldCid)
	}

	return nil, fmt.Errorf("timed out after %dms waiting %s", timeoutMs, waitingFor)
}

func (t *tokenStandardController) BuyMemberTraffic(
	ctx context.Context,
	buyer PartyID,
	ccAmount decimal.Decimal,
	memberId string,
	inputUtxos []string,
	migrationId int,
) (*model.CommandRequest, error) {
	syncID, err := t.GetSynchronizerID()
	if err != nil {
		return nil, err
	}

	amuletRules, err := t.amuletRulesContract(ctx)
	if err != nil {
		return nil, err
	}

	cmd := &damlModel.Command{
		Command: &damlModel.ExerciseCommand{
			TemplateID: amuletRules.TemplateID,
			ContractID: amuletRules.ContractID,
			Choice:     "AmuletRules_BuyMemberTraffic",
			Arguments: map[string]interface{}{
				"buyer":          string(buyer),
				"ccAmount":       ccAmount.String(),
				"memberId":       memberId,
				"inputs":         inputUtxos,
				"migrationId":    migrationId,
				"synchronizerId": string(syncID),
			},
		},
	}

	disclosed := make([]*damlModel.DisclosedContract, 0, 1)
	if dc := amuletRules.ToDisclosed(); dc != nil {
		disclosed = append(disclosed, dc)
	}

	return &model.CommandRequest{
		Command:            cmd,
		DisclosedContracts: disclosed,
	}, nil
}

func (t *tokenStandardController) SelfGrantFeatureAppRights(ctx context.Context) (*model.CommandRequest, error) {
	partyID, err := t.GetPartyID()
	if err != nil {
		return nil, err
	}

	syncID, err := t.GetSynchronizerID()
	if err != nil {
		return nil, err
	}

	amuletRules, err := t.amuletRulesContract(ctx)
	if err != nil {
		return nil, err
	}

	cmd := &damlModel.Command{
		Command: &damlModel.ExerciseCommand{
			TemplateID: amuletRules.TemplateID,
			ContractID: amuletRules.ContractID,
			Choice:     "AmuletRules_DevNet_FeatureApp",
			Arguments: map[string]interface{}{
				"provider":       string(partyID),
				"synchronizerId": string(syncID),
			},
		},
	}

	disclosed := make([]*damlModel.DisclosedContract, 0, 1)
	if dc := amuletRules.ToDisclosed(); dc != nil {
		disclosed = append(disclosed, dc)
	}

	return &model.CommandRequest{
		Command:            cmd,
		DisclosedContracts: disclosed,
	}, nil
}

func (t *tokenStandardController) LookupFeaturedApps(ctx context.Context, maxRetries int, delayMs int) (*model.FeaturedAppRight, error) {
	if maxRetries == 0 {
		maxRetries = 10
	}
	if delayMs == 0 {
		delayMs = 5000
	}

	partyID, err := t.GetPartyID()
	if err != nil {
		return nil, err
	}
	if t.scanProxy == nil {
		return nil, fmt.Errorf("scan-proxy client not configured; call SetScanProxyClient")
	}

	for attempt := 1; attempt <= maxRetries; attempt++ {
		contract, err := t.scanProxy.GetFeaturedAppByProvider(ctx, string(partyID))
		if err == nil && contract != nil {
			var blob []byte
			if dc := contract.ToDisclosed(); dc != nil {
				blob = dc.CreatedEventBlob
			}
			return &model.FeaturedAppRight{
				TemplateID:       contract.TemplateID,
				ContractID:       contract.ContractID,
				Payload:          contract.Payload,
				CreatedEventBlob: blob,
			}, nil
		}

		t.logger.Info().Int("attempt", attempt).Msg("Lookup featured apps returned undefined, retrying again...")

		if attempt < maxRetries {
			time.Sleep(time.Duration(delayMs) * time.Millisecond)
		}
	}

	return nil, nil
}

func (t *tokenStandardController) GrantFeatureAppRightsForInternalParty(ctx context.Context) (*model.FeaturedAppRight, error) {
	featuredAppRights, err := t.LookupFeaturedApps(ctx, 1, 1000)
	if err != nil {
		return nil, err
	}

	if featuredAppRights != nil {
		return featuredAppRights, nil
	}

	result, err := t.SelfGrantFeatureAppRights(ctx)
	if err != nil {
		return nil, err
	}

	partyID, err := t.GetPartyID()
	if err != nil {
		return nil, err
	}

	submitReq := &damlModel.SubmitRequest{
		Commands: &damlModel.Commands{
			UserID:    t.userID,
			CommandID: fmt.Sprintf("feature-app-%d", time.Now().UnixNano()),
			Commands:  []*damlModel.Command{result.Command},
			ActAs:     []string{string(partyID)},
			ReadAs:    []string{},
		},
	}

	_, err = t.damlClient.CommandSubmission.Submit(ctx, submitReq)
	if err != nil {
		return nil, fmt.Errorf("failed to submit feature app grant: %w", err)
	}

	return t.LookupFeaturedApps(ctx, 5, 1000)
}

func (t *tokenStandardController) CreateAndSubmitTapInternal(
	ctx context.Context,
	receiver PartyID,
	amount decimal.Decimal,
	instrumentID string,
	instrumentAdmin string,
) (map[string]interface{}, error) {
	result, err := t.CreateTap(ctx, receiver, amount, instrumentAdmin, instrumentID)
	if err != nil {
		return nil, err
	}

	partyID, err := t.GetPartyID()
	if err != nil {
		return nil, err
	}

	cmdID := fmt.Sprintf("tap-%d", time.Now().UnixNano())
	submitReq := &damlModel.SubmitAndWaitRequest{
		Commands: &damlModel.Commands{
			UserID:    t.userID,
			CommandID: cmdID,
			Commands:  []*damlModel.Command{result.Command},
			ActAs:     []string{string(partyID)},
			ReadAs:    []string{},
		},
	}

	resp, err := t.damlClient.CommandService.SubmitAndWait(ctx, submitReq)
	if err != nil {
		return nil, fmt.Errorf("failed to submit tap: %w", err)
	}

	return map[string]interface{}{
		"commandId":        cmdID,
		"updateId":         resp.UpdateID,
		"completionOffset": resp.CompletionOffset,
	}, nil
}

func (t *tokenStandardController) CreateBatchMergeUtility(ctx context.Context) (*damlModel.Command, error) {
	partyID, err := t.GetPartyID()
	if err != nil {
		return nil, err
	}

	return &damlModel.Command{
		Command: &damlModel.CreateCommand{
			TemplateID: MERGE_DELEGATION_BATCH_MERGE_UTILITY,
			Arguments: map[string]interface{}{
				"operator": string(partyID),
			},
		},
	}, nil
}

func (t *tokenStandardController) findAmuletRulesContract(ctx context.Context) (string, string, error) {
	templateID := t.GetAmuletRulesTemplateID()
	contractID := t.GetAmuletRulesContractID()

	if templateID == "" {
		templateID = os.Getenv("AMULET_RULES_TEMPLATE_ID")
	}
	if contractID == "" {
		contractID = os.Getenv("AMULET_RULES_CONTRACT_ID")
	}

	if templateID != "" && contractID != "" {
		t.logger.Info().
			Str("templateID", templateID).
			Str("contractID", contractID).
			Msg("Using AmuletRules contract from configured values")
		return templateID, contractID, nil
	}

	packages, err := t.damlClient.PackageMng.ListKnownPackages(ctx)
	if err != nil {
		return "", "", fmt.Errorf("failed to list packages: %w", err)
	}

	var spliceAmuletPkgID string
	for _, pkg := range packages {
		if pkg.Name == "splice-amulet" {
			spliceAmuletPkgID = pkg.PackageID
			break
		}
	}

	if spliceAmuletPkgID == "" {
		return "", "", fmt.Errorf("splice-amulet package not found")
	}

	possibleTemplateIDs := []string{
		fmt.Sprintf("%s:Splice.AmuletRules:AmuletRules", spliceAmuletPkgID),
		fmt.Sprintf("%s:Splice.Amulet:AmuletRules", spliceAmuletPkgID),
		fmt.Sprintf("%s:Splice.Amulet.AmuletRules:AmuletRules", spliceAmuletPkgID),
	}

	partyID, err := t.GetPartyID()
	if err != nil {
		return "", "", err
	}

	for _, templateID := range possibleTemplateIDs {
		t.logger.Info().
			Str("templateID", templateID).
			Str("partyID", string(partyID)).
			Msg("Trying to find AmuletRules with template ID")

		req := &damlModel.GetActiveContractsRequest{
			EventFormat: &damlModel.EventFormat{
				Verbose: true,
				FiltersByParty: map[string]*damlModel.Filters{
					string(partyID): {
						Inclusive: &damlModel.InclusiveFilters{
							TemplateFilters: []*damlModel.TemplateFilter{
								{
									TemplateID:              templateID,
									IncludeCreatedEventBlob: false,
								},
							},
						},
					},
				},
			},
		}

		stream, errChan := t.damlClient.StateService.GetActiveContracts(ctx, req)

		var foundContract *damlModel.CreatedEvent
	streamLoop:
		for {
			select {
			case resp, ok := <-stream:
				if !ok {
					if foundContract != nil {
						t.logger.Info().
							Str("templateID", templateID).
							Str("contractID", foundContract.ContractID).
							Msg("Found AmuletRules contract")
						return templateID, foundContract.ContractID, nil
					}
					t.logger.Debug().
						Str("templateID", templateID).
						Str("partyID", string(partyID)).
						Msg("Stream closed, no contract found with this template ID")
					break streamLoop
				}
				if entry, ok := resp.ContractEntry.(*damlModel.ActiveContractEntry); ok {
					if entry.ActiveContract != nil && entry.ActiveContract.CreatedEvent != nil {
						contract := entry.ActiveContract.CreatedEvent
						t.logger.Info().
							Str("templateID", templateID).
							Str("contractID", contract.ContractID).
							Str("partyID", string(partyID)).
							Msg("Received contract from stream")
						foundContract = contract
					}
				}
			case err := <-errChan:
				if err != nil {
					t.logger.Warn().
						Err(err).
						Str("templateID", templateID).
						Str("partyID", string(partyID)).
						Msg("Error querying for template, trying next")
					break streamLoop
				}
			case <-ctx.Done():
				return "", "", ctx.Err()
			}
		}
	}

	return "", "", fmt.Errorf("AmuletRules contract not found - it may need to be initialized first. Attempted template IDs: %v", possibleTemplateIDs)
}

func (t *tokenStandardController) CreateMergeDelegationProposal(ctx context.Context, delegate PartyID, metadata map[string]interface{}) (*damlModel.Command, error) {
	partyID, err := t.GetPartyID()
	if err != nil {
		return nil, err
	}

	if metadata == nil {
		metadata = map[string]interface{}{"values": map[string]interface{}{}}
	}

	return &damlModel.Command{
		Command: &damlModel.CreateCommand{
			TemplateID: MERGE_DELEGATION_PROPOSAL_TEMPLATE_ID,
			Arguments: map[string]interface{}{
				"delegation": map[string]interface{}{
					"operator": string(delegate),
					"owner":    string(partyID),
					"meta":     metadata,
				},
			},
		},
	}, nil
}

func (t *tokenStandardController) LookupMergeDelegationProposal(ctx context.Context, ownerParty PartyID) ([]*damlModel.CreatedEvent, error) {
	if ownerParty == "" {
		var err error
		ownerParty, err = t.GetPartyID()
		if err != nil {
			return nil, err
		}
	}

	req := &damlModel.GetActiveContractsRequest{
		EventFormat: &damlModel.EventFormat{
			Verbose: true,
			FiltersByParty: map[string]*damlModel.Filters{
				string(ownerParty): {
					Inclusive: &damlModel.InclusiveFilters{
						TemplateFilters: []*damlModel.TemplateFilter{
							{TemplateID: MERGE_DELEGATION_PROPOSAL_TEMPLATE_ID},
						},
					},
				},
			},
		},
	}

	stream, errChan := t.damlClient.StateService.GetActiveContracts(ctx, req)

	var contracts []*damlModel.CreatedEvent
	for {
		select {
		case resp, ok := <-stream:
			if !ok {
				return contracts, nil
			}
			if entry, ok := resp.ContractEntry.(*damlModel.ActiveContractEntry); ok {
				if entry.ActiveContract != nil && entry.ActiveContract.CreatedEvent != nil {
					contracts = append(contracts, entry.ActiveContract.CreatedEvent)
				}
			}
		case err := <-errChan:
			if err != nil {
				return nil, err
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func (t *tokenStandardController) ApproveMergeDelegationProposal(ctx context.Context, ownerParty PartyID) (*model.CommandRequest, error) {
	proposals, err := t.LookupMergeDelegationProposal(ctx, ownerParty)
	if err != nil {
		return nil, err
	}

	if len(proposals) == 0 {
		return nil, fmt.Errorf("no merge delegation proposal found")
	}

	proposal := proposals[0]

	cmd := &damlModel.Command{
		Command: &damlModel.ExerciseCommand{
			ContractID: proposal.ContractID,
			TemplateID: MERGE_DELEGATION_PROPOSAL_TEMPLATE_ID,
			Choice:     "MergeDelegationProposal_Accept",
			Arguments:  map[string]interface{}{},
		},
	}

	disclosed := &damlModel.DisclosedContract{
		TemplateID:       proposal.TemplateID,
		ContractID:       proposal.ContractID,
		CreatedEventBlob: proposal.CreatedEventBlob,
	}

	return &model.CommandRequest{
		Command:            cmd,
		DisclosedContracts: []*damlModel.DisclosedContract{disclosed},
	}, nil
}

type PartyID string

type MergeUtxosResult struct {
	Commands           []*damlModel.Command
	DisclosedContracts []*damlModel.DisclosedContract
}

type Allocation struct {
	ContractID       string
	Provider         string
	Receiver         string
	Amount           decimal.Decimal
	CreatedEventBlob []byte
}

type TransferResponse struct {
	SubmissionID string
	Amount       decimal.Decimal
	Receiver     PartyID
}

// ChoiceContext mirrors the registry choice-context payload consumed by the
// *FromContext builders. It is normally produced by a token-standard registry
// HTTP call (see the Deprecated Fetch* methods) and passed in by the caller.
type ChoiceContext struct {
	ChoiceContextData  map[string]interface{}
	DisclosedContracts []*damlModel.DisclosedContract
}

// Beneficiary is a featured-app reward split used by the delegate-proxy choices.
type Beneficiary struct {
	Beneficiary PartyID
	Weight      float64
}

// PrettyContract is a decoded active contract together with its interface view.
type PrettyContract[T any] struct {
	ContractID      string
	InterfaceView   T
	CreatedEvent    *damlModel.CreatedEvent
	FetchedAtOffset int64
}

const (
	// MemoKey is the metadata key used to carry a free-text transfer reason.
	MemoKey = "splice.lfdecentralizedtrust.org/reason"

	// requestedAtSkew back-dates requestedAt to tolerate clock skew, mirroring
	// REQUESTED_AT_SKEW_MS in token-standard-service.ts.
	requestedAtSkew = 60 * time.Second
)

func emptyExtraArgs() gen.ExtraArgs {
	return gen.ExtraArgs{
		Context: gen.ChoiceContext{Values: types.TEXTMAP{}},
		Meta:    gen.Metadata{Values: types.TEXTMAP{}},
	}
}

func newExerciseResult(cmd *damlModel.ExerciseCommand, disclosed []*damlModel.DisclosedContract) *model.CommandRequest {
	if disclosed == nil {
		disclosed = []*damlModel.DisclosedContract{}
	}
	return &model.CommandRequest{
		Command:            &damlModel.Command{Command: cmd},
		DisclosedContracts: disclosed,
	}
}

func toContractIDs(ids []string) []types.CONTRACT_ID {
	out := make([]types.CONTRACT_ID, len(ids))
	for i, id := range ids {
		out[i] = types.CONTRACT_ID(id)
	}
	return out
}

// injectContext reproduces the token-standard-service pattern of assigning the
// registry-provided choice context onto extraArgs before submitting a choice.
func injectContext(ctxData map[string]interface{}) map[string]interface{} {
	if ctxData == nil {
		ctxData = map[string]interface{}{}
	}
	return map[string]interface{}{
		"context": map[string]interface{}{
			"values": map[string]interface{}{"_type": "textmap", "value": ctxData},
		},
		"meta": map[string]interface{}{
			"values": map[string]interface{}{"_type": "textmap", "value": map[string]interface{}{}},
		},
	}
}

func (t *tokenStandardController) isHoldingLocked(view gen.HoldingView, now time.Time) bool {
	if view.Lock == nil {
		return false
	}
	if view.Lock.ExpiresAt == nil {
		return true
	}
	return now.Before(time.Time(*view.Lock.ExpiresAt))
}

// GetInputHoldingsCids returns the contract IDs of the holdings that should fund
// a transfer or allocation for sender, mirroring CoreService.getInputHoldingsCids
// from token-standard-service.ts.
//
// It queries sender's holdings through the token-standard Holding interface; locked
// holdings are already excluded by ListHoldingUtxos. When both instrumentAdmin and
// instrumentID are set, the holdings are further restricted to that instrument
// (matched case-insensitively).
//
// If amount is nil, the IDs of all matching holdings are returned. Otherwise
// largest-first coin selection is applied (see selectForAmount): a holding whose
// amount equals amount exactly is preferred, and it returns an error when the
// available holdings cannot cover amount or when more than 100 utxos would be
// required for a single transaction.
//
// As a side effect, sender is set as the controller's active party.
func (t *tokenStandardController) GetInputHoldingsCids(
	ctx context.Context,
	sender PartyID,
	instrumentAdmin string,
	instrumentID string,
	amount *decimal.Decimal,
) ([]string, error) {
	t.SetPartyID(sender)

	holdings, err := t.ListHoldingUtxos(ctx, false, 0)
	if err != nil {
		return nil, err
	}
	if len(holdings) == 0 {
		return nil, fmt.Errorf("sender has no holdings, so transfer can't be executed")
	}

	if instrumentAdmin != "" && instrumentID != "" {
		filtered := make([]*model.HoldingUTXO, 0)
		for _, h := range holdings {
			if strings.EqualFold(h.InstrumentID, instrumentID) && strings.EqualFold(h.InstrumentAdmin, instrumentAdmin) {
				filtered = append(filtered, h)
			}
		}
		holdings = filtered
	}

	if amount == nil {
		cids := make([]string, 0, len(holdings))
		for _, h := range holdings {
			cids = append(cids, h.ContractID)
		}
		return cids, nil
	}

	return t.GetInputHoldingsCidsForAmount(*amount, holdings)
}

func (t *tokenStandardController) GetInputHoldingsCidsForAmount(amount decimal.Decimal, holdings []*model.HoldingUTXO) ([]string, error) {
	for _, h := range holdings {
		if h.Amount.Equal(amount) {
			return []string{h.ContractID}, nil
		}
	}

	sorted := make([]*model.HoldingUTXO, len(holdings))
	copy(sorted, holdings)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].Amount.LessThan(sorted[j].Amount)
	})

	if len(sorted) == 0 {
		return nil, fmt.Errorf("sender doesn't have any unlocked holdings")
	}

	largest := sorted[len(sorted)-1]
	sorted = sorted[:len(sorted)-1]

	sum := largest.Amount
	cids := []string{largest.ContractID}
	for _, h := range sorted {
		if sum.GreaterThanOrEqual(amount) {
			break
		}
		sum = sum.Add(h.Amount)
		cids = append(cids, h.ContractID)
	}

	if sum.LessThan(amount) {
		return nil, fmt.Errorf("sender doesn't have sufficient funds for this transfer. missing amount: %s", amount.Sub(sum).String())
	}
	if len(cids) > 100 {
		return nil, fmt.Errorf("exceeded the maximum of 100 utxos in 1 transaction")
	}
	return cids, nil
}

func (t *tokenStandardController) BuildTransferChoiceArgs(
	_ context.Context,
	sender PartyID,
	receiver PartyID,
	amount decimal.Decimal,
	instrumentAdmin PartyID,
	instrumentID string,
	inputUtxos []string,
	memo string,
	expiryDate *time.Time,
	meta types.TEXTMAP,
) (*gen.TransferFactoryTransfer, error) {
	if len(inputUtxos) == 0 {
		return nil, fmt.Errorf("no input utxos provided")
	}

	now := time.Now()
	executeBefore := now.Add(24 * time.Hour)
	if expiryDate != nil {
		executeBefore = *expiryDate
	}

	metaVals := types.TEXTMAP{MemoKey: memo}
	for k, v := range meta {
		metaVals[k] = v
	}

	return &gen.TransferFactoryTransfer{
		ExpectedAdmin: types.PARTY(instrumentAdmin),
		Transfer: gen.Transfer2{
			Sender:           types.PARTY(sender),
			Receiver:         types.PARTY(receiver),
			Amount:           types.NewNumericFromDecimal(amount),
			InstrumentId:     gen.InstrumentId{Admin: types.PARTY(instrumentAdmin), Id: types.TEXT(instrumentID)},
			RequestedAt:      types.TIMESTAMP(now.Add(-requestedAtSkew)),
			ExecuteBefore:    types.TIMESTAMP(executeBefore),
			InputHoldingCids: toContractIDs(inputUtxos),
			Meta:             gen.Metadata{Values: metaVals},
		},
		ExtraArgs: emptyExtraArgs(),
	}, nil
}

func (t *tokenStandardController) CreateTransferFromContext(
	factoryID string,
	choiceArgs *gen.TransferFactoryTransfer,
	choiceContext *ChoiceContext,
) (*model.CommandRequest, error) {
	args := choiceArgs.ToMap()
	args["extraArgs"] = injectContext(choiceContext.ChoiceContextData)

	return newExerciseResult(&damlModel.ExerciseCommand{
		TemplateID: gen.ITransferFactoryInterfaceID(nil),
		ContractID: factoryID,
		Choice:     "TransferFactory_Transfer",
		Arguments:  args,
	}, choiceContext.DisclosedContracts), nil
}

func (t *tokenStandardController) CreateTransferInstruction(ctx context.Context, cid string, choice string) (*model.CommandRequest, error) {
	choices := map[string]string{
		"accept":   "TransferInstruction_Accept",
		"reject":   "TransferInstruction_Reject",
		"withdraw": "TransferInstruction_Withdraw",
	}
	damlChoice, found := choices[strings.ToLower(choice)]
	if !found {
		return nil, fmt.Errorf("wrong choice: %s", choice)
	}

	scanChoiceContext, err := t.scanProxy.GetTransferInstructionChoiceContext(ctx, cid, strings.ToLower(choice))
	if err != nil {
		return nil, fmt.Errorf("fetch transfer instruction choice context: %w", err)
	}

	choiceContext := &ChoiceContext{
		ChoiceContextData:  scanChoiceContext.ChoiceContextData,
		DisclosedContracts: scanChoiceContext.DisclosedContracts,
	}
	return newExerciseResult(&damlModel.ExerciseCommand{
		TemplateID: gen.ITransferInstructionInterfaceID(nil),
		ContractID: cid,
		Choice:     damlChoice,
		Arguments:  map[string]interface{}{"extraArgs": injectContext(choiceContext.ChoiceContextData)},
	}, choiceContext.DisclosedContracts), nil
}

// BuildAllocationFactoryChoiceArgs assembles AllocationFactory_Allocate choice
// arguments, mirroring AllocationService.buildAllocationFactoryChoiceArgs.
func (t *tokenStandardController) BuildAllocationFactoryChoiceArgs(
	ctx context.Context,
	spec gen.AllocationSpecification,
	expectedAdmin PartyID,
	inputUtxos []string,
	requestedAt *time.Time,
) (*gen.AllocationFactoryAllocate, error) {
	if spec.Settlement.Meta.Values == nil {
		spec.Settlement.Meta.Values = types.TEXTMAP{}
	}
	if spec.TransferLeg.Meta.Values == nil {
		spec.TransferLeg.Meta.Values = types.TEXTMAP{}
	}

	cids := inputUtxos
	if len(cids) == 0 {
		selected, err := t.GetInputHoldingsCids(
			ctx,
			PartyID(spec.TransferLeg.Sender),
			string(spec.TransferLeg.InstrumentId.Admin),
			string(spec.TransferLeg.InstrumentId.Id),
			nil,
		)
		if err != nil {
			return nil, err
		}
		cids = selected
	}

	reqAt := time.Now().Add(-requestedAtSkew)
	if requestedAt != nil {
		reqAt = *requestedAt
	}

	return &gen.AllocationFactoryAllocate{
		ExpectedAdmin:    types.PARTY(expectedAdmin),
		Allocation:       spec,
		RequestedAt:      types.TIMESTAMP(reqAt),
		InputHoldingCids: toContractIDs(cids),
		ExtraArgs:        emptyExtraArgs(),
	}, nil
}

// CreateAllocationInstructionFromContext mirrors
// AllocationService.createAllocationInstructionFromContext.
func (t *tokenStandardController) CreateAllocationInstructionFromContext(
	factoryID string,
	choiceArgs *gen.AllocationFactoryAllocate,
	choiceContext *ChoiceContext,
) (*model.CommandRequest, error) {
	args := choiceArgs.ToMap()
	args["extraArgs"] = injectContext(choiceContext.ChoiceContextData)

	return newExerciseResult(&damlModel.ExerciseCommand{
		TemplateID: gen.IAllocationFactoryInterfaceID(nil),
		ContractID: factoryID,
		Choice:     "AllocationFactory_Allocate",
		Arguments:  args,
	}, choiceContext.DisclosedContracts), nil
}

func (t *tokenStandardController) allocationFromContext(allocationCid, choice string, choiceContext *ChoiceContext) *model.CommandRequest {
	return newExerciseResult(&damlModel.ExerciseCommand{
		TemplateID: gen.IAllocationInterfaceID(nil),
		ContractID: allocationCid,
		Choice:     choice,
		Arguments:  map[string]interface{}{"extraArgs": injectContext(choiceContext.ChoiceContextData)},
	}, choiceContext.DisclosedContracts)
}

// CreateExecuteTransferAllocationFromContext mirrors
// AllocationService.createExecuteTransferAllocationFromContext.
func (t *tokenStandardController) CreateExecuteTransferAllocationFromContext(allocationCid string, choiceContext *ChoiceContext) *model.CommandRequest {
	return t.allocationFromContext(allocationCid, "Allocation_ExecuteTransfer", choiceContext)
}

// CreateWithdrawAllocationFromContext mirrors
// AllocationService.createWithdrawAllocationFromContext.
func (t *tokenStandardController) CreateWithdrawAllocationFromContext(allocationCid string, choiceContext *ChoiceContext) *model.CommandRequest {
	return t.allocationFromContext(allocationCid, "Allocation_Withdraw", choiceContext)
}

// CreateCancelAllocationFromContext mirrors
// AllocationService.createCancelAllocationFromContext.
func (t *tokenStandardController) CreateCancelAllocationFromContext(allocationCid string, choiceContext *ChoiceContext) *model.CommandRequest {
	return t.allocationFromContext(allocationCid, "Allocation_Cancel", choiceContext)
}

// CreateWithdrawAllocationInstruction mirrors
// AllocationService.createWithdrawAllocationInstruction (no registry context).
func (t *tokenStandardController) CreateWithdrawAllocationInstruction(allocationInstructionCid string) *model.CommandRequest {
	args := gen.AllocationInstructionWithdraw{ExtraArgs: emptyExtraArgs()}
	return newExerciseResult(&damlModel.ExerciseCommand{
		TemplateID: gen.IAllocationInstructionInterfaceID(nil),
		ContractID: allocationInstructionCid,
		Choice:     "AllocationInstruction_Withdraw",
		Arguments:  args.ToMap(),
	}, nil)
}

// CreateUpdateAllocationInstruction mirrors
// AllocationService.createUpdateAllocationInstruction (no registry context).
func (t *tokenStandardController) CreateUpdateAllocationInstruction(
	allocationInstructionCid string,
	extraActors []PartyID,
	extraArgsContext types.TEXTMAP,
	extraArgsMeta types.TEXTMAP,
) *model.CommandRequest {
	actors := make([]types.PARTY, len(extraActors))
	for i, a := range extraActors {
		actors[i] = types.PARTY(a)
	}
	if extraArgsContext == nil {
		extraArgsContext = types.TEXTMAP{}
	}
	if extraArgsMeta == nil {
		extraArgsMeta = types.TEXTMAP{}
	}

	args := gen.AllocationInstructionUpdate{
		ExtraActors: actors,
		ExtraArgs: gen.ExtraArgs{
			Context: gen.ChoiceContext{Values: extraArgsContext},
			Meta:    gen.Metadata{Values: extraArgsMeta},
		},
	}
	return newExerciseResult(&damlModel.ExerciseCommand{
		TemplateID: gen.IAllocationInstructionInterfaceID(nil),
		ContractID: allocationInstructionCid,
		Choice:     "AllocationInstruction_Update",
		Arguments:  args.ToMap(),
	}, nil)
}

// CreateRejectAllocationRequest mirrors
// AllocationService.createRejectAllocationRequest (no registry context).
func (t *tokenStandardController) CreateRejectAllocationRequest(allocationRequestCid string, actor PartyID) *model.CommandRequest {
	return newExerciseResult(&damlModel.ExerciseCommand{
		TemplateID: ALLOCATION_REQUEST_INTERFACE_ID,
		ContractID: allocationRequestCid,
		Choice:     "AllocationRequest_Reject",
		Arguments: map[string]interface{}{
			"actor":     types.PARTY(actor),
			"extraArgs": emptyExtraArgs().ToMap(),
		},
	}, nil)
}

// CreateWithdrawAllocationRequest mirrors
// AllocationService.createWithdrawAllocationRequest (no registry context).
func (t *tokenStandardController) CreateWithdrawAllocationRequest(allocationRequestCid string) *model.CommandRequest {
	return newExerciseResult(&damlModel.ExerciseCommand{
		TemplateID: ALLOCATION_REQUEST_INTERFACE_ID,
		ContractID: allocationRequestCid,
		Choice:     "AllocationRequest_Withdraw",
		Arguments:  map[string]interface{}{"extraArgs": emptyExtraArgs().ToMap()},
	}, nil)
}

// ---------------------------------------------------------------------------
// Registry composite create* methods (mirror the TS create* wrappers).
// Each accepts an optional prefetched choice context: when supplied the command
// is built locally; otherwise the registry HTTP call it would need is reported
// as unavailable.
// ---------------------------------------------------------------------------

// CreateStandardTransfer mirrors TransferService.createTransfer from
// token-standard-service.ts: it builds the TransferFactory_Transfer choice
// arguments (expiryDate -> executeBefore, memo+meta -> transfer.meta.values) and
// exercises TransferFactory_Transfer against factoryID using the supplied
// registry choice context. Fetching that context from the registry is the
// deprecated HTTP path here, so choiceContext must be provided.
func (t *tokenStandardController) CreateStandardTransfer(
	ctx context.Context,
	sender PartyID,
	receiver PartyID,
	amount decimal.Decimal,
	instrumentAdmin PartyID,
	instrumentID string,
	inputUtxos []string,
	memo string,
	expiryDate *time.Time,
	meta types.TEXTMAP,
	factoryID string,
	choiceContext *ChoiceContext,
) (*model.CommandRequest, error) {
	if choiceContext == nil {
		return nil, fmt.Errorf("choice context is required")
	}

	if len(inputUtxos) == 0 {
		return nil, fmt.Errorf("no input utxos provided")
	}

	choiceArgs, err := t.BuildTransferChoiceArgs(ctx, sender, receiver, amount, instrumentAdmin, instrumentID, inputUtxos, memo, expiryDate, meta)
	if err != nil {
		return nil, err
	}

	return t.CreateTransferFromContext(factoryID, choiceArgs, choiceContext)
}

// CreateStandardAllocationInstruction mirrors AllocationService.createAllocationInstruction.
func (t *tokenStandardController) CreateStandardAllocationInstruction(
	ctx context.Context,
	spec gen.AllocationSpecification,
	expectedAdmin PartyID,
	inputUtxos []string,
	requestedAt *time.Time,
	factoryID string,
	choiceContext *ChoiceContext,
) (*model.CommandRequest, error) {
	choiceArgs, err := t.BuildAllocationFactoryChoiceArgs(ctx, spec, expectedAdmin, inputUtxos, requestedAt)
	if err != nil {
		return nil, err
	}
	if choiceContext != nil {
		return t.CreateAllocationInstructionFromContext(factoryID, choiceArgs, choiceContext)
	}
	return nil, t.registryError("allocation-instruction/v1/allocation-factory")
}

// CreateExecuteTransferAllocation mirrors AllocationService.createExecuteTransferAllocation.
func (t *tokenStandardController) CreateExecuteTransferAllocation(allocationCid string, choiceContext *ChoiceContext) (*model.CommandRequest, error) {
	if choiceContext != nil {
		return t.CreateExecuteTransferAllocationFromContext(allocationCid, choiceContext), nil
	}
	return nil, t.registryError("allocations/v1/{id}/choice-contexts/execute-transfer")
}

// CreateWithdrawAllocation mirrors AllocationService.createWithdrawAllocation.
func (t *tokenStandardController) CreateWithdrawAllocation(allocationCid string, choiceContext *ChoiceContext) (*model.CommandRequest, error) {
	if choiceContext != nil {
		return t.CreateWithdrawAllocationFromContext(allocationCid, choiceContext), nil
	}
	return nil, t.registryError("allocations/v1/{id}/choice-contexts/withdraw")
}

// CreateCancelAllocation mirrors AllocationService.createCancelAllocation.
func (t *tokenStandardController) CreateCancelAllocation(allocationCid string, choiceContext *ChoiceContext) (*model.CommandRequest, error) {
	if choiceContext != nil {
		return t.CreateCancelAllocationFromContext(allocationCid, choiceContext), nil
	}
	return nil, t.registryError("allocations/v1/{id}/choice-contexts/cancel")
}

func (t *tokenStandardController) registryError(path string) error {
	url, _ := t.GetTransferFactoryRegistryUrl()
	t.logger.Warn().Str("registryUrl", url).Str("path", path).Msg("token-standard registry HTTP call not available via ledger client")
	return fmt.Errorf("registry API call (%s) not implemented - requires HTTP client", path)
}

// ---------------------------------------------------------------------------
// Featured-app delegate-proxy choices.
// ---------------------------------------------------------------------------

func beneficiariesToArgs(bs []*Beneficiary) []interface{} {
	out := make([]interface{}, len(bs))
	for i, b := range bs {
		out[i] = map[string]interface{}{
			"beneficiary": types.PARTY(b.Beneficiary),
			"weight":      types.NewNumericFromDecimal(decimal.NewFromFloat(b.Weight)),
		}
	}
	return out
}

func (t *tokenStandardController) validateWeight(bs []*Beneficiary) error {
	var sum float64
	for _, b := range bs {
		sum += b.Weight
	}
	if sum > 1.0 {
		return fmt.Errorf("sum of beneficiary weights is larger than 1")
	}
	return nil
}

func unwrapExercise(res *model.CommandRequest) (*damlModel.ExerciseCommand, error) {
	if res == nil || res.Command == nil {
		return nil, fmt.Errorf("nil command")
	}
	ec, ok := res.Command.Command.(*damlModel.ExerciseCommand)
	if !ok {
		return nil, fmt.Errorf("expected an exercise command")
	}
	return ec, nil
}

func (t *tokenStandardController) wrapDelegateProxy(
	proxyCid string,
	choice string,
	inner *damlModel.ExerciseCommand,
	featuredAppRightCid string,
	beneficiaries []*Beneficiary,
	disclosed []*damlModel.DisclosedContract,
) *model.CommandRequest {
	choiceArgs := map[string]interface{}{
		"cid": types.CONTRACT_ID(inner.ContractID),
		"proxyArg": map[string]interface{}{
			"featuredAppRightCid": types.CONTRACT_ID(featuredAppRightCid),
			"beneficiaries":       beneficiariesToArgs(beneficiaries),
			"choiceArg":           inner.Arguments,
		},
	}
	return newExerciseResult(&damlModel.ExerciseCommand{
		TemplateID: FEATURED_APP_DELEGATE_PROXY_INTERFACE_ID,
		ContractID: proxyCid,
		Choice:     choice,
		Arguments:  choiceArgs,
	}, disclosed)
}

func (t *tokenStandardController) CreateDelegateProxyTransfer(
	ctx context.Context,
	proxyCid string,
	featuredAppRightCid string,
	sender PartyID,
	receiver PartyID,
	amount decimal.Decimal,
	instrumentAdmin PartyID,
	instrumentID string,
	beneficiaries []*Beneficiary,
	inputUtxos []string,
	memo string,
	expiryDate *time.Time,
	meta types.TEXTMAP,
	factoryID string,
	choiceContext *ChoiceContext,
) (*model.CommandRequest, error) {
	if choiceContext == nil {
		return nil, t.registryError("transfer-instruction/v1/transfer-factory")
	}

	if len(inputUtxos) == 0 {
		return nil, fmt.Errorf("no input utxos provided")
	}

	if err := t.validateWeight(beneficiaries); err != nil {
		return nil, err
	}

	choiceArgs, err := t.BuildTransferChoiceArgs(ctx, sender, receiver, amount, instrumentAdmin, instrumentID, inputUtxos, memo, expiryDate, meta)
	if err != nil {
		return nil, err
	}
	inner, err := t.CreateTransferFromContext(factoryID, choiceArgs, choiceContext)
	if err != nil {
		return nil, err
	}
	ec, err := unwrapExercise(inner)
	if err != nil {
		return nil, err
	}
	return t.wrapDelegateProxy(proxyCid, "DelegateProxy_TransferFactory_Transfer", ec, featuredAppRightCid, beneficiaries, inner.DisclosedContracts), nil
}
