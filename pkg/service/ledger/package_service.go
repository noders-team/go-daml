package ledger

import (
	"context"

	"google.golang.org/grpc"

	v2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/noders-team/go-daml/pkg/model"
)

type PackageService interface {
	ListPackages(ctx context.Context, req *model.ListPackagesRequest) (*model.ListPackagesResponse, error)
	GetPackage(ctx context.Context, req *model.GetPackageRequest) (*model.GetPackageResponse, error)
	GetPackageStatus(ctx context.Context, req *model.GetPackageStatusRequest) (*model.GetPackageStatusResponse, error)
	ListVettedPackages(ctx context.Context, req *model.ListVettedPackagesRequest) (*model.ListVettedPackagesResponse, error)
}

type packageService struct {
	client v2.PackageServiceClient
}

func NewPackageServiceClient(conn *grpc.ClientConn) *packageService {
	client := v2.NewPackageServiceClient(conn)
	return &packageService{
		client: client,
	}
}

func (c *packageService) ListPackages(ctx context.Context, req *model.ListPackagesRequest) (*model.ListPackagesResponse, error) {
	protoReq := &v2.ListPackagesRequest{}

	resp, err := c.client.ListPackages(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	return &model.ListPackagesResponse{
		PackageIDs: resp.PackageIds,
	}, nil
}

func (c *packageService) GetPackage(ctx context.Context, req *model.GetPackageRequest) (*model.GetPackageResponse, error) {
	protoReq := &v2.GetPackageRequest{
		PackageId: req.PackageID,
	}

	resp, err := c.client.GetPackage(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	return &model.GetPackageResponse{
		ArchivePayload: resp.ArchivePayload,
		HashFunction:   hashFunctionFromProto(resp.HashFunction),
		Hash:           resp.Hash,
	}, nil
}

func (c *packageService) GetPackageStatus(ctx context.Context, req *model.GetPackageStatusRequest) (*model.GetPackageStatusResponse, error) {
	protoReq := &v2.GetPackageStatusRequest{
		PackageId: req.PackageID,
	}

	resp, err := c.client.GetPackageStatus(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	return &model.GetPackageStatusResponse{
		PackageStatus: packageStatusFromProto(resp.PackageStatus),
	}, nil
}

func hashFunctionFromProto(hf v2.HashFunction) model.HashFunction {
	switch hf {
	case v2.HashFunction_HASH_FUNCTION_SHA256:
		return model.HashFunctionSHA256
	default:
		return model.HashFunctionSHA256
	}
}

func (c *packageService) ListVettedPackages(ctx context.Context, req *model.ListVettedPackagesRequest) (*model.ListVettedPackagesResponse, error) {
	protoReq := listVettedPackagesRequestToProto(req)

	resp, err := c.client.ListVettedPackages(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	return listVettedPackagesResponseFromProto(resp), nil
}

func packageStatusFromProto(ps v2.PackageStatus) model.PackageStatus {
	switch ps {
	case v2.PackageStatus_PACKAGE_STATUS_REGISTERED:
		return model.PackageStatusRegistered
	default:
		return model.PackageStatusUnknown
	}
}

func listVettedPackagesRequestToProto(req *model.ListVettedPackagesRequest) *v2.ListVettedPackagesRequest {
	if req == nil {
		return &v2.ListVettedPackagesRequest{}
	}

	protoReq := &v2.ListVettedPackagesRequest{
		PageToken: req.PageToken,
		PageSize:  req.PageSize,
	}

	if req.PackageMetadataFilter != nil {
		protoReq.PackageMetadataFilter = &v2.PackageMetadataFilter{
			PackageIds:          req.PackageMetadataFilter.PackageIDs,
			PackageNamePrefixes: req.PackageMetadataFilter.PackageNamePrefixes,
		}
	}

	if req.TopologyStateFilter != nil {
		protoReq.TopologyStateFilter = &v2.TopologyStateFilter{
			ParticipantIds:  req.TopologyStateFilter.ParticipantIDs,
			SynchronizerIds: req.TopologyStateFilter.SynchronizerIDs,
		}
	}

	return protoReq
}

func listVettedPackagesResponseFromProto(pb *v2.ListVettedPackagesResponse) *model.ListVettedPackagesResponse {
	if pb == nil {
		return nil
	}

	vettedPackages := make([]*model.VettedPackages, len(pb.VettedPackages))
	for i, vp := range pb.VettedPackages {
		vettedPackages[i] = vettedPackagesFromProto(vp)
	}

	return &model.ListVettedPackagesResponse{
		VettedPackages: vettedPackages,
		NextPageToken:  pb.NextPageToken,
	}
}

func vettedPackagesFromProto(pb *v2.VettedPackages) *model.VettedPackages {
	if pb == nil {
		return nil
	}

	packages := make([]*model.VettedPackage, len(pb.Packages))
	for i, pkg := range pb.Packages {
		packages[i] = vettedPackageFromProto(pkg)
	}

	return &model.VettedPackages{
		Packages:       packages,
		ParticipantID:  pb.ParticipantId,
		SynchronizerID: pb.SynchronizerId,
		TopologySerial: pb.TopologySerial,
	}
}

func vettedPackageFromProto(pb *v2.VettedPackage) *model.VettedPackage {
	if pb == nil {
		return nil
	}

	vp := &model.VettedPackage{
		PackageID:      pb.PackageId,
		PackageName:    pb.PackageName,
		PackageVersion: pb.PackageVersion,
	}

	if pb.ValidFromInclusive != nil {
		t := pb.ValidFromInclusive.AsTime()
		vp.ValidFromInclusive = &t
	}

	if pb.ValidUntilExclusive != nil {
		t := pb.ValidUntilExclusive.AsTime()
		vp.ValidUntilExclusive = &t
	}

	return vp
}
