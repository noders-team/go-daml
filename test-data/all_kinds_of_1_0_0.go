package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/noders-team/go-daml/pkg/model"
)

var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
)

const PackageID = "6d7e83e81a0a7960eec37340f5b11e7a61606bd9161f413684bc345c3f387948"

type (
	PARTY     string
	TEXT      string
	INT64     int64
	BOOL      bool
	DECIMAL   *big.Int
	NUMERIC   *big.Int
	DATE      time.Time
	TIMESTAMP time.Time
	UNIT      struct{}
	LIST      []string
	MAP       map[string]interface{}
	OPTIONAL  *interface{}
	GENMAP    map[string]interface{}
)

// argsToMap converts typed arguments to map for ExerciseCommand
func argsToMap(args interface{}) map[string]interface{} {
	// For now, we'll use a simple approach
	// In practice, you might want to implement proper struct-to-map conversion
	if args == nil {
		return map[string]interface{}{}
	}

	// If args is already a map, return it directly
	if m, ok := args.(map[string]interface{}); ok {
		return m
	}

	// For structs, you would typically use reflection or JSON marshaling
	// For simplicity, we'll return the args in a generic wrapper
	return map[string]interface{}{
		"args": args,
	}
}

// Accept is a Record type
type Accept struct{}

// Color is an enum type
type Color string

const (
	ColorRed   Color = "Red"
	ColorGreen Color = "Green"
	ColorBlue  Color = "Blue"
)

// MappyContract is a Template type
type MappyContract struct {
	Operator PARTY  `json:"operator"`
	Value    GENMAP `json:"value"`
}

// GetTemplateID returns the template ID for this template
func (t MappyContract) GetTemplateID() string {
	return fmt.Sprintf("%s:%s:%s", PackageID, "AllKindsOf", "MappyContract")
}

// Choice methods for MappyContract

// Archive exercises the Archive choice on this MappyContract contract
func (t MappyContract) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("%s:%s:%s", PackageID, "AllKindsOf", "MappyContract"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// MyPair is a Record type
type MyPair struct {
	Left  interface{} `json:"left"`
	Right interface{} `json:"right"`
}

// OneOfEverything is a Template type
type OneOfEverything struct {
	Operator        PARTY     `json:"operator"`
	SomeBoolean     BOOL      `json:"someBoolean"`
	SomeInteger     INT64     `json:"someInteger"`
	SomeDecimal     NUMERIC   `json:"someDecimal"`
	SomeMaybe       OPTIONAL  `json:"someMaybe"`
	SomeMaybeNot    OPTIONAL  `json:"someMaybeNot"`
	SomeText        TEXT      `json:"someText"`
	SomeDate        DATE      `json:"someDate"`
	SomeDatetime    TIMESTAMP `json:"someDatetime"`
	SomeSimpleList  LIST      `json:"someSimpleList"`
	SomeSimplePair  MyPair    `json:"someSimplePair"`
	SomeNestedPair  MyPair    `json:"someNestedPair"`
	SomeUglyNesting VPair     `json:"someUglyNesting"`
	SomeMeasurement NUMERIC   `json:"someMeasurement"`
	SomeEnum        Color     `json:"someEnum"`
	TheUnit         UNIT      `json:"theUnit"`
}

// GetTemplateID returns the template ID for this template
func (t OneOfEverything) GetTemplateID() string {
	return fmt.Sprintf("%s:%s:%s", PackageID, "AllKindsOf", "OneOfEverything")
}

// Choice methods for OneOfEverything

// Archive exercises the Archive choice on this OneOfEverything contract
func (t OneOfEverything) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("%s:%s:%s", PackageID, "AllKindsOf", "OneOfEverything"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// Accept exercises the Accept choice on this OneOfEverything contract
func (t OneOfEverything) Accept(contractID string, args Accept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("%s:%s:%s", PackageID, "AllKindsOf", "OneOfEverything"),
		ContractID: contractID,
		Choice:     "Accept",
		Arguments:  argsToMap(args),
	}
}

// VPair is a variant/union type
type VPair struct {
	Left  *interface{} `json:"Left,omitempty"`
	Right *interface{} `json:"Right,omitempty"`
	Both  *VPair       `json:"Both,omitempty"`
}

// MarshalJSON implements custom JSON marshaling for VPair
func (v VPair) MarshalJSON() ([]byte, error) {
	if v.Left != nil {
		return json.Marshal(map[string]interface{}{
			"tag":   "Left",
			"value": v.Left,
		})
	}

	if v.Right != nil {
		return json.Marshal(map[string]interface{}{
			"tag":   "Right",
			"value": v.Right,
		})
	}

	if v.Both != nil {
		return json.Marshal(map[string]interface{}{
			"tag":   "Both",
			"value": v.Both,
		})
	}

	return json.Marshal(map[string]interface{}{})
}

// UnmarshalJSON implements custom JSON unmarshaling for VPair
func (v *VPair) UnmarshalJSON(data []byte) error {
	var tagged struct {
		Tag   string          `json:"tag"`
		Value json.RawMessage `json:"value"`
	}

	if err := json.Unmarshal(data, &tagged); err != nil {
		return err
	}

	switch tagged.Tag {

	case "Left":
		var value interface{}
		if err := json.Unmarshal(tagged.Value, &value); err != nil {
			return err
		}
		v.Left = &value

	case "Right":
		var value interface{}
		if err := json.Unmarshal(tagged.Value, &value); err != nil {
			return err
		}
		v.Right = &value

	case "Both":
		var value VPair
		if err := json.Unmarshal(tagged.Value, &value); err != nil {
			return err
		}
		v.Both = &value

	default:
		return fmt.Errorf("unknown tag: %s", tagged.Tag)
	}

	return nil
}
