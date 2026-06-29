package admin

import (
	"context"

	"google.golang.org/grpc"

	"github.com/noders-team/go-daml/pkg/model"
	adminv2 "github.com/noders-team/go-daml/proto/com/daml/ledger/api/v2/admin"
)

type ParticipantPruning interface {
	Prune(ctx context.Context, pruneRequest *model.PruneRequest) error
}

type participantPruning struct {
	client adminv2.ParticipantPruningServiceClient
}

func NewParticipantPruningClient(conn *grpc.ClientConn) *participantPruning {
	client := adminv2.NewParticipantPruningServiceClient(conn)
	return &participantPruning{
		client: client,
	}
}

func (c *participantPruning) Prune(ctx context.Context, pruneRequest *model.PruneRequest) error {
	req := &adminv2.PruneRequest{
		PruneUpTo:                 pruneRequest.PruneUpTo,
		SubmissionId:              pruneRequest.SubmissionID,
		PruneAllDivulgedContracts: pruneRequest.PruneAllDivulgedContracts,
	}

	_, err := c.client.Prune(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
