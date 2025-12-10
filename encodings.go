package openapi

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"iter"

	"github.com/MarkRosemaker/errpath"
	"github.com/MarkRosemaker/ordmap"
)

// Encodings is a map between a property name and its encoding information.
type Encodings map[string]*Encoding

func (es Encodings) Validate() error {
	for k, e := range es.ByIndex() {
		if err := e.Validate(); err != nil {
			return &errpath.ErrKey{Key: k, Err: err}
		}
	}

	return nil
}

// ByIndex returns a sequence of key-value pairs ordered by index.
func (es Encodings) ByIndex() iter.Seq2[string, *Encoding] {
	return ordmap.ByIndex(es, getIndexEncoding)
}

// Sort sorts the map by key and sets the indices accordingly.
func (es Encodings) Sort() {
	ordmap.Sort(es, setIndexEncoding)
}

// Set sets a value in the map, adding it at the end of the order.
func (es *Encodings) Set(key string, e *Encoding) {
	ordmap.Set(es, key, e, getIndexEncoding, setIndexEncoding)
}

var _ json.MarshalerTo = (*Encodings)(nil)

// MarshalJSONTo marshals the key-value pairs in order.
func (es *Encodings) MarshalJSONTo(enc *jsontext.Encoder) error {
	return ordmap.MarshalJSONTo(es, enc)
}

var _ json.UnmarshalerFrom = (*Encodings)(nil)

// UnmarshalJSONFrom unmarshals the key-value pairs in order and sets the indices.
func (es *Encodings) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	return ordmap.UnmarshalJSONFrom(es, dec, setIndexEncoding)
}

func (l *loader) resolveEncodings(es Encodings) error {
	for k, e := range es.ByIndex() {
		if err := l.resolveEncoding(e); err != nil {
			return &errpath.ErrKey{Key: k, Err: err}
		}
	}

	return nil
}
