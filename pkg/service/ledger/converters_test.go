package ledger

import (
	"math/big"
	"testing"
	"time"

	v2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/noders-team/go-daml/pkg/model"
	"github.com/noders-team/go-daml/pkg/types"
	"github.com/stretchr/testify/require"
)

// Test type definitions matching the generated code structure
type OneOfEverything struct {
	Operator        types.PARTY     `json:"operator"`
	SomeBoolean     types.BOOL      `json:"someBoolean"`
	SomeInteger     types.INT64     `json:"someInteger"`
	SomeDecimal     types.NUMERIC   `json:"someDecimal"`
	SomeMeasurement types.NUMERIC   `json:"someMeasurement"`
	SomeDate        types.DATE      `json:"someDate"`
	SomeDatetime    types.TIMESTAMP `json:"someDatetime"`
	SomeSimpleList  []types.INT64   `json:"someSimpleList"`
	SomeSimplePair  MyPair          `json:"someSimplePair"`
	SomeNestedPair  MyPair          `json:"someNestedPair"`
	SomeUglyNesting VPair           `json:"someUglyNesting"`
	SomeText        types.TEXT      `json:"someText"`
}

type MyPair struct {
	Left  interface{} `json:"left"`
	Right interface{} `json:"right"`
}

type VPair struct {
	Left  *interface{} `json:"Left,omitempty"`
	Right *interface{} `json:"Right,omitempty"`
	Both  *VPair       `json:"Both,omitempty"`
}

func (t OneOfEverything) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})
	args["operator"] = t.Operator
	args["someBoolean"] = t.SomeBoolean
	args["someInteger"] = t.SomeInteger
	args["someDecimal"] = t.SomeDecimal
	args["someMeasurement"] = t.SomeMeasurement
	args["someDate"] = t.SomeDate
	args["someDatetime"] = t.SomeDatetime
	args["someSimpleList"] = t.SomeSimpleList
	args["someSimplePair"] = t.SomeSimplePair
	args["someNestedPair"] = t.SomeNestedPair
	args["someUglyNesting"] = t.SomeUglyNesting
	args["someText"] = t.SomeText
	return &model.CreateCommand{
		TemplateID: "test:template:OneOfEverything",
		Arguments:  args,
	}
}

