package codec

import (
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"time"

	"github.com/noders-team/go-daml/pkg/types"
)

// JsonCodec follows the transcode JsonCodec pattern for encoding DAML values to JSON
type JsonCodec struct {
	// encodeNumericAsString controls whether numeric values are encoded as strings (true) or numbers (false)
	// The latter might be useful for querying and mathematical operations, but can lose precision due to float point errors
	EncodeNumericAsString bool

	// encodeInt64AsString controls whether int64 values are encoded as strings (true) or numbers (false)
	// The latter might be useful for querying and mathematical operations, but can lose precision,
	// as numbers in some json implementations are backed by Double
	EncodeInt64AsString bool

	// excludeNullValuesInRecords controls whether fields with null values in records are excluded from JSON (true)
	// or included with a null value (false)
	ExcludeNullValuesInRecords bool
}

func isTuple2(v reflect.Value) bool {
	if v.Kind() == reflect.Struct && v.NumField() == 2 {
		t := v.Type()
		return t.Field(0).Name == "First" && t.Field(1).Name == "Second"
	}
	return false
}

func isTuple3(v reflect.Value) bool {
	if v.Kind() == reflect.Struct && v.NumField() == 3 {
		t := v.Type()
		return t.Field(0).Name == "First" && t.Field(1).Name == "Second" && t.Field(2).Name == "Third"
	}
	return false
}

// NewJsonCodec creates a new JsonCodec with default settings following transcode patterns
func NewJsonCodec() *JsonCodec {
	return &JsonCodec{
		EncodeNumericAsString:      true,  // Default: encode numeric as string for precision
		EncodeInt64AsString:        true,  // Default: encode int64 as string for precision
		ExcludeNullValuesInRecords: false, // Default: include null values
	}
}

// NewJsonCodecWithOptions creates a JsonCodec with custom options
func NewJsonCodecWithOptions(encodeNumericAsString, encodeInt64AsString, excludeNullValues bool) *JsonCodec {
	return &JsonCodec{
		EncodeNumericAsString:      encodeNumericAsString,
		EncodeInt64AsString:        encodeInt64AsString,
		ExcludeNullValuesInRecords: excludeNullValues,
	}
}

// Marshall converts a DAML structure to JSON bytes following transcode patterns
func (codec *JsonCodec) Marshall(value interface{}) ([]byte, error) {
	jsonValue, err := codec.toDynamicValue(value)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to dynamic value: %w", err)
	}
	return json.Marshal(jsonValue)
}

// Unmarshall converts JSON bytes back to a DAML structure following transcode patterns
func (codec *JsonCodec) Unmarshall(data []byte, target interface{}) error {
	var intermediate interface{}
	if err := json.Unmarshal(data, &intermediate); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return codec.fromDynamicValue(intermediate, target)
}

// toDynamicValue converts Go values to JSON-compatible values following transcode codec patterns
func (codec *JsonCodec) toDynamicValue(value interface{}) (interface{}, error) {
	if value == nil {
		return nil, nil
	}

	switch v := value.(type) {
	case types.PARTY:
		return codec.partyToDynamicValue(v), nil
	case types.TEXT:
		return codec.textToDynamicValue(v), nil
	case types.INT64:
		return codec.int64ToDynamicValue(v), nil
	case types.BOOL:
		return codec.boolToDynamicValue(v), nil
	case types.NUMERIC:
		return codec.numericToDynamicValue(v), nil
	case types.DECIMAL:
		return codec.decimalToDynamicValue(v), nil
	case types.TIMESTAMP:
		return codec.timestampToDynamicValue(v), nil
	case types.DATE:
		return codec.dateToDynamicValue(v), nil
	case types.UNIT:
		return codec.unitToDynamicValue(v), nil
	case types.CONTRACT_ID:
		return codec.contractIdToDynamicValue(v), nil
	case types.RELTIME:
		return codec.reltimeToDynamicValue(v), nil
	case types.SET:
		return codec.setToDynamicValue(v)
	case types.GENMAP:
		return codec.genMapToDynamicValue(v)
	case types.TEXTMAP:
		return codec.textMapToDynamicValue(v)
	case types.MAP:
		return codec.mapToDynamicValue(v)
	case types.LIST:
		return codec.listToDynamicValue(v)
	}

	// Handle interface types
	if variant, ok := value.(types.VARIANT); ok {
		return codec.variantToDynamicValue(variant)
	}
	if enum, ok := value.(types.ENUM); ok {
		return codec.enumToDynamicValue(enum), nil
	}

	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil, nil
		}
		return codec.toDynamicValue(rv.Elem().Interface())
	}

	if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
		return codec.listToDynamicValueFromReflect(rv)
	}

	if isTuple2(rv) {
		return codec.tuple2ToDynamicValueFromReflect(rv)
	}
	if isTuple3(rv) {
		return codec.tuple3ToDynamicValueFromReflect(rv)
	}

	if rv.Kind() == reflect.Struct {
		return codec.recordToDynamicValue(value)
	}

	if rv.Kind() == reflect.Map {
		return codec.mapToDynamicValueFromReflect(rv)
	}

	switch v := value.(type) {
	case string:
		return v, nil
	case int, int8, int16, int32, int64:
		if codec.EncodeInt64AsString {
			return fmt.Sprintf("%d", v), nil
		}
		return v, nil
	case uint, uint8, uint16, uint32, uint64:
		if codec.EncodeInt64AsString {
			return fmt.Sprintf("%d", v), nil
		}
		return v, nil
	case float32, float64:
		return v, nil
	case bool:
		return v, nil
	default:
		return v, nil
	}
}

