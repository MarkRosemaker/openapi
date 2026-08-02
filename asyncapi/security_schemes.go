package asyncapi

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"iter"

	"github.com/MarkRosemaker/errpath"
	"github.com/MarkRosemaker/ordmap"
)

// SecuritySchemes is a map of Security Scheme Objects.
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#componentsSecuritySchemes
type SecuritySchemes map[string]*SecuritySchemeRef

// Validate validates each security scheme.
func (ss SecuritySchemes) Validate() error {
	for name, s := range ss.ByIndex() {
		if err := s.Validate(); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}

// ByIndex returns a sequence of key-value pairs ordered by index.
func (ss SecuritySchemes) ByIndex() iter.Seq2[string, *SecuritySchemeRef] {
	return ordmap.ByIndex(ss, getIndexRef[SecurityScheme, *SecurityScheme])
}

// Sort sorts the map by key and sets the indices accordingly.
func (ss SecuritySchemes) Sort() {
	ordmap.Sort(ss, setIndexRef[SecurityScheme, *SecurityScheme])
}

// Set sets a value in the map, adding it at the end of the order.
func (ss *SecuritySchemes) Set(key string, s *SecuritySchemeRef) {
	ordmap.Set(ss, key, s, getIndexRef[SecurityScheme, *SecurityScheme], setIndexRef[SecurityScheme, *SecurityScheme])
}

var _ json.MarshalerTo = (*SecuritySchemes)(nil)

// MarshalJSONTo marshals the key-value pairs in order.
func (ss *SecuritySchemes) MarshalJSONTo(enc *jsontext.Encoder) error {
	return ordmap.MarshalJSONTo(ss, enc)
}

var _ json.UnmarshalerFrom = (*SecuritySchemes)(nil)

// UnmarshalJSONFrom unmarshals the key-value pairs in order and sets the indices.
func (ss *SecuritySchemes) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	return ordmap.UnmarshalJSONFrom(ss, dec, setIndexRef[SecurityScheme, *SecurityScheme])
}

func (l *loader) collectSecuritySchemes(ss SecuritySchemes, ref ref) {
	for name, s := range ss.ByIndex() {
		l.collectSecuritySchemeRef(s, append(ref, name))
	}
}

func (l *loader) resolveSecuritySchemes(ss SecuritySchemes) error {
	for name, s := range ss.ByIndex() {
		if err := l.resolveSecuritySchemeRef(s); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}
