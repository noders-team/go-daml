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
	SetScanProxyClient(sp *client.ScanProxyClient)
	ListHoldingUtxos(ctx context.Context, includeLocked bool, limit int) ([]*HoldingUTXO, error)
	GetInputHoldingsCids(
		ctx context.Context,
		sender PartyID,
		instrumentAdmin string,
		instrumentID string,
		amount *decimal.Decimal,
	) ([]string, error)
	GetInputHoldingsCidsForAmount(amount decimal.Decimal, holdings []*HoldingUTXO) ([]string, error)
	CreateTransfer(
		ctx context.Context,
		sender PartyID,
		receiver PartyID,
		amount decimal.Decimal,
		inputUtxos []string,
		expiryDate *time.Time,
	) (*CreateTransferResult, error)
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

func (t *tokenStandardController) SetScanProxyClient(sp *client.ScanProxyClient) {
	t.scanProxy = sp
}

func NewTokenStandardController(userID string, damlClient *client.DamlBindingClient) (TokenStandardController, error) {
	logger := log.Logger.With().
		Str("component", "token-standard-controller").
		Str("userID", userID).
		Logger()

	return &tokenStandardController{
		damlClient: damlClient,
		userID:     userID,
		logger:     logger,
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

func (t *tokenStandardController) Lock(ctx context.Context,
	amount decimal.Decimal, expiresAt time.Time,
) (*LockResponse, error) {
	_, err := t.GetPartyID()
	if err != nil {
		return nil, err
	}

	contractID := fmt.Sprintf("lock-%d", time.Now().UnixNano())

	t.logger.Info().
		Str("amount", amount.String()).
		Time("expiresAt", expiresAt).
		Msg("Lock amulet operation")

	return &LockResponse{
		ContractID: contractID,
		Amount:     amount,
		ExpiresAt:  expiresAt,
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

	filterByParty := map[string]*damlModel.Filters{
		string(partyID): {
			Inclusive: &damlModel.InclusiveFilters{
				TemplateFilters: []*damlModel.TemplateFilter{
					{
						TemplateID:              "3ca1343ab26b453d38c8adb70dca5f1ead8440c42b59b68f070786955cbf9ec1:Splice.Amulet:Amulet",
						IncludeCreatedEventBlob: false, // TODO no hardcoded values
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

	balance := decimal.Zero

	for {
		select {
		case resp, ok := <-stream:
			if !ok {
				t.logger.Debug().
					Str("partyID", string(partyID)).
					Str("balance", balance.String()).
					Msg("balance retrieved")
				return balance, nil
			}
			if entry, ok := resp.ContractEntry.(*damlModel.ActiveContractEntry); ok {
				if entry.ActiveContract != nil && entry.ActiveContract.CreatedEvent != nil {
					contract := entry.ActiveContract.CreatedEvent
					if amountVal, ok := contract.CreateArguments.(map[string]interface{})["amount"]; ok {
						if amountStr, ok := amountVal.(string); ok {
							amount, err := decimal.NewFromString(amountStr)
							if err == nil {
								balance = balance.Add(amount)
							}
						}
					}
				}
			}
		case err := <-errChan:
			if err != nil {
				return decimal.Zero, fmt.Errorf("failed to get balance: %w", err)
			}
		case <-ctx.Done():
			return decimal.Zero, ctx.Err()
		}
	}
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

func numericToDecimal(n types.NUMERIC) decimal.Decimal {
	if n == nil {
		return decimal.Zero
	}
	return decimal.NewFromBigInt((*big.Int)(n), -10)
}

func lockToMap(l *gen.Lock) map[string]interface{} {
	if l == nil {
		return nil
	}
	m := map[string]interface{}{}
	if l.ExpiresAt != nil {
		m["expiresAt"] = time.Time(*l.ExpiresAt)
	}
	return m
}

func (t *tokenStandardController) ListHoldingUtxos(ctx context.Context, includeLocked bool, limit int) ([]*HoldingUTXO, error) {
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
	result := make([]*HoldingUTXO, 0, len(contracts))
	for _, c := range contracts {
		if !includeLocked && IsHoldingLocked(c.Data, now) {
			continue
		}
		result = append(result, &HoldingUTXO{
			ContractID:      c.ContractID,
			Amount:          numericToDecimal(c.Data.Amount),
			InstrumentID:    string(c.Data.InstrumentId.Id),
			InstrumentAdmin: string(c.Data.InstrumentId.Admin),
			Owner:           string(c.Data.Owner),
			Lock:            lockToMap(c.Data.Lock),
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

func (t *tokenStandardController) FetchPendingTransferInstructionView(ctx context.Context) ([]*TransferInstruction, error) {
	contracts, err := t.ListContractsByInterface(ctx, "Splice.TransferInstruction:TransferInstruction")
	if err != nil {
		return nil, err
	}

	var instructions []*TransferInstruction
	for _, contract := range contracts {
		args, ok := contract.CreateArguments.(map[string]interface{})
		if !ok {
			continue
		}

		instruction := &TransferInstruction{
			ContractID:       contract.ContractID,
			CreatedEventBlob: contract.CreatedEventBlob,
		}

		if sender, ok := args["sender"].(string); ok {
			instruction.Sender = sender
		}
		if receiver, ok := args["receiver"].(string); ok {
			instruction.Receiver = receiver
		}
		if amountVal, ok := args["amount"]; ok {
			if amountStr, ok := amountVal.(string); ok {
				instruction.Amount, _ = decimal.NewFromString(amountStr)
			}
		}
		if memo, ok := args["memo"].(string); ok {
			instruction.Memo = memo
		}

		instructions = append(instructions, instruction)
	}

	return instructions, nil
}

type CreateTransferResult struct {
	Command            *damlModel.Command
	DisclosedContracts []*damlModel.DisclosedContract
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
) (*CreateTransferResult, error) {
	if t.scanProxy == nil {
		return nil, fmt.Errorf("scan-proxy client not configured; call SetScanProxyClient")
	}

	if len(inputUtxos) == 0 {
		return nil, fmt.Errorf("no utxos available for transfer")
	}

	amuletRules, err := t.scanProxy.GetAmuletRules(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get AmuletRules from scan-proxy: %w", err)
	}
	round, err := t.scanProxy.GetActiveOpenMiningRound(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get open mining round from scan-proxy: %w", err)
	}
	if round == nil {
		return nil, fmt.Errorf("no active open mining round available")
	}
	dsoParty, err := t.scanProxy.GetDSOPartyID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get DSO party from scan-proxy: %w", err)
	}

	inputs := make([]gen.TransferInput, 0, len(inputUtxos))
	for _, utxo := range inputUtxos {
		cid := types.CONTRACT_ID(utxo)
		inputs = append(inputs, gen.TransferInput{InputAmulet: &cid})
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

	disclosed := make([]*damlModel.DisclosedContract, 0, 2)
	if dc := amuletRules.ToDisclosed(); dc != nil {
		disclosed = append(disclosed, dc)
	}
	if dc := round.ToDisclosed(); dc != nil {
		disclosed = append(disclosed, dc)
	}

	return &CreateTransferResult{
		Command:            &damlModel.Command{Command: exercise},
		DisclosedContracts: disclosed,
	}, nil
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
) (*CreateTransferResult, error) {
	exerciseCmd := &damlModel.Command{
		Command: &damlModel.ExerciseCommand{
			ContractID: transferInstructionCid,
			TemplateID: "Splice.TransferInstruction:TransferInstruction",
			Choice:     choice,
			Arguments:  map[string]interface{}{},
		},
	}

	return &CreateTransferResult{
		Command:            exerciseCmd,
		DisclosedContracts: []*damlModel.DisclosedContract{},
	}, nil
}

func (t *tokenStandardController) ListHoldingTransactions(ctx context.Context, beginExclusive int64, endInclusive *int64) ([]*damlModel.GetUpdatesResponse, error) {
	partyID, err := t.GetPartyID()
	if err != nil {
		return nil, err
	}

	filter := &damlModel.TransactionFilter{
		FiltersByParty: map[string]*damlModel.Filters{
			string(partyID): {
				Inclusive: &damlModel.InclusiveFilters{
					InterfaceFilters: []*damlModel.InterfaceFilter{
						{
							InterfaceID:             "Splice.Holding:Holding", // TODO fix it
							IncludeCreatedEventBlob: true,
						},
					},
				},
			},
		},
	}

	req := &damlModel.GetUpdatesRequest{
		Filter:         filter,
		BeginExclusive: beginExclusive,
		EndInclusive:   endInclusive,
		UpdateFormat:   &damlModel.EventFormat{Verbose: true},
	}

	stream, errChan := t.damlClient.UpdateService.GetUpdates(ctx, req)

	var updates []*damlModel.GetUpdatesResponse

	for {
		select {
		case resp, ok := <-stream:
			if !ok {
				return updates, nil
			}
			updates = append(updates, resp)
		case err := <-errChan:
			if err != nil {
				return nil, fmt.Errorf("failed to list holding transactions: %w", err)
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

// Deprecated: requires the token-standard registry HTTP API
// (GET /registry/metadata/v1/info); not reachable through the ledger client.
func (t *tokenStandardController) GetInstrumentAdmin(ctx context.Context) (string, error) {
	registryUrl, err := t.GetTransferFactoryRegistryUrl()
	if err != nil {
		return "", err
	}

	t.logger.Info().Str("registryUrl", registryUrl).Msg("Getting instrument admin from registry")
	return "", fmt.Errorf("registry API call not implemented - requires HTTP client")
}

// Deprecated: requires the token-standard registry HTTP API
// (GET /registry/metadata/v1/instruments/{instrumentId}); not reachable through the ledger client.
func (t *tokenStandardController) GetInstrumentById(ctx context.Context, instrumentId string) (map[string]interface{}, error) {
	registryUrl, err := t.GetTransferFactoryRegistryUrl()
	if err != nil {
		return nil, err
	}

	t.logger.Info().Str("registryUrl", registryUrl).Str("instrumentId", instrumentId).Msg("Getting instrument by ID")
	return nil, fmt.Errorf("registry API call not implemented - requires HTTP client")
}

// Deprecated: requires the token-standard registry HTTP API
// (GET /registry/metadata/v1/instruments); not reachable through the ledger client.
func (t *tokenStandardController) ListInstruments(ctx context.Context, pageSize int, pageToken string) (map[string]interface{}, error) {
	registryUrl, err := t.GetTransferFactoryRegistryUrl()
	if err != nil {
		return nil, err
	}

	t.logger.Info().Str("registryUrl", registryUrl).Int("pageSize", pageSize).Msg("Listing instruments")
	return nil, fmt.Errorf("registry API call not implemented - requires HTTP client")
}

func (t *tokenStandardController) GetTransactionById(ctx context.Context,
	updateId string,
) (*damlModel.GetUpdatesResponse, error) {
	partyID, err := t.GetPartyID()
	if err != nil {
		return nil, err
	}

	t.logger.Info().Str("partyID", string(partyID)).Str("updateId", updateId).Msg("Getting transaction by ID")
	return nil, fmt.Errorf("getTransactionById not fully implemented")
}

func (t *tokenStandardController) ListHoldingUtxo(ctx context.Context,
	contractId string,
) (*HoldingUTXO, error) {
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
	inputUtxos []*HoldingUTXO,
) (*MergeUtxosResult, error) {
	if partyID == "" {
		var err error
		partyID, err = t.GetPartyID()
		if err != nil {
			return nil, err
		}
	}

	var utxos []*HoldingUTXO
	if len(inputUtxos) > 0 {
		utxos = inputUtxos
	} else {
		var err error
		utxos, err = t.ListHoldingUtxos(ctx, false, nodeLimit)
		if err != nil {
			return nil, err
		}
	}

	utxosByInstrument := make(map[string][]*HoldingUTXO)
	for _, utxo := range utxos {
		key := utxo.InstrumentID + "::" + utxo.InstrumentAdmin
		utxosByInstrument[key] = append(utxosByInstrument[key], utxo)
	}

	var allCommands []*damlModel.Command
	var allDisclosed []*damlModel.DisclosedContract
	transferInputUtxoLimit := 100

	for _, group := range utxosByInstrument {
		if len(group) == 0 {
			continue
		}

		instrumentID := group[0].InstrumentID
		instrumentAdmin := group[0].InstrumentAdmin
		transfers := (len(group) + transferInputUtxoLimit - 1) / transferInputUtxoLimit

		for i := 0; i < transfers; i++ {
			start := i * transferInputUtxoLimit
			end := start + transferInputUtxoLimit
			if end > len(group) {
				end = len(group)
			}

			inputUtxosSlice := group[start:end]
			var totalAmount decimal.Decimal
			var inputCids []string

			for _, utxo := range inputUtxosSlice {
				totalAmount = totalAmount.Add(utxo.Amount)
				inputCids = append(inputCids, utxo.ContractID)
			}

			result, err := t.CreateTransfer(ctx, partyID, partyID, totalAmount, instrumentID, instrumentAdmin, inputCids, "merge-utxos", nil, nil)
			if err != nil {
				return nil, err
			}

			allCommands = append(allCommands, result.Command)
			allDisclosed = append(allDisclosed, result.DisclosedContracts...)
		}
	}

	uniqueDisclosed := make(map[string]*damlModel.DisclosedContract)
	for _, dc := range allDisclosed {
		uniqueDisclosed[dc.ContractID] = dc
	}

	disclosed := make([]*damlModel.DisclosedContract, 0, len(uniqueDisclosed))
	for _, dc := range uniqueDisclosed {
		disclosed = append(disclosed, dc)
	}

	return &MergeUtxosResult{
		Commands:           allCommands,
		DisclosedContracts: disclosed,
	}, nil
}

func (t *tokenStandardController) FetchPendingAllocationInstructionView(ctx context.Context) ([]*AllocationInstruction, error) {
	contracts, err := t.ListContractsByInterface(ctx, "Splice.Allocation:AllocationInstruction")
	if err != nil {
		return nil, err
	}

	var instructions []*AllocationInstruction
	for _, contract := range contracts {
		args, ok := contract.CreateArguments.(map[string]interface{})
		if !ok {
			continue
		}

		instruction := &AllocationInstruction{
			ContractID:       contract.ContractID,
			CreatedEventBlob: contract.CreatedEventBlob,
		}

		if provider, ok := args["provider"].(string); ok {
			instruction.Provider = provider
		}
		if spec, ok := args["specification"].(map[string]interface{}); ok {
			instruction.Specification = spec
		}

		instructions = append(instructions, instruction)
	}

	return instructions, nil
}

func (t *tokenStandardController) FetchPendingAllocationRequestView(ctx context.Context) ([]*AllocationRequest, error) {
	contracts, err := t.ListContractsByInterface(ctx, "Splice.Allocation:AllocationRequest")
	if err != nil {
		return nil, err
	}

	var requests []*AllocationRequest
	for _, contract := range contracts {
		args, ok := contract.CreateArguments.(map[string]interface{})
		if !ok {
			continue
		}

		request := &AllocationRequest{
			ContractID:       contract.ContractID,
			CreatedEventBlob: contract.CreatedEventBlob,
		}

		if requester, ok := args["requester"].(string); ok {
			request.Requester = requester
		}
		if spec, ok := args["specification"].(map[string]interface{}); ok {
			request.Specification = spec
		}

		requests = append(requests, request)
	}

	return requests, nil
}

func (t *tokenStandardController) FetchPendingAllocationView(ctx context.Context) ([]*Allocation, error) {
	contracts, err := t.ListContractsByInterface(ctx, "Splice.Allocation:Allocation")
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
) (*CreateTransferResult, error) {
	cmd := &damlModel.Command{
		Command: &damlModel.ExerciseCommand{
			TemplateID: "Splice.Allocation:AllocationFactory",
			Choice:     "CreateAllocationInstruction",
			Arguments: map[string]interface{}{
				"specification": allocationSpecification,
				"expectedAdmin": expectedAdmin,
				"inputs":        inputUtxos,
				"requestedAt":   requestedAt,
			},
		},
	}

	return &CreateTransferResult{
		Command:            cmd,
		DisclosedContracts: []*damlModel.DisclosedContract{},
	}, nil
}

func (t *tokenStandardController) ExerciseAllocationChoice(
	_ context.Context,
	allocationCid string,
	choice string,
) (*CreateTransferResult, error) {
	cmd := &damlModel.Command{
		Command: &damlModel.ExerciseCommand{
			ContractID: allocationCid,
			TemplateID: "Splice.Allocation:Allocation",
			Choice:     choice,
			Arguments:  map[string]interface{}{},
		},
	}

	return &CreateTransferResult{
		Command:            cmd,
		DisclosedContracts: []*damlModel.DisclosedContract{},
	}, nil
}

func (t *tokenStandardController) ExerciseAllocationInstructionChoice(
	ctx context.Context,
	allocationInstructionCid string,
	choice string,
) (*CreateTransferResult, error) {
	cmd := &damlModel.Command{
		Command: &damlModel.ExerciseCommand{
			ContractID: allocationInstructionCid,
			TemplateID: "Splice.Allocation:AllocationInstruction",
			Choice:     choice,
			Arguments:  map[string]interface{}{},
		},
	}

	return &CreateTransferResult{
		Command:            cmd,
		DisclosedContracts: []*damlModel.DisclosedContract{},
	}, nil
}

func (t *tokenStandardController) ExerciseAllocationRequestChoice(
	ctx context.Context,
	allocationRequestCid string,
	choice string,
	actor PartyID,
) (*CreateTransferResult, error) {
	cmd := &damlModel.Command{
		Command: &damlModel.ExerciseCommand{
			ContractID: allocationRequestCid,
			TemplateID: "Splice.Allocation:AllocationRequest",
			Choice:     choice,
			Arguments: map[string]interface{}{
				"actor": string(actor),
			},
		},
	}

	return &CreateTransferResult{
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
) (*CreateTransferResult, error) {
	cmd := &damlModel.Command{
		Command: &damlModel.ExerciseCommand{
			ContractID: proxyCid,
			TemplateID: "Splice.DelegateProxy:DelegateProxy",
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

	return &CreateTransferResult{
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
) (*CreateTransferResult, error) {
	cmd := &damlModel.Command{
		Command: &damlModel.ExerciseCommand{
			ContractID: proxyCid,
			TemplateID: "Splice.DelegateProxy:DelegateProxy",
			Choice:     "ExerciseTransferInstructionChoice",
			Arguments: map[string]interface{}{
				"transferInstructionCid": transferInstructionCid,
				"choice":                 choice,
				"featuredAppRightCid":    featuredAppRightCid,
				"beneficiaries":          beneficiaries,
			},
		},
	}

	return &CreateTransferResult{
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
) (*CreateTransferResult, error) {
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

	return &CreateTransferResult{
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
) (*CreateTransferResult, error) {
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

	return &CreateTransferResult{
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
) (*CreateTransferResult, error) {
	syncID, err := t.GetSynchronizerID()
	if err != nil {
		return nil, err
	}

	cmd := &damlModel.Command{
		Command: &damlModel.ExerciseCommand{
			TemplateID: "Splice.AmuletRules:AmuletRules",
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

	return &CreateTransferResult{
		Command:            cmd,
		DisclosedContracts: []*damlModel.DisclosedContract{},
	}, nil
}

// Deprecated: requires the scan-proxy HTTP API; not reachable through the ledger client.
func (t *tokenStandardController) GetMemberTrafficStatus(ctx context.Context, memberId string) (map[string]interface{}, error) {
	syncID, err := t.GetSynchronizerID()
	if err != nil {
		return nil, err
	}

	t.logger.Info().Str("synchronizerId", string(syncID)).Str("memberId", memberId).Msg("Getting member traffic status")
	return nil, fmt.Errorf("scan proxy API call not implemented - requires HTTP client")
}

type FeaturedAppRight struct {
	TemplateID       string
	ContractID       string
	Payload          map[string]interface{}
	CreatedEventBlob []byte
	CreatedAt        string
}

func (t *tokenStandardController) SelfGrantFeatureAppRights(ctx context.Context) (*CreateTransferResult, error) {
	partyID, err := t.GetPartyID()
	if err != nil {
		return nil, err
	}

	syncID, err := t.GetSynchronizerID()
	if err != nil {
		return nil, err
	}

	cmd := &damlModel.Command{
		Command: &damlModel.ExerciseCommand{
			TemplateID: "Splice.AmuletRules:AmuletRules",
			Choice:     "AmuletRules_DevNet_FeatureApp",
			Arguments: map[string]interface{}{
				"provider":       string(partyID),
				"synchronizerId": string(syncID),
			},
		},
	}

	return &CreateTransferResult{
		Command:            cmd,
		DisclosedContracts: []*damlModel.DisclosedContract{},
	}, nil
}

func (t *tokenStandardController) LookupFeaturedApps(ctx context.Context, maxRetries int, delayMs int) (*FeaturedAppRight, error) {
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

	for attempt := 1; attempt <= maxRetries; attempt++ {
		contracts, err := t.ListContractsByInterface(ctx, "Splice.Amulet:FeaturedAppRight")
		if err == nil && len(contracts) > 0 {
			for _, contract := range contracts {
				args, ok := contract.CreateArguments.(map[string]interface{})
				if !ok {
					continue
				}

				if provider, ok := args["provider"].(string); ok && provider == string(partyID) {
					return &FeaturedAppRight{
						TemplateID:       contract.TemplateID,
						ContractID:       contract.ContractID,
						Payload:          args,
						CreatedEventBlob: contract.CreatedEventBlob,
					}, nil
				}
			}
		}

		t.logger.Info().Int("attempt", attempt).Msg("Lookup featured apps returned undefined, retrying again...")

		if attempt < maxRetries {
			time.Sleep(time.Duration(delayMs) * time.Millisecond)
		}
	}

	return nil, nil
}

func (t *tokenStandardController) GrantFeatureAppRightsForInternalParty(ctx context.Context) (*FeaturedAppRight, error) {
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

func (t *tokenStandardController) UseMergeDelegations(ctx context.Context, walletParty PartyID, nodeLimit int) (*CreateTransferResult, error) {
	if nodeLimit == 0 {
		nodeLimit = 200
	}

	utxos, err := t.ListHoldingUtxos(ctx, true, 100)
	if err != nil {
		return nil, err
	}

	if len(utxos) < 10 {
		return nil, fmt.Errorf("utxos are less than 10, found %d", len(utxos))
	}

	req := &damlModel.GetActiveContractsRequest{
		EventFormat: &damlModel.EventFormat{
			Verbose: true,
			FiltersByParty: map[string]*damlModel.Filters{
				string(walletParty): {
					Inclusive: &damlModel.InclusiveFilters{
						TemplateFilters: []*damlModel.TemplateFilter{
							{TemplateID: "Splice.MergeDelegation:MergeDelegation"},
						},
					},
				},
			},
		},
	}

	stream, errChan := t.damlClient.StateService.GetActiveContracts(ctx, req)

	var mergeDelegationCid string
	select {
	case resp := <-stream:
		if entry, ok := resp.ContractEntry.(*damlModel.ActiveContractEntry); ok {
			if entry.ActiveContract != nil && entry.ActiveContract.CreatedEvent != nil {
				mergeDelegationCid = entry.ActiveContract.CreatedEvent.ContractID
			}
		}
	case err := <-errChan:
		return nil, err
	}

	if mergeDelegationCid == "" {
		return nil, fmt.Errorf("merge delegation contract not found")
	}

	mergeResult, err := t.MergeHoldingUtxos(ctx, nodeLimit, walletParty, utxos)
	if err != nil {
		return nil, err
	}

	var mergeCalls []map[string]interface{}
	for _, cmd := range mergeResult.Commands {
		mergeCalls = append(mergeCalls, map[string]interface{}{
			"delegationCid": mergeDelegationCid,
			"choiceArg":     cmd,
		})
	}

	batchCmd := &damlModel.Command{
		Command: &damlModel.ExerciseCommand{
			TemplateID: "Splice.BatchMergeUtility:BatchMergeUtility",
			Choice:     "BatchMergeUtility_BatchMerge",
			Arguments: map[string]interface{}{
				"mergeCalls": mergeCalls,
			},
		},
	}

	return &CreateTransferResult{
		Command:            batchCmd,
		DisclosedContracts: mergeResult.DisclosedContracts,
	}, nil
}

func (t *tokenStandardController) CreateBatchMergeUtility(ctx context.Context) (*damlModel.Command, error) {
	partyID, err := t.GetPartyID()
	if err != nil {
		return nil, err
	}

	return &damlModel.Command{
		Command: &damlModel.CreateCommand{
			TemplateID: "Splice.BatchMergeUtility:BatchMergeUtility",
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
			TemplateID: "Splice.MergeDelegationProposal:MergeDelegationProposal",
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
							{TemplateID: "Splice.MergeDelegationProposal:MergeDelegationProposal"},
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

func (t *tokenStandardController) ApproveMergeDelegationProposal(ctx context.Context, ownerParty PartyID) (*CreateTransferResult, error) {
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
			TemplateID: "Splice.MergeDelegationProposal:MergeDelegationProposal",
			Choice:     "MergeDelegationProposal_Accept",
			Arguments:  map[string]interface{}{},
		},
	}

	disclosed := &damlModel.DisclosedContract{
		TemplateID:       proposal.TemplateID,
		ContractID:       proposal.ContractID,
		CreatedEventBlob: proposal.CreatedEventBlob,
	}

	return &CreateTransferResult{
		Command:            cmd,
		DisclosedContracts: []*damlModel.DisclosedContract{disclosed},
	}, nil
}

type PartyID string

type HoldingUTXO struct {
	ContractID       string
	Amount           decimal.Decimal
	InstrumentID     string
	InstrumentAdmin  string
	Owner            string
	Lock             map[string]interface{}
	CreatedEventBlob []byte
}

type TransferInstruction struct {
	ContractID       string
	Sender           string
	Receiver         string
	Amount           decimal.Decimal
	Memo             string
	CreatedEventBlob []byte
}

type MergeUtxosResult struct {
	Commands           []*damlModel.Command
	DisclosedContracts []*damlModel.DisclosedContract
}

type AllocationInstruction struct {
	ContractID       string
	Provider         string
	Specification    map[string]interface{}
	CreatedEventBlob []byte
}

type AllocationRequest struct {
	ContractID       string
	Requester        string
	Specification    map[string]interface{}
	CreatedEventBlob []byte
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

type LockResponse struct {
	ContractID string
	Amount     decimal.Decimal
	ExpiresAt  time.Time
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

func newExerciseResult(cmd *damlModel.ExerciseCommand, disclosed []*damlModel.DisclosedContract) *CreateTransferResult {
	if disclosed == nil {
		disclosed = []*damlModel.DisclosedContract{}
	}
	return &CreateTransferResult{
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
		"context": ctxData,
		"meta":    gen.Metadata{Values: types.TEXTMAP{}}.ToMap(),
	}
}

// IsHoldingLocked reports whether a holding is currently locked, mirroring
// TokenStandardService.isHoldingLocked.
func IsHoldingLocked(view gen.HoldingView, now time.Time) bool {
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
		filtered := make([]*HoldingUTXO, 0)
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

func (t *tokenStandardController) GetInputHoldingsCidsForAmount(amount decimal.Decimal, holdings []*HoldingUTXO) ([]string, error) {
	for _, h := range holdings {
		if h.Amount.Equal(amount) {
			return []string{h.ContractID}, nil
		}
	}

	sorted := make([]*HoldingUTXO, len(holdings))
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
) (*gen.TransferFactoryTransfer, error) {
	cids := inputUtxos
	if len(cids) == 0 {
		amt := amount
		selected, err := t.GetInputHoldingsCids(ctx, sender, string(instrumentAdmin), instrumentID, &amt)
		if err != nil {
			return nil, err
		}
		cids = selected
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
			InputHoldingCids: toContractIDs(cids),
			Meta:             gen.Metadata{Values: metaVals},
		},
		ExtraArgs: emptyExtraArgs(),
	}, nil
}

// CreateTransferFromContext builds the TransferFactory_Transfer exercise from a
// prefetched registry choice context, mirroring TransferService.createTransferFromContext.
func (t *tokenStandardController) CreateTransferFromContext(
	factoryID string,
	choiceArgs *gen.TransferFactoryTransfer,
	choiceContext *ChoiceContext,
) (*CreateTransferResult, error) {
	args := choiceArgs.ToMap()
	args["extraArgs"] = injectContext(choiceContext.ChoiceContextData)

	return newExerciseResult(&damlModel.ExerciseCommand{
		TemplateID: gen.ITransferFactoryInterfaceID(nil),
		ContractID: factoryID,
		Choice:     "TransferFactory_Transfer",
		Arguments:  args,
	}, choiceContext.DisclosedContracts), nil
}

func (t *tokenStandardController) transferInstructionFromContext(
	transferInstructionCid string,
	choice string,
	choiceContext *ChoiceContext,
) *CreateTransferResult {
	return newExerciseResult(&damlModel.ExerciseCommand{
		TemplateID: gen.ITransferInstructionInterfaceID(nil),
		ContractID: transferInstructionCid,
		Choice:     choice,
		Arguments:  map[string]interface{}{"extraArgs": injectContext(choiceContext.ChoiceContextData)},
	}, choiceContext.DisclosedContracts)
}

// CreateAcceptTransferInstructionFromContext mirrors
// TransferService.createAcceptTransferInstructionFromContext.
func (t *tokenStandardController) CreateAcceptTransferInstructionFromContext(cid string, choiceContext *ChoiceContext) *CreateTransferResult {
	return t.transferInstructionFromContext(cid, "TransferInstruction_Accept", choiceContext)
}

// CreateRejectTransferInstructionFromContext mirrors
// TransferService.createRejectTransferInstructionFromContext.
func (t *tokenStandardController) CreateRejectTransferInstructionFromContext(cid string, choiceContext *ChoiceContext) *CreateTransferResult {
	return t.transferInstructionFromContext(cid, "TransferInstruction_Reject", choiceContext)
}

// CreateWithdrawTransferInstructionFromContext mirrors
// TransferService.createWithdrawTransferInstructionFromContext.
func (t *tokenStandardController) CreateWithdrawTransferInstructionFromContext(cid string, choiceContext *ChoiceContext) *CreateTransferResult {
	return t.transferInstructionFromContext(cid, "TransferInstruction_Withdraw", choiceContext)
}

// CreateTransferInstruction dispatches to the accept/reject/withdraw builder for
// a prefetched context, mirroring TransferService.createTransferInstruction.
func (t *tokenStandardController) CreateTransferInstruction(cid string, choice string, choiceContext *ChoiceContext) (*CreateTransferResult, error) {
	switch choice {
	case "Accept":
		return t.CreateAcceptTransferInstructionFromContext(cid, choiceContext), nil
	case "Reject":
		return t.CreateRejectTransferInstructionFromContext(cid, choiceContext), nil
	case "Withdraw":
		return t.CreateWithdrawTransferInstructionFromContext(cid, choiceContext), nil
	default:
		return nil, fmt.Errorf("unknown transfer instruction choice %q", choice)
	}
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
) (*CreateTransferResult, error) {
	args := choiceArgs.ToMap()
	args["extraArgs"] = injectContext(choiceContext.ChoiceContextData)

	return newExerciseResult(&damlModel.ExerciseCommand{
		TemplateID: gen.IAllocationFactoryInterfaceID(nil),
		ContractID: factoryID,
		Choice:     "AllocationFactory_Allocate",
		Arguments:  args,
	}, choiceContext.DisclosedContracts), nil
}

func (t *tokenStandardController) allocationFromContext(allocationCid, choice string, choiceContext *ChoiceContext) *CreateTransferResult {
	return newExerciseResult(&damlModel.ExerciseCommand{
		TemplateID: gen.IAllocationInterfaceID(nil),
		ContractID: allocationCid,
		Choice:     choice,
		Arguments:  map[string]interface{}{"extraArgs": injectContext(choiceContext.ChoiceContextData)},
	}, choiceContext.DisclosedContracts)
}

// CreateExecuteTransferAllocationFromContext mirrors
// AllocationService.createExecuteTransferAllocationFromContext.
func (t *tokenStandardController) CreateExecuteTransferAllocationFromContext(allocationCid string, choiceContext *ChoiceContext) *CreateTransferResult {
	return t.allocationFromContext(allocationCid, "Allocation_ExecuteTransfer", choiceContext)
}

// CreateWithdrawAllocationFromContext mirrors
// AllocationService.createWithdrawAllocationFromContext.
func (t *tokenStandardController) CreateWithdrawAllocationFromContext(allocationCid string, choiceContext *ChoiceContext) *CreateTransferResult {
	return t.allocationFromContext(allocationCid, "Allocation_Withdraw", choiceContext)
}

// CreateCancelAllocationFromContext mirrors
// AllocationService.createCancelAllocationFromContext.
func (t *tokenStandardController) CreateCancelAllocationFromContext(allocationCid string, choiceContext *ChoiceContext) *CreateTransferResult {
	return t.allocationFromContext(allocationCid, "Allocation_Cancel", choiceContext)
}

// CreateWithdrawAllocationInstruction mirrors
// AllocationService.createWithdrawAllocationInstruction (no registry context).
func (t *tokenStandardController) CreateWithdrawAllocationInstruction(allocationInstructionCid string) *CreateTransferResult {
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
) *CreateTransferResult {
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
func (t *tokenStandardController) CreateRejectAllocationRequest(allocationRequestCid string, actor PartyID) *CreateTransferResult {
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
func (t *tokenStandardController) CreateWithdrawAllocationRequest(allocationRequestCid string) *CreateTransferResult {
	return newExerciseResult(&damlModel.ExerciseCommand{
		TemplateID: ALLOCATION_REQUEST_INTERFACE_ID,
		ContractID: allocationRequestCid,
		Choice:     "AllocationRequest_Withdraw",
		Arguments:  map[string]interface{}{"extraArgs": emptyExtraArgs().ToMap()},
	}, nil)
}

// Deprecated: requires the token-standard registry HTTP API
// (POST /registry/transfer-instruction/v1/transfer-factory); use
// CreateTransferFromContext with a context obtained out of band.
func (t *tokenStandardController) FetchTransferFactoryChoiceContext(ctx context.Context, choiceArgs *gen.TransferFactoryTransfer) (*ChoiceContext, error) {
	return t.registryUnavailable("transfer-instruction/v1/transfer-factory")
}

// Deprecated: requires the token-standard registry HTTP API
// (POST /registry/transfer-instruction/v1/{id}/choice-contexts/accept).
func (t *tokenStandardController) FetchAcceptTransferInstructionChoiceContext(ctx context.Context, transferInstructionCid string) (*ChoiceContext, error) {
	return t.registryUnavailable("transfer-instruction/v1/{id}/choice-contexts/accept")
}

// Deprecated: requires the token-standard registry HTTP API
// (POST /registry/transfer-instruction/v1/{id}/choice-contexts/reject).
func (t *tokenStandardController) FetchRejectTransferInstructionChoiceContext(ctx context.Context, transferInstructionCid string) (*ChoiceContext, error) {
	return t.registryUnavailable("transfer-instruction/v1/{id}/choice-contexts/reject")
}

// Deprecated: requires the token-standard registry HTTP API
// (POST /registry/transfer-instruction/v1/{id}/choice-contexts/withdraw).
func (t *tokenStandardController) FetchWithdrawTransferInstructionChoiceContext(ctx context.Context, transferInstructionCid string) (*ChoiceContext, error) {
	return t.registryUnavailable("transfer-instruction/v1/{id}/choice-contexts/withdraw")
}

// Deprecated: requires the token-standard registry HTTP API
// (POST /registry/allocation-instruction/v1/allocation-factory).
func (t *tokenStandardController) FetchAllocationFactoryChoiceContext(ctx context.Context, choiceArgs *gen.AllocationFactoryAllocate) (*ChoiceContext, error) {
	return t.registryUnavailable("allocation-instruction/v1/allocation-factory")
}

// Deprecated: requires the token-standard registry HTTP API
// (POST /registry/allocations/v1/{id}/choice-contexts/execute-transfer).
func (t *tokenStandardController) FetchExecuteTransferChoiceContext(ctx context.Context, allocationCid string) (*ChoiceContext, error) {
	return t.registryUnavailable("allocations/v1/{id}/choice-contexts/execute-transfer")
}

// Deprecated: requires the token-standard registry HTTP API
// (POST /registry/allocations/v1/{id}/choice-contexts/withdraw).
func (t *tokenStandardController) FetchWithdrawAllocationChoiceContext(ctx context.Context, allocationCid string) (*ChoiceContext, error) {
	return t.registryUnavailable("allocations/v1/{id}/choice-contexts/withdraw")
}

// Deprecated: requires the token-standard registry HTTP API
// (POST /registry/allocations/v1/{id}/choice-contexts/cancel).
func (t *tokenStandardController) FetchCancelAllocationChoiceContext(ctx context.Context, allocationCid string) (*ChoiceContext, error) {
	return t.registryUnavailable("allocations/v1/{id}/choice-contexts/cancel")
}

// Deprecated: requires the token-standard registry HTTP API; combines
// listInstruments + getInstrumentAdmin which are not reachable via the ledger client.
func (t *tokenStandardController) InstrumentsToAsset(ctx context.Context) ([]map[string]interface{}, error) {
	return nil, t.registryError("metadata/v1/instruments")
}

// Deprecated: requires the token-standard registry HTTP API; iterates
// InstrumentsToAsset across registries.
func (t *tokenStandardController) RegistriesToAssets(ctx context.Context, registryUrls []string) ([]map[string]interface{}, error) {
	return nil, t.registryError("metadata/v1/instruments")
}

func (t *tokenStandardController) registryUnavailable(path string) (*ChoiceContext, error) {
	return nil, t.registryError(path)
}

func (t *tokenStandardController) registryError(path string) error {
	url, _ := t.GetTransferFactoryRegistryUrl()
	t.logger.Warn().Str("registryUrl", url).Str("path", path).Msg("token-standard registry HTTP call not available via ledger client")
	return fmt.Errorf("registry API call (%s) not implemented - requires HTTP client", path)
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
) (*CreateTransferResult, error) {
	if choiceContext == nil {
		return nil, fmt.Errorf("choice context is required")
	}

	choiceArgs, err := t.BuildTransferChoiceArgs(ctx, sender, receiver, amount, instrumentAdmin, instrumentID, inputUtxos, memo, expiryDate, meta)
	if err != nil {
		return nil, err
	}

	return t.CreateTransferFromContext(factoryID, choiceArgs, choiceContext)
}

// CreateAcceptTransferInstruction mirrors TransferService.createAcceptTransferInstruction.
func (t *tokenStandardController) CreateAcceptTransferInstruction(cid string, choiceContext *ChoiceContext) (*CreateTransferResult, error) {
	if choiceContext != nil {
		return t.CreateAcceptTransferInstructionFromContext(cid, choiceContext), nil
	}
	return nil, t.registryError("transfer-instruction/v1/{id}/choice-contexts/accept")
}

// CreateRejectTransferInstruction mirrors TransferService.createRejectTransferInstruction.
func (t *tokenStandardController) CreateRejectTransferInstruction(cid string, choiceContext *ChoiceContext) (*CreateTransferResult, error) {
	if choiceContext != nil {
		return t.CreateRejectTransferInstructionFromContext(cid, choiceContext), nil
	}
	return nil, t.registryError("transfer-instruction/v1/{id}/choice-contexts/reject")
}

// CreateWithdrawTransferInstruction mirrors TransferService.createWithdrawTransferInstruction.
func (t *tokenStandardController) CreateWithdrawTransferInstruction(cid string, choiceContext *ChoiceContext) (*CreateTransferResult, error) {
	if choiceContext != nil {
		return t.CreateWithdrawTransferInstructionFromContext(cid, choiceContext), nil
	}
	return nil, t.registryError("transfer-instruction/v1/{id}/choice-contexts/withdraw")
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
) (*CreateTransferResult, error) {
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
func (t *tokenStandardController) CreateExecuteTransferAllocation(allocationCid string, choiceContext *ChoiceContext) (*CreateTransferResult, error) {
	if choiceContext != nil {
		return t.CreateExecuteTransferAllocationFromContext(allocationCid, choiceContext), nil
	}
	return nil, t.registryError("allocations/v1/{id}/choice-contexts/execute-transfer")
}

// CreateWithdrawAllocation mirrors AllocationService.createWithdrawAllocation.
func (t *tokenStandardController) CreateWithdrawAllocation(allocationCid string, choiceContext *ChoiceContext) (*CreateTransferResult, error) {
	if choiceContext != nil {
		return t.CreateWithdrawAllocationFromContext(allocationCid, choiceContext), nil
	}
	return nil, t.registryError("allocations/v1/{id}/choice-contexts/withdraw")
}

// CreateCancelAllocation mirrors AllocationService.createCancelAllocation.
func (t *tokenStandardController) CreateCancelAllocation(allocationCid string, choiceContext *ChoiceContext) (*CreateTransferResult, error) {
	if choiceContext != nil {
		return t.CreateCancelAllocationFromContext(allocationCid, choiceContext), nil
	}
	return nil, t.registryError("allocations/v1/{id}/choice-contexts/cancel")
}

// ---------------------------------------------------------------------------
// Featured-app delegate-proxy choices.
// ---------------------------------------------------------------------------

func beneficiariesToArgs(bs []Beneficiary) []interface{} {
	out := make([]interface{}, len(bs))
	for i, b := range bs {
		out[i] = map[string]interface{}{
			"beneficiary": types.PARTY(b.Beneficiary),
			"weight":      types.NewNumericFromDecimal(decimal.NewFromFloat(b.Weight)),
		}
	}
	return out
}

func validateBeneficiaryWeights(bs []Beneficiary) error {
	var sum float64
	for _, b := range bs {
		sum += b.Weight
	}
	if sum > 1.0 {
		return fmt.Errorf("sum of beneficiary weights is larger than 1")
	}
	return nil
}

func unwrapExercise(res *CreateTransferResult) (*damlModel.ExerciseCommand, error) {
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
	beneficiaries []Beneficiary,
	disclosed []*damlModel.DisclosedContract,
) *CreateTransferResult {
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

// CreateDelegateProxyTransfer mirrors TokenStandardService.createDelegateProxyTransfer.
// The inner TransferFactory_Transfer requires a registry choice context, so it
// must be supplied via choiceContext (see the Deprecated Fetch* methods).
func (t *tokenStandardController) CreateDelegateProxyTransfer(
	ctx context.Context,
	proxyCid string,
	featuredAppRightCid string,
	sender PartyID,
	receiver PartyID,
	amount decimal.Decimal,
	instrumentAdmin PartyID,
	instrumentID string,
	beneficiaries []Beneficiary,
	inputUtxos []string,
	memo string,
	expiryDate *time.Time,
	meta types.TEXTMAP,
	factoryID string,
	choiceContext *ChoiceContext,
) (*CreateTransferResult, error) {
	if choiceContext == nil {
		return nil, t.registryError("transfer-instruction/v1/transfer-factory")
	}
	if err := validateBeneficiaryWeights(beneficiaries); err != nil {
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

// ExerciseDelegateProxyTransferInstructionAccept mirrors
// TransferService.exerciseDelegateProxyTransferInstructionAccept.
func (t *tokenStandardController) ExerciseDelegateProxyTransferInstructionAccept(
	proxyCid string,
	transferInstructionCid string,
	featuredAppRightCid string,
	beneficiaries []Beneficiary,
	choiceContext *ChoiceContext,
) (*CreateTransferResult, error) {
	if err := validateBeneficiaryWeights(beneficiaries); err != nil {
		return nil, err
	}
	inner, err := t.CreateAcceptTransferInstruction(transferInstructionCid, choiceContext)
	if err != nil {
		return nil, err
	}
	ec, err := unwrapExercise(inner)
	if err != nil {
		return nil, err
	}
	return t.wrapDelegateProxy(proxyCid, "DelegateProxy_TransferInstruction_Accept", ec, featuredAppRightCid, beneficiaries, inner.DisclosedContracts), nil
}

// ExerciseDelegateProxyTransferInstructionReject mirrors
// TransferService.exerciseDelegateProxyTransferInstructionReject.
func (t *tokenStandardController) ExerciseDelegateProxyTransferInstructionReject(
	proxyCid string,
	transferInstructionCid string,
	featuredAppRightCid string,
	beneficiaries []Beneficiary,
	choiceContext *ChoiceContext,
) (*CreateTransferResult, error) {
	if err := validateBeneficiaryWeights(beneficiaries); err != nil {
		return nil, err
	}
	inner, err := t.CreateRejectTransferInstruction(transferInstructionCid, choiceContext)
	if err != nil {
		return nil, err
	}
	ec, err := unwrapExercise(inner)
	if err != nil {
		return nil, err
	}
	return t.wrapDelegateProxy(proxyCid, "DelegateProxy_TransferInstruction_Reject", ec, featuredAppRightCid, beneficiaries, inner.DisclosedContracts), nil
}

// ExerciseDelegateProxyTransferInstructionWithdraw mirrors
// TransferService.exerciseDelegateProxyTransferInstructioWithdraw.
func (t *tokenStandardController) ExerciseDelegateProxyTransferInstructionWithdraw(
	proxyCid string,
	transferInstructionCid string,
	featuredAppRightCid string,
	beneficiaries []Beneficiary,
	choiceContext *ChoiceContext,
) (*CreateTransferResult, error) {
	if err := validateBeneficiaryWeights(beneficiaries); err != nil {
		return nil, err
	}
	inner, err := t.CreateWithdrawTransferInstruction(transferInstructionCid, choiceContext)
	if err != nil {
		return nil, err
	}
	ec, err := unwrapExercise(inner)
	if err != nil {
		return nil, err
	}
	return t.wrapDelegateProxy(proxyCid, "DelegateProxy_TransferInstruction_Withdraw", ec, featuredAppRightCid, beneficiaries, inner.DisclosedContracts), nil
}

// ExerciseDelegateProxyTransferInstructionAcceptForExchange mirrors
// TokenStandardService.exerciseDelegateProxyTransferInstructionAccept, which
// accepts on behalf of a single exchange party (weight 1.0).
func (t *tokenStandardController) ExerciseDelegateProxyTransferInstructionAcceptForExchange(
	exchangeParty PartyID,
	proxyCid string,
	transferInstructionCid string,
	featuredAppRightCid string,
	choiceContext *ChoiceContext,
) (*CreateTransferResult, error) {
	return t.ExerciseDelegateProxyTransferInstructionAccept(
		proxyCid,
		transferInstructionCid,
		featuredAppRightCid,
		[]Beneficiary{{Beneficiary: exchangeParty, Weight: 1.0}},
		choiceContext,
	)
}

// ---------------------------------------------------------------------------
// Pure helpers.
// ---------------------------------------------------------------------------

// FilterHoldingsByInstrument mirrors CoreService.filterHoldingsByInstrument.
func FilterHoldingsByInstrument(holdings []*PrettyContract[gen.HoldingView], instrumentAdmin, instrumentID string) []*PrettyContract[gen.HoldingView] {
	out := make([]*PrettyContract[gen.HoldingView], 0, len(holdings))
	for _, h := range holdings {
		if string(h.InterfaceView.InstrumentId.Id) == instrumentID &&
			string(h.InterfaceView.InstrumentId.Admin) == instrumentAdmin {
			out = append(out, h)
		}
	}
	return out
}

// ToQualifiedMemberId mirrors CoreService.toQualifiedMemberId.
func ToQualifiedMemberId(memberID string) (string, error) {
	if memberID == "" {
		return "", fmt.Errorf("memberId is required")
	}
	if strings.HasPrefix(memberID, "PAR::") || strings.HasPrefix(memberID, "MED::") {
		return memberID, nil
	}
	return "PAR::" + memberID, nil
}

// ---------------------------------------------------------------------------
// Transaction rendering (registry/tx-parser dependent).
// ---------------------------------------------------------------------------

// Deprecated: requires the token-standard registry client (TokenStandardClient);
// not reachable through the ledger client.
func (t *tokenStandardController) GetTokenStandardClient(registryURL string) error {
	return t.registryError("token-standard-client")
}

// Deprecated: requires the core-tx-parser transaction renderer, which has no Go
// port; use ListHoldingTransactions for the raw ledger updates.
func (t *tokenStandardController) ToPrettyTransactions(ctx context.Context, partyID PartyID) ([]map[string]interface{}, error) {
	return nil, fmt.Errorf("pretty transaction rendering requires the core-tx-parser, which has no Go port")
}

// Deprecated: requires the core-tx-parser transaction renderer, which has no Go port.
func (t *tokenStandardController) ToPrettyTransaction(ctx context.Context, updateID string, partyID PartyID) (map[string]interface{}, error) {
	return nil, fmt.Errorf("pretty transaction rendering requires the core-tx-parser, which has no Go port")
}

// Deprecated: requires the core-tx-parser transaction renderer, which has no Go port.
func (t *tokenStandardController) ToPrettyTransferObjects(ctx context.Context, updateID string, partyID PartyID) ([]map[string]interface{}, error) {
	return nil, fmt.Errorf("pretty transfer objects require the core-tx-parser, which has no Go port")
}

// Deprecated: requires the core-tx-parser transaction renderer, which has no Go port.
func (t *tokenStandardController) ToPrettyTransactionsPerParty(ctx context.Context, parties []PartyID) (map[PartyID][]map[string]interface{}, error) {
	return nil, fmt.Errorf("pretty transaction rendering requires the core-tx-parser, which has no Go port")
}

// Deprecated: requires the core-tx-parser transaction renderer, which has no Go
// port; use GetTransactionById for the raw ledger transaction.
func (t *tokenStandardController) GetTransferObjectsById(ctx context.Context, updateID string, partyID PartyID) ([]map[string]interface{}, error) {
	return nil, fmt.Errorf("transfer objects require the core-tx-parser, which has no Go port")
}
