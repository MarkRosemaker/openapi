package openapi

import (
	"iter"

	_json "github.com/MarkRosemaker/openapi/internal/json"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

// Examples is a map of examples.
type Examples map[string]*ExampleRef

// Validate validates the map of examples.
func (ex Examples) Validate() error {
	for k, v := range ex.ByIndex() {
		if err := v.Validate(); err != nil {
			return &ErrKey{Key: k, Err: err}
		}
	}

	return nil
}

// ByIndex returns the keys of the map in the order of the index.
func (ex Examples) ByIndex() iter.Seq2[string, *ExampleRef] {
	return _json.OrderedMapByIndex(ex, getIndexRef[Example, *Example])
}

// UnmarshalJSONV2 unmarshals the map from JSON and sets the index of each variable.
func (ex *Examples) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	return _json.UnmarshalOrderedMap(ex, dec, opts, setIndexRef[Example, *Example])
}

// MarshalJSONV2 marshals the map to JSON in the order of the index.
func (ex *Examples) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	return _json.MarshalOrderedMap(ex, enc, opts)
}
