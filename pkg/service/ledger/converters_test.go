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
			TemplateID: "test:template:OneOfEverything",
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
