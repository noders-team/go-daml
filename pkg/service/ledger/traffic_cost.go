package ledger

import (
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
)

const (
	transactionPaidTrafficCostField  protowire.Number = 11
	reassignmentPaidTrafficCostField protowire.Number = 9
	completionPaidTrafficCostField   protowire.Number = 12
)

func paidTrafficCostFromUnknown(msg proto.Message, field protowire.Number) int64 {
	if msg == nil {
		return 0
	}

	fields := msg.ProtoReflect().GetUnknown()
	for len(fields) > 0 {
		num, typ, tagLen := protowire.ConsumeTag(fields)
		if tagLen < 0 {
			return 0
		}
		fields = fields[tagLen:]

		if num == field && typ == protowire.VarintType {
			value, valueLen := protowire.ConsumeVarint(fields)
			if valueLen < 0 {
				return 0
			}
			return int64(value)
		}

		valueLen := protowire.ConsumeFieldValue(num, typ, fields)
		if valueLen < 0 {
			return 0
		}
		fields = fields[valueLen:]
	}

	return 0
}
