package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type KeycloakConfig struct {
	OIDCURL      string
	TokenURL     string
	ClientID     string
	ClientSecret string
	Audience     string
}

func (c KeycloakConfig) Enabled() bool {
	return c.ClientID != "" && c.ClientSecret != "" && (c.TokenURL != "" || c.OIDCURL != "")
}

type keycloakTokenProvider struct {
	mu           sync.Mutex
	httpClient   *http.Client
	tokenURL     string
	clientID     string
	clientSecret string
	audience     string
	accessToken  string
	refreshToken string
	expiresAt    time.Time
}

func NewKeycloakTokenProvider(cfg KeycloakConfig) (*keycloakTokenProvider, error) {
	tokenURL := cfg.TokenURL
	if tokenURL == "" && cfg.OIDCURL != "" {
		base := strings.TrimRight(cfg.OIDCURL, "/")
		tokenURL = base + "/protocol/openid-connect/token"
	}

	if tokenURL == "" {
		return nil, fmt.Errorf("keycloak token url is empty")
	}

	return &keycloakTokenProvider{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		tokenURL:     tokenURL,
		clientID:     cfg.ClientID,
		clientSecret: cfg.ClientSecret,
		audience:     cfg.Audience,
	}, nil
}

func (p *keycloakTokenProvider) Token() (string, error) {
	if p.accessToken != "" && time.Now().Before(p.expiresAt) {
		return p.accessToken, nil
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.refreshToken != "" {
		if err := p.refreshLocked(); err == nil {
			return p.accessToken, nil
		} else {
			log.Warn().Err(err).Msg("failed to refresh keycloak token, fetching a new access token")
		}
	}

	if err := p.fetchClientCredentialsLocked(); err != nil {
		return "", err
	}

	return p.accessToken, nil
}

func (p *keycloakTokenProvider) refreshLocked() error {
	values := url.Values{}
	values.Set("grant_type", "refresh_token")
	values.Set("refresh_token", p.refreshToken)
	values.Set("client_id", p.clientID)
	values.Set("client_secret", p.clientSecret)
	if p.audience != "" {
		values.Set("audience", p.audience)
	}
	return p.requestTokenLocked(values)
}

func (p *keycloakTokenProvider) fetchClientCredentialsLocked() error {
	values := url.Values{}
	values.Set("grant_type", "client_credentials")
	values.Set("client_id", p.clientID)
	values.Set("client_secret", p.clientSecret)
	if p.audience != "" {
		values.Set("audience", p.audience)
	}
	return p.requestTokenLocked(values)
}

func (p *keycloakTokenProvider) requestTokenLocked(values url.Values) error {
	req, err := http.NewRequest("POST", p.tokenURL, strings.NewReader(values.Encode()))
	if err != nil {
		return fmt.Errorf("failed to build keycloak token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("keycloak token request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if err != nil {
		return fmt.Errorf("failed to read keycloak token response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("keycloak token request failed: status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("failed to parse keycloak token response: %w", err)
	}
	if tokenResp.AccessToken == "" {
		return fmt.Errorf("keycloak token response missing access_token")
	}

	p.accessToken = tokenResp.AccessToken
	if tokenResp.RefreshToken != "" {
		p.refreshToken = tokenResp.RefreshToken
	}

	expiresIn := time.Duration(tokenResp.ExpiresIn) * time.Second
	if expiresIn <= 0 {
		expiresIn = 30 * time.Minute
	}
	p.expiresAt = time.Now().Add(expiresIn - 30*time.Second)
	return nil
}

func (p *keycloakTokenProvider) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	token, err := p.Token()
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"authorization": fmt.Sprintf("Bearer %s", token),
	}, nil
}

func (p *keycloakTokenProvider) RequireTransportSecurity() bool {
	return false
}

func (p *keycloakTokenProvider) UnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		token, err := p.Token()
		if err != nil {
			return err
		}
		if token != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %s", token))
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func (p *keycloakTokenProvider) StreamInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		token, err := p.Token()
		if err != nil {
			return nil, err
		}
		if token != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %s", token))
		}

		return streamer(ctx, desc, cc, method, opts...)
	}
}
