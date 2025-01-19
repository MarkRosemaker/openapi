package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestResponse_JSON(t *testing.T) {
	t.Parallel()

	// Response of an array of a complex type:
	testJSON(t, []byte(`{
  "description": "A complex object array response",
  "content": {
    "application/json": {
      "schema": {
        "type": "array",
        "items": {
          "$ref": "#/components/schemas/VeryComplexType"
        }
      }
    }
  }
}`), &openapi.Response{})

	// Response with a string type:
	testJSON(t, []byte(`{
  "description": "A simple string response",
  "content": {
    "text/plain": {
      "schema": {
        "type": "string"
      }
    }
  }
}`), &openapi.Response{})

	// Plain text response with headers:
	testJSON(t, []byte(`{
  "description": "A simple string response",
  "headers": {
    "X-Rate-Limit-Limit": {
      "description": "The number of allowed requests in the current period",
      "schema": {
        "type": "integer"
      }
    },
    "X-Rate-Limit-Remaining": {
      "description": "The number of remaining requests in the current period",
      "schema": {
        "type": "integer"
      }
    },
    "X-Rate-Limit-Reset": {
      "description": "The number of seconds left in the current period",
      "schema": {
        "type": "integer"
      }
    }
  },
  "content": {
    "text/plain": {
      "schema": {
        "type": "string",
        "example": "whoa!"
      }
    }
  }
}`), &openapi.Response{})

	// Response with no return value:
	testJSON(t, []byte(`{
  "description": "object created"
}`), &openapi.Response{})
}

func TestResponse_Validate_Error(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		p   openapi.Response
		err string
	}{
		{openapi.Response{}, "description is required"},
		{openapi.Response{
			Description: "some description",
			Headers:     openapi.Headers{"foo": {Value: &openapi.Header{}}},
		}, `headers["foo"]: schema or content is required`},
		{openapi.Response{
			Description: "some description",
			Content: openapi.Content{openapi.MediaRangeJSON: {
				Schema: &openapi.SchemaRef{Value: &openapi.Schema{}},
			}},
		}, `content["application/json"].schema.type is required`},
		{openapi.Response{
			Description: "some description",
			Links:       openapi.Links{"address": {Value: &openapi.Link{}}},
		}, `links.address: operationRef or operationId must be set`},
	} {
		t.Run(tc.err, func(t *testing.T) {
			if err := tc.p.Validate(); err == nil || err.Error() != tc.err {
				t.Fatalf("expected: %s, got: %s", tc.err, err)
			}
		})
	}
}
