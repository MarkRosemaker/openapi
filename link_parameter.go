package openapi

import (
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

// LinkParameter is an expression that is the value of a parameter map in a Link Object.
type LinkParameter struct {
	Expression RuntimeExpression

	// an index to the original location of this object
	idx int
}

func getIndexLinkParameter(p *LinkParameter) int                     { return p.idx }
func setIndexLinkParameter(p *LinkParameter, idx int) *LinkParameter { p.idx = idx; return p }

// Validate validates the link parameter.
func (p *LinkParameter) Validate() error { return p.Expression.Validate() }

// UnmarshalJSONFrom unmarschals the link parameter into its appropriate type.
// NOTE: For now, we only implemented the case of it being a runtime expression.
func (p *LinkParameter) UnmarshalJSONFrom(dec *jsontext.Decoder, opts json.Options) error {
	return json.UnmarshalDecode(dec, &p.Expression, opts)
}

// MarshalJSONTo marschals the link parameter into its appropriate type.
// NOTE: For now, we only implemented the case of it being a runtime expression.
func (p *LinkParameter) MarshalJSONTo(enc *jsontext.Encoder, opts json.Options) error {
	return json.MarshalEncode(enc, p.Expression)
}