// DAML type converters following transcode patterns

func (codec *JsonCodec) partyToDynamicValue(party types.PARTY) string {
	return string(party)
}

func (codec *JsonCodec) textToDynamicValue(text types.TEXT) string {
	return string(text)
}

func (codec *JsonCodec) int64ToDynamicValue(i types.INT64) interface{} {
	if codec.EncodeInt64AsString {
		return fmt.Sprintf("%d", int64(i))
	}
	return int64(i)
}

func (codec *JsonCodec) boolToDynamicValue(b types.BOOL) bool {
	return bool(b)
}

func (codec *JsonCodec) numericToDynamicValue(n types.NUMERIC) interface{} {
	return codec.bigIntToDynamicValue((*big.Int)(n))
}

func (codec *JsonCodec) decimalToDynamicValue(d types.DECIMAL) interface{} {
	return codec.bigIntToDynamicValue((*big.Int)(d))
}

func (codec *JsonCodec) bigIntToDynamicValue(bi *big.Int) interface{} {
	if bi == nil {
		return nil
	}
	if codec.EncodeNumericAsString {
		return bi.String()
	}
	f, _ := new(big.Float).SetInt(bi).Float64()
	return f
}

func (codec *JsonCodec) timestampToDynamicValue(t types.TIMESTAMP) string {
	return time.Time(t).Format("2006-01-02T15:04:05.000000Z")
}

func (codec *JsonCodec) dateToDynamicValue(d types.DATE) string {
	return time.Time(d).Format("2006-01-02")
}

func (codec *JsonCodec) unitToDynamicValue(_ types.UNIT) map[string]interface{} {
	return map[string]interface{}{}
}

func (codec *JsonCodec) contractIdToDynamicValue(c types.CONTRACT_ID) string {
	return string(c)
}

func (codec *JsonCodec) reltimeToDynamicValue(r types.RELTIME) interface{} {
	microseconds := int64(time.Duration(r) / time.Microsecond)
	if codec.EncodeInt64AsString {
		return fmt.Sprintf("%d", microseconds)
	}
	return microseconds
}

func (codec *JsonCodec) setToDynamicValue(s types.SET) (interface{}, error) {
	result := make([]interface{}, len(s))
	for i, v := range s {
		converted, err := codec.toDynamicValue(v)
		if err != nil {
			return nil, err
		}
		result[i] = converted
	}
	return result, nil
}

func (codec *JsonCodec) tuple2ToDynamicValueFromReflect(rv reflect.Value) (interface{}, error) {
	first, err := codec.toDynamicValue(rv.Field(0).Interface())
	if err != nil {
		return nil, err
	}
	second, err := codec.toDynamicValue(rv.Field(1).Interface())
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"_1": first,
		"_2": second,
	}, nil
}

