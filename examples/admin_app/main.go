package main

import (
	"context"
	"os"

	"github.com/noders-team/go-daml/pkg/client"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	grpcAddress := os.Getenv("GRPC_ADDRESS")
	if grpcAddress == "" {
		grpcAddress = "localhost:8080"
	}

	bearerToken := os.Getenv("BEARER_TOKEN")
	if bearerToken == "" {
		log.Warn().Msg("BEARER_TOKEN environment variable not set")
	}

	tlsConfig := client.TlsConfig{}

	cl, err := client.NewDamlClient(bearerToken, grpcAddress).
		WithTLSConfig(tlsConfig).
		Build(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to build DAML client")
	}

	log.Info().Msg("=== Starting Users Management ===")
	RunUsersManagement(cl)

	log.Info().Msg("=== Starting Identity Provider Management ===")
	RunIdentityProvider(cl)

	log.Info().Msg("=== Starting Package Management ===")
	RunPackageManagement(cl)

	log.Info().Msg("=== Starting Party Management ===")
	RunPartyManagement(cl)

	log.Info().Msg("=== Starting Pruning ===")
	RunPrunning(cl)
}
