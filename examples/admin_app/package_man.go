package main

import (
	"context"
	"os"
	"time"

	"github.com/noders-team/go-daml/pkg/client"
	"github.com/rs/zerolog/log"
)

func RunPackageManagement(cl *client.DamlBindingClient) {
	packages, err := cl.PackageMng.ListKnownPackages(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to list packages")
	}

	log.Info().Msg("known packages:")
	for _, pkg := range packages {
		log.Info().
			Str("packageID", pkg.PackageID).
			Str("name", pkg.Name).
			Str("version", pkg.Version).
			Uint64("size", pkg.PackageSize).
			Time("knownSince", func() time.Time {
				if pkg.KnownSince != nil {
					return *pkg.KnownSince
				}
				return time.Time{}
			}()).
			Msg("package details")
	}

	darFilePath := "./test-data/test.dar"
	log.Info().Str("path", darFilePath).Msg("testing DAR file upload")

	darContent, err := os.ReadFile(darFilePath)
	if err != nil {
		log.Error().Err(err).Msg("failed to read DAR file")
	} else {
		log.Info().Int("size", len(darContent)).Msg("DAR file size")

		submissionID := "validate-" + time.Now().Format("20060102150405")
		log.Info().Str("submissionID", submissionID).Msg("validating DAR file")

		err = cl.PackageMng.ValidateDarFile(context.Background(), darContent, submissionID)
		if err != nil {
			log.Error().Err(err).Msg("DAR validation failed")
		} else {
			log.Info().Msg("DAR validation successful!")

			uploadSubmissionID := "upload-" + time.Now().Format("20060102150405")
			log.Info().Str("submissionID", uploadSubmissionID).Msg("uploading DAR file")

			err = cl.PackageMng.UploadDarFile(context.Background(), darContent, uploadSubmissionID)
			if err != nil {
				log.Error().Err(err).Msg("DAR upload failed")
			} else {
				log.Info().Msg("DAR upload successful!")

				updatedPackages, err := cl.PackageMng.ListKnownPackages(context.Background())
				if err != nil {
					log.Error().Err(err).Msg("failed to list packages after upload")
				} else {
					log.Info().Int("count", len(updatedPackages)).Msg("total packages after upload")
				}
			}
		}
	}
}
