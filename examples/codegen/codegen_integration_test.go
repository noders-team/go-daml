package codegen_test

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	interfaces "github.com/noders-team/go-daml/examples/codegen/interfaces"
	"github.com/noders-team/go-daml/pkg/client"
	"github.com/noders-team/go-daml/pkg/errors"
	"github.com/noders-team/go-daml/pkg/model"
	"github.com/noders-team/go-daml/pkg/service/ledger"
	. "github.com/noders-team/go-daml/pkg/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

const (
	grpcAddress      = "localhost:3901"
	bearerToken      = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJodHRwczovL2NhbnRvbi5uZXR3b3JrLmdsb2JhbCIsInN1YiI6ImxlZGdlci1hcGktdXNlciJ9.A0VZW69lWWNVsjZmDDpVvr1iQ_dJLga3f-K2bicdtsc"
	darFilePath      = "../../test-data/all-kinds-of-1.0.0.dar"
	interfaceDarPath = "../../test-data/amulets-interface-test-1.0.0.dar"
	user             = "app-provider"
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

	if err = cl.ValidateSDKVersion(ctx, SDKVersion); err != nil {
		log.Warn().Err(err).Msg("failed to validate SDK version, ignoring")
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

	marshalled, err := mappyContract.MarshalJSON()
	require.NoError(t, err)
	log.Info().Msgf("marshalled: %s", string(marshalled))

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

	createUpdateID, err := getUpdateIDFromContractCreate(ctx, party, cl, mappyContract)
	require.NoError(t, err)
	require.NotEmpty(t, createUpdateID, "should have a valid create updateID")

	txResp, err := cl.UpdateService.GetTransactionByID(ctx, &model.GetTransactionByIDRequest{
		UpdateID:          createUpdateID,
		RequestingParties: []string{party},
	})
	require.NoError(t, err, "GetTransactionByID should succeed")
	require.NotNil(t, txResp, "response should not be nil")
	require.NotNil(t, txResp.Transaction, "transaction should not be nil")

	var foundTypedContract bool
	for _, event := range txResp.Transaction.Events {
		if event.Created != nil && event.Created.CreateArguments != nil {
			foundTypedContract = true
			var contract MappyContract
			err = ledger.RecordToStruct(event.Created.CreateArguments, &contract)
			require.NoError(t, err, "RecordToStruct should succeed")

			log.Info().
				Str("operator", string(contract.Operator)).
				Interface("value", contract.Value).
				Msg("successfully retrieved typed MappyContract")

			require.Equal(t, PARTY(party), contract.Operator, "operator should match")
			require.NotNil(t, contract.Value, "value should not be nil")
			require.Equal(t, "value1", contract.Value["key1"], "key1 should have correct value")
			require.Equal(t, "value2", contract.Value["key2"], "key2 should have correct value")
		}
	}
	require.True(t, foundTypedContract, "should find at least one typed created event")
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

	// subscribing to updates
	updRes, errRes := cl.UpdateService.GetUpdates(context.Background(), &model.GetUpdatesRequest{Filter: &model.TransactionFilter{
		FiltersByParty: map[string]*model.Filters{
			party: {},
		},
	}})
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case upd, ok := <-updRes:
				if !ok {
					return
				}
				if upd.Update.Transaction != nil {
					log.Info().Str("updateID", upd.Update.Transaction.UpdateID).Str("workflowID", upd.Update.Transaction.WorkflowID).Int("events", len(upd.Update.Transaction.Events)).Msg("received transaction update")
				} else if upd.Update.Reassignment != nil {
					log.Info().Str("updateID", upd.Update.Reassignment.UpdateID).Msg("received reassignment update")
				} else if upd.Update.OffsetCheckpoint != nil {
					log.Info().Int64("offset", upd.Update.OffsetCheckpoint.Offset).Msg("received offset checkpoint")
				}
			case err := <-errRes:
				log.Fatal().Err(err).Msg("failed to get updates")
			}
		}
	}()

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

	createUpdateID, err := getUpdateIDFromContractCreate(ctx, party, cl, mappyContract)
	require.NoError(t, err)
	require.NotEmpty(t, createUpdateID, "should have a valid create updateID")

	txResp, err := cl.UpdateService.GetTransactionByID(ctx, &model.GetTransactionByIDRequest{
		UpdateID:          createUpdateID,
		RequestingParties: []string{party},
	})
	require.NoError(t, err, "GetTransactionByID should succeed")
	require.NotNil(t, txResp, "response should not be nil")
	require.NotNil(t, txResp.Transaction, "transaction should not be nil")

	var foundTypedContract bool
	for _, event := range txResp.Transaction.Events {
		if event.Created != nil && event.Created.CreateArguments != nil {
			foundTypedContract = true
			var contract OneOfEverything
			err = ledger.RecordToStruct(event.Created.CreateArguments, &contract)
			require.NoError(t, err, "RecordToStruct should succeed")

			log.Info().
				Str("operator", string(contract.Operator)).
				Bool("someBoolean", bool(contract.SomeBoolean)).
				Int64("someInteger", int64(contract.SomeInteger)).
				Str("someText", string(contract.SomeText)).
				Msg("successfully retrieved typed OneOfEverything contract")

			require.Equal(t, PARTY(party), contract.Operator, "operator should match")
			require.True(t, bool(contract.SomeBoolean), "someBoolean should be true")
			require.Equal(t, INT64(190), contract.SomeInteger, "someInteger should be 190")
			require.Equal(t, TEXT("some text"), contract.SomeText, "someText should match")
			require.NotNil(t, contract.SomeDecimal, "someDecimal should not be nil")
			require.NotNil(t, contract.SomeMeasurement, "someMeasurement should not be nil")
			require.NotNil(t, contract.SomeMaybe, "someMaybe should not be nil")
			require.Equal(t, INT64(42), *contract.SomeMaybe, "someMaybe value should be 42")
			require.Nil(t, contract.SomeMaybeNot, "someMaybeNot should be nil")
			require.NotNil(t, contract.SomeSimpleList, "someSimpleList should not be nil")
			require.Len(t, contract.SomeSimpleList, 3, "someSimpleList should have 3 elements")
			require.Equal(t, ColorRed, contract.SomeEnum, "someEnum should be ColorRed")
		}
	}
	require.True(t, foundTypedContract, "should find at least one typed created event")
}

