package ordmap

import (
	"iter"
	"sort"

	"golang.org/x/exp/maps"
)

// ByIndexer is an interface for an ordered map that implements an iterator.
type ByIndexer[K comparable, V any] interface {
	// ByIndex returns a sequence of key-value pairs sorted by index.
	ByIndex() iter.Seq2[K, V]
}

// ByIndex is a helper function for an ordered map to implement an iterator.
func ByIndex[M ~map[K]V, K comparable, V any](m M, getIndex func(V) int) iter.Seq2[K, V] {
	// get the keys and sort them by index
	keys := maps.Keys(m)
	sort.Slice(keys, func(i, j int) bool {
		idxI := getIndex(m[keys[i]])
		idxJ := getIndex(m[keys[j]])
		return idxI != 0 && // if i is not initialized, it should be at the end
			(idxJ == 0 || // if j is not initialized, it should be at the end
				idxI < idxJ) // otherwise, sort by index
	})

	return func(yield func(K, V) bool) {
		for _, k := range keys {
			if !yield(k, m[k]) {
				return
			}
		}
	}
}

// ByIndex returns a sequence of key-value pairs sorted by index.
func (om OrderedMap[K, V]) ByIndex() iter.Seq2[K, V] {
	// get the keys and sort them by index
	keys := maps.Keys(om)
	sort.Slice(keys, func(i, j int) bool {
		return om[keys[i]].idx != 0 && // if i is not initialized, it should be at the end
			(om[keys[j]].idx == 0 || // if j is not initialized, it should be at the end
				om[keys[i]].idx < om[keys[j]].idx) // otherwise, sort by index
	})

	return func(yield func(K, V) bool) {
		for _, k := range keys {
			if !yield(k, om[k].V) {
				return
			}
		}
	}

	// NOTE: The above is equivalent to the following:
	// return func(yield func(K, V) bool) {
	// 	for k, v := range ByIndex(om, getIndex) {
	// 		if !yield(k, v.V) {
	// 			return
	// 		}
	// 	}
	// }
}
