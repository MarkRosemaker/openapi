package openapi

import (
	"iter"

	_json "github.com/MarkRosemaker/openapi/internal/json"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type PathItems map[string]*PathItemRef

func (ps PathItems) Validate() error {
	for name, p := range ps.ByIndex() {
		if err := p.Validate(); err != nil {
			return &ErrKey{Key: name, Err: err}
		}
	}

	return nil
}

// ByIndex returns the keys of the map in the order of the index.
func (ps PathItems) ByIndex() iter.Seq2[string, *PathItemRef] {
	return _json.OrderedMapByIndex(ps, getIndexRef[PathItem, *PathItem])
}

// UnmarshalJSONV2 unmarshals the map from JSON and sets the index of each variable.
func (ps *PathItems) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	return _json.UnmarshalOrderedMap(ps, dec, opts, setIndexRef[PathItem, *PathItem])
}

// MarshalJSONV2 marshals the map to JSON in the order of the index.
func (ps *PathItems) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	return _json.MarshalOrderedMap(ps, enc, opts)
}
