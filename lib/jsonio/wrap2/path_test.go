package wrap2_test

import (
	"fmt"
	"testing"

	"github.com/Jumpaku/api-regression-detector/lib/jsonio/wrap"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/wrap2"
	"github.com/Jumpaku/api-regression-detector/test/assert"
)

func TestKey_String(t *testing.T) {
	v := wrap2.Key("abc").String()
	assert.Equal(t, v, "abc")
}

func TestKey_Integer(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		v, ok := wrap2.Key("123").Integer()
		assert.Equal(t, ok, true)
		assert.Equal(t, v, 123)
	})
	t.Run("ng", func(t *testing.T) {
		_, ok := wrap2.Key("abc").Integer()
		assert.Equal(t, ok, false)
	})
}

func TestPath_Equals(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		p := wrap2.Path([]wrap2.Key{"abc", "123"})
		assert.Equal(t, p.Equals(wrap2.Path([]wrap2.Key{"abc", "123"})), true)
	})
	t.Run("not equal", func(t *testing.T) {
		p := wrap2.Path([]wrap2.Key{"abc", "123"})
		assert.Equal(t, p.Equals(wrap2.Path([]wrap2.Key{"abc", "123", "xyz"})), false)
	})
}

func TestPath_Get(t *testing.T) {
	p := wrap2.Path([]wrap2.Key{"abc", "123"})
	k0, ok := p.Get(0)
	assert.Equal(t, ok, true)
	assert.Equal(t, k0, wrap2.Key("abc"))
	k1, ok := p.Get(1)
	assert.Equal(t, ok, true)
	assert.Equal(t, k1, wrap2.Key("123"))

	_, ok = p.Get(2)
	assert.Equal(t, ok, false)
	_, ok = p.Get(-1)
	assert.Equal(t, ok, false)

}

func TestPath_Len(t *testing.T) {
	p := wrap2.Path([]wrap2.Key{"abc", "123"})
	assert.Equal(t, p.Len(), 2)
}

func TestPath_Append(t *testing.T) {
	p := wrap2.Path([]wrap2.Key{"abc", "123"}).Append("xyz")
	assert.Equal(t, p.Len(), 3)
	k1, ok := p.Get(2)
	assert.Equal(t, ok, true)
	assert.Equal(t, k1, wrap2.Key("xyz"))
}

