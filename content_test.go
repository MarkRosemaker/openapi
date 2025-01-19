package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestContent_JSON(t *testing.T) {
	t.Parallel()

	// Example
	testJSON(t, []byte(`{
  "application/json": {
    "schema": {
         "$ref": "#/components/schemas/Pet"
    },
    "examples": {
      "cat" : {
        "summary": "An example of a cat",
        "value":
          {
            "name": "Fluffy",
            "petType": "Cat",
            "color": "White",
            "gender": "male",
            "breed": "Persian"
          }
      },
      "dog": {
        "summary": "An example of a dog with a cat's name",
        "value" :  {
          "name": "Puma",
          "petType": "Dog",
          "color": "Black",
          "gender": "Female",
          "breed": "Mixed"
        }
      },
      "frog": {
        "$ref": "#/components/examples/frog-example"
      }
    }
  }
}`), &openapi.Content{})

	// NOTE: Content transferred in binary (octet-stream) MAY omit `schema`
	// a PNG image as a binary file
	testJSON(t, []byte(`{"image/png": {}}`), &openapi.Content{})
	// an arbitrary binary file
	testJSON(t, []byte(`{"application/octet-stream": {}}`), &openapi.Content{})

	// Binary content transferred with base64 encoding
	// Note that the `Content-Type` remains `image/png`, describing the semantics of the payload.  The JSON Schema `type` and `contentEncoding` fields explain that the payload is transferred as text.  The JSON Schema `contentMediaType` is technically redundant, but can be used by JSON Schema tools that may not be aware of the OpenAPI context.
	testJSON(t, []byte(`{
    "image/png": {
      "schema": {
        "type": "string",
        "contentMediaType": "image/png",
        "contentEncoding": "base64"
      }
    }
  }`), &openapi.Content{})

	// These examples apply to either input payloads of file uploads or response payloads.
}

func TestContent_Validate_Error(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		c   openapi.Content
		err string
	}{
		{openapi.Content{
			"not a real media type": &openapi.MediaType{},
		}, `["not a real media type"]: mime: expected slash after first token`},
		{openapi.Content{
			openapi.MediaRangeJSON: &openapi.MediaType{
				Schema: &openapi.SchemaRef{Value: &openapi.Schema{}},
			},
		}, `["application/json"].schema.type is required`},
	} {
		t.Run(tc.err, func(t *testing.T) {
			if err := tc.c.Validate(); err == nil || err.Error() != tc.err {
				t.Fatalf("expected %q, got %q", tc.err, err)
			}
		})
	}
}