func (codec *JsonCodec) tuple3ToDynamicValueFromReflect(rv reflect.Value) (interface{}, error) {
	first, err := codec.toDynamicValue(rv.Field(0).Interface())
	if err != nil {
		return nil, err
	}
	second, err := codec.toDynamicValue(rv.Field(1).Interface())
	if err != nil {
		return nil, err
	}
	third, err := codec.toDynamicValue(rv.Field(2).Interface())
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"_1": first,
		"_2": second,
		"_3": third,
	}, nil
}

func (codec *JsonCodec) genMapToDynamicValue(gm types.GENMAP) (interface{}, error) {
	return codec.mapToDynamicValueGeneric(gm)
}

func (codec *JsonCodec) textMapToDynamicValue(tm types.TEXTMAP) (interface{}, error) {
	result := make(map[string]string)
	for k, v := range tm {
		result[k] = v
	}
	return result, nil
}

func (codec *JsonCodec) mapToDynamicValue(m types.MAP) (interface{}, error) {
	return codec.mapToDynamicValueGeneric(m)
}

func (codec *JsonCodec) mapToDynamicValueGeneric(m interface{}) (interface{}, error) {
	result := make(map[string]interface{})
	switch v := m.(type) {
	case types.GENMAP:
		for k, val := range v {
			converted, err := codec.toDynamicValue(val)
			if err != nil {
				return nil, err
			}
			result[k] = converted
		}
	case types.MAP:
		for k, val := range v {
			converted, err := codec.toDynamicValue(val)
			if err != nil {
				return nil, err
			}
			result[k] = converted
		}
	}
	return result, nil
}

func (codec *JsonCodec) listToDynamicValue(l types.LIST) (interface{}, error) {
	result := make([]interface{}, len(l))
	for i, v := range l {
		converted, err := codec.toDynamicValue(v)
		if err != nil {
			return nil, err
		}
		result[i] = converted
	}
	return result, nil
}

func (codec *JsonCodec) listToDynamicValueFromReflect(rv reflect.Value) (interface{}, error) {
	length := rv.Len()
	result := make([]interface{}, length)
	for i := 0; i < length; i++ {
		converted, err := codec.toDynamicValue(rv.Index(i).Interface())
		if err != nil {
			return nil, err
		}
		result[i] = converted
	}
	return result, nil
}

func (codec *JsonCodec) mapToDynamicValueFromReflect(rv reflect.Value) (interface{}, error) {
	result := make(map[string]interface{})
	for _, key := range rv.MapKeys() {
		keyStr := fmt.Sprintf("%v", key.Interface())
		value := rv.MapIndex(key).Interface()
		converted, err := codec.toDynamicValue(value)
		if err != nil {
			return nil, err
		}
		result[keyStr] = converted
	}
	return result, nil
}

// recordToDynamicValue handles struct conversion following transcode record pattern
func (codec *JsonCodec) recordToDynamicValue(value interface{}) (interface{}, error) {
	rv := reflect.ValueOf(value)
	rt := reflect.TypeOf(value)

	if rv.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %v", rv.Kind())
	}

	result := make(map[string]interface{})

	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		fieldType := rt.Field(i)

		// Skip unexported fields
		if !field.CanInterface() {
			continue
		}

		// Get JSON tag or use field name
		jsonTag := fieldType.Tag.Get("json")
		fieldName := fieldType.Name
		if jsonTag != "" && jsonTag != "-" {
			// Parse json tag (handle omitempty, etc.)
			if commaIdx := strings.IndexByte(jsonTag, ','); commaIdx != -1 {
				fieldName = jsonTag[:commaIdx]
			} else {
				fieldName = jsonTag
			}
		}

		fieldValue := field.Interface()

		// Handle optional fields (nil pointers)
		if field.Kind() == reflect.Ptr && field.IsNil() {
			if !codec.ExcludeNullValuesInRecords {
				result[fieldName] = nil
			}
			continue
		}

		converted, err := codec.toDynamicValue(fieldValue)
		if err != nil {
			return nil, fmt.Errorf("failed to convert field %s: %w", fieldName, err)
		}

		// Exclude null values if configured
		if codec.ExcludeNullValuesInRecords && converted == nil {
			continue
		}

		result[fieldName] = converted
	}

	return result, nil
}

