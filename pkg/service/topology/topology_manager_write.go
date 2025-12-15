package topology

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"

	cryptov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/crypto/v30"
	protov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/protocol/v30"
	topov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/topology/admin/v30"
	"github.com/noders-team/go-daml/pkg/model"
)

type TopologyManagerWrite interface {
	Authorize(ctx context.Context, req *model.AuthorizeRequest) (*model.AuthorizeResponse, error)
	AddTransactions(ctx context.Context, req *model.AddTransactionsRequest) (*model.AddTransactionsResponse, error)
	SignTransactions(ctx context.Context, req *model.SignTransactionsRequest) (*model.SignTransactionsResponse, error)
	GenerateTransactions(ctx context.Context, req *model.GenerateTransactionsRequest) (*model.GenerateTransactionsResponse, error)
	CreateTemporaryTopologyStore(ctx context.Context, req *model.CreateTemporaryTopologyStoreRequest) (*model.CreateTemporaryTopologyStoreResponse, error)
	DropTemporaryTopologyStore(ctx context.Context, req *model.DropTemporaryTopologyStoreRequest) (*model.DropTemporaryTopologyStoreResponse, error)
}

type topologyManagerWrite struct {
	client topov30.TopologyManagerWriteServiceClient
}

func NewTopologyManagerWriteClient(conn *grpc.ClientConn) *topologyManagerWrite {
	client := topov30.NewTopologyManagerWriteServiceClient(conn)
	return &topologyManagerWrite{
		client: client,
	}
}

