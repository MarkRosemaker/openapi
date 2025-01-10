package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestSchemaRefs_Set(t *testing.T) {
	var refs openapi.SchemaRefs

	refs.Set("foo", &openapi.SchemaRef{})
	if len(refs) != 1 {
		t.Fatalf("expected 1, got %d", len(refs))
	}

	refs.Set("bar", &openapi.SchemaRef{})
	if len(refs) != 2 {
		t.Fatalf("expected 2, got %d", len(refs))
	}

	refs.Set("baz", &openapi.SchemaRef{})
	if len(refs) != 3 {
		t.Fatalf("expected 3, got %d", len(refs))
	}

	i := 0
	keys := []string{"foo", "bar", "baz"}
	for k := range refs.ByIndex() {
		if k != keys[i] {
			t.Fatalf("expected %s, got %s", keys[i], k)
		}
		i++
	}
}
