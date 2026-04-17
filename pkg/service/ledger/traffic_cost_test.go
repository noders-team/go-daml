package ledger

import (
	"testing"

	v2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
)

func TestPaidTrafficCostFromUnknown(t *testing.T) {
	tests := []struct {
		name  string
		msg   proto.Message
		field protowire.Number
	}{
		{name: "transaction", msg: &v2.Transaction{}, field: transactionPaidTrafficCostField},
		{name: "reassignment", msg: &v2.Reassignment{}, field: reassignmentPaidTrafficCostField},
		{name: "completion", msg: &v2.Completion{}, field: completionPaidTrafficCostField},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw := protowire.AppendTag(nil, tt.field, protowire.VarintType)
			raw = protowire.AppendVarint(raw, 12345)
			tt.msg.ProtoReflect().SetUnknown(raw)

			require.Equal(t, int64(12345), paidTrafficCostFromUnknown(tt.msg, tt.field))
		})
	}
}
