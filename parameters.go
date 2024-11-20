package openapi

import (
	"iter"

	_json "github.com/MarkRosemaker/openapi/internal/json"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type Parameters map[string]*ParameterRef

func (ps Parameters) Validate() error {
	for name, p := range ps.ByIndex() {
		if err := validateKey(name); err != nil {
			return err
		}

		if err := p.Validate(); err != nil {
			return &ErrKey{Key: name, Err: err}
		}
	}

	return nil
}

func (l *loader) collectParameters(ps Parameters, ref ref) {
	for name, p := range ps {
		l.collectParameterRef(p, append(ref, name))
	}
}

func (l *loader) resolveParameters(ps Parameters) error {
	for name, p := range ps.ByIndex() {
		if err := l.resolveParameterRef(p); err != nil {
			return &ErrKey{Key: name, Err: err}
		}
	}

	return nil
}

// ByIndex returns the keys of the map in the order of the index.
func (ss Parameters) ByIndex() iter.Seq2[string, *ParameterRef] {
	return _json.OrderedMapByIndex(ss, getIndexRef[Parameter, *Parameter])
}

// UnmarshalJSONV2 unmarshals the map from JSON and sets the index of each variable.
func (ss *Parameters) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	return _json.UnmarshalOrderedMap(ss, dec, opts, setIndexRef[Parameter, *Parameter])
}

// MarshalJSONV2 marshals the map to JSON in the order of the index.
func (ss *Parameters) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	return _json.MarshalOrderedMap(ss, enc, opts)
}
