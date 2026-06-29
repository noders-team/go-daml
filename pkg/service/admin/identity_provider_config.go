package admin

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	"github.com/noders-team/go-daml/pkg/model"
	adminv2 "github.com/noders-team/go-daml/proto/com/daml/ledger/api/v2/admin"
)

type IdentityProviderConfig interface {
	CreateIdentityProviderConfig(ctx context.Context, config *model.IdentityProviderConfig) (*model.IdentityProviderConfig, error)
	GetIdentityProviderConfig(ctx context.Context, identityProviderID string) (*model.IdentityProviderConfig, error)
	UpdateIdentityProviderConfig(ctx context.Context, config *model.IdentityProviderConfig, updateMask []string) (*model.IdentityProviderConfig, error)
	ListIdentityProviderConfigs(ctx context.Context) ([]*model.IdentityProviderConfig, error)
	DeleteIdentityProviderConfig(ctx context.Context, identityProviderID string) error
}

type identityProviderConfig struct {
	client adminv2.IdentityProviderConfigServiceClient
}

func NewIdentityProviderConfigClient(conn *grpc.ClientConn) *identityProviderConfig {
	client := adminv2.NewIdentityProviderConfigServiceClient(conn)
	return &identityProviderConfig{
		client: client,
	}
}

func (c *identityProviderConfig) CreateIdentityProviderConfig(ctx context.Context, config *model.IdentityProviderConfig) (*model.IdentityProviderConfig, error) {
	req := &adminv2.CreateIdentityProviderConfigRequest{
		IdentityProviderConfig: identityProviderConfigToProto(config),
	}

	resp, err := c.client.CreateIdentityProviderConfig(ctx, req)
	if err != nil {
		return nil, err
	}

	return identityProviderConfigFromProto(resp.IdentityProviderConfig), nil
}

func (c *identityProviderConfig) GetIdentityProviderConfig(ctx context.Context, identityProviderID string) (*model.IdentityProviderConfig, error) {
	req := &adminv2.GetIdentityProviderConfigRequest{
		IdentityProviderId: identityProviderID,
	}

	resp, err := c.client.GetIdentityProviderConfig(ctx, req)
	if err != nil {
		return nil, err
	}

	return identityProviderConfigFromProto(resp.IdentityProviderConfig), nil
}

func (c *identityProviderConfig) UpdateIdentityProviderConfig(ctx context.Context, config *model.IdentityProviderConfig, updateMask []string) (*model.IdentityProviderConfig, error) {
	var fieldMask *fieldmaskpb.FieldMask
	if len(updateMask) > 0 {
		fieldMask = &fieldmaskpb.FieldMask{
			Paths: updateMask,
		}
	}

	req := &adminv2.UpdateIdentityProviderConfigRequest{
		IdentityProviderConfig: identityProviderConfigToProto(config),
		UpdateMask:             fieldMask,
	}

	resp, err := c.client.UpdateIdentityProviderConfig(ctx, req)
	if err != nil {
		return nil, err
	}

	return identityProviderConfigFromProto(resp.IdentityProviderConfig), nil
}

func (c *identityProviderConfig) ListIdentityProviderConfigs(ctx context.Context) ([]*model.IdentityProviderConfig, error) {
	req := &adminv2.ListIdentityProviderConfigsRequest{}

	resp, err := c.client.ListIdentityProviderConfigs(ctx, req)
	if err != nil {
		return nil, err
	}

	return identityProviderConfigsFromProtos(resp.IdentityProviderConfigs), nil
}

func (c *identityProviderConfig) DeleteIdentityProviderConfig(ctx context.Context, identityProviderID string) error {
	req := &adminv2.DeleteIdentityProviderConfigRequest{
		IdentityProviderId: identityProviderID,
	}

	_, err := c.client.DeleteIdentityProviderConfig(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func identityProviderConfigFromProto(pb *adminv2.IdentityProviderConfig) *model.IdentityProviderConfig {
	if pb == nil {
		return nil
	}

	return &model.IdentityProviderConfig{
		IdentityProviderID: pb.IdentityProviderId,
		IsDeactivated:      pb.IsDeactivated,
		Issuer:             pb.Issuer,
		JwksURL:            pb.JwksUrl,
		Audience:           pb.Audience,
	}
}

func identityProviderConfigToProto(config *model.IdentityProviderConfig) *adminv2.IdentityProviderConfig {
	if config == nil {
		return nil
	}

	return &adminv2.IdentityProviderConfig{
		IdentityProviderId: config.IdentityProviderID,
		IsDeactivated:      config.IsDeactivated,
		Issuer:             config.Issuer,
		JwksUrl:            config.JwksURL,
		Audience:           config.Audience,
	}
}

func identityProviderConfigsFromProtos(pbs []*adminv2.IdentityProviderConfig) []*model.IdentityProviderConfig {
	result := make([]*model.IdentityProviderConfig, len(pbs))
	for i, pb := range pbs {
		result[i] = identityProviderConfigFromProto(pb)
	}
	return result
}
