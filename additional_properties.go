package openapi

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
)

// AdditionalProperties is a schema for the values of properties not in "properties", or a bare boolean allowing (true) or forbidding (false) them.
type AdditionalProperties struct {
	// Schema is the schema of the additional properties' values, when one is given.
	Schema *SchemaRef
	// Allowed is the value of the boolean form, only meaningful when Schema is nil.
	Allowed bool
}

// Validate validates the schema of the additional properties, if any.
func (ap *AdditionalProperties) Validate() error {
	if ap.Schema == nil {
		return nil
	}

	return ap.Schema.Validate()
}

var (
	_ json.UnmarshalerFrom = (*AdditionalProperties)(nil)
	_ json.MarshalerTo     = (*AdditionalProperties)(nil)
)

// UnmarshalJSONFrom reads a boolean into Allowed and anything else into Schema.
func (ap *AdditionalProperties) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	if k := dec.PeekKind(); k == jsontext.KindTrue || k == jsontext.KindFalse {
		tok, err := dec.ReadToken()
		if err != nil {
			return err
		}

		*ap = AdditionalProperties{Allowed: tok.Bool()}

		return nil
	}

	s := &SchemaRef{}
	if err := json.UnmarshalDecode(dec, s); err != nil {
		return err
	}

	*ap = AdditionalProperties{Schema: s}

	return nil
}

// MarshalJSONTo writes Schema if set, or else the boolean form.
func (ap *AdditionalProperties) MarshalJSONTo(enc *jsontext.Encoder) error {
	if ap.Schema == nil {
		return enc.WriteToken(jsontext.Bool(ap.Allowed))
	}

	return json.MarshalEncode(enc, ap.Schema)
}
