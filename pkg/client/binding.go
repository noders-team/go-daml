package client

import (
	"github.com/noders-team/go-daml/pkg/service/admin"
	"google.golang.org/grpc"
)

type damlBindingClient struct {
	client  *DamlClient
	grpcCl  *grpc.ClientConn
	UserMng admin.UserManagement
	PartyMng admin.PartyManagement
}

func NewDamlBindingClient(client *DamlClient, grpc *grpc.ClientConn) *damlBindingClient {
	return &damlBindingClient{
		client:   client,
		grpcCl:   grpc,
		UserMng:  admin.NewUserManagementClient(grpc),
		PartyMng: admin.NewPartyManagementClient(grpc),
	}
}

func (c *damlBindingClient) Close() {
	c.grpcCl.Close()
}
