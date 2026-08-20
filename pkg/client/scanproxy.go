package client

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/noders-team/go-daml/pkg/auth"
	"github.com/noders-team/go-daml/pkg/model"
	"github.com/noders-team/go-daml/pkg/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

// ScanProxyClient is an HTTP client for the Splice scan-proxy API. It exposes
// DSO-owned contracts (AmuletRules, open mining rounds, transfer pre-approvals,
// featured apps, ANS entries) that are not directly readable from a wallet
// party's ledger view.
type ScanProxyClient struct {
	baseURL       string
	httpClient    *http.Client
	tokenProvider auth.TokenProvider
	logger        zerolog.Logger
}

type ScanContract struct {
	ContractID       string                 `json:"contract_id"`
	TemplateID       string                 `json:"template_id"`
	Payload          map[string]interface{} `json:"payload"`
	CreatedEventBlob string                 `json:"created_event_blob,omitempty"`
	CreatedAt        string                 `json:"created_at,omitempty"`
	SynchronizerID   string                 `json:"-"`
}

// ToDisclosed converts the scan contract into a ledger DisclosedContract,
// decoding the base64 created_event_blob. Returns nil if the contract carries no
// blob to disclose.
func (sc *ScanContract) ToDisclosed() *model.DisclosedContract {
	if sc == nil || sc.ContractID == "" || sc.CreatedEventBlob == "" {
		return nil
	}
	blob, err := base64.StdEncoding.DecodeString(sc.CreatedEventBlob)
	if err != nil {
		return nil
	}
	return &model.DisclosedContract{
		TemplateID:       sc.TemplateID,
		ContractID:       sc.ContractID,
		CreatedEventBlob: blob,
		SynchronizerID:   sc.SynchronizerID,
	}
}

type ScanContractEntry struct {
	Contract *ScanContract `json:"contract"`
	DomainID string        `json:"domain_id"`
}

func (e ScanContractEntry) contract() *ScanContract {
	if e.Contract != nil {
		e.Contract.SynchronizerID = e.DomainID
	}
	return e.Contract
}

type AmuletRulesResponse struct {
	AmuletRules *ScanContractEntry `json:"amulet_rules"`
}

type OpenMiningRoundsResponse struct {
	OpenRounds    []ScanContractEntry `json:"open_mining_rounds"`
	IssuingRounds []ScanContractEntry `json:"issuing_mining_rounds"`
}

type TransferPreapprovalResponse struct {
	TransferPreapproval *ScanContractEntry `json:"transfer_preapproval"`
}

type FeaturedAppResponse struct {
	FeaturedAppRight *ScanContractEntry `json:"featured_app_right"`
}

type DSOPartyIDResponse struct {
	DSOPartyID string `json:"dso_party_id"`
}

type DSOResponse struct {
	DSO *ScanContractEntry `json:"dso"`
}

type TransferInstructionChoiceContext struct {
	ChoiceContextData  map[string]interface{}
	DisclosedContracts []*model.DisclosedContract
}

type scanAnyValue struct {
	tag   string
	value interface{}
}

func (v scanAnyValue) GetVariantTag() string {
	return v.tag
}

func (v scanAnyValue) GetVariantValue() interface{} {
	return v.value
}

func NewScanProxyClient(baseURL string, provider auth.TokenProvider) *ScanProxyClient {
	logger := log.Logger.With().Str("component", "scan-proxy-client").Logger()

	return &ScanProxyClient{
		baseURL:       baseURL,
		httpClient:    &http.Client{Timeout: 30 * time.Second},
		tokenProvider: provider,
		logger:        logger,
	}
}

func (s *ScanProxyClient) do(ctx context.Context, method, path string, body, result interface{}) error {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, s.baseURL+path, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	token, err := s.tokenProvider.Token()
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	s.logger.Debug().Str("method", method).Str("url", s.baseURL+path).Msg("scan-proxy request")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		s.logger.Error().Int("status", resp.StatusCode).Str("body", string(bodyBytes)).Msg("scan-proxy error response")
		return fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}
	return nil
}

func (s *ScanProxyClient) get(ctx context.Context, path string, result interface{}) error {
	return s.do(ctx, http.MethodGet, path, nil, result)
}

func (s *ScanProxyClient) post(ctx context.Context, path string, body, result interface{}) error {
	return s.do(ctx, http.MethodPost, path, body, result)
}

