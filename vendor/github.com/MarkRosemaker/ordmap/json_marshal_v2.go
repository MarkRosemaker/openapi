package ordmap

import (
	"fmt"

	"github.com/MarkRosemaker/errpath"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

// MarshalJSONV2 marshals the key-value pairs in order.
func (om *OrderedMap[_, _]) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	return MarshalJSONV2(om, enc, opts)
}

// MarshalJSONV2 marshals an ordered map by encoding its key-value pairs in order.
func MarshalJSONV2[M ByIndexer[K, V], K comparable, V any](
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
