package openapi

import (
	"iter"

	"github.com/MarkRosemaker/errpath"
	"github.com/MarkRosemaker/ordmap"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

// ServerVariables is an ordered map of server variable.
type ServerVariables map[string]*ServerVariable

// Validate validates each server variable.
func (vars ServerVariables) Validate() error {
	for k, v := range vars.ByIndex() {
		if err := v.Validate(); err != nil {
			return &errpath.ErrKey{Key: k, Err: err}
		}
	}

	return nil
}

// ByIndex returns a sequence of key-value pairs ordered by index.
func (vars ServerVariables) ByIndex() iter.Seq2[string, *ServerVariable] {
	return ordmap.ByIndex(vars, getIndexServerVariable)
}

// Sort sorts the map by key and sets the indices accordingly.
func (vars ServerVariables) Sort() {
	ordmap.Sort(vars, setIndexServerVariable)
}

// Set sets a value in the map, adding it at the end of the order.
func (vars *ServerVariables) Set(key string, v *ServerVariable) {
	ordmap.Set(vars, key, v, getIndexServerVariable, setIndexServerVariable)
}

// MarshalJSONV2 marshals the key-value pairs in order.
func (vars *ServerVariables) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	return ordmap.MarshalJSONV2(vars, enc, opts)
}

// UnmarshalJSONV2 unmarshals the key-value pairs in order and sets the indices.
func (vars *ServerVariables) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	return ordmap.UnmarshalJSONV2(vars, dec, opts, setIndexServerVariable)
}
