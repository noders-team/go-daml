package ledger

import (
	"context"

	"google.golang.org/grpc"

	"github.com/noders-team/go-daml/pkg/model"
	v2 "github.com/noders-team/go-daml/proto/com/daml/ledger/api/v2"
)

type CommandService interface {
	SubmitAndWait(ctx context.Context, req *model.SubmitAndWaitRequest) (*model.SubmitAndWaitResponse, error)
}

type commandService struct {
	client v2.CommandServiceClient
}

func NewCommandServiceClient(conn *grpc.ClientConn) *commandService {
	client := v2.NewCommandServiceClient(conn)
	return &commandService{
		client: client,
	}
}

func (c *commandService) SubmitAndWait(ctx context.Context, req *model.SubmitAndWaitRequest) (*model.SubmitAndWaitResponse, error) {
	protoReq := &v2.SubmitAndWaitRequest{
		Commands: commandsToProto(req.Commands),
	}

	resp, err := c.client.SubmitAndWait(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	return &model.SubmitAndWaitResponse{
		UpdateID:         resp.UpdateId,
		CompletionOffset: resp.CompletionOffset,
	}, nil
}
