package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/noders-team/go-daml/pkg/client"
	"github.com/noders-team/go-daml/pkg/model"
)

func main() {
	grpcAddress := os.Getenv("GRPC_ADDRESS")
	if grpcAddress == "" {
		grpcAddress = "localhost:8080"
	}

	bearerToken := os.Getenv("BEARER_TOKEN")
	if bearerToken == "" {
		fmt.Println("Warning: BEARER_TOKEN environment variable not set")
	}

	tlsConfig := client.TlsConfig{}

	cl, err := client.NewDamlClient(bearerToken, grpcAddress).
		WithTLSConfig(tlsConfig).
		Build(context.Background())
	if err != nil {
		panic(err)
	}

	pruneUpTo := time.Now().Add(-24 * time.Hour).UnixMicro()

	pruneReq := &model.PruneRequest{
		PruneUpTo:                 pruneUpTo,
		SubmissionID:              fmt.Sprintf("prune-%d", time.Now().Unix()),
		PruneAllDivulgedContracts: false,
	}

	fmt.Printf("Attempting to prune ledger up to: %s (offset: %d)\n",
		time.UnixMicro(pruneUpTo).Format(time.RFC3339), pruneUpTo)

	err = cl.PruningMng.Prune(context.Background(), pruneReq)
	if err != nil {
		fmt.Printf("Prune operation result: %v\n", err)
	} else {
		fmt.Println("Prune operation completed successfully")
	}
}
