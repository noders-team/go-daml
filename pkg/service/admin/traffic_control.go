package admin

import (
	"context"

	participantv30 "github.com/noders-team/go-daml/proto/com/digitalasset/canton/admin/participant/v30"
	"google.golang.org/grpc"
)

type TrafficState struct {
	ExtraTrafficPurchased int64
	ExtraTrafficConsumed  int64
	BaseTrafficRemainder  int64
	LastConsumedCost      int64
}

type TrafficControl interface {
	GetState(ctx context.Context, synchronizerID string) (*TrafficState, error)
}

type trafficControl struct {
	client participantv30.TrafficControlServiceClient
}

func NewTrafficControlClient(conn *grpc.ClientConn) TrafficControl {
	return &trafficControl{
		client: participantv30.NewTrafficControlServiceClient(conn),
	}
}

func (c *trafficControl) GetState(ctx context.Context, synchronizerID string) (*TrafficState, error) {
	resp, err := c.client.TrafficControlState(ctx, &participantv30.TrafficControlStateRequest{
		SynchronizerId: synchronizerID,
	})
	if err != nil {
		return nil, err
	}
	if resp.TrafficState == nil {
		return &TrafficState{}, nil
	}
	return &TrafficState{
		ExtraTrafficPurchased: resp.TrafficState.GetExtraTrafficPurchased(),
		ExtraTrafficConsumed:  resp.TrafficState.GetExtraTrafficConsumed(),
		BaseTrafficRemainder:  resp.TrafficState.GetBaseTrafficRemainder(),
		LastConsumedCost:      int64(resp.TrafficState.GetLastConsumedCost()),
	}, nil
}
