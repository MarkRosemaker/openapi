package ordmap

import (
	"fmt"

	"github.com/MarkRosemaker/errpath"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

// UnmarshalJSONV2 unmarshals the key-value pairs in order and sets the indices.
func (om *OrderedMap[K, V]) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	return UnmarshalJSONV2(om, dec, opts, setIndex)
}

// UnmarshalJSONV2 is a helper function to unmarshal an ordered map setting the indices in order.
func UnmarshalJSONV2[M ~map[K]R, K comparable, R any](
	m *M, dec *jsontext.Decoder, opts json.Options,
	setIndex func(R, int) R,
) error {
	tkn, err := dec.ReadToken()
	if err != nil {
		return err
	}

	if tkn.Kind() != '{' {
		return fmt.Errorf("expected {, got %s", tkn.Kind())
	}

	// create the map
	*m = M{}

	i := 1 // start at 1 to avoid confusion with zero values

	for {
		// check if we reached the end of the object
		if dec.PeekKind() == '}' {
			_, err := dec.ReadToken() // consume '}', should not fail
			return err
		}

		var key K
		if err := json.UnmarshalDecode(dec, &key, opts); err != nil {
			return err
		}

		var v R
		if err := json.UnmarshalDecode(dec, &v, opts); err != nil {
			return &errpath.ErrKey{Key: fmt.Sprint(key), Err: err}
		}

		// set the variable in the map with the proper index
		(*m)[K(key)] = setIndex(v, i)
		i++
	}
}
