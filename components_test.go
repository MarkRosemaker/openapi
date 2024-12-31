package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
	"github.com/go-json-experiment/json/jsontext"
)

func TestComponents_JSON(t *testing.T) {
	t.Parallel()

	testJSON(t, []byte(`{
  "schemas": {
    "GeneralError": {
      "type": "object",
      "properties": {
        "code": {
          "type": "integer",
          "format": "int32"
        },
        "message": {
          "type": "string"
        }
      }
    },
    "Category": {
      "type": "object",
      "properties": {
        "id": {
          "type": "integer",
          "format": "int64"
        },
        "name": {
          "type": "string"
        }
      }
    },
    "Tag": {
      "type": "object",
      "properties": {
        "id": {
          "type": "integer",
          "format": "int64"
        },
        "name": {
          "type": "string"
        }
      }
    }
  },
  "responses": {
    "NotFound": {
      "description": "Entity not found."
    },
    "IllegalInput": {
      "description": "Illegal input for operation."
    },
    "GeneralError": {
      "description": "General Error",
      "content": {
        "application/json": {
          "schema": {
            "$ref": "#/components/schemas/GeneralError"
          }
        }
      }
    }
  },
  "parameters": {
    "skipParam": {
      "name": "skip",
      "in": "query",
      "description": "number of items to skip",
      "required": true,
      "schema": {
        "type": "integer",
        "format": "int32"
      }
    },
    "limitParam": {
      "name": "limit",
      "in": "query",
      "description": "max records to return",
      "required": true,
      "schema" : {
        "type": "integer",
        "format": "int32"
      }
    }
  },
  "securitySchemes": {
    "api_key": {
      "type": "apiKey",
      "name": "api_key",
      "in": "header"
    },
    "petstore_auth": {
      "type": "oauth2",
      "flows": {
        "implicit": {
          "authorizationUrl": "https://example.org/api/oauth/dialog",
          "scopes": {
            "write:pets": "modify pets in your account",
            "read:pets": "read your pets"
          }
        }
      }
    }
  }
}`), &openapi.Components{})
}

func TestComponents_Validate_Error(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		c   openapi.Components
		err string
	}{
		{openapi.Components{
			Schemas: openapi.Schemas{"Pet": &openapi.Schema{}},
		}, `schemas["Pet"].type is required`},
		{openapi.Components{
			Schemas: openapi.Schemas{" ": &openapi.Schema{}},
		}, `schemas[" "] (" ") is invalid: must match the regular expression "^[a-zA-Z0-9\\.\\-_]+$"`},
		{openapi.Components{
			Responses: openapi.ResponsesByName{"PetResponse": &openapi.ResponseRef{
				Value: &openapi.Response{},
			}},
		}, `responses["PetResponse"].description is required`},
		{openapi.Components{
			Responses: openapi.ResponsesByName{" ": &openapi.ResponseRef{}},
		}, `responses[" "] (" ") is invalid: must match the regular expression "^[a-zA-Z0-9\\.\\-_]+$"`},
		{openapi.Components{
			Parameters: openapi.Parameters{"MyParameter": &openapi.ParameterRef{
				Value: &openapi.Parameter{},
			}},
		}, `parameters["MyParameter"].name is required`},
		{openapi.Components{
			Examples: openapi.Examples{"MyExample": invalidExample},
		}, `examples["MyExample"]: value and externalValue are mutually exclusive`},
		{openapi.Components{
			RequestBodies: openapi.RequestBodies{"MyRequestBody": &openapi.RequestBodyRef{
				Value: &openapi.RequestBody{},
			}},
		}, `requestBodies["MyRequestBody"].content is required`},
		{openapi.Components{
			Headers: openapi.Headers{"MyRequestBody": &openapi.HeaderRef{
				Value: &openapi.Header{},
			}},
		}, `headers["MyRequestBody"]: schema or content is required`},
		{openapi.Components{
			SecuritySchemes: openapi.SecuritySchemes{"MyRequestBody": &openapi.SecuritySchemeRef{
				Value: &openapi.SecurityScheme{},
			}},
		}, `securitySchemes["MyRequestBody"].type is invalid, must be one of: "apiKey", "http", "mutualTLS", "oauth2", "openIdConnect"`},
		{openapi.Components{
			Links: openapi.Links{"MyLink": &openapi.LinkRef{
				Value: &openapi.Link{},
			}},
		}, `links.MyLink: operationRef or operationId must be set`},
		{openapi.Components{
			Callbacks: openapi.CallbackRefs{"MyCallback": &openapi.CallbackRef{
				Value: invalidCallback,
			}},
		}, `callbacks["MyCallback"]["{$request.query.callbackUrl}/data"].parameters[0].name is required`},
		{openapi.Components{
			PathItems: openapi.PathItems{" ": &openapi.PathItemRef{}},
		}, `pathItems[" "] (" ") is invalid: must match the regular expression "^[a-zA-Z0-9\\.\\-_]+$"`},
		{openapi.Components{
			PathItems: openapi.PathItems{"MyPathItem": &openapi.PathItemRef{
				Value: &openapi.PathItem{
					Parameters: openapi.ParameterList{{
						Value: &openapi.Parameter{},
					}},
				},
			}},
		}, `pathItems["MyPathItem"].parameters[0].name is required`},
		{openapi.Components{
			Extensions: jsontext.Value(`{"foo": "bar"}`),
		}, `foo: unknown field or extension without "x-" prefix`},
	} {
		t.Run(tc.err, func(t *testing.T) {
			if err := tc.c.Validate(); err == nil || err.Error() != tc.err {
				t.Fatalf("want: %s, got: %s", tc.err, err)
			}
		})
	}
}
