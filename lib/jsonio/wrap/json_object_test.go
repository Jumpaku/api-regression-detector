package wrap_test

import (
	"testing"

	jw "github.com/Jumpaku/api-regression-detector/lib/jsonio/wrap"
	"github.com/Jumpaku/api-regression-detector/test/assert"
	"golang.org/x/exp/slices"
)

func TestObjectKeys(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		v := jw.JsonObject(nil)
		a := v.Keys()
		assert.Equal(t, len(a), 0)
	})
	t.Run("elements", func(t *testing.T) {
		v := jw.JsonObject(map[string]*jw.JsonValue{
			"a": jw.Number(123),
			"b": jw.Number(-123.45),
			"c": jw.String("abc"),
			"d": nil,
			"e": jw.Boolean(true),
			"f": jw.Boolean(false),
			"g": jw.Object(map[string]*jw.JsonValue{}),
			"h": jw.Array(),
		})
		a := v.Keys()
		assert.Equal(t, len(a), 8)
		assert.Equal(t, slices.Contains(a, "a"), true)
		assert.Equal(t, slices.Contains(a, "b"), true)
		assert.Equal(t, slices.Contains(a, "c"), true)
		assert.Equal(t, slices.Contains(a, "d"), true)
		assert.Equal(t, slices.Contains(a, "e"), true)
		assert.Equal(t, slices.Contains(a, "f"), true)
		assert.Equal(t, slices.Contains(a, "g"), true)
		assert.Equal(t, slices.Contains(a, "h"), true)
	})
}

func TestObjectLen(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		v := jw.JsonObject(nil)
		assert.Equal(t, v.Len(), 0)
	})
	t.Run("elements", func(t *testing.T) {
		v := jw.JsonObject(map[string]*jw.JsonValue{
			"a": jw.Number(123),
			"b": jw.Number(-123.45),
			"c": jw.String("abc"),
			"d": nil,
			"e": jw.Boolean(true),
			"f": jw.Boolean(false),
			"g": jw.Object(map[string]*jw.JsonValue{}),
			"h": jw.Array(),
		})
		assert.Equal(t, v.Len(), 8)
	})
}

func TestObjectEmpty(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		v := jw.JsonObject(nil)
		assert.Equal(t, v.Empty(), true)
	})
	t.Run("elements", func(t *testing.T) {
		v := jw.JsonObject(map[string]*jw.JsonValue{
			"a": jw.Number(123),
			"b": jw.Number(-123.45),
			"c": jw.String("abc"),
			"d": nil,
			"e": jw.Boolean(true),
			"f": jw.Boolean(false),
			"g": jw.Object(map[string]*jw.JsonValue{}),
			"h": jw.Array(),
		})
		assert.Equal(t, v.Empty(), false)
	})
}

func TestObjectHas(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		v := jw.JsonObject(nil)
		assert.Equal(t, v.Has("x"), false)
	})
	t.Run("elements", func(t *testing.T) {
		v := jw.JsonObject(map[string]*jw.JsonValue{
			"a": jw.Number(123),
			"b": jw.Number(-123.45),
			"c": jw.String("abc"),
			"d": nil,
			"e": jw.Boolean(true),
			"f": jw.Boolean(false),
			"g": jw.Object(map[string]*jw.JsonValue{}),
			"h": jw.Array(),
		})
		assert.Equal(t, v.Has("a"), true)
		assert.Equal(t, v.Has("b"), true)
		assert.Equal(t, v.Has("c"), true)
		assert.Equal(t, v.Has("d"), true)
		assert.Equal(t, v.Has("e"), true)
		assert.Equal(t, v.Has("f"), true)
		assert.Equal(t, v.Has("g"), true)
		assert.Equal(t, v.Has("h"), true)

		assert.Equal(t, v.Has("x"), false)
	})
}

