package openapi

import (
	"iter"

	_json "github.com/MarkRosemaker/openapi/internal/json"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

// A map representing parameters to pass to an operation as specified with `operationId` or identified via `operationRef`. The key is the parameter name to be used, whereas the value can be a constant or an expression to be evaluated and passed to the linked operation.
// The parameter name can be qualified using the parameter location `[{in}.]{name}` for operations that use the same parameter name in different locations (e.g. path.id).
type LinkParameters map[string]*LinkParameter

func (ps LinkParameters) Validate() error {
	for name, p := range ps.ByIndex() {
		if err := p.Validate(); err != nil {
			return &ErrKey{Key: name, Err: err}
		}
	}

	return nil
}

// ByIndex returns the keys of the map in the order of the index.
func (ps LinkParameters) ByIndex() iter.Seq2[string, *LinkParameter] {
	return _json.OrderedMapByIndex(ps, getIndexLinkParameter)
}

// UnmarshalJSONV2 unmarshals the map from JSON and sets the index of each variable.
func (ps *LinkParameters) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	return _json.UnmarshalOrderedMap(ps, dec, opts, setIndexLinkParameter)
}

// MarshalJSONV2 marshals the map to JSON in the order of the index.
func (ps *LinkParameters) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	return _json.MarshalOrderedMap(ps, enc, opts)
}
