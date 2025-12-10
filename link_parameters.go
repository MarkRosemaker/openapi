package openapi

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"iter"

	"github.com/MarkRosemaker/errpath"
	"github.com/MarkRosemaker/ordmap"
)

// A map representing parameters to pass to an operation as specified with `operationId` or identified via `operationRef`. The key is the parameter name to be used, whereas the value can be a constant or an expression to be evaluated and passed to the linked operation.
// The parameter name can be qualified using the parameter location `[{in}.]{name}` for operations that use the same parameter name in different locations (e.g. path.id).
type LinkParameters map[string]*LinkParameter

func (ps LinkParameters) Validate() error {
	for name, p := range ps.ByIndex() {
		if err := p.Validate(); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}

// ByIndex returns a sequence of key-value pairs ordered by index.
func (ps LinkParameters) ByIndex() iter.Seq2[string, *LinkParameter] {
	return ordmap.ByIndex(ps, getIndexLinkParameter)
}

// Sort sorts the map by key and sets the indices accordingly.
func (ps LinkParameters) Sort() {
	ordmap.Sort(ps, setIndexLinkParameter)
}

// Set sets a value in the map, adding it at the end of the order.
func (ps *LinkParameters) Set(key string, p *LinkParameter) {
	ordmap.Set(ps, key, p, getIndexLinkParameter, setIndexLinkParameter)
}

var _ json.MarshalerTo = (*LinkParameters)(nil)

// MarshalJSONTo marshals the key-value pairs in order.
func (ps *LinkParameters) MarshalJSONTo(enc *jsontext.Encoder) error {
	return ordmap.MarshalJSONTo(ps, enc)
}

var _ json.UnmarshalerFrom = (*LinkParameters)(nil)

// UnmarshalJSONFrom unmarshals the key-value pairs in order and sets the indices.
func (ps *LinkParameters) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	return ordmap.UnmarshalJSONFrom(ps, dec, setIndexLinkParameter)
}
