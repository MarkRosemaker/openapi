package openapi

import (
	"iter"

	"github.com/MarkRosemaker/errpath"
	"github.com/MarkRosemaker/ordmap"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
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

// MarshalJSONTo marshals the key-value pairs in order.
func (es *Encodings) MarshalJSONTo(enc *jsontext.Encoder, opts json.Options) error {
	return ordmap.MarshalJSONTo(es, enc, opts)
}

// UnmarshalJSONFrom unmarshals the key-value pairs in order and sets the indices.
func (es *Encodings) UnmarshalJSONFrom(dec *jsontext.Decoder, opts json.Options) error {
	return ordmap.UnmarshalJSONFrom(es, dec, opts, setIndexEncoding)
}

func (l *loader) resolveEncodings(es Encodings) error {
	for k, e := range es.ByIndex() {
		if err := l.resolveEncoding(e); err != nil {
			return &errpath.ErrKey{Key: k, Err: err}
		}
	}

	return nil
}
