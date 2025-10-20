package codegen_test

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
)

const PackageID = "ddf0d6396a862eaa7f8d647e39d090a6b04c4a3fd6736aa1730ebc9fca6be664"

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
}

func argsToMap(args interface{}) map[string]interface{} {
	if args == nil {
		return map[string]interface{}{}
	}

	if m, ok := args.(map[string]interface{}); ok {
		return m
	}

	// Check if the type has a toMap method
	type mapper interface {
		toMap() map[string]interface{}
	}

	if mapper, ok := args.(mapper); ok {
		return mapper.toMap()
	}

	return map[string]interface{}{
		"args": args,
	}
}

// Accept is a Record type
type Accept struct{}

// toMap converts Accept to a map for DAML arguments
func (t Accept) toMap() map[string]interface{} {
	return map[string]interface{}{}
}

// MarshalJSON implements custom JSON marshaling for Accept using JsonCodec
func (a Accept) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(a)
}

// UnmarshalJSON implements custom JSON unmarshaling for Accept using JsonCodec
func (a *Accept) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, a)
}

// Color is an enum type
type Color string

const (
	ColorRed   Color = "Red"
	ColorGreen Color = "Green"
	ColorBlue  Color = "Blue"
)

// GetEnumConstructor implements types.ENUM interface
func (e Color) GetEnumConstructor() string {
	return string(e)
}

// GetEnumTypeID implements types.ENUM interface
func (e Color) GetEnumTypeID() string {
	return fmt.Sprintf("%s:%s:%s", PackageID, "AllKindsOf", "Color")
}

// MarshalJSON implements custom JSON marshaling for Color using JsonCodec
func (c Color) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(c)
}

// UnmarshalJSON implements custom JSON unmarshaling for Color using JsonCodec
func (c *Color) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, c)
}

// Verify interface implementation
var _ ENUM = Color("")

// MappyContract is a Template type
type MappyContract struct {
	Operator PARTY  `json:"operator"`
	Value    GENMAP `json:"value"`
}

// GetTemplateID returns the template ID for this template
func (t MappyContract) GetTemplateID() string {
	return fmt.Sprintf("%s:%s:%s", PackageID, "AllKindsOf", "MappyContract")
}

// CreateCommand returns a CreateCommand for this template
func (t MappyContract) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["operator"] = t.Operator.ToMap()

	if t.Value != nil && len(t.Value) > 0 {
		args["value"] = map[string]interface{}{"_type": "genmap", "value": t.Value}
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
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

// MarshalJSON implements custom JSON marshaling for MappyContract using JsonCodec
func (m MappyContract) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(m)
}

// UnmarshalJSON implements custom JSON unmarshaling for MappyContract using JsonCodec
func (m *MappyContract) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, m)
}

// MyPair is a Record type
type MyPair struct {
	Left  interface{} `json:"left"`
	Right interface{} `json:"right"`
}

// toMap converts MyPair to a map for DAML arguments
func (t MyPair) toMap() map[string]interface{} {
	return map[string]interface{}{
		"left":  t.Left,
		"right": t.Right,
	}
}

// MarshalJSON implements custom JSON marshaling for MyPair using JsonCodec
func (p MyPair) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(p)
}

// UnmarshalJSON implements custom JSON unmarshaling for MyPair using JsonCodec
func (p *MyPair) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, p)
}

// OneOfEverything is a Template type
type OneOfEverything struct {
	Operator        PARTY     `json:"operator"`
	SomeBoolean     BOOL      `json:"someBoolean"`
	SomeInteger     INT64     `json:"someInteger"`
	SomeDecimal     NUMERIC   `json:"someDecimal"`
	SomeMaybe       *INT64    `json:"someMaybe"`
	SomeMaybeNot    *INT64    `json:"someMaybeNot"`
	SomeText        TEXT      `json:"someText"`
	SomeDate        DATE      `json:"someDate"`
	SomeDatetime    TIMESTAMP `json:"someDatetime"`
	SomeSimpleList  []INT64   `json:"someSimpleList"`
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

// CreateCommand returns a CreateCommand for this template
func (t OneOfEverything) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["operator"] = t.Operator.ToMap()

	args["someBoolean"] = bool(t.SomeBoolean)

	args["someInteger"] = int64(t.SomeInteger)

	if t.SomeDecimal != nil {
		args["someDecimal"] = (*big.Int)(t.SomeDecimal)
	}

	if t.SomeMaybe != nil {
		args["someMaybe"] = map[string]interface{}{
			"_type": "optional",
			"value": int64(*t.SomeMaybe),
		}
	} else {
		args["someMaybe"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	if t.SomeMaybeNot != nil {
		args["someMaybeNot"] = map[string]interface{}{
			"_type": "optional",
			"value": int64(*t.SomeMaybeNot),
		}
	} else {
		args["someMaybeNot"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	args["someText"] = string(t.SomeText)

	args["someDate"] = t.SomeDate

	args["someDatetime"] = t.SomeDatetime

	args["someSimpleList"] = t.SomeSimpleList

	args["someSimplePair"] = t.SomeSimplePair

	args["someNestedPair"] = t.SomeNestedPair

	args["someUglyNesting"] = t.SomeUglyNesting

	if t.SomeMeasurement != nil {
		args["someMeasurement"] = (*big.Int)(t.SomeMeasurement)
	}

	if t.SomeEnum != "" {
		args["someEnum"] = t.SomeEnum
	}

	args["theUnit"] = map[string]interface{}{"_type": "unit"}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
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

// MarshalJSON implements custom JSON marshaling for OneOfEverything using JsonCodec
func (o OneOfEverything) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(o)
}

// UnmarshalJSON implements custom JSON unmarshaling for OneOfEverything using JsonCodec
func (o *OneOfEverything) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, o)
}

// VPair is a variant/union type
type VPair struct {
	Left  *interface{} `json:"Left,omitempty"`
	Right *interface{} `json:"Right,omitempty"`
	Both  *VPair       `json:"Both,omitempty"`
}

// MarshalJSON implements custom JSON marshaling for VPair using JsonCodec
func (v VPair) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(v)
}

// UnmarshalJSON implements custom JSON unmarshaling for VPair using JsonCodec
func (v *VPair) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, v)
}

// GetVariantTag implements types.VARIANT interface
func (v VPair) GetVariantTag() string {
	if v.Left != nil {
		return "Left"
	}

	if v.Right != nil {
		return "Right"
	}

	if v.Both != nil {
		return "Both"
	}

	return ""
}

// GetVariantValue implements types.VARIANT interface
func (v VPair) GetVariantValue() interface{} {
	if v.Left != nil {
		return v.Left
	}

	if v.Right != nil {
		return v.Right
	}

	if v.Both != nil {
		return v.Both
	}

	return nil
}

// Verify interface implementation
var _ VARIANT = (*VPair)(nil)
