package codegen_test

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/noders-team/go-daml/pkg/client"
	"github.com/noders-team/go-daml/pkg/errors"
	"github.com/noders-team/go-daml/pkg/model"
	. "github.com/noders-team/go-daml/pkg/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

const (
	grpcAddress = "localhost:3901"
	bearerToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJodHRwczovL2NhbnRvbi5uZXR3b3JrLmdsb2JhbCIsInN1YiI6ImxlZGdlci1hcGktdXNlciJ9.A0VZW69lWWNVsjZmDDpVvr1iQ_dJLga3f-K2bicdtsc"
	darFilePath = "../../test-data/all-kinds-of-1.0.0.dar"
	user        = "app-provider"
)

func TestCodegenIntegration(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().Str("generatedPackageID", PackageID).Msg("Using package ID from generated code")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	builder := client.NewDamlClient(bearerToken, grpcAddress)
	if strings.HasSuffix(grpcAddress, ":443") {
		tlsConfig := client.TlsConfig{}
		builder = builder.WithTLSConfig(tlsConfig)
	}

	cl, err := builder.
		Build(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to build DAML client")
	}

	if err = cl.Ping(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to ping DAML client")
	}

	uploadedPackageName := "all-kinds-of"
	err = packageUpload(ctx, uploadedPackageName, cl)
	if err != nil {
		log.Panic().Msgf("error: %v", err)
	}

	party := ""

	participantID, err := cl.PartyMng.GetParticipantID(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get participant ID")
	}
	log.Info().Msgf("participantID: %s", participantID)

	users, err := cl.UserMng.ListUsers(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to list users")
	}
	for _, u := range users {
		if u.ID == user {
			party = u.PrimaryParty
			log.Info().Msgf("user %s has primary party %s, using it", u.ID, u.PrimaryParty)
		}
	}

	rights, err := cl.UserMng.ListUserRights(ctx, user)
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
		_, err = cl.UserMng.GrantUserRights(context.Background(), user, "", newRights)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to grant user rights")
		}
	}

	mappyContract := MappyContract{
		Operator: PARTY(party),
		Value: GENMAP{
			"key1": "value1",
			"key2": "value2",
		},
	}

	contractIDs, err := createContract(ctx, party, cl, mappyContract)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to create contract")
	}

	log.Info().Str("templateID", mappyContract.GetTemplateID()).Msg("Using GetTemplateID method")
	if len(contractIDs) == 0 {
		log.Warn().Msg("No contracts were created, cannot demonstrate Archive command")
		return
	}

	var firstContractID string
	if len(contractIDs) > 0 {
		firstContractID = contractIDs[0]
	}
	log.Info().Str("contractID", firstContractID).Msg("Using contract ID from creation for Archive command")
	archiveCmd := mappyContract.Archive(firstContractID)

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

	response, err := cl.CommandService.SubmitAndWait(ctx, submissionReq)
	if err != nil {
		log.Fatal().Err(err).Str("packageId", PackageID).Msg("failed to submit and wait")
	}
	log.Info().Msgf("response.UpdateID: %s", response.UpdateID)

	respUpd, err := cl.UpdateService.GetTransactionByID(ctx, &model.GetTransactionByIDRequest{
		UpdateID:          response.UpdateID,
		RequestingParties: []string{party},
	})
	if err != nil {
		log.Fatal().Err(err).Str("packageId", PackageID).Msg("failed to GetTransactionByID")
	}
	require.NotNil(t, respUpd.Transaction, "expected transaction")
	if respUpd.Transaction != nil {
		for _, event := range respUpd.Transaction.Events {
			if exercisedEvent := event.Exercised; exercisedEvent != nil {
				contractIDs = append(contractIDs, exercisedEvent.ContractID)
				log.Info().
					Str("contractID", exercisedEvent.ContractID).
					Str("templateID", exercisedEvent.TemplateID).
					Msg("found created contract in transaction")
			}
		}
	}
}

