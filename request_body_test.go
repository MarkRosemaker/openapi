package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestRequestBody_JSON(t *testing.T) {
	t.Parallel()

	// A request body with a referenced model definition.
	testJSON(t, []byte(`{
  "description": "user to add to the system",
  "content": {
    "application/json": {
      "schema": {
        "$ref": "#/components/schemas/User"
      },
      "examples": {
          "user" : {
            "summary": "User Example",
            "externalValue": "https://foo.bar/examples/user-example.json"
          }
        }
    },
    "application/xml": {
      "schema": {
        "$ref": "#/components/schemas/User"
      },
      "examples": {
          "user" : {
            "summary": "User example in XML",
            "externalValue": "https://foo.bar/examples/user-example.xml"
          }
        }
    },
    "text/plain": {
      "examples": {
        "user" : {
            "summary": "User example in Plain text",
            "externalValue": "https://foo.bar/examples/user-example.txt"
        }
      }
    },
    "*/*": {
      "examples": {
        "user" : {
            "summary": "User example in other format",
            "externalValue": "https://foo.bar/examples/user-example.whatever"
        }
      }
    }
  }
}`), &openapi.RequestBody{})

	// A body parameter that is an array of string values:
	testJSON(t, []byte(`{
  "description": "user to add to the system",
  "required": true,
  "content": {
    "text/plain": {
      "schema": {
        "type": "array",
        "items": {
          "type": "string"
        }
      }
    }
  }
}`), &openapi.RequestBody{})

	// A `requestBody` for submitting a file in a `POST` operation may look like the following example:
	testJSON(t, []byte(`{
    "content": {
      "application/octet-stream": {}
    }
  }`), &openapi.RequestBody{})

	// In addition, specific media types MAY be specified:
	testJSON(t, []byte(`{
    "content": {
      "image/jpeg": {},
      "image/png": {}
    }
  }`), &openapi.RequestBody{})

	// To upload multiple files, a `multipart` media type MUST be used:
	testJSON(t, []byte(`{
    "content": {
      "multipart/form-data": {
        "schema": {
          "type": "object",
          "properties": {
            "file": {
              "type": "array",
              "items": {}
            }
          }
        }
      }
    }
  }`), &openapi.RequestBody{})

	// To submit content using form url encoding via RFC1866, the following definition may be used:
	testJSON(t, []byte(`{
    "content": {
      "application/x-www-form-urlencoded": {
        "schema": {
          "type": "object",
          "properties": {
            "id": {
              "type": "string",
              "format": "uuid"
            },
            "address": {
              "type": "object",
              "properties": {}
            }
          }
        }
      }
    }
  }`), &openapi.RequestBody{})

	// `contentMediaType` and `contentEncoding`
	testJSON(t, []byte(`{
    "content": {
      "multipart/form-data": {
        "schema": {
          "type": "object",
          "properties": {
            "id": {
              "type": "string",
              "format": "uuid"
            },
            "address": {
              "type": "object",
              "properties": {}
            },
            "profileImage": {
              "type": "string",
              "contentMediaType": "image/png",
              "contentEncoding": "base64"
            },
            "children": {
              "type": "array",
              "items": {
                "type": "string"
              }
            },
            "addresses": {
              "type": "array",
              "items": {
                "$ref": "#/components/schemas/Address"
              }
            }
          }
        }
      }
    }
  }`), &openapi.RequestBody{}) // TODO in the spec, it has  type: object in "items" besides the "$ref"
}

func TestRequestBody_Validate_Error(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		p   openapi.RequestBody
		err string
	}{
		{openapi.RequestBody{}, "content is required"},
		{openapi.RequestBody{
			Content: openapi.Content{"foo; bar": &openapi.MediaType{}},
		}, `content["foo; bar"]: mime: invalid media parameter`},
		{openapi.RequestBody{
			Content:    openapi.Content{"application/json": &openapi.MediaType{}},
			Extensions: []byte(`{"foo": "bar"}`),
		}, "foo: " + openapi.ErrUnknownField.Error()},
	} {
		t.Run(tc.err, func(t *testing.T) {
			if err := tc.p.Validate(); err == nil || err.Error() != tc.err {
				t.Fatalf("expected %q, got %q", tc.err, err)
			}
		})
	}
}
