package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
	"github.com/go-json-experiment/json/jsontext"
)

func TestHeader_JSON(t *testing.T) {
	t.Parallel()

	// A simple header of type `integer`:
	testJSON(t, []byte(`{
  "description": "The number of allowed requests in the current period",
  "schema": {
    "type": "integer"
  }
}`), &openapi.Header{})
}

func TestHeader_Validate_Error(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		header openapi.Header
		err    string
	}{
		{openapi.Header{}, "schema or content is required"},
		{
			openapi.Header{
				Schema:  &openapi.Schema{},
				Content: openapi.Content{},
			},
			"schema and content are mutually exclusive",
		},
		{
			openapi.Header{
				Schema: &openapi.Schema{},
			},
			"schema.type is required",
		},
		{
			openapi.Header{
				Content: openapi.Content{},
			},
			"content is invalid: must contain exactly one entry, got 0",
		},
		{
			openapi.Header{
				Content: openapi.Content{
					"application/json": {
						Schema: &openapi.SchemaRef{
							Value: &openapi.Schema{},
						},
					},
				},
			},
			`content["application/json"].schema.type is required`,
		},
		{openapi.Header{
			Content: openapi.Content{"application/json": {}},
			Style:   "foo",
		}, `style ("foo") is invalid, must be one of: "matrix", "label", "form", "simple", "spaceDelimited", "pipeDelimited", "deepObject"`},
		{openapi.Header{
			Content: openapi.Content{"application/json": {}},
			Explode: true,
		}, `explode (true) is invalid: property has no effect when schema is not present`},
		{openapi.Header{
			Schema:  &openapi.Schema{Type: openapi.TypeString},
			Explode: true,
		}, `explode (true) is invalid: property has no effect when schema type is not array or object, got "string"`},
		{openapi.Header{
			Schema:   &openapi.Schema{Type: openapi.TypeString},
			Example:  jsontext.Value("foo"),
			Examples: openapi.Examples{},
		}, `example and examples are mutually exclusive`},
		{openapi.Header{
			Schema:   &openapi.Schema{Type: openapi.TypeString},
			Examples: openapi.Examples{"foo": invalidExample},
		}, `examples["foo"]: value and externalValue are mutually exclusive`},

		{openapi.Header{
			Schema:     &openapi.Schema{Type: openapi.TypeString},
			Extensions: jsontext.Value(`{"foo": "bar"}`),
		}, `foo: ` + openapi.ErrUnknownField.Error()},
	} {
		t.Run(tc.err, func(t *testing.T) {
			if err := tc.header.Validate(); err == nil {
				t.Fatal("expected error")
			} else if want := tc.err; want != err.Error() {
				t.Fatalf("want: %s, got: %s", want, err)
			}
		})
	}
}
