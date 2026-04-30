package jsonutil

import (
	"cmp"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"maps"
	"slices"
)

// OrderedMapMarshal is a custom marshaler for maps with ordered keys, marshaling them in an ordered fashion.
func OrderedMapMarshal[M ~map[K]V, K cmp.Ordered, V any](enc *jsontext.Encoder, m M) error {
	if m == nil {
		return enc.WriteToken(jsontext.Null)
	}

	if err := enc.WriteToken(jsontext.BeginObject); err != nil {
		return err
	}

	for _, key := range slices.Sorted(maps.Keys(m)) {
		if err := json.MarshalEncode(enc, key); err != nil {
			return err
		}

		if err := json.MarshalEncode(enc, m[key]); err != nil {
			return err
		}
	}

	return enc.WriteToken(jsontext.EndObject)
}
