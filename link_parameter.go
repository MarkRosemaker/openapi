package openapi

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
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

var _ json.UnmarshalerFrom = (*LinkParameter)(nil)

// UnmarshalJSONFrom unmarschals the link parameter into its appropriate type.
// NOTE: For now, we only implemented the case of it being a runtime expression.
func (p *LinkParameter) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	return json.UnmarshalDecode(dec, &p.Expression)
}

var _ json.MarshalerTo = (*LinkParameter)(nil)

// MarshalJSONTo marschals the link parameter into its appropriate type.
// NOTE: For now, we only implemented the case of it being a runtime expression.
func (p *LinkParameter) MarshalJSONTo(enc *jsontext.Encoder) error {
	return json.MarshalEncode(enc, p.Expression)
}
