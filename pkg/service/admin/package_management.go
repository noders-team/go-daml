package admin

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	v2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	adminv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/admin"
	"github.com/noders-team/go-daml/pkg/model"
)

type PackageManagement interface {
	ListKnownPackages(ctx context.Context) ([]*model.PackageDetails, error)
	UploadDarFile(ctx context.Context, darFile []byte, submissionID string) error
	ValidateDarFile(ctx context.Context, darFile []byte, submissionID string) error
	UpdateVettedPackages(ctx context.Context, req *model.UpdateVettedPackagesRequest) (*model.UpdateVettedPackagesResponse, error)
}

type packageManagement struct {
	client adminv2.PackageManagementServiceClient
}

func NewPackageManagementClient(conn *grpc.ClientConn) *packageManagement {
	client := adminv2.NewPackageManagementServiceClient(conn)
	return &packageManagement{
		client: client,
	}
}

func (c *packageManagement) ListKnownPackages(ctx context.Context) ([]*model.PackageDetails, error) {
	req := &adminv2.ListKnownPackagesRequest{}

	resp, err := c.client.ListKnownPackages(ctx, req)
	if err != nil {
		return nil, err
	}

	return packageDetailsFromProtos(resp.PackageDetails), nil
}

func (c *packageManagement) UploadDarFile(ctx context.Context, darFile []byte, submissionID string) error {
	req := &adminv2.UploadDarFileRequest{
		DarFile:      darFile,
		SubmissionId: submissionID,
	}

	_, err := c.client.UploadDarFile(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (c *packageManagement) ValidateDarFile(ctx context.Context, darFile []byte, submissionID string) error {
	req := &adminv2.ValidateDarFileRequest{
		DarFile:      darFile,
		SubmissionId: submissionID,
	}

	_, err := c.client.ValidateDarFile(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func packageDetailsFromProto(pb *adminv2.PackageDetails) *model.PackageDetails {
	if pb == nil {
		return nil
	}

	var knownSince *time.Time
	if pb.KnownSince != nil {
		t := pb.KnownSince.AsTime()
		knownSince = &t
	}

	return &model.PackageDetails{
		PackageID:   pb.PackageId,
		PackageSize: pb.PackageSize,
		KnownSince:  knownSince,
		Name:        pb.Name,
		Version:     pb.Version,
	}
}

func packageDetailsFromProtos(pbs []*adminv2.PackageDetails) []*model.PackageDetails {
	result := make([]*model.PackageDetails, len(pbs))
	for i, pb := range pbs {
		result[i] = packageDetailsFromProto(pb)
	}
	return result
}

func (c *packageManagement) UpdateVettedPackages(ctx context.Context, req *model.UpdateVettedPackagesRequest) (*model.UpdateVettedPackagesResponse, error) {
	protoReq := updateVettedPackagesRequestToProto(req)

	resp, err := c.client.UpdateVettedPackages(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	return updateVettedPackagesResponseFromProto(resp), nil
}

func updateVettedPackagesRequestToProto(req *model.UpdateVettedPackagesRequest) *adminv2.UpdateVettedPackagesRequest {
	if req == nil {
		return &adminv2.UpdateVettedPackagesRequest{}
	}

	protoReq := &adminv2.UpdateVettedPackagesRequest{
		DryRun:         req.DryRun,
		SynchronizerId: req.SynchronizerID,
	}

	if req.ExpectedTopologySerial != nil {
		protoReq.ExpectedTopologySerial = &v2.PriorTopologySerial{
			Serial: &v2.PriorTopologySerial_Prior{
				Prior: req.ExpectedTopologySerial.TopologySerial,
			},
		}
	}

	changes := make([]*adminv2.VettedPackagesChange, len(req.Changes))
	for i, change := range req.Changes {
		changes[i] = vettedPackagesChangeToProto(change)
	}
	protoReq.Changes = changes

	forceFlags := make([]adminv2.UpdateVettedPackagesForceFlag, len(req.UpdateVettedPackagesForceFlags))
	for i, flag := range req.UpdateVettedPackagesForceFlags {
		forceFlags[i] = updateVettedPackagesForceFlagToProto(flag)
	}
	protoReq.UpdateVettedPackagesForceFlags = forceFlags

	return protoReq
}

func vettedPackagesChangeToProto(change *model.VettedPackagesChange) *adminv2.VettedPackagesChange {
	if change == nil {
		return nil
	}

	protoChange := &adminv2.VettedPackagesChange{}

	if change.Vet != nil {
		protoVet := &adminv2.VettedPackagesChange_Vet{
			Packages: vettedPackagesRefsToProto(change.Vet.Packages),
		}

		if change.Vet.NewValidFromInclusive != nil {
			protoVet.NewValidFromInclusive = timestamppb.New(*change.Vet.NewValidFromInclusive)
		}

		if change.Vet.NewValidUntilExclusive != nil {
			protoVet.NewValidUntilExclusive = timestamppb.New(*change.Vet.NewValidUntilExclusive)
		}

		protoChange.Operation = &adminv2.VettedPackagesChange_Vet_{
			Vet: protoVet,
		}
	} else if change.Unvet != nil {
		protoUnvet := &adminv2.VettedPackagesChange_Unvet{
			Packages: vettedPackagesRefsToProto(change.Unvet.Packages),
		}

		protoChange.Operation = &adminv2.VettedPackagesChange_Unvet_{
			Unvet: protoUnvet,
		}
	}

	return protoChange
}

func vettedPackagesRefsToProto(refs []*model.VettedPackagesRef) []*adminv2.VettedPackagesRef {
	protoRefs := make([]*adminv2.VettedPackagesRef, len(refs))
	for i, ref := range refs {
		if ref != nil {
			protoRefs[i] = &adminv2.VettedPackagesRef{
				PackageId:      ref.PackageID,
				PackageName:    ref.PackageName,
				PackageVersion: ref.PackageVersion,
			}
		}
	}
	return protoRefs
}

func updateVettedPackagesForceFlagToProto(flag model.UpdateVettedPackagesForceFlag) adminv2.UpdateVettedPackagesForceFlag {
	switch flag {
	case model.UpdateVettedPackagesForceFlagAllowVetIncompatibleUpgrades:
		return adminv2.UpdateVettedPackagesForceFlag_UPDATE_VETTED_PACKAGES_FORCE_FLAG_ALLOW_VET_INCOMPATIBLE_UPGRADES
	case model.UpdateVettedPackagesForceFlagAllowUnvettedDependencies:
		return adminv2.UpdateVettedPackagesForceFlag_UPDATE_VETTED_PACKAGES_FORCE_FLAG_ALLOW_UNVETTED_DEPENDENCIES
	default:
		return adminv2.UpdateVettedPackagesForceFlag_UPDATE_VETTED_PACKAGES_FORCE_FLAG_UNSPECIFIED
	}
}

func updateVettedPackagesResponseFromProto(pb *adminv2.UpdateVettedPackagesResponse) *model.UpdateVettedPackagesResponse {
	if pb == nil {
		return nil
	}

	return &model.UpdateVettedPackagesResponse{
		PastVettedPackages: vettedPackagesFromProtoLedger(pb.PastVettedPackages),
		NewVettedPackages:  vettedPackagesFromProtoLedger(pb.NewVettedPackages),
	}
}

func vettedPackagesFromProtoLedger(pb *v2.VettedPackages) *model.VettedPackages {
	if pb == nil {
		return nil
	}

	packages := make([]*model.VettedPackage, len(pb.Packages))
	for i, pkg := range pb.Packages {
		vp := &model.VettedPackage{
			PackageID:      pkg.PackageId,
			PackageName:    pkg.PackageName,
			PackageVersion: pkg.PackageVersion,
		}

		if pkg.ValidFromInclusive != nil {
			t := pkg.ValidFromInclusive.AsTime()
			vp.ValidFromInclusive = &t
		}

		if pkg.ValidUntilExclusive != nil {
			t := pkg.ValidUntilExclusive.AsTime()
			vp.ValidUntilExclusive = &t
		}

		packages[i] = vp
	}

	return &model.VettedPackages{
		Packages:       packages,
		ParticipantID:  pb.ParticipantId,
		SynchronizerID: pb.SynchronizerId,
		TopologySerial: pb.TopologySerial,
	}
}
