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
      "style": "simple",
      "schema": {
        "type": "array",
        "items": {
          "type": "string"
        }
      }
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
      "style": "simple",
      "schema": {
        "type": "array",
        "items": {
          "type": "string"
        }
      }
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
				t.Fatalf("want: %s, got: %s", tc.err, err)
			}
		})
	}
}

func TestPathItem_SetOperation(t *testing.T) {
	p := &openapi.PathItem{}

	p.SetOperation("get", &openapi.Operation{})
	if p.Get == nil {
		t.Fatalf("GET is nil")
	}

	p.SetOperation("put", &openapi.Operation{})
	if p.Put == nil {
		t.Fatalf("PUT is nil")
	}

	p.SetOperation("post", &openapi.Operation{})
	if p.Post == nil {
		t.Fatalf("POST is nil")
	}

	p.SetOperation("delete", &openapi.Operation{})
	if p.Delete == nil {
		t.Fatalf("DELETE is nil")
	}

	p.SetOperation("options", &openapi.Operation{})
	if p.Options == nil {
		t.Fatalf("OPTIONS is nil")
	}

	p.SetOperation("head", &openapi.Operation{})
	if p.Head == nil {
		t.Fatalf("HEAD is nil")
	}

	p.SetOperation("patch", &openapi.Operation{})
	if p.Patch == nil {
		t.Fatalf("PATCH is nil")
	}

	p.SetOperation("trace", &openapi.Operation{})
	if p.Trace == nil {
		t.Fatalf("TRACE is nil")
	}
}
