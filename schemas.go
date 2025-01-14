package openapi

import (
	"iter"

	"github.com/MarkRosemaker/errpath"
	"github.com/MarkRosemaker/ordmap"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type Schemas map[string]*Schema

func (ss Schemas) Validate() error {
	for name, s := range ss.ByIndex() {
		if err := validateKey(name); err != nil {
			return err
		}

		if err := s.Validate(); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}

// ByIndex returns a sequence of key-value pairs ordered by index.
func (rs Schemas) ByIndex() iter.Seq2[string, *Schema] {
	return ordmap.ByIndex(rs, getIndexSchema)
}

// Sort sorts the map by key and sets the indices accordingly.
func (rs Schemas) Sort() {
	ordmap.Sort(rs, setIndexSchema)
}

// Set sets a value in the map, adding it at the end of the order.
func (rs *Schemas) Set(key string, v *Schema) {
	ordmap.Set(rs, key, v, getIndexSchema, setIndexSchema)
}

// MarshalJSONTo marshals the key-value pairs in order.
func (rs *Schemas) MarshalJSONTo(enc *jsontext.Encoder, opts json.Options) error {
	return ordmap.MarshalJSONTo(rs, enc, opts)
}

// UnmarshalJSONFrom unmarshals the key-value pairs in order and sets the indices.
func (rs *Schemas) UnmarshalJSONFrom(dec *jsontext.Decoder, opts json.Options) error {
	return ordmap.UnmarshalJSONFrom(rs, dec, opts, setIndexSchema)
}

func (l *loader) collectSchemas(ss Schemas, ref ref) {
	for name, s := range ss.ByIndex() {
		l.collectSchema(s, append(ref, name))
	}
}

func (l *loader) resolveSchemas(ss Schemas) error {
	for name, s := range ss.ByIndex() {
		if err := l.resolveSchema(s); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}
