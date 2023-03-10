package wrap_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/Jumpaku/api-regression-detector/lib/jsonio/wrap"
	"github.com/Jumpaku/api-regression-detector/test/assert"
)

func TestJsonKey_String(t *testing.T) {
	t.Run("number", func(t *testing.T) {
		_, ok := wrap.JsonKey("abc").Integer()
		assert.Equal(t, ok, false)
	})
	t.Run("string", func(t *testing.T) {
		a := wrap.JsonKey("abc").String()
		assert.Equal(t, a, "abc")
	})
}
func TestJsonKey_Integer(t *testing.T) {
	t.Run("number", func(t *testing.T) {
		a, ok := wrap.JsonKeyInteger(123).Integer()
		assert.Equal(t, ok, true)
		assert.Equal(t, a, 123)
	})
	t.Run("string", func(t *testing.T) {
		a := wrap.JsonKeyInteger(123).String()
		assert.Equal(t, a, "123")
	})
}

func TestString(t *testing.T) {
	s := "abc"
	v := wrap.String(s)
	assert.Equal(t, v.Type, wrap.JsonTypeString)
	a := v.MustString()
	assert.Equal(t, a, s)
}

func TestBoolean_True(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		b := true
		v := wrap.Boolean(b)
		assert.Equal(t, v.Type, wrap.JsonTypeBoolean)
		assert.Equal(t, v.MustBool(), b)
	})
	t.Run("false", func(t *testing.T) {
		b := false
		v := wrap.Boolean(b)
		assert.Equal(t, v.Type, wrap.JsonTypeBoolean)
		assert.Equal(t, v.MustBool(), b)
	})
}

func TestNumber_JsonNumber(t *testing.T) {
	t.Run("json.Number", func(t *testing.T) {
		t.Run("int64", func(t *testing.T) {
			n := json.Number("123")
			v := wrap.Number(n)
			assert.Equal(t, v.Type, wrap.JsonTypeNumber)
			a, ok := v.Int64()
			assert.Equal(t, ok, true)
			assert.Equal(t, a, int64(123))
			a, err := v.MustNumber().Int64()
			assert.Equal(t, err, nil)
			assert.Equal(t, a, int64(123))
		})
		t.Run("float64", func(t *testing.T) {
			n := json.Number("-123.45")
			v := wrap.Number(n)
			assert.Equal(t, v.Type, wrap.JsonTypeNumber)
			a, ok := v.Float64()
			assert.Equal(t, ok, true)
			assert.Equal(t, a, -123.45)
			a, err := v.MustNumber().Float64()
			assert.Equal(t, err, nil)
			assert.Equal(t, a, -123.45)
		})
	})
	t.Run("int64", func(t *testing.T) {
		i := int64(123)
		v := wrap.Number(i)
		assert.Equal(t, v.Type, wrap.JsonTypeNumber)
		a, ok := v.Int64()
		assert.Equal(t, ok, true)
		assert.Equal(t, a, i)
		a, err := v.MustNumber().Int64()
		assert.Equal(t, err, nil)
		assert.Equal(t, a, i)
	})
	t.Run("float64", func(t *testing.T) {
		f := -123.45
		v := wrap.Number(f)
		assert.Equal(t, v.Type, wrap.JsonTypeNumber)
		a, ok := v.Float64()
		assert.Equal(t, ok, true)
		assert.Equal(t, a, f)
		a, err := v.MustNumber().Float64()
		assert.Equal(t, err, nil)
		assert.Equal(t, a, f)
	})
}

func TestNull(t *testing.T) {
	v := wrap.Null()
	assert.NotEqual(t, v, nil)
	assert.Equal(t, v.Type, wrap.JsonTypeNull)
}

func TestObject(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		v := wrap.Object(nil)
		assert.Equal(t, v.Type, wrap.JsonTypeObject)
		assert.Equal(t, len(v.MustObject()), 0)
	})
	t.Run("elements", func(t *testing.T) {
		v := wrap.Object(map[string]*wrap.JsonValue{
			"a": nil,
			"g": wrap.Null(),
			"b": wrap.String("xyz"),
			"c": wrap.Number(123),
			"d": wrap.Number(-123.45),
			"e": wrap.Boolean(true),
			"f": wrap.Boolean(false),
		})
		assert.Equal(t, v.Type, wrap.JsonTypeObject)
		assert.Equal(t, len(v.MustObject()), 7)
	})
}

