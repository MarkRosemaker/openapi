package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
	"github.com/MarkRosemaker/ordmap"
)

func TestOrderedMaps(t *testing.T) {
	testSort[*openapi.CallbackRefs](t)
	testSort[*openapi.Callback](t)
	testSort[*openapi.Content](t)
	testSort[*openapi.Encodings](t)
	testSort[*openapi.Examples](t)
	testSort[*openapi.Headers](t)
	testSort[*openapi.LinkParameters](t)
	testSort[*openapi.Links](t)
	testSort[*openapi.Scopes](t)
	testSort[*openapi.Parameters](t)
	testSort[*openapi.PathItems](t)
	testSort[*openapi.Paths](t)
	testSort[*openapi.RequestBodies](t)
	testSort[*openapi.Responses[string]](t)
	testSort[*openapi.Schemas](t)
	testSort[*openapi.SchemaRefs](t)
	testSort[*openapi.ServerVariables](t)
}

func testSort[MP interface {
	Set(K, *V)
	*M
}, M interface {
	Sort()
	ordmap.ByIndexer[K, *V]
}, K ~string, V any](t *testing.T,
) {
	t.Helper()

	var om M

	om.Sort() // no panic

	// set some values
	var a, b, c V
	MP(&om).Set("c", &c)
	MP(&om).Set("a", &a)
	MP(&om).Set("b", &b)

	keys := []K{"c", "a", "b"}

	i := 0
	for k := range om.ByIndex() {
		if k != keys[i] {
			t.Fatalf("got: %v, want: %v", k, keys[i])
		}

		i++
	}

	om.Sort()

	keysSorted := []K{"a", "b", "c"}

	i = 0
	for k := range om.ByIndex() {
		if k != keysSorted[i] {
			t.Fatalf("got: %v, want: %v", k, keysSorted[i])
		}

		i++
	}
}