func TestConvertToRecord(t *testing.T) {
	t.Run("Decimal", func(t *testing.T) {
		decimalValue := types.NUMERIC(big.NewInt(200))
		data := make(map[string]interface{})
		data["someDecimal"] = decimalValue

		record := convertToRecord(data)
		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "someDecimal", record.Fields[0].Label)

		numericStr := record.Fields[0].Value.GetNumeric()
		require.NotEmpty(t, numericStr)

		// DAML expects the fraction format with "/" replaced by "."
		// big.NewInt(200) with scale 10 becomes 1/50000000 -> "1.50000000"
		require.Equal(t, "1.50000000", numericStr)
	})

	t.Run("Integration test scenario", func(t *testing.T) {
		// Simulate the exact scenario from the integration test
		someListInt := []types.INT64{1, 2, 3}
		left := interface{}("a")
		right := interface{}("b")
		leftInterface := left
		rightInterface := right

		oneOfEverything := OneOfEverything{
			Operator:        types.PARTY("test-party"),
			SomeBoolean:     true,
			SomeInteger:     190,
			SomeDecimal:     types.NUMERIC(big.NewInt(200)),
			SomeMeasurement: types.NUMERIC(big.NewInt(300)),
			SomeDate:        types.DATE(time.Now().UTC()),
			SomeDatetime:    types.TIMESTAMP(time.Now().UTC()),
			SomeSimpleList:  someListInt,
			SomeSimplePair:  MyPair{Left: types.INT64(100), Right: types.INT64(200)},
			SomeNestedPair:  MyPair{Left: MyPair{Left: types.INT64(10), Right: types.INT64(20)}, Right: types.INT64(30)},
			SomeUglyNesting: VPair{Both: &VPair{Left: &leftInterface, Right: &rightInterface}, Left: &leftInterface, Right: &rightInterface},
			SomeText:        "some text",
		}

		// Get the CreateCommand that would be sent to DAML
		createCmd := oneOfEverything.CreateCommand()
		require.NotNil(t, createCmd)
		require.NotNil(t, createCmd.Arguments)

		// Convert the arguments to a record (this is what happens in the ledger)
		record := convertToRecord(createCmd.Arguments)
		require.NotNil(t, record)

		// Check that someDecimal and someMeasurement are properly converted
		fieldMap := make(map[string]*v2.RecordField)
		for _, field := range record.Fields {
			fieldMap[field.Label] = field
		}

		require.NotNil(t, fieldMap["someDecimal"])
		require.NotNil(t, fieldMap["someMeasurement"])

		// Check the types - these should be Numeric, not Record
		require.IsType(t, &v2.Value_Numeric{}, fieldMap["someDecimal"].Value.Sum)
		require.IsType(t, &v2.Value_Numeric{}, fieldMap["someMeasurement"].Value.Sum)

		require.Equal(t, "1.50000000", fieldMap["someDecimal"].Value.GetNumeric())
		require.Equal(t, "3.100000000", fieldMap["someMeasurement"].Value.GetNumeric())
	})

	t.Run("Nested with empty values", func(t *testing.T) {
		type VPair struct {
			Left  *interface{} `json:"Left,omitempty"`
			Right *interface{} `json:"Right,omitempty"`
			Both  *VPair       `json:"Both,omitempty"`
		}

		rightVal := interface{}("b")
		pair := VPair{
			Right: &rightVal,
			Both: &VPair{
				Right: &rightVal,
			},
		}

		data := make(map[string]interface{})
		data["myPair"] = pair

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "myPair", record.Fields[0].Label)

		pairRecord := record.Fields[0].Value.GetRecord()
		require.NotNil(t, pairRecord)
		require.Len(t, pairRecord.Fields, 2) //  Right, Both

		fieldMap := make(map[string]*v2.RecordField)
		for _, field := range pairRecord.Fields {
			fieldMap[field.Label] = field
		}

		require.NotNil(t, fieldMap["Right"])
		require.Equal(t, "b", fieldMap["Right"].Value.GetText())

		require.NotNil(t, fieldMap["Both"])
		bothRecord := fieldMap["Both"].Value.GetRecord()
		require.NotNil(t, bothRecord)

		require.Len(t, bothRecord.Fields, 1)

		nestedFieldMap := make(map[string]*v2.RecordField)
		for _, field := range bothRecord.Fields {
			nestedFieldMap[field.Label] = field
		}

		require.Equal(t, "b", nestedFieldMap["Right"].Value.GetText())
	})

	t.Run("Nested structs", func(t *testing.T) {
		type VPair struct {
			Left  *interface{} `json:"Left,omitempty"`
			Right *interface{} `json:"Right,omitempty"`
			Both  *VPair       `json:"Both,omitempty"`
		}

		leftVal := interface{}("a")
		rightVal := interface{}("b")
		pair := VPair{
			Left:  &leftVal,
			Right: &rightVal,
			Both: &VPair{
				Left:  &leftVal,
				Right: &rightVal,
			},
		}

		data := make(map[string]interface{})
		data["myPair"] = pair

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "myPair", record.Fields[0].Label)

		pairRecord := record.Fields[0].Value.GetRecord()
		require.NotNil(t, pairRecord)
		require.Len(t, pairRecord.Fields, 3) // Left, Right, Both

		fieldMap := make(map[string]*v2.RecordField)
		for _, field := range pairRecord.Fields {
			fieldMap[field.Label] = field
		}

		require.NotNil(t, fieldMap["Left"])
		require.Equal(t, "a", fieldMap["Left"].Value.GetText())

		require.NotNil(t, fieldMap["Right"])
		require.Equal(t, "b", fieldMap["Right"].Value.GetText())

		require.NotNil(t, fieldMap["Both"])
		bothRecord := fieldMap["Both"].Value.GetRecord()
		require.NotNil(t, bothRecord)

		require.Len(t, bothRecord.Fields, 2)

		nestedFieldMap := make(map[string]*v2.RecordField)
		for _, field := range bothRecord.Fields {
			nestedFieldMap[field.Label] = field
		}

		require.Equal(t, "a", nestedFieldMap["Left"].Value.GetText())
		require.Equal(t, "b", nestedFieldMap["Right"].Value.GetText())
	})

	t.Run("Basic struct conversion", func(t *testing.T) {
		type MyPair struct {
			Left  interface{} `json:"left"`
			Right interface{} `json:"right"`
		}

		pair := MyPair{
			Left:  "hello",
			Right: types.INT64(42),
		}

		data := make(map[string]interface{})
		data["myPair"] = pair

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "myPair", record.Fields[0].Label)

		pairRecord := record.Fields[0].Value.GetRecord()
		require.NotNil(t, pairRecord)
		require.Len(t, pairRecord.Fields, 2)

		// Find Left and Right fields
		var leftField, rightField *v2.RecordField
		for _, field := range pairRecord.Fields {
			if field.Label == "left" {
				leftField = field
			} else if field.Label == "right" {
				rightField = field
			}
		}

		require.NotNil(t, leftField)
		require.NotNil(t, rightField)
		require.Equal(t, "hello", leftField.Value.GetText())
		require.Equal(t, int64(42), rightField.Value.GetInt64())
	})

	t.Run("Multiple DAML types", func(t *testing.T) {
		type TestStruct struct {
			TextVal     types.TEXT  `json:"textVal"`
			IntVal      types.INT64 `json:"intVal"`
			BoolVal     types.BOOL  `json:"boolVal"`
			PartyVal    types.PARTY `json:"partyVal"`
			RegularInt  int64       `json:"regularInt"`
			RegularText string      `json:"regularText"`
		}

		testData := TestStruct{
			TextVal:     types.TEXT("test"),
			IntVal:      types.INT64(123),
			BoolVal:     types.BOOL(true),
			PartyVal:    types.PARTY("alice"),
			RegularInt:  456,
			RegularText: "regular",
		}

		data := make(map[string]interface{})
		data["testStruct"] = testData

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)

		structRecord := record.Fields[0].Value.GetRecord()
		require.NotNil(t, structRecord)
		require.Len(t, structRecord.Fields, 6)

		fieldMap := make(map[string]*v2.RecordField)
		for _, field := range structRecord.Fields {
			fieldMap[field.Label] = field
		}

		require.Equal(t, "test", fieldMap["textVal"].Value.GetText())
		require.Equal(t, int64(123), fieldMap["intVal"].Value.GetInt64())
		require.Equal(t, true, fieldMap["boolVal"].Value.GetBool())
		require.Equal(t, "alice", fieldMap["partyVal"].Value.GetParty())
		require.Equal(t, int64(456), fieldMap["regularInt"].Value.GetInt64())
		require.Equal(t, "regular", fieldMap["regularText"].Value.GetText())
	})
}