func TestObjectGet(t *testing.T) {
	v := jw.JsonObject(map[string]*jw.JsonValue{
		"x": jw.Object(map[string]*jw.JsonValue{
			"a": jw.Number(123),
			"b": jw.Number(-123.45),
			"c": jw.String("abc"),
			"d": nil,
			"e": jw.Boolean(true),
			"f": jw.Boolean(false),
			"g": jw.Object(nil),
			"h": jw.Array(),
		}),
		"y": jw.Array(
			jw.Number(123),
			jw.Number(-123.45),
			jw.String("abc"),
			nil,
			jw.Boolean(true),
			jw.Boolean(false),
			jw.Object(nil),
			jw.Array(),
		),
	})

	assert.Equal(t, v.Get("x").Type, jw.JsonTypeObject)
	assert.Equal(t, v.Get("x").MustObject().Get("a").Type, jw.JsonTypeNumber)
	assert.Equal(t, v.Get("x").MustObject().Get("b").Type, jw.JsonTypeNumber)
	assert.Equal(t, v.Get("x").MustObject().Get("c").Type, jw.JsonTypeString)
	assert.Equal(t, v.Get("x").MustObject().Get("d").Type, jw.JsonTypeNull)
	assert.Equal(t, v.Get("x").MustObject().Get("e").Type, jw.JsonTypeBoolean)
	assert.Equal(t, v.Get("x").MustObject().Get("f").Type, jw.JsonTypeBoolean)
	assert.Equal(t, v.Get("x").MustObject().Get("g").Type, jw.JsonTypeObject)
	assert.Equal(t, v.Get("x").MustObject().Get("h").Type, jw.JsonTypeArray)

	assert.Equal(t, v.Get("y").Type, jw.JsonTypeArray)
	assert.Equal(t, v.Get("y").MustArray().Get(0).Type, jw.JsonTypeNumber)
	assert.Equal(t, v.Get("y").MustArray().Get(1).Type, jw.JsonTypeNumber)
	assert.Equal(t, v.Get("y").MustArray().Get(2).Type, jw.JsonTypeString)
	assert.Equal(t, v.Get("y").MustArray().Get(3).Type, jw.JsonTypeNull)
	assert.Equal(t, v.Get("y").MustArray().Get(4).Type, jw.JsonTypeBoolean)
	assert.Equal(t, v.Get("y").MustArray().Get(5).Type, jw.JsonTypeBoolean)
	assert.Equal(t, v.Get("y").MustArray().Get(6).Type, jw.JsonTypeObject)
	assert.Equal(t, v.Get("y").MustArray().Get(7).Type, jw.JsonTypeArray)
}

func TestObjectSet(t *testing.T) {
	obj := jw.JsonObject(map[string]*jw.JsonValue{})
	objX := jw.JsonObject(map[string]*jw.JsonValue{})
	objX.Set("a", jw.Number(123))
	objX.Set("b", jw.Number(-123.45))
	objX.Set("c", jw.String("abc"))
	objX.Set("d", nil)
	objX.Set("e", jw.Boolean(true))
	objX.Set("f", jw.Boolean(false))
	objX.Set("g", jw.Object(nil))
	objX.Set("h", jw.Array())
	obj.Set("x", jw.Object(objX))
	obj.Set("y", jw.Array(
		jw.Number(123),
		jw.Number(-123.45),
		jw.String("abc"),
		nil,
		jw.Boolean(true),
		jw.Boolean(false),
		jw.Object(nil),
		jw.Array(),
	))

	assert.Equal(t, obj.Get("x").Type, jw.JsonTypeObject)
	assert.Equal(t, obj.Get("x").MustObject().Get("a").Type, jw.JsonTypeNumber)
	assert.Equal(t, obj.Get("x").MustObject().Get("b").Type, jw.JsonTypeNumber)
	assert.Equal(t, obj.Get("x").MustObject().Get("c").Type, jw.JsonTypeString)
	assert.Equal(t, obj.Get("x").MustObject().Get("d").Type, jw.JsonTypeNull)
	assert.Equal(t, obj.Get("x").MustObject().Get("e").Type, jw.JsonTypeBoolean)
	assert.Equal(t, obj.Get("x").MustObject().Get("f").Type, jw.JsonTypeBoolean)
	assert.Equal(t, obj.Get("x").MustObject().Get("g").Type, jw.JsonTypeObject)
	assert.Equal(t, obj.Get("x").MustObject().Get("h").Type, jw.JsonTypeArray)

	assert.Equal(t, obj.Get("y").Type, jw.JsonTypeArray)
	assert.Equal(t, obj.Get("y").MustArray().Get(0).Type, jw.JsonTypeNumber)
	assert.Equal(t, obj.Get("y").MustArray().Get(1).Type, jw.JsonTypeNumber)
	assert.Equal(t, obj.Get("y").MustArray().Get(2).Type, jw.JsonTypeString)
	assert.Equal(t, obj.Get("y").MustArray().Get(3).Type, jw.JsonTypeNull)
	assert.Equal(t, obj.Get("y").MustArray().Get(4).Type, jw.JsonTypeBoolean)
	assert.Equal(t, obj.Get("y").MustArray().Get(5).Type, jw.JsonTypeBoolean)
	assert.Equal(t, obj.Get("y").MustArray().Get(6).Type, jw.JsonTypeObject)
	assert.Equal(t, obj.Get("y").MustArray().Get(7).Type, jw.JsonTypeArray)
}

func TestObjectDel(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		v := jw.JsonObject(nil)
		v.Del("x")
		assert.Equal(t, v.Len(), 0)
	})
	t.Run("key found", func(t *testing.T) {
		v := jw.JsonObject(map[string]*jw.JsonValue{
			"a": jw.Number(123),
			"b": jw.Number(-123.45),
			"c": jw.String("abc"),
			"d": nil,
			"e": jw.Boolean(true),
			"f": jw.Boolean(false),
			"g": jw.Object(map[string]*jw.JsonValue{}),
			"h": jw.Array(),
		})
		v.Del("a")
		assert.Equal(t, v.Has("a"), false)
	})
}
