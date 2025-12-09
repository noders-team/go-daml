package topology

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	cryptov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/crypto/v30"
	protov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/protocol/v30"
	topov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/topology/admin/v30"
	"github.com/noders-team/go-daml/pkg/model"
)

type TopologyManagerRead interface {
	ListNamespaceDelegation(ctx context.Context, req *model.ListNamespaceDelegationRequest) (*model.ListNamespaceDelegationResponse, error)
	ListPartyToKeyMapping(ctx context.Context, req *model.ListPartyToKeyMappingRequest) (*model.ListPartyToKeyMappingResponse, error)
	ListPartyToParticipant(ctx context.Context, req *model.ListPartyToParticipantRequest) (*model.ListPartyToParticipantResponse, error)
}

type topologyManagerRead struct {
	client topov30.TopologyManagerReadServiceClient
}

func NewTopologyManagerReadClient(conn *grpc.ClientConn) *topologyManagerRead {
	client := topov30.NewTopologyManagerReadServiceClient(conn)
	return &topologyManagerRead{
		client: client,
	}
}

func (c *topologyManagerRead) ListNamespaceDelegation(ctx context.Context, req *model.ListNamespaceDelegationRequest) (*model.ListNamespaceDelegationResponse, error) {
	protoReq := listNamespaceDelegationRequestToProto(req)

	resp, err := c.client.ListNamespaceDelegation(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	return listNamespaceDelegationResponseFromProto(resp), nil
}

func (c *topologyManagerRead) ListPartyToKeyMapping(ctx context.Context, req *model.ListPartyToKeyMappingRequest) (*model.ListPartyToKeyMappingResponse, error) {
	protoReq := listPartyToKeyMappingRequestToProto(req)

	resp, err := c.client.ListPartyToKeyMapping(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	return listPartyToKeyMappingResponseFromProto(resp), nil
}

func (c *topologyManagerRead) ListPartyToParticipant(ctx context.Context, req *model.ListPartyToParticipantRequest) (*model.ListPartyToParticipantResponse, error) {
	protoReq := listPartyToParticipantRequestToProto(req)

	resp, err := c.client.ListPartyToParticipant(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	return listPartyToParticipantResponseFromProto(resp), nil
}

func listNamespaceDelegationRequestToProto(req *model.ListNamespaceDelegationRequest) *topov30.ListNamespaceDelegationRequest {
	if req == nil {
		return nil
	}

	return &topov30.ListNamespaceDelegationRequest{
		BaseQuery:                  baseQueryToProto(req.BaseQuery),
		FilterNamespace:            req.FilterNamespace,
		FilterTargetKeyFingerprint: req.FilterTargetKeyFingerprint,
	}
}

func listNamespaceDelegationResponseFromProto(pb *topov30.ListNamespaceDelegationResponse) *model.ListNamespaceDelegationResponse {
	if pb == nil {
		return nil
	}

	results := make([]*model.NamespaceDelegationResult, len(pb.Results))
	for i, r := range pb.Results {
		results[i] = namespaceDelegationResultFromProto(r)
	}

	return &model.ListNamespaceDelegationResponse{
		Results: results,
	}
}

func listPartyToKeyMappingRequestToProto(req *model.ListPartyToKeyMappingRequest) *topov30.ListPartyToKeyMappingRequest {
	if req == nil {
		return nil
	}

	return &topov30.ListPartyToKeyMappingRequest{
		BaseQuery:   baseQueryToProto(req.BaseQuery),
		FilterParty: req.FilterParty,
	}
}

func listPartyToKeyMappingResponseFromProto(pb *topov30.ListPartyToKeyMappingResponse) *model.ListPartyToKeyMappingResponse {
	if pb == nil {
		return nil
	}

	results := make([]*model.PartyToKeyMappingResult, len(pb.Results))
	for i, r := range pb.Results {
		results[i] = partyToKeyMappingResultFromProto(r)
	}

	return &model.ListPartyToKeyMappingResponse{
		Results: results,
	}
}

func listPartyToParticipantRequestToProto(req *model.ListPartyToParticipantRequest) *topov30.ListPartyToParticipantRequest {
	if req == nil {
		return nil
	}

	return &topov30.ListPartyToParticipantRequest{
		BaseQuery:         baseQueryToProto(req.BaseQuery),
		FilterParty:       req.FilterParty,
		FilterParticipant: req.FilterParticipant,
	}
}

func listPartyToParticipantResponseFromProto(pb *topov30.ListPartyToParticipantResponse) *model.ListPartyToParticipantResponse {
	if pb == nil {
		return nil
	}

	results := make([]*model.PartyToParticipantResult, len(pb.Results))
	for i, r := range pb.Results {
		results[i] = partyToParticipantResultFromProto(r)
	}

	return &model.ListPartyToParticipantResponse{
		Results: results,
	}
}

func baseQueryToProto(query *model.BaseQuery) *topov30.BaseQuery {
	if query == nil {
		return nil
	}

	pbQuery := &topov30.BaseQuery{
		Store:           storeIDToProto(query.Store),
		Proposals:       query.Proposals,
		Operation:       operationToProto(query.Operation),
		FilterSignedKey: query.FilterSignedKey,
	}

	if query.ProtocolVersion != nil {
		pbQuery.ProtocolVersion = query.ProtocolVersion
	}

	if query.TimeQuery != nil {
		if query.TimeQuery.Serial != nil {
			pbQuery.TimeQuery = &topov30.BaseQuery_Snapshot{
				Snapshot: timestamppb.New(time.Unix(*query.TimeQuery.Serial, 0)),
			}
		} else if query.TimeQuery.Range != nil {
			pbRange := &topov30.BaseQuery_TimeRange{}
			if query.TimeQuery.Range.From != nil {
				pbRange.From = timestamppb.New(*query.TimeQuery.Range.From)
			}
			if query.TimeQuery.Range.Until != nil {
				pbRange.Until = timestamppb.New(*query.TimeQuery.Range.Until)
			}
			pbQuery.TimeQuery = &topov30.BaseQuery_Range{
				Range: pbRange,
			}
		}
	}

	return pbQuery
}

func operationToProto(op model.Operation) protov30.Enums_TopologyChangeOp {
	switch op {
	case model.OperationAddReplace:
		return protov30.Enums_TOPOLOGY_CHANGE_OP_ADD_REPLACE
	case model.OperationRemove:
		return protov30.Enums_TOPOLOGY_CHANGE_OP_REMOVE
	default:
		return protov30.Enums_TOPOLOGY_CHANGE_OP_UNSPECIFIED
	}
}

func namespaceDelegationResultFromProto(pb *topov30.ListNamespaceDelegationResponse_Result) *model.NamespaceDelegationResult {
	if pb == nil {
		return nil
	}

	return &model.NamespaceDelegationResult{
		Context: baseResultFromProto(pb.Context),
		Item:    namespaceDelegationFromProto(pb.Item),
	}
}

func partyToKeyMappingResultFromProto(pb *topov30.ListPartyToKeyMappingResponse_Result) *model.PartyToKeyMappingResult {
	if pb == nil {
		return nil
	}

	return &model.PartyToKeyMappingResult{
		Context: baseResultFromProto(pb.Context),
		Item:    partyToKeyMappingFromProto(pb.Item),
	}
}

func partyToParticipantResultFromProto(pb *topov30.ListPartyToParticipantResponse_Result) *model.PartyToParticipantResult {
	if pb == nil {
		return nil
	}

	return &model.PartyToParticipantResult{
		Context: baseResultFromProto(pb.Context),
		Item:    partyToParticipantMappingFromProto(pb.Item),
	}
}

func baseResultFromProto(pb *topov30.BaseResult) *model.BaseResult {
	if pb == nil {
		return nil
	}

	result := &model.BaseResult{
		Store:                storeIDFromProto(pb.Store),
		Operation:            operationFromProto(pb.Operation),
		TransactionHash:      pb.TransactionHash,
		Serial:               pb.Serial,
		SignedByFingerprints: pb.SignedByFingerprints,
	}

	if pb.Sequenced != nil {
		t := pb.Sequenced.AsTime()
		result.Sequenced = &t
	}

	if pb.ValidFrom != nil {
		t := pb.ValidFrom.AsTime()
		result.ValidFrom = &t
	}

	if pb.ValidUntil != nil {
		t := pb.ValidUntil.AsTime()
		result.ValidUntil = &t
	}

	return result
}

func storeIDFromProto(pb *topov30.StoreId) *model.StoreID {
	if pb == nil {
		return nil
	}

	var value string
	switch store := pb.Store.(type) {
	case *topov30.StoreId_Authorized_:
		value = "authorized"
	case *topov30.StoreId_Synchronizer_:
		if store.Synchronizer != nil {
			value = "synchronizer:" + store.Synchronizer.Id
		}
	case *topov30.StoreId_Temporary_:
		if store.Temporary != nil {
			value = "temporary:" + store.Temporary.Name
		}
	}

	return &model.StoreID{Value: value}
}

func operationFromProto(op protov30.Enums_TopologyChangeOp) model.Operation {
	switch op {
	case protov30.Enums_TOPOLOGY_CHANGE_OP_ADD_REPLACE:
		return model.OperationAddReplace
	case protov30.Enums_TOPOLOGY_CHANGE_OP_REMOVE:
		return model.OperationRemove
	default:
		return model.OperationUnspecified
	}
}

func namespaceDelegationFromProto(pb *protov30.NamespaceDelegation) *model.NamespaceDelegationMapping {
	if pb == nil {
		return nil
	}

	return &model.NamespaceDelegationMapping{
		Namespace:        pb.Namespace,
		TargetKey:        signingPublicKeyFromProto(pb.TargetKey),
		IsRootDelegation: pb.IsRootDelegation,
	}
}

func partyToKeyMappingFromProto(pb *protov30.PartyToKeyMapping) *model.PartyToKeyMapping {
	if pb == nil {
		return nil
	}

	keys := make([]model.PublicKey, len(pb.SigningKeys))
	for i, k := range pb.SigningKeys {
		keys[i] = signingPublicKeyFromProto(k)
	}

	return &model.PartyToKeyMapping{
		Party:       pb.Party,
		Threshold:   pb.Threshold,
		SigningKeys: keys,
	}
}

func partyToParticipantMappingFromProto(pb *protov30.PartyToParticipant) *model.PartyToParticipantMapping {
	if pb == nil {
		return nil
	}

	participants := make([]model.HostingParticipant, len(pb.Participants))
	for i, p := range pb.Participants {
		participants[i] = model.HostingParticipant{
			ParticipantUID: p.ParticipantUid,
			Permission:     participantPermissionFromProto(p.Permission),
		}
	}

	return &model.PartyToParticipantMapping{
		Party:        pb.Party,
		Threshold:    pb.Threshold,
		Participants: participants,
	}
}

func participantPermissionFromProto(pp protov30.Enums_ParticipantPermission) model.ParticipantPermission {
	switch pp {
	case protov30.Enums_PARTICIPANT_PERMISSION_CONFIRMATION:
		return model.ParticipantPermissionConfirmation
	case protov30.Enums_PARTICIPANT_PERMISSION_OBSERVATION:
		return model.ParticipantPermissionObservation
	default:
		return model.ParticipantPermissionSubmission
	}
}

func signingPublicKeyFromProto(pb *cryptov30.SigningPublicKey) model.PublicKey {
	if pb == nil {
		return model.PublicKey{}
	}
	return model.PublicKey{
		Format: int32(pb.Format),
		Key:    pb.PublicKey,
		ID:     "",
	}
}
