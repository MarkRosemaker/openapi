package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func pointer[T any](v T) *T { return &v }

func TestSchema_Validate_Error(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		s   openapi.Schema
		err string
	}{
		{openapi.Schema{}, "type is required"},
		{openapi.Schema{
			Type: "foo",
		}, `type ("foo") is invalid, must be one of: "integer", "number", "string", "array", "boolean", "object"`},
		{openapi.Schema{
			Type: openapi.TypeArray,
		}, `items is required`},
		{openapi.Schema{
			Type:   openapi.TypeString,
			Format: "foo",
		}, `format ("foo") is invalid, must be one of: "int32", "int64", "float", "double", "byte", "binary", "date", "date-time", "password", "duration", "uuid", "email", "uri", "zip-code"`},
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
			Format: openapi.FormatPassword,
		}, `format ("password") is invalid: only valid for string type, got boolean`},
		{openapi.Schema{
			Type:  openapi.TypeBoolean,
			Items: &openapi.SchemaRef{},
		}, `items is invalid: only valid for array type, got boolean`},
		{openapi.Schema{
			Type: openapi.TypeArray,
			Items: &openapi.SchemaRef{
				Value: &openapi.Schema{
					Type: openapi.TypeNumber,
					Min:  pointer(4.0),
					Max:  pointer(3.0),
				},
			},
		}, `items.minimum (4) is invalid: minimum is greater than maximum (4 > 3)`},
		{openapi.Schema{
			Type: openapi.TypeBoolean,
			Min:  pointer(3.0),
		}, `minimum (3) is invalid: only valid for number type, got boolean`},
		{openapi.Schema{
			Type: openapi.TypeBoolean,
			Max:  pointer(4.0),
		}, `maximum (4) is invalid: only valid for number type, got boolean`},
		{openapi.Schema{
			Type: openapi.TypeInteger,
			Min:  pointer(5.3),
		}, `minimum (5.3) is invalid: not an integer`},
		{openapi.Schema{
			Type: openapi.TypeInteger,
			Max:  pointer(4.2),
		}, `maximum (4.2) is invalid: not an integer`},
		{openapi.Schema{
			Type: openapi.TypeInteger,
			Min:  pointer(5.0),
			Max:  pointer(4.0),
		}, `minimum (5) is invalid: minimum is greater than maximum (5 > 4)`},
		{openapi.Schema{
			Type: openapi.TypeNumber,
			Min:  pointer(5.6),
			Max:  pointer(4.2),
		}, `minimum (5.6) is invalid: minimum is greater than maximum (5.6 > 4.2)`},
		{openapi.Schema{
			Type:     openapi.TypeNumber,
			MinItems: 3,
		}, `minItems (3) is invalid: only valid for array type, got number`},
		{openapi.Schema{
			Type:     openapi.TypeNumber,
			MaxItems: pointer[uint](4),
		}, `maxItems (4) is invalid: only valid for array type, got number`},
		{openapi.Schema{
			Type:     openapi.TypeArray,
			MinItems: 5,
			MaxItems: pointer[uint](4),
			Items:    &openapi.SchemaRef{},
		}, `minItems (5) is invalid: minItems is greater than maxItems (5 > 4)`},
		{openapi.Schema{
			AllOf: openapi.SchemaRefs{
				{Value: &openapi.Schema{}},
			},
		}, `allOf[0].type is required`},
		{openapi.Schema{
			AllOf: openapi.SchemaRefs{
				{Value: &openapi.Schema{}},
			},
		}, `allOf[0].type is required`},
		{openapi.Schema{
			Type: openapi.TypeObject,
			Properties: openapi.Schemas{
				"foo": {Value: &openapi.Schema{}},
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
			Properties: openapi.Schemas{},
		}, `properties is invalid: only valid for object type, got boolean`},
		{openapi.Schema{
			Type: openapi.TypeBoolean,
			AdditionalProperties: &openapi.SchemaRef{
				Value: &openapi.Schema{},
			},
		}, `additionalProperties is invalid: only valid for object type, got boolean`},
	} {
		t.Run(tc.err, func(t *testing.T) {
			if err := tc.s.Validate(); err == nil || err.Error() != tc.err {
				t.Fatalf("want: %s, got: %s", tc.err, err)
			}
		})
	}
}