// Handle variant types (implements VARIANT interface)
func (codec *JsonCodec) variantToDynamicValue(variant types.VARIANT) (interface{}, error) {
	tag := variant.GetVariantTag()
	value := variant.GetVariantValue()

	if tag == "" {
		return map[string]interface{}{}, nil
	}

	converted, err := codec.toDynamicValue(value)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"tag":   tag,
		"value": converted,
	}, nil
}

// Handle enum types (implements ENUM interface)
func (codec *JsonCodec) enumToDynamicValue(enum types.ENUM) string {
	return enum.GetEnumConstructor()
}

// fromDynamicValue converts JSON-compatible values back to Go values following transcode codec patterns
func (codec *JsonCodec) fromDynamicValue(value interface{}, target interface{}) error {
	if value == nil && target == nil {
		return nil
	}

	rv := reflect.ValueOf(target)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("target must be a pointer, got %v", rv.Kind())
	}

	if rv.IsNil() {
		return fmt.Errorf("target pointer is nil")
	}

	elem := rv.Elem()
	return codec.assignValue(value, elem)
}

// assignValue assigns a JSON value to a reflect.Value following DAML type patterns
func (codec *JsonCodec) assignValue(jsonValue interface{}, target reflect.Value) error {
	if jsonValue == nil {
		if target.Kind() == reflect.Ptr {
			target.Set(reflect.Zero(target.Type()))
			return nil
		}
		if target.Kind() == reflect.Struct {
			target.Set(reflect.Zero(target.Type()))
			return nil
		}
		if target.Kind() == reflect.String {
			target.Set(reflect.Zero(target.Type()))
			return nil
		}
		return fmt.Errorf("cannot assign nil to non-pointer type %v", target.Type())
	}

	targetType := target.Type()

	switch targetType {
	case reflect.TypeOf(types.PARTY("")):
		if str, ok := jsonValue.(string); ok {
			target.Set(reflect.ValueOf(types.PARTY(str)))
			return nil
		}
		return fmt.Errorf("expected string for PARTY, got %T", jsonValue)

	case reflect.TypeOf(types.TEXT("")):
		if str, ok := jsonValue.(string); ok {
			target.Set(reflect.ValueOf(types.TEXT(str)))
			return nil
		}
		return fmt.Errorf("expected string for TEXT, got %T", jsonValue)

	case reflect.TypeOf(types.INT64(0)):
		return codec.assignInt64Value(jsonValue, target)

	case reflect.TypeOf(types.BOOL(false)):
		if b, ok := jsonValue.(bool); ok {
			target.Set(reflect.ValueOf(types.BOOL(b)))
			return nil
		}
		return fmt.Errorf("expected bool for BOOL, got %T", jsonValue)

	case reflect.TypeOf(types.NUMERIC(nil)):
		return codec.assignNumericValue(jsonValue, target)

	case reflect.TypeOf(types.DECIMAL(nil)):
		return codec.assignDecimalValue(jsonValue, target)

	case reflect.TypeOf(types.TIMESTAMP{}):
		return codec.assignTimestampValue(jsonValue, target)

	case reflect.TypeOf(types.DATE{}):
		return codec.assignDateValue(jsonValue, target)

	case reflect.TypeOf(types.UNIT{}):
		target.Set(reflect.ValueOf(types.UNIT{}))
		return nil

	case reflect.TypeOf(types.CONTRACT_ID("")):
		if str, ok := jsonValue.(string); ok {
			target.Set(reflect.ValueOf(types.CONTRACT_ID(str)))
			return nil
		}
		return fmt.Errorf("expected string for CONTRACT_ID, got %T", jsonValue)

	case reflect.TypeOf(types.RELTIME(0)):
		return codec.assignReltimeValue(jsonValue, target)

	case reflect.TypeOf(types.SET{}):
		return codec.assignSetValue(jsonValue, target)

	case reflect.TypeOf(types.GENMAP{}):
		return codec.assignGenMapValue(jsonValue, target)

	case reflect.TypeOf(types.TEXTMAP{}):
		return codec.assignTextMapValue(jsonValue, target)

	case reflect.TypeOf(types.MAP{}):
		return codec.assignMapValue(jsonValue, target)

	case reflect.TypeOf(types.LIST{}):
		return codec.assignListValue(jsonValue, target)
	}
	if target.Kind() == reflect.Ptr {
		if jsonValue == nil {
			target.Set(reflect.Zero(target.Type()))
			return nil
		}

		newElem := reflect.New(target.Type().Elem())
		if err := codec.assignValue(jsonValue, newElem.Elem()); err != nil {
			return err
		}
		target.Set(newElem)
		return nil
	}

	if target.Kind() == reflect.Slice {
		return codec.assignSliceValue(jsonValue, target)
	}

	if isTuple2(target) {
		return codec.assignTuple2Value(jsonValue, target)
	}
	if isTuple3(target) {
		return codec.assignTuple3Value(jsonValue, target)
	}

	if target.Kind() == reflect.Struct {
		return codec.assignStructValue(jsonValue, target)
	}

	if target.Type().Implements(reflect.TypeOf((*types.VARIANT)(nil)).Elem()) {
		return codec.assignVariantValue(jsonValue, target)
	}

	if target.Type().Implements(reflect.TypeOf((*types.ENUM)(nil)).Elem()) {
		return codec.assignEnumValue(jsonValue, target)
	}

	switch target.Kind() {
	case reflect.String:
		if str, ok := jsonValue.(string); ok {
			target.SetString(str)
			return nil
		}
		return fmt.Errorf("expected string, got %T", jsonValue)

	case reflect.Bool:
		if b, ok := jsonValue.(bool); ok {
			target.SetBool(b)
			return nil
		}
		return fmt.Errorf("expected bool, got %T", jsonValue)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return codec.assignIntValue(jsonValue, target)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return codec.assignUintValue(jsonValue, target)

	case reflect.Float32, reflect.Float64:
		return codec.assignFloatValue(jsonValue, target)

	case reflect.Interface:
		if target.Type().NumMethod() == 0 {
			target.Set(reflect.ValueOf(jsonValue))
			return nil
		}
		return fmt.Errorf("cannot assign to non-empty interface type: %v", target.Type())

	default:
		return fmt.Errorf("unsupported target type: %v", target.Type())
	}
}

