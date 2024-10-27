package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestExample_JSON(t *testing.T) {
	t.Parallel()

	// In a request body:
	testJSON(t, []byte(`{
    "content": {
      "application/json": {
        "schema": {
          "$ref": "#/components/schemas/Address"
        },
        "examples": {
          "foo": {
            "summary": "A foo example",
            "value": {
              "foo": "bar"
            }
          },
          "bar": {
            "summary": "A bar example",
            "value": {
              "bar": "baz"
            }
          }
        }
      },
      "application/xml": {
        "examples": {
          "xmlExample": {
            "summary": "This is an example in XML",
            "externalValue": "https://example.org/examples/address-example.xml"
          }
        }
      },
      "text/plain": {
        "examples": {
          "textExample": {
            "summary": "This is a text example",
            "externalValue": "https://foo.bar/examples/address-example.txt"
          }
        }
      }
    }
  }`), &openapi.RequestBody{})

	// In a parameter:
	testJSON(t, []byte(`[
    {
      "name": "zipCode",
      "in": "query",
      "schema": {
        "type": "string",
        "format": "zip-code"
      },
      "examples": {
        "zip-example": {
          "$ref": "#/components/examples/zip-example"
        }
      }
    }
  ]`), &openapi.ParameterList{})

	// In a response:
	testJSON(t, []byte(`{
    "200": {
      "description": "your car appointment has been booked",
      "content": {
        "application/json": {
          "schema": {
            "$ref": "#/components/schemas/SuccessResponse"
          },
          "examples": {
            "confirmation-success": {
              "$ref": "#/components/examples/confirmation-success"
            }
          }
        }
      }
    }
  }`), &openapi.OperationResponses{})
}
