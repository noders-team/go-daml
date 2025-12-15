package admin_test

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/noders-team/go-daml/pkg/client"
	"github.com/noders-team/go-daml/pkg/model"
	"github.com/noders-team/go-daml/pkg/testutil"
	"github.com/stretchr/testify/require"
)

func TestAllocateExternalParty(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	cl := testutil.GetClient()
	require.NotNil(t, cl)
	require.NotNil(t, cl.PartyMng)
	require.NotNil(t, cl.TopologyManagerWrite)

	syncResp, err := cl.StateService.GetConnectedSynchronizers(ctx, &model.GetConnectedSynchronizersRequest{})
	require.NoError(t, err)
	require.NotEmpty(t, syncResp.ConnectedSynchronizers)

	synchronizerID := syncResp.ConnectedSynchronizers[0].SynchronizerID

	participantID, err := cl.PartyMng.GetParticipantID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, participantID)

	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)

	keyFingerprint := computeCantonFingerprint(publicKey, 12)
	namespace := keyFingerprint
	partyID := fmt.Sprintf("external-party-%d", time.Now().Unix())
	fullPartyID := fmt.Sprintf("%s::%s", partyID, namespace)

	onboardingTxs, multiHashSigs, err := createValidOnboardingTransactions(
		ctx,
		cl,
		fullPartyID,
		publicKey,
		privateKey,
		participantID,
	)
	require.NoError(t, err)
	require.NotEmpty(t, onboardingTxs)

	allocatedPartyID, err := cl.PartyMng.AllocateExternalParty(
		ctx,
		synchronizerID,
		onboardingTxs,
		multiHashSigs,
		"",
	)
	require.NoError(t, err)
	require.NotEmpty(t, allocatedPartyID)

	time.Sleep(2 * time.Second)

	parties, err := cl.PartyMng.ListKnownParties(ctx, "", 100, "")
	require.NoError(t, err)
	require.NotEmpty(t, parties.PartyDetails)

	found := false
	for _, party := range parties.PartyDetails {
		if party.Party == allocatedPartyID {
			found = true
			require.True(t, party.IsLocal, "External party should be marked as local after allocation")
			break
		}
	}
	require.True(t, found, "Allocated external party should appear in ListKnownParties")
}

