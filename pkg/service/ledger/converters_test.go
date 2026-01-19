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

// VPair struct for variant tests (needs to be at package level for interface implementation)
type VPairTest struct {
	Left  *interface{} `json:"Left,omitempty"`
	Right *interface{} `json:"Right,omitempty"`
	Both  *VPairTest   `json:"Both,omitempty"`
}

func (v VPairTest) GetVariantTag() string {
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

func (v VPairTest) GetVariantValue() interface{} {
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

// Color enum type for enum tests (needs to be at package level for interface implementation)
type ColorTest string

const (
	ColorTestRed   ColorTest = "Red"
	ColorTestGreen ColorTest = "Green"
	ColorTestBlue  ColorTest = "Blue"
)

func (e ColorTest) GetEnumConstructor() string {
	return string(e)
}

func (e ColorTest) GetEnumTypeID() string {
	return "test-package:TestModule:Color"
}

// StatusTest enum type for enum tests (needs to be at package level for interface implementation)
type StatusTest string

const (
	StatusTestActive   StatusTest = "Active"
	StatusTestInactive StatusTest = "Inactive"
	StatusTestPending  StatusTest = "Pending"
)

func (s StatusTest) GetEnumConstructor() string {
	return string(s)
}

func (s StatusTest) GetEnumTypeID() string {
	return "test-package:TestModule:Status"
}

// VPairIntegration struct for integration tests (needs to be at package level for interface implementation)
type VPairIntegration struct {
	Left  *interface{}      `json:"Left,omitempty"`
	Right *interface{}      `json:"Right,omitempty"`
	Both  *VPairIntegration `json:"Both,omitempty"`
}

func (v VPairIntegration) GetVariantTag() string {
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

func (v VPairIntegration) GetVariantValue() interface{} {
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
var (
	_ types.VARIANT = (*VPairTest)(nil)
	_ types.ENUM    = ColorTest("")
	_ types.ENUM    = StatusTest("")
	_ types.VARIANT = (*VPairIntegration)(nil)
)

func TestConvertToRecordBasic(t *testing.T) {
	t.Run("Numeric", func(t *testing.T) {
		decimalValue := types.NUMERIC(big.NewInt(200))
		data := make(map[string]interface{})
		data["someNumeric"] = decimalValue

		record := convertToRecord(data)
		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "someNumeric", record.Fields[0].Label)

		numericStr := record.Fields[0].Value.GetNumeric()
		require.NotEmpty(t, numericStr)

		require.Equal(t, "0.0000000200", numericStr)
	})

	t.Run("Decimal", func(t *testing.T) {
		decimalValue := types.DECIMAL(big.NewInt(200))
		data := make(map[string]interface{})
		data["someDecimal"] = decimalValue

		record := convertToRecord(data)
		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "someDecimal", record.Fields[0].Label)

		numericStr := record.Fields[0].Value.GetNumeric()
		require.NotEmpty(t, numericStr)

		require.Equal(t, "0.0000000200", numericStr)
	})

	t.Run("*big.Int tests", func(t *testing.T) {
		decimalValue := big.NewInt(200)
		data := make(map[string]interface{})
		data["someBigInt"] = decimalValue

		record := convertToRecord(data)
		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "someBigInt", record.Fields[0].Label)

		numericStr := record.Fields[0].Value.GetNumeric()
		require.NotEmpty(t, numericStr)

		require.Equal(t, "0.0000000200", numericStr)
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

func TestConvertToRecordContractID(t *testing.T) {
	t.Run("CONTRACT_ID type conversion", func(t *testing.T) {
		contractID := types.CONTRACT_ID("00000123456789abcdef")

		data := make(map[string]interface{})
		data["contractId"] = contractID

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "contractId", record.Fields[0].Label)

		// Verify that CONTRACT_ID is converted to proto Value_ContractId
		contractIdValue := record.Fields[0].Value.GetContractId()
		require.Equal(t, "00000123456789abcdef", contractIdValue)
	})

	t.Run("CONTRACT_ID in struct", func(t *testing.T) {
		type TestContractStruct struct {
			Owner      types.PARTY       `json:"owner"`
			ContractID types.CONTRACT_ID `json:"contractId"`
			Name       types.TEXT        `json:"name"`
		}

		testData := TestContractStruct{
			Owner:      types.PARTY("alice"),
			ContractID: types.CONTRACT_ID("00000123456789abcdef"),
			Name:       types.TEXT("test contract"),
		}

		data := make(map[string]interface{})
		data["testStruct"] = testData

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)

		structRecord := record.Fields[0].Value.GetRecord()
		require.NotNil(t, structRecord)
		require.Len(t, structRecord.Fields, 3)

		fieldMap := make(map[string]*v2.RecordField)
		for _, field := range structRecord.Fields {
			fieldMap[field.Label] = field
		}

		require.Equal(t, "alice", fieldMap["owner"].Value.GetParty())
		require.Equal(t, "00000123456789abcdef", fieldMap["contractId"].Value.GetContractId())
		require.Equal(t, "test contract", fieldMap["name"].Value.GetText())
	})
}

func TestConvertToRecordVariant(t *testing.T) {
	t.Run("VARIANT type conversion - Left", func(t *testing.T) {
		leftValue := interface{}("test value")
		variant := VPairTest{
			Left: &leftValue,
		}

		data := make(map[string]interface{})
		data["variant"] = variant

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "variant", record.Fields[0].Label)

		variantValue := record.Fields[0].Value.GetVariant()
		require.NotNil(t, variantValue)
		require.Equal(t, "Left", variantValue.Constructor)
		require.Equal(t, "test value", variantValue.Value.GetText())
	})

	t.Run("VARIANT type conversion - Right", func(t *testing.T) {
		rightValue := interface{}(types.INT64(42))
		variant := VPairTest{
			Right: &rightValue,
		}

		data := make(map[string]interface{})
		data["variant"] = variant

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "variant", record.Fields[0].Label)

		variantValue := record.Fields[0].Value.GetVariant()
		require.NotNil(t, variantValue)
		require.Equal(t, "Right", variantValue.Constructor)
		require.Equal(t, int64(42), variantValue.Value.GetInt64())
	})

	t.Run("VARIANT type conversion - Both (nested)", func(t *testing.T) {
		leftValue := interface{}("nested left")
		rightValue := interface{}("nested right")
		nestedVariant := &VPairTest{
			Left:  &leftValue,
			Right: &rightValue,
		}
		variant := VPairTest{
			Both: nestedVariant,
		}

		data := make(map[string]interface{})
		data["variant"] = variant

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "variant", record.Fields[0].Label)

		variantValue := record.Fields[0].Value.GetVariant()
		require.NotNil(t, variantValue)
		require.Equal(t, "Both", variantValue.Constructor)

		nestedVariantValue := variantValue.Value.GetVariant()
		require.NotNil(t, nestedVariantValue)

		// The nested variant should have both Left and Right, but VARIANT interface
		// returns only the first non-nil value, which should be Left
		require.Equal(t, "Left", nestedVariantValue.Constructor)
		require.Equal(t, "nested left", nestedVariantValue.Value.GetText())
	})

	t.Run("VARIANT in struct", func(t *testing.T) {
		type TestVariantStruct struct {
			Owner   types.PARTY `json:"owner"`
			Variant VPairTest   `json:"variant"`
			Name    types.TEXT  `json:"name"`
		}

		leftValue := interface{}("variant value")
		testData := TestVariantStruct{
			Owner: types.PARTY("alice"),
			Variant: VPairTest{
				Left: &leftValue,
			},
			Name: types.TEXT("test variant"),
		}

		data := make(map[string]interface{})
		data["testStruct"] = testData

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)

		structRecord := record.Fields[0].Value.GetRecord()
		require.NotNil(t, structRecord)
		require.Len(t, structRecord.Fields, 3)

		fieldMap := make(map[string]*v2.RecordField)
		for _, field := range structRecord.Fields {
			fieldMap[field.Label] = field
		}

		require.Equal(t, "alice", fieldMap["owner"].Value.GetParty())
		require.Equal(t, "test variant", fieldMap["name"].Value.GetText())

		variantValue := fieldMap["variant"].Value.GetVariant()
		require.NotNil(t, variantValue)
		require.Equal(t, "Left", variantValue.Constructor)
		require.Equal(t, "variant value", variantValue.Value.GetText())
	})

	t.Run("Nested with empty values", func(t *testing.T) {
		type VPairStruct struct {
			Left  *interface{} `json:"Left,omitempty"`
			Right *interface{} `json:"Right,omitempty"`
			Both  *VPairStruct `json:"Both,omitempty"`
		}

		rightVal := interface{}("b")
		pair := VPairStruct{
			Right: &rightVal,
			Both: &VPairStruct{
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
		type VPairStruct struct {
			Left  *interface{} `json:"Left,omitempty"`
			Right *interface{} `json:"Right,omitempty"`
			Both  *VPairStruct `json:"Both,omitempty"`
		}

		leftVal := interface{}("a")
		rightVal := interface{}("b")
		pair := VPairStruct{
			Left:  &leftVal,
			Right: &rightVal,
			Both: &VPairStruct{
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
}

func TestConvertToRecordEnum(t *testing.T) {
	t.Run("ENUM type conversion - Red", func(t *testing.T) {
		enumValue := ColorTestRed

		data := make(map[string]interface{})
		data["color"] = enumValue

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "color", record.Fields[0].Label)

		// Verify that ENUM is converted to proto Value_Enum
		enumProtoValue := record.Fields[0].Value.GetEnum()
		require.NotNil(t, enumProtoValue)
		require.Equal(t, "Red", enumProtoValue.Constructor)
	})

	t.Run("ENUM type conversion - Green", func(t *testing.T) {
		enumValue := ColorTestGreen

		data := make(map[string]interface{})
		data["color"] = enumValue

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "color", record.Fields[0].Label)

		// Verify that ENUM is converted to proto Value_Enum
		enumProtoValue := record.Fields[0].Value.GetEnum()
		require.NotNil(t, enumProtoValue)
		require.Equal(t, "Green", enumProtoValue.Constructor)
	})

	t.Run("ENUM type conversion - Blue", func(t *testing.T) {
		enumValue := ColorTestBlue

		data := make(map[string]interface{})
		data["color"] = enumValue

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "color", record.Fields[0].Label)

		// Verify that ENUM is converted to proto Value_Enum
		enumProtoValue := record.Fields[0].Value.GetEnum()
		require.NotNil(t, enumProtoValue)
		require.Equal(t, "Blue", enumProtoValue.Constructor)
	})

	t.Run("ENUM in struct", func(t *testing.T) {
		type TestEnumStruct struct {
			Owner     types.PARTY `json:"owner"`
			Color     ColorTest   `json:"color"`
			Name      types.TEXT  `json:"name"`
			IsEnabled types.BOOL  `json:"isEnabled"`
		}

		testData := TestEnumStruct{
			Owner:     types.PARTY("alice"),
			Color:     ColorTestRed,
			Name:      types.TEXT("test enum"),
			IsEnabled: types.BOOL(true),
		}

		data := make(map[string]interface{})
		data["testStruct"] = testData

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)

		structRecord := record.Fields[0].Value.GetRecord()
		require.NotNil(t, structRecord)
		require.Len(t, structRecord.Fields, 4)

		fieldMap := make(map[string]*v2.RecordField)
		for _, field := range structRecord.Fields {
			fieldMap[field.Label] = field
		}

		require.Equal(t, "alice", fieldMap["owner"].Value.GetParty())
		require.Equal(t, "test enum", fieldMap["name"].Value.GetText())
		require.Equal(t, true, fieldMap["isEnabled"].Value.GetBool())

		// Check enum field
		enumProtoValue := fieldMap["color"].Value.GetEnum()
		require.NotNil(t, enumProtoValue)
		require.Equal(t, "Red", enumProtoValue.Constructor)
	})

	t.Run("Multiple ENUMs in struct", func(t *testing.T) {
		type TestMultiEnumStruct struct {
			Color  ColorTest  `json:"color"`
			Status StatusTest `json:"status"`
			Name   types.TEXT `json:"name"`
		}

		testData := TestMultiEnumStruct{
			Color:  ColorTestBlue,
			Status: StatusTestActive,
			Name:   types.TEXT("multi enum test"),
		}

		data := make(map[string]interface{})
		data["testStruct"] = testData

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)

		structRecord := record.Fields[0].Value.GetRecord()
		require.NotNil(t, structRecord)
		require.Len(t, structRecord.Fields, 3)

		fieldMap := make(map[string]*v2.RecordField)
		for _, field := range structRecord.Fields {
			fieldMap[field.Label] = field
		}

		require.Equal(t, "multi enum test", fieldMap["name"].Value.GetText())

		// Check first enum field
		colorEnumValue := fieldMap["color"].Value.GetEnum()
		require.NotNil(t, colorEnumValue)
		require.Equal(t, "Blue", colorEnumValue.Constructor)

		// Check second enum field
		statusEnumValue := fieldMap["status"].Value.GetEnum()
		require.NotNil(t, statusEnumValue)
		require.Equal(t, "Active", statusEnumValue.Constructor)
	})
}

