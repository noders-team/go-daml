package client

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/noders-team/go-daml/pkg/service/ledger"
	v2 "github.com/noders-team/go-daml/proto/com/daml/ledger/api/v2"
)

type stateStub struct {
	v2.UnimplementedStateServiceServer
	sent    int
	failure error
}

func (s *stateStub) GetLedgerEnd(context.Context, *v2.GetLedgerEndRequest) (*v2.GetLedgerEndResponse, error) {
	return &v2.GetLedgerEndResponse{Offset: 1}, nil
}

func (s *stateStub) GetActiveContracts(_ *v2.GetActiveContractsRequest, stream grpc.ServerStreamingServer[v2.GetActiveContractsResponse]) error {
	for i := 0; i < s.sent; i++ {
		resp := &v2.GetActiveContractsResponse{
			ContractEntry: &v2.GetActiveContractsResponse_ActiveContract{
				ActiveContract: &v2.ActiveContract{
					CreatedEvent: &v2.CreatedEvent{
						ContractId: fmt.Sprintf("cid-%d", i),
						TemplateId: &v2.Identifier{PackageId: "pkg", ModuleName: "Mod", EntityName: "Tpl"},
					},
				},
			},
		}
		if err := stream.Send(resp); err != nil {
			return err
		}
	}
	return s.failure
}

func newStubClient(t *testing.T, stub *stateStub) *DamlBindingClient {
	t.Helper()

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	srv := grpc.NewServer()
	v2.RegisterStateServiceServer(srv, stub)
	go func() { _ = srv.Serve(lis) }()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = conn.Close()
		srv.Stop()
	})

	return &DamlBindingClient{StateService: ledger.NewStateServiceClient(conn)}
}

func TestFindContractsByTemplateStreamError(t *testing.T) {
	const calls = 500

	tests := []struct {
		name      string
		sent      int
		maxEntry  int
		failure   error
		wantCount int
		wantErr   bool
	}{
		{name: "fails after five contracts", sent: 5, failure: status.Error(codes.Unavailable, "stream broken"), wantCount: 0, wantErr: true},
		{name: "fails after one contract", sent: 1, failure: status.Error(codes.Unavailable, "stream broken"), wantCount: 0, wantErr: true},
		{name: "fails before any contract", sent: 0, failure: status.Error(codes.Unavailable, "stream broken"), wantCount: 0, wantErr: true},
		{name: "completes after five contracts", sent: 5, failure: nil, wantCount: 5, wantErr: false},
		{name: "completes without contracts", sent: 0, failure: nil, wantCount: 0, wantErr: false},
		{name: "stops at max entries", sent: 5, maxEntry: 2, failure: nil, wantCount: 2, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub := &stateStub{sent: tt.sent, failure: tt.failure}
			query := NewContractQuery[map[string]any](newStubClient(t, stub))

			lost := 0
			for range calls {
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				got, err := query.FindContractsByTemplate(ctx, "party", "pkg:Mod:Tpl", tt.maxEntry)
				require.NoError(t, ctx.Err())
				cancel()

				if tt.wantErr && err == nil {
					lost++
					continue
				}
				if !tt.wantErr {
					require.NoError(t, err)
				}
				require.Len(t, got, tt.wantCount)
			}

			require.Zerof(t, lost, "stream error was swallowed in %d of %d calls", lost, calls)
		})
	}
}
