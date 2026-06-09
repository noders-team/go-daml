package auth

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type BearerTokenAuth struct {
	token string
}

func NewBearerTokenProvider(token string) *BearerTokenAuth {
	return &BearerTokenAuth{
		token: token,
	}
}

func (b *BearerTokenAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	if b.token == "" {
		return nil, nil
	}

	return map[string]string{
		"authorization": fmt.Sprintf("Bearer %s", b.token),
	}, nil
}

func (b *BearerTokenAuth) RequireTransportSecurity() bool {
	return false
}

func (b *BearerTokenAuth) UnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if b.token != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %s", b.token))
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func (b *BearerTokenAuth) StreamInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		if b.token != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %s", b.token))
		}

		return streamer(ctx, desc, cc, method, opts...)
	}
}

func (b *BearerTokenAuth) Token() (string, error) {
	return b.token, nil
}