func TestConvertToRecordIntegration(t *testing.T) {
	// Integration test structs
	type MyPairIntegration struct {
		Left  interface{} `json:"left"`
		Right interface{} `json:"right"`
	}

	type OneOfEverythingIntegration struct {
		Operator        types.PARTY       `json:"operator"`
		SomeBoolean     types.BOOL        `json:"someBoolean"`
		SomeInteger     types.INT64       `json:"someInteger"`
		SomeDecimal     types.NUMERIC     `json:"someDecimal"`
		SomeMeasurement types.NUMERIC     `json:"someMeasurement"`
		SomeDate        types.DATE        `json:"someDate"`
		SomeDatetime    types.TIMESTAMP   `json:"someDatetime"`
		SomeSimpleList  []types.INT64     `json:"someSimpleList"`
		SomeSimplePair  MyPairIntegration `json:"someSimplePair"`
		SomeNestedPair  MyPairIntegration `json:"someNestedPair"`
		SomeUglyNesting VPairIntegration  `json:"someUglyNesting"`
		SomeText        types.TEXT        `json:"someText"`
	}

	// CreateCommand method for OneOfEverythingIntegration
	createCommand := func(t OneOfEverythingIntegration) *model.CreateCommand {
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
			TemplateID: "#test:template:OneOfEverything",
			Arguments:  args,
		}
	}

	t.Run("Integration test scenario", func(t *testing.T) {
		someListInt := []types.INT64{1, 2, 3}
		left := interface{}("a")
		right := interface{}("b")
		leftInterface := left
		rightInterface := right

		oneOfEverything := OneOfEverythingIntegration{
			Operator:        types.PARTY("test-party"),
			SomeBoolean:     true,
			SomeInteger:     190,
			SomeDecimal:     types.NUMERIC(big.NewInt(200)),
			SomeMeasurement: types.NUMERIC(big.NewInt(300)),
			SomeDate:        types.DATE(time.Now().UTC()),
			SomeDatetime:    types.TIMESTAMP(time.Now().UTC()),
			SomeSimpleList:  someListInt,
			SomeSimplePair:  MyPairIntegration{Left: types.INT64(100), Right: types.INT64(200)},
			SomeNestedPair:  MyPairIntegration{Left: MyPairIntegration{Left: types.INT64(10), Right: types.INT64(20)}, Right: types.INT64(30)},
			SomeUglyNesting: VPairIntegration{Both: &VPairIntegration{Left: &leftInterface, Right: &rightInterface}, Left: &leftInterface, Right: &rightInterface},
			SomeText:        "some text",
		}

		createCmd := createCommand(oneOfEverything)
		require.NotNil(t, createCmd)
		require.NotNil(t, createCmd.Arguments)

		record := convertToRecord(createCmd.Arguments)
		require.NotNil(t, record)

		fieldMap := make(map[string]*v2.RecordField)
		for _, field := range record.Fields {
			fieldMap[field.Label] = field
		}

		require.NotNil(t, fieldMap["someDecimal"])
		require.NotNil(t, fieldMap["someMeasurement"])

		require.IsType(t, &v2.Value_Numeric{}, fieldMap["someDecimal"].Value.Sum)
		require.IsType(t, &v2.Value_Numeric{}, fieldMap["someMeasurement"].Value.Sum)

		require.Equal(t, "0.0000000200", fieldMap["someDecimal"].Value.GetNumeric())
		require.Equal(t, "0.0000000300", fieldMap["someMeasurement"].Value.GetNumeric())
	})
}

func TestConvertToRecordOptional(t *testing.T) {
	t.Run("Optional with value", func(t *testing.T) {
		optionalData := map[string]interface{}{
			"_type": "optional",
			"value": int64(42),
		}

		data := make(map[string]interface{})
		data["optionalField"] = optionalData

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "optionalField", record.Fields[0].Label)

		optionalValue := record.Fields[0].Value.GetOptional()
		require.NotNil(t, optionalValue)
		require.NotNil(t, optionalValue.Value)
		require.Equal(t, int64(42), optionalValue.Value.GetInt64())
	})

	t.Run("Optional without value (None)", func(t *testing.T) {
		optionalData := map[string]interface{}{
			"_type": "optional",
		}

		data := make(map[string]interface{})
		data["optionalField"] = optionalData

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "optionalField", record.Fields[0].Label)

		optionalValue := record.Fields[0].Value.GetOptional()
		require.NotNil(t, optionalValue)
		require.Nil(t, optionalValue.Value)
	})

	t.Run("Multiple optional fields in struct", func(t *testing.T) {
		type TestOptionalStruct struct {
			Name      types.TEXT   `json:"name"`
			MaybeInt  *types.INT64 `json:"maybeInt"`
			MaybeText *types.TEXT  `json:"maybeText"`
			MaybeBool *types.BOOL  `json:"maybeBool"`
		}

		data := make(map[string]interface{})
		data["name"] = types.TEXT("test")
		data["maybeInt"] = map[string]interface{}{
			"_type": "optional",
			"value": int64(100),
		}
		data["maybeText"] = map[string]interface{}{
			"_type": "optional",
			"value": "optional text",
		}
		data["maybeBool"] = map[string]interface{}{
			"_type": "optional",
		}

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 4)

		fieldMap := make(map[string]*v2.RecordField)
		for _, field := range record.Fields {
			fieldMap[field.Label] = field
		}

		require.Equal(t, "test", fieldMap["name"].Value.GetText())

		optInt := fieldMap["maybeInt"].Value.GetOptional()
		require.NotNil(t, optInt)
		require.NotNil(t, optInt.Value)
		require.Equal(t, int64(100), optInt.Value.GetInt64())

		optText := fieldMap["maybeText"].Value.GetOptional()
		require.NotNil(t, optText)
		require.NotNil(t, optText.Value)
		require.Equal(t, "optional text", optText.Value.GetText())

		optBool := fieldMap["maybeBool"].Value.GetOptional()
		require.NotNil(t, optBool)
		require.Nil(t, optBool.Value)
	})
}