// Type-specific assignment methods
func (codec *JsonCodec) assignInt64Value(jsonValue interface{}, target reflect.Value) error {
	switch v := jsonValue.(type) {
	case string:
		if i, err := parseIntFromString(v); err == nil {
			target.Set(reflect.ValueOf(types.INT64(i)))
			return nil
		}
		return fmt.Errorf("invalid string format for INT64: %s", v)
	case float64:
		target.Set(reflect.ValueOf(types.INT64(int64(v))))
		return nil
	case int64:
		target.Set(reflect.ValueOf(types.INT64(v)))
		return nil
	default:
		return fmt.Errorf("expected string or number for INT64, got %T", jsonValue)
	}
}

func (codec *JsonCodec) assignNumericValue(jsonValue interface{}, target reflect.Value) error {
	return codec.assignBigIntValue(jsonValue, target, "NUMERIC", func(bi *big.Int) reflect.Value {
		return reflect.ValueOf(types.NUMERIC(bi))
	})
}

func (codec *JsonCodec) assignDecimalValue(jsonValue interface{}, target reflect.Value) error {
	return codec.assignBigIntValue(jsonValue, target, "DECIMAL", func(bi *big.Int) reflect.Value {
		return reflect.ValueOf(types.DECIMAL(bi))
	})
}

func (codec *JsonCodec) assignBigIntValue(jsonValue interface{}, target reflect.Value, typeName string, converter func(*big.Int) reflect.Value) error {
	switch v := jsonValue.(type) {
	case string:
		if bi, ok := new(big.Int).SetString(v, 10); ok {
			target.Set(converter(bi))
			return nil
		}
		if rat, ok := new(big.Rat).SetString(v); ok {
			scaledInt := new(big.Int)
			scaledInt.Mul(rat.Num(), big.NewInt(10000000000))
			scaledInt.Div(scaledInt, rat.Denom())
			target.Set(converter(scaledInt))
			return nil
		}
		return fmt.Errorf("invalid string format for %s: %s", typeName, v)
	case float64:
		if v == float64(int64(v)) {
			bi := big.NewInt(int64(v))
			target.Set(converter(bi))
			return nil
		}
		rat := new(big.Rat).SetFloat64(v)
		scaledInt := new(big.Int)
		scaledInt.Mul(rat.Num(), big.NewInt(10000000000))
		scaledInt.Div(scaledInt, rat.Denom())
		target.Set(converter(scaledInt))
		return nil
	case int64:
		bi := big.NewInt(v)
		target.Set(converter(bi))
		return nil
	default:
		return fmt.Errorf("expected string or number for %s, got %T", typeName, jsonValue)
	}
}

