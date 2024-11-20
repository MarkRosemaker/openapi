package openapi

func setToMap[M ~map[K]V, K comparable, V any](
	m *M, key K, v V,
	getIndex func(V) int,
	setIndex func(V, int),
) {
	if *m == nil {
		setIndex(v, 1)
		*m = M{key: v}
		return
	}

	highestIdx := 0
	for _, v := range *m {
		if idx := getIndex(v); idx > highestIdx {
			highestIdx = idx
		}
	}

	setIndex(v, highestIdx+1)
	(*m)[key] = v
}
