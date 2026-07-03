package client

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/noders-team/go-daml/pkg/auth"
	"github.com/noders-team/go-daml/pkg/model"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
		bodyBytes, _ := io.ReadAll(resp.Body)
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