func TestAllocateExternalParty_ErrorHandling(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	cl := testutil.GetClient()
	require.NotNil(t, cl)
	require.NotNil(t, cl.PartyMng)

	syncResp, err := cl.StateService.GetConnectedSynchronizers(ctx, &model.GetConnectedSynchronizersRequest{})
	require.NoError(t, err)
	require.NotEmpty(t, syncResp.ConnectedSynchronizers)

	synchronizerID := syncResp.ConnectedSynchronizers[0].SynchronizerID

	tests := []struct {
		name                   string
		synchronizer           string
		onboardingTransactions []model.SignedTransaction
		multiHashSignatures    []model.Signature
		identityProviderID     string
		expectError            bool
		errorContains          string
	}{
		{
			name:                   "empty synchronizer should fail",
			synchronizer:           "",
			onboardingTransactions: []model.SignedTransaction{createTestSignedTransaction(t)},
			multiHashSignatures:    createTestMultiHashSignatures(t),
			identityProviderID:     "",
			expectError:            true,
			errorContains:          "",
		},
		{
			name:                   "empty onboarding transactions should fail",
			synchronizer:           synchronizerID,
			onboardingTransactions: []model.SignedTransaction{},
			multiHashSignatures:    createTestMultiHashSignatures(t),
			identityProviderID:     "",
			expectError:            true,
			errorContains:          "",
		},
		{
			name:         "invalid transaction format should fail",
			synchronizer: synchronizerID,
			onboardingTransactions: []model.SignedTransaction{
				createTestSignedTransaction(t),
			},
			multiHashSignatures: createTestMultiHashSignatures(t),
			identityProviderID:  "",
			expectError:         true,
			errorContains:       "Invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := cl.PartyMng.AllocateExternalParty(
				ctx,
				tt.synchronizer,
				tt.onboardingTransactions,
				tt.multiHashSignatures,
				tt.identityProviderID,
			)

			if tt.expectError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func createValidOnboardingTransactions(
	ctx context.Context,
	adminCl *client.DamlBindingClient,
	partyID string,
	publicKey ed25519.PublicKey,
	privateKey ed25519.PrivateKey,
	participantID string,
) ([]model.SignedTransaction, []model.Signature, error) {
	pubKey := &model.PublicKey{
		Format:  3,
		Key:     publicKey,
		Scheme:  int32(model.SigningKeySchemeED25519),
		KeySpec: int32(model.SigningKeySpecCurve25519),
		Usage:   []int32{},
	}

	keyFingerprint := computeCantonFingerprint(publicKey, 12)

	storeID := &model.StoreID{Value: "authorized"}

	proposals := []*model.GenerateTransactionProposal{
		{
			Operation: model.OperationAddReplace,
			Serial:    1,
			Mapping: &model.NamespaceDelegationMapping{
				Namespace:        keyFingerprint,
				TargetKey:        *pubKey,
				IsRootDelegation: true,
			},
			Store: storeID,
		},
		{
			Operation: model.OperationAddReplace,
			Serial:    1,
			Mapping: &model.PartyToKeyMapping{
				Party:       partyID,
				Threshold:   1,
				SigningKeys: []model.PublicKey{*pubKey},
			},
			Store: storeID,
		},
		{
			Operation: model.OperationAddReplace,
			Serial:    1,
			Mapping: &model.PartyToParticipantMapping{
				Party:     partyID,
				Threshold: 1,
				Participants: []model.HostingParticipant{
					{
						ParticipantUID: participantID,
						Permission:     model.ParticipantPermissionConfirmation,
					},
				},
			},
			Store: storeID,
		},
	}

	genResp, err := adminCl.TopologyManagerWrite.GenerateTransactions(ctx, &model.GenerateTransactionsRequest{
		Proposals: proposals,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate transactions: %w", err)
	}

	if len(genResp.GeneratedTransactions) != len(proposals) {
		return nil, nil, fmt.Errorf("expected %d transactions, got %d", len(proposals), len(genResp.GeneratedTransactions))
	}

	onboardingTxs := make([]model.SignedTransaction, len(genResp.GeneratedTransactions))
	transactionHashes := make([][]byte, len(genResp.GeneratedTransactions))

	for i, genTx := range genResp.GeneratedTransactions {
		signature := ed25519.Sign(privateKey, genTx.TransactionHash)

		onboardingTxs[i] = model.SignedTransaction{
			Transaction: genTx.SerializedTransaction,
			Signatures: []model.Signature{
				{
					Format:               model.SignatureFormatConcat,
					Signature:            signature,
					SignedBy:             keyFingerprint,
					SigningAlgorithmSpec: model.SigningAlgorithmSpecED25519,
				},
			},
		}

		transactionHashes[i] = genTx.TransactionHash
	}

	multiHash := computeMultiHash(transactionHashes)
	multiHashSignature := ed25519.Sign(privateKey, multiHash)

	multiHashSigs := []model.Signature{
		{
			Format:               model.SignatureFormatConcat,
			Signature:            multiHashSignature,
			SignedBy:             keyFingerprint,
			SigningAlgorithmSpec: model.SigningAlgorithmSpecED25519,
		},
	}

	return onboardingTxs, multiHashSigs, nil
}

func createTestSignedTransaction(t *testing.T) model.SignedTransaction {
	t.Helper()

	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)

	transactionBytes := []byte("test-transaction-data-" + time.Now().Format(time.RFC3339Nano))
	signature := ed25519.Sign(privateKey, transactionBytes)

	return model.SignedTransaction{
		Transaction: transactionBytes,
		Signatures: []model.Signature{
			{
				Format:               model.SignatureFormatConcat,
				Signature:            signature,
				SignedBy:             hex.EncodeToString(publicKey),
				SigningAlgorithmSpec: model.SigningAlgorithmSpecED25519,
			},
		},
	}
}

func createTestMultiHashSignatures(t *testing.T) []model.Signature {
	t.Helper()

	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)

	multiHashData := []byte("test-multi-hash-data-" + time.Now().Format(time.RFC3339Nano))
	signature := ed25519.Sign(privateKey, multiHashData)

	return []model.Signature{
		{
			Format:               model.SignatureFormatConcat,
			Signature:            signature,
			SignedBy:             hex.EncodeToString(publicKey),
			SigningAlgorithmSpec: model.SigningAlgorithmSpecED25519,
		},
	}
}

func computeCantonFingerprint(data []byte, purpose uint32) string {
	buf := make([]byte, 4+len(data))
	buf[0] = byte(purpose >> 24)
	buf[1] = byte(purpose >> 16)
	buf[2] = byte(purpose >> 8)
	buf[3] = byte(purpose)
	copy(buf[4:], data)

	hash := sha256.Sum256(buf)
	multihash := append([]byte{0x12, 0x20}, hash[:]...)
	return hex.EncodeToString(multihash)
}

func computeMultiHash(hashes [][]byte) []byte {
	hasher := sha256.New()
	for _, h := range hashes {
		hasher.Write(h)
	}
	return hasher.Sum(nil)
}
