package openapi

import (
	"iter"
	"testing"

	_json "github.com/MarkRosemaker/openapi/internal/json"
)

type (
	mapType  map[string]*mapValue
	mapValue struct{ idx int }
)

func getIndex(v *mapValue) int      { return v.idx }
func setIndex(v *mapValue, idx int) { v.idx = idx }

func (m *mapType) Set(key string, v *mapValue) {
	setToMap(m, key, v, getIndex, setIndex)
}

func (cs mapType) ByIndex() iter.Seq2[string, *mapValue] {
	return _json.OrderedMapByIndex(cs, getIndex)
}

func TestSetToMap(t *testing.T) {
	var m mapType

	m.Set("foo", &mapValue{})
	if len(m) != 1 {
		t.Fatalf("expected 1, got %d", len(m))
	}

	m.Set("bar", &mapValue{})
	if len(m) != 2 {
		t.Fatalf("expected 2, got %d", len(m))
	}

	m.Set("baz", &mapValue{})
	if len(m) != 3 {
		t.Fatalf("expected 3, got %d", len(m))
	}

	i := 0
	keys := []string{"foo", "bar", "baz"}
	for k := range m.ByIndex() {
		if k != keys[i] {
			t.Fatalf("expected %s, got %s", keys[i], k)
		}
		i++
	}
}
