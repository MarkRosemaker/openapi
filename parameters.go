package openapi

import (
	"iter"

	_json "github.com/MarkRosemaker/openapi/internal/json"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type Parameters map[string]*ParameterRef

func (ss Parameters) Validate() error {
	for name, value := range ss.ByIndex() {
		if err := validateKey(name); err != nil {
			return err
		}

		if err := value.Validate(); err != nil {
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
