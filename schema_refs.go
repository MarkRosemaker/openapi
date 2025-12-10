package openapi

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"iter"

	"github.com/MarkRosemaker/errpath"
	"github.com/MarkRosemaker/ordmap"
)

type SchemaRefs map[string]*SchemaRef

func (ss SchemaRefs) Validate() error {
	for name, s := range ss.ByIndex() {
		if err := s.Validate(); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}

// ByIndex returns a sequence of key-value pairs ordered by index.
func (ss SchemaRefs) ByIndex() iter.Seq2[string, *SchemaRef] {
	return ordmap.ByIndex(ss, getIndexRef[Schema, *Schema])
}

// Sort sorts the map by key and sets the indices accordingly.
func (ss SchemaRefs) Sort() {
	ordmap.Sort(ss, setIndexRef[Schema, *Schema])
}

// Set sets a value in the map, adding it at the end of the order.
func (ss *SchemaRefs) Set(key string, v *SchemaRef) {
	ordmap.Set(ss, key, v, getIndexRef[Schema, *Schema], setIndexRef[Schema, *Schema])
}

var _ json.MarshalerTo = (*SchemaRefs)(nil)

// MarshalJSONTo marshals the key-value pairs in order.
func (ss *SchemaRefs) MarshalJSONTo(enc *jsontext.Encoder) error {
	return ordmap.MarshalJSONTo(ss, enc)
}

var _ json.UnmarshalerFrom = (*SchemaRefs)(nil)

// UnmarshalJSONFrom unmarshals the key-value pairs in order and sets the indices.
func (ss *SchemaRefs) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	return ordmap.UnmarshalJSONFrom(ss, dec, setIndexRef[Schema, *Schema])
}

func (l *loader) resolveSchemaRefs(ss SchemaRefs) error {
	for name, value := range ss.ByIndex() {
		if err := l.resolveSchemaRef(value); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}