func (codec *JsonCodec) assignTimestampValue(jsonValue interface{}, target reflect.Value) error {
	switch v := jsonValue.(type) {
	case string:
		if t, err := time.Parse("2006-01-02T15:04:05.000000Z", v); err == nil {
			target.Set(reflect.ValueOf(types.TIMESTAMP(t)))
			return nil
		}
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			target.Set(reflect.ValueOf(types.TIMESTAMP(t)))
			return nil
		}
		return fmt.Errorf("invalid timestamp format: %s", v)
	case float64:
		t := time.Unix(int64(v)/1000000, (int64(v)%1000000)*1000)
		target.Set(reflect.ValueOf(types.TIMESTAMP(t)))
		return nil
	case int64:
		t := time.Unix(v/1000000, (v%1000000)*1000)
		target.Set(reflect.ValueOf(types.TIMESTAMP(t)))
		return nil
	default:
		return fmt.Errorf("expected string or number for TIMESTAMP, got %T", jsonValue)
	}
}

func (codec *JsonCodec) assignDateValue(jsonValue interface{}, target reflect.Value) error {
	switch v := jsonValue.(type) {
	case string:
		if t, err := time.Parse("2006-01-02", v); err == nil {
			target.Set(reflect.ValueOf(types.DATE(t)))
			return nil
		}
		return fmt.Errorf("invalid date format: %s", v)
	case float64:
		t := time.Unix(int64(v)*86400, 0).UTC()
		target.Set(reflect.ValueOf(types.DATE(t)))
		return nil
	case int64:
		t := time.Unix(v*86400, 0).UTC()
		target.Set(reflect.ValueOf(types.DATE(t)))
		return nil
	case int32:
		t := time.Unix(int64(v)*86400, 0).UTC()
		target.Set(reflect.ValueOf(types.DATE(t)))
		return nil
	default:
		return fmt.Errorf("expected string or number for DATE, got %T", jsonValue)
	}
}

func (codec *JsonCodec) assignGenMapValue(jsonValue interface{}, target reflect.Value) error {
	if m, ok := jsonValue.(map[string]interface{}); ok {
		result := make(types.GENMAP)
		for k, v := range m {
			result[k] = v
		}
		target.Set(reflect.ValueOf(result))
		return nil
	}
	return fmt.Errorf("expected object for GENMAP, got %T", jsonValue)
}

func (codec *JsonCodec) assignTextMapValue(jsonValue interface{}, target reflect.Value) error {
	if m, ok := jsonValue.(map[string]interface{}); ok {
		result := make(types.TEXTMAP)
		for k, v := range m {
			if str, ok := v.(string); ok {
				result[k] = str
			} else {
				result[k] = fmt.Sprintf("%v", v)
			}
		}
		target.Set(reflect.ValueOf(result))
		return nil
	}
	return fmt.Errorf("expected object for TEXTMAP, got %T", jsonValue)
}

func (codec *JsonCodec) assignMapValue(jsonValue interface{}, target reflect.Value) error {
	if m, ok := jsonValue.(map[string]interface{}); ok {
		result := make(types.MAP)
		for k, v := range m {
			result[k] = v
		}
		target.Set(reflect.ValueOf(result))
		return nil
	}
	return fmt.Errorf("expected object for MAP, got %T", jsonValue)
}

func (codec *JsonCodec) assignListValue(jsonValue interface{}, target reflect.Value) error {
	if arr, ok := jsonValue.([]interface{}); ok {
		result := make(types.LIST, len(arr))
		for i, v := range arr {
			if str, ok := v.(string); ok {
				result[i] = str
			} else {
				result[i] = fmt.Sprintf("%v", v)
			}
		}
		target.Set(reflect.ValueOf(result))
		return nil
	}
	return fmt.Errorf("expected array for LIST, got %T", jsonValue)
}

