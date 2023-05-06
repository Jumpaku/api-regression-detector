package wrap2_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Jumpaku/api-regression-detector/lib/jsonio/wrap"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/wrap2"
	"github.com/Jumpaku/api-regression-detector/test/assert"
	"golang.org/x/exp/slices"
)

func TestType(t *testing.T) {
	t.Run(`null`, func(t *testing.T) {
		assert.Equal(t, wrap2.Null().Type(), wrap.JsonTypeNull)
	})
	t.Run(`number`, func(t *testing.T) {
		assert.Equal(t, wrap2.Number(123).Type(), wrap.JsonTypeNumber)
	})
	t.Run(`string`, func(t *testing.T) {
		assert.Equal(t, wrap2.String("abc").Type(), wrap.JsonTypeString)
	})
	t.Run(`boolean`, func(t *testing.T) {
		assert.Equal(t, wrap2.Boolean(true).Type(), wrap.JsonTypeBoolean)
	})
	t.Run(`array`, func(t *testing.T) {
		assert.Equal(t, wrap2.Array().Type(), wrap.JsonTypeArray)
	})
	t.Run(`object`, func(t *testing.T) {
		assert.Equal(t, wrap2.Object().Type(), wrap.JsonTypeObject)
	})
}

func TestNumberGet(t *testing.T) {
	t.Run(`float`, func(t *testing.T) {
		t.Run(`json.Number`, func(t *testing.T) {
			a, err := wrap2.Number(json.Number("-123.45")).NumberGet().Float64()
			assert.Equal(t, err, nil)
			assert.CloseTo(t, a, float64(-123.45), 1e-9)
		})

		t.Run(`float32`, func(t *testing.T) {
			a, err := wrap2.Number(float32(-123.45)).NumberGet().Float64()
			assert.Equal(t, err, nil)
			assert.CloseTo(t, a, -123.45, 1e-9)
		})

		t.Run(`float64`, func(t *testing.T) {
			a, err := wrap2.Number(float64(-123.45)).NumberGet().Float64()
			assert.Equal(t, err, nil)
			assert.CloseTo(t, a, -123.45, 1e-9)
		})
	})
	t.Run(`integer`, func(t *testing.T) {
		t.Run(`json.Number`, func(t *testing.T) {
			a, err := wrap2.Number(json.Number("123")).NumberGet().Int64()
			assert.Equal(t, err, nil)
			assert.Equal(t, a, 123)
		})

		t.Run(`int`, func(t *testing.T) {
			a, err := wrap2.Number(int(123)).NumberGet().Int64()
			assert.Equal(t, err, nil)
			assert.Equal(t, a, 123)
		})

		t.Run(`int8`, func(t *testing.T) {
			a, err := wrap2.Number(int8(123)).NumberGet().Int64()
			assert.Equal(t, err, nil)
			assert.Equal(t, a, 123)
		})

		t.Run(`int16`, func(t *testing.T) {
			a, err := wrap2.Number(int16(123)).NumberGet().Int64()
			assert.Equal(t, err, nil)
			assert.Equal(t, a, 123)
		})

		t.Run(`int32`, func(t *testing.T) {
			a, err := wrap2.Number(int32(123)).NumberGet().Int64()
			assert.Equal(t, err, nil)
			assert.Equal(t, a, 123)
		})

		t.Run(`int64`, func(t *testing.T) {
			a, err := wrap2.Number(int64(123)).NumberGet().Int64()
			assert.Equal(t, err, nil)
			assert.Equal(t, a, 123)
		})
	})
}

func TestStringGet(t *testing.T) {
	assert.Equal(t, wrap2.String("abc").StringGet(), "abc")
}

func TestBooleanGet(t *testing.T) {
	t.Run(`true`, func(t *testing.T) {
		assert.Equal(t, wrap2.Boolean(true).BooleanGet(), true)
	})

	t.Run(`false`, func(t *testing.T) {
		assert.Equal(t, wrap2.Boolean(false).BooleanGet(), false)
	})
}

func sliceContains[T comparable](val T) func(a []T) (matches bool, message string) {
	return func(a []T) (matches bool, message string) {
		return slices.Contains(a, val), fmt.Sprintf("slice contains %v:%T", val, val)
	}
}
func TestObjectKeys(t *testing.T) {
	o := wrap2.Object(map[string]wrap2.JsonValue{
		"a": wrap2.Null(),
		"b": wrap2.Number(123),
		"c": wrap2.String("abc"),
		"d": wrap2.Boolean(true),
		"e": wrap2.Object(),
		"f": wrap2.Array(),
	})
	aKeys := o.ObjectKeys()

	assert.Equal(t, len(aKeys), 6)
	assert.Match(t, aKeys, sliceContains("a"))
	assert.Match(t, aKeys, sliceContains("b"))
	assert.Match(t, aKeys, sliceContains("c"))
	assert.Match(t, aKeys, sliceContains("d"))
	assert.Match(t, aKeys, sliceContains("e"))
	assert.Match(t, aKeys, sliceContains("f"))
}

