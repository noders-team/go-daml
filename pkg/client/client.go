package client

import (
	"context"
	"crypto/tls"
	"fmt"

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
	opts := c.buildDialOptions()

	conn, err := grpc.NewClient(c.config.Address, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DAML ledger: %w", err)
	}

	c.conn = conn

	var adminConn *grpc.ClientConn
	if c.config.AdminAddress != "" {
		adminConn, err = grpc.NewClient(c.config.AdminAddress, opts...)
		if err != nil {
			_ = c.conn.Close()
			return nil, fmt.Errorf("failed to connect to DAML admin endpoint: %w", err)
		}
		c.adminConn = adminConn
	}

	return NewConnection(c, conn, adminConn), nil
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

func (c *Client) buildDialOptions() []grpc.DialOption {
	opts := append([]grpc.DialOption{}, c.config.GRPCDialOptions...)

	if c.config.TLS != nil {
		tlsConfig := c.buildTLSConfig()
		creds := credentials.NewTLS(tlsConfig)
		opts = append(opts, grpc.WithTransportCredentials(creds))

		auth := c.config.Auth
		if auth != nil {
			opts = append(opts, grpc.WithPerRPCCredentials(auth.TokenProvider))
		}
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

		auth := c.config.Auth
		if auth != nil {
			opts = append(opts,
				grpc.WithUnaryInterceptor(auth.TokenProvider.UnaryInterceptor()),
				grpc.WithStreamInterceptor(auth.TokenProvider.StreamInterceptor()),
			)
		}
	}

	return opts
}

func (c *Client) buildTLSConfig() *tls.Config {
	tlsConfig := &tls.Config{
		ServerName:         c.config.TLS.ServerName,
		InsecureSkipVerify: c.config.TLS.InsecureSkipVerify,
	}

	if c.config.TLS.CertFile != "" {
	}

	return tlsConfig
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
