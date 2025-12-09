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
		pbStore.Store = &topov30.StoreId_Synchronizer_{
			Synchronizer: &topov30.StoreId_Synchronizer{
				Id: store.Value[13:],
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

	return &protov30.SignedTopologyTransaction{
		Transaction: tx.Transaction,
		Signatures:  signatures,
		Proposal:    tx.Proposal,
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

	return &model.SignedTopologyTransaction{
		Transaction: pb.Transaction,
		Signatures:  signatures,
		Proposal:    pb.Proposal,
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
	}

	return pbMapping
}

func signingPublicKeyToProto(key *model.PublicKey) *cryptov30.SigningPublicKey {
	if key == nil {
		return nil
	}
	return &cryptov30.SigningPublicKey{
		Format:    cryptov30.CryptoKeyFormat(key.Format),
		PublicKey: key.Key,
	}
}
