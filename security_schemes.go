package openapi

import (
	"iter"

	"github.com/MarkRosemaker/errpath"
	"github.com/MarkRosemaker/ordmap"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type SecuritySchemes map[string]*SecuritySchemeRef

func (ss SecuritySchemes) Validate() error {
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
func (ss SecuritySchemes) ByIndex() iter.Seq2[string, *SecuritySchemeRef] {
	return ordmap.ByIndex(ss, getIndexRef[SecurityScheme, *SecurityScheme])
}

// Sort sorts the map by key and sets the indices accordingly.
func (ss SecuritySchemes) Sort() {
	ordmap.Sort(ss, setIndexRef[SecurityScheme, *SecurityScheme])
}

// Set sets a value in the map, adding it at the end of the order.
func (ss *SecuritySchemes) Set(key string, v *SecuritySchemeRef) {
	ordmap.Set(ss, key, v, getIndexRef[SecurityScheme, *SecurityScheme], setIndexRef[SecurityScheme, *SecurityScheme])
}

// MarshalJSONV2 marshals the key-value pairs in order.
func (ss *SecuritySchemes) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	return ordmap.MarshalJSONV2(ss, enc, opts)
}

// UnmarshalJSONV2 unmarshals the key-value pairs in order and sets the indices.
func (ss *SecuritySchemes) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	return ordmap.UnmarshalJSONV2(ss, dec, opts, setIndexRef[SecurityScheme, *SecurityScheme])
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
