package ordmap

import (
	"fmt"

	"github.com/MarkRosemaker/errpath"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

// MarshalJSONTo marshals the key-value pairs in order.
func (om *OrderedMap[_, _]) MarshalJSONTo(enc *jsontext.Encoder, opts json.Options) error {
	return MarshalJSONTo(om, enc, opts)
}

// MarshalJSONTo marshals an ordered map by encoding its key-value pairs in order.
func MarshalJSONTo[M ByIndexer[K, V], K comparable, V any](
	m M, enc *jsontext.Encoder, opts json.Options,
) error {
	if err := enc.WriteToken(jsontext.ObjectStart); err != nil {
		return err // should never fail
	}

	for k, v := range m.ByIndex() {
		if err := json.MarshalEncode(enc, k, opts); err != nil {
			return err
		}

		if err := json.MarshalEncode(enc, v, opts); err != nil {
			return &errpath.ErrKey{Key: fmt.Sprint(k), Err: err}
		}
	}

	return enc.WriteToken(jsontext.ObjectEnd)
}