func TestWalk(t *testing.T) {
	t.Run(`error`, func(t *testing.T) {
		v := wrap2.Null()
		err := wrap2.Walk(v, func(path wrap2.Path, val wrap2.JsonValue) error {
			return fmt.Errorf("")
		})
		assert.IsNotNil(t, err)
	})
	t.Run(`null`, func(t *testing.T) {
		v := wrap2.Null()
		p := []wrap2.Path{}
		_ = wrap2.Walk(v, func(path wrap2.Path, val wrap2.JsonValue) error {
			p = append(p, path)
			return nil
		})
		assert.Equal(t, len(p), 1)
	})
	t.Run(`object`, func(t *testing.T) {
		v, _ := wrap2.FromGo(map[string]any{
			"a": nil,
			"b": map[string]any{"x": nil, "y": map[string]any{"w": nil}, "z": []any{nil}},
			"c": []any{nil, map[string]any{"w": nil}, []any{nil}},
		})
		p := []wrap2.Path{}
		_ = wrap2.Walk(v, func(path wrap2.Path, val wrap2.JsonValue) error {
			p = append(p, path)
			return nil
		})
		assert.Equal(t, len(p), 14)
	})
	t.Run(`array`, func(t *testing.T) {
		v, _ := wrap2.FromGo([]any{
			nil,
			map[string]any{"x": nil, "y": map[string]any{"w": nil}, "z": []any{nil}},
			[]any{nil, map[string]any{"w": nil}, []any{nil}},
		})
		p := []wrap2.Path{}
		_ = wrap2.Walk(v, func(path wrap2.Path, val wrap2.JsonValue) error {
			p = append(p, path)
			return nil
		})
		assert.Equal(t, len(p), 14)
	})
}
func TestFind(t *testing.T) {
	t.Run(`not found`, func(t *testing.T) {
		t.Run(`null`, func(t *testing.T) {
			v := wrap2.Null()
			_, ok := wrap2.Find(v, wrap2.Path{"xxx"})
			assert.Equal(t, ok, false)
		})
		t.Run(`object`, func(t *testing.T) {
			v, _ := wrap2.FromGo(map[string]any{
				"a": nil,
				"b": map[string]any{"x": nil, "y": map[string]any{"w": nil}, "z": []any{nil}},
				"c": []any{nil, map[string]any{"w": nil}, []any{nil}},
			})
			_, ok := wrap2.Find(v, wrap2.Path{"xxx"})
			assert.Equal(t, ok, false)
		})
		t.Run(`array`, func(t *testing.T) {
			v, _ := wrap2.FromGo([]any{
				nil,
				map[string]any{"x": nil, "y": map[string]any{"w": nil}, "z": []any{nil}},
				[]any{nil, map[string]any{"w": nil}, []any{nil}},
			})
			_, ok := wrap2.Find(v, wrap2.Path{"xxx"})
			assert.Equal(t, ok, false)
		})
	})
	t.Run(`null`, func(t *testing.T) {
		v := wrap2.Null()
		a, ok := wrap2.Find(v, wrap2.Path{})
		assert.Equal(t, ok, true)
		assert.Equal(t, a.Type(), wrap.JsonTypeNull)
	})
	t.Run(`object`, func(t *testing.T) {
		v, _ := wrap2.FromGo(map[string]any{
			"a": nil,
			"b": map[string]any{"x": nil, "y": map[string]any{"w": nil}, "z": []any{nil}},
			"c": []any{nil, map[string]any{"w": nil}, []any{nil}},
		})
		t.Run(".", func(t *testing.T) {
			a, ok := wrap2.Find(v, wrap2.Path{})
			assert.Equal(t, ok, true)
			assert.Equal(t, a.Type(), wrap.JsonTypeObject)
		})
		t.Run(".a", func(t *testing.T) {
			a, ok := wrap2.Find(v, wrap2.Path{"a"})
			assert.Equal(t, ok, true)
			assert.Equal(t, a.Type(), wrap.JsonTypeNull)
		})
		t.Run(".b", func(t *testing.T) {
			a, ok := wrap2.Find(v, wrap2.Path{"b"})
			assert.Equal(t, ok, true)
			assert.Equal(t, a.Type(), wrap.JsonTypeObject)
		})
		t.Run(".b.x", func(t *testing.T) {
			a, ok := wrap2.Find(v, wrap2.Path{"b", "x"})
			assert.Equal(t, ok, true)
			assert.Equal(t, a.Type(), wrap.JsonTypeNull)
		})
		t.Run(".b.y", func(t *testing.T) {
			a, ok := wrap2.Find(v, wrap2.Path{"b", "y"})
			assert.Equal(t, ok, true)
			assert.Equal(t, a.Type(), wrap.JsonTypeObject)
		})
		t.Run(".b.y.w", func(t *testing.T) {
			a, ok := wrap2.Find(v, wrap2.Path{"b", "y", "w"})
			assert.Equal(t, ok, true)
			assert.Equal(t, a.Type(), wrap.JsonTypeNull)
		})
		t.Run(".b.z", func(t *testing.T) {
			a, ok := wrap2.Find(v, wrap2.Path{"b", "z"})
			assert.Equal(t, ok, true)
			assert.Equal(t, a.Type(), wrap.JsonTypeArray)
		})
		t.Run(".b.z.0", func(t *testing.T) {
			a, ok := wrap2.Find(v, wrap2.Path{"b", "z", "0"})
			assert.Equal(t, ok, true)
			assert.Equal(t, a.Type(), wrap.JsonTypeNull)
		})
		t.Run(".c", func(t *testing.T) {
			a, ok := wrap2.Find(v, wrap2.Path{"c"})
			assert.Equal(t, ok, true)
			assert.Equal(t, a.Type(), wrap.JsonTypeArray)
		})
		t.Run(".c.0", func(t *testing.T) {
			a, ok := wrap2.Find(v, wrap2.Path{"c", "0"})
			assert.Equal(t, ok, true)
			assert.Equal(t, a.Type(), wrap.JsonTypeNull)
		})
		t.Run(".c.1", func(t *testing.T) {
			a, ok := wrap2.Find(v, wrap2.Path{"c", "1"})
			assert.Equal(t, ok, true)
			assert.Equal(t, a.Type(), wrap.JsonTypeObject)
		})
		t.Run(".c.1.w", func(t *testing.T) {
			a, ok := wrap2.Find(v, wrap2.Path{"c", "1", "w"})
			assert.Equal(t, ok, true)
			assert.Equal(t, a.Type(), wrap.JsonTypeNull)
		})
		t.Run(".c.2", func(t *testing.T) {
			a, ok := wrap2.Find(v, wrap2.Path{"c", "2"})
			assert.Equal(t, ok, true)
			assert.Equal(t, a.Type(), wrap.JsonTypeArray)
		})
		t.Run(".c.2.0", func(t *testing.T) {
			a, ok := wrap2.Find(v, wrap2.Path{"c", "2", "0"})
			assert.Equal(t, ok, true)
			assert.Equal(t, a.Type(), wrap.JsonTypeNull)
		})
	})
	t.Run(`array`, func(t *testing.T) {
		v, _ := wrap2.FromGo([]any{
			nil,
			map[string]any{"x": nil, "y": map[string]any{"w": nil}, "z": []any{nil}},
			[]any{nil, map[string]any{"w": nil}, []any{nil}},
		})
		t.Run(".", func(t *testing.T) {
			a, ok := wrap2.Find(v, wrap2.Path{})
			assert.Equal(t, ok, true)
			assert.Equal(t, a.Type(), wrap.JsonTypeArray)
		})
		t.Run(".0", func(t *testing.T) {
			a, ok := wrap2.Find(v, wrap2.Path{"0"})
			assert.Equal(t, ok, true)
			assert.Equal(t, a.Type(), wrap.JsonTypeNull)
		})
		t.Run(".1", func(t *testing.T) {
			a, ok := wrap2.Find(v, wrap2.Path{"1"})
			assert.Equal(t, ok, true)
			assert.Equal(t, a.Type(), wrap.JsonTypeObject)
		})
		t.Run(".1.x", func(t *testing.T) {
			a, ok := wrap2.Find(v, wrap2.Path{"1", "x"})
			assert.Equal(t, ok, true)
			assert.Equal(t, a.Type(), wrap.JsonTypeNull)
		})
		t.Run(".1.y", func(t *testing.T) {
			a, ok := wrap2.Find(v, wrap2.Path{"1", "y"})
			assert.Equal(t, ok, true)
			assert.Equal(t, a.Type(), wrap.JsonTypeObject)
		})
		t.Run(".1.y.w", func(t *testing.T) {
			a, ok := wrap2.Find(v, wrap2.Path{"1", "y", "w"})
			assert.Equal(t, ok, true)
			assert.Equal(t, a.Type(), wrap.JsonTypeNull)
		})
		t.Run(".1.z", func(t *testing.T) {
			a, ok := wrap2.Find(v, wrap2.Path{"1", "z"})
			assert.Equal(t, ok, true)
			assert.Equal(t, a.Type(), wrap.JsonTypeArray)
		})
		t.Run(".1.z.0", func(t *testing.T) {
			a, ok := wrap2.Find(v, wrap2.Path{"1", "z", "0"})
			assert.Equal(t, ok, true)
			assert.Equal(t, a.Type(), wrap.JsonTypeNull)
		})
		t.Run(".2", func(t *testing.T) {
			a, ok := wrap2.Find(v, wrap2.Path{"2"})
			assert.Equal(t, ok, true)
			assert.Equal(t, a.Type(), wrap.JsonTypeArray)
		})
		t.Run(".2.0", func(t *testing.T) {
			a, ok := wrap2.Find(v, wrap2.Path{"2", "0"})
			assert.Equal(t, ok, true)
			assert.Equal(t, a.Type(), wrap.JsonTypeNull)
		})
		t.Run(".2.1", func(t *testing.T) {
			a, ok := wrap2.Find(v, wrap2.Path{"2", "1"})
			assert.Equal(t, ok, true)
			assert.Equal(t, a.Type(), wrap.JsonTypeObject)
		})
		t.Run(".2.1.w", func(t *testing.T) {
			a, ok := wrap2.Find(v, wrap2.Path{"2", "1", "w"})
			assert.Equal(t, ok, true)
			assert.Equal(t, a.Type(), wrap.JsonTypeNull)
		})
		t.Run(".2.2", func(t *testing.T) {
			a, ok := wrap2.Find(v, wrap2.Path{"2", "2"})
			assert.Equal(t, ok, true)
			assert.Equal(t, a.Type(), wrap.JsonTypeArray)
		})
		t.Run(".2.2.0", func(t *testing.T) {
			a, ok := wrap2.Find(v, wrap2.Path{"2", "2", "0"})
			assert.Equal(t, ok, true)
			assert.Equal(t, a.Type(), wrap.JsonTypeNull)
		})
	})
}
