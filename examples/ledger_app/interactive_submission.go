package main

import (
	"context"
	"time"

	"github.com/noders-team/go-daml/pkg/client"
	"github.com/noders-team/go-daml/pkg/model"
	"github.com/rs/zerolog/log"
)

func RunInteractiveSubmission(cl *client.DamlBindingClient) {
	userID, party := getAvailableUserAndParty(cl)
	packageVersionReq := &model.GetPreferredPackageVersionRequest{
		PackageName: "DamlScript",
		Parties:     []string{party},
	}

	packageVersionResp, err := cl.InteractiveSubmissionService.GetPreferredPackageVersion(context.Background(), packageVersionReq)
	if err != nil {
		log.Warn().Err(err).Msg("failed to get preferred package version")
	} else {
		log.Info().
			Interface("packageReference", packageVersionResp.PackageReference).
			Str("synchronizerId", packageVersionResp.SynchronizerID).
			Msg("got preferred package version")
	}

	prepareReq := &model.PrepareSubmissionRequest{
		UserID:    userID,
		CommandID: "interactive-cmd-" + time.Now().Format("20060102150405"),
		ActAs:     []string{party},
		Commands:  []*model.Command{},
	}

	prepareResp, err := cl.InteractiveSubmissionService.PrepareSubmission(context.Background(), prepareReq)
	if err != nil {
		log.Warn().Err(err).Msg("failed to prepare submission")
	} else {
		log.Info().
			Int("transactionSize", len(prepareResp.PreparedTransaction)).
			Int("hashSize", len(prepareResp.PreparedTransactionHash)).
			Interface("hashingSchemeVersion", prepareResp.HashingSchemeVersion).
			Str("hashingDetails", prepareResp.HashingDetails).
			Msg("prepared interactive submission")

		if len(prepareResp.PreparedTransaction) > 0 {
			executeReq := &model.ExecuteSubmissionRequest{
				PreparedTransaction: prepareResp.PreparedTransaction,
				UserID:              userID,
				SubmissionID:        "exec-" + time.Now().Format("20060102150405"),
				DeduplicationPeriod: model.DeduplicationDuration{
					Duration: 60 * time.Second,
				},
				HashingSchemeVersion: prepareResp.HashingSchemeVersion,
			}

			_, err := cl.InteractiveSubmissionService.ExecuteSubmission(context.Background(), executeReq)
			if err != nil {
				log.Warn().Err(err).Msg("failed to execute submission")
			} else {
				log.Info().Msg("executed interactive submission")
			}
		}
	}
}