func TestAmuletsTransfer(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().Str("interfacePackageID", interfaces.PackageID).Msg("Using interface package ID")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	builder := client.NewDamlClient(bearerToken, grpcAddress)
	if strings.HasSuffix(grpcAddress, ":443") {
		tlsConfig := client.TlsConfig{}
		builder = builder.WithTLSConfig(tlsConfig)
	}

	cl, err := builder.Build(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to build DAML client")
	}

	if err = cl.Ping(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to ping DAML client")
	}

	uploadedPackageName := "amulets-interface-test"
	err = packageUploadWithPath(ctx, uploadedPackageName, interfaceDarPath, interfaces.PackageID, cl)
	if err != nil {
		log.Panic().Msgf("error: %v", err)
	}

	party := ""
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

	assetContract := interfaces.Asset{
		Owner: PARTY(party),
		Name:  TEXT("Test Asset"),
		Value: INT64(100),
	}

	contractIDs, err := createContract(ctx, party, cl, assetContract)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create Asset contract")
	}

	require.Greater(t, len(contractIDs), 0, "Expected at least one contract to be created")
	assetContractID := contractIDs[0]
	log.Info().Str("assetContractID", assetContractID).Msg("Created Asset contract")

	transferArgs := interfaces.Transfer{NewOwner: PARTY(party)}
	transferCmd := assetContract.Transfer(assetContractID, transferArgs)

	transferSubmissionReq := &model.SubmitAndWaitRequest{
		Commands: &model.Commands{
			WorkflowID:   "transfer-workflow-" + time.Now().Format("20060102150405"),
			CommandID:    "transfer-" + time.Now().Format("20060102150405"),
			ActAs:        []string{party},
			SubmissionID: "transfer-sub-" + time.Now().Format("20060102150405"),
			DeduplicationPeriod: model.DeduplicationDuration{
				Duration: 60 * time.Second,
			},
			Commands: []*model.Command{{Command: transferCmd}},
		},
	}

	transferResponse, err := cl.CommandService.SubmitAndWait(ctx, transferSubmissionReq)
	require.NoError(t, err, "Transfer command should succeed")
	log.Info().Str("updateID", transferResponse.UpdateID).Msg("Transfer executed successfully")

	newContractIDs, err := getContractIDsFromUpdate(ctx, party, transferResponse.UpdateID, cl)
	require.NoError(t, err, "Should be able to get contract IDs from transfer transaction")
	require.Greater(t, len(newContractIDs), 0, "Transfer should create at least one contract")

	newAssetContractID := newContractIDs[0]
	log.Info().Str("newAssetContractID", newAssetContractID).Msg("Got new Asset contract from Transfer")

	archiveCmd := assetContract.Archive(newAssetContractID)

	archiveSubmissionReq := &model.SubmitAndWaitRequest{
		Commands: &model.Commands{
			WorkflowID:   "archive-workflow-" + time.Now().Format("20060102150405"),
			CommandID:    "archive-" + time.Now().Format("20060102150405"),
			ActAs:        []string{party},
			SubmissionID: "archive-sub-" + time.Now().Format("20060102150405"),
			DeduplicationPeriod: model.DeduplicationDuration{
				Duration: 60 * time.Second,
			},
			Commands: []*model.Command{{Command: archiveCmd}},
		},
	}

	archiveResponse, err := cl.CommandService.SubmitAndWait(ctx, archiveSubmissionReq)
	require.NoError(t, err, "Archive command should succeed")
	log.Info().Str("updateID", archiveResponse.UpdateID).Msg("Archive executed successfully")

	log.Info().Msg("TestAmuletsTransfer completed successfully")
}