func TestConvertToRecordSlices(t *testing.T) {
	t.Run("[]types.INT64 conversion", func(t *testing.T) {
		sliceData := []types.INT64{1, 2, 3, 42, 100}

		data := make(map[string]interface{})
		data["intList"] = sliceData

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "intList", record.Fields[0].Label)

		listValue := record.Fields[0].Value.GetList()
		require.NotNil(t, listValue)
		require.Len(t, listValue.Elements, 5)

		require.Equal(t, int64(1), listValue.Elements[0].GetInt64())
		require.Equal(t, int64(2), listValue.Elements[1].GetInt64())
		require.Equal(t, int64(3), listValue.Elements[2].GetInt64())
		require.Equal(t, int64(42), listValue.Elements[3].GetInt64())
		require.Equal(t, int64(100), listValue.Elements[4].GetInt64())
	})

	t.Run("[]types.TEXT conversion", func(t *testing.T) {
		sliceData := []types.TEXT{"hello", "world", "test", "slice"}

		data := make(map[string]interface{})
		data["textList"] = sliceData

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "textList", record.Fields[0].Label)

		listValue := record.Fields[0].Value.GetList()
		require.NotNil(t, listValue)
		require.Len(t, listValue.Elements, 4)

		require.Equal(t, "hello", listValue.Elements[0].GetText())
		require.Equal(t, "world", listValue.Elements[1].GetText())
		require.Equal(t, "test", listValue.Elements[2].GetText())
		require.Equal(t, "slice", listValue.Elements[3].GetText())
	})

	t.Run("[]types.BOOL conversion", func(t *testing.T) {
		sliceData := []types.BOOL{true, false, true, false, true}

		data := make(map[string]interface{})
		data["boolList"] = sliceData

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "boolList", record.Fields[0].Label)

		listValue := record.Fields[0].Value.GetList()
		require.NotNil(t, listValue)
		require.Len(t, listValue.Elements, 5)

		require.Equal(t, true, listValue.Elements[0].GetBool())
		require.Equal(t, false, listValue.Elements[1].GetBool())
		require.Equal(t, true, listValue.Elements[2].GetBool())
		require.Equal(t, false, listValue.Elements[3].GetBool())
		require.Equal(t, true, listValue.Elements[4].GetBool())
	})

	t.Run("[]int64 conversion", func(t *testing.T) {
		sliceData := []int64{10, 20, 30, 40, 50}

		data := make(map[string]interface{})
		data["regularIntList"] = sliceData

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "regularIntList", record.Fields[0].Label)

		listValue := record.Fields[0].Value.GetList()
		require.NotNil(t, listValue)
		require.Len(t, listValue.Elements, 5)

		require.Equal(t, int64(10), listValue.Elements[0].GetInt64())
		require.Equal(t, int64(20), listValue.Elements[1].GetInt64())
		require.Equal(t, int64(30), listValue.Elements[2].GetInt64())
		require.Equal(t, int64(40), listValue.Elements[3].GetInt64())
		require.Equal(t, int64(50), listValue.Elements[4].GetInt64())
	})

	t.Run("[]string conversion", func(t *testing.T) {
		sliceData := []string{"apple", "banana", "cherry", "date"}

		data := make(map[string]interface{})
		data["stringList"] = sliceData

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "stringList", record.Fields[0].Label)

		listValue := record.Fields[0].Value.GetList()
		require.NotNil(t, listValue)
		require.Len(t, listValue.Elements, 4)

		require.Equal(t, "apple", listValue.Elements[0].GetText())
		require.Equal(t, "banana", listValue.Elements[1].GetText())
		require.Equal(t, "cherry", listValue.Elements[2].GetText())
		require.Equal(t, "date", listValue.Elements[3].GetText())
	})

	t.Run("[]interface{} conversion", func(t *testing.T) {
		sliceData := []interface{}{"text", int64(42), true, types.PARTY("alice"), types.TEXT("daml-text")}

		data := make(map[string]interface{})
		data["interfaceList"] = sliceData

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "interfaceList", record.Fields[0].Label)

		listValue := record.Fields[0].Value.GetList()
		require.NotNil(t, listValue)
		require.Len(t, listValue.Elements, 5)

		require.Equal(t, "text", listValue.Elements[0].GetText())
		require.Equal(t, int64(42), listValue.Elements[1].GetInt64())
		require.Equal(t, true, listValue.Elements[2].GetBool())
		require.Equal(t, "alice", listValue.Elements[3].GetParty())
		require.Equal(t, "daml-text", listValue.Elements[4].GetText())
	})

	t.Run("types.LIST conversion", func(t *testing.T) {
		listData := types.LIST{"first", "second", "third", "fourth"}

		data := make(map[string]interface{})
		data["damlList"] = listData

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "damlList", record.Fields[0].Label)

		listValue := record.Fields[0].Value.GetList()
		require.NotNil(t, listValue)
		require.Len(t, listValue.Elements, 4)

		require.Equal(t, "first", listValue.Elements[0].GetText())
		require.Equal(t, "second", listValue.Elements[1].GetText())
		require.Equal(t, "third", listValue.Elements[2].GetText())
		require.Equal(t, "fourth", listValue.Elements[3].GetText())
	})

	t.Run("Mixed slices in struct", func(t *testing.T) {
		type TestSliceStruct struct {
			IntList       []types.INT64 `json:"intList"`
			TextList      []types.TEXT  `json:"textList"`
			BoolList      []types.BOOL  `json:"boolList"`
			RegIntList    []int64       `json:"regIntList"`
			StrList       []string      `json:"strList"`
			InterfaceList []interface{} `json:"interfaceList"`
			DamlList      types.LIST    `json:"damlList"`
		}

		testData := TestSliceStruct{
			IntList:       []types.INT64{1, 2, 3},
			TextList:      []types.TEXT{"a", "b", "c"},
			BoolList:      []types.BOOL{true, false, true},
			RegIntList:    []int64{10, 20, 30},
			StrList:       []string{"x", "y", "z"},
			InterfaceList: []interface{}{"mixed", int64(99), false},
			DamlList:      types.LIST{"daml1", "daml2", "daml3"},
		}

		data := make(map[string]interface{})
		data["testStruct"] = testData

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)

		structRecord := record.Fields[0].Value.GetRecord()
		require.NotNil(t, structRecord)
		require.Len(t, structRecord.Fields, 7)

		fieldMap := make(map[string]*v2.RecordField)
		for _, field := range structRecord.Fields {
			fieldMap[field.Label] = field
		}

		// Verify intList
		intList := fieldMap["intList"].Value.GetList()
		require.NotNil(t, intList)
		require.Len(t, intList.Elements, 3)
		require.Equal(t, int64(1), intList.Elements[0].GetInt64())
		require.Equal(t, int64(2), intList.Elements[1].GetInt64())
		require.Equal(t, int64(3), intList.Elements[2].GetInt64())

		// Verify textList
		textList := fieldMap["textList"].Value.GetList()
		require.NotNil(t, textList)
		require.Len(t, textList.Elements, 3)
		require.Equal(t, "a", textList.Elements[0].GetText())
		require.Equal(t, "b", textList.Elements[1].GetText())
		require.Equal(t, "c", textList.Elements[2].GetText())

		// Verify boolList
		boolList := fieldMap["boolList"].Value.GetList()
		require.NotNil(t, boolList)
		require.Len(t, boolList.Elements, 3)
		require.Equal(t, true, boolList.Elements[0].GetBool())
		require.Equal(t, false, boolList.Elements[1].GetBool())
		require.Equal(t, true, boolList.Elements[2].GetBool())

		// Verify regIntList
		regIntList := fieldMap["regIntList"].Value.GetList()
		require.NotNil(t, regIntList)
		require.Len(t, regIntList.Elements, 3)
		require.Equal(t, int64(10), regIntList.Elements[0].GetInt64())
		require.Equal(t, int64(20), regIntList.Elements[1].GetInt64())
		require.Equal(t, int64(30), regIntList.Elements[2].GetInt64())

		// Verify strList
		strList := fieldMap["strList"].Value.GetList()
		require.NotNil(t, strList)
		require.Len(t, strList.Elements, 3)
		require.Equal(t, "x", strList.Elements[0].GetText())
		require.Equal(t, "y", strList.Elements[1].GetText())
		require.Equal(t, "z", strList.Elements[2].GetText())

		// Verify interfaceList
		interfaceList := fieldMap["interfaceList"].Value.GetList()
		require.NotNil(t, interfaceList)
		require.Len(t, interfaceList.Elements, 3)
		require.Equal(t, "mixed", interfaceList.Elements[0].GetText())
		require.Equal(t, int64(99), interfaceList.Elements[1].GetInt64())
		require.Equal(t, false, interfaceList.Elements[2].GetBool())

		// Verify damlList
		damlList := fieldMap["damlList"].Value.GetList()
		require.NotNil(t, damlList)
		require.Len(t, damlList.Elements, 3)
		require.Equal(t, "daml1", damlList.Elements[0].GetText())
		require.Equal(t, "daml2", damlList.Elements[1].GetText())
		require.Equal(t, "daml3", damlList.Elements[2].GetText())
	})
}

