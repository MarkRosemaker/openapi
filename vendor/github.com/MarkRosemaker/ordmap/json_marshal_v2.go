package ordmap

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"

	"github.com/MarkRosemaker/errpath"
)

var _ json.MarshalerTo = (*OrderedMap[string, any])(nil)

// MarshalJSONTo marshals the key-value pairs in order.
func (om *OrderedMap[_, _]) MarshalJSONTo(enc *jsontext.Encoder) error {
	return MarshalJSONTo(om, enc)
}

// MarshalJSONTo marshals an ordered map by encoding its key-value pairs in order.
func MarshalJSONTo[M ByIndexer[K, V], K comparable, V any](
	m M, enc *jsontext.Encoder,
) error {
	if err := enc.WriteToken(jsontext.BeginObject); err != nil {
		return err // should never fail
	}

	for k, v := range m.ByIndex() {
		if err := json.MarshalEncode(enc, k, enc.Options()); err != nil {
			return err
		}

		if err := json.MarshalEncode(enc, v, enc.Options()); err != nil {
			return &errpath.ErrKey{Key: fmt.Sprint(k), Err: err}
		}
	}

	return enc.WriteToken(jsontext.EndObject)
}
