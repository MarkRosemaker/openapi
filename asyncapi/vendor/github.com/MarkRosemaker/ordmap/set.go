package ordmap

// Set sets a value in the map, adding it at the end of the order.
func (om *OrderedMap[K, V]) Set(key K, v V) {
	Set(om, key, Value[V]{V: v}, getIndex[V], setIndex[V])
}

// Set is a helper function to set a value in the map, adding it at the end of the order.
func Set[M ~map[K]V, K comparable, V any](
	m *M, key K, v V,
	getIndex func(V) int,
	setIndex func(V, int) V,
) {
	// check if the map is nil and create it if it is
	if *m == nil {
		*m = M{key: setIndex(v, 1)}
		return
	}

	highestIdx := 0
	for _, v := range *m {
		if idx := getIndex(v); idx > highestIdx {
			highestIdx = idx
		}
	}

	(*m)[key] = setIndex(v, highestIdx+1)
}