func TestConvertToRecordRELTIME(t *testing.T) {
	t.Run("RELTIME type conversion - 1 second", func(t *testing.T) {
		reltimeValue := types.RELTIME(1 * time.Second)

		data := make(map[string]interface{})
		data["duration"] = reltimeValue

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "duration", record.Fields[0].Label)

		int64Value := record.Fields[0].Value.GetInt64()
		require.Equal(t, int64(1000000), int64Value)
	})

	t.Run("RELTIME type conversion - 5 minutes", func(t *testing.T) {
		reltimeValue := types.RELTIME(5 * time.Minute)

		data := make(map[string]interface{})
		data["duration"] = reltimeValue

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "duration", record.Fields[0].Label)

		int64Value := record.Fields[0].Value.GetInt64()
		require.Equal(t, int64(300000000), int64Value)
	})

	t.Run("RELTIME type conversion - 100 microseconds", func(t *testing.T) {
		reltimeValue := types.RELTIME(100 * time.Microsecond)

		data := make(map[string]interface{})
		data["duration"] = reltimeValue

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "duration", record.Fields[0].Label)

		int64Value := record.Fields[0].Value.GetInt64()
		require.Equal(t, int64(100), int64Value)
	})

	t.Run("RELTIME in struct", func(t *testing.T) {
		type TestReltimeStruct struct {
			Owner       types.PARTY   `json:"owner"`
			Duration    types.RELTIME `json:"duration"`
			Name        types.TEXT    `json:"name"`
			MaxDuration types.RELTIME `json:"maxDuration"`
		}

		testData := TestReltimeStruct{
			Owner:       types.PARTY("alice"),
			Duration:    types.RELTIME(30 * time.Second),
			Name:        types.TEXT("test reltime"),
			MaxDuration: types.RELTIME(1 * time.Hour),
		}

		data := make(map[string]interface{})
		data["testStruct"] = testData

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)

		structRecord := record.Fields[0].Value.GetRecord()
		require.NotNil(t, structRecord)
		require.Len(t, structRecord.Fields, 4)

		fieldMap := make(map[string]*v2.RecordField)
		for _, field := range structRecord.Fields {
			fieldMap[field.Label] = field
		}

		require.Equal(t, "alice", fieldMap["owner"].Value.GetParty())
		require.Equal(t, "test reltime", fieldMap["name"].Value.GetText())
		require.Equal(t, int64(30000000), fieldMap["duration"].Value.GetInt64())
		require.Equal(t, int64(3600000000), fieldMap["maxDuration"].Value.GetInt64())
	})
}

func TestConvertToRecordSET(t *testing.T) {
	t.Run("SET type conversion - strings", func(t *testing.T) {
		setValue := types.SET{"item1", "item2", "item3"}

		data := make(map[string]interface{})
		data["set"] = setValue

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "set", record.Fields[0].Label)

		listValue := record.Fields[0].Value.GetList()
		require.NotNil(t, listValue)
		require.Len(t, listValue.Elements, 3)
		require.Equal(t, "item1", listValue.Elements[0].GetText())
		require.Equal(t, "item2", listValue.Elements[1].GetText())
		require.Equal(t, "item3", listValue.Elements[2].GetText())
	})

	t.Run("SET type conversion - integers", func(t *testing.T) {
		setValue := types.SET{int64(1), int64(2), int64(3)}

		data := make(map[string]interface{})
		data["set"] = setValue

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "set", record.Fields[0].Label)

		listValue := record.Fields[0].Value.GetList()
		require.NotNil(t, listValue)
		require.Len(t, listValue.Elements, 3)
		require.Equal(t, int64(1), listValue.Elements[0].GetInt64())
		require.Equal(t, int64(2), listValue.Elements[1].GetInt64())
		require.Equal(t, int64(3), listValue.Elements[2].GetInt64())
	})

	t.Run("SET type conversion - mixed types", func(t *testing.T) {
		setValue := types.SET{"text", int64(42), true}

		data := make(map[string]interface{})
		data["set"] = setValue

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "set", record.Fields[0].Label)

		listValue := record.Fields[0].Value.GetList()
		require.NotNil(t, listValue)
		require.Len(t, listValue.Elements, 3)
		require.Equal(t, "text", listValue.Elements[0].GetText())
		require.Equal(t, int64(42), listValue.Elements[1].GetInt64())
		require.Equal(t, true, listValue.Elements[2].GetBool())
	})

	t.Run("SET type conversion - empty set", func(t *testing.T) {
		setValue := types.SET{}

		data := make(map[string]interface{})
		data["set"] = setValue

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "set", record.Fields[0].Label)

		listValue := record.Fields[0].Value.GetList()
		require.NotNil(t, listValue)
		require.Len(t, listValue.Elements, 0)
	})

	t.Run("SET in struct", func(t *testing.T) {
		type TestSetStruct struct {
			Owner           types.PARTY `json:"owner"`
			RequiredParties types.SET   `json:"requiredParties"`
			Name            types.TEXT  `json:"name"`
			AllowedValues   types.SET   `json:"allowedValues"`
		}

		testData := TestSetStruct{
			Owner:           types.PARTY("alice"),
			RequiredParties: types.SET{"alice", "bob", "charlie"},
			Name:            types.TEXT("test set"),
			AllowedValues:   types.SET{int64(1), int64(2), int64(3)},
		}

		data := make(map[string]interface{})
		data["testStruct"] = testData

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)

		structRecord := record.Fields[0].Value.GetRecord()
		require.NotNil(t, structRecord)
		require.Len(t, structRecord.Fields, 4)

		fieldMap := make(map[string]*v2.RecordField)
		for _, field := range structRecord.Fields {
			fieldMap[field.Label] = field
		}

		require.Equal(t, "alice", fieldMap["owner"].Value.GetParty())
		require.Equal(t, "test set", fieldMap["name"].Value.GetText())

		partiesSet := fieldMap["requiredParties"].Value.GetList()
		require.NotNil(t, partiesSet)
		require.Len(t, partiesSet.Elements, 3)
		require.Equal(t, "alice", partiesSet.Elements[0].GetText())
		require.Equal(t, "bob", partiesSet.Elements[1].GetText())
		require.Equal(t, "charlie", partiesSet.Elements[2].GetText())

		valuesSet := fieldMap["allowedValues"].Value.GetList()
		require.NotNil(t, valuesSet)
		require.Len(t, valuesSet.Elements, 3)
		require.Equal(t, int64(1), valuesSet.Elements[0].GetInt64())
		require.Equal(t, int64(2), valuesSet.Elements[1].GetInt64())
		require.Equal(t, int64(3), valuesSet.Elements[2].GetInt64())
	})

	t.Run("SET with DAML types", func(t *testing.T) {
		setValue := types.SET{
			types.PARTY("alice"),
			types.PARTY("bob"),
			types.PARTY("charlie"),
		}

		data := make(map[string]interface{})
		data["parties"] = setValue

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "parties", record.Fields[0].Label)

		listValue := record.Fields[0].Value.GetList()
		require.NotNil(t, listValue)
		require.Len(t, listValue.Elements, 3)
		require.Equal(t, "alice", listValue.Elements[0].GetParty())
		require.Equal(t, "bob", listValue.Elements[1].GetParty())
		require.Equal(t, "charlie", listValue.Elements[2].GetParty())
	})
}

