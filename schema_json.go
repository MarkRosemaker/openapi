package openapi

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
)

// schemaFields is Schema without its methods, so (un)marshaling it doesn't recurse.
type schemaFields Schema

// schemaJSON overrides the embedded "type" with a shallower one, and writes "$ref" after the keywords beside it; Title and Description only keep Schema's field order.
type schemaJSON struct {
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	Type        schemaType `json:"type,omitzero"`

	*schemaFields

	Ref string `json:"$ref,omitempty"`
}

// schemaType is a schema's "type": a single type, or an array of one type and optionally "null".
type schemaType struct {
	Type     DataType
	Nullable bool
}

var errMultipleTypes = errors.New(`multiple types other than "null" are not supported`)

func (t *schemaType) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	if dec.PeekKind() != jsontext.KindBeginArray {
		return json.UnmarshalDecode(dec, &t.Type)
	}

	var types []DataType
	if err := json.UnmarshalDecode(dec, &types); err != nil {
		return err
	}

	*t = schemaType{}

	for _, tp := range types {
		switch {
		case tp == TypeNull:
			t.Nullable = true
		case t.Type != "":
			return errMultipleTypes
		default:
			t.Type = tp
		}
	}

	// ["null"] on its own is just the null type.
	if t.Type == "" && t.Nullable {
		*t = schemaType{Type: TypeNull}
	}

	return nil
}

func (t *schemaType) MarshalJSONTo(enc *jsontext.Encoder) error {
	if !t.Nullable {
		return json.MarshalEncode(enc, t.Type)
	}

	return json.MarshalEncode(enc, []DataType{t.Type, TypeNull})
}

var (
	_ json.UnmarshalerFrom = (*Schema)(nil)
	_ json.MarshalerTo     = (*Schema)(nil)
)

// UnmarshalJSONFrom unmarshals a schema, reading a "type" of [X, "null"] as Type X with Nullable set.
func (s *Schema) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	v := schemaJSON{schemaFields: (*schemaFields)(s)}
	if err := json.UnmarshalDecode(dec, &v); err != nil {
		return err
	}

	s.Title, s.Description = v.Title, v.Description
	s.Ref = nil
	if v.Ref != "" {
		s.Ref = &SchemaRef{Identifier: v.Ref}
	}

	s.Type, s.Nullable = v.Type.Type, v.Type.Nullable

	return nil
}

// MarshalJSONTo marshals a schema, writing a nullable Type X as [X, "null"].
func (s *Schema) MarshalJSONTo(enc *jsontext.Encoder) error {
	v := &schemaJSON{
		Title:        s.Title,
		Description:  s.Description,
		Type:         schemaType{Type: s.Type, Nullable: s.Nullable},
		schemaFields: (*schemaFields)(s),
	}
	if s.Ref != nil {
		v.Ref = s.Ref.Identifier
	}

	return json.MarshalEncode(enc, v)
}
