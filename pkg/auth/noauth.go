package auth

import (
	"context"

	"google.golang.org/grpc"
)

type NoAuth struct{}

func NewNoAuth() *NoAuth {
	return &NoAuth{}
}

func (b *NoAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return make(map[string]string), nil
}

func (b *NoAuth) RequireTransportSecurity() bool {
	return false
}

func (b *NoAuth) UnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func (b *NoAuth) StreamInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return streamer(ctx, desc, cc, method, opts...)
	}
}

func (b *NoAuth) Token() (string, error) {
	return "", nil
}
