package ordmap

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
)

var (
	_ json.UnmarshalerFrom = (*Value[any])(nil)
	_ json.MarshalerTo     = (*Value[any])(nil)
)

// Value is a value with an index.
type Value[V any] struct {
	V   V
	idx int
}

// UnmarshalJSONFrom unmarshals a value by just decoding the value.
// The index is set by the caller.
func (cs *Value[_]) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	return json.UnmarshalDecode(dec, &cs.V, dec.Options())
}

// MarshalJSONTo marshals a value by encoding just the value and ignoring the index.
func (v Value[_]) MarshalJSONTo(enc *jsontext.Encoder) error {
	return json.MarshalEncode(enc, v.V, enc.Options())
}

// getIndex returns the index of a value.
func getIndex[V any](v Value[V]) int { return v.idx }

// setIndex sets the index of a value.
func setIndex[V any](v Value[V], i int) Value[V] { v.idx = i; return v }
