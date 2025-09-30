package admin

import (
	"context"

	"google.golang.org/grpc"

	adminv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/admin"
	"github.com/noders-team/go-daml/pkg/model"
)

type UserManagement interface {
	CreateUser(ctx context.Context, user *model.User, rights []*model.Right) (*model.User, error)
	GetUser(ctx context.Context, userID string) (*model.User, error)
	DeleteUser(ctx context.Context, userID string) error
	GrantUserRights(ctx context.Context, userID, identityProviderID string, rights []*model.Right) ([]*model.Right, error)
	RevokeUserRights(ctx context.Context, userID string, rights []*model.Right) ([]*model.Right, error)
	ListUserRights(ctx context.Context, userID string) ([]*model.Right, error)
	ListUsers(ctx context.Context) ([]*model.User, error)
}

type userManagement struct {
	client adminv2.UserManagementServiceClient
}

func NewUserManagementClient(conn *grpc.ClientConn) *userManagement {
	client := adminv2.NewUserManagementServiceClient(conn)
	return &userManagement{
		client: client,
	}
}

func (c *userManagement) CreateUser(ctx context.Context, user *model.User, rights []*model.Right) (*model.User, error) {
	request := &adminv2.CreateUserRequest{
		User:   userToProto(user),
		Rights: rightsToProto(rights),
	}

	resp, err := c.client.CreateUser(ctx, request)
	if err != nil {
		return nil, err
	}

	return userFromProto(resp.GetUser()), nil
}

func (c *userManagement) GetUser(ctx context.Context, userID string) (*model.User, error) {
	req := &adminv2.GetUserRequest{
		UserId: userID,
	}

	resp, err := c.client.GetUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return userFromProto(resp.User), nil
}

func (c *userManagement) ListUsers(ctx context.Context) ([]*model.User, error) {
	req := &adminv2.ListUsersRequest{}

	resp, err := c.client.ListUsers(ctx, req)
	if err != nil {
		return nil, err
	}

	return usersFromProto(resp.Users), nil
}

func (c *userManagement) DeleteUser(ctx context.Context, userID string) error {
	req := &adminv2.DeleteUserRequest{
		UserId: userID,
	}

	_, err := c.client.DeleteUser(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (c *userManagement) GrantUserRights(ctx context.Context, userID, identityProviderID string, rights []*model.Right) ([]*model.Right, error) {
	req := &adminv2.GrantUserRightsRequest{
		UserId:             userID,
		IdentityProviderId: identityProviderID,
		Rights:             rightsToProto(rights),
	}

	resp, err := c.client.GrantUserRights(ctx, req)
	if err != nil {
		return nil, err
	}

	return rightsFromProto(resp.NewlyGrantedRights), nil
}

func (c *userManagement) RevokeUserRights(ctx context.Context, userID string, rights []*model.Right) ([]*model.Right, error) {
	req := &adminv2.RevokeUserRightsRequest{
		UserId: userID,
		Rights: rightsToProto(rights),
	}

	resp, err := c.client.RevokeUserRights(ctx, req)
	if err != nil {
		return nil, err
	}

	return rightsFromProto(resp.NewlyRevokedRights), nil
}

func (c *userManagement) ListUserRights(ctx context.Context, userID string) ([]*model.Right, error) {
	req := &adminv2.ListUserRightsRequest{
		UserId: userID,
	}

	resp, err := c.client.ListUserRights(ctx, req)
	if err != nil {
		return nil, err
	}

	return rightsFromProto(resp.Rights), nil
}

func userFromProto(pb *adminv2.User) *model.User {
	if pb == nil {
		return nil
	}
	metadata := make(map[string]string)
	if pb.Metadata != nil {
		metadata = pb.Metadata.Annotations
	}
	return &model.User{
		ID:                 pb.Id,
		PrimaryParty:       pb.PrimaryParty,
		IsDeactivated:      pb.IsDeactivated,
		Metadata:           metadata,
		IdentityProviderID: pb.IdentityProviderId,
	}
}

func userToProto(u *model.User) *adminv2.User {
	if u == nil {
		return nil
	}
	var metadata *adminv2.ObjectMeta
	if len(u.Metadata) > 0 {
		metadata = &adminv2.ObjectMeta{
			Annotations: u.Metadata,
		}
	}
	return &adminv2.User{
		Id:                 u.ID,
		PrimaryParty:       u.PrimaryParty,
		IsDeactivated:      u.IsDeactivated,
		Metadata:           metadata,
		IdentityProviderId: u.IdentityProviderID,
	}
}

func rightFromProto(pb *adminv2.Right) *model.Right {
	if pb == nil {
		return nil
	}
	r := &model.Right{}
	switch rt := pb.Kind.(type) {
	case *adminv2.Right_CanActAs_:
		r.Type = model.CanActAs{Party: rt.CanActAs.Party}
	case *adminv2.Right_CanReadAs_:
		r.Type = model.CanReadAs{Party: rt.CanReadAs.Party}
	case *adminv2.Right_ParticipantAdmin_:
		r.Type = model.ParticipantAdmin{}
	case *adminv2.Right_IdentityProviderAdmin_:
		r.Type = model.IdentityProviderAdmin{}
	}
	return r
}

func rightToProto(r *model.Right) *adminv2.Right {
	if r == nil {
		return nil
	}
	pb := &adminv2.Right{}
	switch rt := r.Type.(type) {
	case model.CanActAs:
		pb.Kind = &adminv2.Right_CanActAs_{
			CanActAs: &adminv2.Right_CanActAs{Party: rt.Party},
		}
	case model.CanReadAs:
		pb.Kind = &adminv2.Right_CanReadAs_{
			CanReadAs: &adminv2.Right_CanReadAs{Party: rt.Party},
		}
	case model.ParticipantAdmin:
		pb.Kind = &adminv2.Right_ParticipantAdmin_{
			ParticipantAdmin: &adminv2.Right_ParticipantAdmin{},
		}
	case model.IdentityProviderAdmin:
		pb.Kind = &adminv2.Right_IdentityProviderAdmin_{
			IdentityProviderAdmin: &adminv2.Right_IdentityProviderAdmin{},
		}
	}
	return pb
}

func rightsFromProto(pbs []*adminv2.Right) []*model.Right {
	rights := make([]*model.Right, len(pbs))
	for i, pb := range pbs {
		rights[i] = rightFromProto(pb)
	}
	return rights
}

func rightsToProto(rights []*model.Right) []*adminv2.Right {
	pbs := make([]*adminv2.Right, len(rights))
	for i, r := range rights {
		pbs[i] = rightToProto(r)
	}
	return pbs
}

func usersFromProto(pbs []*adminv2.User) []*model.User {
	users := make([]*model.User, len(pbs))
	for i, pb := range pbs {
		users[i] = userFromProto(pb)
	}
	return users
}
