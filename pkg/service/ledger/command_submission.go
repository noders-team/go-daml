package ledger

import (
	"context"

	"google.golang.org/grpc"

	"github.com/noders-team/go-daml/pkg/model"
	v2 "github.com/noders-team/go-daml/proto/com/daml/ledger/api/v2"
)

type CommandSubmission interface {
	Submit(ctx context.Context, req *model.SubmitRequest) (*model.SubmitResponse, error)
}

type commandSubmission struct {
	client v2.CommandSubmissionServiceClient
}

func NewCommandSubmissionClient(conn *grpc.ClientConn) *commandSubmission {
	client := v2.NewCommandSubmissionServiceClient(conn)
	return &commandSubmission{
		client: client,
	}
}

func (c *commandSubmission) Submit(ctx context.Context, req *model.SubmitRequest) (*model.SubmitResponse, error) {
	protoReq := &v2.SubmitRequest{
		Commands: commandsToProto(req.Commands),
	}

	_, err := c.client.Submit(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	return &model.SubmitResponse{}, nil
}
