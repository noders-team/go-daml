package ledger

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/noders-team/go-daml/pkg/model"
	"github.com/noders-team/go-daml/proto/com/daml/ledger/api/v2/interactive"
)

type InteractiveSubmissionService interface {
	PrepareSubmission(ctx context.Context, req *model.PrepareSubmissionRequest) (*model.PrepareSubmissionResponse, error)
	ExecuteSubmission(ctx context.Context, req *model.ExecuteSubmissionRequest) (*model.ExecuteSubmissionResponse, error)
	GetPreferredPackageVersion(ctx context.Context, req *model.GetPreferredPackageVersionRequest) (*model.GetPreferredPackageVersionResponse, error)
}

type interactiveSubmissionService struct {
	client interactive.InteractiveSubmissionServiceClient
}

func NewInteractiveSubmissionServiceClient(conn *grpc.ClientConn) InteractiveSubmissionService {
	client := interactive.NewInteractiveSubmissionServiceClient(conn)
	return &interactiveSubmissionService{
		client: client,
	}
}

func (c *interactiveSubmissionService) PrepareSubmission(ctx context.Context, req *model.PrepareSubmissionRequest) (*model.PrepareSubmissionResponse, error) {
	pbReq := prepareSubmissionRequestToProto(req)
	pbResp, err := c.client.PrepareSubmission(ctx, pbReq)
	if err != nil {
		return nil, err
	}
	return prepareSubmissionResponseFromProto(pbResp), nil
}

func (c *interactiveSubmissionService) ExecuteSubmission(ctx context.Context, req *model.ExecuteSubmissionRequest) (*model.ExecuteSubmissionResponse, error) {
	pbReq := executeSubmissionRequestToProto(req)
	_, err := c.client.ExecuteSubmission(ctx, pbReq)
	if err != nil {
		return nil, err
	}
	return &model.ExecuteSubmissionResponse{}, nil
}

func (c *interactiveSubmissionService) GetPreferredPackageVersion(ctx context.Context, req *model.GetPreferredPackageVersionRequest) (*model.GetPreferredPackageVersionResponse, error) {
	pbReq := &interactive.GetPreferredPackageVersionRequest{
		Parties:        req.Parties,
		PackageName:    req.PackageName,
		SynchronizerId: req.SynchronizerID,
	}

	if req.VettingValidAt != nil {
		pbReq.VettingValidAt = timestamppb.New(*req.VettingValidAt)
	}

	pbResp, err := c.client.GetPreferredPackageVersion(ctx, pbReq)
	if err != nil {
		return nil, err
	}

	resp := &model.GetPreferredPackageVersionResponse{}
	if pbResp.PackagePreference != nil {
		resp.SynchronizerID = pbResp.PackagePreference.SynchronizerId
		if pbResp.PackagePreference.PackageReference != nil {
			resp.PackageReference = &model.PackageReference{
				PackageID:      pbResp.PackagePreference.PackageReference.PackageId,
				PackageName:    pbResp.PackagePreference.PackageReference.PackageName,
				PackageVersion: pbResp.PackagePreference.PackageReference.PackageVersion,
			}
		}
	}

	return resp, nil
}
