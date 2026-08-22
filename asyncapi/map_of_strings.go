package asyncapi

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"iter"

	"github.com/MarkRosemaker/ordmap"
)

// MapOfStrings is an ordered map of strings, e.g. the available scopes of an OAuth flow,
// which are "a map between the scope name and a short description for it".
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#oauthFlowObject
type MapOfStrings map[string]String

// ByIndex returns a sequence of key-value pairs ordered by index.
func (ss MapOfStrings) ByIndex() iter.Seq2[string, String] {
	return ordmap.ByIndex(ss, getIndexString)
}

// Sort sorts the map by key and sets the indices accordingly.
func (ss MapOfStrings) Sort() {
	ordmap.Sort(ss, setIndexString)
}

// Set sets a value in the map, adding it at the end of the order.
func (ss *MapOfStrings) Set(key string, s String) {
	ordmap.Set(ss, key, s, getIndexString, setIndexString)
}

var _ json.MarshalerTo = (*MapOfStrings)(nil)

// MarshalJSONTo marshals the key-value pairs in order.
func (ss *MapOfStrings) MarshalJSONTo(enc *jsontext.Encoder) error {
	return ordmap.MarshalJSONTo(ss, enc)
}

var _ json.UnmarshalerFrom = (*MapOfStrings)(nil)

// UnmarshalJSONFrom unmarshals the key-value pairs in order and sets the indices.
func (ss *MapOfStrings) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	return ordmap.UnmarshalJSONFrom(ss, dec, setIndexString)
}

// String is a string value that remembers its position in an ordered map.
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

func getIndexString(s String) int           { return s.idx }
func setIndexString(s String, i int) String { s.idx = i; return s }