func (codec *JsonCodec) assignReltimeValue(jsonValue interface{}, target reflect.Value) error {
	switch v := jsonValue.(type) {
	case string:
		if i, err := parseIntFromString(v); err == nil {
			target.Set(reflect.ValueOf(types.RELTIME(time.Duration(i) * time.Microsecond)))
			return nil
		}
		return fmt.Errorf("invalid string format for RELTIME: %s", v)
	case float64:
		target.Set(reflect.ValueOf(types.RELTIME(time.Duration(int64(v)) * time.Microsecond)))
		return nil
	case int64:
		target.Set(reflect.ValueOf(types.RELTIME(time.Duration(v) * time.Microsecond)))
		return nil
	default:
		return fmt.Errorf("expected string or number for RELTIME, got %T", jsonValue)
	}
}

func (codec *JsonCodec) assignSetValue(jsonValue interface{}, target reflect.Value) error {
	if arr, ok := jsonValue.([]interface{}); ok {
		result := make(types.SET, len(arr))
		for i, v := range arr {
			result[i] = v
		}
		target.Set(reflect.ValueOf(result))
		return nil
	}
	return fmt.Errorf("expected array for SET, got %T", jsonValue)
}

func (codec *JsonCodec) assignTuple2Value(jsonValue interface{}, target reflect.Value) error {
	if m, ok := jsonValue.(map[string]interface{}); ok {
		first, hasFirst := m["_1"]
		second, hasSecond := m["_2"]
		if !hasFirst || !hasSecond {
			return fmt.Errorf("TUPLE2 missing _1 or _2 fields")
		}
		if err := codec.assignValue(first, target.Field(0)); err != nil {
			return fmt.Errorf("failed to assign TUPLE2 first field: %w", err)
		}
		if err := codec.assignValue(second, target.Field(1)); err != nil {
			return fmt.Errorf("failed to assign TUPLE2 second field: %w", err)
		}
		return nil
	}
	return fmt.Errorf("expected object for TUPLE2, got %T", jsonValue)
}

func (codec *JsonCodec) assignTuple3Value(jsonValue interface{}, target reflect.Value) error {
	if m, ok := jsonValue.(map[string]interface{}); ok {
		first, hasFirst := m["_1"]
		second, hasSecond := m["_2"]
		third, hasThird := m["_3"]
		if !hasFirst || !hasSecond || !hasThird {
			return fmt.Errorf("TUPLE3 missing _1, _2, or _3 fields")
		}
		if err := codec.assignValue(first, target.Field(0)); err != nil {
			return fmt.Errorf("failed to assign TUPLE3 first field: %w", err)
		}
		if err := codec.assignValue(second, target.Field(1)); err != nil {
			return fmt.Errorf("failed to assign TUPLE3 second field: %w", err)
		}
		if err := codec.assignValue(third, target.Field(2)); err != nil {
			return fmt.Errorf("failed to assign TUPLE3 third field: %w", err)
		}
		return nil
	}
	return fmt.Errorf("expected object for TUPLE3, got %T", jsonValue)
}

func (codec *JsonCodec) assignSliceValue(jsonValue interface{}, target reflect.Value) error {
	if arr, ok := jsonValue.([]interface{}); ok {
		slice := reflect.MakeSlice(target.Type(), len(arr), len(arr))

		for i, v := range arr {
			elem := slice.Index(i)
			if err := codec.assignValue(v, elem); err != nil {
				return fmt.Errorf("failed to assign slice element %d: %w", i, err)
			}
		}

		target.Set(slice)
		return nil
	}
	return fmt.Errorf("expected array for slice, got %T", jsonValue)
}

