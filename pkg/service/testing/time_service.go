package testing

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/noders-team/go-daml/pkg/model"
	testingv2 "github.com/noders-team/go-daml/proto/com/daml/ledger/api/v2/testing"
)

type TimeService interface {
	GetTime(ctx context.Context, req *model.GetTimeRequest) (*model.GetTimeResponse, error)
	SetTime(ctx context.Context, req *model.SetTimeRequest) (*model.SetTimeResponse, error)
}

type timeService struct {
	client testingv2.TimeServiceClient
}

func NewTimeServiceClient(conn *grpc.ClientConn) TimeService {
	client := testingv2.NewTimeServiceClient(conn)
	return &timeService{
		client: client,
	}
}

func (c *timeService) GetTime(ctx context.Context, req *model.GetTimeRequest) (*model.GetTimeResponse, error) {
	protoReq := &testingv2.GetTimeRequest{}

	resp, err := c.client.GetTime(ctx, protoReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get time: %w", err)
	}

	response := &model.GetTimeResponse{}
	if resp.CurrentTime != nil {
		response.CurrentTime = resp.CurrentTime.AsTime()
	}

	return response, nil
}

func (c *timeService) SetTime(ctx context.Context, req *model.SetTimeRequest) (*model.SetTimeResponse, error) {
	protoReq := &testingv2.SetTimeRequest{
		CurrentTime: timestamppb.New(req.CurrentTime),
		NewTime:     timestamppb.New(req.NewTime),
	}

	_, err := c.client.SetTime(ctx, protoReq)
	if err != nil {
		return nil, fmt.Errorf("failed to set time: %w", err)
	}

	return &model.SetTimeResponse{}, nil
}
