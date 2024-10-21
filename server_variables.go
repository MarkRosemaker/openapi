package openapi

import (
	"iter"

	_json "github.com/MarkRosemaker/openapi/internal/json"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

// ServerVariables is an ordered map of server variable.
type ServerVariables map[string]*ServerVariable

// ByIndex returns the keys of the map in the order of the index.
func (vars ServerVariables) ByIndex() iter.Seq2[string, *ServerVariable] {
	return _json.OrderedMapByIndex(vars, getIndexServerVariable)
}

// Validate validates each server variable.
func (vars ServerVariables) Validate() error {
	for k, v := range vars.ByIndex() {
		if err := v.Validate(); err != nil {
			return &ErrKey{Key: k, Err: err}
		}
	}

	return nil
}

// UnmarshalJSONV2 unmarshals the map from JSON and sets the index of each variable.
func (sv *ServerVariables) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	return _json.UnmarshalOrderedMap(sv, dec, opts, setIndexServerVariable)
}

// MarshalJSONV2 marshals the map to JSON in the order of the index.
func (sv ServerVariables) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	return _json.MarshalOrderedMap(&sv, enc, opts)
}
