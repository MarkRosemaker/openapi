package openapi_test

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"regexp"
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestSchema_UnmarshalTypeArray(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		json     string
		want     openapi.DataType
		nullable bool
	}{
		{`{"type": "string"}`, openapi.TypeString, false},
		{`{"type": ["string", "null"]}`, openapi.TypeString, true},
		{`{"type": ["null", "integer"]}`, openapi.TypeInteger, true},
		{`{"type": ["object"]}`, openapi.TypeObject, false},
		{`{"type": ["null"]}`, openapi.TypeNull, false},
	} {
		var s openapi.Schema
		if err := json.Unmarshal([]byte(tc.json), &s); err != nil {
			t.Fatalf("%s: %v", tc.json, err)
		}

		if s.Type != tc.want || s.Nullable != tc.nullable {
			t.Errorf("%s: got Type %q, Nullable %v; want %q, %v", tc.json, s.Type, s.Nullable, tc.want, tc.nullable)
		}
	}

	// two types other than "null" have no representation, so they fail loudly.
	var s openapi.Schema
	if err := json.Unmarshal([]byte(`{"type": ["string", "integer"]}`), &s); err == nil {
		t.Fatal("expected an error for multiple non-null types")
	}
}

func TestSchema_Validate(t *testing.T) {
	t.Parallel()

	str := &openapi.Schema{Type: openapi.TypeString}
	num := &openapi.Schema{Type: openapi.TypeNumber}

	for i, tc := range []openapi.Schema{
		{Type: openapi.TypeNumber, Default: jsontext.Value("3.14")},
		{Type: openapi.TypeInteger, Default: jsontext.Value("3")},
		{Type: openapi.TypeInteger, Format: openapi.FormatDuration, Default: jsontext.Value("3")}, // e.g. seconds
		{Type: openapi.TypeString, Format: openapi.FormatByte},                                    // base64-encoded data
		// type is optional (JSON Schema 2020-12): the empty schema accepts
		// any value, and enum or const alone constrain it.
		{},
		{Description: "Inference output."},
		{Enum: []jsontext.Value{jsontext.Value(`"error"`)}},
		{Type: openapi.TypeArray, MaxItems: new(uint(0))},
		{Const: jsontext.Value("401")},
		// oneOf, anyOf, not allow type to be omitted
		// See: https://spec.openapis.org/oas/v3.2.0.html#schema-object
		{OneOf: openapi.SchemaList{str, num}},
		{AnyOf: openapi.SchemaList{str, num}},
		{Not: str},
		// combining with a type is also valid
		{Type: openapi.TypeString, OneOf: openapi.SchemaList{str}},
		// enum accepts any JSON type per JSON Schema 2020-12
		{Type: openapi.TypeInteger, Enum: []jsontext.Value{jsontext.Value("4"), jsontext.Value("6"), jsontext.Value("8")}},
		{Type: openapi.TypeString, Enum: []jsontext.Value{jsontext.Value(`"foo"`), jsontext.Value(`"bar"`)}},
		{Type: openapi.TypeInteger, Const: jsontext.Value("401")},
		// a nullable schema's enum and const may hold null too
		{Type: openapi.TypeString, Nullable: true, Enum: []jsontext.Value{jsontext.Value(`"foo"`), jsontext.Value("null")}},
		{Type: openapi.TypeString, Nullable: true, Const: jsontext.Value("null")},
		// prefixItems alone satisfies array's items requirement
		{Type: openapi.TypeArray, PrefixItems: openapi.SchemaList{str, num}},
		// prefixItems together with items for elements beyond it
		{Type: openapi.TypeArray, PrefixItems: openapi.SchemaList{str, num}, Items: str},
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
		{openapi.Schema{
			Type: "foo",
		}, `type ("foo") is invalid, must be one of: "integer", "number", "string", "array", "boolean", "object", "null"`},
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
			Items: &openapi.Schema{},
		}, `items is invalid: only valid for array type, got boolean`},
		{openapi.Schema{
			Type: openapi.TypeArray,
			Items: &openapi.Schema{
				Type: openapi.TypeNumber,
				Min:  new(4.0),
				Max:  new(3.0),
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
			Items:    &openapi.Schema{},
		}, `minItems (5) is invalid: minItems is greater than maxItems (5 > 4)`},
		{openapi.Schema{
			Type: openapi.TypeBoolean,
			PrefixItems: openapi.SchemaList{
				&openapi.Schema{Type: openapi.TypeString},
			},
		}, `prefixItems is invalid: only valid for array type, got boolean`},
		{openapi.Schema{
			Type: openapi.TypeArray,
			PrefixItems: openapi.SchemaList{
				&openapi.Schema{Type: openapi.TypeString, Required: []string{"id"}},
			},
			Items: &openapi.Schema{Type: openapi.TypeBoolean},
		}, `prefixItems[0].required is invalid: only valid for object type, got string`},
		{openapi.Schema{
			AllOf: openapi.SchemaList{
				&openapi.Schema{Type: openapi.TypeString, Required: []string{"id"}},
			},
		}, `allOf[0].required is invalid: only valid for object type, got string`},
		{openapi.Schema{
			OneOf: openapi.SchemaList{
				&openapi.Schema{Type: openapi.TypeString, Required: []string{"id"}},
			},
		}, `oneOf[0].required is invalid: only valid for object type, got string`},
		{openapi.Schema{
			AnyOf: openapi.SchemaList{
				&openapi.Schema{Type: openapi.TypeString, Required: []string{"id"}},
			},
		}, `anyOf[0].required is invalid: only valid for object type, got string`},
		{openapi.Schema{
			Not: &openapi.Schema{Type: openapi.TypeString, Required: []string{"id"}},
		}, `not.required is invalid: only valid for object type, got string`},
		{openapi.Schema{
			Type: openapi.TypeObject,
			Properties: openapi.Schemas{
				"foo": &openapi.Schema{Type: openapi.TypeString, Required: []string{"id"}},
			},
		}, `properties["foo"].required is invalid: only valid for object type, got string`},
		{openapi.Schema{
			Type:     openapi.TypeObject,
			Required: []string{"foo"},
		}, `required[0] ("foo") is invalid: property does not exist`},
		{openapi.Schema{
			Type: openapi.TypeObject,
			AdditionalProperties: &openapi.AdditionalProperties{
				Schema: &openapi.Schema{Type: openapi.TypeString, Required: []string{"id"}},
			},
		}, `additionalProperties.required is invalid: only valid for object type, got string`},
		{openapi.Schema{
			Type:       openapi.TypeBoolean,
			Properties: openapi.Schemas{},
		}, `properties is invalid: only valid for object type, got boolean`},
		{openapi.Schema{
			Type: openapi.TypeBoolean,
			AdditionalProperties: &openapi.AdditionalProperties{
				Schema: &openapi.Schema{},
			},
		}, `additionalProperties is invalid: only valid for object type, got boolean`},
		{openapi.Schema{
			Type: openapi.TypeBoolean,
			Enum: []jsontext.Value{jsontext.Value(`"not-a-bool"`)},
		}, `enum[0] ("not-a-bool") is invalid: must be a boolean value`},
		{openapi.Schema{
			Type: openapi.TypeInteger,
			Enum: []jsontext.Value{jsontext.Value("3.14")},
		}, `enum[0] (3.14) is invalid: must be a integer value`},
		{openapi.Schema{
			Type:  openapi.TypeInteger,
			Const: jsontext.Value(`"401"`),
		}, `const ("401") is invalid: must be a integer value`},
		{openapi.Schema{
			Type:    openapi.TypeInteger,
			Pattern: regexp.MustCompile(`^\d+$`),
		}, `pattern is invalid: only valid for string type, got integer`},
		{openapi.Schema{
			Type:            openapi.TypeObject,
			ContentEncoding: "base64",
		}, `contentEncoding is invalid: only valid for string type, got object`},
		{openapi.Schema{
			Type:  openapi.TypeString,
			Const: jsontext.Value("null"),
		}, `const ("null") is invalid: must be a string value`},
		{openapi.Schema{
			Type:    openapi.TypeBoolean,
			Default: jsontext.Value(`"foo"`),
		}, `default ("foo") is invalid: does not match schema type, got boolean`},
		{openapi.Schema{
			Type:    openapi.TypeString,
			Default: jsontext.Value(`"foo"`),
			Enum:    []jsontext.Value{jsontext.Value(`"bar"`), jsontext.Value(`"buz"`)},
		}, `default ("foo") is invalid: is not one of the enums (["bar" "buz"])`},
		{openapi.Schema{
			Type:    openapi.TypeInteger,
			Default: jsontext.Value("3.14"),
		}, `default (3.14) is invalid: does not match schema type, got integer`},
		{openapi.Schema{
			Type:    openapi.TypeString,
			Default: jsontext.Value("3.14"),
		}, `default (3.14) is invalid: does not match schema type, got string`},
		{openapi.Schema{
			Type:    openapi.TypeString,
			Default: jsontext.Value("3"),
		}, `default (3) is invalid: does not match schema type, got string`},
		{openapi.Schema{
			Type:      openapi.TypeString,
			MinLength: 5,
			MaxLength: new(uint(4)),
		}, `minLength (5) is invalid: minLength is greater than maxLength (5 > 4)`},
		{openapi.Schema{
			Type:      openapi.TypeInteger,
			MaxLength: new(uint(4)),
		}, `maxLength is invalid: only valid for string type, got integer`},
		{openapi.Schema{
			Type:         openapi.TypeString,
			ExclusiveMin: new(0.0),
		}, `exclusiveMinimum (0) is invalid: only valid for number type, got string`},
		{openapi.Schema{
			Type:         openapi.TypeString,
			ExclusiveMax: new(1.0),
		}, `exclusiveMaximum (1) is invalid: only valid for number type, got string`},
		{openapi.Schema{
			Type:        openapi.TypeObject,
			UniqueItems: true,
		}, `uniqueItems (true) is invalid: only valid for array type, got object`},
		{openapi.Schema{
			Type:          openapi.TypeArray,
			Items:         &openapi.Schema{},
			MaxProperties: new(uint(2)),
		}, `maxProperties (2) is invalid: only valid for object type, got array`},
		{openapi.Schema{
			Discriminator: &openapi.Discriminator{PropertyName: "kind"},
		}, `discriminator is invalid: only valid with oneOf, anyOf or allOf, or on a component schema another extends through allOf`},
		{openapi.Schema{
			OneOf:         openapi.SchemaList{{Type: openapi.TypeString}},
			Discriminator: &openapi.Discriminator{},
		}, `discriminator.propertyName is required`},
		{openapi.Schema{
			Ref: &openapi.SchemaRef{Identifier: "#/components/schemas/Pet"},
		}, `$ref: "#/components/schemas/Pet" was not resolved`},
		{openapi.Schema{
			Extensions: jsontext.Value(`{"minContains":1}`),
		}, `minContains: unknown field or extension without "x-" prefix`},
		{openapi.Schema{
			Type:  openapi.TypeArray,
			Items: &openapi.Schema{Extensions: jsontext.Value(`{"minContains":1}`)},
		}, `items.minContains: unknown field or extension without "x-" prefix`},
		{openapi.Schema{
			Type:         openapi.TypeInteger,
			ExclusiveMin: new(0.5),
		}, `exclusiveMinimum (0.5) is invalid: not an integer`},
		{openapi.Schema{
			Type:         openapi.TypeInteger,
			ExclusiveMax: new(1.5),
		}, `exclusiveMaximum (1.5) is invalid: not an integer`},
		{openapi.Schema{
			Type:         openapi.TypeNumber,
			Min:          new(1.0),
			ExclusiveMax: new(1.0),
		}, `minimum (1) is invalid: minimum is not less than exclusiveMaximum (1 >= 1)`},
		{openapi.Schema{
			Type:         openapi.TypeNumber,
			ExclusiveMin: new(1.0),
			Max:          new(1.0),
		}, `exclusiveMinimum (1) is invalid: exclusiveMinimum is not less than maximum (1 >= 1)`},
		{openapi.Schema{
			Type:         openapi.TypeNumber,
			ExclusiveMin: new(2.0),
			ExclusiveMax: new(1.0),
		}, `exclusiveMinimum (2) is invalid: exclusiveMinimum is not less than exclusiveMaximum (2 >= 1)`},
		{openapi.Schema{
			Type:          openapi.TypeObject,
			Properties:    openapi.Schemas{"a": {}, "b": {}},
			Required:      []string{"a", "b"},
			MaxProperties: new(uint(1)),
		}, `maxProperties (1) is invalid: fewer than the 2 required properties`},
		{openapi.Schema{
			Type:    openapi.TypeString,
			Example: jsontext.Value(`1`),
		}, `example (1) is invalid: must be a string value`},
		{openapi.Schema{
			Type:     openapi.TypeInteger,
			Examples: []jsontext.Value{jsontext.Value(`1`), jsontext.Value(`"two"`)},
		}, `examples[1] ("two") is invalid: must be a integer value`},
		{openapi.Schema{
			Type:    openapi.TypeObject,
			Example: jsontext.Value(`null`),
		}, `example ("null") is invalid: must be a object value`},
		{openapi.Schema{
			Type:       openapi.TypeNumber,
			MultipleOf: new(0.0),
		}, `multipleOf (0) is invalid: must be greater than 0`},
		{openapi.Schema{
			Type:       openapi.TypeInteger,
			MultipleOf: new(0.5),
		}, `multipleOf (0.5) is invalid: not an integer`},
		{openapi.Schema{
			Type:       openapi.TypeString,
			MultipleOf: new(4.0),
		}, `multipleOf (4) is invalid: only valid for number type, got string`},
		{openapi.Schema{
			Type:          openapi.TypeArray,
			PropertyNames: &openapi.Schema{Type: openapi.TypeString},
		}, `propertyNames is invalid: only valid for object type, got array`},
		{openapi.Schema{
			Type:          openapi.TypeObject,
			PropertyNames: &openapi.Schema{Type: openapi.TypeInteger, MaxLength: new(uint(2))},
		}, `propertyNames.maxLength is invalid: only valid for string type, got integer`},
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

	want := []string{"4", "6", "8", "10", "12", "16"}

	if len(s.Enum) != len(want) {
		t.Fatalf("len(Enum) = %d, want %d", len(s.Enum), len(want))
	}

	for i, v := range s.Enum {
		if v.Kind() != jsontext.KindNumber {
			t.Errorf("Enum[%d].Kind() = %v, want number", i, v.Kind())
		}

		if v.String() != want[i] {
			t.Errorf("Enum[%d] = %s, want %s", i, v.String(), want[i])
		}
	}
}