func (s *ScanProxyClient) GetTransferInstructionChoiceContext(ctx context.Context, transferInstructionID, choicePath string) (*TransferInstructionChoiceContext, error) {
	path := fmt.Sprintf(
		"/v0/scan-proxy/registry/transfer-instruction/v1/%s/choice-contexts/%s",
		url.PathEscape(transferInstructionID),
		choicePath,
	)
	var resp struct {
		ChoiceContextData struct {
			Values map[string]json.RawMessage `json:"values"`
		} `json:"choiceContextData"`
		DisclosedContracts []struct {
			TemplateID       string `json:"templateId"`
			ContractID       string `json:"contractId"`
			CreatedEventBlob string `json:"createdEventBlob"`
			SynchronizerID   string `json:"synchronizerId"`
		} `json:"disclosedContracts"`
	}
	if err := s.post(ctx, path, map[string]interface{}{}, &resp); err != nil {
		return nil, fmt.Errorf("failed to get transfer instruction choice context: %w", err)
	}

	values := make(map[string]interface{}, len(resp.ChoiceContextData.Values))
	for key, raw := range resp.ChoiceContextData.Values {
		value, err := decodeScanAnyValue(raw)
		if err != nil {
			return nil, fmt.Errorf("decode transfer instruction choice context %q: %w", key, err)
		}
		values[key] = value
	}

	disclosed := make([]*model.DisclosedContract, 0, len(resp.DisclosedContracts))
	for _, dc := range resp.DisclosedContracts {
		if dc.TemplateID == "" || dc.ContractID == "" || dc.CreatedEventBlob == "" {
			return nil, fmt.Errorf("invalid disclosed contract in transfer instruction choice context")
		}
		blob, err := base64.StdEncoding.DecodeString(dc.CreatedEventBlob)
		if err != nil {
			return nil, fmt.Errorf("decode disclosed contract blob for %s: %w", dc.ContractID, err)
		}
		disclosed = append(disclosed, &model.DisclosedContract{
			TemplateID:       dc.TemplateID,
			ContractID:       dc.ContractID,
			CreatedEventBlob: blob,
			SynchronizerID:   dc.SynchronizerID,
		})
	}

	return &TransferInstructionChoiceContext{
		ChoiceContextData:  values,
		DisclosedContracts: disclosed,
	}, nil
}

