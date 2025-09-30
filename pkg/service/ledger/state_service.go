package ledger

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/grpc"

	v2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/noders-team/go-daml/pkg/model"
)

type StateService interface {
	GetActiveContracts(ctx context.Context, req *model.GetActiveContractsRequest) (<-chan *model.GetActiveContractsResponse, <-chan error)
	GetConnectedSynchronizers(ctx context.Context, req *model.GetConnectedSynchronizersRequest) (*model.GetConnectedSynchronizersResponse, error)
	GetLedgerEnd(ctx context.Context, req *model.GetLedgerEndRequest) (*model.GetLedgerEndResponse, error)
	GetLatestPrunedOffsets(ctx context.Context, req *model.GetLatestPrunedOffsetsRequest) (*model.GetLatestPrunedOffsetsResponse, error)
}

type stateService struct {
	client v2.StateServiceClient
}

func NewStateServiceClient(conn *grpc.ClientConn) *stateService {
	client := v2.NewStateServiceClient(conn)
	return &stateService{
		client: client,
	}
}

func (c *stateService) GetActiveContracts(ctx context.Context, req *model.GetActiveContractsRequest) (<-chan *model.GetActiveContractsResponse, <-chan error) {
	protoReq := &v2.GetActiveContractsRequest{
		Filter:         transactionFilterToProto(req.Filter),
		Verbose:        req.Verbose,
		ActiveAtOffset: req.ActiveAtOffset,
		EventFormat:    eventFormatToProto(req.EventFormat),
	}

	stream, err := c.client.GetActiveContracts(ctx, protoReq)
	if err != nil {
		errCh := make(chan error, 1)
		errCh <- err
		close(errCh)
		return nil, errCh
	}

	responseCh := make(chan *model.GetActiveContractsResponse)
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

			modelResp := getActiveContractsResponseFromProto(resp)
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

func (c *stateService) GetConnectedSynchronizers(ctx context.Context, req *model.GetConnectedSynchronizersRequest) (*model.GetConnectedSynchronizersResponse, error) {
	protoReq := &v2.GetConnectedSynchronizersRequest{}

	resp, err := c.client.GetConnectedSynchronizers(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	return getConnectedSynchronizersResponseFromProto(resp), nil
}

func (c *stateService) GetLedgerEnd(ctx context.Context, req *model.GetLedgerEndRequest) (*model.GetLedgerEndResponse, error) {
	protoReq := &v2.GetLedgerEndRequest{}

	resp, err := c.client.GetLedgerEnd(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	return &model.GetLedgerEndResponse{
		Offset: resp.Offset,
	}, nil
}

func (c *stateService) GetLatestPrunedOffsets(ctx context.Context, req *model.GetLatestPrunedOffsetsRequest) (*model.GetLatestPrunedOffsetsResponse, error) {
	protoReq := &v2.GetLatestPrunedOffsetsRequest{}

	resp, err := c.client.GetLatestPrunedOffsets(ctx, protoReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest pruned offsets: %w", err)
	}

	return &model.GetLatestPrunedOffsetsResponse{
		ParticipantPrunedUpToInclusive:          resp.ParticipantPrunedUpToInclusive,
		AllDivulgedContractsPrunedUpToInclusive: resp.AllDivulgedContractsPrunedUpToInclusive,
	}, nil
}

func getActiveContractsResponseFromProto(pb *v2.GetActiveContractsResponse) *model.GetActiveContractsResponse {
	if pb == nil {
		return nil
	}

	resp := &model.GetActiveContractsResponse{
		WorkflowID: pb.WorkflowId,
	}

	// Handle the oneof contract_entry field
	switch entry := pb.ContractEntry.(type) {
	case *v2.GetActiveContractsResponse_ActiveContract:
		if entry.ActiveContract != nil && entry.ActiveContract.CreatedEvent != nil {
			resp.ActiveContracts = append(resp.ActiveContracts, createdEventFromProto(entry.ActiveContract.CreatedEvent))
		}
	case *v2.GetActiveContractsResponse_IncompleteUnassigned:
		if entry.IncompleteUnassigned != nil {
			resp.IncompleteUnassigned = &model.IncompleteUnassigned{
				CreatedEvent:    createdEventFromProto(entry.IncompleteUnassigned.CreatedEvent),
				UnassignedEvent: unassignedEventFromProto(entry.IncompleteUnassigned.UnassignedEvent),
			}
			// Set offset from the unassigned event if available
			if entry.IncompleteUnassigned.UnassignedEvent != nil && entry.IncompleteUnassigned.UnassignedEvent.Offset > 0 {
				resp.Offset = entry.IncompleteUnassigned.UnassignedEvent.Offset
			}
		}
	case *v2.GetActiveContractsResponse_IncompleteAssigned:
		if entry.IncompleteAssigned != nil {
			resp.IncompleteAssigned = &model.IncompleteAssigned{
				AssignedEvent: assignedEventFromProto(entry.IncompleteAssigned.AssignedEvent),
			}
		}
	}

	return resp
}

func getConnectedSynchronizersResponseFromProto(pb *v2.GetConnectedSynchronizersResponse) *model.GetConnectedSynchronizersResponse {
	if pb == nil {
		return nil
	}

	resp := &model.GetConnectedSynchronizersResponse{}

	for _, sync := range pb.ConnectedSynchronizers {
		resp.ConnectedSynchronizers = append(resp.ConnectedSynchronizers, connectedSynchronizerFromProto(sync))
	}

	return resp
}

func connectedSynchronizerFromProto(pb *v2.GetConnectedSynchronizersResponse_ConnectedSynchronizer) *model.ConnectedSynchronizer {
	if pb == nil {
		return nil
	}

	return &model.ConnectedSynchronizer{
		SynchronizerID:        pb.SynchronizerId,
		ParticipantPermission: participantPermissionFromProto(pb.Permission),
	}
}

func participantPermissionFromProto(pp v2.ParticipantPermission) model.ParticipantPermission {
	switch pp {
	case v2.ParticipantPermission_PARTICIPANT_PERMISSION_CONFIRMATION:
		return model.ParticipantPermissionConfirmation
	case v2.ParticipantPermission_PARTICIPANT_PERMISSION_OBSERVATION:
		return model.ParticipantPermissionObservation
	default:
		return model.ParticipantPermissionSubmission
	}
}