func (codec *JsonCodec) assignStructValue(jsonValue interface{}, target reflect.Value) error {
	if m, ok := jsonValue.(map[string]interface{}); ok {
		targetType := target.Type()

		for i := 0; i < target.NumField(); i++ {
			field := target.Field(i)
			fieldType := targetType.Field(i)

			if !field.CanSet() {
				continue
			}

			// Get JSON tag or use field name
			jsonTag := fieldType.Tag.Get("json")
			fieldName := fieldType.Name
			if jsonTag != "" && jsonTag != "-" {
				if commaIdx := strings.IndexByte(jsonTag, ','); commaIdx != -1 {
					fieldName = jsonTag[:commaIdx]
				} else {
					fieldName = jsonTag
				}
			}

			if jsonVal, exists := m[fieldName]; exists {
				if err := codec.assignValue(jsonVal, field); err != nil {
					return fmt.Errorf("failed to assign field %s: %w", fieldName, err)
				}
			}
		}

		return nil
	}
	return fmt.Errorf("expected object for struct, got %T", jsonValue)
}

func (codec *JsonCodec) assignVariantValue(jsonValue interface{}, target reflect.Value) error {
	if m, ok := jsonValue.(map[string]interface{}); ok {
		tag, hasTag := m["tag"].(string)
		if !hasTag {
			return fmt.Errorf("variant missing tag field")
		}

		// For other variant types, try to use built-in UnmarshalJSON if available
		if target.Type().Implements(reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()) {
			jsonBytes, err := json.Marshal(jsonValue)
			if err != nil {
				return fmt.Errorf("failed to re-marshal variant for custom unmarshaling: %w", err)
			}

			newValue := reflect.New(target.Type())
			if unmarshaler, ok := newValue.Interface().(json.Unmarshaler); ok {
				if err := unmarshaler.UnmarshalJSON(jsonBytes); err != nil {
					return fmt.Errorf("custom variant unmarshaling failed: %w", err)
				}
				target.Set(newValue.Elem())
				return nil
			}
		}

		return fmt.Errorf("variant unmarshalling not implemented for type %v with tag %s", target.Type(), tag)
	}
	return fmt.Errorf("expected object for variant, got %T", jsonValue)
}

func (codec *JsonCodec) assignEnumValue(jsonValue interface{}, target reflect.Value) error {
	if str, ok := jsonValue.(string); ok {
		if target.Type().Kind() == reflect.String {
			target.Set(reflect.ValueOf(str).Convert(target.Type()))
			return nil
		}

		return fmt.Errorf("enum unmarshalling not implemented for type %v with value %s", target.Type(), str)
	}
	return fmt.Errorf("expected string for enum, got %T", jsonValue)
}

func (codec *JsonCodec) assignIntValue(jsonValue interface{}, target reflect.Value) error {
	switch v := jsonValue.(type) {
	case string:
		if i, err := parseIntFromString(v); err == nil {
			target.SetInt(i)
			return nil
		}
		return fmt.Errorf("invalid string format for int: %s", v)
	case float64:
		target.SetInt(int64(v))
		return nil
	case int64:
		target.SetInt(v)
		return nil
	default:
		return fmt.Errorf("expected string or number for int, got %T", jsonValue)
	}
}

func (codec *JsonCodec) assignUintValue(jsonValue interface{}, target reflect.Value) error {
	switch v := jsonValue.(type) {
	case string:
		if i, err := parseIntFromString(v); err == nil && i >= 0 {
			target.SetUint(uint64(i))
			return nil
		}
		return fmt.Errorf("invalid string format for uint: %s", v)
	case float64:
		if v >= 0 {
			target.SetUint(uint64(v))
			return nil
		}
		return fmt.Errorf("negative value for uint: %f", v)
	case int64:
		if v >= 0 {
			target.SetUint(uint64(v))
			return nil
		}
		return fmt.Errorf("negative value for uint: %d", v)
	default:
		return fmt.Errorf("expected string or number for uint, got %T", jsonValue)
	}
}

func (codec *JsonCodec) assignFloatValue(jsonValue interface{}, target reflect.Value) error {
	switch v := jsonValue.(type) {
	case string:
		if f, err := parseFloatFromString(v); err == nil {
			target.SetFloat(f)
			return nil
		}
		return fmt.Errorf("invalid string format for float: %s", v)
	case float64:
		target.SetFloat(v)
		return nil
	case int64:
		target.SetFloat(float64(v))
		return nil
	default:
		return fmt.Errorf("expected string or number for float, got %T", jsonValue)
	}
}

func parseIntFromString(s string) (int64, error) {
	var i int64
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}

func parseFloatFromString(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}