func TestObjectGetElm(t *testing.T) {
	o := wrap2.Object(map[string]wrap2.JsonValue{
		"a": wrap2.Null(),
		"b": wrap2.Number(123),
		"c": wrap2.String("abc"),
		"d": wrap2.Boolean(true),
		"e": wrap2.Object(),
		"f": wrap2.Array(),
	})

	t.Run("get null", func(t *testing.T) {
		v, ok := o.ObjectGetElm("a")
		assert.Equal(t, ok, true)
		assert.Equal(t, v.Type(), wrap.JsonTypeNull)
	})

	t.Run("get number", func(t *testing.T) {
		v, ok := o.ObjectGetElm("b")
		assert.Equal(t, ok, true)
		assert.Equal(t, v.Type(), wrap.JsonTypeNumber)
	})

	t.Run("get string", func(t *testing.T) {
		v, ok := o.ObjectGetElm("c")
		assert.Equal(t, ok, true)
		assert.Equal(t, v.Type(), wrap.JsonTypeString)
	})

	t.Run("get boolean", func(t *testing.T) {
		v, ok := o.ObjectGetElm("d")
		assert.Equal(t, ok, true)
		assert.Equal(t, v.Type(), wrap.JsonTypeBoolean)
	})

	t.Run("get object", func(t *testing.T) {
		v, ok := o.ObjectGetElm("e")
		assert.Equal(t, ok, true)
		assert.Equal(t, v.Type(), wrap.JsonTypeObject)
	})

	t.Run("get array", func(t *testing.T) {
		v, ok := o.ObjectGetElm("f")
		assert.Equal(t, ok, true)
		assert.Equal(t, v.Type(), wrap.JsonTypeArray)
	})
}

func TestObjectSetElm(t *testing.T) {
	t.Run("set null", func(t *testing.T) {
		o := wrap2.Object()
		o.ObjectSetElm("a", wrap2.Null())
		v, ok := o.ObjectGetElm("a")
		assert.Equal(t, ok, true)
		assert.Equal(t, v.Type(), wrap.JsonTypeNull)
	})

	t.Run("set number", func(t *testing.T) {
		o := wrap2.Object()
		o.ObjectSetElm("b", wrap2.Number(123))
		v, ok := o.ObjectGetElm("b")
		assert.Equal(t, ok, true)
		assert.Equal(t, v.Type(), wrap.JsonTypeNumber)
	})

	t.Run("set string", func(t *testing.T) {
		o := wrap2.Object()
		o.ObjectSetElm("c", wrap2.String("abc"))
		v, ok := o.ObjectGetElm("c")
		assert.Equal(t, ok, true)
		assert.Equal(t, v.Type(), wrap.JsonTypeString)
	})

	t.Run("set boolean", func(t *testing.T) {
		o := wrap2.Object()
		o.ObjectSetElm("d", wrap2.Boolean(true))
		v, ok := o.ObjectGetElm("d")
		assert.Equal(t, ok, true)
		assert.Equal(t, v.Type(), wrap.JsonTypeBoolean)
	})

	t.Run("set object", func(t *testing.T) {
		o := wrap2.Object()
		o.ObjectSetElm("e", wrap2.Object())
		v, ok := o.ObjectGetElm("e")
		assert.Equal(t, ok, true)
		assert.Equal(t, v.Type(), wrap.JsonTypeObject)
	})

	t.Run("set array", func(t *testing.T) {
		o := wrap2.Object()
		o.ObjectSetElm("f", wrap2.Array())
		v, ok := o.ObjectGetElm("f")
		assert.Equal(t, ok, true)
		assert.Equal(t, v.Type(), wrap.JsonTypeArray)
	})
}

