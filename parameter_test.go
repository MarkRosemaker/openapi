package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestParameter_JSON(t *testing.T) {
	t.Parallel()

	// A header parameter with an array of 64 bit integer numbers:
	testJSON(t, []byte(`{
  "name": "token",
  "in": "header",
  "description": "token to be passed as a header",
  "required": true,
  "style": "simple",
  "schema": {
    "type": "array",
    "items": {
      "type": "integer",
      "format": "int64"
    }
  }
}`), &openapi.Parameter{})

	// A path parameter of a string value:
	testJSON(t, []byte(`{
  "name": "username",
  "in": "path",
  "description": "username to fetch",
  "required": true,
  "schema": {
    "type": "string"
  }
}`), &openapi.Parameter{})

	// An optional query parameter of a string value, allowing multiple values by repeating the query parameter:
	testJSON(t, []byte(`{
  "name": "id",
  "in": "query",
  "description": "ID of the object to fetch",
  "style": "form",
  "explode": true,
  "schema": {
    "type": "array",
    "items": {
      "type": "string"
    }
  }
}`), &openapi.Parameter{})

	// A free-form query parameter, allowing undefined parameters of a specific type:
	testJSON(t, []byte(`{
  "name": "freeForm",
  "in": "query",
  "style": "form",
  "schema": {
    "type": "object",
    "additionalProperties": {
      "type": "integer"
    }
  }
}`), &openapi.Parameter{})

	// A complex parameter using `content` to define serialization:
	testJSON(t, []byte(`{
  "name": "coordinates",
  "in": "query",
  "content": {
    "application/json": {
      "schema": {
        "type": "object",
        "properties": {
          "lat": {
            "type": "number"
          },
          "long": {
            "type": "number"
          }
        },
        "required": [
          "lat",
          "long"
        ]
      }
    }
  }
}`), &openapi.Parameter{})
}

func TestParameter_Validate_Error(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		p   openapi.Parameter
		err string
	}{
		{openapi.Parameter{}, "name is required"},
		{openapi.Parameter{
			Name: "myname",
		}, "in is required"},
		{openapi.Parameter{
			Name: "myname",
			In:   "foo",
		}, `in ("foo") is invalid, must be one of: "path", "query", "header", "cookie"`},
		{openapi.Parameter{
			Name: "myname",
			In:   openapi.ParameterLocationPath,
		}, "required (false) is invalid: must be true for path parameters"},
		{openapi.Parameter{
			Name:     "myname",
			In:       openapi.ParameterLocationPath,
			Required: true,
		}, "schema or content is required"},
		{openapi.Parameter{
			Name:     "myname",
			In:       openapi.ParameterLocationPath,
			Required: true,
			Schema:   &openapi.Schema{},
		}, "schema.type is required"},
		{openapi.Parameter{
			Name:            "myname",
			In:              openapi.ParameterLocationPath,
			Required:        true,
			Schema:          &openapi.Schema{Type: openapi.TypeString},
			AllowEmptyValue: true,
		}, `allowEmptyValue (true) is invalid: can only be true for query parameters, got "path"`},
		{openapi.Parameter{
			Name:          "myname",
			In:            openapi.ParameterLocationPath,
			Required:      true,
			Schema:        &openapi.Schema{Type: openapi.TypeString},
			AllowReserved: true,
		}, `allowReserved (true) is invalid: only applies to query parameters, got "path"`},
		{openapi.Parameter{
			Name:     "myname",
			In:       openapi.ParameterLocationPath,
			Required: true,
			Schema:   &openapi.Schema{Type: openapi.TypeString},
			Content:  openapi.Content{},
		}, `schema and content are mutually exclusive`},
		{openapi.Parameter{
			Name:     "myname",
			In:       openapi.ParameterLocationPath,
			Required: true,
			Content:  openapi.Content{},
		}, `content is invalid: must contain exactly one entry, got 0`},
		{openapi.Parameter{
			Name:     "myname",
			In:       openapi.ParameterLocationPath,
			Required: true,
			Content:  openapi.Content{"foo/bar; baz": {}},
		}, `content["foo/bar; baz"]: mime: invalid media parameter`},
		{openapi.Parameter{
			Name:     "myname",
			In:       openapi.ParameterLocationPath,
			Required: true,
			Content:  openapi.Content{"foo": {}},
			Style:    "foo",
		}, `style ("foo") is invalid, must be one of: "matrix", "label", "form", "simple", "spaceDelimited", "pipeDelimited", "deepObject"`},
		{openapi.Parameter{
			Name:     "myname",
			In:       openapi.ParameterLocationPath,
			Required: true,
			Content:  openapi.Content{"foo": {}},
			Explode:  true,
		}, `explode (true) is invalid: property has no effect when schema is not present`},
		{openapi.Parameter{
			Name:     "myname",
			In:       openapi.ParameterLocationPath,
			Required: true,
			Schema:   &openapi.Schema{Type: openapi.TypeString},
			Explode:  true,
		}, `explode (true) is invalid: property has no effect when schema type is not array or object, got "string"`},
		{openapi.Parameter{
			Name:     "myname",
			In:       openapi.ParameterLocationPath,
			Required: true,
			Schema:   &openapi.Schema{Type: openapi.TypeString},
			Example:  "foo",
			Examples: openapi.Examples{},
		}, `example and examples are mutually exclusive`},
		{openapi.Parameter{
			Name:       "myname",
			In:         openapi.ParameterLocationPath,
			Required:   true,
			Schema:     &openapi.Schema{Type: openapi.TypeString},
			Extensions: []byte(`{"foo": "bar"}`),
		}, `foo: ` + openapi.ErrUnknownField.Error()},
		{openapi.Parameter{
			Name:   "myname",
			In:     openapi.ParameterLocationQuery,
			Schema: &openapi.Schema{Type: openapi.TypeString},
			Examples: openapi.Examples{
				"foo": invalidExample,
			},
		}, `examples["foo"]: value and externalValue are mutually exclusive`},
	} {
		t.Run(tc.err, func(t *testing.T) {
			if err := tc.p.Validate(); err == nil || err.Error() != tc.err {
				t.Fatalf("want: %s, got: %s", tc.err, err)
			}
		})
	}
}
