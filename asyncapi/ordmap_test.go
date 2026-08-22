package asyncapi_test

import (
	"testing"

	"github.com/MarkRosemaker/asyncapi"
	"github.com/MarkRosemaker/ordmap"
)

func TestOrderedMaps(t *testing.T) {
	t.Parallel()

	testSort[*asyncapi.Schemas](t)
	testSort[*asyncapi.Servers](t)
	testSort[*asyncapi.ServerVariables](t)
	testSort[*asyncapi.Channels](t)
	testSort[*asyncapi.Operations](t)
	testSort[*asyncapi.Messages](t)
	testSort[*asyncapi.MessageTraits](t)
	testSort[*asyncapi.OperationTraits](t)
	testSort[*asyncapi.Parameters](t)
	testSort[*asyncapi.CorrelationIDs](t)
	testSort[*asyncapi.Replies](t)
	testSort[*asyncapi.ReplyAddresses](t)
	testSort[*asyncapi.ExternalDocsByName](t)
	testSort[*asyncapi.TagsByName](t)
	testSort[*asyncapi.SecuritySchemes](t)
	testSort[*asyncapi.BindingsByName](t)
	testSort[*asyncapi.Bindings](t)
	testSortValues[*asyncapi.MapOfStrings](t)
}

// testSort checks that an ordered map of pointers keeps the order in which the keys were set
// and that it can be sorted by key.
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

	checkOrder(t, om, []K{"c", "a", "b"})

	om.Sort()

	checkOrder(t, om, []K{"a", "b", "c"})
}

// testSortValues is like testSort for an ordered map that holds values instead of pointers.
func testSortValues[MP interface {
	Set(K, V)
	*M
}, M interface {
	Sort()
	ordmap.ByIndexer[K, V]
}, K ~string, V any](t *testing.T,
) {
	t.Helper()

	var om M

	om.Sort() // no panic

	// set some values
	var a, b, c V
	MP(&om).Set("c", c)
	MP(&om).Set("a", a)
	MP(&om).Set("b", b)

	checkOrder(t, om, []K{"c", "a", "b"})

	om.Sort()

	checkOrder(t, om, []K{"a", "b", "c"})
}

func checkOrder[M ordmap.ByIndexer[K, V], K ~string, V any](t *testing.T, om M, want []K) {
	t.Helper()

	i := 0
	for k := range om.ByIndex() {
		if k != want[i] {
			t.Fatalf("got: %v, want: %v", k, want[i])
		}

		i++
	}

	if i != len(want) {
		t.Fatalf("got: %d keys, want: %d", i, len(want))
	}
}