func TestConvertToRecordRELTIMEAndSETIntegration(t *testing.T) {
	t.Run("RELTIME and SET in same struct", func(t *testing.T) {
		type TestIntegrationStruct struct {
			Owner                types.PARTY   `json:"owner"`
			TickDuration         types.RELTIME `json:"tickDuration"`
			RequiredParties      types.SET     `json:"requiredParties"`
			Name                 types.TEXT    `json:"name"`
			MaxProcessingTime    types.RELTIME `json:"maxProcessingTime"`
			AllowedSynchronizers types.SET     `json:"allowedSynchronizers"`
		}

		testData := TestIntegrationStruct{
			Owner:                types.PARTY("alice"),
			TickDuration:         types.RELTIME(10 * time.Second),
			RequiredParties:      types.SET{"alice", "bob", "charlie"},
			Name:                 types.TEXT("integration test"),
			MaxProcessingTime:    types.RELTIME(5 * time.Minute),
			AllowedSynchronizers: types.SET{"sync1", "sync2"},
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

		require.Equal(t, "alice", fieldMap["owner"].Value.GetParty())
		require.Equal(t, "integration test", fieldMap["name"].Value.GetText())

		require.Equal(t, int64(10000000), fieldMap["tickDuration"].Value.GetInt64())
		require.Equal(t, int64(300000000), fieldMap["maxProcessingTime"].Value.GetInt64())

		partiesSet := fieldMap["requiredParties"].Value.GetList()
		require.NotNil(t, partiesSet)
		require.Len(t, partiesSet.Elements, 3)
		require.Equal(t, "alice", partiesSet.Elements[0].GetText())
		require.Equal(t, "bob", partiesSet.Elements[1].GetText())
		require.Equal(t, "charlie", partiesSet.Elements[2].GetText())

		syncSet := fieldMap["allowedSynchronizers"].Value.GetList()
		require.NotNil(t, syncSet)
		require.Len(t, syncSet.Elements, 2)
		require.Equal(t, "sync1", syncSet.Elements[0].GetText())
		require.Equal(t, "sync2", syncSet.Elements[1].GetText())
	})
}

func TestConvertToRecordTUPLE2(t *testing.T) {
	t.Run("TUPLE2 type conversion - strings", func(t *testing.T) {
		tuple2Value := types.TUPLE2[string, string]{
			First:  "hello",
			Second: "world",
		}

		data := make(map[string]interface{})
		data["tuple"] = tuple2Value

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "tuple", record.Fields[0].Label)

		tupleRecord := record.Fields[0].Value.GetRecord()
		require.NotNil(t, tupleRecord)
		require.Len(t, tupleRecord.Fields, 2)

		fieldMap := make(map[string]*v2.RecordField)
		for _, field := range tupleRecord.Fields {
			fieldMap[field.Label] = field
		}

		require.Equal(t, "hello", fieldMap["_1"].Value.GetText())
		require.Equal(t, "world", fieldMap["_2"].Value.GetText())
	})

	t.Run("TUPLE2 type conversion - integers", func(t *testing.T) {
		tuple2Value := types.TUPLE2[int64, int64]{
			First:  int64(42),
			Second: int64(100),
		}

		data := make(map[string]interface{})
		data["tuple"] = tuple2Value

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "tuple", record.Fields[0].Label)

		tupleRecord := record.Fields[0].Value.GetRecord()
		require.NotNil(t, tupleRecord)
		require.Len(t, tupleRecord.Fields, 2)

		fieldMap := make(map[string]*v2.RecordField)
		for _, field := range tupleRecord.Fields {
			fieldMap[field.Label] = field
		}

		require.Equal(t, int64(42), fieldMap["_1"].Value.GetInt64())
		require.Equal(t, int64(100), fieldMap["_2"].Value.GetInt64())
	})

	t.Run("TUPLE2 type conversion - mixed types", func(t *testing.T) {
		tuple2Value := types.TUPLE2[string, int64]{
			First:  "text",
			Second: int64(42),
		}

		data := make(map[string]interface{})
		data["tuple"] = tuple2Value

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "tuple", record.Fields[0].Label)

		tupleRecord := record.Fields[0].Value.GetRecord()
		require.NotNil(t, tupleRecord)
		require.Len(t, tupleRecord.Fields, 2)

		fieldMap := make(map[string]*v2.RecordField)
		for _, field := range tupleRecord.Fields {
			fieldMap[field.Label] = field
		}

		require.Equal(t, "text", fieldMap["_1"].Value.GetText())
		require.Equal(t, int64(42), fieldMap["_2"].Value.GetInt64())
	})

	t.Run("TUPLE2 type conversion - DAML types", func(t *testing.T) {
		tuple2Value := types.TUPLE2[types.PARTY, types.TEXT]{
			First:  types.PARTY("alice"),
			Second: types.TEXT("test"),
		}

		data := make(map[string]interface{})
		data["tuple"] = tuple2Value

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "tuple", record.Fields[0].Label)

		tupleRecord := record.Fields[0].Value.GetRecord()
		require.NotNil(t, tupleRecord)
		require.Len(t, tupleRecord.Fields, 2)

		fieldMap := make(map[string]*v2.RecordField)
		for _, field := range tupleRecord.Fields {
			fieldMap[field.Label] = field
		}

		require.Equal(t, "alice", fieldMap["_1"].Value.GetParty())
		require.Equal(t, "test", fieldMap["_2"].Value.GetText())
	})

	t.Run("TUPLE2 in struct", func(t *testing.T) {
		type TestTuple2Struct struct {
			Owner      types.PARTY                  `json:"owner"`
			Coordinate types.TUPLE2[int64, int64]   `json:"coordinate"`
			Name       types.TEXT                   `json:"name"`
			Pair       types.TUPLE2[string, string] `json:"pair"`
		}

		testData := TestTuple2Struct{
			Owner: types.PARTY("alice"),
			Coordinate: types.TUPLE2[int64, int64]{
				First:  int64(100),
				Second: int64(200),
			},
			Name: types.TEXT("test tuple2"),
			Pair: types.TUPLE2[string, string]{
				First:  "key",
				Second: "value",
			},
		}

		data := make(map[string]interface{})
		data["testStruct"] = testData

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)

		structRecord := record.Fields[0].Value.GetRecord()
		require.NotNil(t, structRecord)
		require.Len(t, structRecord.Fields, 4)

		fieldMap := make(map[string]*v2.RecordField)
		for _, field := range structRecord.Fields {
			fieldMap[field.Label] = field
		}

		require.Equal(t, "alice", fieldMap["owner"].Value.GetParty())
		require.Equal(t, "test tuple2", fieldMap["name"].Value.GetText())

		coordinateRecord := fieldMap["coordinate"].Value.GetRecord()
		require.NotNil(t, coordinateRecord)
		require.Len(t, coordinateRecord.Fields, 2)
		coordFieldMap := make(map[string]*v2.RecordField)
		for _, field := range coordinateRecord.Fields {
			coordFieldMap[field.Label] = field
		}
		require.Equal(t, int64(100), coordFieldMap["_1"].Value.GetInt64())
		require.Equal(t, int64(200), coordFieldMap["_2"].Value.GetInt64())

		pairRecord := fieldMap["pair"].Value.GetRecord()
		require.NotNil(t, pairRecord)
		require.Len(t, pairRecord.Fields, 2)
		pairFieldMap := make(map[string]*v2.RecordField)
		for _, field := range pairRecord.Fields {
			pairFieldMap[field.Label] = field
		}
		require.Equal(t, "key", pairFieldMap["_1"].Value.GetText())
		require.Equal(t, "value", pairFieldMap["_2"].Value.GetText())
	})

	t.Run("Nested TUPLE2", func(t *testing.T) {
		innerTuple := types.TUPLE2[string, string]{
			First:  "inner1",
			Second: "inner2",
		}

		outerTuple := types.TUPLE2[types.TUPLE2[string, string], int64]{
			First:  innerTuple,
			Second: int64(999),
		}

		data := make(map[string]interface{})
		data["nestedTuple"] = outerTuple

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "nestedTuple", record.Fields[0].Label)

		outerRecord := record.Fields[0].Value.GetRecord()
		require.NotNil(t, outerRecord)
		require.Len(t, outerRecord.Fields, 2)

		outerFieldMap := make(map[string]*v2.RecordField)
		for _, field := range outerRecord.Fields {
			outerFieldMap[field.Label] = field
		}

		innerRecord := outerFieldMap["_1"].Value.GetRecord()
		require.NotNil(t, innerRecord)
		require.Len(t, innerRecord.Fields, 2)

		innerFieldMap := make(map[string]*v2.RecordField)
		for _, field := range innerRecord.Fields {
			innerFieldMap[field.Label] = field
		}

		require.Equal(t, "inner1", innerFieldMap["_1"].Value.GetText())
		require.Equal(t, "inner2", innerFieldMap["_2"].Value.GetText())
		require.Equal(t, int64(999), outerFieldMap["_2"].Value.GetInt64())
	})

	t.Run("TUPLE2 with complex DAML types", func(t *testing.T) {
		tuple2Value := types.TUPLE2[types.RELTIME, types.SET]{
			First:  types.RELTIME(30 * time.Second),
			Second: types.SET{"alice", "bob", "charlie"},
		}

		data := make(map[string]interface{})
		data["complexTuple"] = tuple2Value

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "complexTuple", record.Fields[0].Label)

		tupleRecord := record.Fields[0].Value.GetRecord()
		require.NotNil(t, tupleRecord)
		require.Len(t, tupleRecord.Fields, 2)

		fieldMap := make(map[string]*v2.RecordField)
		for _, field := range tupleRecord.Fields {
			fieldMap[field.Label] = field
		}

		require.Equal(t, int64(30000000), fieldMap["_1"].Value.GetInt64())

		setList := fieldMap["_2"].Value.GetList()
		require.NotNil(t, setList)
		require.Len(t, setList.Elements, 3)
		require.Equal(t, "alice", setList.Elements[0].GetText())
		require.Equal(t, "bob", setList.Elements[1].GetText())
		require.Equal(t, "charlie", setList.Elements[2].GetText())
	})
}

