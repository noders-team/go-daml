package client

import (
	"context"
	"crypto/tls"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type DamlClient struct {
	grpcAddress string
	token       string
	tlsConfig   *TlsConfig
}

type TlsConfig struct {
	Certificate string
}

func NewDamlClient(token string, grpcAddress string) *DamlClient {
	return &DamlClient{
		grpcAddress: grpcAddress,
		token:       token,
	}
}

func (c *DamlClient) WithTLSConfig(cfg TlsConfig) *DamlClient {
	c.tlsConfig = &cfg
	return c
}

func (c *DamlClient) Build(ctx context.Context) (*damlBindingClient, error) {
	var opts []grpc.DialOption

	if c.tlsConfig != nil {
		if c.tlsConfig.Certificate != "" {
			creds, err := credentials.NewClientTLSFromFile(c.tlsConfig.Certificate, "")
			if err != nil {
				return nil, fmt.Errorf("failed to load TLS credentials: %w", err)
			}
			opts = append(opts, grpc.WithTransportCredentials(creds))
		} else {
			opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
		}
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	conn, err := grpc.Dial(c.grpcAddress, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DAML ledger: %w", err)
	}

	return NewDamlBindingClient(c, conn), nil
}