func TestCodegenIntegrationAllFieldsContract(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().Str("generatedPackageID", PackageID).Msg("Using package ID from generated code")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	builder := client.NewDamlClient(bearerToken, grpcAddress)
	if strings.HasSuffix(grpcAddress, ":443") {
		tlsConfig := client.TlsConfig{}
		builder = builder.WithTLSConfig(tlsConfig)
	}

	cl, err := builder.
		Build(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to build DAML client")
	}

	if err = cl.Ping(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to ping DAML client")
	}

	uploadedPackageName := "all-kinds-of"
	err = packageUpload(ctx, uploadedPackageName, cl)
	if err != nil {
		log.Panic().Msgf("error: %v", err)
	}

	party := ""

	participantID, err := cl.PartyMng.GetParticipantID(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get participant ID")
	}
	log.Info().Msgf("participantID: %s", participantID)

	users, err := cl.UserMng.ListUsers(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to list users")
	}
	for _, u := range users {
		if u.ID == user {
			party = u.PrimaryParty
			log.Info().Msgf("user %s has primary party %s, using it", u.ID, u.PrimaryParty)
		}
	}

	rights, err := cl.UserMng.ListUserRights(ctx, user)
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
		log.Info().Msg("granting rights")
		newRights := make([]*model.Right, 0)
		newRights = append(newRights, &model.Right{Type: model.CanReadAs{Party: party}})
		_, err = cl.UserMng.GrantUserRights(context.Background(), user, "", newRights)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to grant user rights")
		}
	}

	someListInt := []INT64{1, 2, 3}
	someMaybe := INT64(42)
	mappyContract := OneOfEverything{
		Operator:        PARTY(party),
		SomeBoolean:     true,
		SomeInteger:     190,
		SomeDecimal:     NUMERIC(big.NewInt(200)),
		SomeMeasurement: NUMERIC(big.NewInt(300)),
		SomeMaybe:       &someMaybe,
		SomeMaybeNot:    nil, // Testing optional None case
		SomeDate:        DATE(time.Now().UTC()),
		SomeDatetime:    TIMESTAMP(time.Now().UTC()),
		SomeSimpleList:  someListInt,
		SomeSimplePair:  MyPair{Left: INT64(100), Right: INT64(200)},
		SomeNestedPair: MyPair{
			Left:  MyPair{Left: INT64(10), Right: INT64(20)},
			Right: MyPair{Left: INT64(30), Right: INT64(40)},
		},
		SomeUglyNesting: VPair{
			Both: &VPair{
				Left: func() *interface{} {
					val := interface{}(MyPair{
						Left:  MyPair{Left: INT64(10), Right: INT64(20)},
						Right: MyPair{Left: INT64(30), Right: INT64(40)},
					})
					return &val
				}(),
			},
		},
		SomeText: "some text",
		SomeEnum: ColorRed,
	}

	contractIDs, err := createContract(ctx, party, cl, mappyContract)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to create contract")
		damlErr := errors.AsDamlError(err)
		if strings.EqualFold(damlErr.ErrorCode, "COMMAND_PREPROCESSING_FAILED") {
			log.Fatal().Msgf("failed to create contract: %s", damlErr.Message)
		}
	}

	log.Info().Str("templateID", mappyContract.GetTemplateID()).Msg("Using GetTemplateID method")
	if len(contractIDs) == 0 {
		log.Warn().Msg("No contracts were created, cannot demonstrate Archive command")
		return
	}

	var firstContractID string
	if len(contractIDs) > 0 {
		firstContractID = contractIDs[0]
	}
	log.Info().Str("contractID", firstContractID).Msg("Using contract ID from creation for Archive command")
	archiveCmd := mappyContract.Archive(firstContractID)

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

	response, err := cl.CommandService.SubmitAndWait(ctx, submissionReq)
	if err != nil {
		log.Fatal().Err(err).Str("packageId", PackageID).Msg("failed to submit and wait")
	}
	log.Info().Msgf("response.UpdateID: %s", response.UpdateID)

	respUpd, err := cl.UpdateService.GetTransactionByID(ctx, &model.GetTransactionByIDRequest{
		UpdateID:          response.UpdateID,
		RequestingParties: []string{party},
	})
	if err != nil {
		log.Fatal().Err(err).Str("packageId", PackageID).Msg("failed to GetTransactionByID")
	}
	require.NotNil(t, respUpd.Transaction, "expected transaction")
	if respUpd.Transaction != nil {
		for _, event := range respUpd.Transaction.Events {
			if exercisedEvent := event.Exercised; exercisedEvent != nil {
				contractIDs = append(contractIDs, exercisedEvent.ContractID)
				log.Info().
					Str("contractID", exercisedEvent.ContractID).
					Str("templateID", exercisedEvent.TemplateID).
					Msg("found created contract in transaction")
			}
		}
	}
}

