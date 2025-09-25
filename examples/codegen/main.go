package main

import (
	"context"
	"os"
	"time"

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

	darFilePath := "./test-data/all-kinds-of-1.0.0.dar"
	darContent, err := os.ReadFile(darFilePath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read DAR file")
	}

	submissionID := "validate-" + time.Now().Format("20060102150405")
	log.Info().Str("submissionID", submissionID).Msg("validating DAR file")

	err = cl.PackageMng.ValidateDarFile(context.Background(), darContent, submissionID)
	if err != nil {
		log.Fatal().Err(err).Msgf("DAR validation failed for %s", darFilePath)
	}

	uploadSubmissionID := "upload-" + time.Now().Format("20060102150405")
	log.Info().Str("submissionID", uploadSubmissionID).Msg("uploading DAR file")

	err = cl.PackageMng.UploadDarFile(context.Background(), darContent, uploadSubmissionID)
	if err != nil {
		log.Fatal().Err(err).Msg("DAR upload failed")
	}

	updatedPackages, err := cl.PackageMng.ListKnownPackages(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("failed to list packages after upload")
	} else {
		log.Info().Int("count", len(updatedPackages)).Msg("total packages after upload")
	}
}