func TestObjectDelElm(t *testing.T) {
	t.Run("delete null", func(t *testing.T) {
		o := wrap2.Object(map[string]wrap2.JsonValue{"a": wrap2.Null()})
		o.ObjectDelElm("a")
		_, ok := o.ObjectGetElm("a")
		assert.Equal(t, ok, false)
	})

	t.Run("delete number", func(t *testing.T) {
		o := wrap2.Object(map[string]wrap2.JsonValue{"b": wrap2.Number(123)})
		o.ObjectDelElm("b")
		_, ok := o.ObjectGetElm("b")
		assert.Equal(t, ok, false)
	})

	t.Run("delete string", func(t *testing.T) {
		o := wrap2.Object(map[string]wrap2.JsonValue{"c": wrap2.String("abc")})
		o.ObjectDelElm("c")
		_, ok := o.ObjectGetElm("c")
		assert.Equal(t, ok, false)
	})

	t.Run("delete boolean", func(t *testing.T) {
		o := wrap2.Object(map[string]wrap2.JsonValue{"d": wrap2.Boolean(true)})
		o.ObjectDelElm("d")
		_, ok := o.ObjectGetElm("d")
		assert.Equal(t, ok, false)
	})

	t.Run("delete object", func(t *testing.T) {
		o := wrap2.Object(map[string]wrap2.JsonValue{"e": wrap2.Object()})
		o.ObjectDelElm("e")
		_, ok := o.ObjectGetElm("e")
		assert.Equal(t, ok, false)
	})

	t.Run("delete array", func(t *testing.T) {
		o := wrap2.Object(map[string]wrap2.JsonValue{"f": wrap2.Array()})
		o.ObjectDelElm("f")
		_, ok := o.ObjectGetElm("f")
		assert.Equal(t, ok, false)
	})
}

func TestObjectLen(t *testing.T) {
	o := wrap2.Object(map[string]wrap2.JsonValue{
		"a": wrap2.Null(),
		"b": wrap2.Number(123),
		"c": wrap2.String("abc"),
		"d": wrap2.Boolean(true),
		"e": wrap2.Object(),
		"f": wrap2.Array(),
	})

	assert.Equal(t, o.ObjectLen(), 6)
}

func TestArrayGetElm(t *testing.T) {
	o := wrap2.Array(
		wrap2.Null(),
		wrap2.Number(123),
		wrap2.String("abc"),
		wrap2.Boolean(true),
		wrap2.Object(),
		wrap2.Array(),
	)

	t.Run("get null", func(t *testing.T) {
		v, ok := o.ArrayGetElm(0)
		assert.Equal(t, ok, true)
		assert.Equal(t, v.Type(), wrap.JsonTypeNull)
	})

	t.Run("get number", func(t *testing.T) {
		v, ok := o.ArrayGetElm(1)
		assert.Equal(t, ok, true)
		assert.Equal(t, v.Type(), wrap.JsonTypeNumber)
	})

	t.Run("get string", func(t *testing.T) {
		v, ok := o.ArrayGetElm(2)
		assert.Equal(t, ok, true)
		assert.Equal(t, v.Type(), wrap.JsonTypeString)
	})

	t.Run("get boolean", func(t *testing.T) {
		v, ok := o.ArrayGetElm(3)
		assert.Equal(t, ok, true)
		assert.Equal(t, v.Type(), wrap.JsonTypeBoolean)
	})

	t.Run("get object", func(t *testing.T) {
		v, ok := o.ArrayGetElm(4)
		assert.Equal(t, ok, true)
		assert.Equal(t, v.Type(), wrap.JsonTypeObject)
	})

	t.Run("get array", func(t *testing.T) {
		v, ok := o.ArrayGetElm(5)
		assert.Equal(t, ok, true)
		assert.Equal(t, v.Type(), wrap.JsonTypeArray)
	})
}

