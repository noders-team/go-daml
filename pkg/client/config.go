package client

import (
	"github.com/noders-team/go-daml/pkg/auth"
)

type Config struct {
	Address      string
	AdminAddress string
	TLS          *TLSConfig
	Auth         *AuthConfig
}

type TLSConfig struct {
	CertFile           string
	ServerName         string
	InsecureSkipVerify bool
}

type AuthConfig struct {
	Token         string
	TokenProvider auth.TokenProvider
}

type ConfigOption func(*Config)

func WithAddress(addr string) ConfigOption {
	return func(c *Config) {
		c.Address = addr
	}
}

func WithAdminAddress(addr string) ConfigOption {
	return func(c *Config) {
		c.AdminAddress = addr
	}
}

func WithTLS(tls *TLSConfig) ConfigOption {
	return func(c *Config) {
		c.TLS = tls
	}
}

func WithToken(token string) ConfigOption {
	return func(c *Config) {
		if c.Auth == nil {
			c.Auth = &AuthConfig{}
		}
		c.Auth.Token = token
	}
}

func WithTokenProvider(provider auth.TokenProvider) ConfigOption {
	return func(c *Config) {
		if c.Auth == nil {
			c.Auth = &AuthConfig{}
		}
		c.Auth.TokenProvider = provider
	}
}

func NewConfig(opts ...ConfigOption) *Config {
	cfg := &Config{}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}
