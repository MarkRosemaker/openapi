package openapi

import (
	"iter"

	_json "github.com/MarkRosemaker/openapi/internal/json"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type Links map[string]*LinkRef

func (l Links) Validate() error {
	for expr, v := range l.ByIndex() {
		if err := v.Validate(); err != nil {
			return &ErrField{Field: expr, Err: err}
		}
	}

	return nil
}

// ByIndex returns the keys of the map in the order of the index.
func (l Links) ByIndex() iter.Seq2[string, *LinkRef] {
	return _json.OrderedMapByIndex(l, getIndexRef[Link, *Link])
}

// UnmarshalJSONV2 unmarshals the map from JSON and sets the index of each variable.
func (l *Links) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	return _json.UnmarshalOrderedMap(l, dec, opts, setIndexRef[Link, *Link])
}

// MarshalJSONV2 marshals the map to JSON in the order of the index.
func (l *Links) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	return _json.MarshalOrderedMap(l, enc, opts)
}
