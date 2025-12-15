package client

import (
	"context"

	"github.com/noders-team/go-daml/pkg/auth"
)

type DamlClient struct {
	config *Config
}

func NewDamlClient(token string, grpcAddress string) *DamlClient {
	config := &Config{
		Address: grpcAddress,
	}
	if token != "" {
		config.Auth = &AuthConfig{
			Token: token,
		}
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

type TlsConfig struct {
	Certificate string
}

func Connect(ctx context.Context, address string, opts ...ConfigOption) (*Connection, error) {
	allOpts := append([]ConfigOption{WithAddress(address)}, opts...)
	config := NewConfig(allOpts...)
	client := NewClient(config)
	return client.Connect(ctx)
}

func ConnectWithToken(ctx context.Context, address, token string, opts ...ConfigOption) (*Connection, error) {
	allOpts := append([]ConfigOption{WithAddress(address), WithToken(token)}, opts...)
	config := NewConfig(allOpts...)
	client := NewClient(config)
	return client.Connect(ctx)
}

func ConnectWithTokenProvider(ctx context.Context, address string, provider auth.TokenProvider, opts ...ConfigOption) (*Connection, error) {
	allOpts := append([]ConfigOption{WithAddress(address), WithTokenProvider(provider)}, opts...)
	config := NewConfig(allOpts...)
	client := NewClient(config)
	return client.Connect(ctx)
}
