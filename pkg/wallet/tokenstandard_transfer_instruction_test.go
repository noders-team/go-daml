package wallet

import (
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/noders-team/go-daml/pkg/auth"
	damlclient "github.com/noders-team/go-daml/pkg/client"
	damlModel "github.com/noders-team/go-daml/pkg/model"
	"github.com/noders-team/go-daml/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestCreateTransferInstructionFetchesScanProxyChoiceContext(t *testing.T) {
	blob := []byte("created-event-blob")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v0/scan-proxy/registry/transfer-instruction/v1/transfer-cid/choice-contexts/accept", r.URL.Path)
		require.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.JSONEq(t, `{}`, string(body))

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write([]byte(`{
			"choiceContextData": {
				"values": {
					"external-party-config-state": {"tag": "AV_Text", "value": "state"}
				}
			},
			"disclosedContracts": [{
				"templateId": "#pkg:Module:Template",
				"contractId": "disclosed-cid",
				"createdEventBlob": "` + base64.StdEncoding.EncodeToString(blob) + `",
				"synchronizerId": "sync-id"
			}]
		}`))
		require.NoError(t, err)
	}))
	defer server.Close()

	ctrl := &tokenStandardController{}
	ctrl.SetScanProxyClient(damlclient.NewScanProxyClient(server.URL, auth.NewBearerTokenProvider("test-token")))

	result, err := ctrl.CreateTransferInstruction(context.Background(), "transfer-cid", "accept")
	require.NoError(t, err)
	require.Equal(t, []*damlModel.DisclosedContract{{
		TemplateID:       "#pkg:Module:Template",
		ContractID:       "disclosed-cid",
		CreatedEventBlob: blob,
		SynchronizerID:   "sync-id",
	}}, result.DisclosedContracts)

	exercise, ok := result.Command.Command.(*damlModel.ExerciseCommand)
	require.True(t, ok)
	require.Equal(t, TRANSFER_INSTRUCTION_INTERFACE_ID, exercise.TemplateID)
	require.Equal(t, "transfer-cid", exercise.ContractID)
	require.Equal(t, "TransferInstruction_Accept", exercise.Choice)

	extraArgs, ok := exercise.Arguments["extraArgs"].(map[string]interface{})
	require.True(t, ok)
	contextArg, ok := extraArgs["context"].(map[string]interface{})
	require.True(t, ok)
	valuesArg, ok := contextArg["values"].(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, "textmap", valuesArg["_type"])
	values, ok := valuesArg["value"].(map[string]interface{})
	require.True(t, ok)
	anyValue, ok := values["external-party-config-state"].(types.VARIANT)
	require.True(t, ok)
	require.Equal(t, "AV_Text", anyValue.GetVariantTag())
	require.Equal(t, types.TEXT("state"), anyValue.GetVariantValue())

	metaArg, ok := extraArgs["meta"].(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, map[string]interface{}{"_type": "textmap", "value": map[string]interface{}{}}, metaArg["values"])
}

func TestCreateTransferInstructionRequiresScanProxy(t *testing.T) {
	ctrl := &tokenStandardController{}

	_, err := ctrl.CreateTransferInstruction(context.Background(), "transfer-cid", "accept")
	require.ErrorContains(t, err, "scan-proxy client not configured")
}
