package main

import (
	"context"

	"github.com/noders-team/go-daml/pkg/client"
	"github.com/noders-team/go-daml/pkg/model"
	"github.com/rs/zerolog/log"
)

func RunPackageService(cl *client.DamlBindingClient) {
	listReq := &model.ListPackagesRequest{}
	packages, err := cl.PackageService.ListPackages(context.Background(), listReq)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to list packages")
	}

	log.Info().
		Int("count", len(packages.PackageIDs)).
		Interface("packageIds", packages.PackageIDs).
		Msg("listed packages")

	if len(packages.PackageIDs) > 0 {
		packageID := packages.PackageIDs[0]

		getReq := &model.GetPackageRequest{
			PackageID: packageID,
		}
		packageData, err := cl.PackageService.GetPackage(context.Background(), getReq)
		if err != nil {
			log.Warn().Err(err).Str("packageId", packageID).Msg("failed to get package")
		} else {
			log.Info().
				Str("packageId", packageID).
				Int("archiveSize", len(packageData.ArchivePayload)).
				Str("hash", packageData.Hash).
				Msg("got package details")
		}

		statusReq := &model.GetPackageStatusRequest{
			PackageID: packageID,
		}
		status, err := cl.PackageService.GetPackageStatus(context.Background(), statusReq)
		if err != nil {
			log.Warn().Err(err).Str("packageId", packageID).Msg("failed to get package status")
		} else {
			log.Info().
				Str("packageId", packageID).
				Interface("status", status.PackageStatus).
				Msg("got package status")
		}
	}
}
