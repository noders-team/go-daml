package ledger

import (
	"context"
	"io"

	rpcstatus "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/noders-team/go-daml/pkg/model"
	v2 "github.com/noders-team/go-daml/proto/com/daml/ledger/api/v2"
)

type CommandCompletion interface {
	CompletionStream(ctx context.Context, req *model.CompletionStreamRequest) (<-chan *model.CompletionStreamResponse, <-chan error)
}

type commandCompletion struct {
	client v2.CommandCompletionServiceClient
}

func NewCommandCompletionClient(conn *grpc.ClientConn) *commandCompletion {
	client := v2.NewCommandCompletionServiceClient(conn)
	return &commandCompletion{
		client: client,
	}
}

func (c *commandCompletion) CompletionStream(ctx context.Context, req *model.CompletionStreamRequest) (<-chan *model.CompletionStreamResponse, <-chan error) {
	streamReq := &v2.CompletionStreamRequest{
		UserId:         req.UserID,
		Parties:        req.Parties,
		BeginExclusive: req.BeginExclusive,
	}

	stream, err := c.client.CompletionStream(ctx, streamReq)
	if err != nil {
		errCh := make(chan error, 1)
		errCh <- err
		close(errCh)
		return nil, errCh
	}

	responseCh := make(chan *model.CompletionStreamResponse)
	errCh := make(chan error, 1)

	go func() {
		defer close(responseCh)
		defer close(errCh)

		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				errCh <- err
				return
			}

			modelResp := completionStreamResponseFromProto(resp)
			if modelResp != nil {
				select {
				case responseCh <- modelResp:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return responseCh, errCh
}

func completionStreamResponseFromProto(pb *v2.CompletionStreamResponse) *model.CompletionStreamResponse {
	if pb == nil {
		return nil
	}

	resp := &model.CompletionStreamResponse{}

	switch cr := pb.CompletionResponse.(type) {
	case *v2.CompletionStreamResponse_Completion:
		if cr.Completion != nil {
			resp.Response = completionFromProto(cr.Completion)
		}
	case *v2.CompletionStreamResponse_OffsetCheckpoint:
		resp.Response = model.OffsetCheckpoint{
			Offset: cr.OffsetCheckpoint.Offset,
		}
	}

	return resp
}

func completionFromProto(pb *v2.Completion) model.Completion {
	comp := model.Completion{
		CommandID:       pb.CommandId,
		UpdateID:        pb.UpdateId,
		SubmissionID:    pb.SubmissionId,
		Offset:          pb.Offset,
		PaidTrafficCost: paidTrafficCostFromUnknown(pb, completionPaidTrafficCostField),
	}

	// Note: proto Completion doesn't have CompletedAt or TransactionID fields
	// TransactionID would come from a separate transaction response

	if pb.Status != nil {
		comp.Status = statusFromProto(pb.Status)
	}

	return comp
}

func statusFromProto(pb *rpcstatus.Status) model.Status {
	if pb == nil {
		return model.StatusOK{}
	}

	code := codes.Code(pb.Code)
	if code == codes.OK {
		return model.StatusOK{}
	}

	return model.StatusError{
		Code:    pb.Code,
		Message: pb.Message,
	}
}
