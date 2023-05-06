package wrap2

import (
	"bytes"
	"encoding/json"

	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/wrap"
)

func ToGo(v JsonValue) any {
	switch v.Type() {
	case wrap.JsonTypeNull:
		return nil
	case wrap.JsonTypeBoolean:
		return v.BooleanGet()
	case wrap.JsonTypeNumber:
		return v.NumberGet()
	case wrap.JsonTypeString:
		return v.StringGet()
	case wrap.JsonTypeArray:
		l := v.ArrayLen()
		a := make([]any, l)
		for i := 0; i < l; i++ {
			val, _ := v.ArrayGetElm(i)
			a[i] = ToGo(val)
		}
		return a
	case wrap.JsonTypeObject:
		m := map[string]any{}
		keys := v.ObjectKeys()
		for _, key := range keys {
			val, _ := v.ObjectGetElm(key)
			m[key] = ToGo(val)
		}
		return m
	default:
		return errors.Unreachable[any]("Unreachable")
	}
}

func FromGo(valAny any) (JsonValue, error) {
	if v, ok := valAny.(JsonValue); ok {
		return v.Clone(), nil
	}

	errInfo := errors.Info{"valAny": valAny}

	switch val := valAny.(type) {
	case nil:
		return Null(), nil
	case string:
		return String(val), nil
	case json.Number:
		return Number(val), nil
	case bool:
		return Boolean(val), nil
	case map[string]any:
		if val == nil {
			return Null(), nil
		}
		o := Object()
		for k, v := range val {
			v, err := FromGo(v)
			if err != nil {
				return nil, errors.Wrap(errors.BadConversion.Err(err), errInfo.AppendTo("fail to marshal to JsonValue"))
			}
			o.ObjectSetElm(k, v)
		}
		return o, nil
	case []any:
		if val == nil {
			return Null(), nil
		}
		a := Array()
		for _, v := range val {
			v, err := FromGo(v)
			if err != nil {
				return nil, errors.Wrap(errors.BadConversion.Err(err), errInfo.AppendTo("fail to marshal to JsonValue"))
			}
			a.ArrayAddElm(v)
		}
		return a, nil
	}

	buf := bytes.NewBuffer(nil)

	encoder := json.NewEncoder(buf)
	if err := encoder.Encode(valAny); err != nil {
		return nil, errors.Wrap(errors.BadConversion.Err(err), errInfo.AppendTo("fail to marshal to JsonValue"))
	}

	decoder := json.NewDecoder(buf)
	decoder.UseNumber()
	var val any
	if err := decoder.Decode(&val); err != nil {
		return nil, errors.Wrap(errors.BadConversion.Err(err), errInfo.AppendTo("fail to marshal to JsonValue"))
	}

	return FromGo(val)
}