func newExampleArray(n int, v wrap2.JsonValue) wrap2.JsonValue {
	var a = make([]wrap2.JsonValue, n)
	for i := 0; i < n; i++ {
		a[i] = v
	}
	return wrap2.Array(a...)
}
func TestArraySetElm(t *testing.T) {
	t.Run("get null", func(t *testing.T) {
		o := newExampleArray(6, wrap2.Array())
		ok := o.ArraySetElm(0, wrap2.Null())
		assert.Equal(t, ok, true)
		v, ok := o.ArrayGetElm(0)
		assert.Equal(t, ok, true)
		assert.Equal(t, v.Type(), wrap.JsonTypeNull)
	})

	t.Run("get number", func(t *testing.T) {
		o := newExampleArray(6, wrap2.Null())
		ok := o.ArraySetElm(1, wrap2.Number(123))
		assert.Equal(t, ok, true)
		v, ok := o.ArrayGetElm(1)
		assert.Equal(t, ok, true)
		assert.Equal(t, v.Type(), wrap.JsonTypeNumber)
	})

	t.Run("get string", func(t *testing.T) {
		o := newExampleArray(6, wrap2.Null())
		ok := o.ArraySetElm(2, wrap2.String("abc"))
		assert.Equal(t, ok, true)
		v, ok := o.ArrayGetElm(2)
		assert.Equal(t, ok, true)
		assert.Equal(t, v.Type(), wrap.JsonTypeString)
	})

	t.Run("get boolean", func(t *testing.T) {
		o := newExampleArray(6, wrap2.Null())
		ok := o.ArraySetElm(3, wrap2.Boolean(true))
		assert.Equal(t, ok, true)
		v, ok := o.ArrayGetElm(3)
		assert.Equal(t, ok, true)
		assert.Equal(t, v.Type(), wrap.JsonTypeBoolean)
	})

	t.Run("get object", func(t *testing.T) {
		o := newExampleArray(6, wrap2.Null())
		ok := o.ArraySetElm(4, wrap2.Object())
		assert.Equal(t, ok, true)
		v, ok := o.ArrayGetElm(4)
		assert.Equal(t, ok, true)
		assert.Equal(t, v.Type(), wrap.JsonTypeObject)
	})

	t.Run("get array", func(t *testing.T) {
		o := newExampleArray(6, wrap2.Null())
		ok := o.ArraySetElm(5, wrap2.Array())
		assert.Equal(t, ok, true)
		v, ok := o.ArrayGetElm(5)
		assert.Equal(t, ok, true)
		assert.Equal(t, v.Type(), wrap.JsonTypeArray)
	})
}

func TestArrayLen(t *testing.T) {
	o := wrap2.Array(
		wrap2.Null(),
		wrap2.Number(123),
		wrap2.String("abc"),
		wrap2.Boolean(true),
		wrap2.Object(),
		wrap2.Array(),
	)

	assert.Equal(t, o.ArrayLen(), 6)
}

func TestArrayAddElm(t *testing.T) {
	o := wrap2.Array()

	o.ArrayAddElm(wrap2.Null())
	o.ArrayAddElm(wrap2.Number(123))
	o.ArrayAddElm(wrap2.String("abc"))
	o.ArrayAddElm(wrap2.Boolean(true))
	o.ArrayAddElm(wrap2.Object())
	o.ArrayAddElm(wrap2.Array())

	t.Run("get null", func(t *testing.T) {
		v, ok := o.ArrayGetElm(0)
		assert.Equal(t, ok, true)
		assert.Equal(t, v.Type(), wrap.JsonTypeNull)
	})

	t.Run("get number", func(t *testing.T) {
		v, ok := o.ArrayGetElm(1)
		assert.Equal(t, ok, true)
		assert.Equal(t, v.Type(), wrap.JsonTypeNumber)
	})

	t.Run("get string", func(t *testing.T) {
		v, ok := o.ArrayGetElm(2)
		assert.Equal(t, ok, true)
		assert.Equal(t, v.Type(), wrap.JsonTypeString)
	})

	t.Run("get boolean", func(t *testing.T) {
		v, ok := o.ArrayGetElm(3)
		assert.Equal(t, ok, true)
		assert.Equal(t, v.Type(), wrap.JsonTypeBoolean)
	})

	t.Run("get object", func(t *testing.T) {
		v, ok := o.ArrayGetElm(4)
		assert.Equal(t, ok, true)
		assert.Equal(t, v.Type(), wrap.JsonTypeObject)
	})

	t.Run("get array", func(t *testing.T) {
		v, ok := o.ArrayGetElm(5)
		assert.Equal(t, ok, true)
		assert.Equal(t, v.Type(), wrap.JsonTypeArray)
	})
}

func TestArraySlice(t *testing.T) {
	o := wrap2.Array(
		wrap2.Null(),
		wrap2.Number(123),
		wrap2.String("abc"),
		wrap2.Boolean(true),
		wrap2.Object(),
		wrap2.Array(),
	)

	t.Run("whole", func(t *testing.T) {
		a, ok := o.ArraySlice(0, o.ArrayLen())
		assert.Equal(t, ok, true)
		assert.Equal(t, a.ArrayLen(), 6)
	})

	t.Run("empty", func(t *testing.T) {
		a, ok := o.ArraySlice(3, 3)
		assert.Equal(t, ok, true)
		assert.Equal(t, a.ArrayLen(), 0)
	})

	t.Run("sub-array", func(t *testing.T) {
		a, ok := o.ArraySlice(2, 4)
		assert.Equal(t, ok, true)
		assert.Equal(t, a.ArrayLen(), 2)
	})
}
