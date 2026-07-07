package gen_clients

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/noders-team/go-daml/pkg/codec"
	"github.com/noders-team/go-daml/pkg/model"
	. "github.com/noders-team/go-daml/pkg/types"
)

var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = fmt.Sprintf
	_ = codec.NewJsonCodec
	_ = model.CreateCommand{}
	_ = NewNumericFromDecimal
)

const packageNameSpliceApiTokenHoldingV1100 = "splice-api-token-holding-v1"

// IHolding is a DAML interface
type IHolding interface {
	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand
}

// HoldingView is a Record type
type HoldingView struct {
	Owner        PARTY        `json:"owner"`
	InstrumentId InstrumentId `json:"instrumentId"`
	Amount       NUMERIC      `json:"amount"`
	Lock         *Lock        `json:"lock"`
	Meta         Metadata     `json:"meta"`
}

// ToMap converts HoldingView to a map for DAML arguments
func (t HoldingView) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["owner"] = t.Owner.ToMap()

	m["instrumentId"] = t.InstrumentId

	m["amount"] = (*big.Int)(t.Amount)

	if t.Lock != nil {
		m["lock"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.Lock,
		}
	} else {
		m["lock"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	m["meta"] = t.Meta

	return m
}

// MarshalJSON implements custom JSON marshaling for HoldingView using JsonCodec
func (t HoldingView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for HoldingView using JsonCodec
func (t *HoldingView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// InstrumentId is a Record type
type InstrumentId struct {
	Admin PARTY `json:"admin"`
	Id    TEXT  `json:"id"`
}

// ToMap converts InstrumentId to a map for DAML arguments
func (t InstrumentId) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["admin"] = t.Admin.ToMap()

	m["id"] = string(t.Id)

	return m
}

// MarshalJSON implements custom JSON marshaling for InstrumentId using JsonCodec
func (t InstrumentId) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for InstrumentId using JsonCodec
func (t *InstrumentId) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Lock is a Record type
type Lock struct {
	Holders      []PARTY    `json:"holders"`
	ExpiresAt    *TIMESTAMP `json:"expiresAt"`
	ExpiresAfter RELTIME    `json:"expiresAfter"`
	Context      *TEXT      `json:"context"`
}

// ToMap converts Lock to a map for DAML arguments
func (t Lock) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["holders"] = t.Holders

	if t.ExpiresAt != nil {
		m["expiresAt"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.ExpiresAt,
		}
	} else {
		m["expiresAt"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	m["expiresAfter"] = t.ExpiresAfter

	if t.Context != nil {
		m["context"] = map[string]interface{}{
			"_type": "optional",
			"value": string(*t.Context),
		}
	} else {
		m["context"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for Lock using JsonCodec
func (t Lock) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for Lock using JsonCodec
func (t *Lock) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// IHoldingInterfaceID returns the interface ID for the IHolding interface
func IHoldingInterfaceID(packageID *string) string {
	pkgName := packageNameSpliceApiTokenHoldingV1100
	if packageID != nil {
		pkgName = *packageID
	}
	return fmt.Sprintf("#%s:%s:%s", pkgName, "Splice.Api.Token.HoldingV1", "Holding")
}
