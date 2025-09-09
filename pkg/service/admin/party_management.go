package admin

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	adminv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/admin"
	"github.com/noders-team/go-daml/pkg/model"
)

type PartyManagement interface {
	GetParticipantID(ctx context.Context) (string, error)
	GetParties(ctx context.Context, parties []string, identityProviderID string) ([]*model.PartyDetails, error)
	ListKnownParties(ctx context.Context, pageToken string, pageSize int32, identityProviderID string) (*model.ListKnownPartiesResponse, error)
	AllocateParty(ctx context.Context, partyIDHint string, localMetadata map[string]string, identityProviderID string) (*model.PartyDetails, error)
	UpdatePartyDetails(ctx context.Context, party *model.PartyDetails) (*model.PartyDetails, error)
	UpdatePartyIdentityProviderID(ctx context.Context, party string, sourceIdentityProviderID string, targetIdentityProviderID string) error
}

type partyManagement struct {
	client adminv2.PartyManagementServiceClient
}

func NewPartyManagementClient(conn *grpc.ClientConn) *partyManagement {
	client := adminv2.NewPartyManagementServiceClient(conn)
	return &partyManagement{
		client: client,
	}
}

func (c *partyManagement) GetParticipantID(ctx context.Context) (string, error) {
	req := &adminv2.GetParticipantIdRequest{}

	resp, err := c.client.GetParticipantId(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to get participant ID: %w", err)
	}

	return resp.ParticipantId, nil
}

func (c *partyManagement) GetParties(ctx context.Context, parties []string, identityProviderID string) ([]*model.PartyDetails, error) {
	req := &adminv2.GetPartiesRequest{
		Parties:            parties,
		IdentityProviderId: identityProviderID,
	}

	resp, err := c.client.GetParties(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get parties: %w", err)
	}

	return partyDetailsFromProtos(resp.PartyDetails), nil
}

func (c *partyManagement) ListKnownParties(ctx context.Context, pageToken string, pageSize int32, identityProviderID string) (*model.ListKnownPartiesResponse, error) {
	req := &adminv2.ListKnownPartiesRequest{
		PageToken:          pageToken,
		PageSize:           pageSize,
		IdentityProviderId: identityProviderID,
	}

	resp, err := c.client.ListKnownParties(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list known parties: %w", err)
	}

	return &model.ListKnownPartiesResponse{
		PartyDetails:  partyDetailsFromProtos(resp.PartyDetails),
		NextPageToken: resp.NextPageToken,
	}, nil
}

func (c *partyManagement) AllocateParty(ctx context.Context, partyIDHint string, localMetadata map[string]string, identityProviderID string) (*model.PartyDetails, error) {
	var metadata *adminv2.ObjectMeta
	if len(localMetadata) > 0 {
		metadata = &adminv2.ObjectMeta{
			Annotations: localMetadata,
		}
	}

	req := &adminv2.AllocatePartyRequest{
		PartyIdHint:        partyIDHint,
		LocalMetadata:      metadata,
		IdentityProviderId: identityProviderID,
	}

	resp, err := c.client.AllocateParty(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to allocate party: %w", err)
	}

	return partyDetailsFromProto(resp.PartyDetails), nil
}

func (c *partyManagement) UpdatePartyDetails(ctx context.Context, party *model.PartyDetails) (*model.PartyDetails, error) {
	req := &adminv2.UpdatePartyDetailsRequest{
		PartyDetails: partyDetailsToProto(party),
	}

	resp, err := c.client.UpdatePartyDetails(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update party details: %w", err)
	}

	return partyDetailsFromProto(resp.PartyDetails), nil
}

func (c *partyManagement) UpdatePartyIdentityProviderID(ctx context.Context, party string, sourceIdentityProviderID string, targetIdentityProviderID string) error {
	req := &adminv2.UpdatePartyIdentityProviderIdRequest{
		Party:                    party,
		SourceIdentityProviderId: sourceIdentityProviderID,
		TargetIdentityProviderId: targetIdentityProviderID,
	}

	_, err := c.client.UpdatePartyIdentityProviderId(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to update party identity provider ID: %w", err)
	}

	return nil
}

func partyDetailsFromProto(pb *adminv2.PartyDetails) *model.PartyDetails {
	if pb == nil {
		return nil
	}
	
	localMetadata := make(map[string]string)
	if pb.LocalMetadata != nil {
		localMetadata = pb.LocalMetadata.Annotations
	}
	
	return &model.PartyDetails{
		Party:              pb.Party,
		IsLocal:            pb.IsLocal,
		LocalMetadata:      localMetadata,
		IdentityProviderID: pb.IdentityProviderId,
	}
}

func partyDetailsToProto(pd *model.PartyDetails) *adminv2.PartyDetails {
	if pd == nil {
		return nil
	}
	
	var metadata *adminv2.ObjectMeta
	if len(pd.LocalMetadata) > 0 {
		metadata = &adminv2.ObjectMeta{
			Annotations: pd.LocalMetadata,
		}
	}
	
	return &adminv2.PartyDetails{
		Party:              pd.Party,
		IsLocal:            pd.IsLocal,
		LocalMetadata:      metadata,
		IdentityProviderId: pd.IdentityProviderID,
	}
}

func partyDetailsFromProtos(pbs []*adminv2.PartyDetails) []*model.PartyDetails {
	result := make([]*model.PartyDetails, len(pbs))
	for i, pb := range pbs {
		result[i] = partyDetailsFromProto(pb)
	}
	return result
}