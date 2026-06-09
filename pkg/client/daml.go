package client

import (
	"context"

	"github.com/noders-team/go-daml/pkg/auth"
)

type DamlClient struct {
	config *Config
}

type TlsConfig struct {
	Certificate string
}

func NewDamlClient(grpcAddress string, provider auth.TokenProvider) *DamlClient {
	config := &Config{
		Address: grpcAddress,
	}
	config.Auth = &AuthConfig{
		TokenProvider: provider,
	}
	return &DamlClient{
		config: config,
	}
}

func (c *DamlClient) WithTLSConfig(cfg TlsConfig) *DamlClient {
	c.config.TLS = &TLSConfig{
		CertFile: cfg.Certificate,
	}
	return c
}

func (c *DamlClient) WithAdminAddress(addr string) *DamlClient {
	c.config.AdminAddress = addr
	return c
}

func (c *DamlClient) Build(ctx context.Context) (*DamlBindingClient, error) {
	client := NewClient(c.config)
	conn, err := client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return NewDamlBindingClient(c, conn), nil
}