func TestConvertToRecordRELTIMEAndSETAndTUPLE2Integration(t *testing.T) {
	t.Run("RELTIME, SET, and TUPLE2 in same struct", func(t *testing.T) {
		type TestFullIntegrationStruct struct {
			Owner             types.PARTY                  `json:"owner"`
			TickDuration      types.RELTIME                `json:"tickDuration"`
			RequiredParties   types.SET                    `json:"requiredParties"`
			Name              types.TEXT                   `json:"name"`
			Coordinate        types.TUPLE2[int64, int64]   `json:"coordinate"`
			MaxProcessingTime types.RELTIME                `json:"maxProcessingTime"`
			Metadata          types.TUPLE2[string, string] `json:"metadata"`
		}

		testData := TestFullIntegrationStruct{
			Owner:           types.PARTY("alice"),
			TickDuration:    types.RELTIME(10 * time.Second),
			RequiredParties: types.SET{"alice", "bob", "charlie"},
			Name:            types.TEXT("full integration test"),
			Coordinate: types.TUPLE2[int64, int64]{
				First:  int64(100),
				Second: int64(200),
			},
			MaxProcessingTime: types.RELTIME(5 * time.Minute),
			Metadata: types.TUPLE2[string, string]{
				First:  "version",
				Second: "1.0.0",
			},
		}

		data := make(map[string]interface{})
		data["testStruct"] = testData

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)

		structRecord := record.Fields[0].Value.GetRecord()
		require.NotNil(t, structRecord)
		require.Len(t, structRecord.Fields, 7)

		fieldMap := make(map[string]*v2.RecordField)
		for _, field := range structRecord.Fields {
			fieldMap[field.Label] = field
		}

		require.Equal(t, "alice", fieldMap["owner"].Value.GetParty())
		require.Equal(t, "full integration test", fieldMap["name"].Value.GetText())

		require.Equal(t, int64(10000000), fieldMap["tickDuration"].Value.GetInt64())
		require.Equal(t, int64(300000000), fieldMap["maxProcessingTime"].Value.GetInt64())

		partiesSet := fieldMap["requiredParties"].Value.GetList()
		require.NotNil(t, partiesSet)
		require.Len(t, partiesSet.Elements, 3)
		require.Equal(t, "alice", partiesSet.Elements[0].GetText())
		require.Equal(t, "bob", partiesSet.Elements[1].GetText())
		require.Equal(t, "charlie", partiesSet.Elements[2].GetText())

		coordinateRecord := fieldMap["coordinate"].Value.GetRecord()
		require.NotNil(t, coordinateRecord)
		require.Len(t, coordinateRecord.Fields, 2)
		coordFieldMap := make(map[string]*v2.RecordField)
		for _, field := range coordinateRecord.Fields {
			coordFieldMap[field.Label] = field
		}
		require.Equal(t, int64(100), coordFieldMap["_1"].Value.GetInt64())
		require.Equal(t, int64(200), coordFieldMap["_2"].Value.GetInt64())

		metadataRecord := fieldMap["metadata"].Value.GetRecord()
		require.NotNil(t, metadataRecord)
		require.Len(t, metadataRecord.Fields, 2)
		metaFieldMap := make(map[string]*v2.RecordField)
		for _, field := range metadataRecord.Fields {
			metaFieldMap[field.Label] = field
		}
		require.Equal(t, "version", metaFieldMap["_1"].Value.GetText())
		require.Equal(t, "1.0.0", metaFieldMap["_2"].Value.GetText())
	})
}

func TestConvertToRecordGENMAP(t *testing.T) {
	t.Run("GENMAP with string values converts to GenMap", func(t *testing.T) {
		genMapData := map[string]interface{}{
			"_type": "genmap",
			"value": map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
		}

		data := make(map[string]interface{})
		data["genmap"] = genMapData

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "genmap", record.Fields[0].Label)

		genMap := record.Fields[0].Value.GetGenMap()
		require.NotNil(t, genMap)
		require.Len(t, genMap.Entries, 3)

		entriesMap := make(map[string]string)
		for _, entry := range genMap.Entries {
			key := entry.Key.GetText()
			entriesMap[key] = entry.Value.GetText()
		}

		require.Equal(t, "value1", entriesMap["key1"])
		require.Equal(t, "value2", entriesMap["key2"])
		require.Equal(t, "value3", entriesMap["key3"])
	})

	t.Run("GENMAP with mixed types converts to GenMap", func(t *testing.T) {
		genMapData := map[string]interface{}{
			"_type": "genmap",
			"value": map[string]interface{}{
				"text":   "hello",
				"number": int64(42),
				"bool":   true,
			},
		}

		data := make(map[string]interface{})
		data["genmap"] = genMapData

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "genmap", record.Fields[0].Label)

		genMap := record.Fields[0].Value.GetGenMap()
		require.NotNil(t, genMap)
		require.Len(t, genMap.Entries, 3)

		entriesMap := make(map[string]*v2.Value)
		for _, entry := range genMap.Entries {
			key := entry.Key.GetText()
			entriesMap[key] = entry.Value
		}

		require.Equal(t, "hello", entriesMap["text"].GetText())
		require.Equal(t, int64(42), entriesMap["number"].GetInt64())
		require.Equal(t, true, entriesMap["bool"].GetBool())
	})

	t.Run("GENMAP with DAML types", func(t *testing.T) {
		genMapData := map[string]interface{}{
			"_type": "genmap",
			"value": map[string]interface{}{
				"party":   types.PARTY("alice"),
				"int":     types.INT64(100),
				"text":    types.TEXT("test"),
				"bool":    types.BOOL(false),
				"numeric": types.NUMERIC(big.NewInt(500)),
			},
		}

		data := make(map[string]interface{})
		data["genmap"] = genMapData

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "genmap", record.Fields[0].Label)

		genMap := record.Fields[0].Value.GetGenMap()
		require.NotNil(t, genMap)
		require.Len(t, genMap.Entries, 5)

		entriesMap := make(map[string]*v2.Value)
		for _, entry := range genMap.Entries {
			key := entry.Key.GetText()
			entriesMap[key] = entry.Value
		}

		require.Equal(t, "alice", entriesMap["party"].GetParty())
		require.Equal(t, int64(100), entriesMap["int"].GetInt64())
		require.Equal(t, "test", entriesMap["text"].GetText())
		require.Equal(t, false, entriesMap["bool"].GetBool())
		require.Equal(t, "0.0000000500", entriesMap["numeric"].GetNumeric())
	})

	t.Run("GENMAP with empty map", func(t *testing.T) {
		genMapData := map[string]interface{}{
			"_type": "genmap",
			"value": map[string]interface{}{},
		}

		data := make(map[string]interface{})
		data["emptyMap"] = genMapData

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "emptyMap", record.Fields[0].Label)

		genMap := record.Fields[0].Value.GetGenMap()
		require.NotNil(t, genMap)
		require.Len(t, genMap.Entries, 0)
	})

	t.Run("GENMAP in struct", func(t *testing.T) {
		type TestGenmapStruct struct {
			Owner     types.PARTY `json:"owner"`
			StringMap interface{} `json:"stringMap"`
			MixedMap  interface{} `json:"mixedMap"`
			Name      types.TEXT  `json:"name"`
		}

		testData := TestGenmapStruct{
			Owner: types.PARTY("alice"),
			StringMap: map[string]interface{}{
				"_type": "genmap",
				"value": map[string]interface{}{
					"a": "alpha",
					"b": "beta",
				},
			},
			MixedMap: map[string]interface{}{
				"_type": "genmap",
				"value": map[string]interface{}{
					"count": int64(10),
					"name":  "test",
				},
			},
			Name: types.TEXT("genmap test"),
		}

		data := make(map[string]interface{})
		data["testStruct"] = testData

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)

		structRecord := record.Fields[0].Value.GetRecord()
		require.NotNil(t, structRecord)
		require.Len(t, structRecord.Fields, 4)

		fieldMap := make(map[string]*v2.RecordField)
		for _, field := range structRecord.Fields {
			fieldMap[field.Label] = field
		}

		require.Equal(t, "alice", fieldMap["owner"].Value.GetParty())
		require.Equal(t, "genmap test", fieldMap["name"].Value.GetText())

		stringGenMap := fieldMap["stringMap"].Value.GetGenMap()
		require.NotNil(t, stringGenMap)
		require.Len(t, stringGenMap.Entries, 2)

		mixedGenMap := fieldMap["mixedMap"].Value.GetGenMap()
		require.NotNil(t, mixedGenMap)
		require.Len(t, mixedGenMap.Entries, 2)
	})

	t.Run("GENMAP with types.GENMAP type", func(t *testing.T) {
		genMapValue := types.GENMAP{
			"name":   "alice",
			"age":    int64(30),
			"active": true,
		}

		genMapData := map[string]interface{}{
			"_type": "genmap",
			"value": genMapValue,
		}

		data := make(map[string]interface{})
		data["genmap"] = genMapData

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "genmap", record.Fields[0].Label)

		genMap := record.Fields[0].Value.GetGenMap()
		require.NotNil(t, genMap)
		require.Len(t, genMap.Entries, 3)

		entriesMap := make(map[string]*v2.Value)
		for _, entry := range genMap.Entries {
			key := entry.Key.GetText()
			entriesMap[key] = entry.Value
		}

		require.Equal(t, "alice", entriesMap["name"].GetText())
		require.Equal(t, int64(30), entriesMap["age"].GetInt64())
		require.Equal(t, true, entriesMap["active"].GetBool())
	})
}

