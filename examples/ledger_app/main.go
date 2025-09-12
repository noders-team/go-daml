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

	log.Info().Msg("=== Starting Version Service ==.")
	RunVersionService(cl)

	log.Info().Msg("=== Starting Package Service ===")
	RunPackageService(cl)

	log.Info().Msg("=== Starting State Service ===")
	RunStateService(cl)

	log.Info().Msg("=== Starting Update Service ===")
	RunUpdateService(cl)

	log.Info().Msg("=== Starting Command Completion ===")
	RunCommandCompletion(cl)

	log.Info().Msg("=== Starting Command Service ===")
	RunCommandService(cl)

	log.Info().Msg("=== Starting Command Submission ===")
	RunCommandSubmission(cl)

	log.Info().Msg("=== Starting Event Query ===")
	RunEventQuery(cl)

	log.Info().Msg("=== Starting Interactive Submission ===")
	RunInteractiveSubmission(cl)

	log.Info().Msg("=== Starting Time Service ===")
	RunTimeService(cl)
}
