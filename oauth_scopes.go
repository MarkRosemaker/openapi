package openapi

import (
	"iter"

	"github.com/MarkRosemaker/ordmap"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type Scopes map[string]*String

// ByIndex returns a sequence of key-value pairs ordered by index.
func (scs Scopes) ByIndex() iter.Seq2[string, *String] {
	return ordmap.ByIndex(scs, getIndexScope)
}

// Sort sorts the map by key and sets the indices accordingly.
func (scs Scopes) Sort() {
	ordmap.Sort(scs, setIndexScope)
}

// Set sets a value in the map, adding it at the end of the order.
func (scs *Scopes) Set(key string, s *String) {
	ordmap.Set(scs, key, s, getIndexScope, setIndexScope)
}

// MarshalJSONTo marshals the key-value pairs in order.
func (scs *Scopes) MarshalJSONTo(enc *jsontext.Encoder, opts json.Options) error {
	return ordmap.MarshalJSONTo(scs, enc, opts)
}

// UnmarshalJSONFrom unmarshals the key-value pairs in order and sets the indices.
func (scs *Scopes) UnmarshalJSONFrom(dec *jsontext.Decoder, opts json.Options) error {
	return ordmap.UnmarshalJSONFrom(scs, dec, opts, setIndexScope)
}

type String struct {
	Value string

	idx int
}

// UnmarshalJSONFrom unmarshals the value of the String.
func (s *String) UnmarshalJSONFrom(dec *jsontext.Decoder, opts json.Options) error {
	return json.UnmarshalDecode(dec, &s.Value, opts)
}

// MarshalJSONTo marshals the value of the String.
func (s *String) MarshalJSONTo(enc *jsontext.Encoder, opts json.Options) error {
	return json.MarshalEncode(enc, s.Value, opts)
}

func getIndexScope(s *String) int            { return s.idx }
func setIndexScope(s *String, i int) *String { s.idx = i; return s }
