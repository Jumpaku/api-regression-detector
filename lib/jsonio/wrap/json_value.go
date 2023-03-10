package wrap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Jumpaku/api-regression-detector/lib/errors"
)

type JsonType int

const (
	JsonTypeNull JsonType = iota
	JsonTypeString
	JsonTypeNumber
	JsonTypeBoolean
	JsonTypeArray
	JsonTypeObject
)

type JsonKey string

func JsonKeyInteger(index int) JsonKey {
	return JsonKey(strconv.FormatInt(int64(index), 10))
}

func (k JsonKey) String() string {
	return string(k)
}

func (k JsonKey) Integer() (int, bool) {
	v, err := strconv.ParseInt(string(k), 10, 64)
	return int(v), err == nil
}

type JsonValue struct {
	Type         JsonType
	StringValue  string
	NumberValue  json.Number
	BooleanValue bool
	ArrayValue   JsonArray
	ObjectValue  JsonObject
}

func String(v string) *JsonValue {
	return &JsonValue{Type: JsonTypeString, StringValue: v}
}

func mustToInt64(v any) int64 {
	switch v := v.(type) {
	case int:
		return int64(v)
	case int8:
		return int64(v)
	case int16:
		return int64(v)
	case int32:
		return int64(v)
	case int64:
		return int64(v)
	case uint:
		return int64(v)
	case uint8:
		return int64(v)
	case uint16:
		return int64(v)
	case uint32:
		return int64(v)
	case uint64:
		return int64(v)
	}
	panic(fmt.Sprintf("must be compatible to int64 but %v", v))
}

func mustToFloat64(v any) float64 {
	switch v := v.(type) {
	case float32:
		return float64(v)
	case float64:
		return float64(v)
	}
	panic(fmt.Sprintf("must be compatible to float64 but %v", v))
}

func Number[T json.Number | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64](v T) *JsonValue {
	switch v := any(v).(type) {
	case json.Number:
		return &JsonValue{Type: JsonTypeNumber, NumberValue: v}
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return Number(json.Number(strconv.FormatInt(mustToInt64(v), 10)))
	case float32, float64:
		return Number(json.Number(strconv.FormatFloat(mustToFloat64(v), 'f', 16, 64)))
	default:
		return errors.Unreachable[*JsonValue]("never match this case %#v", v)
	}
}

func Boolean(v bool) *JsonValue {
	return &JsonValue{Type: JsonTypeBoolean, BooleanValue: v}
}

func Null() *JsonValue {
	return &JsonValue{Type: JsonTypeNull}
}

func Object(mv map[string]*JsonValue) *JsonValue {
	for k, v := range mv {
		if v == nil {
			mv[k] = Null()
		}
	}
	return &JsonValue{Type: JsonTypeObject, ObjectValue: JsonObject(mv)}
}

func Array(av ...*JsonValue) *JsonValue {
	for i, v := range av {
		if v == nil {
			av[i] = Null()
		}
	}
	return &JsonValue{Type: JsonTypeArray, ArrayValue: JsonArray(av)}
}

var (
	_ json.Marshaler   = (*JsonValue)(nil)
	_ json.Unmarshaler = (*JsonValue)(nil)
)

func (v *JsonValue) MarshalJSON() ([]byte, error) {
	a := ToAny(v)

	b, err := json.Marshal(a)
	if err != nil {
		return nil, errors.Wrap(errors.BadConversion.Err(err), errors.Info{"valAny": a}.AppendTo("fail to convert JSON to []byte"))
	}

	return b, nil
}

func (v *JsonValue) UnmarshalJSON(b []byte) error {
	var a any
	decoder := json.NewDecoder(bytes.NewBuffer(b))
	decoder.UseNumber()

	if err := decoder.Decode(&a); err != nil {
		return errors.Wrap(errors.BadConversion.Err(err), errors.Info{"data": string(b)}.AppendTo("fail to convert JSON to JsonValue"))
	}

	u, err := FromAny(a)
	if err != nil {
		return errors.Wrap(errors.BadConversion.Err(err), errors.Info{"value": a}.AppendTo("fail to convert JSON to JsonValue"))
	}

	*v = *u

	return nil
}

func (v *JsonValue) MustString() string {
	errors.Assert(v.Type == JsonTypeString, "MustString() must be called with JsonTypeString (val=%#v)", v)

	return v.StringValue
}

func (v *JsonValue) MustBool() bool {
	errors.Assert(v.Type == JsonTypeBoolean, "MustBool() must be called with JsonTypeBoolean (val=%#v)", v)

	return v.BooleanValue
}

func (v *JsonValue) MustNumber() json.Number {
	errors.Assert(v.Type == JsonTypeNumber, "MustNumber() must be called with JsonTypeNumber (val=%#v)", v)

	return v.NumberValue
}

func (v *JsonValue) Int64() (int64, bool) {
	errors.Assert(v.Type == JsonTypeNumber, "Int64() must be called with JsonTypeNumber (val=%#v)", v)

	n, err := v.NumberValue.Int64()

	return n, err == nil
}

func (v *JsonValue) Float64() (float64, bool) {
	errors.Assert(v.Type == JsonTypeNumber, "Float64() must be called with JsonTypeNumber (val=%#v)", v)

	n, err := v.NumberValue.Float64()

	return n, err == nil
}

func (v *JsonValue) MustObject() JsonObject {
	errors.Assert(v.Type == JsonTypeObject, "MustObject() must be called with JsonTypeObject (val=%#v)", v)

	return v.ObjectValue
}

func (v *JsonValue) MustArray() JsonArray {
	errors.Assert(v.Type == JsonTypeArray, "MustArray() must be called with JsonTypeArray (val=%#v)", v)

	return v.ArrayValue
}

func (v *JsonValue) Walk(visitor func(keys []JsonKey, val *JsonValue) error) error {
	return walkImpl([]JsonKey{}, v, visitor)
}

func walkImpl(parentKey []JsonKey, val *JsonValue, walkFunc func(keys []JsonKey, val *JsonValue) error) error {
	if err := walkFunc(parentKey, val); err != nil {
		return err
	}
	switch val.Type {
	case JsonTypeObject:
		for k, vk := range val.MustObject() {
			if err := walkImpl(append(parentKey, JsonKey(k)), vk, walkFunc); err != nil {
				return err
			}
		}
	case JsonTypeArray:
		for i, vi := range val.MustArray() {
			if err := walkImpl(append(parentKey, JsonKeyInteger(i)), vi, walkFunc); err != nil {
				return err
			}
		}
	}
	return nil
}
