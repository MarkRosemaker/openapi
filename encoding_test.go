package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
	"github.com/go-json-experiment/json/jsontext"
)

func TestEncoding_JSON(t *testing.T) {
	t.Parallel()

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
            "historyMetadata": {
              "description": "metadata in XML format",
              "type": "object",
              "properties": {}
            },
            "profileImage": {}
          }
        },
        "encoding": {
          "historyMetadata": {
            "contentType": "application/xml; charset=utf-8"
          },
          "profileImage": {
            "contentType": "image/png, image/jpeg",
            "headers": {
              "X-Rate-Limit-Limit": {
                "description": "The number of allowed requests in the current period",
                "schema": {
                  "type": "integer"
                }
              }
            }
          }
        }
      }
    }
  }`), &openapi.RequestBody{})
}

func TestEncoding_Validate_Error(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		c   openapi.Encoding
		err string
	}{
		{openapi.Encoding{
			Headers: openapi.Headers{
				"foo": {Value: &openapi.Header{}},
			},
		}, `headers["foo"]: schema or content is required`},
		{openapi.Encoding{
			Style: "not a valid style",
		}, `style ("not a valid style") is invalid, must be one of: "matrix", "label", "form", "simple", "spaceDelimited", "pipeDelimited", "deepObject"`},
		{openapi.Encoding{
			Extensions: jsontext.Value(`{"foo": "bar"}`),
		}, `foo: ` + openapi.ErrUnknownField.Error()},
	} {
		t.Run(tc.err, func(t *testing.T) {
			if err := tc.c.Validate(); err == nil || err.Error() != tc.err {
				t.Fatalf("want: %s, got: %s", tc.err, err)
			}
		})
	}
}