func TestConvertToRecordOptionalWithMaps(t *testing.T) {
	t.Run("Optional with GenMap", func(t *testing.T) {
		optionalData := map[string]interface{}{
			"_type": "optional",
			"value": map[string]interface{}{
				"_type": "genmap",
				"value": map[string]interface{}{
					"count": int64(5),
					"name":  "test",
				},
			},
		}

		data := make(map[string]interface{})
		data["optionalMap"] = optionalData

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "optionalMap", record.Fields[0].Label)

		optionalValue := record.Fields[0].Value.GetOptional()
		require.NotNil(t, optionalValue)
		require.NotNil(t, optionalValue.Value)

		genMap := optionalValue.Value.GetGenMap()
		require.NotNil(t, genMap)
		require.Len(t, genMap.Entries, 2)

		entriesMap := make(map[string]*v2.Value)
		for _, entry := range genMap.Entries {
			key := entry.Key.GetText()
			entriesMap[key] = entry.Value
		}

		require.Equal(t, int64(5), entriesMap["count"].GetInt64())
		require.Equal(t, "test", entriesMap["name"].GetText())
	})

	t.Run("Optional with TextMap", func(t *testing.T) {
		textMapValue := map[string]string{
			"key1": "value1",
			"key2": "value2",
		}

		optionalData := map[string]interface{}{
			"_type": "optional",
			"value": map[string]interface{}{
				"_type": "textmap",
				"value": textMapValue,
			},
		}

		data := make(map[string]interface{})
		data["optionalTextMap"] = optionalData

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "optionalTextMap", record.Fields[0].Label)

		optionalValue := record.Fields[0].Value.GetOptional()
		require.NotNil(t, optionalValue)
		require.NotNil(t, optionalValue.Value)

		textMap := optionalValue.Value.GetTextMap()
		require.NotNil(t, textMap)
		require.Len(t, textMap.Entries, 2)

		entriesMap := make(map[string]string)
		for _, entry := range textMap.Entries {
			entriesMap[entry.Key] = entry.Value.GetText()
		}

		require.Equal(t, "value1", entriesMap["key1"])
		require.Equal(t, "value2", entriesMap["key2"])
	})

	t.Run("Optional without map value (None)", func(t *testing.T) {
		optionalData := map[string]interface{}{
			"_type": "optional",
		}

		data := make(map[string]interface{})
		data["optionalEmptyMap"] = optionalData

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "optionalEmptyMap", record.Fields[0].Label)

		optionalValue := record.Fields[0].Value.GetOptional()
		require.NotNil(t, optionalValue)
		require.Nil(t, optionalValue.Value)
	})

	t.Run("Multiple optional maps in struct", func(t *testing.T) {
		type TestOptionalMapsStruct struct {
			Name           types.TEXT  `json:"name"`
			MaybeStringMap interface{} `json:"maybeStringMap"`
			MaybeMixedMap  interface{} `json:"maybeMixedMap"`
			MaybeEmptyMap  interface{} `json:"maybeEmptyMap"`
		}

		testData := TestOptionalMapsStruct{
			Name: types.TEXT("optional maps test"),
			MaybeStringMap: map[string]interface{}{
				"_type": "optional",
				"value": map[string]interface{}{
					"_type": "genmap",
					"value": map[string]interface{}{
						"a": "alpha",
						"b": "beta",
					},
				},
			},
			MaybeMixedMap: map[string]interface{}{
				"_type": "optional",
				"value": map[string]interface{}{
					"_type": "genmap",
					"value": map[string]interface{}{
						"count": int64(100),
						"name":  "test",
					},
				},
			},
			MaybeEmptyMap: map[string]interface{}{
				"_type": "optional",
			},
		}

		data := make(map[string]interface{})
		data["testStruct"] = testData

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)

		structRecord := record.Fields[0].Value.GetRecord()
		require.NotNil(t, structRecord)
		require.Len(t, structRecord.Fields, 4)

		fieldMap := make(map[string]*v2.RecordField)
		for _, field := range structRecord.Fields {
			fieldMap[field.Label] = field
		}

		require.Equal(t, "optional maps test", fieldMap["name"].Value.GetText())

		optStringMap := fieldMap["maybeStringMap"].Value.GetOptional()
		require.NotNil(t, optStringMap)
		require.NotNil(t, optStringMap.Value)
		stringGenMap := optStringMap.Value.GetGenMap()
		require.NotNil(t, stringGenMap)
		require.Len(t, stringGenMap.Entries, 2)

		optMixedMap := fieldMap["maybeMixedMap"].Value.GetOptional()
		require.NotNil(t, optMixedMap)
		require.NotNil(t, optMixedMap.Value)
		mixedGenMap := optMixedMap.Value.GetGenMap()
		require.NotNil(t, mixedGenMap)
		require.Len(t, mixedGenMap.Entries, 2)

		optEmptyMap := fieldMap["maybeEmptyMap"].Value.GetOptional()
		require.NotNil(t, optEmptyMap)
		require.Nil(t, optEmptyMap.Value)
	})
}

func TestConvertToRecordNestedMaps(t *testing.T) {
	t.Run("GenMap with nested maps", func(t *testing.T) {
		nestedMapData := map[string]interface{}{
			"_type": "genmap",
			"value": map[string]interface{}{
				"config": map[string]interface{}{
					"_type": "genmap",
					"value": map[string]interface{}{
						"timeout": int64(30),
						"retry":   int64(3),
					},
				},
				"metadata": map[string]interface{}{
					"_type": "genmap",
					"value": map[string]interface{}{
						"version": "2.0",
						"env":     "staging",
					},
				},
			},
		}

		data := make(map[string]interface{})
		data["nestedMap"] = nestedMapData

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "nestedMap", record.Fields[0].Label)

		genMap := record.Fields[0].Value.GetGenMap()
		require.NotNil(t, genMap)
		require.Len(t, genMap.Entries, 2)

		entriesMap := make(map[string]*v2.Value)
		for _, entry := range genMap.Entries {
			key := entry.Key.GetText()
			entriesMap[key] = entry.Value
		}

		configGenMap := entriesMap["config"].GetGenMap()
		require.NotNil(t, configGenMap)
		require.Len(t, configGenMap.Entries, 2)

		metadataGenMap := entriesMap["metadata"].GetGenMap()
		require.NotNil(t, metadataGenMap)
		require.Len(t, metadataGenMap.Entries, 2)
	})

	t.Run("Map with complex nested structure", func(t *testing.T) {
		complexMapData := map[string]interface{}{
			"_type": "genmap",
			"value": map[string]interface{}{
				"list": []interface{}{int64(1), int64(2), int64(3)},
				"map": map[string]interface{}{
					"_type": "genmap",
					"value": map[string]interface{}{
						"inner": "value",
					},
				},
				"simple": "text",
			},
		}

		data := make(map[string]interface{})
		data["complexMap"] = complexMapData

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "complexMap", record.Fields[0].Label)

		genMap := record.Fields[0].Value.GetGenMap()
		require.NotNil(t, genMap)
		require.Len(t, genMap.Entries, 3)

		entriesMap := make(map[string]*v2.Value)
		for _, entry := range genMap.Entries {
			key := entry.Key.GetText()
			entriesMap[key] = entry.Value
		}

		listValue := entriesMap["list"].GetList()
		require.NotNil(t, listValue)
		require.Len(t, listValue.Elements, 3)
		require.Equal(t, int64(1), listValue.Elements[0].GetInt64())
		require.Equal(t, int64(2), listValue.Elements[1].GetInt64())
		require.Equal(t, int64(3), listValue.Elements[2].GetInt64())

		innerMap := entriesMap["map"].GetGenMap()
		require.NotNil(t, innerMap)
		require.Len(t, innerMap.Entries, 1)

		require.Equal(t, "text", entriesMap["simple"].GetText())
	})

	t.Run("Map in Optional in Map", func(t *testing.T) {
		complexMapData := map[string]interface{}{
			"_type": "genmap",
			"value": map[string]interface{}{
				"optionalMap": map[string]interface{}{
					"_type": "optional",
					"value": map[string]interface{}{
						"_type": "genmap",
						"value": map[string]interface{}{
							"key": "value",
						},
					},
				},
			},
		}

		data := make(map[string]interface{})
		data["complexMap"] = complexMapData

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "complexMap", record.Fields[0].Label)

		genMap := record.Fields[0].Value.GetGenMap()
		require.NotNil(t, genMap)
		require.Len(t, genMap.Entries, 1)

		entriesMap := make(map[string]*v2.Value)
		for _, entry := range genMap.Entries {
			key := entry.Key.GetText()
			entriesMap[key] = entry.Value
		}

		optionalValue := entriesMap["optionalMap"].GetOptional()
		require.NotNil(t, optionalValue)
		require.NotNil(t, optionalValue.Value)

		innerGenMap := optionalValue.Value.GetGenMap()
		require.NotNil(t, innerGenMap)
		require.Len(t, innerGenMap.Entries, 1)

		innerEntries := make(map[string]*v2.Value)
		for _, entry := range innerGenMap.Entries {
			key := entry.Key.GetText()
			innerEntries[key] = entry.Value
		}
		require.Equal(t, "value", innerEntries["key"].GetText())
	})
}

