package asyncapi

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"iter"

	"github.com/MarkRosemaker/errpath"
	"github.com/MarkRosemaker/ordmap"
)

// Schemas is a map of schema definitions, each of which is either a schema object,
// a multi format schema object or a reference to one of them.
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#componentsSchemas
type Schemas map[string]*AnySchemaRef

// Validate validates each schema.
func (ss Schemas) Validate() error {
	for name, s := range ss.ByIndex() {
		if err := s.Validate(); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}

// ByIndex returns a sequence of key-value pairs ordered by index.
func (ss Schemas) ByIndex() iter.Seq2[string, *AnySchemaRef] {
	return ordmap.ByIndex(ss, getIndexRef[AnySchema, *AnySchema])
}

// Sort sorts the map by key and sets the indices accordingly.
func (ss Schemas) Sort() {
	ordmap.Sort(ss, setIndexRef[AnySchema, *AnySchema])
}

// Set sets a value in the map, adding it at the end of the order.
func (ss *Schemas) Set(key string, s *AnySchemaRef) {
	ordmap.Set(ss, key, s, getIndexRef[AnySchema, *AnySchema], setIndexRef[AnySchema, *AnySchema])
}

var _ json.MarshalerTo = (*Schemas)(nil)

// MarshalJSONTo marshals the key-value pairs in order.
func (ss *Schemas) MarshalJSONTo(enc *jsontext.Encoder) error {
	return ordmap.MarshalJSONTo(ss, enc)
}

var _ json.UnmarshalerFrom = (*Schemas)(nil)

// UnmarshalJSONFrom unmarshals the key-value pairs in order and sets the indices.
func (ss *Schemas) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	return ordmap.UnmarshalJSONFrom(ss, dec, setIndexRef[AnySchema, *AnySchema])
}

func (l *loader) collectSchemas(ss Schemas, ref ref) {
	for name, s := range ss.ByIndex() {
		l.collectAnySchemaRef(s, append(ref, name))
	}
}

func (l *loader) resolveSchemas(ss Schemas) error {
	for name, s := range ss.ByIndex() {
		if err := l.resolveAnySchemaRef(s); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}
