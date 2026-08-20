package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	config    *Config
	conn      *grpc.ClientConn
	adminConn *grpc.ClientConn
}

func NewClient(config *Config) *Client {
	return &Client{
		config: config,
	}
}

func (c *Client) Connect(ctx context.Context) (*Connection, error) {
	opts, err := c.buildDialOptions(c.config.TLS, c.config.Auth)
	if err != nil {
		return nil, fmt.Errorf("failed to build ledger dial options: %w", err)
	}

	conn, err := grpc.NewClient(c.config.Address, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DAML ledger: %w", err)
	}

	c.conn = conn

	var adminConn *grpc.ClientConn
	if c.config.AdminAddress != "" {
		adminOpts := opts
		if c.config.AdminTLS != nil || c.config.AdminAuth != nil {
			adminOpts, err = c.buildDialOptions(firstNonNilTLS(c.config.AdminTLS, c.config.TLS), firstNonNilAuth(c.config.AdminAuth, c.config.Auth))
			if err != nil {
				_ = c.conn.Close()
				return nil, fmt.Errorf("failed to build admin dial options: %w", err)
			}
		}
		adminConn, err = grpc.NewClient(c.config.AdminAddress, adminOpts...)
		if err != nil {
			_ = c.conn.Close()
			return nil, fmt.Errorf("failed to connect to DAML admin endpoint: %w", err)
		}
		c.adminConn = adminConn
	}

	return NewConnection(c, conn, adminConn), nil
}

func firstNonNilTLS(primary, fallback *TLSConfig) *TLSConfig {
	if primary != nil {
		return primary
	}
	return fallback
}

func firstNonNilAuth(primary, fallback *AuthConfig) *AuthConfig {
	if primary != nil {
		return primary
	}
	return fallback
}

func (c *Client) Close() error {
	var err error
	if c.conn != nil {
		err = c.conn.Close()
	}
	if c.adminConn != nil {
		if adminErr := c.adminConn.Close(); adminErr != nil && err == nil {
			err = adminErr
		}
	}
	return err
}

func (c *Client) buildDialOptions(tlsCfg *TLSConfig, authCfg *AuthConfig) ([]grpc.DialOption, error) {
	opts := append([]grpc.DialOption{}, c.config.GRPCDialOptions...)

	if tlsCfg != nil {
		builtTLS, err := buildTLSConfig(tlsCfg)
		if err != nil {
			return nil, err
		}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(builtTLS)))

		if authCfg != nil {
			opts = append(opts, grpc.WithPerRPCCredentials(authCfg.TokenProvider))
		}
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

		if authCfg != nil {
			opts = append(opts,
				grpc.WithUnaryInterceptor(authCfg.TokenProvider.UnaryInterceptor()),
				grpc.WithStreamInterceptor(authCfg.TokenProvider.StreamInterceptor()),
			)
		}
	}

	return opts, nil
}

func buildTLSConfig(cfg *TLSConfig) (*tls.Config, error) {
	tlsConfig := &tls.Config{
		ServerName:         cfg.ServerName,
		InsecureSkipVerify: cfg.InsecureSkipVerify,
	}

	if cfg.CertFile != "" {
		pemBytes, err := os.ReadFile(cfg.CertFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read TLS CertFile %q: %w", cfg.CertFile, err)
		}
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(pemBytes) {
			return nil, fmt.Errorf("failed to parse TLS CertFile %q as PEM", cfg.CertFile)
		}
		tlsConfig.RootCAs = pool
	}

	if cfg.ClientCertFile != "" && cfg.ClientKeyFile != "" {
		cert, err := tls.LoadX509KeyPair(cfg.ClientCertFile, cfg.ClientKeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load mTLS client key pair: %w", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	return tlsConfig, nil
}

type Connection struct {
	client    *Client
	conn      *grpc.ClientConn
	adminConn *grpc.ClientConn
}

func NewConnection(client *Client, conn *grpc.ClientConn, adminConn *grpc.ClientConn) *Connection {
	return &Connection{
		client:    client,
		conn:      conn,
		adminConn: adminConn,
	}
}

func (c *Connection) GRPCConn() *grpc.ClientConn {
	return c.conn
}

func (c *Connection) AdminGRPCConn() *grpc.ClientConn {
	if c.adminConn != nil {
		return c.adminConn
	}
	return c.conn
}
