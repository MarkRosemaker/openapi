package openapi

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"iter"

	"github.com/MarkRosemaker/errpath"
	"github.com/MarkRosemaker/ordmap"
)

type SecuritySchemes map[SecuritySchemeName]*SecuritySchemeRef

func (ss SecuritySchemes) Validate() error {
	for name, s := range ss.ByIndex() {
		if err := validateKey(string(name)); err != nil {
			return err
		}

		if err := s.Validate(); err != nil {
			return &errpath.ErrKey{Key: string(name), Err: err}
		}
	}

	return nil
}

// KeyFunc returns the first key k satisfying f(ss[k]),
// or "" if none do.
func (ss SecuritySchemes) KeyFunc(f func(*SecurityScheme) bool) SecuritySchemeName {
	for name, s := range ss.ByIndex() {
		if f(s.Value) {
			return name
		}
	}

	return ""
}

// ByIndex returns a sequence of key-value pairs ordered by index.
func (ss SecuritySchemes) ByIndex() iter.Seq2[SecuritySchemeName, *SecuritySchemeRef] {
	return ordmap.ByIndex(ss, getIndexRef[SecurityScheme, *SecurityScheme])
}

// Sort sorts the map by key and sets the indices accordingly.
func (ss SecuritySchemes) Sort() {
	ordmap.Sort(ss, setIndexRef[SecurityScheme, *SecurityScheme])
}

// Set sets a value in the map, adding it at the end of the order.
func (ss *SecuritySchemes) Set(key SecuritySchemeName, v *SecuritySchemeRef) {
	ordmap.Set(ss, key, v, getIndexRef[SecurityScheme, *SecurityScheme], setIndexRef[SecurityScheme, *SecurityScheme])
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
		l.collectSecuritySchemeRef(s, append(ref, string(name)))
	}
}

func (l *loader) resolveSecuritySchemes(ss SecuritySchemes) error {
	for name, s := range ss.ByIndex() {
		if err := l.resolveSecuritySchemeRef(s); err != nil {
			return &errpath.ErrKey{Key: string(name), Err: err}
		}
	}

	return nil
}
