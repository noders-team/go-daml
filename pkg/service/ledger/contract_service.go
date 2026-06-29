package ledger

import (
	"context"

	"google.golang.org/grpc"

	"github.com/noders-team/go-daml/pkg/model"
	v2 "github.com/noders-team/go-daml/proto/com/daml/ledger/api/v2"
)

type ContractService interface {
	GetContract(ctx context.Context, req *model.GetContractRequest) (*model.GetContractResponse, error)
}

type contractService struct {
	client v2.ContractServiceClient
}

func NewContractServiceClient(conn *grpc.ClientConn) *contractService {
	client := v2.NewContractServiceClient(conn)
	return &contractService{
		client: client,
	}
}

func (c *contractService) GetContract(ctx context.Context, req *model.GetContractRequest) (*model.GetContractResponse, error) {
	protoReq := &v2.GetContractRequest{
		ContractId:      req.ContractID,
		QueryingParties: req.QueryingParties,
	}

	resp, err := c.client.GetContract(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	return &model.GetContractResponse{
		CreatedEvent: createdEventFromProto(resp.CreatedEvent),
	}, nil
}
