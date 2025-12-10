package openapi

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"iter"

	"github.com/MarkRosemaker/ordmap"
)

type MapOfStrings map[string]String

// ByIndex returns a sequence of key-value pairs ordered by index.
func (scs MapOfStrings) ByIndex() iter.Seq2[string, String] {
	return ordmap.ByIndex(scs, getIndexScope)
}

// Sort sorts the map by key and sets the indices accordingly.
func (scs MapOfStrings) Sort() {
	ordmap.Sort(scs, setIndexScope)
}

// Set sets a value in the map, adding it at the end of the order.
func (scs *MapOfStrings) Set(key string, s String) {
	ordmap.Set(scs, key, s, getIndexScope, setIndexScope)
}

var _ json.MarshalerTo = (*MapOfStrings)(nil)

// MarshalJSONTo marshals the key-value pairs in order.
func (scs *MapOfStrings) MarshalJSONTo(enc *jsontext.Encoder) error {
	return ordmap.MarshalJSONTo(scs, enc)
}

var _ json.UnmarshalerFrom = (*MapOfStrings)(nil)

// UnmarshalJSONFrom unmarshals the key-value pairs in order and sets the indices.
func (scs *MapOfStrings) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	return ordmap.UnmarshalJSONFrom(scs, dec, setIndexScope)
}

type String struct {
	Value string

	idx int
}

var _ json.UnmarshalerFrom = (*String)(nil)

// UnmarshalJSONFrom unmarshals the value of the String.
func (s *String) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	return json.UnmarshalDecode(dec, &s.Value)
}

var _ json.MarshalerTo = (*String)(nil)

// MarshalJSONTo marshals the value of the String.
func (s *String) MarshalJSONTo(enc *jsontext.Encoder) error {
	return json.MarshalEncode(enc, s.Value)
}

func getIndexScope(s String) int           { return s.idx }
func setIndexScope(s String, i int) String { s.idx = i; return s }
