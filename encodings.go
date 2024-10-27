package openapi

import (
	"iter"

	_json "github.com/MarkRosemaker/openapi/internal/json"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

// Encodings is a map between a property name and its encoding information.
type Encodings map[string]*Encoding

func (es Encodings) Validate() error {
	for k, e := range es.ByIndex() {
		if err := e.Validate(); err != nil {
			return &ErrKey{Key: k, Err: err}
		}
	}

	return nil
}

// ByIndex returns the keys of the map in the order of the index.
func (e Encodings) ByIndex() iter.Seq2[string, *Encoding] {
	return _json.OrderedMapByIndex(e, getIndexEncoding)
}

// UnmarshalJSONV2 unmarshals the map from JSON and sets the index of each variable.
func (c *Encodings) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	return _json.UnmarshalOrderedMap(c, dec, opts, setIndexEncoding)
}

// MarshalJSONV2 marshals the map to JSON in the order of the index.
func (c *Encodings) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	return _json.MarshalOrderedMap(c, enc, opts)
}
