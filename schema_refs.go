package openapi

import (
	"iter"

	_json "github.com/MarkRosemaker/openapi/internal/json"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type SchemaRefs map[string]*SchemaRef

func (ss SchemaRefs) Validate() error {
	for name, value := range ss.ByIndex() {
		if err := value.Validate(); err != nil {
			return &ErrKey{Key: name, Err: err}
		}
	}

	return nil
}

// ByIndex returns the keys of the map in the order of the index.
func (ss SchemaRefs) ByIndex() iter.Seq2[string, *SchemaRef] {
	return _json.OrderedMapByIndex(ss, getIndexRef[Schema, *Schema])
}

// UnmarshalJSONV2 unmarshals the map from JSON and sets the index of each variable.
func (ss *SchemaRefs) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	return _json.UnmarshalOrderedMap(ss, dec, opts, setIndexRef[Schema, *Schema])
}

// MarshalJSONV2 marshals the map to JSON in the order of the index.
func (ss *SchemaRefs) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	return _json.MarshalOrderedMap(ss, enc, opts)
}
