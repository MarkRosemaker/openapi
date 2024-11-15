package openapi

import (
	"iter"

	_json "github.com/MarkRosemaker/openapi/internal/json"
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
			return &ErrKey{Key: name, Err: err}
		}
	}

	return nil
}

// ByIndex returns the keys of the map in the order of the index.
func (ss SecuritySchemes) ByIndex() iter.Seq2[string, *SecuritySchemeRef] {
	return _json.OrderedMapByIndex(ss, getIndexRef[SecurityScheme, *SecurityScheme])
}

// UnmarshalJSONV2 unmarshals the map from JSON and sets the index of each variable.
func (ss *SecuritySchemes) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	return _json.UnmarshalOrderedMap(ss, dec, opts, setIndexRef[SecurityScheme, *SecurityScheme])
}

// MarshalJSONV2 marshals the map to JSON in the order of the index.
func (c *SecuritySchemes) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	return _json.MarshalOrderedMap(c, enc, opts)
}