func (c *topologyManagerWrite) Authorize(ctx context.Context, req *model.AuthorizeRequest) (*model.AuthorizeResponse, error) {
	protoReq := authorizeRequestToProto(req)

	resp, err := c.client.Authorize(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	return authorizeResponseFromProto(resp), nil
}

func (c *topologyManagerWrite) AddTransactions(ctx context.Context, req *model.AddTransactionsRequest) (*model.AddTransactionsResponse, error) {
	protoReq := addTransactionsRequestToProto(req)

	resp, err := c.client.AddTransactions(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	return addTransactionsResponseFromProto(resp), nil
}

func authorizeRequestToProto(req *model.AuthorizeRequest) *topov30.AuthorizeRequest {
	if req == nil {
		return nil
	}

	protoReq := &topov30.AuthorizeRequest{
		MustFullyAuthorize: req.MustFullyAuthorize,
		ForceChanges:       forceFlagsToProto(req.ForceChanges),
		SignedBy:           req.SignedBy,
		Store:              storeIDToProto(req.Store),
	}

	if req.WaitToBecomeEffective != nil {
		protoReq.WaitToBecomeEffective = durationpb.New(*req.WaitToBecomeEffective)
	}

	if req.Proposal != nil {
		protoReq.Type = &topov30.AuthorizeRequest_Proposal_{
			Proposal: &topov30.AuthorizeRequest_Proposal{
				Change:  operationToProto(req.Proposal.Operation),
				Mapping: topologyMappingToProto(req.Proposal.Mapping),
				Serial:  req.Proposal.Serial,
			},
		}
	} else if req.TransactionHash != "" {
		protoReq.Type = &topov30.AuthorizeRequest_TransactionHash{
			TransactionHash: req.TransactionHash,
		}
	}

	return protoReq
}

func authorizeResponseFromProto(pb *topov30.AuthorizeResponse) *model.AuthorizeResponse {
	if pb == nil {
		return nil
	}

	return &model.AuthorizeResponse{
		Transaction: signedTopologyTransactionFromProto(pb.Transaction),
	}
}

func addTransactionsRequestToProto(req *model.AddTransactionsRequest) *topov30.AddTransactionsRequest {
	if req == nil {
		return nil
	}

	protoReq := &topov30.AddTransactionsRequest{
		Transactions: signedTopologyTransactionsToProto(req.Transactions),
		ForceChanges: forceFlagsToProto(req.ForceChanges),
		Store:        storeIDToProto(req.Store),
	}

	if req.WaitToBecomeEffective != nil {
		protoReq.WaitToBecomeEffective = durationpb.New(*req.WaitToBecomeEffective)
	}

	return protoReq
}

func addTransactionsResponseFromProto(pb *topov30.AddTransactionsResponse) *model.AddTransactionsResponse {
	if pb == nil {
		return nil
	}

	return &model.AddTransactionsResponse{}
}

func forceFlagsToProto(flags []model.ForceFlag) []topov30.ForceFlag {
	result := make([]topov30.ForceFlag, len(flags))
	for i, flag := range flags {
		result[i] = forceFlagToProto(flag)
	}
	return result
}

func forceFlagToProto(flag model.ForceFlag) topov30.ForceFlag {
	switch flag {
	case model.ForceFlagAlienMember:
		return topov30.ForceFlag_FORCE_FLAG_ALIEN_MEMBER
	case model.ForceFlagLedgerTimeRecordTimeToleranceIncrease:
		return topov30.ForceFlag_FORCE_FLAG_LEDGER_TIME_RECORD_TIME_TOLERANCE_INCREASE
	default:
		return topov30.ForceFlag_FORCE_FLAG_UNSPECIFIED
	}
}

func storeIDToProto(store *model.StoreID) *topov30.StoreId {
	if store == nil {
		return nil
	}

	pbStore := &topov30.StoreId{}
	if store.Value == "authorized" {
		pbStore.Store = &topov30.StoreId_Authorized_{
			Authorized: &topov30.StoreId_Authorized{},
		}
	} else if strings.HasPrefix(store.Value, "synchronizer:") {
		pbStore.Store = &topov30.StoreId_Synchronizer{
			Synchronizer: &topov30.Synchronizer{
				Kind: &topov30.Synchronizer_Id{
					Id: store.Value[13:],
				},
			},
		}
	} else if strings.HasPrefix(store.Value, "temporary:") {
		pbStore.Store = &topov30.StoreId_Temporary_{
			Temporary: &topov30.StoreId_Temporary{
				Name: store.Value[10:],
			},
		}
	}

	return pbStore
}

func signedTopologyTransactionToProto(tx *model.SignedTopologyTransaction) *protov30.SignedTopologyTransaction {
	if tx == nil {
		return nil
	}

	signatures := make([]*cryptov30.Signature, len(tx.Signatures))
	for i, sig := range tx.Signatures {
		signatures[i] = &cryptov30.Signature{
			SignedBy:  sig.SignedBy,
			Signature: sig.Signature,
			Format:    cryptov30.SignatureFormat(sig.SignatureFormat),
		}
	}

	multiTxSigs := make([]*protov30.MultiTransactionSignatures, len(tx.MultiTransactionSignatures))
	for i, mts := range tx.MultiTransactionSignatures {
		sigs := make([]*cryptov30.Signature, len(mts.Signatures))
		for j, sig := range mts.Signatures {
			sigs[j] = &cryptov30.Signature{
				SignedBy:  sig.SignedBy,
				Signature: sig.Signature,
				Format:    cryptov30.SignatureFormat(sig.SignatureFormat),
			}
		}
		multiTxSigs[i] = &protov30.MultiTransactionSignatures{
			TransactionHashes: mts.TransactionHashes,
			Signatures:        sigs,
		}
	}

	return &protov30.SignedTopologyTransaction{
		Transaction:                tx.Transaction,
		Signatures:                 signatures,
		MultiTransactionSignatures: multiTxSigs,
		Proposal:                   tx.Proposal,
	}
}

func signedTopologyTransactionsToProto(txs []*model.SignedTopologyTransaction) []*protov30.SignedTopologyTransaction {
	result := make([]*protov30.SignedTopologyTransaction, len(txs))
	for i, tx := range txs {
		result[i] = signedTopologyTransactionToProto(tx)
	}
	return result
}

func signedTopologyTransactionFromProto(pb *protov30.SignedTopologyTransaction) *model.SignedTopologyTransaction {
	if pb == nil {
		return nil
	}

	signatures := make([]model.TopologyTransactionSignature, len(pb.Signatures))
	for i, sig := range pb.Signatures {
		signatures[i] = model.TopologyTransactionSignature{
			SignedBy:        sig.SignedBy,
			Signature:       sig.Signature,
			SignatureFormat: int32(sig.Format),
		}
	}

	multiTxSigs := make([]*model.MultiTransactionSignatures, len(pb.MultiTransactionSignatures))
	for i, mts := range pb.MultiTransactionSignatures {
		sigs := make([]model.TopologyTransactionSignature, len(mts.Signatures))
		for j, sig := range mts.Signatures {
			sigs[j] = model.TopologyTransactionSignature{
				SignedBy:        sig.SignedBy,
				Signature:       sig.Signature,
				SignatureFormat: int32(sig.Format),
			}
		}
		multiTxSigs[i] = &model.MultiTransactionSignatures{
			TransactionHashes: mts.TransactionHashes,
			Signatures:        sigs,
		}
	}

	return &model.SignedTopologyTransaction{
		Transaction:                pb.Transaction,
		Signatures:                 signatures,
		MultiTransactionSignatures: multiTxSigs,
		Proposal:                   pb.Proposal,
	}
}

func signedTopologyTransactionsFromProto(pbs []*protov30.SignedTopologyTransaction) []*model.SignedTopologyTransaction {
	result := make([]*model.SignedTopologyTransaction, len(pbs))
	for i, pb := range pbs {
		result[i] = signedTopologyTransactionFromProto(pb)
	}
	return result
}

func topologyMappingToProto(mapping model.TopologyMapping) *protov30.TopologyMapping {
	if mapping == nil {
		return nil
	}

	pbMapping := &protov30.TopologyMapping{}

	switch m := mapping.(type) {
	case *model.NamespaceDelegationMapping:
		pbMapping.Mapping = &protov30.TopologyMapping_NamespaceDelegation{
			NamespaceDelegation: &protov30.NamespaceDelegation{
				Namespace:        m.Namespace,
				TargetKey:        signingPublicKeyToProto(&m.TargetKey),
				IsRootDelegation: m.IsRootDelegation,
			},
		}
	case *model.PartyToKeyMapping:
		keys := make([]*cryptov30.SigningPublicKey, len(m.SigningKeys))
		for i, k := range m.SigningKeys {
			keys[i] = signingPublicKeyToProto(&k)
		}
		pbMapping.Mapping = &protov30.TopologyMapping_PartyToKeyMapping{
			PartyToKeyMapping: &protov30.PartyToKeyMapping{
				Party:       m.Party,
				Threshold:   m.Threshold,
				SigningKeys: keys,
			},
		}
	case *model.PartyToParticipantMapping:
		participants := make([]*protov30.PartyToParticipant_HostingParticipant, len(m.Participants))
		for i, p := range m.Participants {
			participants[i] = &protov30.PartyToParticipant_HostingParticipant{
				ParticipantUid: p.ParticipantUID,
				Permission:     participantPermissionToProto(p.Permission),
			}
		}
		pbMapping.Mapping = &protov30.TopologyMapping_PartyToParticipant{
			PartyToParticipant: &protov30.PartyToParticipant{
				Party:        m.Party,
				Threshold:    m.Threshold,
				Participants: participants,
			},
		}
	}

	return pbMapping
}

func signingPublicKeyToProto(key *model.PublicKey) *cryptov30.SigningPublicKey {
	if key == nil {
		return nil
	}
	usage := make([]cryptov30.SigningKeyUsage, len(key.Usage))
	for i, u := range key.Usage {
		usage[i] = cryptov30.SigningKeyUsage(u)
	}
	return &cryptov30.SigningPublicKey{
		Format:    cryptov30.CryptoKeyFormat(key.Format),
		PublicKey: key.Key,
		Scheme:    cryptov30.SigningKeyScheme(key.Scheme),
		KeySpec:   cryptov30.SigningKeySpec(key.KeySpec),
		Usage:     usage,
	}
}

func participantPermissionToProto(permission model.ParticipantPermission) protov30.Enums_ParticipantPermission {
	switch permission {
	case model.ParticipantPermissionConfirmation:
		return protov30.Enums_PARTICIPANT_PERMISSION_CONFIRMATION
	case model.ParticipantPermissionObservation:
		return protov30.Enums_PARTICIPANT_PERMISSION_OBSERVATION
	default:
		return protov30.Enums_PARTICIPANT_PERMISSION_SUBMISSION
	}
}

func (c *topologyManagerWrite) SignTransactions(ctx context.Context, req *model.SignTransactionsRequest) (*model.SignTransactionsResponse, error) {
	protoReq := signTransactionsRequestToProto(req)

	resp, err := c.client.SignTransactions(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	return signTransactionsResponseFromProto(resp), nil
}

func (c *topologyManagerWrite) GenerateTransactions(ctx context.Context, req *model.GenerateTransactionsRequest) (*model.GenerateTransactionsResponse, error) {
	protoReq := generateTransactionsRequestToProto(req)

	resp, err := c.client.GenerateTransactions(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	return generateTransactionsResponseFromProto(resp), nil
}

func (c *topologyManagerWrite) CreateTemporaryTopologyStore(ctx context.Context, req *model.CreateTemporaryTopologyStoreRequest) (*model.CreateTemporaryTopologyStoreResponse, error) {
	protoReq := &topov30.CreateTemporaryTopologyStoreRequest{
		Name:            req.Name,
		ProtocolVersion: req.ProtocolVersion,
	}

	resp, err := c.client.CreateTemporaryTopologyStore(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	return &model.CreateTemporaryTopologyStoreResponse{
		StoreID: &model.StoreID{Value: "temporary:" + resp.StoreId.Name},
	}, nil
}

func (c *topologyManagerWrite) DropTemporaryTopologyStore(ctx context.Context, req *model.DropTemporaryTopologyStoreRequest) (*model.DropTemporaryTopologyStoreResponse, error) {
	protoReq := &topov30.DropTemporaryTopologyStoreRequest{
		StoreId: &topov30.StoreId_Temporary{
			Name: req.StoreID.Value[10:],
		},
	}

	_, err := c.client.DropTemporaryTopologyStore(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	return &model.DropTemporaryTopologyStoreResponse{}, nil
}

func signTransactionsRequestToProto(req *model.SignTransactionsRequest) *topov30.SignTransactionsRequest {
	if req == nil {
		return nil
	}

	return &topov30.SignTransactionsRequest{
		Transactions: signedTopologyTransactionsToProto(req.Transactions),
		SignedBy:     req.SignedBy,
		Store:        storeIDToProto(req.Store),
		ForceFlags:   forceFlagsToProto(req.ForceFlags),
	}
}

func signTransactionsResponseFromProto(pb *topov30.SignTransactionsResponse) *model.SignTransactionsResponse {
	if pb == nil {
		return nil
	}

	return &model.SignTransactionsResponse{
		Transactions: signedTopologyTransactionsFromProto(pb.Transactions),
	}
}

func generateTransactionsRequestToProto(req *model.GenerateTransactionsRequest) *topov30.GenerateTransactionsRequest {
	if req == nil {
		return nil
	}

	proposals := make([]*topov30.GenerateTransactionsRequest_Proposal, len(req.Proposals))
	for i, p := range req.Proposals {
		proposals[i] = &topov30.GenerateTransactionsRequest_Proposal{
			Operation: operationToProto(p.Operation),
			Serial:    p.Serial,
			Mapping:   topologyMappingToProto(p.Mapping),
			Store:     storeIDToProto(p.Store),
		}
	}

	return &topov30.GenerateTransactionsRequest{
		Proposals: proposals,
	}
}

func generateTransactionsResponseFromProto(pb *topov30.GenerateTransactionsResponse) *model.GenerateTransactionsResponse {
	if pb == nil {
		return nil
	}

	genTxs := make([]*model.GeneratedTransaction, len(pb.GeneratedTransactions))
	for i, tx := range pb.GeneratedTransactions {
		genTxs[i] = &model.GeneratedTransaction{
			SerializedTransaction: tx.SerializedTransaction,
			TransactionHash:       tx.TransactionHash,
		}
	}

	return &model.GenerateTransactionsResponse{
		GeneratedTransactions: genTxs,
	}
}