func TestArray(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		v := wrap.Array()
		assert.Equal(t, v.Type, wrap.JsonTypeArray)
		assert.Equal(t, len(v.MustArray()), 0)
	})
	t.Run("elements", func(t *testing.T) {
		v := wrap.Array(
			nil,
			wrap.Null(),
			wrap.String("xyz"),
			wrap.Number(123),
			wrap.Number(-123.45),
			wrap.Boolean(true),
			wrap.Boolean(false),
		)
		assert.Equal(t, v.Type, wrap.JsonTypeArray)
		assert.Equal(t, len(v.MustArray()), 7)
	})
}

func TestMarshalJSON(t *testing.T) {
	t.Run("null", func(t *testing.T) {
		a, err := wrap.Null().MarshalJSON()
		assert.Equal(t, err, nil)
		assert.Equal(t, string(a), "null")
	})
	t.Run("integer", func(t *testing.T) {
		a, err := wrap.Number(123).MarshalJSON()
		assert.Equal(t, err, nil)
		assert.Equal(t, string(a), "123")
	})
	t.Run("float", func(t *testing.T) {
		a, err := wrap.Number(-123.45).MarshalJSON()
		assert.Equal(t, err, nil)
		assert.Equal(t, strings.HasPrefix(string(a), "-123.45"), true)
	})
	t.Run("string", func(t *testing.T) {
		a, err := wrap.String("abc").MarshalJSON()
		assert.Equal(t, err, nil)
		assert.Equal(t, string(a), `"abc"`)
	})
	t.Run("bool", func(t *testing.T) {
		t.Run("true", func(t *testing.T) {
			a, err := wrap.Boolean(true).MarshalJSON()
			assert.Equal(t, err, nil)
			assert.Equal(t, string(a), "true")
		})
		t.Run("false", func(t *testing.T) {
			a, err := wrap.Boolean(false).MarshalJSON()
			assert.Equal(t, err, nil)
			assert.Equal(t, string(a), "false")
		})
	})
	t.Run("array", func(t *testing.T) {
		a, err := wrap.Array(wrap.Boolean(true), wrap.Boolean(false), wrap.Null(), wrap.Number(123), wrap.String("abc"), wrap.Array(), wrap.Object(nil)).MarshalJSON()
		assert.Equal(t, err, nil)
		assert.Equal(t, string(a), `[true,false,null,123,"abc",[],{}]`)
	})
	t.Run("object", func(t *testing.T) {
		a, err := wrap.Object(map[string]*wrap.JsonValue{
			"a": wrap.Boolean(true),
			"b": wrap.Boolean(false),
			"c": wrap.Null(),
			"d": wrap.Number(123),
			"e": wrap.String("abc"),
			"f": wrap.Array(),
			"g": wrap.Object(nil),
		}).MarshalJSON()
		assert.Equal(t, err, nil)
		assert.Equal(t, string(a), `{"a":true,"b":false,"c":null,"d":123,"e":"abc","f":[],"g":{}}`)
	})
}

func TestUnmarshalJSON(t *testing.T) {
	t.Run("null", func(t *testing.T) {
		a := wrap.JsonValue{}
		err := a.UnmarshalJSON([]byte("null"))
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Type, wrap.JsonTypeNull)
	})
	t.Run("integer", func(t *testing.T) {
		a := wrap.JsonValue{}
		err := a.UnmarshalJSON([]byte("123"))
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Type, wrap.JsonTypeNumber)
	})
	t.Run("float", func(t *testing.T) {
		a := wrap.JsonValue{}
		err := a.UnmarshalJSON([]byte("-123.45"))
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Type, wrap.JsonTypeNumber)
	})
	t.Run("string", func(t *testing.T) {
		a := wrap.JsonValue{}
		err := a.UnmarshalJSON([]byte(`"abc"`))
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Type, wrap.JsonTypeString)
	})
	t.Run("bool", func(t *testing.T) {
		t.Run("true", func(t *testing.T) {
			a := wrap.JsonValue{}
			err := a.UnmarshalJSON([]byte("true"))
			assert.Equal(t, err, nil)
			assert.Equal(t, a.Type, wrap.JsonTypeBoolean)
		})
		t.Run("false", func(t *testing.T) {
			a := wrap.JsonValue{}
			err := a.UnmarshalJSON([]byte("false"))
			assert.Equal(t, err, nil)
			assert.Equal(t, a.Type, wrap.JsonTypeBoolean)
		})
	})
	t.Run("array", func(t *testing.T) {
		a := wrap.JsonValue{}
		err := a.UnmarshalJSON([]byte(`[true,false,null,123,"abc",[],{}]`))
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Type, wrap.JsonTypeArray)
	})
	t.Run("object", func(t *testing.T) {
		a := wrap.JsonValue{}
		err := a.UnmarshalJSON([]byte(`{"a":true,"b":false,"c":null,"d":123,"e":"abc","f":[],"g":{}}`))
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Type, wrap.JsonTypeObject)
	})
}

