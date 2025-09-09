package admin

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	adminv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/admin"
)

type UserManagement interface {
	CreateUser(ctx context.Context, userID, primaryParty string, rights []*adminv2.Right) (*adminv2.User, error)
	GetUser(ctx context.Context, userID string) (*adminv2.User, error)
	DeleteUser(ctx context.Context, userID string) error
	GrantUserRights(ctx context.Context, userID string, rights []*adminv2.Right) ([]*adminv2.Right, error)
	RevokeUserRights(ctx context.Context, userID string, rights []*adminv2.Right) ([]*adminv2.Right, error)
	ListUserRights(ctx context.Context, userID string) ([]*adminv2.Right, error)
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

func (c *userManagement) CreateUser(ctx context.Context, userID, primaryParty string, rights []*adminv2.Right) (*adminv2.User, error) {
	request := &adminv2.CreateUserRequest{
		User: &adminv2.User{
			Id:           userID,
			PrimaryParty: primaryParty,
		},
		Rights: rights,
	}

	resp, err := c.client.CreateUser(ctx, request)
	if err != nil {
		return nil, err
	}

	return resp.GetUser(), nil
}

func (c *userManagement) GetUser(ctx context.Context, userID string) (*adminv2.User, error) {
	req := &adminv2.GetUserRequest{
		UserId: userID,
	}

	resp, err := c.client.GetUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return resp.User, nil
}

func (c *userManagement) ListUsers(ctx context.Context) ([]*adminv2.User, error) {
	req := &adminv2.ListUsersRequest{}

	resp, err := c.client.ListUsers(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return resp.Users, nil
}

func (c *userManagement) DeleteUser(ctx context.Context, userID string) error {
	req := &adminv2.DeleteUserRequest{
		UserId: userID,
	}

	_, err := c.client.DeleteUser(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (c *userManagement) GrantUserRights(ctx context.Context, userID string, rights []*adminv2.Right) ([]*adminv2.Right, error) {
	req := &adminv2.GrantUserRightsRequest{
		UserId: userID,
		Rights: rights,
	}

	resp, err := c.client.GrantUserRights(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to grant user rights: %w", err)
	}

	return resp.NewlyGrantedRights, nil
}

func (c *userManagement) RevokeUserRights(ctx context.Context, userID string, rights []*adminv2.Right) ([]*adminv2.Right, error) {
	req := &adminv2.RevokeUserRightsRequest{
		UserId: userID,
		Rights: rights,
	}

	resp, err := c.client.RevokeUserRights(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to revoke user rights: %w", err)
	}

	return resp.NewlyRevokedRights, nil
}

func (c *userManagement) ListUserRights(ctx context.Context, userID string) ([]*adminv2.Right, error) {
	req := &adminv2.ListUserRightsRequest{
		UserId: userID,
	}

	resp, err := c.client.ListUserRights(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list user rights: %w", err)
	}

	return resp.Rights, nil
}
