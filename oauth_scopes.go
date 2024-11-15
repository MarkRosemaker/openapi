package openapi

import (
	"iter"

	_json "github.com/MarkRosemaker/openapi/internal/json"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type Scopes map[string]*String

// ByIndex returns the keys of the map in the order of the index.
func (ps Scopes) ByIndex() iter.Seq2[string, *String] {
	return _json.OrderedMapByIndex(ps, getIndexScope)
}

// UnmarshalJSONV2 unmarshals the map from JSON and sets the index of each variable.
func (ps *Scopes) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	return _json.UnmarshalOrderedMap(ps, dec, opts, setIndexScope)
}

// MarshalJSONV2 marshals the map to JSON in the order of the index.
func (ps *Scopes) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	return _json.MarshalOrderedMap(ps, enc, opts)
}

type String struct {
	Value string

	idx int
}

// UnmarshalJSONV2 unmarshals the value of the String.
func (s *String) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	return json.UnmarshalDecode(dec, &s.Value, opts)
}

// MarshalJSONV2 marshals the value of the String.
func (s *String) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	return json.MarshalEncode(enc, s.Value, opts)
}

func getIndexScope(s *String) int    { return s.idx }
func setIndexScope(s *String, i int) { s.idx = i }
