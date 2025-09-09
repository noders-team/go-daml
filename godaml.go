package godaml

import "github.com/noders-team/go-daml/pkg/client"

func Client(token string, grpcAddress string) *client.DamlClient {
	return client.NewDamlClient(token, grpcAddress)
}
