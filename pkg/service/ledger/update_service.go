package ledger

import (
	"context"
	"io"

	"google.golang.org/grpc"

	v2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/noders-team/go-daml/pkg/model"
)

type UpdateService interface {
	GetUpdates(ctx context.Context, req *model.GetUpdatesRequest) (<-chan *model.GetUpdatesResponse, <-chan error)
	GetUpdateById(ctx context.Context, req *model.GetUpdateByIDRequest) (*model.GetUpdateResponse, error)
	GetTransactionByID(ctx context.Context, req *model.GetTransactionByIDRequest) (*model.GetTransactionResponse, error)
	GetTransactionByOffset(ctx context.Context, req *model.GetTransactionByOffsetRequest) (*model.GetTransactionResponse, error)
}

type updateService struct {
	client v2.UpdateServiceClient
}

func NewUpdateServiceClient(conn *grpc.ClientConn) *updateService {
	client := v2.NewUpdateServiceClient(conn)
	return &updateService{
		client: client,
	}
}

func (c *updateService) Client() v2.UpdateServiceClient {
	return c.client
}

func (c *updateService) GetUpdates(ctx context.Context, req *model.GetUpdatesRequest) (<-chan *model.GetUpdatesResponse, <-chan error) {
	protoReq := &v2.GetUpdatesRequest{
		BeginExclusive: req.BeginExclusive,
		Filter:         transactionFilterToProto(req.Filter),
		UpdateFormat:   updateFormatToProto(req.UpdateFormat),
		Verbose:        req.Verbose,
	}

	if req.EndInclusive != nil {
		protoReq.EndInclusive = req.EndInclusive
	}

	stream, err := c.client.GetUpdates(ctx, protoReq)
	if err != nil {
		errCh := make(chan error, 1)
		errCh <- err
		close(errCh)
		return nil, errCh
	}

	responseCh := make(chan *model.GetUpdatesResponse)
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

			modelResp := getUpdatesResponseFromProto(resp)
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

func (c *updateService) GetUpdateById(ctx context.Context, req *model.GetUpdateByIDRequest) (*model.GetUpdateResponse, error) {
	protoReq := &v2.GetUpdateByIdRequest{
		UpdateId:     req.UpdateID,
		UpdateFormat: updateFormatToProto(req.UpdateFormat),
	}

	resp, err := c.client.GetUpdateById(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	return getUpdateResponseFromProto(resp), nil
}

func (c *updateService) GetTransactionByID(ctx context.Context, req *model.GetTransactionByIDRequest) (*model.GetTransactionResponse, error) {
	protoReq := &v2.GetTransactionByIdRequest{
		UpdateId:          req.UpdateID,
		RequestingParties: req.RequestingParties,
	}

	resp, err := c.client.GetTransactionById(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	return getTransactionResponseFromProto(resp), nil
}

func (c *updateService) GetTransactionByOffset(ctx context.Context, req *model.GetTransactionByOffsetRequest) (*model.GetTransactionResponse, error) {
	protoReq := &v2.GetTransactionByOffsetRequest{
		Offset:            req.Offset,
		RequestingParties: req.RequestingParties,
	}

	resp, err := c.client.GetTransactionByOffset(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	return getTransactionResponseFromProto(resp), nil
}

func getUpdatesResponseFromProto(pb *v2.GetUpdatesResponse) *model.GetUpdatesResponse {
	if pb == nil {
		return nil
	}

	resp := &model.GetUpdatesResponse{
		Update: &model.Update{},
	}

	switch update := pb.Update.(type) {
	case *v2.GetUpdatesResponse_Transaction:
		if update.Transaction != nil {
			resp.Update.Transaction = transactionFromProto(update.Transaction)
		}
	case *v2.GetUpdatesResponse_Reassignment:
		if update.Reassignment != nil {
			resp.Update.Reassignment = reassignmentFromProto(update.Reassignment)
		}
	case *v2.GetUpdatesResponse_OffsetCheckpoint:
		resp.Update.OffsetCheckpoint = &model.OffsetCheckpoint{
			Offset: update.OffsetCheckpoint.Offset,
		}
	}

	return resp
}

func getUpdateResponseFromProto(pb *v2.GetUpdateResponse) *model.GetUpdateResponse {
	if pb == nil {
		return nil
	}

	return &model.GetUpdateResponse{
		Transaction: transactionFromProto(pb.GetTransaction()),
	}
}

func getTransactionResponseFromProto(pb *v2.GetTransactionResponse) *model.GetTransactionResponse {
	if pb == nil {
		return nil
	}

	return &model.GetTransactionResponse{
		Transaction: transactionFromProto(pb.Transaction),
	}
}

func transactionFromProto(pb *v2.Transaction) *model.Transaction {
	if pb == nil {
		return nil
	}

	tx := &model.Transaction{
		UpdateID:   pb.UpdateId,
		CommandID:  pb.CommandId,
		WorkflowID: pb.WorkflowId,
		Offset:     pb.Offset,
	}

	if pb.EffectiveAt != nil {
		t := pb.EffectiveAt.AsTime()
		tx.EffectiveAt = &t
	}

	for _, event := range pb.Events {
		tx.Events = append(tx.Events, eventFromProto(event))
	}

	return tx
}

func eventFromProto(pb *v2.Event) *model.Event {
	if pb == nil {
		return nil
	}

	event := &model.Event{}

	switch e := pb.Event.(type) {
	case *v2.Event_Created:
		event.Created = createdEventFromProto(e.Created)
	case *v2.Event_Archived:
		event.Archived = archivedEventFromProto(e.Archived)
	case *v2.Event_Exercised:
		event.Exercised = exercisedEventFromProto(e.Exercised)
	}

	return event
}

func reassignmentFromProto(pb *v2.Reassignment) *model.Reassignment {
	if pb == nil {
		return nil
	}

	r := &model.Reassignment{
		UpdateID: pb.UpdateId,
		Offset:   pb.Offset,
	}

	if pb.RecordTime != nil {
		t := pb.RecordTime.AsTime()
		r.SubmittedAt = &t
	}

	for _, event := range pb.Events {
		switch e := event.Event.(type) {
		case *v2.ReassignmentEvent_Unassigned:
			if e.Unassigned != nil {
				r.UnassignID = e.Unassigned.UnassignId
				r.Source = e.Unassigned.Source
				r.Target = e.Unassigned.Target
				r.Counter = int64(e.Unassigned.ReassignmentCounter)
				if e.Unassigned.AssignmentExclusivity != nil {
					t := e.Unassigned.AssignmentExclusivity.AsTime()
					r.Unassigned = &t
				}
			}
		case *v2.ReassignmentEvent_Assigned:
			if e.Assigned != nil {
				if r.UnassignID == "" {
					r.UnassignID = e.Assigned.UnassignId
				}
				if r.Source == "" {
					r.Source = e.Assigned.Source
				}
				if r.Target == "" {
					r.Target = e.Assigned.Target
				}
				if r.Counter == 0 {
					r.Counter = int64(e.Assigned.ReassignmentCounter)
				}
				r.Reassigned = r.SubmittedAt
			}
		}
	}

	return r
}
