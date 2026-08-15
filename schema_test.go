package openapi_test

import (
	"encoding/json/v2"
	"fmt"
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestSchema_JSON(t *testing.T) {
	t.Parallel()

	testJSON(t, []byte(`{
		"type": "object",
		"example": null
	}`), &openapi.Schema{})
}

func TestSchema_Validate(t *testing.T) {
	t.Parallel()

	str := &openapi.SchemaRef{Value: &openapi.Schema{Type: openapi.TypeString}}
	num := &openapi.SchemaRef{Value: &openapi.Schema{Type: openapi.TypeNumber}}

	for i, tc := range []openapi.Schema{
		{Type: openapi.TypeNumber, Default: 3.14},
		{Type: openapi.TypeInteger, Default: 3.0},
		{Type: openapi.TypeInteger, Format: openapi.FormatDuration, Default: 3}, // e.g. seconds
		{Type: openapi.TypeString, Format: openapi.FormatByte},                  // base64-encoded data
		// oneOf, anyOf, not allow type to be omitted
		// See: https://spec.openapis.org/oas/v3.2.0.html#schema-object
		{OneOf: openapi.SchemaRefList{str, num}},
		{AnyOf: openapi.SchemaRefList{str, num}},
		{Not: str},
		// combining with a type is also valid
		{Type: openapi.TypeString, OneOf: openapi.SchemaRefList{str}},
		// enum accepts any JSON type per JSON Schema 2020-12
		{Type: openapi.TypeInteger, Enum: []any{float64(4), float64(6), float64(8)}},
		{Type: openapi.TypeString, Enum: []any{"foo", "bar"}},
	} {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			if err := tc.Validate(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestSchema_Validate_Error(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		s   openapi.Schema
		err string
	}{
		{openapi.Schema{}, "type is required"},
		{openapi.Schema{
			Type: "foo",
		}, `type ("foo") is invalid, must be one of: "integer", "number", "string", "array", "boolean", "object", "null"`},
		{openapi.Schema{
			Type: openapi.TypeArray,
		}, `items is required`},
		{openapi.Schema{
			Type:   openapi.TypeString,
			Format: "foo",
		}, `format ("foo") is invalid, must be one of: ` + validFormats},
		{openapi.Schema{
			Type:   openapi.TypeString,
			Format: openapi.FormatInt64,
		}, `format ("int64") is invalid: only valid for integer type, got string`},
		{openapi.Schema{
			Type:   openapi.TypeString,
			Format: openapi.FormatDouble,
		}, `format ("double") is invalid: only valid for number type, got string`},
		{openapi.Schema{
			Type:   openapi.TypeBoolean,
			Format: openapi.FormatByte,
		}, `format ("byte") is invalid: only valid for string type, got boolean`},
		{openapi.Schema{
			Type:   openapi.TypeBoolean,
			Format: openapi.FormatPassword,
		}, `format ("password") is invalid: only valid for string type, got boolean`},
		{openapi.Schema{
			Type:   openapi.TypeBoolean,
			Format: openapi.FormatDuration,
		}, `format ("duration") is invalid: only valid for integer or string type, got boolean`},
		{openapi.Schema{
			Type:  openapi.TypeBoolean,
			Items: &openapi.SchemaRef{},
		}, `items is invalid: only valid for array type, got boolean`},
		{openapi.Schema{
			Type: openapi.TypeArray,
			Items: &openapi.SchemaRef{
				Value: &openapi.Schema{
					Type: openapi.TypeNumber,
					Min:  new(4.0),
					Max:  new(3.0),
				},
			},
		}, `items.minimum (4) is invalid: minimum is greater than maximum (4 > 3)`},
		{openapi.Schema{
			Type: openapi.TypeBoolean,
			Min:  new(3.0),
		}, `minimum (3) is invalid: only valid for number type, got boolean`},
		{openapi.Schema{
			Type: openapi.TypeBoolean,
			Max:  new(4.0),
		}, `maximum (4) is invalid: only valid for number type, got boolean`},
		{openapi.Schema{
			Type: openapi.TypeInteger,
			Min:  new(5.3),
		}, `minimum (5.3) is invalid: not an integer`},
		{openapi.Schema{
			Type: openapi.TypeInteger,
			Max:  new(4.2),
		}, `maximum (4.2) is invalid: not an integer`},
		{openapi.Schema{
			Type: openapi.TypeInteger,
			Min:  new(5.0),
			Max:  new(4.0),
		}, `minimum (5) is invalid: minimum is greater than maximum (5 > 4)`},
		{openapi.Schema{
			Type: openapi.TypeNumber,
			Min:  new(5.6),
			Max:  new(4.2),
		}, `minimum (5.6) is invalid: minimum is greater than maximum (5.6 > 4.2)`},
		{openapi.Schema{
			Type:     openapi.TypeNumber,
			MinItems: 3,
		}, `minItems (3) is invalid: only valid for array type, got number`},
		{openapi.Schema{
			Type:     openapi.TypeNumber,
			MaxItems: new(uint(4)),
		}, `maxItems (4) is invalid: only valid for array type, got number`},
		{openapi.Schema{
			Type:     openapi.TypeArray,
			MinItems: 5,
			MaxItems: new(uint(4)),
			Items:    &openapi.SchemaRef{},
		}, `minItems (5) is invalid: minItems is greater than maxItems (5 > 4)`},
		{openapi.Schema{
			AllOf: openapi.SchemaRefList{
				{Value: &openapi.Schema{}},
			},
		}, `allOf[0].type is required`},
		{openapi.Schema{
			OneOf: openapi.SchemaRefList{
				{Value: &openapi.Schema{}},
			},
		}, `oneOf[0].type is required`},
		{openapi.Schema{
			AnyOf: openapi.SchemaRefList{
				{Value: &openapi.Schema{}},
			},
		}, `anyOf[0].type is required`},
		{openapi.Schema{
			Not: &openapi.SchemaRef{Value: &openapi.Schema{}},
		}, `not.type is required`},
		{openapi.Schema{
			Type: openapi.TypeObject,
			Properties: openapi.SchemaRefs{
				"foo": &openapi.SchemaRef{Value: &openapi.Schema{}},
			},
		}, `properties["foo"].type is required`},
		{openapi.Schema{
			Type:     openapi.TypeObject,
			Required: []string{"foo"},
		}, `required[0] ("foo") is invalid: property does not exist`},
		{openapi.Schema{
			Type: openapi.TypeObject,
			AdditionalProperties: &openapi.SchemaRef{
				Value: &openapi.Schema{},
			},
		}, `additionalProperties.type is required`},
		{openapi.Schema{
			Type:       openapi.TypeBoolean,
			Properties: openapi.SchemaRefs{},
		}, `properties is invalid: only valid for object type, got boolean`},
		{openapi.Schema{
			Type: openapi.TypeBoolean,
			AdditionalProperties: &openapi.SchemaRef{
				Value: &openapi.Schema{},
			},
		}, `additionalProperties is invalid: only valid for object type, got boolean`},
		{openapi.Schema{
			Type: openapi.TypeBoolean,
			Enum: []any{"not-a-bool"},
		}, `enum[0] ("not-a-bool") is invalid: must be a boolean value`},
		{openapi.Schema{
			Type: openapi.TypeInteger,
			Enum: []any{float64(3.14)},
		}, `enum[0] (3.14) is invalid: must be a integer value`},
		{openapi.Schema{
			Type:    openapi.TypeBoolean,
			Default: "foo",
		}, `default ("foo") is invalid: does not match schema type, got boolean`},
		{openapi.Schema{
			Type:    openapi.TypeString,
			Default: "foo",
			Enum:    []any{"bar", "buz"},
		}, `default ("foo") is invalid: is not one of the enums (["bar" "buz"])`},
		{openapi.Schema{
			Type:    openapi.TypeInteger,
			Default: 3.14,
		}, `default (3.14) is invalid: does not match schema type, got integer`},
		{openapi.Schema{
			Type:    openapi.TypeString,
			Default: 3.14,
		}, `default (3.14) is invalid: does not match schema type, got string`},
		{openapi.Schema{
			Type:    openapi.TypeString,
			Default: 3,
		}, `default (3) is invalid: does not match schema type, got string`},
		{openapi.Schema{
			Type:    openapi.TypeString,
			Default: struct{}{},
		}, `default is invalid: unknown type struct {}`},
	} {
		t.Run(tc.err, func(t *testing.T) {
			if err := tc.s.Validate(); err == nil || err.Error() != tc.err {
				t.Fatalf("want: %s, got: %s", tc.err, err)
			}
		})
	}
}

func TestSchema_UnmarshalNumericEnum(t *testing.T) {
	const src = `{
		"type": "integer",
		"enum": [4, 6, 8, 10, 12, 16]
	}`

	s := &openapi.Schema{}
	if err := json.Unmarshal([]byte(src), s); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if s.Type != openapi.TypeInteger {
		t.Errorf("Type = %q, want integer", s.Type)
	}

	want := []any{float64(4), float64(6), float64(8), float64(10), float64(12), float64(16)}
	// (json numbers become float64 by default)

	if len(s.Enum) != len(want) {
		t.Fatalf("len(Enum) = %d, want %d", len(s.Enum), len(want))
	}

	for i, v := range s.Enum {
		if v != want[i] {
			t.Errorf("Enum[%d] = %v (%T), want %v", i, v, v, want[i])
		}
	}
}