func TestMapToValue_TIMESTAMP(t *testing.T) {
	tests := []struct {
		name     string
		input    types.TIMESTAMP
		expected int64
	}{
		{
			name:     "epoch timestamp",
			input:    types.TIMESTAMP(time.Unix(0, 0)),
			expected: 0,
		},
		{
			name:     "standard timestamp",
			input:    types.TIMESTAMP(time.Date(2024, 1, 15, 12, 30, 45, 123456000, time.UTC)),
			expected: 1705321845123456,
		},
		{
			name:     "timestamp with microseconds",
			input:    types.TIMESTAMP(time.Date(2025, 6, 1, 10, 20, 30, 500000000, time.UTC)),
			expected: 1748773230500000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapToValue(tt.input)

			require.NotNil(t, result)
			timestampValue, ok := result.Sum.(*v2.Value_Timestamp)
			require.True(t, ok)
			require.Equal(t, tt.expected, timestampValue.Timestamp)
		})
	}
}

func TestValueFromProto_TIMESTAMP(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected time.Time
	}{
		{
			name:     "epoch timestamp",
			input:    0,
			expected: time.Unix(0, 0),
		},
		{
			name:     "standard timestamp",
			input:    1705321845123456,
			expected: time.Date(2024, 1, 15, 12, 30, 45, 123456000, time.UTC),
		},
		{
			name:     "timestamp with microseconds",
			input:    1748773230500000,
			expected: time.Date(2025, 6, 1, 10, 20, 30, 500000000, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			protoValue := &v2.Value{
				Sum: &v2.Value_Timestamp{
					Timestamp: tt.input,
				},
			}

			result := valueFromProto(protoValue)

			require.NotNil(t, result)
			timestamp, ok := result.(time.Time)
			require.True(t, ok)
			require.Equal(t, tt.expected.Unix(), timestamp.Unix())
			require.Equal(t, tt.expected.Nanosecond(), timestamp.Nanosecond())
		})
	}
}

func TestTimestampRoundTrip(t *testing.T) {
	tests := []struct {
		name      string
		timestamp types.TIMESTAMP
	}{
		{
			name:      "epoch",
			timestamp: types.TIMESTAMP(time.Unix(0, 0)),
		},
		{
			name:      "recent date",
			timestamp: types.TIMESTAMP(time.Date(2024, 1, 15, 12, 30, 45, 123456000, time.UTC)),
		},
		{
			name:      "future date",
			timestamp: types.TIMESTAMP(time.Date(2030, 12, 31, 23, 59, 59, 999999000, time.UTC)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			protoValue := mapToValue(tt.timestamp)
			result := valueFromProto(protoValue)

			originalTime := time.Time(tt.timestamp)
			resultTime, ok := result.(time.Time)
			require.True(t, ok)
			require.Equal(t, originalTime.Unix(), resultTime.Unix())
			require.Equal(t, originalTime.Nanosecond()/1000*1000, resultTime.Nanosecond())
		})
	}
}

func TestConvertToRecordTIMESTAMP(t *testing.T) {
	t.Run("TIMESTAMP type conversion - epoch", func(t *testing.T) {
		timestampValue := types.TIMESTAMP(time.Unix(0, 0))

		data := make(map[string]interface{})
		data["timestamp"] = timestampValue

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "timestamp", record.Fields[0].Label)

		int64Value := record.Fields[0].Value.GetTimestamp()
		require.Equal(t, int64(0), int64Value)
	})

	t.Run("TIMESTAMP type conversion - standard date", func(t *testing.T) {
		timestampValue := types.TIMESTAMP(time.Date(2024, 1, 15, 12, 30, 45, 123456000, time.UTC))

		data := make(map[string]interface{})
		data["timestamp"] = timestampValue

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "timestamp", record.Fields[0].Label)

		int64Value := record.Fields[0].Value.GetTimestamp()
		require.Equal(t, int64(1705321845123456), int64Value)
	})

	t.Run("TIMESTAMP type conversion - with microseconds", func(t *testing.T) {
		timestampValue := types.TIMESTAMP(time.Date(2025, 6, 1, 10, 20, 30, 500000000, time.UTC))

		data := make(map[string]interface{})
		data["timestamp"] = timestampValue

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)
		require.Equal(t, "timestamp", record.Fields[0].Label)

		int64Value := record.Fields[0].Value.GetTimestamp()
		require.Equal(t, int64(1748773230500000), int64Value)
	})

	t.Run("TIMESTAMP in struct", func(t *testing.T) {
		type TestTimestampStruct struct {
			Owner     types.PARTY     `json:"owner"`
			CreatedAt types.TIMESTAMP `json:"createdAt"`
			UpdatedAt types.TIMESTAMP `json:"updatedAt"`
			Name      types.TEXT      `json:"name"`
		}

		testData := TestTimestampStruct{
			Owner:     types.PARTY("alice"),
			CreatedAt: types.TIMESTAMP(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
			UpdatedAt: types.TIMESTAMP(time.Date(2024, 12, 31, 23, 59, 59, 999999000, time.UTC)),
			Name:      types.TEXT("test timestamp"),
		}

		data := make(map[string]interface{})
		data["testStruct"] = testData

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)

		structRecord := record.Fields[0].Value.GetRecord()
		require.NotNil(t, structRecord)
		require.Len(t, structRecord.Fields, 4)

		fieldMap := make(map[string]*v2.RecordField)
		for _, field := range structRecord.Fields {
			fieldMap[field.Label] = field
		}

		require.Equal(t, "alice", fieldMap["owner"].Value.GetParty())
		require.Equal(t, "test timestamp", fieldMap["name"].Value.GetText())
		require.Equal(t, int64(1704067200000000), fieldMap["createdAt"].Value.GetTimestamp())
		require.Equal(t, int64(1735689599999999), fieldMap["updatedAt"].Value.GetTimestamp())
	})

	t.Run("Multiple TIMESTAMP fields", func(t *testing.T) {
		type TestMultiTimestampStruct struct {
			StartTime types.TIMESTAMP `json:"startTime"`
			EndTime   types.TIMESTAMP `json:"endTime"`
			LastCheck types.TIMESTAMP `json:"lastCheck"`
		}

		testData := TestMultiTimestampStruct{
			StartTime: types.TIMESTAMP(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
			EndTime:   types.TIMESTAMP(time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)),
			LastCheck: types.TIMESTAMP(time.Date(2024, 6, 15, 12, 30, 0, 0, time.UTC)),
		}

		data := make(map[string]interface{})
		data["testStruct"] = testData

		record := convertToRecord(data)

		require.NotNil(t, record)
		require.Len(t, record.Fields, 1)

		structRecord := record.Fields[0].Value.GetRecord()
		require.NotNil(t, structRecord)
		require.Len(t, structRecord.Fields, 3)

		fieldMap := make(map[string]*v2.RecordField)
		for _, field := range structRecord.Fields {
			fieldMap[field.Label] = field
		}

		require.Equal(t, int64(1704067200000000), fieldMap["startTime"].Value.GetTimestamp())
		require.Equal(t, int64(1735689599000000), fieldMap["endTime"].Value.GetTimestamp())
		require.Equal(t, int64(1718454600000000), fieldMap["lastCheck"].Value.GetTimestamp())
	})
}