func decodeScanAnyValue(raw json.RawMessage) (scanAnyValue, error) {
	var tagged struct {
		Tag   string          `json:"tag"`
		Value json.RawMessage `json:"value"`
	}
	if err := json.Unmarshal(raw, &tagged); err != nil {
		return scanAnyValue{}, fmt.Errorf("unmarshal tagged AnyValue: %w", err)
	}
	if tagged.Tag == "" {
		return scanAnyValue{}, fmt.Errorf("missing AnyValue tag")
	}

	switch tagged.Tag {
	case "AV_Text":
		var v string
		if err := json.Unmarshal(tagged.Value, &v); err != nil {
			return scanAnyValue{}, err
		}
		return scanAnyValue{tag: tagged.Tag, value: types.TEXT(v)}, nil
	case "AV_ContractId":
		var v string
		if err := json.Unmarshal(tagged.Value, &v); err != nil {
			return scanAnyValue{}, err
		}
		return scanAnyValue{tag: tagged.Tag, value: types.CONTRACT_ID(v)}, nil
	case "AV_Party":
		var v string
		if err := json.Unmarshal(tagged.Value, &v); err != nil {
			return scanAnyValue{}, err
		}
		return scanAnyValue{tag: tagged.Tag, value: types.PARTY(v)}, nil
	case "AV_Bool":
		var v bool
		if err := json.Unmarshal(tagged.Value, &v); err != nil {
			return scanAnyValue{}, err
		}
		return scanAnyValue{tag: tagged.Tag, value: types.BOOL(v)}, nil
	case "AV_Int":
		var v int64
		if err := json.Unmarshal(tagged.Value, &v); err != nil {
			return scanAnyValue{}, err
		}
		return scanAnyValue{tag: tagged.Tag, value: types.INT64(v)}, nil
	case "AV_Decimal":
		var v string
		if err := json.Unmarshal(tagged.Value, &v); err != nil {
			return scanAnyValue{}, err
		}
		d, err := decimal.NewFromString(v)
		if err != nil {
			return scanAnyValue{}, err
		}
		return scanAnyValue{tag: tagged.Tag, value: types.NewNumericFromDecimal(d)}, nil
	case "AV_Date":
		var v string
		if err := json.Unmarshal(tagged.Value, &v); err != nil {
			return scanAnyValue{}, err
		}
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			return scanAnyValue{}, err
		}
		return scanAnyValue{tag: tagged.Tag, value: types.DATE(t)}, nil
	case "AV_Time":
		var v string
		if err := json.Unmarshal(tagged.Value, &v); err != nil {
			return scanAnyValue{}, err
		}
		t, err := time.Parse(time.RFC3339Nano, v)
		if err != nil {
			return scanAnyValue{}, err
		}
		return scanAnyValue{tag: tagged.Tag, value: types.TIMESTAMP(t)}, nil
	case "AV_RelTime":
		var micros int64
		if err := json.Unmarshal(tagged.Value, &micros); err != nil {
			var text string
			if textErr := json.Unmarshal(tagged.Value, &text); textErr != nil {
				return scanAnyValue{}, err
			}
			duration, parseErr := time.ParseDuration(text)
			if parseErr != nil {
				return scanAnyValue{}, parseErr
			}
			return scanAnyValue{tag: tagged.Tag, value: types.RELTIME(duration)}, nil
		}
		return scanAnyValue{tag: tagged.Tag, value: types.RELTIME(time.Duration(micros) * time.Microsecond)}, nil
	case "AV_List":
		var rawItems []json.RawMessage
		if err := json.Unmarshal(tagged.Value, &rawItems); err != nil {
			return scanAnyValue{}, err
		}
		values := make([]interface{}, len(rawItems))
		for i, rawItem := range rawItems {
			value, err := decodeScanAnyValue(rawItem)
			if err != nil {
				return scanAnyValue{}, err
			}
			values[i] = value
		}
		return scanAnyValue{tag: tagged.Tag, value: values}, nil
	case "AV_Map":
		values := types.TEXTMAP{}
		if err := json.Unmarshal(tagged.Value, &values); err != nil {
			return scanAnyValue{}, err
		}
		return scanAnyValue{tag: tagged.Tag, value: values}, nil
	default:
		return scanAnyValue{}, fmt.Errorf("unsupported AnyValue tag %q", tagged.Tag)
	}
}

// GetAmuletRules returns the current AmuletRules contract.
func (s *ScanProxyClient) GetAmuletRules(ctx context.Context) (*ScanContract, error) {
	var resp AmuletRulesResponse
	if err := s.get(ctx, "/v0/scan-proxy/amulet-rules", &resp); err != nil {
		return nil, fmt.Errorf("failed to get amulet rules: %w", err)
	}
	if resp.AmuletRules == nil || resp.AmuletRules.Contract == nil {
		return nil, fmt.Errorf("amulet rules not found in response")
	}
	contract := resp.AmuletRules.contract()
	if contract.ContractID == "" || contract.TemplateID == "" {
		return nil, fmt.Errorf("invalid amulet rules contract structure")
	}
	return contract, nil
}

// GetAmuletSynchronizerID returns the synchronizer id carried by AmuletRules.
func (s *ScanProxyClient) GetAmuletSynchronizerID(ctx context.Context) (string, error) {
	rules, err := s.GetAmuletRules(ctx)
	if err != nil {
		return "", err
	}
	if syncID, ok := rules.Payload["synchronizerId"].(string); ok {
		return syncID, nil
	}
	return "", fmt.Errorf("synchronizerId not found in amulet rules")
}

func (s *ScanProxyClient) getOpenAndIssuingRounds(ctx context.Context) (*OpenMiningRoundsResponse, error) {
	var resp OpenMiningRoundsResponse
	if err := s.get(ctx, "/v0/scan-proxy/open-and-issuing-mining-rounds", &resp); err != nil {
		return nil, fmt.Errorf("failed to get mining rounds: %w", err)
	}
	return &resp, nil
}

func unwrapEntries(entries []ScanContractEntry) ([]*ScanContract, error) {
	contracts := make([]*ScanContract, 0, len(entries))
	for i := range entries {
		c := entries[i].contract()
		if c == nil || c.ContractID == "" || c.TemplateID == "" {
			return nil, fmt.Errorf("invalid mining round contract structure at index %d", i)
		}
		contracts = append(contracts, c)
	}
	return contracts, nil
}

