package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

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
				Value: &openapi.Schema{},
			},
		}, `items.type is required`},
	} {
		t.Run(tc.err, func(t *testing.T) {
			if err := tc.s.Validate(); err == nil || err.Error() != tc.err {
				t.Errorf("want: %s, got: %s", tc.err, err)
			}
		})
	}
}