func packageUpload(ctx context.Context, uploadedPackageName string, cl *client.DamlBindingClient) error {
	darContent, err := os.ReadFile(darFilePath)
	if err != nil {
		return fmt.Errorf("error reading darFilePath")
	}

	if !packageExists(uploadedPackageName, cl) {
		log.Info().Msg("package not found, uploading")

		submissionID := "validate-" + time.Now().Format("20060102150405")
		log.Info().Str("submissionID", submissionID).Msg("validating DAR file")

		err = cl.PackageMng.ValidateDarFile(ctx, darContent, submissionID)
		if err != nil {
			return fmt.Errorf("DAR validation failed for %s %v", darFilePath, err)
		}

		uploadSubmissionID := "upload-" + time.Now().Format("20060102150405")
		log.Info().Str("submissionID", uploadSubmissionID).Msg("uploading DAR file")

		err = cl.PackageMng.UploadDarFile(ctx, darContent, uploadSubmissionID)
		if err != nil {
			return fmt.Errorf("DAR upload failed %v", err)
		}

		if !packageExists(uploadedPackageName, cl) {
			return fmt.Errorf("package not found")
		}
	}
	status, err := cl.PackageService.GetPackageStatus(ctx,
		&model.GetPackageStatusRequest{PackageID: PackageID})
	if err != nil {
		return fmt.Errorf("failed to get package status %v", err)
	}
	log.Info().Msgf("package status: %v", status.PackageStatus)

	return nil
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
			pkgInspect, err := cl.PackageService.GetPackage(context.Background(),
				&model.GetPackageRequest{PackageID: pkg.PackageID})
			if err != nil {
				log.Warn().Err(err).Msgf("failed to get package details for %s", pkg.Name)
				return true
			}
			log.Warn().Msgf("package details: Hash: %s HashFunction: %d", pkgInspect.Hash, pkgInspect.HashFunction)
			return true
		}
	}

	return false
}

func createContract(ctx context.Context, party string, cl *client.DamlBindingClient, template Template) ([]string, error) {
	log.Info().Msg("creating sample contracts...")

	createCommands := []*model.Command{
		{
			Command: template.CreateCommand(),
		},
	}

	createSubmissionReq := &model.SubmitAndWaitRequest{
		Commands: &model.Commands{
			WorkflowID:   "create-contracts-" + time.Now().Format("20060102150405"),
			CommandID:    "create-" + time.Now().Format("20060102150405"),
			ActAs:        []string{party},
			SubmissionID: "create-sub-" + time.Now().Format("20060102150405"),
			DeduplicationPeriod: model.DeduplicationDuration{
				Duration: 60 * time.Second,
			},
			Commands: createCommands,
		},
	}

	log.Info().Msg("submitting contract creation commands...")
	createResponse, err := cl.CommandService.SubmitAndWait(context.Background(), createSubmissionReq)
	if err != nil {
		log.Err(err).Msg("failed to create contracts")
		return nil, err
	}
	log.Info().Str("updateID", createResponse.UpdateID).Msg("Successfully created contracts")

	// Use the updateID to get transaction details and extract contract IDs
	contractIDs, err := getContractIDsFromUpdate(ctx, party, createResponse.UpdateID, cl)
	if err != nil {
		log.Err(err).Msg("failed to get contract IDs from update")
		return nil, err
	}

	log.Info().Strs("contractIDs", contractIDs).Msg("extracted contract IDs from transaction")

	return contractIDs, nil
}

func getContractIDsFromUpdate(ctx context.Context, party, updateID string, cl *client.DamlBindingClient) ([]string, error) {
	response, err := cl.UpdateService.GetTransactionByID(ctx, &model.GetTransactionByIDRequest{
		UpdateID:          updateID,
		RequestingParties: []string{party},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction by ID: %w", err)
	}

	var contractIDs []string
	if response.Transaction != nil {
		for _, event := range response.Transaction.Events {
			if createdEvent := event.Created; createdEvent != nil {
				contractIDs = append(contractIDs, createdEvent.ContractID)
				log.Info().
					Str("contractID", createdEvent.ContractID).
					Str("templateID", createdEvent.TemplateID).
					Msg("found created contract in transaction")
			}
		}
	}

	return contractIDs, nil
}