func TestWalk(t *testing.T) {
	t.Run("null", func(t *testing.T) {
		a := [][]wrap.JsonKey{}
		err := wrap.Null().Walk(func(keys []wrap.JsonKey, val *wrap.JsonValue) error {
			a = append(a, keys)
			return nil
		})
		assert.Equal(t, err, nil)
		assert.Equal(t, len(a), 1)
	})
	t.Run("integer", func(t *testing.T) {
		a := [][]wrap.JsonKey{}
		err := wrap.Number(123).Walk(func(keys []wrap.JsonKey, val *wrap.JsonValue) error {
			a = append(a, keys)
			return nil
		})
		assert.Equal(t, err, nil)
		assert.Equal(t, len(a), 1)
	})
	t.Run("float", func(t *testing.T) {
		a := [][]wrap.JsonKey{}
		err := wrap.Number(-123.45).Walk(func(keys []wrap.JsonKey, val *wrap.JsonValue) error {
			a = append(a, keys)
			return nil
		})
		assert.Equal(t, err, nil)
		assert.Equal(t, len(a), 1)
	})
	t.Run("string", func(t *testing.T) {
		a := [][]wrap.JsonKey{}
		err := wrap.String("abc").Walk(func(keys []wrap.JsonKey, val *wrap.JsonValue) error {
			a = append(a, keys)
			return nil
		})
		assert.Equal(t, err, nil)
		assert.Equal(t, len(a), 1)
	})
	t.Run("bool", func(t *testing.T) {
		t.Run("true", func(t *testing.T) {
			a := [][]wrap.JsonKey{}
			err := wrap.Boolean(true).Walk(func(keys []wrap.JsonKey, val *wrap.JsonValue) error {
				a = append(a, keys)
				return nil
			})
			assert.Equal(t, err, nil)
			assert.Equal(t, len(a), 1)
		})
		t.Run("false", func(t *testing.T) {
			a := [][]wrap.JsonKey{}
			err := wrap.Boolean(false).Walk(func(keys []wrap.JsonKey, val *wrap.JsonValue) error {
				a = append(a, keys)
				return nil
			})
			assert.Equal(t, err, nil)
			assert.Equal(t, len(a), 1)
		})
	})
	t.Run("array", func(t *testing.T) {
		a := [][]wrap.JsonKey{}
		err := wrap.Array(
			wrap.Boolean(true),
			wrap.Boolean(false),
			wrap.Null(),
			wrap.Number(123),
			wrap.String("abc"),
			wrap.Array(),
			wrap.Object(nil),
		).Walk(func(keys []wrap.JsonKey, val *wrap.JsonValue) error {
			a = append(a, keys)
			return nil
		})
		assert.Equal(t, err, nil)
		assert.Equal(t, len(a), 8)
	})
	t.Run("object", func(t *testing.T) {
		a := [][]wrap.JsonKey{}
		err := wrap.Object(map[string]*wrap.JsonValue{
			"a": wrap.Boolean(true),
			"b": wrap.Boolean(false),
			"c": wrap.Null(),
			"d": wrap.Number(123),
			"e": wrap.String("abc"),
			"f": wrap.Array(),
			"g": wrap.Object(nil),
		}).Walk(func(keys []wrap.JsonKey, val *wrap.JsonValue) error {
			a = append(a, keys)
			return nil
		})
		assert.Equal(t, err, nil)
		assert.Equal(t, len(a), 8)
	})
}
