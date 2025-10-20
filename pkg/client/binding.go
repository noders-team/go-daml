package client

import (
	"context"

	"github.com/noders-team/go-daml/pkg/model"
	"github.com/noders-team/go-daml/pkg/service/admin"
	"github.com/noders-team/go-daml/pkg/service/ledger"
	"github.com/noders-team/go-daml/pkg/service/testing"
	"google.golang.org/grpc"
)

type DamlBindingClient struct {
	client                       *DamlClient
	grpcCl                       *grpc.ClientConn
	UserMng                      admin.UserManagement
	PartyMng                     admin.PartyManagement
	PruningMng                   admin.ParticipantPruning
	PackageMng                   admin.PackageManagement
	CommandInspectionMng         admin.CommandInspection
	IdentityProviderMng          admin.IdentityProviderConfig
	CommandCompletion            ledger.CommandCompletion
	CommandService               ledger.CommandService
	CommandSubmission            ledger.CommandSubmission
	EventQuery                   ledger.EventQuery
	PackageService               ledger.PackageService
	StateService                 ledger.StateService
	UpdateService                ledger.UpdateService
	VersionService               ledger.VersionService
	InteractiveSubmissionService ledger.InteractiveSubmissionService
	TimeService                  testing.TimeService
}

func NewDamlBindingClient(client *DamlClient, grpc *grpc.ClientConn) *DamlBindingClient {
	return &DamlBindingClient{
		client:                       client,
		grpcCl:                       grpc,
		UserMng:                      admin.NewUserManagementClient(grpc),
		PartyMng:                     admin.NewPartyManagementClient(grpc),
		PruningMng:                   admin.NewParticipantPruningClient(grpc),
		PackageMng:                   admin.NewPackageManagementClient(grpc),
		CommandInspectionMng:         admin.NewCommandInspectionClient(grpc),
		IdentityProviderMng:          admin.NewIdentityProviderConfigClient(grpc),
		CommandCompletion:            ledger.NewCommandCompletionClient(grpc),
		CommandService:               ledger.NewCommandServiceClient(grpc),
		CommandSubmission:            ledger.NewCommandSubmissionClient(grpc),
		EventQuery:                   ledger.NewEventQueryClient(grpc),
		PackageService:               ledger.NewPackageServiceClient(grpc),
		StateService:                 ledger.NewStateServiceClient(grpc),
		UpdateService:                ledger.NewUpdateServiceClient(grpc),
		VersionService:               ledger.NewVersionServiceClient(grpc),
		InteractiveSubmissionService: ledger.NewInteractiveSubmissionServiceClient(grpc),
		TimeService:                  testing.NewTimeServiceClient(grpc),
	}
}

func (c *DamlBindingClient) Close() {
	c.grpcCl.Close()
}

func (c *DamlBindingClient) Ping(ctx context.Context) error {
	_, err := c.VersionService.GetLedgerAPIVersion(ctx, &model.GetLedgerAPIVersionRequest{})
	return err
}

func (c *DamlBindingClient) ValidateSDKVersion(ctx context.Context, contractSDKVersion string) error {
	if contractSDKVersion == "" {
		return nil
	}

	resp, err := c.VersionService.GetLedgerAPIVersion(ctx, &model.GetLedgerAPIVersionRequest{})
	if err != nil {
		return err
	}

	if resp.Version != contractSDKVersion {
		return &model.SDKVersionMismatchError{
			NodeVersion:     resp.Version,
			ContractVersion: contractSDKVersion,
		}
	}

	return nil
}
