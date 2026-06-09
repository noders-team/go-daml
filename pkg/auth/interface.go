package auth

import (
	"context"

	"google.golang.org/grpc"
)

type TokenProvider interface {
	Token() (string, error)
	UnaryInterceptor() grpc.UnaryClientInterceptor
	StreamInterceptor() grpc.StreamClientInterceptor
	GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error)
	RequireTransportSecurity() bool
}
