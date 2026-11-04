package admin

import (
	"context"

	"google.golang.org/grpc"

	"github.com/noders-team/go-daml/pkg/model"
	sequencerv30 "github.com/noders-team/go-daml/proto/com/digitalasset/canton/sequencer/admin/v30"
)

type SequencerBFT interface {
	ListConfiguredEndpoints(ctx context.Context, req *model.ListConfiguredEndpointsRequest) (*model.ListConfiguredEndpointsResponse, error)
}

type sequencerBFT struct {
	client sequencerv30.SequencerBftAdministrationServiceClient
}

func NewSequencerBFTClient(conn *grpc.ClientConn) SequencerBFT {
	return &sequencerBFT{
		client: sequencerv30.NewSequencerBftAdministrationServiceClient(conn),
	}
}

func (c *sequencerBFT) ListConfiguredEndpoints(ctx context.Context, req *model.ListConfiguredEndpointsRequest) (*model.ListConfiguredEndpointsResponse, error) {
	resp, err := c.client.ListConfiguredEndpoints(ctx, &sequencerv30.ListConfiguredEndpointsRequest{})
	if err != nil {
		return nil, err
	}

	endpoints := make([]*model.PeerEndpoint, len(resp.Endpoints))
	for i, endpoint := range resp.Endpoints {
		endpoints[i] = peerEndpointFromProto(endpoint)
	}

	return &model.ListConfiguredEndpointsResponse{Endpoints: endpoints}, nil
}

func peerEndpointFromProto(pb *sequencerv30.PeerEndpoint) *model.PeerEndpoint {
	if pb == nil {
		return nil
	}

	endpoint := &model.PeerEndpoint{
		Address:     pb.Address,
		Port:        pb.Port,
		SequencerID: pb.GetSequencerId(),
	}

	switch security := pb.Security.(type) {
	case *sequencerv30.PeerEndpoint_PlainText:
		if security.PlainText != nil {
			endpoint.Security = model.PlainTextPeerEndpoint{}
		}
	case *sequencerv30.PeerEndpoint_Tls:
		if security.Tls != nil {
			tls := &model.TLSPeerEndpoint{
				CustomServerTrustCertificate: security.Tls.CustomServerTrustCertificate,
			}
			if security.Tls.ClientCertificate != nil {
				tls.ClientCertificate = &model.TLSClientCertificate{
					CertificateChain: security.Tls.ClientCertificate.CertificateChain,
					PrivateKeyFile:   security.Tls.ClientCertificate.PrivateKeyFile,
				}
			}
			endpoint.Security = *tls
		}
	}

	return endpoint
}
