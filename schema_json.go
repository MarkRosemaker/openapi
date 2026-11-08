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

	s.migrateXNullable()

	return nil
}

// migrateXNullable turns the x-nullable extension of Swagger 2.0 tools, which OpenAPI 3.1 has no use for, into how
// JSON Schema 2020-12 allows null: a type of [X, "null"], or null as an alternative of its own. x-nullable false is
// the default and dropped; a value other than a boolean is left as it is.
//
// See https://spec.openapis.org/oas/v3.1.0#data-types and https://json-schema.org/draft/2020-12/json-schema-validation#name-type.
func (s *Schema) migrateXNullable() {
	val, rest, ok := takeExtension(s.Extensions, "x-nullable")
	if !ok {
		return
	}

	var nullable bool
	if err := json.Unmarshal(val, &nullable); err != nil {
		return
	}

	s.Extensions = rest

	if !nullable {
		return
	}

	null := &Schema{Type: TypeNull}

	switch {
	case s.Type == TypeNull:
	case s.Type != "":
		s.Nullable = true
	case s.Ref != nil:
		s.AnyOf = SchemaList{{Ref: s.Ref}, null}
		s.Ref = nil
	case len(s.OneOf) > 0:
		s.OneOf = append(s.OneOf, null)
	case len(s.AnyOf) > 0:
		s.AnyOf = append(s.AnyOf, null)
	}
	// a schema of no type and no alternatives allows null already
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
