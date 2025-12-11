package ledger

import (
	"context"

	"google.golang.org/grpc"

	v2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/noders-team/go-daml/pkg/model"
)

type EventQuery interface {
	GetEventsByContractID(ctx context.Context, req *model.GetEventsByContractIDRequest) (*model.GetEventsByContractIDResponse, error)
}

type eventQuery struct {
	client v2.EventQueryServiceClient
}

func NewEventQueryClient(conn *grpc.ClientConn) *eventQuery {
	client := v2.NewEventQueryServiceClient(conn)
	return &eventQuery{
		client: client,
	}
}

func (c *eventQuery) GetEventsByContractID(ctx context.Context, req *model.GetEventsByContractIDRequest) (*model.GetEventsByContractIDResponse, error) {
	protoReq := &v2.GetEventsByContractIdRequest{
		ContractId:  req.ContractID,
		EventFormat: eventFormatToProto(req.EventFormat),
	}

	resp, err := c.client.GetEventsByContractId(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	return getEventsByContractIDResponseFromProto(resp), nil
}

func getEventsByContractIDResponseFromProto(pb *v2.GetEventsByContractIdResponse) *model.GetEventsByContractIDResponse {
	if pb == nil {
		return nil
	}

	resp := &model.GetEventsByContractIDResponse{}

	if pb.Created != nil && pb.Created.CreatedEvent != nil {
		resp.CreateEvent = createdEventFromProto(pb.Created.CreatedEvent)
	}

	if pb.Archived != nil && pb.Archived.ArchivedEvent != nil {
		resp.ArchiveEvent = archivedEventFromProto(pb.Archived.ArchivedEvent)
	}

	return resp
}
