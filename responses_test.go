package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestResponses_JSON(t *testing.T) {
	t.Parallel()

	// A 200 response for a successful operation and a default response for others (implying an error):
	testJSON(t, []byte(`{
  "200": {
    "description": "a pet to be returned",
    "content": {
      "application/json": {
        "schema": {
          "$ref": "#/components/schemas/Pet"
        }
      }
    }
  },
  "default": {
    "description": "Unexpected error",
    "content": {
      "application/json": {
        "schema": {
          "$ref": "#/components/schemas/ErrorModel"
        }
      }
    }
  }
}`), &openapi.OperationResponses{})
}
