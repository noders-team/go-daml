package wallet

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/noders-team/go-daml/pkg/client"
	"github.com/noders-team/go-daml/pkg/model"
	"github.com/noders-team/go-daml/pkg/service/ledger"
)

// When both a closed response channel and a queued error are ready, select
// picks one at random. Each call on code that ignores the queued error drops
// it with probability 1/2, so all streamCalls passing by chance has
// probability 2^-100.
const streamCalls = 100

// stubStateService returns channels in the state GetActiveContracts leaves
// them in once its stream goroutine has exited.
type stubStateService struct {
	ledger.StateService
	responses []*model.GetActiveContractsResponse
	err       error
}

func (s stubStateService) GetActiveContracts(context.Context, *model.GetActiveContractsRequest) (<-chan *model.GetActiveContractsResponse, <-chan error) {
	respCh := make(chan *model.GetActiveContractsResponse, len(s.responses))
	for _, r := range s.responses {
		respCh <- r
	}
	close(respCh)
	errCh := make(chan error, 1)
	if s.err != nil {
		errCh <- s.err
	}
	close(errCh)
	return respCh, errCh
}

func newStubController(t *testing.T, stub stubStateService) *tokenStandardController {
	t.Helper()
	ctrl, err := NewTokenStandardController("user", &client.DamlBindingClient{StateService: stub}, nil)
	require.NoError(t, err)
	ctrl.SetPartyID("Alice")
	return ctrl.(*tokenStandardController)
}

func TestListContractsByInterfaceReportsStreamError(t *testing.T) {
	streamErr := errors.New("stream broke")
	ctrl := newStubController(t, stubStateService{err: streamErr})

	dropped := 0
	for i := 0; i < streamCalls; i++ {
		events, err := ctrl.ListContractsByInterface(context.Background(), "#pkg:Mod:Iface", 0)
		if err == nil {
			dropped++
			continue
		}
		require.ErrorIs(t, err, streamErr)
		require.Nil(t, events)
	}
	t.Logf("%d of %d calls returned a nil error", dropped, streamCalls)
	require.Zero(t, dropped)
}

func TestListContractsByInterfaceCompleteStream(t *testing.T) {
	const sent = 5
	responses := make([]*model.GetActiveContractsResponse, sent)
	for i := range responses {
		responses[i] = &model.GetActiveContractsResponse{
			ContractEntry: &model.ActiveContractEntry{
				ActiveContract: &model.ActiveContract{CreatedEvent: &model.CreatedEvent{ContractID: "cid"}},
			},
		}
	}
	ctrl := newStubController(t, stubStateService{responses: responses})

	for i := 0; i < streamCalls; i++ {
		events, err := ctrl.ListContractsByInterface(context.Background(), "#pkg:Mod:Iface", 0)
		require.NoError(t, err)
		require.Len(t, events, sent)
	}
}