// GetOpenMiningRounds returns all open mining rounds.
func (s *ScanProxyClient) GetOpenMiningRounds(ctx context.Context) ([]*ScanContract, error) {
	resp, err := s.getOpenAndIssuingRounds(ctx)
	if err != nil {
		return nil, err
	}
	return unwrapEntries(resp.OpenRounds)
}

// GetIssuingMiningRounds returns all issuing mining rounds.
func (s *ScanProxyClient) GetIssuingMiningRounds(ctx context.Context) ([]*ScanContract, error) {
	resp, err := s.getOpenAndIssuingRounds(ctx)
	if err != nil {
		return nil, err
	}
	return unwrapEntries(resp.IssuingRounds)
}

// GetActiveOpenMiningRound returns the open mining round currently within its
// opensAt..targetClosesAt window, or nil if none is active.
func (s *ScanProxyClient) GetActiveOpenMiningRound(ctx context.Context) (*ScanContract, error) {
	rounds, err := s.GetOpenMiningRounds(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	for _, round := range rounds {
		opensAtStr, ok := round.Payload["opensAt"].(string)
		if !ok {
			continue
		}
		opensAt, err := time.Parse(time.RFC3339, opensAtStr)
		if err != nil {
			continue
		}
		targetClosesAtStr, ok := round.Payload["targetClosesAt"].(string)
		if !ok {
			continue
		}
		targetClosesAt, err := time.Parse(time.RFC3339, targetClosesAtStr)
		if err != nil {
			continue
		}
		if now.After(opensAt) && now.Before(targetClosesAt) {
			return round, nil
		}
	}
	return nil, nil
}

// GetTransferPreApprovalByParty returns the transfer pre-approval for a receiver
// party, or nil if none exists.
func (s *ScanProxyClient) GetTransferPreApprovalByParty(ctx context.Context, party string) (*ScanContract, error) {
	path := fmt.Sprintf("/v0/scan-proxy/transfer-preapprovals/by-party/%s", party)
	var resp TransferPreapprovalResponse
	if err := s.get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to get transfer preapproval: %w", err)
	}
	if resp.TransferPreapproval == nil || resp.TransferPreapproval.Contract == nil {
		return nil, nil
	}
	contract := resp.TransferPreapproval.contract()
	if contract.ContractID == "" || contract.TemplateID == "" {
		return nil, fmt.Errorf("invalid transfer preapproval contract structure")
	}
	return contract, nil
}

// GetFeaturedAppByProvider returns the featured-app right for a provider party,
// or nil if none exists.
func (s *ScanProxyClient) GetFeaturedAppByProvider(ctx context.Context, providerParty string) (*ScanContract, error) {
	path := fmt.Sprintf("/v0/scan-proxy/featured-apps/%s", providerParty)
	var resp FeaturedAppResponse
	if err := s.get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to get featured app: %w", err)
	}
	if resp.FeaturedAppRight == nil || resp.FeaturedAppRight.Contract == nil {
		return nil, nil
	}
	contract := resp.FeaturedAppRight.contract()
	if contract.ContractID == "" || contract.TemplateID == "" {
		return nil, fmt.Errorf("invalid featured app contract structure")
	}
	return contract, nil
}

// GetDSOPartyID returns the DSO party id.
func (s *ScanProxyClient) GetDSOPartyID(ctx context.Context) (string, error) {
	var resp DSOPartyIDResponse
	if err := s.get(ctx, "/v0/scan-proxy/dso-party-id", &resp); err != nil {
		return "", fmt.Errorf("failed to get DSO party ID: %w", err)
	}
	if resp.DSOPartyID == "" {
		return "", fmt.Errorf("DSO party ID not found in response")
	}
	return resp.DSOPartyID, nil
}

// GetDSO returns the DSO contract.
func (s *ScanProxyClient) GetDSO(ctx context.Context) (*ScanContract, error) {
	var resp DSOResponse
	if err := s.get(ctx, "/v0/scan-proxy/dso", &resp); err != nil {
		return nil, fmt.Errorf("failed to get DSO: %w", err)
	}
	if resp.DSO == nil || resp.DSO.Contract == nil {
		return nil, fmt.Errorf("DSO not found in response")
	}
	contract := resp.DSO.contract()
	if contract.ContractID == "" || contract.TemplateID == "" {
		return nil, fmt.Errorf("invalid DSO contract structure")
	}
	return contract, nil
}
