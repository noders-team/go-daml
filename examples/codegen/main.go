package main

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/noders-team/go-daml/pkg/client"
	"github.com/noders-team/go-daml/pkg/model"
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

	uploadedPackageName := "all-kinds-of"
	if !packageExists(uploadedPackageName, cl) {
		log.Info().Msg("package not found, uploading")

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

		if !packageExists(uploadedPackageName, cl) {
			log.Fatal().Msg("package not found")
		}
	}

	party := "app_provider_localnet-localparty-1::1220716cdae4d7884d468f02b30eb826a7ef54e98f3eb5f875b52a0ef8728ed98c3a"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	contractsCh, errCh := cl.StateService.GetActiveContracts(ctx,
		&model.GetActiveContractsRequest{
			Filter: &model.TransactionFilter{
				FiltersByParty: map[string]*model.Filters{
					party: {
						Inclusive: &model.InclusiveFilters{
							TemplateFilters: []*model.TemplateFilter{
								{
									TemplateID:              "all-kinds-of:ddf0d6396a862eaa7f8d647e39d090a6b04c4a3fd6736aa1730ebc9fca6be664",
									IncludeCreatedEventBlob: true,
								},
							},
						},
					},
				},
			},
			Verbose: false,
		})

	contractCount := 0
	done := false
	for !done {
		select {
		case response, ok := <-contractsCh:
			if !ok {
				log.Info().Int("totalContracts", contractCount).Msg("active contracts stream completed")
				cancel()
				done = true
				break
			}
			if response != nil && len(response.ActiveContracts) > 0 {
				contractCount += len(response.ActiveContracts)
				log.Info().
					Int("activeContracts", len(response.ActiveContracts)).
					Int64("offset", response.Offset).
					Msg("received active contracts batch")
			}
		case err := <-errCh:
			if err != nil {
				log.Error().Err(err).Msg("active contracts stream error")
				cancel()
				done = true
				break
			}
		case <-ctx.Done():
			log.Info().Int("totalContracts", contractCount).Msg("active contracts stream timeout")
			done = true
			break
		}
	}

	log.Info().Msg("finished reading active contracts")
	if contractCount == 0 {
		log.Error().Msgf("no active contracts found for party %s", party)
		return
	}

	packageID := "ddf0d6396a862eaa7f8d647e39d090a6b04c4a3fd6736aa1730ebc9fca6be664"
	status, err := cl.PackageService.GetPackageStatus(context.Background(),
		&model.GetPackageStatusRequest{PackageID: packageID})
	if err != nil {
		log.Fatal().Err(err).Str("packageId", packageID).Msg("failed to get package status")
	}
	log.Info().Msgf("package status: %v", status.PackageStatus)

	/*
		parties, err := cl.PartyMng.ListKnownParties(context.Background(), "", 0, "")
		if err != nil {
			log.Fatal().Err(err).Msg("failed to list parties")
		}

		for _, party := range parties.PartyDetails {
			log.Info().Msgf("party: %+v", party)
		}
	*/
	participantID, err := cl.PartyMng.GetParticipantID(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get participant ID")
	}
	log.Info().Msgf("participantID: %s", participantID)

	users, err := cl.UserMng.ListUsers(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to list users")
	}
	for _, u := range users {
		log.Info().Msgf("user: %+v", u)
	}

	rights, err := cl.UserMng.ListUserRights(context.Background(), "app-provider")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to list user rights")
	}
	rightsGranded := false
	for _, r := range rights {
		canAct, ok := r.Type.(model.RightType).(model.CanActAs)
		if ok && canAct.Party == party {
			rightsGranded = true
		}
	}

	if !rightsGranded {
		log.Info().Msg("grant rights")
		newRights := make([]*model.Right, 0)
		newRights = append(newRights, &model.Right{Type: model.CanReadAs{Party: party}})
		_, err = cl.UserMng.GrantUserRights(context.Background(), "app-provider", "", newRights)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to grant user rights")
		}
	}

	mappyContract := MappyContract{
		Operator: "Alice",
		Value: GENMAP{
			"a": "b",
			"c": "d",
		},
	}

	// Create Archive command
	archiveCmd := mappyContract.Archive(packageID)

	// Submit the Archive command
	commandID := "archive-" + time.Now().Format("20060102150405")
	submissionReq := &model.SubmitAndWaitRequest{
		Commands: &model.Commands{
			WorkflowID:   "archive-workflow-" + time.Now().Format("20060102150405"),
			CommandID:    commandID,
			ActAs:        []string{party},
			SubmissionID: "sub-" + time.Now().Format("20060102150405"),
			DeduplicationPeriod: model.DeduplicationDuration{
				Duration: 60 * time.Second,
			},
			Commands: []*model.Command{{Command: archiveCmd}},
		},
	}

	response, err := cl.CommandService.SubmitAndWait(context.Background(), submissionReq)
	if err != nil {
		log.Fatal().Err(err).Str("packageId", packageID).Msg("failed to get package status")
	}
	log.Info().Msgf("response.UpdateID: %s", response.UpdateID)
}

func packageExists(pkgName string, cl *client.DamlBindingClient) bool {
	updatedPackages, err := cl.PackageMng.ListKnownPackages(context.Background())
	if err != nil {
		log.Warn().Err(err).Msg("failed to list packages after upload")
		return false
	}

	for _, pkg := range updatedPackages {
		if strings.EqualFold(pkg.Name, pkgName) {
			log.Warn().Msgf("package already exists %+v", pkg)
			return true
		}
	}

	return false
}