func packageUpload(ctx context.Context, uploadedPackageName string, cl *client.DamlBindingClient) error {
	return packageUploadWithPath(ctx, uploadedPackageName, darFilePath, PackageID, cl)
}

func packageUploadWithPath(ctx context.Context, uploadedPackageName, darPath, packageID string, cl *client.DamlBindingClient) error {
	darContent, err := os.ReadFile(darPath)
	if err != nil {
		return fmt.Errorf("error reading dar file %s: %v", darPath, err)
	}

	if !packageExists(uploadedPackageName, cl) {
		log.Info().Msg("package not found, uploading")

		submissionID := "validate-" + time.Now().Format("20060102150405")
		log.Info().Str("submissionID", submissionID).Msg("validating DAR file")

		err = cl.PackageMng.ValidateDarFile(ctx, darContent, submissionID)
		if err != nil {
			return fmt.Errorf("DAR validation failed for %s %v", darPath, err)
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
		&model.GetPackageStatusRequest{PackageID: packageID})
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

func getUpdateIDFromContractCreate(ctx context.Context, party string, cl *client.DamlBindingClient, template Template) (string, error) {
	createCommands := []*model.Command{
		{
			Command: template.CreateCommand(),
		},
	}

	createSubmissionReq := &model.SubmitAndWaitRequest{
		Commands: &model.Commands{
			WorkflowID:   "create-for-typed-test-" + time.Now().Format("20060102150405"),
			CommandID:    "create-typed-" + time.Now().Format("20060102150405"),
			ActAs:        []string{party},
			SubmissionID: "create-typed-sub-" + time.Now().Format("20060102150405"),
			DeduplicationPeriod: model.DeduplicationDuration{
				Duration: 60 * time.Second,
			},
			Commands: createCommands,
		},
	}

	createResponse, err := cl.CommandService.SubmitAndWait(ctx, createSubmissionReq)
	if err != nil {
		return "", fmt.Errorf("failed to create contract for typed test: %w", err)
	}

	return createResponse.UpdateID, nil
}
