package main

import (
	"context"
	"fmt"
	"os"

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
		fmt.Println("warning: BEARER_TOKEN environment variable not set")
	}

	tlsConfig := client.TlsConfig{}

	cl, err := client.NewDamlClient(bearerToken, grpcAddress).
		WithTLSConfig(tlsConfig).
		Build(context.Background())
	if err != nil {
		panic(err)
	}

	configs, err := cl.IdentityProviderMng.ListIdentityProviderConfigs(context.Background())
	if err != nil {
		panic(err)
	}
	for _, cfg := range configs {
		println(fmt.Sprintf("identity provider config: %+v", cfg))
	}

	newConfig := &model.IdentityProviderConfig{
		IdentityProviderID: fmt.Sprintf("test-provider-%d", os.Getpid()),
		IsDeactivated:      false,
		Issuer:             "https://example.com",
		JwksURL:            "https://example.com/.well-known/jwks.json",
		Audience:           "https://daml.network",
	}

	createdConfig, err := cl.IdentityProviderMng.CreateIdentityProviderConfig(context.Background(), newConfig)
	if err != nil {
		fmt.Printf("create identity provider error: %v\n", err)
	} else {
		println(fmt.Sprintf("created identity provider: %+v", createdConfig))
	}

	if createdConfig != nil {
		retrievedConfig, err := cl.IdentityProviderMng.GetIdentityProviderConfig(context.Background(), createdConfig.IdentityProviderID)
		if err != nil {
			panic(err)
		}
		println(fmt.Sprintf("retrieved identity provider: %+v", retrievedConfig))

		updatedConfig := &model.IdentityProviderConfig{
			IdentityProviderID: createdConfig.IdentityProviderID,
			IsDeactivated:      false,
			Issuer:             "https://updated.example.com",
			JwksURL:            "https://updated.example.com/.well-known/jwks.json",
			Audience:           "https://daml.network.updated",
		}

		finalConfig, err := cl.IdentityProviderMng.UpdateIdentityProviderConfig(context.Background(), updatedConfig, []string{"issuer", "jwks_url", "audience"})
		if err != nil {
			fmt.Printf("update identity provider error: %v\n", err)
		} else {
			println(fmt.Sprintf("updated identity provider: %+v", finalConfig))
		}

		err = cl.IdentityProviderMng.DeleteIdentityProviderConfig(context.Background(), createdConfig.IdentityProviderID)
		if err != nil {
			fmt.Printf("delete identity provider error: %v\n", err)
		} else {
			fmt.Println("identity provider deleted successfully")
		}
	}

	finalConfigs, err := cl.IdentityProviderMng.ListIdentityProviderConfigs(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println("final identity provider configs:")
	for _, cfg := range finalConfigs {
		println(fmt.Sprintf("  - %+v", cfg))
	}
}
