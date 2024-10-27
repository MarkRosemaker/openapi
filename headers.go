package openapi

import (
	"iter"

	_json "github.com/MarkRosemaker/openapi/internal/json"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type Headers map[string]*HeaderRef

func (hs Headers) Validate() error {
	for k, h := range hs.ByIndex() {
		if err := h.Validate(); err != nil {
			return &ErrKey{Key: k, Err: err}
		}
	}

	return nil
}

// ByIndex returns the keys of the map in the order of the index.
func (h Headers) ByIndex() iter.Seq2[string, *HeaderRef] {
	return _json.OrderedMapByIndex(h, getIndexRef[Header, *Header])
}

// UnmarshalJSONV2 unmarshals the map from JSON and sets the index of each variable.
func (h *Headers) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	return _json.UnmarshalOrderedMap(h, dec, opts, setIndexRef[Header, *Header])
}

// MarshalJSONV2 marshals the map to JSON in the order of the index.
func (h *Headers) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	return _json.MarshalOrderedMap(h, enc, opts)
}