// TestSchema_UnmarshalPrefixItems guards against prefixItems silently
// landing in Extensions as an unrecognized field again: before this field
// existed, it parsed without error either way (Extensions has its own
// catch-all for unrecognized members), so the schema loaded fine but
// prefixItems carried no meaning -- invisible to Validate, $ref resolution,
// or anything else that reads Schema's own fields.
func TestSchema_UnmarshalPrefixItems(t *testing.T) {
	const src = `{
		"type": "array",
		"prefixItems": [
			{"type": "string"},
			{"type": "integer"}
		],
		"items": {"type": "boolean"}
	}`

	s := &openapi.Schema{}
	if err := json.Unmarshal([]byte(src), s); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if len(s.PrefixItems) != 2 {
		t.Fatalf("len(PrefixItems) = %d, want 2", len(s.PrefixItems))
	}

	if got, want := s.PrefixItems[0].Type, openapi.TypeString; got != want {
		t.Errorf("PrefixItems[0].Type = %q, want %q", got, want)
	}

	if got, want := s.PrefixItems[1].Type, openapi.TypeInteger; got != want {
		t.Errorf("PrefixItems[1].Type = %q, want %q", got, want)
	}

	if len(s.Extensions) != 0 {
		t.Errorf("Extensions = %s, want empty", s.Extensions)
	}
}

func TestSchema_Validate_Untyped(t *testing.T) {
	t.Parallel()

	maxItems := uint(3)
	multipleOf := 4.0

	// without a type, each keyword applies to instances of its own type (JSON Schema 2020-12 core, §7.6.1)
	s := &openapi.Schema{
		AnyOf:      openapi.SchemaList{{Type: openapi.TypeObject}, {Type: openapi.TypeNull}},
		Properties: openapi.Schemas{"width": {Type: openapi.TypeInteger, MultipleOf: &multipleOf}},
		Required:   []string{"width", "height"},
		MinLength:  1,
		Format:     openapi.FormatUUID,
		MultipleOf: &multipleOf,
		MaxItems:   &maxItems,
	}
	if err := s.Validate(); err != nil {
		t.Fatal(err)
	}

	// the keywords themselves are still checked
	s.Properties["width"].Type = openapi.TypeString
	if err := s.Validate(); err == nil || err.Error() != `properties["width"].multipleOf (4) is invalid: only valid for number type, got string` {
		t.Fatalf("got %v", err)
	}
}
