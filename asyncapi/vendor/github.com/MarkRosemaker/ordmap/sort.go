package ordmap

import (
	"cmp"
	"slices"

	"golang.org/x/exp/maps"
)

// Sort sorts the map by key using a custom comparison function and sets the indices accordingly.
func (om OrderedMap[K, V]) Sort(less func(K, K) int) {
	SortFunc(om, setIndex, less)
}

// Sort is a helper function to sort a map by key and set the indices accordingly.
func Sort[M ~map[K]V, K cmp.Ordered, V any](m M, setIndex func(V, int) V) {
	doSort(m, setIndex, slices.Sort)
}

// SortFunc is a helper function to sort a map by key using a custom comparison function and set the indices accordingly.
func SortFunc[M ~map[K]V, K comparable, V any](m M, setIndex func(V, int) V, less func(K, K) int) {
	doSort(m, setIndex, func(keys []K) { slices.SortFunc(keys, less) })
}

func doSort[M ~map[K]V, K comparable, V any](m M, setIndex func(V, int) V, sortKeys func([]K)) {
	keys := maps.Keys(m)
	sortKeys(keys)

	for i, key := range keys {
		m[key] = setIndex(m[key], i+1)
	}
}
