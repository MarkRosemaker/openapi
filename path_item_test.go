package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestPathItem_JSON(t *testing.T) {
	t.Parallel()

	testJSON(t, []byte(`{
  "parameters": [
    {
      "name": "id",
      "in": "path",
      "description": "ID of pet to use",
      "required": true,
      "schema": {
        "type": "array",
        "items": {
          "type": "string"
        }
      },
      "style": "simple"
    }
  ],
  "get": {
    "summary": "Find pets by ID",
    "description": "Returns pets based on ID",
    "operationId": "getPetsById",
    "responses": {
      "200": {
        "description": "pet response",
        "content": {
          "*/*": {
            "schema": {
              "type": "array",
              "items": {
                "$ref": "#/components/schemas/Pet"
              }
            }
          }
        }
      },
      "default": {
        "description": "error payload",
        "content": {
          "text/html": {
            "schema": {
              "$ref": "#/components/schemas/ErrorModel"
            }
          }
        }
      }
    }
  }
}`), &openapi.PathItem{})

	testJSON(t, []byte(`{
  "parameters": [
    {
      "name": "id",
      "in": "path",
      "description": "ID of pet to use",
      "required": true,
      "schema": {
        "type": "array",
        "items": {
          "type": "string"
        }
      },
      "style": "simple"
    }
  ],
  "get": {
    "summary": "Find pets by ID",
    "description": "Returns pets based on ID",
    "operationId": "getPetsById",
    "responses": {
      "200": {
        "description": "pet response",
        "content": {
          "*/*": {
            "schema": {
              "type": "array",
              "items": {
                "$ref": "#/components/schemas/Pet"
              }
            }
          }
        }
      },
      "default": {
        "description": "error payload",
        "content": {
          "text/html": {
            "schema": {
              "$ref": "#/components/schemas/ErrorModel"
            }
          }
        }
      }
    }
  },
  "x-foo": "bar",
  "x-bar": 42
}`), &openapi.PathItem{})
}

func TestPathItem_Validate_Error(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		p   openapi.PathItem
		err string
	}{
		{openapi.PathItem{
			Get: &openapi.Operation{Servers: openapi.Servers{{}}},
		}, `GET.servers[0].url is required`},
		{openapi.PathItem{
			Put: &openapi.Operation{Servers: openapi.Servers{{}}},
		}, `PUT.servers[0].url is required`},
		{openapi.PathItem{
			Post: &openapi.Operation{Servers: openapi.Servers{{}}},
		}, `POST.servers[0].url is required`},
		{openapi.PathItem{
			Delete: &openapi.Operation{Servers: openapi.Servers{{}}},
		}, `DELETE.servers[0].url is required`},
		{openapi.PathItem{
			Options: &openapi.Operation{Servers: openapi.Servers{{}}},
		}, `OPTIONS.servers[0].url is required`},
		{openapi.PathItem{
			Head: &openapi.Operation{Servers: openapi.Servers{{}}},
		}, `HEAD.servers[0].url is required`},
		{openapi.PathItem{
			Patch: &openapi.Operation{Servers: openapi.Servers{{}}},
		}, `PATCH.servers[0].url is required`},
		{openapi.PathItem{
			Trace: &openapi.Operation{Servers: openapi.Servers{{}}},
		}, `TRACE.servers[0].url is required`},
	} {
		t.Run(tc.err, func(t *testing.T) {
			if err := tc.p.Validate(); err == nil || err.Error() != tc.err {
				t.Errorf("want: %s, got: %s", tc.err, err)
			}
		})
	}
}
